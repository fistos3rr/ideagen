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
	Prompts interface {
		Insert(p *Prompt) error
		Get(id int64) (*Prompt, error)
		GetAll(text string, typeID int64, activeOnly bool, filters Filters) ([]*Prompt, Metadata, error)
		Delete(id int64) error
		Update(prompt *Prompt) error
		GetRandom(typeID int64, limit int, activeOnly bool) ([]*Prompt, error)
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Types:   TypeModel{DB: db},
		Prompts: PromptModel{DB: db},
	}
}
