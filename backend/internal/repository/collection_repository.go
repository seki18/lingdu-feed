package repository

import (
	"community-backend/internal/common"
	"community-backend/internal/model"
	"fmt"

	"errors"
)

// GetCollectionByUserID retrieves all Collections by user ID, with pagination.
func GetCollectionByUserID(userID int, page, pageSize int) ([]model.Collection, int, error) {
	var collections []model.Collection
	var total int

	if err := common.DB.Get(&total, `SELECT COUNT(1) FROM collections WHERE user_id = $1`, userID); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := common.DB.Select(&collections, `
		SELECT c.user_id, c.post_id, c.created_time
		FROM collections c
		WHERE c.user_id = $1
		ORDER BY c.created_time DESC
		LIMIT $2 OFFSET $3
	`, userID, pageSize, offset)

	return collections, total, err
}

func IsCollectionExist(collection model.Collection) (bool, error) {
	var exists bool
	err := common.DB.Get(&exists, `
		SELECT EXISTS (
			SELECT 1
			FROM collections
			WHERE user_id = $1 AND post_id = $2
		)
	`, collection.UserID, collection.PostID)
	return exists, err
}

// CreateCollection inserts a new Collection and returns the created record.
func CreateCollection(collection model.Collection) (model.Collection, error) {
	err := common.DB.QueryRowx(`
		INSERT INTO collections (post_id, user_id, created_time)
		VALUES ($1, $2, NOW())
		RETURNING user_id, post_id, created_time
	`, collection.PostID, collection.UserID).
		StructScan(&collection)

	return collection, err
}

// DeleteCollection delete collection by primary key.
func DeleteCollection(collection model.Collection) error {
	result, err := common.DB.Exec(`
		DELETE
		FROM collections
		WHERE user_id = $1
		AND post_id = $2
	`, collection.UserID, collection.PostID)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("collection not found, user_id: " + fmt.Sprint(collection.UserID) + ", post_id: " + fmt.Sprint(collection.PostID))
	}

	return nil
}

// GetCollectionCountByPostID returns the total number of collections for a given post.
func GetCollectionCountByPostID(postID int) (int, error) {
	var count int
	err := common.DB.Get(&count, `
		SELECT COUNT(1)
		FROM collections
		WHERE post_id = $1
	`, postID)
	return count, err
}
