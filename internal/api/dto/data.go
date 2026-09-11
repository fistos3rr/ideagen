package dto

import (
	"time"

	"github.com/fistos3rr/ideagen/internal/data"
)

type TypeResponse struct {
	Type *data.Type `json:"type"`
}

type TypeRequest struct {
	Name     string `json:"name"`
	IsActive *bool  `json:"is_active"`
}

type TypeListRequest struct {
	Name string `json:"name,omitempty"`
	ActiveOnly bool `json:"bool,omitempty"`
	data.Filters
}

type TypeListResponse struct {
	Types []*data.Type `json:"types"`
	Metadata data.Metadata `json:"metadata"`
}

type TypeUpdateRequest struct {
	Name     *string `json:"name,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
}

type IdeaResponse struct {
	Idea *data.Idea `json:"idea"`
}

type IdeaRequest struct {
	TypeID int64 `json:"type_id"`
	Text string `json:"text"`
}

type IdeaListRequest struct {
	Text       string `json:"text,omitempty"`
	TypeID     int64 `json:"type_id,omitempty"`
	ActiveOnly bool `json:"active_only,omitempty"`
	data.Filters
}

type IdeaListResponse struct {
	Ideas []*data.Idea `json:"ideas"`
	Metadata data.Metadata `json:"metadata"`
}

type IdeaUpdateRequest struct {
	Text   *string `json:"name,omitempty"`
	TypeID *int64  `json:"type_id,omitempty"`
}

type UserIdeaRequest struct {
	UserID int64 `json:"user_id"`
	IdeaID int64 `json:"idea_id"`
}

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
}

type UserResponse struct {
	User User `json:"user"`
}

func NewUserResponse(user *data.User) UserResponse {
	return UserResponse{
		User: User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			Email: user.Email,
			Role: user.Role,
		},
	}
}

type UserCredentials struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type IdeaUserListRequest struct {
	UserID     int64 `json:"user_id"`
	Text       string `json:"text"`
	TypeID     int64 `json:"type_id"`
	ActiveOnly bool `json:"active_only"`
	Status     data.UserIdeaStatus `json:"status"`
	data.Filters
}

type BufferIdeaResponse struct {
	BufferIdea *data.BufferIdea `json:"buffer_idea"`	
}

type BufferIdeaListResponse struct {
	BufferIdeas []*data.BufferIdea `json:"buffer_ideas"`	
}

type BufferIdeaRequest struct {
	BufferIdeaID string `json:"buffer_idea_id"`
}
