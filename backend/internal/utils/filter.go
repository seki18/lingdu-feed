package utils

import (
	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/model"
)

// FilterPostIDs removes IDs that the user has already consumed
// (state > StateDelivered). Pass checkStatus=false to skip the state filter.
// Returns the filtered slice, preserving order.
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

	// Step 2: state filter — batch query states for these IDs
	if checkStatus && len(ids) > 0 {
		query, args, err := common.DB.BindNamed(`
			SELECT post_id FROM states
			WHERE user_id = :uid AND post_id = ANY(:pids) AND status > :max
		`, map[string]interface{}{
			"uid":  userID,
			"pids": ids,
			"max":  model.StateDelivered,
		})
		if err != nil {
			return nil, err
		}

		var consumedIDs []int
		if err := common.DB.Select(&consumedIDs, query, args...); err != nil {
			return nil, err
		}

		consumed := make(map[int]bool, len(consumedIDs))
		for _, id := range consumedIDs {
			consumed[id] = true
		}

		filtered := make([]int, 0, len(ids))
		for _, id := range ids {
			if !consumed[id] {
				filtered = append(filtered, id)
			}
		}
		ids = filtered
	}

	return ids, nil
}
