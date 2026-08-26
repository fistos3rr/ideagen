package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fistos3rr/ideagen/internal/data"
	"github.com/fistos3rr/ideagen/internal/validator"
)

func (app *application) showTypeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	t, err := app.models.Types.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"type": t}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createTypeHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
		IsActive *bool `json:"is_active"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	t := &data.Type{
		Name: input.Name,
	}

	if input.IsActive == nil {
		t.IsActive = true
	} else {
		t.IsActive = *input.IsActive
	}

	v := validator.New()

	if data.ValidateType(v, t); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Types.Insert(t)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateType):
			v.AddError("name", "a type with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/types/%d", t.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"type": t}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteTypeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Types.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrForeignKeyViolation):
			app.badRequestResponse(w, r, err)
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "type successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listTypesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string
		ActiveOnly bool
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()
	
	input.Name = app.readString(qs, "name", "")
	input.ActiveOnly = app.readBool(qs, "active_only", v)
	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Sort = app.readString(qs, "sort", "id")
	input.SortSafelist = []string{"id", "name", "-id", "-name"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	types, metadata, err := app.models.Types.GetAll(input.Name, input.ActiveOnly, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"types": types, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
