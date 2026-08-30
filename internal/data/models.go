package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound      = errors.New("record not found")
	ErrEditConflict        = errors.New("edit conflict")
	ErrForeignKeyViolation = errors.New("foreign key violation")
)

type Models struct {
	Types interface {
		Insert(t *Type) error
		Get(id int64) (*Type, error)
		GetAll(name string, activeOnly bool, filters Filters) ([]*Type, Metadata, error)
		Delete(id int64) error
		Update(t *Type) error
		GetRandom(limit int, activeOnly bool) ([]*Type, error)
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Types:   TypeModel{DB: db},
	}
}
