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

	router.HandleFunc("/v1/ask", app.requireUserRole("Admin", app.aiHandler)).Methods("POST")

	router.HandleFunc("/v1/types", app.requireUserRole("Admin", app.createTypeHandler)).Methods("POST")
	router.HandleFunc("/v1/types/{id}", app.requireUserRole("Admin", app.showTypeHandler)).Methods("GET")
	router.HandleFunc("/v1/types/{id}", app.requireUserRole("Admin", app.deleteTypeHandler)).Methods("DELETE")
	router.HandleFunc("/v1/types", app.requireUserRole("Admin", app.listTypesHandler)).Methods("GET")
	router.HandleFunc("/v1/types/{id}", app.requireUserRole("Admin", app.updateTypeHandler)).Methods("UPDATE")
	//router.HandleFunc("/v1/random/types", app.requireUserRole("Admin", app.randomTypesHandler)).Methods("GET")

	router.HandleFunc("/v1/ideas", app.requireUserRole("Admin", app.createIdeaHandler)).Methods("POST")
	router.HandleFunc("/v1/ideas/{id}", app.requireUserRole("Admin", app.showIdeaHandler)).Methods("GET")
	router.HandleFunc("/v1/ideas/{id}", app.requireUserRole("Admin", app.deleteIdeaHandler)).Methods("DELETE")
	router.HandleFunc("/v1/ideas", app.requireUserRole("Admin", app.listIdeasHandler)).Methods("GET")
	router.HandleFunc("/v1/ideas/{id}", app.requireUserRole("Admin", app.updateIdeaHandler)).Methods("UPDATE")
	//router.HandleFunc("/v1/idea", app.requireUserRole("Admin", app.generateIdeaHandler)).Methods("GET")

	router.HandleFunc("/v1/user-idea", app.requireUserRole("Admin", app.createUserIdeaHandler)).Methods("POST")

	router.HandleFunc("/v1/service/me", app.requireAuthenticatedUser(app.showMeHandler)).Methods("GET")
	router.HandleFunc("/v1/service/ideas", app.requireAuthenticatedUser(app.listMyIdeasHandler)).Methods("GET")
	router.HandleFunc("/v1/service/{id}", app.requireAuthenticatedUser(app.showMyIdeaHandler)).Methods("GET")
	router.HandleFunc("/v1/service/{id}", app.requireAuthenticatedUser(app.deleteMyIdeaHandler)).Methods("DELETE")

	router.HandleFunc("/v1/register", app.registerUserHandler).Methods("POST")
	router.HandleFunc("/v1/auth/login", app.loginUserHandler).Methods("POST")
	router.HandleFunc("/v1/auth/logout", app.requireAuthenticatedUser(app.logoutHandler)).Methods("POST")
	router.HandleFunc("/v1/auth/refresh", app.refreshHandler).Methods("POST")

	return app.recoverPanic(app.authenticate(router))
}
