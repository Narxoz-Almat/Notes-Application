package models

import "time"

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int       `json:"user_id"`
	TagIDs    []int     `json:"tag_ids"`
	CreatedAt time.Time `json:"created_at"`
}
