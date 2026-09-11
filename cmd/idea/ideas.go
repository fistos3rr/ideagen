package main

import (
	"errors"
	"net/http"
	"fmt"

	"github.com/fistos3rr/ideagen/internal/data"
	"github.com/fistos3rr/ideagen/internal/api/dto"
	"github.com/fistos3rr/ideagen/internal/validator"
)

// showIdeaHandler godoc
//
// @Summary Get Idea by id
// @Tags ideas
// @Produce json
// @Param id path int true "Idea ID" example(42)
// @Success 200 {object} dto.IdeaResponse
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Not enough privilleges"
// @Failure 404 {object} dto.ErrorResponse "Not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /ideas/{id} [get]
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

	resp := dto.IdeaResponse{
		Idea: idea,
	}

	err = app.writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// createIdeaHandler godoc
//
// @Summary Create Idea
// @Description Create Idea manually (no random)
// @Tags ideas
// @Accept json
// @Produce json
// @Param request body dto.IdeaRequest true "Idea data"
// @Success 201 {object} dto.IdeaResponse "Idea created"
// @Header 201 {string} Location "URL of created resource, for example: /v1/ideas/42"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Not enough privilleges"
// @Failure 422 {object} dto.ValidationErrorResponse "Validation failed"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /ideas [post]
func (app *application) createIdeaHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.IdeaRequest

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

	resp := dto.IdeaResponse{
		Idea: idea,
	}

	err = app.writeJSON(w, http.StatusCreated, resp, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteIdeaHandler godoc
//
// @Summary Delete idea
// @Description Delete idea by ID
// @Tags ideas
// @Produce json
// @Param id path int true "Idea ID" example(42)
// @Success 200 {object} dto.MessageResponse "Idea deleted successfully"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Not enough privilleges"
// @Failure 404 {object} dto.ErrorResponse "Not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /ideas/{id} [delete]
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

	resp := dto.MessageResponse{
		Message: "idea successfully deleted",
	}

	err = app.writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listIdeasHandler godoc
//
// @Summary Get Idea list
// @Description Get all ideas pages with filters
// @Tags ideas
// @Produce json
// @Param text query string false "Idea text"
// @Param type_id query int false "Type ID"
// @Param active_only query boolean false "Active only"
// @Param page query int false "Page number" 
// @Param page_size query int false "Page size"
// @Param sort query string false "Sort by" Enum("id", "name", "type_id", "-id", "-name", "-type_id")
// @Success 200 {object} dto.IdeaListResponse
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Not enough privilleges"
// @Failure 422 {object} dto.ValidationErrorResponse "Validation failed"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /ideas [get]
func (app *application) listIdeasHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.IdeaListRequest

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

	resp := dto.IdeaListResponse{
		Ideas: ideas,
		Metadata: metadata,
	}

	err = app.writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateIdeaHandler godoc
//
// @Summary Update idea 
// @Description Patch idea's fields
// @Tags ideas
// @Accept json
// @Produce json
// @Param id path int true "Idea ID" example(42)
// @Param request body dto.IdeaUpdateRequest false "Idea update data"
// @Success 200 {object} dto.IdeaResponse "Idea updated"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Not enough privilleges"
// @Failure 422 {object} dto.ValidationErrorResponse "Validation failed"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /ideas/{id} [patch]
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

	var input dto.TypeUpdateRequest

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

	resp := dto.IdeaResponse{
		Idea: idea,
	}

	err = app.writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
