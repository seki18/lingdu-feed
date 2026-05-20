package repository

import (
	"strconv"

	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/model"

	"errors"
)

// CreateLike inserts a new Like and returns the created record.
func CreateLike(like model.Like) (model.Like, error) {
	err := common.DB.QueryRowx(`
		INSERT INTO likes (post_id, user_id, created_time)
		VALUES ($1, $2, NOW())
		RETURNING user_id, post_id, created_time
	`, like.PostID, like.UserID).
		StructScan(&like)

	return like, err
}

// DeleteLike delete like by primary key.
func DeleteLike(like model.Like) error {
	result, err := common.DB.Exec(`
		DELETE
		FROM likes
		WHERE user_id = $1
		AND post_id = $2
	`, like.UserID, like.PostID)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("like not found, user_id: " + strconv.Itoa(like.UserID) + ", post_id: " + strconv.Itoa(like.PostID))
	}

	return nil
}

// CheckLiked returns whether the given user has liked the given post.
func CheckLiked(userID, postID int) (bool, error) {
	var exists bool
	err := common.DB.Get(&exists, `
		SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND post_id = $2)
	`, userID, postID)
	return exists, err
}
