package repository

import (
	"fmt"

	"github.com/lib/pq"
	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/model"
)

// InsertPostImages inserts multiple image records for a post.
func InsertPostImages(postID int, imageURLs []string) ([]model.PostImage, error) {
	if len(imageURLs) == 0 {
		return []model.PostImage{}, nil
	}

	tx, err := common.DB.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var result []model.PostImage
	for i, url := range imageURLs {
		var img model.PostImage
		err := tx.QueryRowx(`
			INSERT INTO post_images (post_id, image_url, sort_order)
			VALUES ($1, $2, $3)
			RETURNING id, post_id, image_url, sort_order
		`, postID, url, i).StructScan(&img)
		if err != nil {
			return nil, fmt.Errorf("failed to insert post image: %w", err)
		}
		result = append(result, img)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, nil
}

// GetImagesByPostID returns all images for a post, ordered by sort_order.
func GetImagesByPostID(postID int) ([]model.PostImage, error) {
	var images []model.PostImage
	err := common.DB.Select(&images, `
		SELECT id, post_id, image_url, sort_order
		FROM post_images
		WHERE post_id = $1
		ORDER BY sort_order ASC
	`, postID)
	return images, err
}

// GetFirstImagesByPostIDs returns the first image (sort_order=0) for each post ID.
// Uses pq.Array for a clean single-parameter query.
func GetFirstImagesByPostIDs(postIDs []int) (map[int]string, error) {
	if len(postIDs) == 0 {
		return map[int]string{}, nil
	}

	rows, err := common.DB.Query(`
		SELECT post_id, image_url
		FROM post_images
		WHERE post_id = ANY($1) AND sort_order = 0
	`, pq.Array(postIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to query first images: %w", err)
	}
	defer rows.Close()

	result := make(map[int]string)
	for rows.Next() {
		var postID int
		var imageURL string
		if err := rows.Scan(&postID, &imageURL); err != nil {
			return nil, fmt.Errorf("failed to scan first image row: %w", err)
		}
		result[postID] = imageURL
	}
	return result, rows.Err()
}
