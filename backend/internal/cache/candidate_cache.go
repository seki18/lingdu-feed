package cache

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/seki18/lingdu-feed/internal/common"
)

const candidateKey = "candidate"

// GetLatestPostIDs returns up to count newest post IDs from the candidate cache.
func GetLatestPostIDs(count int) ([]int, error) {
	if common.Redis == nil {
		return nil, redis.Nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	members, err := common.Redis.ZRevRange(ctx, candidateKey, 0, int64(count-1)).Result()
	if err != nil {
		return nil, err
	}

	ids := make([]int, 0, len(members))
	for _, m := range members {
		id, err := strconv.Atoi(m)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// AddCandidate adds a newly created post to the candidate cache.
// The score is the Unix timestamp of creation.
// Prunes entries older than 3 days to keep the ZSET from growing unbounded.
func AddCandidate(postID int, createdUnix int64) error {
	if common.Redis == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	member := strconv.Itoa(postID)
	if err := common.Redis.ZAdd(ctx, candidateKey, redis.Z{
		Score:  float64(createdUnix),
		Member: member,
	}).Err(); err != nil {
		log.Printf("[CandidateCache] ZAdd failed: %v", err)
		return err
	}

	// Prune entries older than 3 days
	threeDaysAgo := float64(time.Now().AddDate(0, 0, -3).Unix())
	if err := common.Redis.ZRemRangeByScore(ctx, candidateKey, "-inf", strconv.FormatFloat(threeDaysAgo, 'f', 0, 64)).Err(); err != nil {
		log.Printf("[CandidateCache] Prune failed: %v", err)
	}
	return nil
}
