package data

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Prompt struct {
	ID int64 `json:"id"`
	Type Type `json:"type"`
	Text string `json:"text"`
}

type PromptModel struct {
	DB *sql.DB
}

func (m PromptModel) Insert(p *Prompt) error {
	query := `
		INSERT INTO prompts (type_id, text)
		VALUES ($1, $2)
		RETURNING id
	`

	args := []any{p.Type.ID, p.Text}
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).
		Scan(&p.ID)
}

func (m PromptModel) GetById(id int64) (Prompt, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	
	query = `
		SELECT p.id, p.text, t.id, t.name
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
		&t.ID,
		&t.Name,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	s
	prompt.Type = t
	
	return &prompt, nil
}

func (m PromptModel) GetAll(
	text string,
	typeName string,
	filters Filters,
) ([]*Prompt, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), p.id, p.text, t.id, t.name
		FROM prompts p
		JOIN types t ON p.type_id = t.id
		WHERE 
			(to_tsvector('simple', p.text) @@ plainto_tsquery('simple', $1) OR $1 = '')
			AND
			(to_tsvector('simple', t.name) @@ plainto_tsquery('simple', $2) OR $2 = '')
		ORDER BY %s %s, p.id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	args := []any{
		text,
		typeName,
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
			&t.ID,
			&t.Name,
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

func (m *PromptModel) Delete(id int64) error {
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