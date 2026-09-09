package models

// Interaction models a unique upvote/downvote entity entry.
// Guarantees strict constraints: Value is 1 (like) or -1 (dislike)
type Interaction struct {
	ID         string `json:"id"`          // UUID
	UserID     string `json:"user_id"`     // Voter identity context
	TargetID   string `json:"target_id"`   // References Post ID or Comment ID
	TargetType string `json:"target_type"` // Evaluates strictly to "post" or "comment"
	Value      int    `json:"value"`       // 1 or -1 interaction value
}
