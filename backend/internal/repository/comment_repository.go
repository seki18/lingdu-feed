package repository

import (
	"community-backend/internal/common"
	"community-backend/internal/model"
)

// GetCommentByID retrieves a single Comment by primary key.
func GetCommentByID(id int) (model.Comment, error) {
	var Comment model.Comment

	err := common.DB.Get(&Comment, `
		SELECT c.id, c.post_id, c.user_id, u.username, c.reply_id,
		       ru.username AS reply_username, c.content, c.created_time
		FROM Comments c
		JOIN users u ON u.id = c.user_id
		LEFT JOIN users ru ON ru.id = c.reply_id
		WHERE c.id = $1
	`, id)

	return Comment, err
}

// CreateComment inserts a new Comment and returns the created record.
func CreateComment(Comment model.Comment) (model.Comment, error) {
	err := common.DB.QueryRowx(`
		INSERT INTO Comments (post_id, user_id, reply_id, content, created_time)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, post_id, user_id, reply_id, content, created_time
	`, Comment.PostID, Comment.UserID, Comment.ReplyID, Comment.Content).
		StructScan(&Comment)

	// Populate username with a follow-up query
	var username string
	_ = common.DB.Get(&username, `SELECT username FROM users WHERE id = $1`, Comment.UserID)
	Comment.Username = username

	return Comment, err
}

// GetCommentsByPostID retrieves all comments for a given post, ordered by creation time.
func GetCommentsByPostID(postID int) ([]model.Comment, error) {
	var comments []model.Comment

	err := common.DB.Select(&comments, `
		SELECT c.id, c.post_id, c.user_id, u.username, c.reply_id,
		       ru.username AS reply_username, c.content, c.created_time
		FROM Comments c
		JOIN users u ON u.id = c.user_id
		LEFT JOIN users ru ON ru.id = c.reply_id
		WHERE c.post_id = $1
		ORDER BY c.created_time ASC
	`, postID)

	return comments, err
}
