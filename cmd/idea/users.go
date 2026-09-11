package main

import (
	"errors"
	"net/http"

	"github.com/fistos3rr/ideagen/internal/data"
	"github.com/fistos3rr/ideagen/internal/api/dto"
	"github.com/fistos3rr/ideagen/internal/validator"
)

// registerUserHandler godoc
//
// @Summary Register user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.UserCredentials true "Credentials"
// @Success 201 {object} dto.MessageResponse "User created"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 422 {object} dto.ValidationErrorResponse "Validation failed"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /register [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.UserCredentials

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Email: input.Email,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	resp := dto.MessageResponse{
		Message: "user created successfully",
	}
	err = app.writeJSON(w, http.StatusCreated, resp, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
