package service

import (
	"community-backend/internal/model"
	"community-backend/internal/repository"
)

// GetInteractionStatus retrieves a single InteractionStatus by its ID.
func GetInteractionStatus(req model.CreateInteractionStatusRequest) (model.InteractionStatus, error) {
	interactionStatusCreate := model.InteractionStatus{
		PostID:  req.PostID,
		UserID:  req.UserID,
		Status:  req.Status,
	}
	return repository.GetInteractionStatus(interactionStatusCreate)
}

// GetInteractionStatusByUserID retrieves InteractionStatus by user ID.
func GetInteractionStatusByUserID(id int) ([]model.InteractionStatus, error) {
	return repository.GetInteractionStatusByUserID(id)
}

// UpsertInteractionStatus inserts a new InteractionStatus into the database.
func UpsertInteractionStatus(req model.CreateInteractionStatusRequest) error {
	interactionStatusCreate := model.InteractionStatus{
		PostID:  req.PostID,
		UserID:  req.UserID,
		Status:  req.Status,
	}
	return repository.UpsertInteractionStatus(interactionStatusCreate)
}