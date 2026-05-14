package repository

import (
	"community-backend/internal/common"
	"community-backend/internal/model"
)

func GetPostByID(id int) (model.Post, error) {
	var post model.Post

	err := common.DB.Get(&post, `
		SELECT id, user_id, title, content, created_time, updated_time
		FROM posts 
		WHERE id = $1
	`, id)

	return post, err
}

func CreatePost(post model.Post) (model.Post, error) {
	err := common.DB.QueryRowx(`
		INSERT INTO posts (user_id, title, content, created_time, updated_time)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, user_id, title, content, created_time, updated_time
	`, post.UserID, post.Title, post.Content).
		StructScan(&post)

	return post, err
}

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
