package repository

import (
	"fmt"
	"community-backend/internal/common"
	"community-backend/internal/model"

	"errors"
)

func IsPraiseExist(praise model.Praise) (bool, error) {
	var exists bool
	err := common.DB.Get(&exists, `
		SELECT EXISTS (
			SELECT 1
			FROM praises
			WHERE user_id = $1 AND post_id = $2
		)
	`, praise.UserID, praise.PostID)
	return exists, err
}

// CreatePraise inserts a new Praise and returns the created record.
func CreatePraise(praise model.Praise) (model.Praise, error) {
	err := common.DB.QueryRowx(`
		INSERT INTO praises (post_id, user_id, created_time)
		VALUES ($1, $2, NOW())
		RETURNING user_id, post_id, created_time
	`, praise.PostID, praise.UserID).
		StructScan(&praise)

	return praise, err
}

// DeletePraise delete praise by primary key.
func DeletePraise(praise model.Praise) error {
	result, err := common.DB.Exec(`
		DELETE
		FROM praises
		WHERE user_id = $1
		AND post_id = $2
	`, praise.UserID, praise.PostID)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("praise not found, user_id: " + fmt.Sprint(praise.UserID) + ", post_id: " + fmt.Sprint(praise.PostID))
	}

	return nil
}

// GetPraiseCountByPostID returns the total number of praises for a given post.
func GetPraiseCountByPostID(postID int) (int, error) {
	var count int
	err := common.DB.Get(&count, `
		SELECT COUNT(1)
		FROM praises
		WHERE post_id = $1
	`, postID)
	return count, err
}