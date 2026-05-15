package repository

import (
	"community-backend/internal/common"
	"community-backend/internal/model"

	"errors"
)

// GetPostByID retrieves a single post by primary key.
func GetPostByID(id int) (model.Post, error) {
	var post model.Post

	err := common.DB.Get(&post, `
		SELECT id, user_id, title, content, created_time, updated_time
		FROM posts 
		WHERE id = $1
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
func GetRecentPosts() ([]model.Posts, error) {
	var posts []model.Posts

	err := common.DB.Select(&posts, `
		SELECT p.id, p.user_id, u.username, p.title, p.created_time
		FROM posts as p
		LEFT JOIN users as u
		ON p.user_id = u.id
		ORDER BY created_time DESC
		LIMIT 5
	`)

	return posts, err
}

// GetPostsByUserID returns all posts authored by the given user, newest first.
func GetPostsByUserID(userID int) ([]model.Posts, error) {
	var posts []model.Posts

	err := common.DB.Select(&posts, `
		SELECT p.id, p.user_id, u.username, p.title, p.created_time
		FROM posts as p
		LEFT JOIN users as u
		ON p.user_id = u.id
		WHERE p.user_id = $1
		ORDER BY created_time DESC
	`, userID)

	return posts, err
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