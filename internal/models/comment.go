package models

import "time"

// Comment represents a direct response to a parent post entity.
type Comment struct {
	ID        string    `json:"id"`      // UUID
	PostID    string    `json:"post_id"` // Attached thread context target
	UserID    string    `json:"user_id"` // Author foreign key
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	CreatedAt time.Time `json:"created_at"`
}
