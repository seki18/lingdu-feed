package repository

import (
	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/model"

	"github.com/jmoiron/sqlx"
)

// GetRecentPostIDs returns post IDs that should be shown to the user.
// When allowDegraded is true, skips the state filter to include already-seen posts.
func GetRecentPostIDs(count int, excludeIDs []int, userID int, allowDegraded bool) ([]int, error) {
	query := `
		SELECT p.id
		FROM posts as p`
	var args []any

	if userID != -1 {
		query += `
		LEFT JOIN states as s 
		ON p.id = s.post_id AND s.user_id = ?
		WHERE (s.status IS NULL OR s.status <= ?)`
		args = append(args, userID, model.StateDelivered)
	} else {
		query += ` WHERE 1=1`
	}

	if len(excludeIDs) > 0 {
		query += ` AND p.id NOT IN (?)`
		args = append(args, excludeIDs)
	}

	query += `
		ORDER BY p.created_time DESC
		LIMIT ?
	`
	args = append(args, count)

	query, args, err := sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}
	query = common.DB.Rebind(query)

	var ids []int
	if err := common.DB.Select(&ids, query, args...); err != nil {
		return nil, err
	}
	return ids, nil
}

// GetFollowingPostIDs returns post IDs from users that the given user follows.
// If checkStatus is true, applies the state filter to skip viewed posts.
func GetFollowingPostIDs(count int, excludeIDs []int, userID int, checkStatus bool) ([]int, error) {
	query := `
		SELECT p.id
		FROM posts as p
		JOIN follows as f ON p.user_id = f.following_id`
	var args []any

	if checkStatus {
		query += `
		LEFT JOIN states as s ON p.id = s.post_id AND s.user_id = ?
		WHERE f.follower_id = ?
		AND (s.status IS NULL OR s.status <= ?)`
		args = append(args, userID, userID, model.StateDelivered)
	} else {
		query += `
		WHERE f.follower_id = ?`
		args = append(args, userID)
	}

	if len(excludeIDs) > 0 {
		query += ` AND p.id NOT IN (?)`
		args = append(args, excludeIDs)
	}

	query += `
		ORDER BY p.created_time DESC
		LIMIT ?
	`
	args = append(args, count)

	query, args, err := sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}
	query = common.DB.Rebind(query)

	var ids []int
	if err := common.DB.Select(&ids, query, args...); err != nil {
		return nil, err
	}
	return ids, nil
}

// GetHistoryPostIDs returns post IDs that the user has viewed (clicked).
func GetHistoryPostIDs(userID, page, pageSize int) ([]int, int, error) {
	var total int
	if err := common.DB.Get(&total, `SELECT COUNT(1) FROM states WHERE user_id = $1 AND status = $2`, userID, model.StateClicked); err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	var ids []int
	err := common.DB.Select(&ids, `
		SELECT post_id FROM states
		WHERE user_id = $1 AND status = $2
		ORDER BY updated_time DESC LIMIT $3 OFFSET $4
	`, userID, model.StateClicked, pageSize, offset)
	return ids, total, err
}

// GetFavoritePostIDs returns post IDs that the user has favorited.
func GetFavoritePostIDs(userID, page, pageSize int) ([]int, int, error) {
	var total int
	if err := common.DB.Get(&total, `SELECT COUNT(1) FROM favorites WHERE user_id = $1`, userID); err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	var ids []int
	err := common.DB.Select(&ids, `
		SELECT post_id FROM favorites
		WHERE user_id = $1
		ORDER BY created_time DESC LIMIT $2 OFFSET $3
	`, userID, pageSize, offset)
	return ids, total, err
}

// GetRecommendPostIDs returns post IDs ranked by weighted score (recency * 0.1
// + views*3 + likes*5 + favorites*4 + comments*4). Returns count posts.
// When allowDegraded is true, skips the state filter to include already-seen posts.
func GetRecommendPostIDs(count int, excludeIDs []int, userID int, allowDegraded bool) ([]int, error) {
	query := `
		SELECT p.id FROM posts p`
	var args []any

	if userID != -1 && !allowDegraded {
		query += `
		LEFT JOIN states as s 
		ON p.id = s.post_id AND s.user_id = ?
		WHERE (s.status IS NULL OR s.status <= ?)`
		args = append(args, userID, model.StateDelivered)
	} else {
		query += ` WHERE 1=1`
	}

	if len(excludeIDs) > 0 {
		query += ` AND p.id NOT IN (?)`
		args = append(args, excludeIDs)
	}

	query += `
		ORDER BY (
			EXTRACT(EPOCH FROM p.created_time) * 0.1 +
			p.view_count * 3 +
			p.like_count * 5 +
			p.favorite_count * 4 +
			p.comment_count * 4
		) DESC
		LIMIT ?
	`
	args = append(args, count)

	query, args, err := sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}
	query = common.DB.Rebind(query)

	var ids []int
	if err := common.DB.Select(&ids, query, args...); err != nil {
		return nil, err
	}
	return ids, nil
}
