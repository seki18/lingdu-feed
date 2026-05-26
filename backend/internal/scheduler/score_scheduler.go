package scheduler

import (
	"log"
	"time"

	"github.com/seki18/lingdu-feed/internal/cache"
	"github.com/seki18/lingdu-feed/internal/common"
)

// CalculateAndUpdateScores updates the score column for posts.
// Score is always in [0, 1], computed from recency, CTR, and engagement rates.
// When fullUpdate is true, updates ALL posts (used at startup).
// When fullUpdate is false, updates posts whose stats changed within 24h
// OR whose age is within the recency decay window (21 days, ~3 half-lives).
func CalculateAndUpdateScores(fullUpdate bool) {
	// Score formula (each term ∈ [0, 1], weights sum to 1.0):
	//   0.15 × recency_decay                            — EXP decay, 7-day half-life
	//   0.35 × tanh(view_count / 200)                   — absolute popularity
	//   0.20 × ctr (view_count / expose_count)          — click-through rate
	//   0.15 × tanh(like_count / 50)                    — absolute likes
	//   0.10 × tanh(comment_count / 30)                 — absolute comments
	//   0.05 × tanh(favorite_count / 30)                — absolute favorites
	//
	// tanh(x) = (e^2x - 1) / (e^2x + 1), smooth saturation to [0, 1)
	query := `
		UPDATE post_stats ps SET score = ROUND((
			0.15 * EXP(-EXTRACT(EPOCH FROM (NOW() - p.created_time)) / 604800.0) +
			0.35 * (
				(EXP(2.0 * ps.view_count / 200.0) - 1) /
				(EXP(2.0 * ps.view_count / 200.0) + 1)
			) +
			0.20 * COALESCE(ps.view_count::float / NULLIF(ps.expose_count, 0), 0) +
			0.15 * (
				(EXP(2.0 * ps.like_count / 50.0) - 1) /
				(EXP(2.0 * ps.like_count / 50.0) + 1)
			) +
			0.10 * (
				(EXP(2.0 * ps.comment_count / 30.0) - 1) /
				(EXP(2.0 * ps.comment_count / 30.0) + 1)
			) +
			0.05 * (
				(EXP(2.0 * ps.favorite_count / 30.0) - 1) /
				(EXP(2.0 * ps.favorite_count / 30.0) + 1)
			)
		)::numeric, 6)::double precision
		FROM posts p
		WHERE ps.id = p.id
	`
	if !fullUpdate {
		query += ` AND (
			ps.updated_time >= NOW() - INTERVAL '24 hours'
			OR p.created_time >= NOW() - INTERVAL '21 days'
		)`
	}

	result, err := common.DB.Exec(query)
	if err != nil {
		log.Printf("[ScoreScheduler] Failed to update scores: %v", err)
		return
	}
	rows, _ := result.RowsAffected()
	log.Printf("[ScoreScheduler] Score update completed, rows affected: %d, fullUpdate: %v", rows, fullUpdate)

	// Refresh ranking cache with new scores
	if err := cache.RefreshRanking(); err != nil {
		log.Printf("[ScoreScheduler] Failed to refresh ranking cache: %v", err)
	}
}

// RunScoreScheduler starts the score calculation loop.
// It runs immediately on startup (full update), then every 1 minute (incremental, 24h window).
func RunScoreScheduler() {
	log.Println("[ScoreScheduler] Starting score scheduler...")

	// Startup: full table update to initialize all scores
	CalculateAndUpdateScores(true)

	// Then tick every 1 minute, only updating posts modified within 24h
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		CalculateAndUpdateScores(false)
	}
}
