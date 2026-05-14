package service

import (
	"community-backend/internal/model"
	"community-backend/internal/repository"

	"errors"
)

func GetPostByID(id int) (model.Post, error) {
	return repository.GetPostByID(id)
}

func CreatePost(req model.CreatePostRequest) (model.Post, error) {

	post_create := model.Post{
		UserID: req.UserID,
		Title: req.Title,
		Content: req.Content,
	}

	return repository.CreatePost(post_create)
}

func UpdatePost(req model.UpdatePostRequest) (model.Post, error) {

	post, err := repository.GetPostByID(req.ID)
	if err != nil {
		return model.Post{}, err
	}
	if req.UserID != post.UserID {
		return model.Post{}, errors.New("no power")
	}
	postUpdate := model.Post{
		ID: req.ID,
		UserID: req.UserID,
		Title: req.Title,
		Content: req.Content,
	}

	return repository.UpdatePost(postUpdate)
}

func GetRecentPosts() ([]model.Posts, error) {
	return repository.GetRecentPosts()
}
