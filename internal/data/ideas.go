package data

import "time"

type Idea struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Types     []Type    `json:"types,omitempty"`
	Message   string    `json:"message"`
}