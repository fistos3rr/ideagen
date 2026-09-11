package dto

type MessageResponse struct {
	Message string `json:"message" example:"string"`
}

type HealthResponse struct {
	Status string `json:"status"`	
	Environment string `json:"environment"`
}
