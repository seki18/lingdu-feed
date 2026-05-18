package service

import (
	"log"
	"math/rand"

	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/repository"
)

// GetRecommendPosts returns the recommended feed.
// Recommend posts take more than half of total, the rest is 2/3 recent + 1/3 following.
func GetRecommendPosts(requestType string, excludeIDs []int, userID int) ([]model.Posts, error) {
	count := feedCount(requestType)

	// Part A: Recommend > half
	recCount := count/2 + 1
	recommendIDs, _ := repository.GetRecommendPostIDs(recCount, excludeIDs)

	// Part B: Recent 2/3 + Following 1/3 of remaining
	remainder := count - len(recommendIDs)
	recentCount := remainder * 2 / 3
	followingCount := remainder - recentCount

	recentIDs, _ := repository.GetRecentPostIDs(recentCount, excludeIDs, userID)
	followingIDs, _ := repository.GetFollowingPostIDs(followingCount, excludeIDs, userID, true)

	// Build result: recommend first, then randomly interleave recent+following
	result := make([]int, len(recommendIDs))
	copy(result, recommendIDs)

	rest := append(recentIDs, followingIDs...)
	rand.Shuffle(len(rest), func(i, j int) { rest[i], rest[j] = rest[j], rest[i] })
	result = append(result, rest...)

	posts, err := repository.GetPostsByIDs(result)
	if err != nil {
		return nil, err
	}
	log.Printf("[GetRecommendPosts] rec=%d recent=%d following=%d posts=%d",
		len(recommendIDs), len(recentIDs), len(followingIDs), len(posts))
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
	log.Printf("[GetFollowingPosts] ids=%v posts=%d", ids, len(posts))
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
