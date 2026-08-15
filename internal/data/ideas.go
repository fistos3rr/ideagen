package data

import "time"

type Idea struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Type      string    `json:"type"`
}
