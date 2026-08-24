package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()
	//test health
	router.HandleFunc("/v1/health", app.healthcheckHandler).Methods("GET")

	return router
}
