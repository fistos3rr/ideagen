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

	router.HandleFunc("/v1/prompts", app.createPromptHandler).Methods("POST")
	router.HandleFunc("/v1/prompts/{id}", app.showPromptHandler).Methods("GET")
	router.HandleFunc("/v1/prompts/{id}", app.deletePromptHandler).Methods("DELETE")
	router.HandleFunc("/v1/prompts", app.listPromptsHandler).Methods("GET")
	router.HandleFunc("/v1/prompts/{id}", app.updatePromptHandler).Methods("UPDATE")

	router.HandleFunc("/v1/random/types", app.randomTypesHandler).Methods("GET")
	router.HandleFunc("/v1/random/prompts", app.randomPromptsHandler).Methods("GET")

	return router
}
