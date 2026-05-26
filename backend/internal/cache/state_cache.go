package cache

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/seki18/lingdu-feed/internal/common"
)

const consumedKeyPrefix = "consumed:"
const consumedTTL = 30 * 60 // 30 minutes, reset on every SADD

func consumedKey(userID int) string {
	return fmt.Sprintf("%s%d", consumedKeyPrefix, userID)
}

// MarkConsumed adds post IDs to the user's consumed SET and resets TTL.
// Used by BatchUpsertState to track which posts the user has seen.
func MarkConsumed(userID int, postIDs []int) error {
	if common.Redis == nil || len(postIDs) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	members := make([]interface{}, len(postIDs))
	for i, id := range postIDs {
		members[i] = strconv.Itoa(id)
	}

	pipe := common.Redis.Pipeline()
	pipe.SAdd(ctx, consumedKey(userID), members...)
	pipe.Expire(ctx, consumedKey(userID), consumedTTL*time.Second)
	_, err := pipe.Exec(ctx)
	return err
}

// GetConsumedIDs returns all consumed post IDs for a user from the cache SET.
// Also refreshes the TTL on hit (sliding TTL: active users keep their cache).
// Returns (ids, cacheHit). On cache miss, caller should fall back to DB.
func GetConsumedIDs(userID int) ([]int, bool, error) {
	if common.Redis == nil {
		return nil, false, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	key := consumedKey(userID)
	members, err := common.Redis.SMembers(ctx, key).Result()
	if err != nil {
		return nil, false, err
	}
	if len(members) == 0 {
		return nil, false, nil // cache miss
	}

	// Refresh TTL (sliding: active users keep cache alive)
	common.Redis.Expire(ctx, key, consumedTTL*time.Second)

	ids := make([]int, 0, len(members))
	for _, m := range members {
		id, err := strconv.Atoi(m)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids, true, nil
}

// SyncConsumedFromDB loads consumed IDs from DB into Redis.
// Called on cache miss to rebuild the SET.
func SyncConsumedFromDB(userID int, postIDs []int, statusThreshold int) error {
	if common.Redis == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Query DB for consumed post IDs among the given candidates
	query, args, err := common.DB.BindNamed(`
		SELECT post_id FROM states
		WHERE user_id = :uid AND post_id = ANY(:pids) AND status > :min
	`, map[string]interface{}{
		"uid":  userID,
		"pids": postIDs,
		"min":  statusThreshold,
	})
	if err != nil {
		return err
	}

	var consumedDB []int
	if err := common.DB.Select(&consumedDB, query, args...); err != nil {
		return err
	}

	if len(consumedDB) == 0 {
		// Mark as seen (empty set) to avoid repeated DB queries
		common.Redis.SAdd(ctx, consumedKey(userID), "_placeholder")
		common.Redis.Expire(ctx, consumedKey(userID), consumedTTL*time.Second)
		common.Redis.SRem(ctx, consumedKey(userID), "_placeholder")
		return nil
	}

	members := make([]interface{}, len(consumedDB))
	for i, id := range consumedDB {
		members[i] = strconv.Itoa(id)
	}

	pipe := common.Redis.Pipeline()
	pipe.SAdd(ctx, consumedKey(userID), members...)
	pipe.Expire(ctx, consumedKey(userID), consumedTTL*time.Second)
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Printf("[StateCache] SyncConsumedFromDB pipeline error: %v", err)
	}
	return err
}
