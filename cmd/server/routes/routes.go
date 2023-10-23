package routes

import (
	"github.com/denelop/potok/pkg/streaming"
	"github.com/gorilla/mux"
)

func API(router *mux.Router) {
	router.HandleFunc(
		"/streaming/{streamName}/{file}",
		streaming.Handler,
	).Methods("GET")
}
