package models

import "time"

type NewsFullDetailed struct {
	ID       int       `json:"id"` // Было string, стало int
	Title    string    `json:"title"`
	Content  string    `json:"content"` // В GONEWS нет Author/CreatedAt, есть PubTime/Link/Source
	PubTime  time.Time `json:"pub_time"`
	Link     string    `json:"link"`
	Source   string    `json:"source"`
	Comments []Comment `json:"comments,omitempty"`
}

type NewsShortDetailed struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	PubTime time.Time `json:"pub_time"`
	Link    string    `json:"link"`
	Source  string    `json:"source"`
}
