package dto

type AskRequest struct {
	Message string `json:"message" example:"What year now?"`
}

type AskResponse struct {
	Answer string `json:"answer" example:"2026"`
}
