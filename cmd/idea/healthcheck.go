package main

import (
	"net/http"

	"github.com/fistos3rr/ideagen/internal/api/dto"
)

// healthcheckHandler godoc
//
// @Summary Check health
// @Tags metrics
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /health [get]
func (app *application) healthcheckHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	resp := dto.HealthResponse{
		Status: "available",
		Environment: app.config.env,
	}

	err := app.writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
