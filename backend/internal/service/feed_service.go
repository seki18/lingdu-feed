package service

import (
	"log"

	"github.com/seki18/lingdu-feed/internal/cache"
	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/repository"
	post_stats_repository "github.com/seki18/lingdu-feed/internal/repository"
	"github.com/seki18/lingdu-feed/internal/utils"
)

// GetRecommendPosts returns the recommended feed using cursor-based pagination.
// Retrieves post IDs from three recall sources (recommend, recent, following),
// applies state-based deduplication, and assembles a hybrid feed.
//
// cache-first: each source queries Redis first; on miss or insufficient results,
// falls back to PostgreSQL. cursorID=0 means first page; >0 means subsequent page.
func GetRecommendPosts(userID int, requestType string, cursorID int) ([]model.FeedItem, int, error) {
	count := feedCount(requestType)

	// ── 1. Recent posts (cache → cursor filter → DB fallback) ──
	recentIDs, _ := cache.GetLatestPostIDs(count * 3)
	if len(recentIDs) > 0 {
		log.Printf("[GetRecommendPosts] recent=HIT(cache) count=%d", len(recentIDs))
	} else {
		log.Printf("[GetRecommendPosts] recent=MISS(cache)")
	}
	if cursorID > 0 && len(recentIDs) > 0 {
		before := len(recentIDs)
		recentIDs = filterByCursor(recentIDs, cursorID)
		log.Printf("[GetRecommendPosts] recent filterByCursor cursor=%d before=%d after=%d", cursorID, before, len(recentIDs))
	}
	if len(recentIDs) == 0 {
		recentIDs, _ = repository.GetRecentPostIDsCursor(count*3, cursorID)
		log.Printf("[GetRecommendPosts] recent=HIT(db) count=%d cursor=%d", len(recentIDs), cursorID)
	}

	// ── 2. Following posts (cache user list → DB fallback) ──
	followingUserIDs, _ := cache.GetFollowingIDs(userID)
	if len(followingUserIDs) > 0 {
		log.Printf("[GetRecommendPosts] followingUserIds=HIT(cache) count=%d", len(followingUserIDs))
	} else {
		log.Printf("[GetRecommendPosts] followingUserIds=MISS(cache)")
	}
	if len(followingUserIDs) == 0 {
		followingUserIDs, _ = repository.GetAllFollowingIDs(userID)
		if len(followingUserIDs) > 0 {
			_ = cache.SetFollowingIDs(userID, followingUserIDs)
		}
		log.Printf("[GetRecommendPosts] followingUserIds=HIT(db) count=%d", len(followingUserIDs))
	}
	followingPostIDs, _ := repository.GetPostsByFollowingIDsCursor(count*3, followingUserIDs, cursorID)

	// ── 3. Recommend posts (cache → cursor filter → DB fallback) ──
	recommendIDs, _ := cache.GetTopRankedPostIDs(count * 3)
	if len(recommendIDs) > 0 {
		log.Printf("[GetRecommendPosts] recommend=HIT(cache) count=%d", len(recommendIDs))
	} else {
		log.Printf("[GetRecommendPosts] recommend=MISS(cache)")
	}
	if cursorID > 0 && len(recommendIDs) > 0 {
		before := len(recommendIDs)
		recommendIDs = filterByCursor(recommendIDs, cursorID)
		log.Printf("[GetRecommendPosts] recommend filterByCursor cursor=%d before=%d after=%d", cursorID, before, len(recommendIDs))
	}
	if len(recommendIDs) == 0 {
		recommendIDs, _ = post_stats_repository.GetRecommendPostIDs(count*3, cursorID)
		log.Printf("[GetRecommendPosts] recommend=HIT(db) count=%d cursor=%d", len(recommendIDs), cursorID)
	}

	// ── State filter ──
	if userID != -1 {
		recentIDs, _ = utils.FilterPostIDs(recentIDs, userID, true)
		followingPostIDs, _ = utils.FilterPostIDs(followingPostIDs, userID, true)
		recommendIDs, _ = utils.FilterPostIDs(recommendIDs, userID, true)
	}

	// ── Assemble: recommend first, then recent+following shuffled ──
	result := assembleFeedIDs(recommendIDs, recentIDs, followingPostIDs, count)

	// ── Degrade: fill remaining with recommend (no status filter) ──
	if len(result) < count {
		remaining := count - len(result)
		degradeCursorID := 0
		if len(result) > 0 {
			degradeCursorID = result[len(result)-1]
		}
		degradedIDs, _ := post_stats_repository.GetRecommendPostIDs(remaining*3, degradeCursorID)
		if userID != -1 {
			degradedIDs, _ = utils.FilterPostIDs(degradedIDs, userID, false)
		}
		result = append(result, take(degradedIDs, remaining)...)
		log.Printf("[GetRecommendPosts] degrade=%d remaining=%d", len(degradedIDs), remaining)
	}

	posts, err := cache.GetFeedItemsBatch(result)
	if err != nil {
		return nil, 0, err
	}
	attachStatsToFeedItems(posts)

	newCursorID := 0
	if len(posts) > 0 {
		newCursorID = posts[len(posts)-1].ID
	}

	log.Printf("[GetRecommendPosts] requestType=%s cursor=%d result=%v posts=%d count=%d newCursor=%d",
		requestType, cursorID, result, len(posts), count, newCursorID)
	return posts, newCursorID, nil
}

