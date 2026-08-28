package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fistos3rr/ideagen/internal/data"
	"github.com/fistos3rr/ideagen/internal/validator"
)

func (app *application) randomPromptsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TypeID     int64
		Limit      int
		ActiveOnly bool
	}

	v := validator.New()
	qs := r.URL.Query()

	input.TypeID = int64(app.readInt(qs, "type_id", 0, v))
	input.Limit = app.readInt(qs, "limit", 1, v)
	input.ActiveOnly = app.readBool(qs, "active_only", v)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	prompts, err := app.models.Prompts.GetRandom(input.TypeID, input.Limit, input.ActiveOnly)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"prompts": prompts, "size": len(prompts)}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showPromptHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	prompt, err := app.models.Prompts.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"prompt": prompt}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createPromptHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TypeID   int64  `json:"type_id"`
		Text     string `json:"text"`
		IsActive *bool  `json:"is_active"`
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

	prompt := &data.Prompt{
		Type: *typeObj,
		Text: input.Text,
	}

	if input.IsActive == nil {
		prompt.IsActive = true
	} else {
		prompt.IsActive = *input.IsActive
	}

	v := validator.New()
	if data.ValidatePrompt(v, prompt); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Prompts.Insert(prompt)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/prompts/%d", prompt.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"prompt": prompt}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deletePromptHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Prompts.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "prompt successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listPromptsHandler(w http.ResponseWriter, r *http.Request) {
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
	input.SortSafelist = []string{"id", "text", "type_id", "-id", "-name", "-type_id"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	prompts, metadata, err := app.models.Prompts.GetAll(input.Text, input.TypeID, input.ActiveOnly, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"prompts": prompts, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updatePromptHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	prompt, err := app.models.Prompts.Get(id)
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
		Text     *string `json:"name"`
		TypeID   *int64  `json:"type_id"`
		IsActive *bool   `json:"is_active"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Text != nil {
		prompt.Text = *input.Text
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
		prompt.Type = *typeObj
	}
	if input.IsActive != nil {
		prompt.IsActive = *input.IsActive
	}

	v := validator.New()
	if data.ValidatePrompt(v, prompt); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Prompts.Update(prompt)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"prompt": prompt}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
