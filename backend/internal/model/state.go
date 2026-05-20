package model

import "time"

// StateStatus represents the delivery/engagement status of a post for a user.
type StateStatus int

const (
	StateUnknown   StateStatus = iota // 0 - unknown
	StateDelivered                    // 1 - delivered (post shown in feed)
	StateExposed                      // 2 - exposed (post displayed on screen)
	StateClicked                      // 3 - clicked (post opened / detail viewed)
)

// State represents a row in the states table (formerly interaction_status).
type State struct {
	PostID      int         `db:"post_id" json:"post_id"`
	UserID      int         `db:"user_id" json:"user_id"`
	Status      StateStatus `db:"status" json:"status"`
	UpdatedTime time.Time   `db:"updated_time" json:"updated_time"`
}

// StateRequest is the JSON body for state reporting.
type StateRequest struct {
	PostID int         `db:"post_id" json:"post_id"`
	UserID int         `db:"user_id" json:"user_id"`
	Status StateStatus `db:"status" json:"status"`
}
