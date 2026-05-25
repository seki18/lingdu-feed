package repository

import (
	"fmt"

	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/model"

	"errors"
)

// GetFavoritesByUserID retrieves all Favorites by user ID, with pagination.
func GetFavoritesByUserID(userID int, page, pageSize int) ([]model.Favorite, int, error) {
	var favorites []model.Favorite
	var total int

	if err := common.DB.Get(&total, `SELECT COUNT(1) FROM favorites WHERE user_id = $1`, userID); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := common.DB.Select(&favorites, `
		SELECT f.user_id, f.post_id, f.created_time
		FROM favorites f
		WHERE f.user_id = $1
		ORDER BY f.created_time DESC
		LIMIT $2 OFFSET $3
	`, userID, pageSize, offset)

	return favorites, total, err
}

// CreateFavorite inserts a new Favorite and returns the created record.
func CreateFavorite(favorite model.Favorite) (model.Favorite, error) {
	err := common.DB.QueryRowx(`
		INSERT INTO favorites (post_id, user_id, created_time)
		VALUES ($1, $2, NOW())
		RETURNING user_id, post_id, created_time
	`, favorite.PostID, favorite.UserID).
		StructScan(&favorite)

	return favorite, err
}

// DeleteFavorite delete favorite by primary key.
func DeleteFavorite(favorite model.Favorite) error {
	result, err := common.DB.Exec(`
		DELETE
		FROM favorites
		WHERE user_id = $1
		AND post_id = $2
	`, favorite.UserID, favorite.PostID)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("favorite not found, user_id: " + fmt.Sprint(favorite.UserID) + ", post_id: " + fmt.Sprint(favorite.PostID))
	}

	return nil
}

// CheckFavorited returns whether the given user has favorited the given post.
func CheckFavorited(userID, postID int) (bool, error) {
	var exists bool
	err := common.DB.Get(&exists, `
		SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id = $1 AND post_id = $2)
	`, userID, postID)
	return exists, err
}

// GetFavoritePostIDs returns post IDs that the user has favorited.
func GetFavoritePostIDs(userID, page, pageSize int) ([]int, int, error) {
	var total int
	if err := common.DB.Get(&total, `SELECT COUNT(1) FROM favorites WHERE user_id = $1`, userID); err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	var ids []int
	err := common.DB.Select(&ids, `
		SELECT post_id FROM favorites
		WHERE user_id = $1
		ORDER BY created_time DESC LIMIT $2 OFFSET $3
	`, userID, pageSize, offset)
	return ids, total, err
}
