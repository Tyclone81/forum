package models

import "time"

// Post represents a thread created by a authenticated architect user.
type Post struct {
	ID         string    `json:"id"`       // UUID
	UserID     string    `json:"user_id"`  // Author foreign key
	Username   string    `json:"username"` // Joined from users table for simple UI rendering
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Categories []string  `json:"categories"` // Fulfills "associate one or more categories to a post"
	Likes      int       `json:"likes"`      // Computed metric visible to all users
	Dislikes   int       `json:"dislikes"`   // Computed metric visible to all users
	CreatedAt  time.Time `json:"created_at"`
}

// Category represents a distinct system subforum topic tag.
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"` // e.g., "Golang", "Docker", "AI"
}
