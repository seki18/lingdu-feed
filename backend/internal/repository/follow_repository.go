package repository

import (
	"fmt"

	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/model"

	"errors"
)

func IsFollowExist(follow model.Follow) (bool, error) {
	var exists bool
	err := common.DB.Get(&exists, `
		SELECT EXISTS (
			SELECT 1
			FROM follows
			WHERE follower_id = $1 AND following_id = $2
		)
	`, follow.FollowerID, follow.FollowingID)
	return exists, err
}

// CreateFollow inserts a new Follow and returns the created record.
func CreateFollow(follow model.Follow) (model.Follow, error) {
	err := common.DB.QueryRowx(`
		INSERT INTO follows (follower_id, following_id, created_time)
		VALUES ($1, $2, NOW())
		RETURNING follower_id, following_id, created_time
	`, follow.FollowerID, follow.FollowingID).
		StructScan(&follow)

	return follow, err
}

// DeleteFollow delete follow by primary key.
func DeleteFollow(follow model.Follow) error {
	result, err := common.DB.Exec(`
		DELETE
		FROM follows
		WHERE follower_id = $1
		AND following_id = $2
	`, follow.FollowerID, follow.FollowingID)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("follow not found, follower_id: " + fmt.Sprint(follow.FollowerID) + ", following_id: " + fmt.Sprint(follow.FollowingID))
	}

	return nil
}

// IncrFollowingCount atomically increments the following_count for a user.
func IncrFollowingCount(userID int) error {
	_, err := common.DB.Exec(`UPDATE users SET following_count = following_count + 1 WHERE id = $1`, userID)
	return err
}

// DecrFollowingCount atomically decrements the following_count for a user (floor 0).
func DecrFollowingCount(userID int) error {
	_, err := common.DB.Exec(`UPDATE users SET following_count = GREATEST(following_count - 1, 0) WHERE id = $1`, userID)
	return err
}

// IncrFollowerCount atomically increments the follower_count for a user.
func IncrFollowerCount(userID int) error {
	_, err := common.DB.Exec(`UPDATE users SET follower_count = follower_count + 1 WHERE id = $1`, userID)
	return err
}

// DecrFollowerCount atomically decrements the follower_count for a user (floor 0).
func DecrFollowerCount(userID int) error {
	_, err := common.DB.Exec(`UPDATE users SET follower_count = GREATEST(follower_count - 1, 0) WHERE id = $1`, userID)
	return err
}

// GetFollowingListByFollowerID returns a paginated list of users that a given user is following.
func GetFollowingListByFollowerID(followerID int, page, pageSize int) ([]model.Follow, int, error) {
	var follows []model.Follow
	var total int

	if err := common.DB.Get(&total, `SELECT COUNT(1) FROM follows WHERE follower_id = $1`, followerID); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := common.DB.Select(&follows, `
		SELECT c.follower_id, c.following_id, u.username, c.created_time
		FROM follows c
		JOIN users u ON u.id = c.following_id
		WHERE c.follower_id = $1
		ORDER BY c.created_time DESC
		LIMIT $2 OFFSET $3
	`, followerID, pageSize, offset)

	return follows, total, err
}

// GetFollowerListByFollowingID returns a paginated list of followers for a given user, including usernames.
func GetFollowerListByFollowingID(followingID int, page, pageSize int) ([]model.Follow, int, error) {
	var follows []model.Follow
	var total int

	if err := common.DB.Get(&total, `SELECT COUNT(1) FROM follows WHERE following_id = $1`, followingID); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := common.DB.Select(&follows, `
		SELECT c.follower_id, c.following_id, u.username, c.created_time
		FROM follows c
		JOIN users u ON u.id = c.follower_id
		WHERE c.following_id = $1
		ORDER BY c.created_time DESC
		LIMIT $2 OFFSET $3
	`, followingID, pageSize, offset)

	return follows, total, err
}
