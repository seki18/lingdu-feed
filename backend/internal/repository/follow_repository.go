package repository

import (
	"community-backend/internal/common"
	"community-backend/internal/model"
	"fmt"

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

// GetFollowingCountByFollowerID returns the total number of users that a given user is following.
func GetFollowingCountByFollowerID(followerID int) (int, error) {
	var count int
	err := common.DB.Get(&count, `
		SELECT COUNT(1)
		FROM follows
		WHERE follower_id = $1
	`, followerID)
	return count, err
}

// GetFollowerCountByFollowingID returns the total number of followers for a given user.
func GetFollowerCountByFollowingID(followingID int) (int, error) {
	var count int
	err := common.DB.Get(&count, `
		SELECT COUNT(1)
		FROM follows
		WHERE following_id = $1
	`, followingID)
	return count, err
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
