package dto

type ErrorResponse struct {
	Error string `json:"error" example:"something went wrong"`
}

type ValidationErrorResponse struct {
	Error map[string]string `json:"error" example:"{\"email\":\"must be a valid email address\"}"`
}
