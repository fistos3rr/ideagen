package main

import (
	"net/http"
	"encoding/json"
)

func (app *application) aiHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.serverErrorResponse(w, r, err)
	}
	
	answer, err := app.aiProvider.SendMessage(r.Context(), req.Message)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	
	err = app.writeJSON(w, http.StatusOK, envelope{"answer": answer}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}