package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/model"
)

// GetStatsByPostIDs returns stats for the given post IDs in the same order.
func GetStatsByPostIDs(ids []int) ([]model.PostStats, error) {
	if len(ids) == 0 {
		return []model.PostStats{}, nil
	}

	query, args, err := sqlx.In(`
		SELECT id, like_count, comment_count, favorite_count,
			view_count, expose_count, score, updated_time
		FROM post_stats
		WHERE id IN (?)
	`, ids)
	if err != nil {
		return nil, err
	}
	query = common.DB.Rebind(query)

	var rows []model.PostStats
	if err := common.DB.Select(&rows, query, args...); err != nil {
		return nil, err
	}

	// Re-sort to match input ID order (deduplicated)
	rowMap := make(map[int]model.PostStats, len(rows))
	for _, row := range rows {
		rowMap[row.ID] = row
	}
	seen := make(map[int]bool, len(ids))
	result := make([]model.PostStats, 0, len(ids))
	for _, id := range ids {
		if seen[id] {
			continue
		}
		seen[id] = true
		if row, ok := rowMap[id]; ok {
			result = append(result, row)
		}
	}
	return result, nil
}

// GetStatsByPostID returns stats for a single post.
func GetStatsByPostID(id int) (*model.PostStats, error) {
	var stats model.PostStats
	err := common.DB.Get(&stats, `
		SELECT id, like_count, comment_count, favorite_count,
			view_count, expose_count, score, updated_time
		FROM post_stats
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// CreateStats inserts a default row into post_stats after post creation.
func CreateStats(id int) error {
	_, err := common.DB.Exec(`INSERT INTO post_stats (id) VALUES ($1)`, id)
	return err
}

// GetExistingPostIDsByStats returns the subset of postIDs that actually exist in the posts table.
// Used to filter out orphan stats cache entries before upsert.
func GetExistingPostIDsByStats(ids []int) ([]int, error) {
	if len(ids) == 0 {
		return []int{}, nil
	}
	query, args, err := sqlx.In(`SELECT id FROM posts WHERE id IN (?)`, ids)
	if err != nil {
		return nil, err
	}
	query = common.DB.Rebind(query)
	var existing []int
	if err := common.DB.Select(&existing, query, args...); err != nil {
		return nil, err
	}
	return existing, nil
}

// BatchUpsertStats bulk upserts stats rows (used by periodic cache sync).
func BatchUpsertStats(stats []model.PostStats) error {
	if len(stats) == 0 {
		return nil
	}

	// Build batched INSERT ... ON CONFLICT DO UPDATE
	query := `
		INSERT INTO post_stats (id, like_count, comment_count, favorite_count,
			view_count, expose_count, score, updated_time)
		VALUES (:id, :like_count, :comment_count, :favorite_count,
			:view_count, :expose_count, :score, NOW())
		ON CONFLICT (id) DO UPDATE SET
			like_count     = EXCLUDED.like_count,
			comment_count  = EXCLUDED.comment_count,
			favorite_count = EXCLUDED.favorite_count,
			view_count     = EXCLUDED.view_count,
			expose_count   = EXCLUDED.expose_count,
			score          = EXCLUDED.score,
			updated_time   = NOW()
	`

	tx, err := common.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, s := range stats {
		if _, err := tx.NamedExec(query, s); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetRecommendPostIDs returns post IDs from post_stats ranked by score DESC.
// cursorID is the id of the last item from the previous page; pass 0 for first page.
func GetRecommendPostIDs(count int, cursorID int) ([]int, error) {
	query := `
		SELECT s.id FROM post_stats s
		WHERE ($1 = 0 OR s.id < $1)
		ORDER BY s.score DESC, s.id DESC
		LIMIT $2
	`
	var ids []int
	if err := common.DB.Select(&ids, query, cursorID, count); err != nil {
		return nil, err
	}
	return ids, nil
}
