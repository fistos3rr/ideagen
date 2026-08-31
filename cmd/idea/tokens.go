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
)

func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

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

	app.writeJSON(w, http.StatusOK, envelope{"access_token": accessToken}, nil)
}

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
		app.serverErrorResponse(w, r, err)
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

	app.writeJSON(w, http.StatusOK, envelope{"access_token": newAccess}, nil)
}

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

	app.writeJSON(w, http.StatusOK, envelope{"message": "successfully logout"}, nil)
}
