package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/v1/health", app.healthcheckHandler).Methods("GET")
	router.HandleFunc("/v1/ask", app.aiHandler).Methods("POST")

	return router
}
