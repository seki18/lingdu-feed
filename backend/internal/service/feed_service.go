package service

import (
	"log"

	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/repository"
)

// GetRecommendPosts returns the recommended feed.
// Recommend posts take more than half of total, the rest is 2/3 recent + 1/3 following.
func GetRecommendPosts(requestType string, excludeIDs []int, userID int) ([]model.Posts, error) {
	count := feedCount(requestType)

	// Part A: Recommend 1/2 of total, excluding already seen posts. This is the core of the "recommend" feed.
	recentCount := count / 3
	followingCount := count * 1 / 6

	recentIDs, _ := repository.GetRecentPostIDs(recentCount, excludeIDs, userID)
	followingIDs, _ := repository.GetFollowingPostIDs(followingCount, excludeIDs, userID, true)

	// Part B: Fill the rest with more recommend posts, excluding all already seen IDs from A. This ensures we always return "count" posts if available.
	recCount := count - len(recentIDs) - len(followingIDs)
	recExcludeIDs := append(excludeIDs, recentIDs...)
	recExcludeIDs = append(recExcludeIDs, followingIDs...)
	recommendIDs, _ := repository.GetRecommendPostIDs(recCount, recExcludeIDs, userID)

	// Build result: recommend first, then randomly interleave recent+following
	// Use a seen set to prevent duplicates across the three sources
	seen := make(map[int]bool, count)
	result := make([]int, 0, count)

	for _, id := range recommendIDs {
		if !seen[id] {
			seen[id] = true
			result = append(result, id)
		}
	}
	for _, id := range recentIDs {
		if !seen[id] {
			seen[id] = true
			result = append(result, id)
		}
	}
	for _, id := range followingIDs {
		if !seen[id] {
			seen[id] = true
			result = append(result, id)
		}
	}

	posts, err := repository.GetPostsByIDs(result)
	if err != nil {
		return nil, err
	}
	log.Printf("[GetRecommendPosts][recommendIDs=%v recentIDs=%v followingIDs=%v] posts=%d count=%d",
		recommendIDs, recentIDs, followingIDs, len(posts), count)
	return posts, nil
}

// GetFollowingPosts returns the following feed.
func GetFollowingPosts(requestType string, excludeIDs []int, userID int) ([]model.Posts, error) {
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
func GetHistoryPosts(userID, page, pageSize int) ([]model.Posts, int, error) {
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

// GetCollectionPosts returns the user's collected posts as feed posts.
func GetCollectionPosts(userID, page, pageSize int) ([]model.Posts, int, error) {
	ids, total, err := repository.GetCollectionPostIDs(userID, page, pageSize)
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
func GetAuthorPosts(userID, page, pageSize int) ([]model.Posts, int, error) {
	return repository.GetPostsByUserID(userID, page, pageSize)
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
