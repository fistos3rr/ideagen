package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type BufferRefreshTokenModel struct {
	Client *redis.Client
	TTL    time.Duration
}

func tokenKey(hash string) string {
	return fmt.Sprintf("refresh_token:%s", hash)
}

func userTokensKey(userID int64) string {
	return fmt.Sprintf("user_tokens:%d", userID)
}

func (m BufferRefreshTokenModel) Insert(token *RefreshToken) error {
	bufTokenJSON, err := json.Marshal(token)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = m.Client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Set(ctx, tokenKey(token.TokenHash), bufTokenJSON, m.TTL)
		pipe.SAdd(ctx, userTokensKey(token.UserID), token.TokenHash)
		return nil
	})
	if err != nil {
		return fmt.Errorf("redis pipeline execution: %w", err)
	}

	return nil
}

func (m BufferRefreshTokenModel) GetByHash(hash string) (*RefreshToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	bufToken, err := m.Client.Get(ctx, tokenKey(hash)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	var token RefreshToken
	if err := json.Unmarshal([]byte(bufToken), &token); err != nil {
		return nil, err
	}

	return &token, nil
}

func (m BufferRefreshTokenModel) DeleteByHash(hash string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	token, err := m.GetByHash(hash)
	if err != nil {
		return err
	}

	_, err = m.Client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Del(ctx, tokenKey(hash))
		pipe.SRem(ctx, userTokensKey(token.UserID), hash)
		return nil
	})
	if err != nil {
		return fmt.Errorf("redis pipeline execution: %w", err)
	}

	return nil
}

func (m BufferRefreshTokenModel) RevokeByHash(hash string) error {
	return m.DeleteByHash(hash)
}

func (m BufferRefreshTokenModel) RevokeAllByUserID(userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	uKey := userTokensKey(userID)

	hashes, err := m.Client.SMembers(ctx, uKey).Result()
	if err != nil {
		return err
	}

	_, err = m.Client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, hash := range hashes {
			pipe.Del(ctx, tokenKey(hash))
		}
		pipe.Del(ctx, uKey)
		return nil
	})
	if err != nil {
		return fmt.Errorf("redis pipeline execution: %w", err)
	}

	return nil
}

func (m BufferRefreshTokenModel) DeleteExpired() error {
	return nil
}
