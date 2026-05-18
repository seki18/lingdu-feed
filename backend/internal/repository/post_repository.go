package repository

import (
	"community-backend/internal/common"
	"community-backend/internal/model"
	"errors"

	"github.com/jmoiron/sqlx"
)

// GetPostByID retrieves a single post by primary key, including author username.
func GetPostByID(id int) (model.Post, error) {
	var post model.Post

	err := common.DB.Get(&post, `
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_time, p.updated_time
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.id = $1
	`, id)

	return post, err
}

// CreatePost inserts a new post and returns the created record.
func CreatePost(post model.Post) (model.Post, error) {
	err := common.DB.QueryRowx(`
		INSERT INTO posts (user_id, title, content, created_time, updated_time)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, user_id, title, content, created_time, updated_time
	`, post.UserID, post.Title, post.Content).
		StructScan(&post)

	return post, err
}

// UpdatePost updates an existing post and returns the updated record.
func UpdatePost(post model.Post) (model.Post, error) {
	err := common.DB.QueryRowx(`
		UPDATE posts
		SET title = $3, content = $4, updated_time = now()
		WHERE id = $1 and user_id = $2
		RETURNING id, user_id, title, content, created_time, updated_time
	`, post.ID, post.UserID, post.Title, post.Content).
		StructScan(&post)

	return post, err
}

// GetRecentPosts returns the most recent posts with author usernames.
func GetRecentPosts(count int, excludeIDs []int, userID int) ([]model.Posts, error) {
	var posts []model.Posts

	query := `
		SELECT p.id, p.user_id, u.username, p.title, COALESCE(s.status, 0) as status, p.created_time
		FROM posts as p
		LEFT JOIN users as u
		ON p.user_id = u.id
		LEFT JOIN interaction_status as s
		ON p.id = s.post_id
		WHERE s.status IS NULL OR s.status <= ?`
	args := []any{model.FeedDisplay}

	if len(excludeIDs) > 0 {
		query += ` AND p.id NOT IN (?)`
		args = append(args, excludeIDs)
	}

	if userID != -1 {
		query += ` AND s.user_id = ?`
		args = append(args, userID)
	}

	query += `
		ORDER BY created_time DESC
		LIMIT ?
	`
	args = append(args, count+3)

	query, args, err := sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}
	query = common.DB.Rebind(query)

	err = common.DB.Select(&posts, query, args...)
	return posts, err
}

// GetPostsByUserID returns posts authored by the given user, newest first, with pagination.
func GetPostsByUserID(userID int, page, pageSize int) ([]model.Posts, int, error) {
	var posts []model.Posts
	var total int

	// Count total
	if err := common.DB.Get(&total, `SELECT COUNT(1) FROM posts WHERE user_id = $1`, userID); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := common.DB.Select(&posts, `
		SELECT p.id, p.user_id, u.username, p.title, p.created_time
		FROM posts as p
		LEFT JOIN users as u
		ON p.user_id = u.id
		WHERE p.user_id = $1
		ORDER BY created_time DESC
		LIMIT $2 OFFSET $3
	`, userID, pageSize, offset)

	return posts, total, err
}

// DeletePostByID delete a single post by primary key.
func DeletePostByID(id int64) error {
	result, err := common.DB.Exec(`
		DELETE
		FROM posts
		WHERE id = $1
	`, id)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("post not found")
	}

	return nil
}
