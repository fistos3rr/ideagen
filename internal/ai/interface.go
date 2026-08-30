// Package ai
package ai

import (
	"context"
	"errors"
)

var (
	ErrSetRequest = errors.New("ai setting request error")
)

type Config struct {
	APIKey string
	Model  string
	APIURL string
}

type Request interface {
	Clear(message string)
	AddMessage(message string)
	AddSystemMessage(message string)
}

type Provider interface {
	SendMessage(ctx context.Context, message string) (string, error)
	SendRequest(ctx context.Context) (string, error)
	SetRequest(req Request) error
}
