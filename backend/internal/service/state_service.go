package service

import (
	"github.com/seki18/lingdu-feed/internal/cache"
	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/repository"
)

// GetState retrieves a single State by post/user.
func GetState(req model.StateRequest) (model.State, error) {
	return repository.GetState(model.State{
		PostID: req.PostID,
		UserID: req.UserID,
	})
}

// UpsertState inserts a new State into the database.
func UpsertState(req model.StateRequest) error {
	stateCreate := model.State{
		PostID: req.PostID,
		UserID: req.UserID,
		Status: req.Status,
	}
	return repository.UpsertState(stateCreate)
}

// BatchUpsertState inserts or updates multiple State records.
func BatchUpsertState(reqs []model.StateRequest, userID int) error {
	for i := range reqs {
		reqs[i].UserID = userID
		if err := UpsertState(reqs[i]); err != nil {
			return err
		}
	}
	return nil
}

// GetViewCountByPostID returns the total number of views (clicks) for a given post.
func GetViewCountByPostID(postID int) (int, error) {
	return repository.GetViewCountByPostID(postID)
}

// IncrExposeCount increments the expose_count for a post via cache.
func IncrExposeCount(postID int) error {
	return cache.IncrExposeCount(postID)
}

// IncrViewCount increments the view_count for a post via cache.
func IncrViewCount(postID int) error {
	return cache.IncrViewCount(postID)
}
