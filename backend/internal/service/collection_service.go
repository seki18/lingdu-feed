package service

import (
	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/repository"
)

// GetCollectionByUserID retrieves Collections by user ID, with pagination.
func GetCollectionByUserID(userID int, page, pageSize int) ([]model.Collection, int, error) {
	return repository.GetCollectionByUserID(userID, page, pageSize)
}

// CreateCollection inserts a new Collection into the database.
func CreateCollection(req model.CreateCollectionRequest) (model.Collection, error) {
	praiseCreate := model.Collection{
		PostID: req.PostID,
		UserID: req.UserID,
	}

	return repository.CreateCollection(praiseCreate)
}

// DeleteCollection deletes a Collection by its user_id and post_id.
func DeleteCollection(req model.CreateCollectionRequest) error {
	praiseDelete := model.Collection{
		PostID: req.PostID,
		UserID: req.UserID,
	}
	return repository.DeleteCollection(praiseDelete)
}
