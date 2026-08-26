package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/fistos3rr/ideagen/internal/validator"
)

var (
	ErrDuplicateType = errors.New("duplicate type")
)

type Type struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

func ValidateType(v *validator.Validator, t *Type) {
	v.Check(t.Name != "", "name", "must be provided")
	v.Check(len(t.Name) <= 80, "name", "must not be more than 80 bytes long")
}

type TypeModel struct {
	DB *sql.DB
}

func (m TypeModel) Insert(t *Type) error {
	query := `
		INSERT INTO types (name, is_active)
		VALUES ($1, $2)
		RETURNING id
	`

	args := []any{t.Name, t.IsActive}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&t.ID)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "types_name_key" (23505)`:
			return ErrDuplicateType
		default:
			return err
		}
	}

	return nil
}

func (m TypeModel) Get(id int64) (*Type, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, name, is_active
		FROM types
		WHERE id = $1
	`

	var t Type

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&t.ID, &t.Name, &t.IsActive)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &t, nil
}

func (m TypeModel) GetAll(
	name string,
	activeOnly bool,
	filters Filters,
) ([]*Type, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, name, is_active 
		FROM types 
		WHERE 
			(to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
			AND
			($2 = false OR is_active = $2)
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		name,
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
	types := []*Type{}

	for rows.Next() {
		var t Type
		err := rows.Scan(
			&totalRecords,
			&t.ID,
			&t.Name,
			&t.IsActive,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		types = append(types, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return types, metadata, nil
}

func (m TypeModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM types
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		if isForeignKeyViolation(err) {
			return ErrForeignKeyViolation
		}
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

func (m TypeModel) Update(t *Type) error {
	query := `
		UPDATE types
		SET name = $1, is_active = $2
		WHERE id = $3
	`

	args := []any{
		t.Name,
		t.IsActive,
		t.ID,
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
