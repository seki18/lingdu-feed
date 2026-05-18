package service

import (
	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/repository"

	"errors"
)

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

	return repository.UpdatePost(postUpdate)
}

// GetPostsByUserID returns all posts authored by the given user, with pagination.
func GetPostsByUserID(userID int, page, pageSize int) ([]model.Posts, int, error) {
	return repository.GetPostsByUserID(userID, page, pageSize)
}

// DeletePostByID deletes a single post by its ID.
func DeletePostByID(id int64) error {
	return repository.DeletePostByID(id)
}

// GetPostDetail returns the full detail for a post, including interaction
// status and comments, all fetched concurrently.
func GetPostDetail(id, userID int) (*model.PostDetail, error) {
	var (
		post      model.Post
		praised   bool
		collected bool
		comments  []model.Comment
	)
	errCh := make(chan error, 4)

	go func() {
		var e error
		post, e = repository.GetPostContentByID(id)
		errCh <- e
	}()
	go func() {
		var e error
		praised, e = repository.CheckPraised(userID, id)
		errCh <- e
	}()
	go func() {
		var e error
		collected, e = repository.CheckCollected(userID, id)
		errCh <- e
	}()
	go func() {
		var e error
		comments, e = repository.GetCommentsByPostID(id)
		errCh <- e
	}()

	for i := 0; i < 4; i++ {
		if err := <-errCh; err != nil {
			return nil, err
		}
	}

	return &model.PostDetail{
		Post:         post,
		HasPraised:   praised,
		HasCollected: collected,
		Comments:     comments,
	}, nil
}
