package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fistos3rr/ideagen/internal/data"
	"github.com/fistos3rr/ideagen/internal/api/dto"
	"github.com/fistos3rr/ideagen/internal/validator"
)

// showMeHandler godoc
//
// @Summary Get me
// @Tags service
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /service/me [get]
func (app *application) showMeHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	resp := dto.NewUserResponse(user)

	err := app.writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listMyIdeasHandler godoc
//
// @Summary Get Idea list
// @Description Get all ideas pages with filters
// @Tags service
// @Produce json
// @Param status query int false "Idea status (0 - all, 1 - active, 2 - completed, 3 - archived)"
// @Param text query string false "Idea text"
// @Param type_id query int false "Type ID"
// @Param active_only query boolean false "Active only"
// @Param page query int false "Page number" 
// @Param page_size query int false "Page size"
// @Param sort query string false "Sort by" Enum("id", "name", "type_id", "-id", "-name", "-type_id")
// @Success 200 {object} dto.IdeaListResponse
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 422 {object} dto.ValidationErrorResponse "Validation failed"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /service/ideas [get]
func (app *application) listMyIdeasHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	var input dto.IdeaUserListRequest

	v := validator.New()
	qs := r.URL.Query()

	input.Status = data.UserIdeaStatus(app.readInt(qs, "status", 0, v))
	input.UserID = user.ID
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

	ideas, metadata, err := app.models.UserIdeas.GetIdeasByUserID(
		input.UserID,
		input.Text,
		input.TypeID,
		input.ActiveOnly,
		input.Status,
		input.Filters,
	)
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


// showMyIdeaHandler godoc
//
// @Summary Get My Idea by id
// @Tags service
// @Produce json
// @Param id path int true "Idea ID" example(42)
// @Success 200 {object} dto.IdeaResponse
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 404 {object} dto.ErrorResponse "Not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /service/ideas/{id} [get]
func (app *application) showMyIdeaHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	exists, err := app.models.UserIdeas.ExistsUserIdea(user.ID, id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !exists {
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

// deleteIdeaHandler godoc
//
// @Summary Delete my idea
// @Description Delete my idea by ID
// @Tags service
// @Produce json
// @Param id path int true "Idea ID" example(42)
// @Success 200 {object} dto.MessageResponse "Idea deleted successfully"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 404 {object} dto.ErrorResponse "Not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /service/ideas/{id} [delete]
func (app *application) deleteMyIdeaHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var resp dto.MessageResponse

	err = app.models.UserIdeas.DeleteById(user.ID, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.writeJSON(w, http.StatusNoContent, resp, nil)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	
	resp.Message = "idea successfully deleted"
	err = app.writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}


// createUserIdeaHandler godoc
//
// @Summary Create user-idea bound
// @Description Create bound between user and idea, so idea became user's
// @Tags service
// @Accept json
// @Produce json
// @Param request body dto.UserIdeaRequest true "User and idea ID"
// @Success 201 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Not enough privilleges"
// @Failure 422 {object} dto.ValidationErrorResponse "Validation failed"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /service/useridea [post]
func (app *application) createUserIdeaHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.UserIdeaRequest

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	user, err := app.models.Users.Get(input.UserID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			v := validator.New()
			v.AddError("user_id", "user with this id does not exists")
			app.failedValidationResponse(w, r, v.Errors)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	idea, err := app.models.Ideas.Get(input.IdeaID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v := validator.New()
			v.AddError("idea_id", "idea with this id does not exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.UserIdeas.Insert(user, idea)
	if err != nil {
		if errors.Is(err, data.ErrDuplicateRecord) {
			app.badRequestResponse(w, r, err)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	resp := dto.MessageResponse{
		Message: "User-Idea bound created",
	}

	err = app.writeJSON(w, http.StatusCreated, resp, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// generateMyIdeaHandler godoc
//
// @Summary Generate random idea
// @Description Generates random ideas and push them into buffer
// @Tags service
// @Produce json
// @Success 200 {object} dto.BufferIdeaResponse "Idea generated in buffer"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /service/idea/generate [post]
func (app *application) generateMyIdeaHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	var t *data.Type
	types, err := app.models.Types.GetRandom(1, true)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if len(types) == 0 {
		app.notFoundResponse(w, r)
		return
	}
	t = types[0]

	idea, err := app.generateIdea(r.Context(), t)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	bufIdea, err := app.models.BufferIdeas.Add(r.Context(), user.ID, idea)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	resp := dto.BufferIdeaResponse{
		BufferIdea: bufIdea,
	}

	err = app.writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}


// listMyBufferIdeasHandler godoc
//
// @Summary Get my generated ideas from buffer
// @Tags service
// @Produce json
// @Success 200 {object} dto.BufferIdeaListResponse
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /service/idea/buffer [get]
func (app *application) listMyBufferIdeasHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	bufIdeas, err := app.models.BufferIdeas.GetAll(r.Context(), user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	resp := dto.BufferIdeaListResponse{
		BufferIdeas: bufIdeas,
	}

	err = app.writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// chooseMyBufferIdeaHandler godoc
//
// @Summary Choose buffer idea by UUID
// @Description Choose idea from buffer by UUID, pushes that idea into database, clears buffer
// @Tags service
// @Accept json
// @Produce json
// @Param request body dto.BufferIdeaRequest true "Buffer Idea UUID"
// @Success 201 {object} dto.IdeaResponse "Idea created"
// @Header 201 {string} Location "URL of created resource, for example: /service/ideas/{id}"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /service/idea/buffer [post]
func (app *application) chooseMyBufferIdeaHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	var input dto.BufferIdeaRequest

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	bufIdea, err := app.models.BufferIdeas.Get(r.Context(), user.ID, input.BufferIdeaID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	idea := bufIdea.Idea

	err = app.models.Ideas.Insert(idea)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.UserIdeas.Insert(user, idea)
	if err != nil {
		if errors.Is(err, data.ErrDuplicateRecord) {
			app.badRequestResponse(w, r, err)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.BufferIdeas.Clear(r.Context(), user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/service/ideas/%d", idea.ID))

	resp := dto.IdeaResponse{
		Idea: idea,
	}

	err = app.writeJSON(w, http.StatusCreated, resp, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}


// showMyBufferIdeaHandler godoc
//
// @Summary Get Idea from buffer by UUID
// @Tags service
// @Produce json
// @Param id path string true "Buffer Idea UUID"
// @Success 200 {object} dto.BufferIdeaResponse
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 404 {object} dto.ErrorResponse "Not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /service/idea/buffer/{id} [get]
func (app *application) showMyBufferIdeaHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	id, err := app.readStringIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	bufIdea, err := app.models.BufferIdeas.Get(r.Context(), user.ID, id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	resp := dto.BufferIdeaResponse{
		BufferIdea: bufIdea,
	}
	err = app.writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
