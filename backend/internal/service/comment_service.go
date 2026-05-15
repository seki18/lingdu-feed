package service

import (
	"community-backend/internal/model"
	"community-backend/internal/repository"
)

// GetCommentByID retrieves a single Comment by its ID.
func GetCommentByID(id int) (model.Comment, error) {
	return repository.GetCommentByID(id)
}

// CreateComment inserts a new Comment into the database.
func CreateComment(req model.CreateCommentRequest) (model.Comment, error) {
	commentCreate := model.Comment{
		PostID:  req.PostID,
		UserID:  req.UserID,
		Content: req.Content,
	}
	if req.ReplyID != nil {
		rid := *req.ReplyID
		commentCreate.ReplyID = &rid
	}

	return repository.CreateComment(commentCreate)
}

// GetCommentsByPostID retrieves all comments for a given post.
func GetCommentsByPostID(postID int) ([]model.Comment, error) {
	return repository.GetCommentsByPostID(postID)
}
