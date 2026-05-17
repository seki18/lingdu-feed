package service

import (
	"community-backend/internal/model"
	"community-backend/internal/repository"

	"errors"
	"log"
)

// GetPostByID retrieves a single post by its ID.
func GetPostByID(id int) (model.Post, error) {
	return repository.GetPostByID(id)
}

// CreatePost inserts a new post into the database.
func CreatePost(req model.CreatePostRequest) (model.Post, error) {
	postCreate := model.Post{
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
	}

	return repository.CreatePost(postCreate)
}

// UpdatePost updates an existing post. Returns an error if the user
// does not own the post.
func UpdatePost(req model.UpdatePostRequest) (model.Post, error) {
	post, err := repository.GetPostByID(req.ID)
	if err != nil {
		return model.Post{}, err
	}
	if req.UserID != post.UserID {
		return model.Post{}, errors.New("no power")
	}
	postUpdate := model.Post{
		ID:      req.ID,
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
	}

	return repository.UpdatePost(postUpdate)
}

// GetRecentPosts returns posts ordered by creation time descending with request type controlling count.
func GetRecentPosts(requestType string, excludeIDs []int) ([]model.Posts, error) {
	count := 3
	switch requestType {
	case "initial", "refresh":
		count = 5
	case "subsequent", "next", "more":
		count = 2
	default:
		count = 2
	}

	posts, err := repository.GetRecentPosts(count, excludeIDs)
	if err != nil {
		return nil, err
	}
	log.Printf("[GetRecentPosts] Request type=%s count=%d fetched=%d excludeIDs=%v posts: %+v", requestType, count, len(posts), excludeIDs, posts)
	if len(posts) > count {
		for i := len(posts) - 1; i >= 0; i-- {
			if len(posts) <= count {
				break
			}
			log.Printf("[GetRecentPosts] Checking post ID %d with status %d", posts[i].ID, posts[i].Status)
			if posts[i].Status >= model.FeedDisplay {
				posts = append(posts[:i], posts[i+1:]...)
			}
		}
	}
	log.Printf("[GetRecentPosts] Returning %d posts after filtering, count: %d posts: %+v", len(posts), count, posts)
	return posts, nil
}

// GetPostsByUserID returns all posts authored by the given user, with pagination.
func GetPostsByUserID(userID int, page, pageSize int) ([]model.Posts, int, error) {
	return repository.GetPostsByUserID(userID, page, pageSize)
}

// DeletePostByID deletes a single post by its ID.
func DeletePostByID(id int64) error {
	return repository.DeletePostByID(id)
}
