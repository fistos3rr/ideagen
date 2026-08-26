package data

import (
	"errors"
	"database/sql"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict = errors.New("edit conflict")
)

type Models struct {
	Types interface {
		Insert(t *Type) error		
		GetById(id int64) (*Type, error)
		GetAll(name string, filters Filters) ([]*Type, Metadata, error) 
		Delete(id int64) error
	}
	Prompts interface {
		Insert(p *Prompt) error
		GetById(id int64) (*Prompt, error)
		GetAll(text string,	typeName string, filters Filters) ([]*Prompt, Metadata, error)
		Delete(id int64) error
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Types: TypeModel{DB: db},
		Prompts: PromptModel{DB: db},
	}
}
