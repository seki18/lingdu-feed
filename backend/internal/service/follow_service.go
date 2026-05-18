package service

import (
	"community-backend/internal/model"
	"community-backend/internal/repository"
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

// GetFollowingCountByFollowerID returns the total number of users that a given user is following.
func GetFollowingCountByFollowerID(followerID int) (int, error) {
	return repository.GetFollowingCountByFollowerID(followerID)
}

// GetFollowerCountByFollowingID returns the total number of followers for a given user.
func GetFollowerCountByFollowingID(followingID int) (int, error) {
	return repository.GetFollowerCountByFollowingID(followingID)
}

// GetFollowingListByFollowerID returns a paginated list of users that a given user is following.
func GetFollowingListByFollowerID(followerID int, page, pageSize int) ([]model.Follow, int, error) {
	return repository.GetFollowingListByFollowerID(followerID, page, pageSize)
}

// GetFollowerListByFollowingID returns a paginated list of followers for a given user.
func GetFollowerListByFollowingID(followingID int, page, pageSize int) ([]model.Follow, int, error) {
	return repository.GetFollowerListByFollowingID(followingID, page, pageSize)
}
