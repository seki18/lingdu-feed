package repository

import (
	"strconv"

	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/model"

	"errors"
)

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
		return errors.New("praise not found, user_id: " + strconv.Itoa(praise.UserID) + ", post_id: " + strconv.Itoa(praise.PostID))
	}

	return nil
}

// CheckPraised returns whether the given user has praised the given post.
func CheckPraised(userID, postID int) (bool, error) {
	var exists bool
	err := common.DB.Get(&exists, `
		SELECT EXISTS(SELECT 1 FROM praises WHERE user_id = $1 AND post_id = $2)
	`, userID, postID)
	return exists, err
}
