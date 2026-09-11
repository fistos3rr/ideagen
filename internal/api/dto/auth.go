package dto

type LoginRequest struct {
	Email string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"secret"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token" example:"eyJ..."`
}
