package service

import (
	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/repository"
)

// GetInteractionStatus retrieves a single InteractionStatus by post/user.
func GetInteractionStatus(req model.CreateInteractionStatusRequest) (model.InteractionStatus, error) {
	return repository.GetInteractionStatus(model.InteractionStatus{
		PostID: req.PostID,
		UserID: req.UserID,
	})
}

// UpsertInteractionStatus inserts a new InteractionStatus into the database.
func UpsertInteractionStatus(req model.CreateInteractionStatusRequest) error {
	interactionStatusCreate := model.InteractionStatus{
		PostID: req.PostID,
		UserID: req.UserID,
		Status: req.Status,
	}
	return repository.UpsertInteractionStatus(interactionStatusCreate)
}

// BatchUpsertInteractionStatus inserts or updates multiple InteractionStatus records.
func BatchUpsertInteractionStatus(reqs []model.CreateInteractionStatusRequest, userID int) error {
	for i := range reqs {
		reqs[i].UserID = userID
		if err := UpsertInteractionStatus(reqs[i]); err != nil {
			return err
		}
	}
	return nil
}

// GetViewCountByPostID returns the total number of views (clicks) for a given post.
func GetViewCountByPostID(postID int) (int, error) {
	return repository.GetViewCountByPostID(postID)
}
