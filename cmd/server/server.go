package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/denelop/potok/cmd/server/routes"
	"github.com/denelop/potok/pkg/env"
	"github.com/denelop/potok/pkg/streaming"
	"github.com/domonda/go-errs"
	"github.com/domonda/golog/log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/ungerik/go-fs"
	"github.com/ungerik/go-httpx"
	"github.com/ungerik/go-httpx/httperr"
	"golang.org/x/crypto/acme/autocert"
)

var config struct {
	Port           int      `env:"SERVER_PORT"`
	TLSPort        int      `env:"SERVER_TLS_PORT"`
	CertDir        fs.File  `env:"CERT_DIR"`
	Domains        []string `env:"SERVER_DOMAINS"`
	AllowedOrigins []string `env:"ALLOWED_ORIGINS"`
}

func main() {
	err := env.Parse(&config)
	if err != nil {
		panic(err)
	}

	err = streaming.StartAll(context.Background())
	if err != nil {
		panic(err)
	}

	httperr.DefaultHandler = httperr.HandlerFunc(func(err error, writer http.ResponseWriter, request *http.Request) (handled bool) {
		if errors.Is(err, context.Canceled) || errs.IsContextCanceled(request.Context()) {
			return httperr.Handled
		}
		if httperr.WriteHandler(err, writer, request) {
			// dont log handled errors
			return httperr.Handled
		}
		if httperr.ShouldLog(err) {
			log.Error("Internal error").
				Request(request).
				Err(err).
				Log()
		}
		httperr.WriteInternalServerError(err, writer)
		return httperr.Handled
	})

	log.Debug("Configurated").
		StructFields(config).
		Log()

	router := mux.NewRouter().StrictSlash(true)
	router.Use(SecHeaders)

	routes.API(router)

	var handler http.Handler = router
	if len(config.AllowedOrigins) > 0 {
		allowedMethods := []string{"GET", "POST"}
		exposedHeaders := []string{"X-Request-ID"}
		handler = handlers.CORS(
			handlers.AllowedOrigins(config.AllowedOrigins),
			handlers.AllowedMethods(allowedMethods),
			handlers.ExposedHeaders(exposedHeaders),
		)(router)

		log.Info("CORS enabled because allowed origins are configured").
			Strs("allowedOrigins", config.AllowedOrigins).
			Strs("allowedMethods", allowedMethods).
			Strs("exposedHeaders", exposedHeaders).
			Log()
	}

	serverAddr := fmt.Sprintf(":%d", config.Port)
	if config.TLSPort != 0 {
		serverAddr = fmt.Sprintf(":%d", config.TLSPort)
	}

	server := &http.Server{
		Addr:         serverAddr,
		Handler:      handler,
		ErrorLog:     log.ErrorWriter().StdLogger(),
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	httpx.GracefulShutdownServerOnSignal(server, nil, log.ErrorWriter(), time.Minute)

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		panic(err)
	}

	if config.TLSPort != 0 {
		certManager := autocert.Manager{
			Email:      "hi@denelop.com",
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(config.Domains...),
			Cache:      autocert.DirCache(config.CertDir),
		}

		go func() {
			log.Info("autocert manager listening").
				Int("port", config.Port).
				Log()

			err := http.ListenAndServe(
				fmt.Sprintf(":%d", config.Port),
				// redirect to https
				certManager.HTTPHandler(http.HandlerFunc((func(w http.ResponseWriter, req *http.Request) {
					target := "https://" + req.Host + req.URL.Path
					if len(req.URL.RawQuery) > 0 {
						target += "?" + req.URL.RawQuery + req.URL.Fragment
					}
					http.Redirect(w, req, target, http.StatusPermanentRedirect)
				}))),
			)
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Error("autocert manager server error").
					Err(err).
					Log()
			}
		}()

		server.TLSConfig = &tls.Config{GetCertificate: certManager.GetCertificate, MinVersion: tls.VersionTLS12}
		log.Info("TLS listening").
			Log()
		err = server.ServeTLS(listener, "", "")
	} else {
		log.Info("Listening").
			Log()
		err = server.Serve(listener)
	}

	if err == nil || errors.Is(err, http.ErrServerClosed) {
		log.Info("Shutting down gracefully").
			Log()
	} else {
		log.Fatal("Server error").
			Err(err).
			LogAndPanic()
	}
}

// SecHeaders injects security headers to every response.
// Read more: https://gist.github.com/mikesamuel/f7c7caed42413396e4d3e61dae6f5712.
func SecHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Specifies interactions between frames and affects click-jacking.
		w.Header().Set("X-Frame-Options", "DENY")

		// Affects leakage of sensitive information via URL parameters included in the referrer header.
		w.Header().Set("Referrer-Policy", "no-referrer")

		// Addresses leakage of sensitive information and MITM attacks via HTTPS->HTTP downgrade attacks.
		if r.TLS != nil {
			w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		}

		// Addresses polyglot attacks by forbidding content-type sniffing.
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// No search engine should index anything about the app
		w.Header().Set("X-Robots-Tag", "noindex")

		h.ServeHTTP(w, r)
	})
}
