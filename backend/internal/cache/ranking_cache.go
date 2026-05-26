package cache

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/seki18/lingdu-feed/internal/common"
)

const rankingKey = "ranking"

// GetTopRankedPostIDs returns up to count post IDs from the ranking cache, newest first.
func GetTopRankedPostIDs(count int) ([]int, error) {
	if common.Redis == nil {
		return nil, redis.Nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	members, err := common.Redis.ZRevRangeByScore(ctx, rankingKey, &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Result()
	if err != nil {
		return nil, err
	}

	ids := make([]int, 0, count)
	for _, m := range members {
		id, err := strconv.Atoi(m)
		if err != nil {
			continue
		}
		ids = append(ids, id)
		if len(ids) >= count {
			break
		}
	}
	return ids, nil
}

// RefreshRanking rebuilds the ranking cache with the top 1000 posts by score.
func RefreshRanking() error {
	if common.Redis == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Fetch top 1000 post IDs with scores from DB
	var rows []struct {
		ID    int     `db:"id"`
		Score float64 `db:"score"`
	}
	if err := common.DB.Select(&rows, `SELECT id, score FROM post_stats ORDER BY score DESC LIMIT 1000`); err != nil {
		log.Printf("[RankingCache] DB query failed: %v", err)
		return err
	}

	// Clear and rebuild the ZSET
	if err := common.Redis.Del(ctx, rankingKey).Err(); err != nil {
		log.Printf("[RankingCache] Delete failed: %v", err)
		return err
	}

	members := make([]redis.Z, len(rows))
	for i, r := range rows {
		members[i] = redis.Z{Score: r.Score, Member: strconv.Itoa(r.ID)}
	}
	if err := common.Redis.ZAdd(ctx, rankingKey, members...).Err(); err != nil {
		log.Printf("[RankingCache] ZAdd failed: %v", err)
		return err
	}
	log.Printf("[RankingCache] Refreshed with %d posts", len(rows))
	return nil
}
