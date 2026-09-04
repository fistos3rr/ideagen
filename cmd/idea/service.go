package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fistos3rr/ideagen/internal/data"
	"github.com/fistos3rr/ideagen/internal/redis"
	"github.com/fistos3rr/ideagen/internal/validator"
)

func (app *application) showMeHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listMyIdeasHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	var input struct {
		UserID     int64
		Text       string
		TypeID     int64
		ActiveOnly bool
		Status     data.UserIdeaStatus
		data.Filters
	}

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

	err = app.writeJSON(w, http.StatusOK, envelope{"ideas": ideas, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

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

	err = app.writeJSON(w, http.StatusOK, envelope{"idea": idea}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

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

	err = app.models.UserIdeas.DeleteById(user.ID, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.writeJSON(w, http.StatusNoContent, envelope{}, nil)
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

func (app *application) readUserIdeaIDParams(r *http.Request) (int64, int64, error) {
	pair, err := app.readIDParams(r, "user_id", "idea_id")
	if err != nil {
		return 0, 0, err
	}

	if _, ok := pair["user_id"]; !ok {
		panic("no user_id provided")
	}

	if _, ok := pair["idea_id"]; !ok {
		panic("no idea_id provided")
	}

	return pair["user_id"], pair["idea_id"], nil
}

func (app *application) createUserIdeaHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID int64 `json:"user_id"`
		IdeaID int64 `json:"idea_id"`
	}

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

	err = app.writeJSON(w, http.StatusCreated, envelope{}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) readUserIDParam(r *http.Request) (int64, error) {
	ids, err := app.readIDParams(r, "user_id")
	if err != nil {
		return 0, err
	}

	if _, ok := ids["user_id"]; !ok {
		panic("no user_id provided")
	}

	return ids["user_id"], nil
}

func (app *application) readIdeaIDParam(r *http.Request) (int64, error) {
	ids, err := app.readIDParams(r, "idea_id")
	if err != nil {
		return 0, err
	}

	if _, ok := ids["idea_id"]; !ok {
		panic("no idea_id provided")
	}

	return ids["idea_id"], nil
}

func (app *application) listIdeasByUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID     int64
		Text       string
		TypeID     int64
		ActiveOnly bool
		Status     data.UserIdeaStatus
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	var err error
	input.UserID, err = app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	input.Status = data.UserIdeaStatus(app.readInt(qs, "status", 0, v))
	input.Text = app.readString(qs, "text", "")
	input.TypeID = int64(app.readInt(qs, "type_id", 0, v))
	input.ActiveOnly = app.readBool(qs, "active_only", v)
	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Sort = app.readString(qs, "sort", "id")
	input.SortSafelist = []string{"role", "email", "created_at", "-email", "-created_at", "-role"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	ideas, metadata, err := app.models.UserIdeas.GetIdeasByUserID(input.UserID, input.Text, input.TypeID, input.ActiveOnly, input.Status, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"ideas": ideas, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listUsersByIdeaHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		IdeaID int64
		Role   string
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	var err error
	input.IdeaID, err = app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	input.Role = app.readString(qs, "role", "")
	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Sort = app.readString(qs, "sort", "id")
	input.SortSafelist = []string{"role", "email", "created_at", "-email", "-created_at", "-role"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	users, metadata, err := app.models.UserIdeas.GetUsersByIdeaID(input.IdeaID, input.Role, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"users": users, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

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

	bufIdea, err := app.redisRepository.BufferIdeas.Add(r.Context(), user.ID, idea)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"buffer_idea": bufIdea}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) listMyBufferIdeasHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	bufIdeas, err := app.redisRepository.BufferIdeas.GetAll(r.Context(), user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"buffer_ideas": bufIdeas}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) chooseMyBufferIdeaHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	var input struct {
		BufIdeaID string `json:"buffer_idea_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	bufIdea, err := app.redisRepository.BufferIdeas.Get(r.Context(), user.ID, input.BufIdeaID)
	if err != nil {
		switch {
		case errors.Is(err, redis.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	idea := bufIdea.Idea

	app.models.Ideas.Insert(idea)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.redisRepository.BufferIdeas.Clear(r.Context(), user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/service/ideas/%d", idea.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"idea": idea}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

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

	bufIdea, err := app.redisRepository.BufferIdeas.Get(r.Context(), user.ID, id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"buffer_idea": bufIdea}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
