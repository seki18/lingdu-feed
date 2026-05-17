package service

import (
	"community-backend/internal/model"
	"community-backend/internal/repository"
)

// GetCollectionByUserID retrieves Collections by user ID, with pagination.
func GetCollectionByUserID(userID int, page, pageSize int) ([]model.Collection, int, error) {
	return repository.GetCollectionByUserID(userID, page, pageSize)
}

// IsCollectionExist checks if a Collection already exists for the given user_id and post_id.
func IsCollectionExist(req model.CreateCollectionRequest) (bool, error) {
	Collection := model.Collection{
		PostID: req.PostID,
		UserID: req.UserID,
	}
	return repository.IsCollectionExist(Collection)
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

// GetCollectionCountByPostID returns the total number of praises for a given post.
func GetCollectionCountByPostID(postID int) (int, error) {
	return repository.GetCollectionCountByPostID(postID)
}
