package service

import (
	"github.com/seki18/lingdu-feed/internal/cache"
	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/repository"
	post_stats_repository "github.com/seki18/lingdu-feed/internal/repository"

	"errors"
)

// CreatePost inserts a new post into the database and initializes stats.
func CreatePost(req model.CreatePostRequest) (model.Post, error) {
	postCreate := model.Post{
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
	}

	result, err := repository.CreatePost(postCreate)
	if err != nil {
		return result, err
	}
	// Create stats row and warm cache
	_ = post_stats_repository.CreateStats(result.ID)
	_, _ = cache.GetStats(result.ID)             // warm stats cache
	_ = cache.SetContent(result.ID, req.Content) // warm content cache
	// Warm feeditem cache
	_ = cache.SetFeedItem(&model.FeedItem{
		ID:          result.ID,
		UserID:      result.UserID,
		Username:    result.Username,
		Title:       result.Title,
		CreatedTime: result.CreatedTime,
	})
	result.Stats = &model.PostStats{ID: result.ID}
	return result, nil
}

// UpdatePost updates an existing post. Returns an error if the user
// does not own the post. Also updates caches.
func UpdatePost(req model.UpdatePostRequest) (model.Post, error) {
	exists, err := repository.PostExists(req.ID)
	if err != nil || !exists {
		return model.Post{}, errors.New("post not found")
	}

	postUpdate := model.Post{
		ID:      req.ID,
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
	}

	result, err := repository.UpdatePost(postUpdate)
	if err != nil {
		return result, err
	}
	// Update caches
	_ = cache.SetContent(result.ID, result.Content)
	_ = cache.SetFeedItem(&model.FeedItem{
		ID:          result.ID,
		UserID:      result.UserID,
		Username:    result.Username,
		Title:       result.Title,
		CreatedTime: result.CreatedTime,
	})
	return result, nil
}

// GetPostsByUserID returns all posts authored by the given user, with pagination.
func GetPostsByUserID(userID int, page, pageSize int) ([]model.FeedItem, int, error) {
	posts, total, err := repository.GetPostsByUserID(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	// Attach stats from cache
	attachStatsToFeedItems(posts)
	return posts, total, nil
}

// DeletePostByID deletes a single post by its ID. Also cleans up caches.
func DeletePostByID(id int64) error {
	if err := repository.DeletePostByID(id); err != nil {
		return err
	}
	// Clean up caches
	cache.DeleteFeedItem(int(id))
	cache.DeleteContent(int(id))
	return nil
}

// GetPostDetail returns the full detail for a post, including interaction
// status and comments, all fetched concurrently.
func GetPostDetail(id, userID int) (*model.PostDetailResponse, error) {
	var (
		post      model.Post
		liked     bool
		favorited bool
		comments  []model.Comment
	)
	errCh := make(chan error, 5)

	go func() {
		var e error
		post, e = repository.GetPostContentByID(id)
		// Warm content cache in background (best-effort)
		if e == nil && post.Content != "" {
			_ = cache.SetContent(id, post.Content)
		}
		errCh <- e
	}()
	go func() {
		var e error
		liked, e = repository.CheckLiked(userID, id)
		errCh <- e
	}()
	go func() {
		var e error
		favorited, e = repository.CheckFavorited(userID, id)
		errCh <- e
	}()
	go func() {
		var e error
		comments, _, e = repository.GetCommentsByPostID(id, 1, 50) // first page, up to 50
		errCh <- e
	}()
	go func() {
		stats, e := cache.GetStats(id)
		if e == nil && stats != nil {
			post.Stats = stats
		}
		errCh <- e
	}()

	for i := 0; i < 5; i++ {
		if err := <-errCh; err != nil {
			return nil, err
		}
	}

	return &model.PostDetailResponse{
		Post:         post,
		HasLiked:     liked,
		HasFavorited: favorited,
		Comments:     comments,
	}, nil
}
