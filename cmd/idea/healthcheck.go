package main

import (
	"net/http"
)

// healthcheckHandler godoc
//
// @Summary Check health
// @Tags helper
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (app *application) healthcheckHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	env := envelope{
		"status": "available",
		"env":    app.config.env,
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
