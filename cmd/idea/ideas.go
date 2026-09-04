package main

import (
	"context"
	"errors"
	"fmt"
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

func (app *application) showIdeaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	idea, err := app.models.Ideas.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"idea": idea}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createIdeaHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TypeID int64  `json:"type_id"`
		Text   string `json:"text"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	typeObj, err := app.models.Types.Get(input.TypeID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			v := validator.New()
			v.AddError("type_id", "type with this id does not exist")
			app.failedValidationResponse(w, r, v.Errors)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	idea := &data.Idea{
		Type: typeObj,
		Text: input.Text,
	}

	v := validator.New()
	if data.ValidateIdea(v, idea); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Ideas.Insert(idea)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/ideas/%d", idea.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"idea": idea}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteIdeaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Ideas.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "idea successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listIdeasHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Text       string
		TypeID     int64
		ActiveOnly bool
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Text = app.readString(qs, "text", "")
	input.TypeID = int64(app.readInt(qs, "type_id", 0, v))
	input.ActiveOnly = app.readBool(qs, "active_only", v)
	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Sort = app.readString(qs, "sort", "id")
	input.SortSafelist = []string{"id", "name", "type_id", "-id", "-name", "-type_id"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	ideas, metadata, err := app.models.Ideas.GetAll(input.Text, input.TypeID, input.ActiveOnly, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"ideas": ideas, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateIdeaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	idea, err := app.models.Ideas.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Text   *string `json:"name"`
		TypeID *int64  `json:"type_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Text != nil {
		idea.Text = *input.Text
	}
	if input.TypeID != nil {
		typeObj, err := app.models.Types.Get(*input.TypeID)
		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				v := validator.New()
				v.AddError("type_id", "type with this id does not exist")
				app.failedValidationResponse(w, r, v.Errors)
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}
		idea.Type = typeObj
	}

	v := validator.New()
	if data.ValidateIdea(v, idea); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Ideas.Update(idea)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"idea": idea}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
