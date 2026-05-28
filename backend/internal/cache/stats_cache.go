package cache

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/repository"
)

const statsKeyPrefix = "stats:"
const statsTTL = 3600 // 1 hour

// statsKey returns the Redis key for a post's stats HASH.
func statsKey(postID int) string {
	return fmt.Sprintf("%s%d", statsKeyPrefix, postID)
}

// ensureCache loads stats from DB into Redis if not already cached.
func ensureCache(ctx context.Context, postID int) error {
	exists, err := common.Redis.Exists(ctx, statsKey(postID)).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return nil
	}

	s, err := repository.GetStatsByPostID(postID)
	if err != nil {
		return err
	}

	pipe := common.Redis.Pipeline()
	pipe.HSet(ctx, statsKey(postID),
		"like_count", s.LikeCount,
		"comment_count", s.CommentCount,
		"favorite_count", s.FavoriteCount,
		"view_count", s.ViewCount,
		"expose_count", s.ExposeCount,
		"score", fmt.Sprintf("%.6f", s.Score),
	)
	pipe.Expire(ctx, statsKey(postID), statsTTL*time.Second)
	_, err = pipe.Exec(ctx)
	return err
}

// GetStats returns stats for a single post (cache-first, DB fallback).
func GetStats(postID int) (*model.PostStats, error) {
	if common.Redis == nil {
		return repository.GetStatsByPostID(postID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := ensureCache(ctx, postID); err != nil {
		log.Printf("[StatsCache] ensureCache failed for post %d: %v", postID, err)
		return repository.GetStatsByPostID(postID)
	}

	fields, err := common.Redis.HGetAll(ctx, statsKey(postID)).Result()
	if err != nil {
		return repository.GetStatsByPostID(postID)
	}
	if len(fields) == 0 {
		return repository.GetStatsByPostID(postID)
	}

	return parseStats(postID, fields), nil
}

// GetStatsBatch returns stats for multiple posts efficiently.
// Uses pipeline for cache reads; falls back to single DB query for misses.
func GetStatsBatch(postIDs []int) (map[int]*model.PostStats, error) {
	result := make(map[int]*model.PostStats, len(postIDs))
	if len(postIDs) == 0 {
		return result, nil
	}

	if common.Redis == nil {
		stats, err := repository.GetStatsByPostIDs(postIDs)
		if err != nil {
			return nil, err
		}
		for i := range stats {
			result[stats[i].ID] = &stats[i]
		}
		return result, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Pipeline HGETALL for all keys
	pipe := common.Redis.Pipeline()
	cmds := make([]*redis.MapStringStringCmd, len(postIDs))
	for i, id := range postIDs {
		cmds[i] = pipe.HGetAll(ctx, statsKey(id))
	}
	_, _ = pipe.Exec(ctx)

	// Collect misses
	var missedIDs []int
	for i, id := range postIDs {
		fields, err := cmds[i].Result()
		if err != nil || len(fields) == 0 {
			missedIDs = append(missedIDs, id)
		} else {
			result[id] = parseStats(id, fields)
		}
	}

	// Single DB query for all misses
	if len(missedIDs) > 0 {
		dbStats, err := repository.GetStatsByPostIDs(missedIDs)
		if err != nil {
			log.Printf("[StatsCache] DB batch lookup failed: %v", err)
			return result, nil // return what we have
		}

		pipe := common.Redis.Pipeline()
		for i := range dbStats {
			s := &dbStats[i]
			result[s.ID] = s
			pipe.HSet(ctx, statsKey(s.ID),
				"like_count", s.LikeCount,
				"comment_count", s.CommentCount,
				"favorite_count", s.FavoriteCount,
				"view_count", s.ViewCount,
				"expose_count", s.ExposeCount,
				"score", fmt.Sprintf("%.6f", s.Score),
			)
			pipe.Expire(ctx, statsKey(s.ID), statsTTL*time.Second)
		}
		_, _ = pipe.Exec(ctx)
	}

	return result, nil
}

// ── Atomic increment/decrement helpers ──

func incrField(postID int, field string, delta int64) error {
	if common.Redis == nil {
		return nil // degraded mode
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := ensureCache(ctx, postID); err != nil {
		return err
	}

	newVal, err := common.Redis.HIncrBy(ctx, statsKey(postID), field, delta).Result()
	if err != nil {
		return err
	}

	// Floor guard: prevent negative counts
	if newVal < 0 {
		_ = common.Redis.HSet(ctx, statsKey(postID), field, 0).Err()
	}
	return nil
}

// IncrLikeCount increments the like_count for a post in Redis.
func IncrLikeCount(postID int) error { return incrField(postID, "like_count", 1) }

// DecrLikeCount decrements the like_count for a post in Redis (floor 0).
func DecrLikeCount(postID int) error { return incrField(postID, "like_count", -1) }

// IncrCommentCount increments the comment_count for a post in Redis.
func IncrCommentCount(postID int) error { return incrField(postID, "comment_count", 1) }

// DecrCommentCount decrements the comment_count for a post in Redis (floor 0).
func DecrCommentCount(postID int) error { return incrField(postID, "comment_count", -1) }

// IncrFavoriteCount increments the favorite_count for a post in Redis.
func IncrFavoriteCount(postID int) error { return incrField(postID, "favorite_count", 1) }

// DecrFavoriteCount decrements the favorite_count for a post in Redis (floor 0).
func DecrFavoriteCount(postID int) error { return incrField(postID, "favorite_count", -1) }

// IncrExposeCount increments the expose_count for a post in Redis.
func IncrExposeCount(postID int) error { return incrField(postID, "expose_count", 1) }

// IncrViewCount increments the view_count for a post in Redis.
func IncrViewCount(postID int) error { return incrField(postID, "view_count", 1) }

// UpdateScore sets the score for a post in Redis cache (used by scheduler).
func UpdateScore(postID int, score float64) error {
	if common.Redis == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return common.Redis.HSet(ctx, statsKey(postID), "score", fmt.Sprintf("%.6f", score)).Err()
}

// SyncAllToDB scans all stats:* keys and batch-upserts them to the DB.
func SyncAllToDB() {
	if common.Redis == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()
	var stats []model.PostStats
	var cursor uint64
	var keys []string

	for {
		var err error
		keys, cursor, err = common.Redis.Scan(ctx, cursor, statsKeyPrefix+"*", 100).Result()
		if err != nil {
			log.Printf("[StatsCache] Sync scan error: %v", err)
			return
		}

		// Pipeline HGETALL for this batch
		if len(keys) > 0 {
			pipe := common.Redis.Pipeline()
			cmds := make([]*redis.MapStringStringCmd, len(keys))
			for i, key := range keys {
				cmds[i] = pipe.HGetAll(ctx, key)
			}
			_, _ = pipe.Exec(ctx)

			for i, cmd := range cmds {
				fields, err := cmd.Result()
				if err != nil || len(fields) == 0 {
					continue
				}
				id, err := strconv.Atoi(keys[i][len(statsKeyPrefix):])
				if err != nil {
					continue
				}
				stats = append(stats, *parseStats(id, fields))
			}
		}

		if cursor == 0 {
			break
		}
	}

	if len(stats) > 0 {
		// Filter out stats for posts that don't exist (stale/orphan cache entries)
		validIDs, err := repository.GetExistingPostIDsByStats(extractIDs(stats))
		if err != nil {
			log.Printf("[StatsCache] Sync validate error: %v", err)
			return
		}
		validSet := make(map[int]bool, len(validIDs))
		for _, id := range validIDs {
			validSet[id] = true
		}
		filtered := make([]model.PostStats, 0, len(stats))
		for _, s := range stats {
			if validSet[s.ID] {
				filtered = append(filtered, s)
			} else {
				// Clean up orphan Redis key
				common.Redis.Del(ctx, statsKey(s.ID))
			}
		}

		if len(filtered) > 0 {
			if err := repository.BatchUpsertStats(filtered); err != nil {
				log.Printf("[StatsCache] Sync upsert error: %v", err)
				return
			}
			log.Printf("[StatsCache] Synced %d keys to DB in %v", len(filtered), time.Since(start))
		}
	}
}

// parseStats converts Redis HASH fields into a PostStats struct.
func parseStats(id int, fields map[string]string) *model.PostStats {
	s := &model.PostStats{ID: id}
	if v, err := strconv.Atoi(fields["like_count"]); err == nil {
		s.LikeCount = v
	}
	if v, err := strconv.Atoi(fields["comment_count"]); err == nil {
		s.CommentCount = v
	}
	if v, err := strconv.Atoi(fields["favorite_count"]); err == nil {
		s.FavoriteCount = v
	}
	if v, err := strconv.Atoi(fields["view_count"]); err == nil {
		s.ViewCount = v
	}
	if v, err := strconv.Atoi(fields["expose_count"]); err == nil {
		s.ExposeCount = v
	}
	if v, err := strconv.ParseFloat(fields["score"], 64); err == nil {
		s.Score = v
	}
	return s
}

func extractIDs(stats []model.PostStats) []int {
	ids := make([]int, len(stats))
	for i, s := range stats {
		ids[i] = s.ID
	}
	return ids
}
