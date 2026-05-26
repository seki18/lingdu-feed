package utils

import (
	"log"

	"github.com/seki18/lingdu-feed/internal/cache"
	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/model"
)

// FilterPostIDs removes IDs that the user has already consumed
// (state > StateDelivered). Pass checkStatus=false to skip the state filter.
// Returns the filtered slice, preserving order.
// Uses Redis SET cache with sliding TTL when available, falls back to DB.
func FilterPostIDs(ids []int, userID int, checkStatus bool) ([]int, error) {
	if len(ids) == 0 {
		return ids, nil
	}

	// Step 1: deduplicate input (preserve first occurrence)
	seen := make(map[int]bool, len(ids))
	deduped := make([]int, 0, len(ids))
	for _, id := range ids {
		if !seen[id] {
			seen[id] = true
			deduped = append(deduped, id)
		}
	}
	ids = deduped

	// Step 2: state filter — try Redis SET cache first, then DB
	if checkStatus && len(ids) > 0 {
		consumed := make(map[int]bool, len(ids))

		// Try Redis SET cache
		cachedIDs, cacheHit, _ := cache.GetConsumedIDs(userID)
		if cacheHit {
			for _, cid := range cachedIDs {
				consumed[cid] = true
			}
			log.Printf("[FilterPostIDs] user=%d cache=HIT cached=%d input=%d",
				userID, len(cachedIDs), len(ids))
		} else {
			// Cache miss: query DB once and rebuild cache
			log.Printf("[FilterPostIDs] user=%d cache=MISS querying DB...", userID)
			_ = cache.SyncConsumedFromDB(userID, ids, int(model.StateDelivered))

			// Read back from the just-rebuilt cache (avoids a second DB query)
			cachedIDs, cacheHit, _ = cache.GetConsumedIDs(userID)
			if cacheHit {
				for _, cid := range cachedIDs {
					consumed[cid] = true
				}
				log.Printf("[FilterPostIDs] user=%d rebuilt cache, cached=%d", userID, len(cachedIDs))
			} else {
				// Fallback: direct DB query if cache rebuild failed
				query, args, _ := common.DB.BindNamed(`
					SELECT post_id FROM states
					WHERE user_id = :uid AND post_id = ANY(:pids) AND status > :max
				`, map[string]interface{}{
					"uid":  userID,
					"pids": ids,
					"max":  model.StateDelivered,
				})
				var consumedIDs []int
				common.DB.Select(&consumedIDs, query, args...)
				for _, id := range consumedIDs {
					consumed[id] = true
				}
			}
		}

		filtered := make([]int, 0, len(ids))
		for _, id := range ids {
			if !consumed[id] {
				filtered = append(filtered, id)
			}
		}
		log.Printf("[FilterPostIDs] user=%d before=%d after=%d removed=%d",
			userID, len(ids), len(filtered), len(ids)-len(filtered))
		ids = filtered
	}

	return ids, nil
}
