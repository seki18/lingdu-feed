package service

import (
	"log"
	"time"

	"github.com/seki18/lingdu-feed/internal/cache"
	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/repository"
	"github.com/seki18/lingdu-feed/internal/utils"
)

// FeedCursors holds pagination cursors for each recall source.
type FeedCursors struct {
	RecommendScore float64   `json:"recommend_score"`
	RecommendID    int       `json:"recommend_id"`
	RecentTime     time.Time `json:"recent_time"`
	FollowingTime  time.Time `json:"following_time"`
}

// GetRecommendPosts returns the recommended feed using cursor-based pagination.
func GetRecommendPosts(userID int, cursors FeedCursors) ([]model.FeedItem, FeedCursors, error) {
	count := feedCount("subsequent")

	// Part A: 3 recall sources with cursors — try cache first, fallback to DB
	recentIDs, _ := cache.GetLatestPostIDs(count * 3)
	recentIDs, _ = repository.GetRecentPostIDs(count*3, cursors.RecentTime)
	if len(recentIDs) == 0 {
		recentIDs, _ = repository.GetRecentPostIDs(count*3, cursors.RecentTime)
	}

	followingUserIDs, _ := repository.GetAllFollowingIDs(userID)
	followingPostIDs, _ := repository.GetPostsByFollowingIDs(count*3, followingUserIDs, cursors.FollowingTime)

	recommendIDs, _ := cache.GetTopRankedPostIDs(count * 3)
	if len(recommendIDs) == 0 {
		recommendIDs, _ = repository.GetRecommendPostIDs(count*3, cursors.RecommendScore, cursors.RecommendID)
	}

	// Part B: Go-level filter (excludeIDs + status) — empty excludeIDs for cursor model
	recentIDs, _ = utils.FilterPostIDs(recentIDs, nil, userID, true)
	followingPostIDs, _ = utils.FilterPostIDs(followingPostIDs, nil, userID, true)
	recommendIDs, _ = utils.FilterPostIDs(recommendIDs, nil, userID, true)

	// Part C: assemble result: recommend first, then recent+following shuffled
	result := assembleFeedIDs(recommendIDs, recentIDs, followingPostIDs, count)

	// Degrade: fill remaining with recommend (no status filter)
	if len(result) < count {
		remaining := count - len(result)
		degradeCursorScore, degradeCursorID := lastCursor(result, "recommend")
		degradedIDs, _ := repository.GetRecommendPostIDs(remaining*3, degradeCursorScore, degradeCursorID)
		degradedIDs, _ = utils.FilterPostIDs(degradedIDs, nil, userID, false)
		tail := take(degradedIDs, remaining)
		result = append(result, tail...)
	}

	posts, err := repository.GetPostsByIDs(result)
	if err != nil {
		return nil, FeedCursors{}, err
	}

	// Build new cursors from the last item of each source
	newCursors := FeedCursors{
		RecommendScore: cursorScore(result),
		RecommendID:    cursorID(result),
		RecentTime:     cursorTime(result),
		FollowingTime:  cursorTime(result),
	}

	log.Printf("[GetRecommendPosts] result=%v posts=%d count=%d", result, len(posts), count)
	return posts, newCursors, nil
}

// GetFollowingPosts returns the following-only feed.
func GetFollowingPosts(userID int, cursor time.Time) ([]model.FeedItem, time.Time, error) {
	count := feedCount("subsequent")
	followingUserIDs, _ := repository.GetAllFollowingIDs(userID)
	ids, err := repository.GetPostsByFollowingIDs(count*3, followingUserIDs, cursor)
	if err != nil {
		return nil, time.Time{}, err
	}
	ids, _ = utils.FilterPostIDs(ids, nil, userID, true)
	ids = take(ids, count)
	posts, err := repository.GetPostsByIDs(ids)
	if err != nil {
		return nil, time.Time{}, err
	}
	var newCursor time.Time
	if len(ids) > 0 {
		// cursor is the created_time of the last post in this batch
		row := posts[len(posts)-1]
		newCursor = row.CreatedTime
	}
	log.Printf("[GetFollowingPosts] ids=%v posts=%d count=%d", ids, len(posts), count)
	return posts, newCursor, nil
}

// ── helpers ──

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

func take(ids []int, n int) []int {
	if len(ids) <= n {
		return ids
	}
	return ids[:n]
}

func lastCursor(ids []int, source string) (float64, int) {
	// For degradation: use the last ID as cursor for recommend source
	if len(ids) == 0 {
		return 0, 0
	}
	return 0, ids[len(ids)-1]
}

func cursorScore(ids []int) float64 { return 0 } // post.CreatedTime for recent/following cursors
func cursorID(ids []int) int        { return 0 }
func cursorTime(ids []int) time.Time {
	return time.Time{} // stub — real implementation needs to look up post created_time
}

// GetHistoryPosts returns the user's browsing history as feed posts.
func GetHistoryPosts(userID, page, pageSize int) ([]model.FeedItem, int, error) {
	ids, total, err := repository.GetHistoryPostIDs(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	posts, err := repository.GetPostsByIDs(ids)
	if err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

// GetFavoriteFeed returns the user's favorited posts as feed posts.
func GetFavoriteFeed(userID, page, pageSize int) ([]model.FeedItem, int, error) {
	ids, total, err := repository.GetFavoritePostIDs(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	posts, err := repository.GetPostsByIDs(ids)
	if err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

// GetAuthorPosts returns posts authored by the given user with pagination.
func GetAuthorPosts(userID, page, pageSize int) ([]model.FeedItem, int, error) {
	return repository.GetPostsByUserID(userID, page, pageSize)
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

	return user, posts, total, nil
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
