package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fistos3rr/ideagen/internal/data"
	"github.com/fistos3rr/ideagen/internal/validator"
	"github.com/fistos3rr/ideagen/internal/api/dto"
)

// showTypeHandler godoc
//
// @Summary Get Type by id
// @Tags types
// @Produce json
// @Param id path int true "Type ID" example(42)
// @Success 200 {object} dto.TypeResponse
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Not enough privilleges"
// @Failure 404 {objcet} dto.ErrorResponse "Not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /types/{id} [get]
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

	resp := dto.TypeResponse{
		Type: t,
	}

	err = app.writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// createTypeHandler godoc
//
// @Summary Create Type
// @Description Create Type for Ideas
// @Tags types
// @Accept json
// @Produce json
// @Param request body dto.TypeRequest true "Type data"
// @Success 201 {object} dto.TypeResponse "Type created"
// @Header 201 {string} Location "URL of created resource, for example: /v1/types/42"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Not enough privilleges"
// @Failure 422 {object} dto.ValidationErrorResponse "Validation failed"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /types [post]
func (app *application) createTypeHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.TypeRequest

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

	resp := dto.TypeResponse{
		Type: t,
	}
	err = app.writeJSON(w, http.StatusCreated, resp, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteTypeHandler godoc
//
// @Summary Delete type
// @Description Delete type by ID
// @Tags types
// @Produce json
// @Param id path int true "Type ID" example(42)
// @Success 200 {object} dto.MessageResponse "Type deleted successfully"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Not enough privilleges"
// @Failure 404 {objcet} dto.ErrorResponse "Not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /types/{id} [delete]
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

	resp := dto.MessageResponse{
		Message: "type successfully deleted",
	}
	err = app.writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listTypesHandler godoc
//
// @Summary Get Type list
// @Description Get all types pages with filters
// @Tags types
// @Produce json
// @Param name query string false "Type name"
// @Param active_only query boolean false "Active only"
// @Param page query int false "Page number" 
// @Param page_size query int false "Page size"
// @Param sort query string false "Sort by" Enum(id, -id, name, -name)
// @Success 200 {object} dto.TypeListResponse
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Not enough privilleges"
// @Failure 404 {objcet} dto.ErrorResponse "Not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /types [get]
func (app *application) listTypesHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.TypeListRequest

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

	resp := dto.TypeListResponse{
		Types: types,
		Metadata: metadata,
	}

	err = app.writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateTypeHandler godoc
//
// @Summary Update type 
// @Description Update type fields, for example change IsActive status
// @Tags types
// @Accept json
// @Produce json
// @Param id path int true "Type ID" example(42)
// @Param request body dto.TypeUpdateRequest false "Type update data"
// @Success 200 {object} dto.TypeResponse "Type updated"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Not enough privilleges"
// @Failure 422 {object} dto.ValidationErrorResponse "Validation failed"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /types/{id} [patch]
func (app *application) updateTypeHandler(w http.ResponseWriter, r *http.Request) {
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

	var input dto.TypeUpdateRequest

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		t.Name = *input.Name
	}
	if input.IsActive != nil {
		t.IsActive = *input.IsActive
	}

	v := validator.New()

	if data.ValidateType(v, t); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Types.Update(t)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	resp := TypeResponse{
		Type: t,
	}

	err = app.writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
