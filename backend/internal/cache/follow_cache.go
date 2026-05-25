package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/seki18/lingdu-feed/internal/common"
)

const (
	followKeyPrefix = "follow:"
	followTTL       = 24 * time.Hour
)

func followKey(userID int) string {
	return fmt.Sprintf("%s%d", followKeyPrefix, userID)
}

// GetFollowingIDs returns the cached following user IDs for a follower.
func GetFollowingIDs(followerID int) ([]int, error) {
	if common.Redis == nil {
		return nil, redis.Nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	raw, err := common.Redis.Get(ctx, followKey(followerID)).Result()
	if err != nil {
		return nil, err
	}

	var ids []int
	if err := json.Unmarshal([]byte(raw), &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// SetFollowingIDs writes the following user IDs for a follower with 24h TTL.
func SetFollowingIDs(followerID int, ids []int) error {
	if common.Redis == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	data, err := json.Marshal(ids)
	if err != nil {
		return err
	}

	if err := common.Redis.Set(ctx, followKey(followerID), data, followTTL).Err(); err != nil {
		log.Printf("[FollowCache] Set failed: %v", err)
		return err
	}
	return nil
}

// InvalidateFollow removes the cached following list for a follower.
func InvalidateFollow(followerID int) {
	if common.Redis == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := common.Redis.Del(ctx, followKey(followerID)).Err(); err != nil {
		log.Printf("[FollowCache] Delete failed: %v", err)
	}
}
