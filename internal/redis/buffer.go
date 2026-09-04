package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fistos3rr/ideagen/internal/data"

	"github.com/redis/go-redis/v9"
)

type BufferIdea struct {
	ID      string     `json:"id"`
	Idea    *data.Idea `json:"idea"`
	Created time.Time  `json:"created"`
}

type BufferIdeasRepository struct {
	Client        *redis.Client
	TTL           time.Duration
	MaxBufferSize int64
}

func (repo *BufferIdeasRepository) Add(ctx context.Context, userID int64, idea *data.Idea) (*BufferIdea, error) {
	bufIdea := BufferIdea{
		ID:      data.GenerateUUID(),
		Idea:    idea,
		Created: time.Now(),
	}
	bufIdeaJSON, err := json.Marshal(bufIdea)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("buffer:idea:%d", userID)

	// Redis transaction analogue
	_, err = repo.Client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.RPush(ctx, key, bufIdeaJSON)
		pipe.LTrim(ctx, key, -repo.MaxBufferSize, -1)
		pipe.Expire(ctx, key, repo.TTL)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("redis pipeline execution: %w", err)
	}

	return &bufIdea, nil
}

func (repo *BufferIdeasRepository) GetAll(
	ctx context.Context,
	userID int64,
) ([]*BufferIdea, error) {
	key := fmt.Sprintf("buffer:idea:%d", userID)

	values, err := repo.Client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	ideas := make([]*BufferIdea, 0, len(values))

	for _, value := range values {
		var bufIdea BufferIdea

		if err := json.Unmarshal([]byte(value), &bufIdea); err != nil {
			return nil, fmt.Errorf("unmarshal buffer idea: %w", err)
		}

		ideas = append(ideas, &bufIdea)
	}

	return ideas, nil
}

func (repo *BufferIdeasRepository) Get(
	ctx context.Context,
	userID int64,
	bufIdeaID string,
) (*BufferIdea, error) {
	key := fmt.Sprintf("buffer:idea:%d", userID)

	values, err := repo.Client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("redis lrange: %w", err)
	}

	for _, value := range values {
		var bufIdea BufferIdea
		if err := json.Unmarshal([]byte(value), &bufIdea); err != nil {
			continue
		}

		if bufIdea.ID == bufIdeaID {
			return &bufIdea, nil
		}
	}

	return nil, ErrRecordNotFound
}

func (repo *BufferIdeasRepository) Clear(
	ctx context.Context,
	userID int64,
) error {
	key := fmt.Sprintf("buffer:idea:%d", userID)

	if err := repo.Client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis delete buffer for user %d: %w", userID, err)
	}

	return nil
}
