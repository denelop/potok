package routes

import (
	"net/http"

	"github.com/domonda/go-function"
	"github.com/gorilla/mux"
)

func API(router *mux.Router) {
	router.HandleFunc(
		"/stream/hls/{stream}",
		function.HTTPHandler(
			getMuxVars,
			streamingGetHLSBytes,
			function.RespondDetectContentType,
		),
	).Methods("GET")
}

func getMuxVars(request *http.Request) (map[string]string, error) {
	return mux.Vars(request), nil
}
