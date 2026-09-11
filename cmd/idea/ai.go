package main

import (
	"encoding/json"
	"net/http"

	"github.com/fistos3rr/ideagen/internal/api/dto"
)


// aiHandler godoc
//
// @Summary Ask AI
// @Description Ask AI using AI provider
// @Tags ai
// @Accept json
// @Produce json
// @Param request body dto.AskRequest true "Message to AI"
// @Success 200 {object} dto.AskResponse "Success answer"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Not enough privilegies"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /ask [post]
func (app *application) aiHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.AskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	answer, err := app.aiProvider.SendMessage(r.Context(), req.Message)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	resp := dto.AskResponse{
		Answer: answer,
	}
	err = app.writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
