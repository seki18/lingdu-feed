package service

import (
	"community-backend/internal/model"
	"community-backend/internal/repository"

	"errors"
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

// GetRecentPosts returns all posts ordered by creation time descending.
func GetRecentPosts() ([]model.Posts, error) {
	return repository.GetRecentPosts()
}

// GetPostsByUserID returns all posts authored by the given user.
func GetPostsByUserID(userID int) ([]model.Posts, error) {
	return repository.GetPostsByUserID(userID)
}

// DeletePostByID deletes a single post by its ID.
func DeletePostByID(id int64) error {
	return repository.DeletePostByID(id)
}
