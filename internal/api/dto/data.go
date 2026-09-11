package dto

import "github.com/fistos3rr/ideagen/internal/data"

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
