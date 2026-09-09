package models

import "time"

// Session holds authentication validation values for active user tokens.
type Session struct {
	ID        string    `json:"id"` // UUID Session Cookie Value
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"` // Validated expiration date
}
