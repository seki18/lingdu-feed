package service

import (
	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/repository"
)

// ── Like ──

// CreateLike inserts a new Like into the database.
func CreateLike(req model.LikeRequest) (model.Like, error) {
	likeCreate := model.Like{
		PostID: req.PostID,
		UserID: req.UserID,
	}
	return repository.CreateLike(likeCreate)
}

// DeleteLike deletes a Like by its user_id and post_id.
func DeleteLike(req model.LikeRequest) error {
	likeDelete := model.Like{
		PostID: req.PostID,
		UserID: req.UserID,
	}
	return repository.DeleteLike(likeDelete)
}

// ── Favorite ──

// GetFavoritesByUserID retrieves Favorites by user ID, with pagination.
func GetFavoritesByUserID(userID int, page, pageSize int) ([]model.Favorite, int, error) {
	return repository.GetFavoritesByUserID(userID, page, pageSize)
}

// CreateFavorite inserts a new Favorite into the database.
func CreateFavorite(req model.FavoriteRequest) (model.Favorite, error) {
	favoriteCreate := model.Favorite{
		PostID: req.PostID,
		UserID: req.UserID,
	}
	return repository.CreateFavorite(favoriteCreate)
}

// DeleteFavorite deletes a Favorite by its user_id and post_id.
func DeleteFavorite(req model.FavoriteRequest) error {
	favoriteDelete := model.Favorite{
		PostID: req.PostID,
		UserID: req.UserID,
	}
	return repository.DeleteFavorite(favoriteDelete)
}

// ── Comment ──

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
