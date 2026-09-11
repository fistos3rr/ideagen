package main

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/fistos3rr/ideagen/docs"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandleFunc("/v1/health", app.healthcheckHandler).Methods("GET")//

	router.HandleFunc("/v1/ask", app.requireUserRole("Admin", app.aiHandler)).Methods("POST")//

	router.HandleFunc("/v1/types", app.requireUserRole("Admin", app.createTypeHandler)).Methods("POST")
	router.HandleFunc("/v1/types/{id}", app.requireUserRole("Admin", app.showTypeHandler)).Methods("GET")
	router.HandleFunc("/v1/types/{id}", app.requireUserRole("Admin", app.deleteTypeHandler)).Methods("DELETE")
	router.HandleFunc("/v1/types", app.requireUserRole("Admin", app.listTypesHandler)).Methods("GET")
	router.HandleFunc("/v1/types/{id}", app.requireUserRole("Admin", app.updateTypeHandler)).Methods("PATCH")

	router.HandleFunc("/v1/ideas", app.requireUserRole("Admin", app.createIdeaHandler)).Methods("POST")
	router.HandleFunc("/v1/ideas/{id}", app.requireUserRole("Admin", app.showIdeaHandler)).Methods("GET")
	router.HandleFunc("/v1/ideas/{id}", app.requireUserRole("Admin", app.deleteIdeaHandler)).Methods("DELETE")
	router.HandleFunc("/v1/ideas", app.requireUserRole("Admin", app.listIdeasHandler)).Methods("GET")
	router.HandleFunc("/v1/ideas/{id}", app.requireUserRole("Admin", app.updateIdeaHandler)).Methods("PATCH")

	router.HandleFunc("/v1/user-idea", app.requireUserRole("Admin", app.createUserIdeaHandler)).Methods("POST")

	router.HandleFunc("/v1/service/me", app.requireAuthenticatedUser(app.showMeHandler)).Methods("GET")
	router.HandleFunc("/v1/service/ideas", app.requireAuthenticatedUser(app.listMyIdeasHandler)).Methods("GET")
	router.HandleFunc("/v1/service/ideas/{id}", app.requireAuthenticatedUser(app.showMyIdeaHandler)).Methods("GET")
	router.HandleFunc("/v1/service/ideas/{id}", app.requireAuthenticatedUser(app.deleteMyIdeaHandler)).Methods("DELETE")

	router.HandleFunc("/v1/service/idea/generate", app.requireAuthenticatedUser(app.generateMyIdeaHandler)).Methods("POST")
	router.HandleFunc("/v1/service/idea/buffer/{id}", app.requireAuthenticatedUser(app.showMyBufferIdeaHandler)).Methods("GET")
	router.HandleFunc("/v1/service/idea/buffer", app.requireAuthenticatedUser(app.listMyBufferIdeasHandler)).Methods("GET")
	router.HandleFunc("/v1/service/idea/buffer", app.requireAuthenticatedUser(app.chooseMyBufferIdeaHandler)).Methods("POST")

	router.HandleFunc("/v1/register", app.registerUserHandler).Methods("POST")
	router.HandleFunc("/v1/auth/login", app.loginUserHandler).Methods("POST")//
	router.HandleFunc("/v1/auth/logout", app.requireAuthenticatedUser(app.logoutHandler)).Methods("POST")
	router.HandleFunc("/v1/auth/refresh", app.refreshHandler).Methods("POST")

	apiHandler := app.recoverPanic(app.authenticate(router))

	swaggerMux := http.NewServeMux()
	swaggerMux.Handle("/swagger/", httpSwagger.WrapHandler)

	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/swagger/") {
			swaggerMux.ServeHTTP(w, r)
			return
		}
		apiHandler.ServeHTTP(w, r)
	})
}
