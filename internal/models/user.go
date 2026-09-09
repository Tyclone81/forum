package models

import "time"

// User represents a registered forum member account.
type User struct {
	ID           string    `json:"id"`       // UUID
	Username     string    `json:"username"` // Unique public screen name
	Email        string    `json:"email"`    // Unique registration email address
	PasswordHash string    `json:"-"`        // Securely hashed using bcrypt (excluded from JSON outputs)
	CreatedAt    time.Time `json:"created_at"`
}
