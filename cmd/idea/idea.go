package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/fistos3rr/ideagen/internal/ai"
	"github.com/fistos3rr/ideagen/internal/data"
	"github.com/fistos3rr/ideagen/internal/validator"
)

func (app *application) newRequest(pr string, sysPr string) error {
	var reqErr error
	providerType := app.config.aiProviderType
	var req ai.Request
	switch providerType{
	case "groq":
		req = ai.NewGroqRequest(app.aiConfig)
		req.AddMessage(pr)
		req.AddSystemMessage(sysPr)
		reqErr = app.aiProvider.SetRequest(req)
	default:
		return ai.ErrSetRequest
	}

	return reqErr
}

func (app *application) generateIdea(ctx context.Context, t string) (string, error) {
	sysPr, pr, err := app.promptManager.GetPrompts(t)
	if err != nil {
		return "", err
	}

	err = app.newRequest(pr, sysPr)
	if errors.Is(err, ai.ErrSetRequest) {
		app.logger.PrintFatal(err, nil)
	} else if err != nil {
		return "", err
	}

	answer, err := app.aiProvider.SendRequest(ctx)

	return answer, err
}

func (app *application) generateIdeaHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TypeID int
	}

	v := validator.New()
	qs := r.URL.Query()

	input.TypeID = app.readInt(qs, "type_id", 0, v)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	t, err := app.models.Types.Get(int64(input.TypeID))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.failedValidationResponse(w, r, map[string]string{
				"type_id": "type with this id does not exists",
			})
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	ideaText, err := app.generateIdea(r.Context(), t.Name)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"idea": ideaText}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
