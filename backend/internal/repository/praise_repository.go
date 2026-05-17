package repository

import (
	"community-backend/internal/common"
	"community-backend/internal/model"

	"errors"
)

// GetPraiseByID retrieves a single Praise by primary key.
func GetPraiseByID(id int) (model.Praise, error) {
	var Praise model.Praise

	err := common.DB.Get(&Praise, `
		SELECT user_id, post_id, created_time
		FROM praises
		WHERE id = $1
	`, id)

	return Praise, err
}

func IsPraiseExist(Praise model.Praise) (bool, error) {
	var exists bool
	err := common.DB.Get(&exists, `
		SELECT EXISTS (
			SELECT 1
			FROM praises
			WHERE user_id = $1 AND post_id = $2
		)
	`, Praise.UserID, Praise.PostID)
	return exists, err
}

// CreatePraise inserts a new Praise and returns the created record.
func CreatePraise(Praise model.Praise) (model.Praise, error) {
	err := common.DB.QueryRowx(`
		INSERT INTO praises (post_id, user_id, created_time)
		VALUES ($1, $2, NOW())
		RETURNING user_id, post_id, created_time
	`, Praise.PostID, Praise.UserID).
		StructScan(&Praise)

	return Praise, err
}

// DeletePraise delete praise by primary key.
func DeletePraise(Praise model.Praise) error {
	result, err := common.DB.Exec(`
		DELETE
		FROM praises
		WHERE user_id = $1
		AND post_id = $2
	`, Praise.UserID, Praise.PostID)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("praise not found")
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