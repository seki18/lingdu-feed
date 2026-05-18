package service

import (
	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/repository"
)

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

// DeleteCommentByID deletes a comment by its ID. If the comment has replies, they will also be deleted.
func DeleteCommentByID(req model.DeleteCommentRequest) error {
	commentDelete := model.Comment{
		ID:     req.PostID,
		UserID: req.UserID,
	}
	return repository.DeleteCommentByID(commentDelete)
}
