package service

import (
	"github.com/seki18/lingdu-feed/internal/cache"
	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/repository"
)

// ── Like ──

// CreateLike inserts a new Like into the database and increments stats.
func CreateLike(req model.LikeRequest) (model.Like, error) {
	likeCreate := model.Like{
		PostID: req.PostID,
		UserID: req.UserID,
	}
	result, err := repository.CreateLike(likeCreate)
	if err != nil {
		return result, err
	}
	_ = cache.IncrLikeCount(req.PostID)
	return result, nil
}

// DeleteLike deletes a Like by its user_id and post_id and decrements stats.
func DeleteLike(req model.LikeRequest) error {
	likeDelete := model.Like{
		PostID: req.PostID,
		UserID: req.UserID,
	}
	if err := repository.DeleteLike(likeDelete); err != nil {
		return err
	}
	_ = cache.DecrLikeCount(req.PostID)
	return nil
}

// ── Favorite ──

// GetFavoritesByUserID retrieves Favorites by user ID, with pagination.
func GetFavoritesByUserID(userID int, page, pageSize int) ([]model.Favorite, int, error) {
	return repository.GetFavoritesByUserID(userID, page, pageSize)
}

// CreateFavorite inserts a new Favorite into the database and increments stats.
func CreateFavorite(req model.FavoriteRequest) (model.Favorite, error) {
	favoriteCreate := model.Favorite{
		PostID: req.PostID,
		UserID: req.UserID,
	}
	result, err := repository.CreateFavorite(favoriteCreate)
	if err != nil {
		return result, err
	}
	_ = cache.IncrFavoriteCount(req.PostID)
	return result, nil
}

// DeleteFavorite deletes a Favorite by its user_id and post_id and decrements stats.
func DeleteFavorite(req model.FavoriteRequest) error {
	favoriteDelete := model.Favorite{
		PostID: req.PostID,
		UserID: req.UserID,
	}
	if err := repository.DeleteFavorite(favoriteDelete); err != nil {
		return err
	}
	_ = cache.DecrFavoriteCount(req.PostID)
	return nil
}

// ── Comment ──

// CreateComment inserts a new Comment into the database and increments stats.
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
	result, err := repository.CreateComment(commentCreate)
	if err != nil {
		return result, err
	}
	_ = cache.IncrCommentCount(req.PostID)
	return result, nil
}

// DeleteCommentByID deletes a comment by its ID and decrements stats.
func DeleteCommentByID(req model.DeleteCommentRequest) error {
	commentDelete := model.Comment{
		ID:     req.PostID,
		UserID: req.UserID,
	}
	if err := repository.DeleteCommentByID(commentDelete); err != nil {
		return err
	}
	_ = cache.DecrCommentCount(req.PostID)
	return nil
}
