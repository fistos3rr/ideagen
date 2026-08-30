// Package data
package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/fistos3rr/ideagen/internal/validator"
)

type Idea struct {
	ID   int64  `json:"id"`
	Text string `json:"text"`
	Type *Type  `json:"type,omitempty"`
}

func ValidateIdea(v *validator.Validator, idea *Idea) {
	v.Check(idea.Text != "", "text", "must be provided")
	v.Check(len(idea.Text) <= 5000, "text", "must not be more than 5000 bytes long")
	v.Check(idea.Type.ID > 0, "type_id", "must be a valid type")
}

type IdeaModel struct {
	DB *sql.DB
}

func (m IdeaModel) Insert(idea *Idea) error {
	query := `
		INSERT INTO ideas (text, type_id)
		VALUES($1, $2, $3)
		RETURNING id
	`

	args := []any{idea.Type.ID, idea.Text}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).
		Scan(&idea.ID)
}

func (m IdeaModel) Get(id int64) (*Idea, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT i.id, i.text, t.id, t.name, t.is_active
		FROM ideas i
		JOIN types t ON i.type_id = t.id
		WHERE i.id = $1
	`

	var idea Idea
	var t Type

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&idea.ID,
		&idea.Text,
		&t.ID,
		&t.Name,
		&t.IsActive,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	idea.Type = &t

	return &idea, nil
}

func (m IdeaModel) GetAll(
	text string,
	typeID int64,
	activeOnly bool,
	filters Filters,
) ([]*Idea, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), i.id, i.text, t.id, t.name, t.is_active
		FROM ideas i
		JOIN types t ON i.type_id = t.id
		WHERE
			(to_tsvector('simple', i.text) @@ plainto_tsquery('simple', $1) OR $1 = '')
			AND
			($2 = 0 OR i.type_id = $2)
			AND
			($3 = false OR t.is_active = $3)
		ORDER BY i.%s %s, i.id ASC
		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		text,
		typeID,
		activeOnly,
		filters.limit(),
		filters.offset(),
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

func (m IdeaModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM ideas
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
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

func (m IdeaModel) Update(idea *Idea) error {
	query := `
		UPDATE ideas
		SET type_id = $1, text = $2
		WHERE id = $3
	`

	args := []any{
		idea.Type.ID,
		idea.Text,
		idea.ID,
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

func (m IdeaModel) GetRandom(
	typeID int64,
	limit int,
	activeOnly bool,
) ([]*Idea, error) {
	query := `
		SELECT i.id, i.text, t.id, t.name, t.is_active
		FROM ideas i
		JOIN types t ON i.type_id = t.id
		WHERE
			($1 = 0 OR i.type_id = $1)
			AND
			($2 = false OR t.is_active = $2)
		ORDER BY RANDOM()
		LIMIT $3
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		typeID,
		activeOnly,
		limit,
	}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ideas := []*Idea{}

	for rows.Next() {
		var idea Idea
		var t Type
		err := rows.Scan(
			&idea.ID,
			&idea.Text,
			&t.ID,
			&t.Name,
			&t.IsActive,
		)
		if err != nil {
			return nil, err
		}
		idea.Type = &t
		ideas = append(ideas, &idea)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ideas, nil
}
