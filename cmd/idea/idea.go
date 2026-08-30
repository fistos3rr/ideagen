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
	switch providerType {
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

func (app *application) generateIdea(ctx context.Context, t *data.Type) (*data.Idea, error) {
	sysPr, pr, err := app.promptManager.GetPrompts(t.Name)
	if err != nil {
		return nil, err
	}

	err = app.newRequest(pr, sysPr)
	if errors.Is(err, ai.ErrSetRequest) {
		app.logger.PrintFatal(err, nil)
	} else if err != nil {
		return nil, err
	}

	answer, err := app.aiProvider.SendRequest(ctx)

	idea := &data.Idea{
		Type: t,
		Text: answer,
	}

	return idea, err
}

func (app *application) generateIdeaHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TypeID     int
		ActiveOnly bool
	}

	v := validator.New()
	qs := r.URL.Query()

	input.TypeID = app.readInt(qs, "type_id", -1, v)
	input.ActiveOnly = app.readBool(qs, "active_only", v)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	var t *data.Type
	var err error
	if input.TypeID == -1 {
		types, err := app.models.Types.GetRandom(1, input.ActiveOnly)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		if len(types) == 0 {
			app.notFoundResponse(w, r)
			return
		}
		t = types[0]
	} else {
		t, err = app.models.Types.Get(int64(input.TypeID))
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
		if !t.IsActive {
			app.failedValidationResponse(w, r, map[string]string{
				"active_only": "type is not active, while active_only query set",
			})
			return
		}
	}

	idea, err := app.generateIdea(r.Context(), t)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"idea": idea.Text, "type_id": t.ID}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
