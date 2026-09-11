package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/fistos3rr/ideagen/internal/auth"
	"github.com/fistos3rr/ideagen/internal/data"
	"github.com/fistos3rr/ideagen/internal/validator"
	"github.com/fistos3rr/ideagen/internal/api/dto"
)

// loginUserHandler godoc
//
// @Summary Login user
// @Description User authentication using email and password, and returns JWT-token.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User login data"
// @Success 200 {object} dto.LoginResponse "Successful authentication"
// @Header 200 {string} Set-Cookie "refresh_token=<jwt>; HttpOnly; Secure; SameSite=Strict; Path=/"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Invalid credentials"
// @Failure 422 {object} dto.ValidationErrorResponse "Validation error"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.LoginRequest

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	accessToken, refreshToken, refreshRecord, err := app.generateJWTTokenPair(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.models.RefreshTokens.Insert(refreshRecord); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  refreshRecord.ExpiresAt,
	})

	resp := dto.LoginResponse{
		AccessToken: accessToken,
	}

	app.writeJSON(w, http.StatusOK, resp, nil)
}

// refreshHandler godoc
//
// @Summary Refresh access token
// @Description Trade refresh-token from HttpOnly-cookie on new access token.
// @Description
// @Desctiption Read cookie `refresh_token`. If there is no cookie or token expired - returns 401.
// @Tags auth
// @Produce json
// @Success 200 {object} dto.LoginResponse "New access-token"
// @Failure 401 {object} dto.ErrorResponse "Invalid token"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /auth/refresh [post]
func (app *application) refreshHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}
	refreshToken := cookie.Value

	claims, err := auth.ParseRefreshToken(refreshToken, app.config.jwt.secret)
	if err != nil {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	hash := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hash[:])

	stored, err := app.models.RefreshTokens.GetByHash(tokenHash)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidAuthenticationTokenResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if stored == nil || stored.Revoked || stored.ExpiresAt.Before(time.Now()) {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	if err := app.models.RefreshTokens.DeleteByHash(tokenHash); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user, err := app.models.Users.Get(claims.UserID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	newAccess, newRefresh, newRecord, err := app.generateJWTTokenPair(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.models.RefreshTokens.Insert(newRecord); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefresh,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  newRecord.ExpiresAt,
	})

	resp := dto.LoginResponse {
		AccessToken: newAccess,	
	}

	app.writeJSON(w, http.StatusOK, resp, nil)
}

// logoutHandler godoc
//
// @Summary Logout
// @Desctiption Clears refresh token from cookie, delete session
// @Tags auth
// @Produce json
// @Success 200 {object} dto.MessageResponse
// @Failuer 400 {object} dto.ErrorResponse "Bad request"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /auth/logout [post]
func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		hash := sha256.Sum256([]byte(cookie.Value))
		tokenHash := hex.EncodeToString(hash[:])
		_ = app.models.RefreshTokens.DeleteByHash(tokenHash)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
	})

	resp := dto.MessageResponse{
		Message: "successfully logout",
	}
	app.writeJSON(w, http.StatusOK, resp, nil)
}
