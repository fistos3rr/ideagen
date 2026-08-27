package ai

import "context"

type Config struct {
	APIKey string
	Model  string
	APIURL string
}

type Provider interface {
	SendMessage(ctx context.Context, message string, temperature float64, topP float64) (string, error)
}
