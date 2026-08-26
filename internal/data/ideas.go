package data

type Idea struct {
	ID int64 `json:"id"`
	Type Type `json:"type"`
	Prompt Prompt `json:"prompt"`
}
