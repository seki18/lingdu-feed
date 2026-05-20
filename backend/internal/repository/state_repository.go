package repository

import (
	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/model"

	"github.com/jmoiron/sqlx"
)

// GetState retrieves a single states row.
func GetState(state model.State) (model.State, error) {
	var s model.State
	err := common.DB.Get(&s, `
		SELECT post_id, user_id, status, updated_time
		FROM states
		WHERE post_id = $1 AND user_id = $2
	`, state.PostID, state.UserID)
	return s, err
}

// UpsertState inserts a new State if it doesn't exist, or updates it if it does.
func UpsertState(state model.State) error {

	_, err := common.DB.Exec(`
		INSERT INTO states (
			user_id,
			post_id,
			status,
			updated_time
		)
		VALUES (
			$1,
			$2,
			$3,
			NOW()
		)
		ON CONFLICT (user_id, post_id)
		DO UPDATE SET
			status = GREATEST(
				states.status,
				EXCLUDED.status
			),
			updated_time = NOW()
	`,
		state.UserID,
		state.PostID,
		state.Status,
	)

	return err
}

// GetViewCounts returns a map of post_id → view count for the given post IDs.
// Views are counted from states rows where status = 3 (StateClicked).
func GetViewCounts(postIDs []int) (map[int]int, error) {
	if len(postIDs) == 0 {
		return map[int]int{}, nil
	}
	type row struct {
		PostID int `db:"post_id"`
		Cnt    int `db:"cnt"`
	}
	var rows []row
	query, args, err := sqlx.In(`
		SELECT post_id, COUNT(*) AS cnt
		FROM states
		WHERE post_id IN (?) AND status = ?
		GROUP BY post_id
	`, postIDs, model.StateClicked)
	if err != nil {
		return nil, err
	}
	query = common.DB.Rebind(query)
	if err := common.DB.Select(&rows, query, args...); err != nil {
		return nil, err
	}
	result := make(map[int]int, len(rows))
	for _, r := range rows {
		result[r.PostID] = r.Cnt
	}
	return result, nil
}

// GetViewCountByPostID returns the total number of views (clicks) for a single post.
func GetViewCountByPostID(postID int) (int, error) {
	var count int
	err := common.DB.Get(&count, `
		SELECT COUNT(1)
		FROM states
		WHERE post_id = $1 AND status = $2
	`, postID, model.StateClicked)
	return count, err
}