// GetFollowingPosts returns the following-only feed.
func GetFollowingPosts(userID int, cursorID int) ([]model.FeedItem, int, error) {
	count := feedCount("subsequent")
	followingUserIDs, _ := repository.GetAllFollowingIDs(userID)
	ids, err := repository.GetPostsByFollowingIDsCursor(count*3, followingUserIDs, cursorID)
	if err != nil {
		return nil, 0, err
	}
	ids, _ = utils.FilterPostIDs(ids, userID, true)
	ids = take(ids, count)
	posts, err := cache.GetFeedItemsBatch(ids)
	if err != nil {
		return nil, 0, err
	}
	attachStatsToFeedItems(posts)
	newCursorID := 0
	if len(posts) > 0 {
		newCursorID = posts[len(posts)-1].ID
	}
	log.Printf("[GetFollowingPosts] cursor=%d ids=%v posts=%d count=%d", cursorID, ids, len(posts), count)
	return posts, newCursorID, nil
}

// assembleFeedIDs merges three ID slices into a single deduplicated result.
// recommend IDs are placed first, followed by interleaved recent + following IDs.
func assembleFeedIDs(recommend, recent, following []int, count int) []int {
	seen := make(map[int]bool, count)
	result := make([]int, 0, count)
	for _, id := range recommend {
		if !seen[id] && len(result) < count {
			seen[id] = true
			result = append(result, id)
		}
	}
	shuffled := make([]int, 0, len(recent)+len(following))
	shuffled = append(shuffled, recent...)
	shuffled = append(shuffled, following...)
	for _, id := range shuffled {
		if !seen[id] && len(result) < count {
			seen[id] = true
			result = append(result, id)
		}
	}
	return result
}

// take returns the first n elements of ids (or all if n ≥ len(ids)).
func take(ids []int, n int) []int {
	if len(ids) <= n {
		return ids
	}
	return ids[:n]
}

// filterByCursor keeps only IDs strictly less than cursorID.
// Used to filter cached IDs when paginating beyond the first page.
func filterByCursor(ids []int, cursorID int) []int {
	result := make([]int, 0, len(ids))
	for _, id := range ids {
		if id < cursorID {
			result = append(result, id)
		}
	}
	return result
}

// GetHistoryPosts returns the user's browsing history as feed posts.
func GetHistoryPosts(userID, page, pageSize int) ([]model.FeedItem, int, error) {
	ids, total, err := repository.GetHistoryPostIDs(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	posts, err := cache.GetFeedItemsBatch(ids)
	if err != nil {
		return nil, 0, err
	}
	attachStatsToFeedItems(posts)
	return posts, total, nil
}

// GetFavoriteFeed returns the user's favorited posts as feed posts.
func GetFavoriteFeed(userID, page, pageSize int) ([]model.FeedItem, int, error) {
	ids, total, err := repository.GetFavoritePostIDs(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	posts, err := cache.GetFeedItemsBatch(ids)
	if err != nil {
		return nil, 0, err
	}
	attachStatsToFeedItems(posts)
	return posts, total, nil
}

// GetAuthorPosts returns posts authored by the given user with pagination.
func GetAuthorPosts(userID, page, pageSize int) ([]model.FeedItem, int, error) {
	posts, total, err := repository.GetPostsByUserID(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	attachStatsToFeedItems(posts)
	return posts, total, nil
}

// GetAuthorPage returns author profile and paginated authored posts in one response.
func GetAuthorPage(userID, currentUserID, page, pageSize int) (model.User, []model.FeedItem, int, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return model.User{}, nil, 0, err
	}
	if currentUserID > 0 && currentUserID != userID {
		following, _ := repository.IsFollowExist(model.Follow{
			FollowerID:  currentUserID,
			FollowingID: userID,
		})
		user.IsFollowing = following
	}

	posts, total, err := repository.GetPostsByUserID(userID, page, pageSize)
	if err != nil {
		return model.User{}, nil, 0, err
	}
	attachStatsToFeedItems(posts)
	return user, posts, total, nil
}

func attachStatsToFeedItems(items []model.FeedItem) {
	for i := range items {
		stats, err := cache.GetStats(items[i].ID)
		if err == nil && stats != nil {
			items[i].Stats = stats
		}
	}
}

func feedCount(requestType string) int {
	switch requestType {
	case "initial", "refresh":
		return 6
	case "subsequent", "next", "more":
		return 10
	default:
		return 10
	}
}
