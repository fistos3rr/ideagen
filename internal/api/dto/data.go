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
	Name string `json:"name"`
	ActiveOnly bool `json:"bool"`
	data.Filters
}

type TypeListResponse struct {
	Types []*data.Type `json:"types"`
	Metadata data.Metadata `json:"metadata"`
}

type TypeUpdateRequest struct {
	Name     *string `json:"name"`
	IsActive *bool  `json:"is_active"`
}
