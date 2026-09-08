package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrRecordNotFound      = errors.New("record not found")
	ErrEditConflict        = errors.New("edit conflict")
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrDuplicateRecord     = errors.New("duplicate record")
)

type Config struct {
	IdeaTTL         time.Duration
	RefreshTokenTTL time.Duration
}

type Models struct {
	Types interface {
		Insert(t *Type) error
		Get(id int64) (*Type, error)
		GetAll(name string, activeOnly bool, filters Filters) ([]*Type, Metadata, error)
		Delete(id int64) error
		Update(t *Type) error
		GetRandom(limit int, activeOnly bool) ([]*Type, error)
	}
	Ideas interface {
		Insert(idea *Idea) error
		Get(id int64) (*Idea, error)
		GetAll(text string, typeID int64, activeOnly bool, filters Filters) ([]*Idea, Metadata, error)
		Delete(id int64) error
		Update(idea *Idea) error
	}
	Users interface {
		Insert(user *User) error
		Get(id int64) (*User, error)
		GetByEmail(email string) (*User, error)
		Update(user *User) error
	}
	RefreshTokens interface {
		Insert(token *RefreshToken) error
		GetByHash(hash string) (*RefreshToken, error)
		DeleteByHash(hash string) error
		RevokeByHash(hash string) error
		RevokeAllByUserID(userID int64) error
		DeleteExpired() error
	}
	UserIdeas interface {
		Insert(user *User, idea *Idea) error
		Delete(user *User, idea *Idea) error
		DeleteById(userID int64, ideaID int64) error
		ExistsUserIdea(userID int64, ideaID int64) (bool, error)
		UpdateUserIdea(ui *UserIdea) error
		GetIdeasByUserID(
			userID int64,
			text string,
			typeID int64,
			activeOnly bool,
			status UserIdeaStatus,
			filters Filters,
		) ([]*Idea, Metadata, error)
		GetUsersByIdeaID(
			ideaID int64,
			role string,
			filters Filters,
		) ([]*User, Metadata, error)
	}
	BufferIdeas interface {
		Add(ctx context.Context, userID int64, idea *Idea) (*BufferIdea, error)
		GetAll(ctx context.Context, userID int64) ([]*BufferIdea, error)
		Get(ctx context.Context, userID int64, bufIdeaID string) (*BufferIdea, error)
		Clear(ctx context.Context, userID int64) error
	}
}

func NewModels(db *sql.DB, client *redis.Client, config Config) Models {
	return Models{
		Types: TypeModel{DB: db},
		Ideas: IdeaModel{DB: db},
		Users: UserModel{DB: db},
		// RefreshTokens: RefreshTokenModel{DB: db},
		RefreshTokens: BufferRefreshTokenModel{
			Client: client,
			TTL:    config.RefreshTokenTTL,
		},
		UserIdeas: UserIdeasModel{DB: db},
		BufferIdeas: BufferIdeasModel{
			Client:        client,
			TTL:           config.IdeaTTL,
			MaxBufferSize: 10,
		},
	}
}
