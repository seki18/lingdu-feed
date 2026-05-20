package service

import (
	"log"

	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/repository"
)

// GetRecommendPosts returns the recommended feed.
// Recommend posts take more than half of total, the rest is 2/3 recent + 1/3 following.
func GetRecommendPosts(requestType string, excludeIDs []int, userID int) ([]model.FeedItem, error) {
	count := feedCount(requestType)

	// Part A: Recommend 1/2 of total, excluding already seen posts. This is the core of the "recommend" feed.
	recentCount := count / 3
	followingCount := count * 1 / 6

	recentIDs, _ := repository.GetRecentPostIDs(recentCount, excludeIDs, userID, false)
	followingIDs, _ := repository.GetFollowingPostIDs(followingCount, excludeIDs, userID, true)

	// Part B: Fill the rest with more recommend posts, excluding all already seen IDs from A. This ensures we always return "count" posts if available.
	recCount := count - len(recentIDs) - len(followingIDs)
	recExcludeIDs := append(excludeIDs, recentIDs...)
	recExcludeIDs = append(recExcludeIDs, followingIDs...)
	recommendIDs, _ := repository.GetRecommendPostIDs(recCount, recExcludeIDs, userID, false)

	// Build result: recommend first, then recent+following (shuffled together)
	seen := make(map[int]bool, count)
	result := make([]int, 0, count)

	for _, id := range recommendIDs {
		if !seen[id] {
			seen[id] = true
			result = append(result, id)
		}
	}
	// Shuffle recent and following together for variety
	shuffled := make([]int, 0, len(recentIDs)+len(followingIDs))
	shuffled = append(shuffled, recentIDs...)
	shuffled = append(shuffled, followingIDs...)
	for _, id := range shuffled {
		if !seen[id] {
			seen[id] = true
			result = append(result, id)
		}
	}
	var degradedIDs []int
	// Degrade: if we didn't get enough posts, fill remaining slots without state filter
	if len(result) < count {
		remaining := count - len(result)
		degradedIDs, _ = repository.GetRecommendPostIDs(remaining, result, userID, true)
		for _, id := range degradedIDs {
			if !seen[id] {
				seen[id] = true
				result = append(result, id)
			}
		}
	}

	posts, err := repository.GetPostsByIDs(result)
	if err != nil {
		return nil, err
	}
	log.Printf("[GetRecommendPosts][recommendIDs=%v recentIDs=%v followingIDs=%v degradedIDs=%v result=%v] posts=%d count=%d",
		recommendIDs, recentIDs, followingIDs, degradedIDs, result, len(posts), count)
	return posts, nil
}

// GetFollowingPosts returns the following feed.
func GetFollowingPosts(requestType string, excludeIDs []int, userID int) ([]model.FeedItem, error) {
	count := feedCount(requestType)
	ids, err := repository.GetFollowingPostIDs(count, excludeIDs, userID, false)
	if err != nil {
		return nil, err
	}
	posts, err := repository.GetPostsByIDs(ids)
	if err != nil {
		return nil, err
	}
	log.Printf("[GetFollowingPosts] ids=%v posts=%d count=%d", ids, len(posts), count)
	return posts, nil
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
