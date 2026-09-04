package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fistos3rr/ideagen/internal/validator"
)

type UserIdeaStatus int

const (
	UserIdeaActive UserIdeaStatus = iota + 1
	UserIdeaCompleted
	UserIdeaArchived
)

func (s UserIdeaStatus) IsValid() bool {
	switch s {
	case UserIdeaActive, UserIdeaCompleted, UserIdeaArchived:
		return true
	default:
		return false
	}
}

type UserIdea struct {
	UserID int64          `json:"user_id"`
	IdeaID int64          `json:"idea_id"`
	Status UserIdeaStatus `json:"is_active"`
}

type UserIdeasModel struct {
	DB *sql.DB
}

func ValidateUserStatus(v *validator.Validator, status UserIdeaStatus) {
	v.Check(status.IsValid(), "status", "status must be valid")
}

func ValidateUserIdea(v *validator.Validator, ui *UserIdea) {
	ValidateUserStatus(v, ui.Status)
}

func (m UserIdeasModel) Insert(user *User, idea *Idea) error {
	if user == nil || idea == nil {
		return ErrRecordNotFound
	}

	query := `
		INSERT INTO user_ideas (user_id, idea_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, idea_id) DO NOTHING;
	`

	args := []any{user.ID, idea.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDuplicateRecord
	}

	return nil
}

func (m UserIdeasModel) DeleteById(userID int64, ideaID int64) error {
	query := `DELETE FROM user_ideas WHERE user_id = $1 AND idea_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, userID, ideaID)
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

func (m UserIdeasModel) Delete(user *User, idea *Idea) error {
	if user == nil {
		return fmt.Errorf("user: %w", ErrRecordNotFound)
	}
	if idea == nil {
		return fmt.Errorf("idea: %w", ErrRecordNotFound)
	}

	query := `DELETE FROM user_ideas WHERE user_id = $1 AND idea_id = $2`

	args := []any{
		user.ID,
		idea.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, args)
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

func (m UserIdeasModel) GetIdeasByUserID(
	userID int64,
	text string,
	typeID int64,
	activeOnly bool,
	status UserIdeaStatus,
	filters Filters,
) ([]*Idea, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), i.id, i.text, t.id, t.name, t.is_active
		FROM ideas i
		JOIN types t ON i.type_id = t.id
		JOIN user_ideas ui ON i.id = ui.idea_id
		WHERE ui.user_id = $1
			AND (to_tsvector('simple', i.text) @@ plainto_tsquery('simple', $2) OR $2 = '')
			AND ($3 = 0 OR i.type_id = $3)
			AND ($4 = false OR t.is_active = $4)
			AND ($7 = 0 OR ui.status = $7)
		ORDER BY i.%s %s, i.id ASC
		LIMIT $5 OFFSET $6`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		userID,
		text,
		typeID,
		activeOnly,
		filters.limit(),
		filters.offset(),
		status,
	}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	ideas := []*Idea{}

	for rows.Next() {
		var idea Idea
		var t Type
		err := rows.Scan(
			&totalRecords,
			&idea.ID,
			&idea.Text,
			&t.ID,
			&t.Name,
			&t.IsActive,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		idea.Type = &t
		ideas = append(ideas, &idea)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return ideas, metadata, nil
}

func (m UserIdeasModel) GetUsersByIdeaID(
	ideaID int64,
	role string,
	filters Filters,
) ([]*User, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), u.id, u.created_at, u.email, u.role
		FROM users u
		JOIN user_ideas ui ON u.id = ui.idea_id
		WHERE ui.idea_id = $1
			AND (to_tsvector('simple', u.role) @@ plainto_tsquery('simple', $2) OR $2 = '')
		ORDER BY i.%s %s, i.id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		ideaID,
		role,
		filters.limit(),
		filters.offset(),
	}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	users := []*User{}

	for rows.Next() {
		var user User
		err := rows.Scan(
			&totalRecords,
			&user.ID,
			&user.CreatedAt,
			&user.Email,
			&user.Role,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return users, metadata, nil
}

func (m UserIdeasModel) ExistsUserIdea(userID int64, ideaID int64) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM user_ideas
			WHERE user_id = $1 AND idea_id = $2
		)
	`

	var exists bool
	args := []any{
		userID,
		ideaID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (m UserIdeasModel) UpdateUserIdea(ui *UserIdea) error {
	exists, err := m.ExistsUserIdea(ui.UserID, ui.IdeaID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrRecordNotFound
	}

	query := `
		UPDATE user_ideas
		SET status = $3
		WHERE user_id = $1 AND idea_id = $2
	`

	args := []any{
		ui.UserID,
		ui.IdeaID,
		ui.Status,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrEditConflict
	}

	return nil
}
