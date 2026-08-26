package data

type Idea struct {
	ID     int64  `json:"id"`
	Text   string `json:"text"`
	Type   Type   `json:"type"`
	Prompt Prompt `json:"prompt"`
}
