package service

import (
	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/repository"
)

// CreatePraise inserts a new Praise into the database.
func CreatePraise(req model.CreatePraiseRequest) (model.Praise, error) {
	praiseCreate := model.Praise{
		PostID: req.PostID,
		UserID: req.UserID,
	}

	return repository.CreatePraise(praiseCreate)
}

// DeletePraise deletes a Praise by its user_id and post_id.
func DeletePraise(req model.CreatePraiseRequest) error {
	praiseDelete := model.Praise{
		PostID: req.PostID,
		UserID: req.UserID,
	}
	return repository.DeletePraise(praiseDelete)
}
