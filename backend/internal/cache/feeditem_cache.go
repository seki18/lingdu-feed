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

const feeditemKeyPrefix = "feeditem:"
const feeditemTTL = 3600 // 1 hour

func feeditemKey(postID int) string {
	return fmt.Sprintf("%s%d", feeditemKeyPrefix, postID)
}

// SetFeedItem writes a single FeedItem into Redis HASH.
func SetFeedItem(item *model.FeedItem) error {
	if common.Redis == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return common.Redis.HSet(ctx, feeditemKey(item.ID),
		"id", item.ID,
		"user_id", item.UserID,
		"username", item.Username,
		"title", item.Title,
		"created_time", item.CreatedTime.Format(time.RFC3339Nano),
	).Err()
}

// GetFeedItemsBatch returns FeedItems for given IDs using pipeline cache reads.
// Misses are fetched from DB in a single query and backfilled to cache.
func GetFeedItemsBatch(ids []int) ([]model.FeedItem, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	// Without Redis, fall back to DB directly
	if common.Redis == nil {
		return repository.GetPostsByIDs(ids)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Pipeline HGETALL for all keys
	pipe := common.Redis.Pipeline()
	cmds := make([]*redis.MapStringStringCmd, len(ids))
	for i, id := range ids {
		cmds[i] = pipe.HGetAll(ctx, feeditemKey(id))
	}
	_, _ = pipe.Exec(ctx)

	// Collect hits and misses
	var missedIDs []int
	hitMap := make(map[int]model.FeedItem, len(ids))
	for i, cmd := range cmds {
		fields, err := cmd.Result()
		if err != nil || len(fields) == 0 {
			missedIDs = append(missedIDs, ids[i])
			continue
		}
		item := parseFeedItem(ids[i], fields)
		if item != nil {
			hitMap[ids[i]] = *item
		} else {
			missedIDs = append(missedIDs, ids[i])
		}
	}

	// Fetch misses from DB and backfill cache
	if len(missedIDs) > 0 {
		dbItems, err := repository.GetPostsByIDs(missedIDs)
		if err != nil {
			return nil, err
		}
		backfillPipe := common.Redis.Pipeline()
		for i := range dbItems {
			item := &dbItems[i]
			hitMap[item.ID] = *item
			backfillPipe.HSet(ctx, feeditemKey(item.ID),
				"id", item.ID,
				"user_id", item.UserID,
				"username", item.Username,
				"title", item.Title,
				"created_time", item.CreatedTime.Format(time.RFC3339Nano),
			)
			backfillPipe.Expire(ctx, feeditemKey(item.ID), feeditemTTL*time.Second)
		}
		if _, err := backfillPipe.Exec(ctx); err != nil {
			log.Printf("[FeedItemCache] backfill pipeline error: %v", err)
		}
	}

	// Reassemble in input order
	result := make([]model.FeedItem, 0, len(ids))
	seen := make(map[int]bool, len(ids))
	for _, id := range ids {
		if seen[id] {
			continue
		}
		seen[id] = true
		if item, ok := hitMap[id]; ok {
			result = append(result, item)
		}
	}
	return result, nil
}

// DeleteFeedItem removes a FeedItem from cache (e.g., on post deletion).
func DeleteFeedItem(postID int) {
	if common.Redis == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	common.Redis.Del(ctx, feeditemKey(postID))
}

func parseFeedItem(id int, fields map[string]string) *model.FeedItem {
	uid, _ := strconv.Atoi(fields["user_id"])
	createdTime, err := time.Parse(time.RFC3339Nano, fields["created_time"])
	if err != nil {
		createdTime = time.Now()
	}
	return &model.FeedItem{
		ID:          id,
		UserID:      uid,
		Username:    fields["username"],
		Title:       fields["title"],
		CreatedTime: createdTime,
	}
}
