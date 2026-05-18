package service

import (
	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/repository"
)

// IsFollowExist checks if a Follow already exists for the given user_id and post_id.
func IsFollowExist(req model.CreateFollowRequest) (bool, error) {
	Follow := model.Follow{
		FollowerID:  req.FollowerID,
		FollowingID: req.FollowingID,
	}
	return repository.IsFollowExist(Follow)
}

// CreateFollow inserts a new Follow into the database.
func CreateFollow(req model.CreateFollowRequest) (model.Follow, error) {
	followCreate := model.Follow{
		FollowerID:  req.FollowerID,
		FollowingID: req.FollowingID,
	}

	return repository.CreateFollow(followCreate)
}

// DeleteFollow deletes a Follow by its user_id and post_id.
func DeleteFollow(req model.CreateFollowRequest) error {
	followDelete := model.Follow{
		FollowerID:  req.FollowerID,
		FollowingID: req.FollowingID,
	}
	return repository.DeleteFollow(followDelete)
}

// GetFollowingListByFollowerID returns a paginated list of users that a given user is following.
func GetFollowingListByFollowerID(followerID int, page, pageSize int) ([]model.Follow, int, error) {
	return repository.GetFollowingListByFollowerID(followerID, page, pageSize)
}

// GetFollowerListByFollowingID returns a paginated list of followers for a given user.
func GetFollowerListByFollowingID(followingID int, page, pageSize int) ([]model.Follow, int, error) {
	return repository.GetFollowerListByFollowingID(followingID, page, pageSize)
}
