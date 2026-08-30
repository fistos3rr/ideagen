package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/v1/health", app.healthcheckHandler).Methods("GET")
	router.HandleFunc("/v1/ask", app.aiHandler).Methods("POST")

	router.HandleFunc("/v1/types", app.createTypeHandler).Methods("POST")
	router.HandleFunc("/v1/types/{id}", app.showTypeHandler).Methods("GET")
	router.HandleFunc("/v1/types/{id}", app.deleteTypeHandler).Methods("DELETE")
	router.HandleFunc("/v1/types", app.listTypesHandler).Methods("GET")
	router.HandleFunc("/v1/types/{id}", app.updateTypeHandler).Methods("UPDATE")

	router.HandleFunc("/v1/ideas", app.createIdeaHandler).Methods("POST")
	router.HandleFunc("/v1/ideas/{id}", app.showIdeaHandler).Methods("GET")
	router.HandleFunc("/v1/ideas/{id}", app.deleteIdeaHandler).Methods("DELETE")
	router.HandleFunc("/v1/ideas", app.listIdeasHandler).Methods("GET")
	router.HandleFunc("/v1/ideas/{id}", app.updateIdeaHandler).Methods("UPDATE")

	router.HandleFunc("/v1/idea", app.generateIdeaHandler).Methods("GET")

	router.HandleFunc("/v1/random/types", app.randomTypesHandler).Methods("GET")

	return router
}
