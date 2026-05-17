package service

import (
	"community-backend/internal/model"
	"community-backend/internal/repository"
)

// GetHistoryPostsByUserID retrieves History Posts by user ID, with pagination.
func GetHistoryPostsByUserID(id int, page, pageSize int) ([]model.Post, int, error) {
	return repository.GetHistoryPostsByUserID(id, page, pageSize)
}
