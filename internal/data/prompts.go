package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/fistos3rr/ideagen/internal/validator"
)

type Prompt struct {
	ID   int64  `json:"id"`
	Type Type   `json:"type"`
	Text string `json:"text"`
	IsActive bool `json:"is_active"`
}

func ValidatePrompt(v *validator.Validator, prompt *Prompt) {
	v.Check(prompt.Text != "", "text", "must be provided")
	v.Check(len(prompt.Text) <= 5000, "text", "must not be more than 5000 bytes long")
	v.Check(prompt.Type.ID > 0, "type_id", "must be a valid type")
}

type PromptModel struct {
	DB *sql.DB
}

func (m PromptModel) Insert(p *Prompt) error {
	query := `
		INSERT INTO prompts (type_id, text, is_active)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	args := []any{p.Type.ID, p.Text, p.IsActive}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).
		Scan(&p.ID)
}

func (m PromptModel) Get(id int64) (*Prompt, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT p.id, p.text, p.is_active, t.id, t.name, t.is_active
		FROM prompts p
		JOIN types t ON p.type_id = t.id
		WHERE p.id = $1
	`

	var prompt Prompt
	var t Type

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&prompt.ID,
		&prompt.Text,
		&prompt.IsActive,
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
	prompt.Type = t

	return &prompt, nil
}

func (m PromptModel) GetAll(
	text string,
	typeID int64,
	activeOnly bool,
	filters Filters,
) ([]*Prompt, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), p.id, p.text, p.is_active, t.id, t.name, t.is_active
		FROM prompts p
		JOIN types t ON p.type_id = t.id
		WHERE 
			(to_tsvector('simple', p.text) @@ plainto_tsquery('simple', $1) OR $1 = '')
			AND
			($2 = 0 OR p.type_id = $2)
			AND
			($3 = false OR (p.is_active = $3 AND t.is_active = $3))
		ORDER BY p.%s %s, p.id ASC
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
	prompts := []*Prompt{}

	for rows.Next() {
		var p Prompt
		var t Type
		err := rows.Scan(
			&totalRecords,
			&p.ID,
			&p.Text,
			&p.IsActive,
			&t.ID,
			&t.Name,
			&t.IsActive,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		p.Type = t
		prompts = append(prompts, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return prompts, metadata, nil
}

func (m PromptModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM prompts
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

func (m PromptModel) Update(prompt *Prompt) error {
	query := `
		UPDATE prompts
		SET type_id = $1, text = $2, is_active = $3
		WHERE id = $4
	`

	args := []any{
		prompt.Type.ID,
		prompt.Text,
		prompt.IsActive,
		prompt.ID,
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
