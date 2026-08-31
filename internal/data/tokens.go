package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type RefreshToken struct {
	ID        string    `json:"id"`
	UserID    int64     `json:"user_id"`
	TokenHash string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
}

type RefreshTokenModel struct {
	DB *sql.DB
}

func (m RefreshTokenModel) Insert(token *RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, revoked)
		VALUES ($1, $2, $3, $4, $5)
	`

	args := []any{token.ID, token.UserID, token.TokenHash, token.ExpiresAt, token.Revoked}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m RefreshTokenModel) GetByHash(hash string) (*RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, revoked
		FROM refresh_tokens
		WHERE token_hash = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var token RefreshToken
	err := m.DB.QueryRowContext(ctx, query, hash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.Revoked,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &token, nil
}

func (m RefreshTokenModel) DeleteByHash(hash string) error {
	query := `DELETE FROM refresh_tokens WHERE token_hash = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, hash)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m RefreshTokenModel) RevokeByHash(hash string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = true
		WHERE token_hash = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, hash)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m RefreshTokenModel) RevokeAllByUserID(userID int64) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = true
		WHERE user_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m RefreshTokenModel) DeleteExpired() error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query)
	return err
}
