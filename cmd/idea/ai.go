package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) aiHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Message string `json:"message"`
		Temperature *float64 `json:"temperature,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	temperature := 1.0
	if req.Temperature != nil {
		temperature = *req.Temperature
	}

	answer, err := app.aiProvider.SendMessage(r.Context(), req.Message, temperature)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"answer": answer}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
