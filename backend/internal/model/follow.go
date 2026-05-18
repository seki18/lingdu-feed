package model

import "time"

// Follow represents a single Follow row in the database.
type Follow struct {
	FollowerID  int       `db:"follower_id" json:"follower_id"`
	FollowingID int       `db:"following_id" json:"following_id"`
	Username    string    `db:"username" json:"username"`
	CreatedTime time.Time `db:"created_time" json:"created_time"`
}

// CreateFollowRequest is the JSON body for POST /Follows.
type CreateFollowRequest struct {
	FollowerID  int `json:"follower_id"`
	FollowingID int `json:"following_id"`
}
