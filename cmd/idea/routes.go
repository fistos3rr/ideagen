package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandleFunc("/v1/health", app.healthcheckHandler).Methods("GET")
	router.HandleFunc("/v1/ask", app.requireUserRole("admin", app.aiHandler)).Methods("POST")

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

	router.HandleFunc("/v1/users", app.registerUserHandler).Methods("POST")

	router.HandleFunc("/v1/idea", app.requireAuthenticatedUser(app.generateIdeaHandler)).Methods("GET")

	router.HandleFunc("/v1/random/types", app.randomTypesHandler).Methods("GET")

	router.HandleFunc("/v1/tokens/auth", app.createJwtTokenHandler).Methods("POST")

	return app.recoverPanic(app.authenticate(router))
}
