package repository

import (
	"community-backend/internal/common"
	"community-backend/internal/model"

	"errors"
)

// GetCommentByID retrieves a single Comment by primary key.
func GetCommentByID(id int) (model.Comment, error) {
	var comment model.Comment

	err := common.DB.Get(&comment, `
		SELECT c.id, c.post_id, c.user_id, u.username, c.reply_id,
		       ru.username AS reply_username, c.content, c.created_time
		FROM comments c
		JOIN users u ON u.id = c.user_id
		LEFT JOIN users ru ON ru.id = c.reply_id
		WHERE c.id = $1
	`, id)

	return comment, err
}

// CreateComment inserts a new Comment and returns the created record.
func CreateComment(comment model.Comment) (model.Comment, error) {
	err := common.DB.QueryRowx(`
		INSERT INTO comments (post_id, user_id, reply_id, content, created_time)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, post_id, user_id, reply_id, content, created_time
	`, comment.PostID, comment.UserID, comment.ReplyID, comment.Content).
		StructScan(&comment)

	// Populate username with a follow-up query
	var username string
	_ = common.DB.Get(&username, `SELECT username FROM users WHERE id = $1`, comment.UserID)
	comment.Username = username

	// Populate reply_username if this is a reply
	if comment.ReplyID != nil {
		var replyUsername string
		_ = common.DB.Get(&replyUsername, `SELECT u.username FROM comments c JOIN users u ON u.id = c.user_id WHERE c.id = $1`, *comment.ReplyID)
		comment.ReplyUsername = &replyUsername
	}

	return comment, err
}

// GetCommentsByPostID retrieves all comments for a given post, ordered by creation time.
func GetCommentsByPostID(postID int) ([]model.Comment, error) {
	var comments []model.Comment

	err := common.DB.Select(&comments, `
		SELECT c.id, c.post_id, c.user_id, u.username, c.reply_id,
		       ru.username AS reply_username, c.content, c.created_time
		FROM comments c
		JOIN users u ON u.id = c.user_id
		LEFT JOIN users ru ON ru.id = c.reply_id
		WHERE c.post_id = $1
		ORDER BY c.created_time ASC
	`, postID)

	return comments, err
}

// DeleteCommentByID deletes a comment by its ID. If the comment has replies, they will also be deleted.
func DeleteCommentByID(comment model.Comment) error {
	result, err := common.DB.Exec(`
		DELETE
		FROM comments
		WHERE id = $1 AND user_id = $2
		OR reply_id = $1
	`, comment.ID, comment.UserID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("comment not found")
	}
	return nil
}

// GetCommentCountByPostID returns the total number of comments for a given post.
func GetCommentCountByPostID(postID int) (int, error) {
	var count int
	err := common.DB.Get(&count, `
		SELECT COUNT(1)
		FROM comments
		WHERE post_id = $1
	`, postID)
	return count, err
}
