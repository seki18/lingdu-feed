package service

import (
	"community-backend/internal/model"
	"community-backend/internal/repository"
)

// IsPraiseExist checks if a Praise already exists for the given user_id and post_id.
func IsPraiseExist(req model.CreatePraiseRequest) (bool, error) {
	Praise := model.Praise{
		PostID:  req.PostID,
		UserID:  req.UserID,
	}
	return repository.IsPraiseExist(Praise)
}

// CreatePraise inserts a new Praise into the database.
func CreatePraise(req model.CreatePraiseRequest) (model.Praise, error) {
	praiseCreate := model.Praise{
		PostID:  req.PostID,
		UserID:  req.UserID,
	}

	return repository.CreatePraise(praiseCreate)
}

// DeletePraise deletes a Praise by its user_id and post_id.
func DeletePraise(req model.CreatePraiseRequest) error {
	praiseDelete := model.Praise{
		PostID:  req.PostID,
		UserID:  req.UserID,
	}
	return repository.DeletePraise(praiseDelete)
}

// GetPraiseCountByPostID returns the total number of praises for a given post.
func GetPraiseCountByPostID(postID int) (int, error) {
	return repository.GetPraiseCountByPostID(postID)
}