package model

import "time"

type FeedStatus int

const (
	FeedUnknown FeedStatus = iota
	FeedDelivery
	FeedDisplay
	FeedClick
)

// Status represents a single Status row in the database.
type InteractionStatus struct {
	PostID        int       `db:"post_id" json:"post_id"`
	UserID        int       `db:"user_id" json:"user_id"`
	Status      FeedStatus    `db:"status" json:"status"`
	UpdatedTime   time.Time `db:"updated_time" json:"updated_time"`
}

// CreateInteractionStatusRequest is the JSON body for POST /InteractionStatuss		.
type CreateInteractionStatusRequest struct {
	PostID        int       `db:"post_id" json:"post_id"`
	UserID        int       `db:"user_id" json:"user_id"`
	Status      FeedStatus    `db:"status" json:"status"`
}
