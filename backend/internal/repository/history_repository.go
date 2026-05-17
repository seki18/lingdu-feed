package repository

import (
	"community-backend/internal/common"
	"community-backend/internal/model"
)

// GetHistoryPostsByUserID retrieves History Posts by user ID, with pagination.
func GetHistoryPostsByUserID(userId int, page, pageSize int) ([]model.Post, int, error) {
	var posts []model.Post
	var total int

	if err := common.DB.Get(&total, `
		SELECT COUNT(1)
		FROM interaction_status
		WHERE user_id = $1 AND status = 3
	`, userId); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := common.DB.Select(&posts, `
		SELECT p.id, p.user_id, p.title, p.created_time, p.updated_time
		FROM posts p
		JOIN interaction_status s ON p.id = s.post_id
		WHERE s.user_id = $1
		AND s.status = 3
		ORDER BY s.updated_time DESC
		LIMIT $2 OFFSET $3
	`, userId, pageSize, offset)

	return posts, total, err
}
