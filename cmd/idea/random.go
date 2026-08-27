package main

import (
	"net/http"
	"errors"

	"github.com/fistos3rr/ideagen/internal/validator"
	"github.com/fistos3rr/ideagen/internal/data"
)

func (app *application) generatePromptHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TypeID int64 `json:"type_id"`
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
	
	metaPrompt := app.metaGenerator.GenerateMetaPrompt(typeObj.Name)

	prompt, err := app.aiProvider.SendMessage(r.Context(), metaPrompt, 1.0, 1.0)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"prompt": prompt}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) generateMetaPromptHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TypeID   int64  `json:"type_id"`
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

	metaPrompt := app.metaGenerator.GenerateMetaPrompt(typeObj.Name)

	err = app.writeJSON(w, http.StatusOK, envelope{"meta_prompt": metaPrompt}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) randomTypesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Limit      int
		ActiveOnly bool
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Limit = app.readInt(qs, "limit", 1, v)
	input.ActiveOnly = app.readBool(qs, "active_only", v)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	types, err := app.models.Types.GetRandom(input.Limit, input.ActiveOnly)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"types": types, "size": len(types)}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
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
