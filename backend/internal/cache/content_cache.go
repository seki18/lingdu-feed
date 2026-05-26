package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/seki18/lingdu-feed/internal/common"
)

const contentKeyPrefix = "content:"
const contentTTL = 3600 // 1 hour

func contentKey(postID int) string {
	return fmt.Sprintf("%s%d", contentKeyPrefix, postID)
}

// GetContent returns the post content (cache-first, DB fallback via caller).
func GetContent(postID int) (string, error) {
	if common.Redis == nil {
		return "", nil // caller should fall back to DB
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	val, err := common.Redis.Get(ctx, contentKey(postID)).Result()
	if err != nil {
		return "", err // redis.Nil means cache miss
	}
	return val, nil
}

// SetContent writes post content to cache with TTL.
func SetContent(postID int, content string) error {
	if common.Redis == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return common.Redis.SetEx(ctx, contentKey(postID), content, contentTTL*time.Second).Err()
}

// DeleteContent removes post content from cache.
func DeleteContent(postID int) {
	if common.Redis == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	common.Redis.Del(ctx, contentKey(postID))
}
