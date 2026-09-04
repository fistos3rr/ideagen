package redis

import (
	"context"
	"errors"
	"time"

	"github.com/fistos3rr/ideagen/internal/data"

	"github.com/redis/go-redis/v9"
)

var ErrRecordNotFound = errors.New("record not found")

type Config struct {
	IdeaTTL time.Duration
}

type Repository struct {
	BufferIdeas interface {
		Add(ctx context.Context, userID int64, idea *data.Idea) (*BufferIdea, error)
		GetAll(ctx context.Context, userID int64) ([]*BufferIdea, error)
		Get(ctx context.Context, userID int64, bufIdeaID string) (*BufferIdea, error)
		Clear(ctx context.Context, userID int64) error
	}
}

func NewRepository(client *redis.Client, config Config) Repository {
	return Repository{
		BufferIdeas: &BufferIdeasRepository{
			Client:        client,
			TTL:           config.IdeaTTL,
			MaxBufferSize: 10,
		},
	}
}
