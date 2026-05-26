package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/model"

	"github.com/jmoiron/sqlx"
)

// GetPostContentByID retrieves post content fields by primary key (no stats).
func GetPostContentByID(id int) (model.Post, error) {
	var post model.Post

	err := common.DB.Get(&post, `
		SELECT p.id, p.user_id, u.username, p.title, p.content,
			p.created_time, p.updated_time
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.id = $1
	`, id)

	return post, err
}

// PostExists checks whether a post with the given ID exists.
func PostExists(id int) (bool, error) {
	var exists bool
	err := common.DB.Get(&exists, `SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)`, id)
	return exists, err
}

// CreatePost inserts a new post and returns the created record.
func CreatePost(post model.Post) (model.Post, error) {
	err := common.DB.QueryRowx(`
		INSERT INTO posts (user_id, title, content, created_time, updated_time)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, user_id, title, content, created_time, updated_time
	`, post.UserID, post.Title, post.Content).
		StructScan(&post)

	return post, err
}

// UpdatePost updates an existing post and returns the updated record.
func UpdatePost(post model.Post) (model.Post, error) {
	err := common.DB.QueryRowx(`
		UPDATE posts
		SET title = $3, content = $4, updated_time = now()
		WHERE id = $1 and user_id = $2
		RETURNING id, user_id, title, content, created_time, updated_time
	`, post.ID, post.UserID, post.Title, post.Content).
		StructScan(&post)

	return post, err
}

// DeletePostByID delete a single post by primary key.
func DeletePostByID(id int64) error {
	result, err := common.DB.Exec(`
		DELETE
		FROM posts
		WHERE id = $1
	`, id)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("post not found")
	}

	return nil
}

// GetPostsByIDs returns post summary rows for the given IDs (no stats),
// preserving the order of the supplied slice by re-sorting in Go.
func GetPostsByIDs(ids []int) ([]model.FeedItem, error) {
	if len(ids) == 0 {
		return []model.FeedItem{}, nil
	}

	query := `
		SELECT p.id, p.user_id, u.username, p.title,
			p.created_time
		FROM posts as p
		LEFT JOIN users as u ON p.user_id = u.id
		WHERE p.id IN (?)
	`
	query, args, err := sqlx.In(query, ids)
	if err != nil {
		return nil, err
	}
	query = common.DB.Rebind(query)

	var rows []model.FeedItem
	if err := common.DB.Select(&rows, query, args...); err != nil {
		return nil, err
	}

	// Re-sort rows to match input ID order
	rowMap := make(map[int]model.FeedItem, len(rows))
	for _, row := range rows {
		rowMap[row.ID] = row
	}
	seen := make(map[int]bool, len(ids))
	result := make([]model.FeedItem, 0, len(ids))
	for _, id := range ids {
		if seen[id] {
			continue
		}
		seen[id] = true
		if row, ok := rowMap[id]; ok {
			result = append(result, row)
		}
	}
	return result, nil
}

// GetPostsByUserID returns posts authored by the given user (no stats), newest first, with pagination.
func GetPostsByUserID(userID int, page, pageSize int) ([]model.FeedItem, int, error) {
	var posts []model.FeedItem
	var total int

	if err := common.DB.Get(&total, `SELECT COUNT(1) FROM posts WHERE user_id = $1`, userID); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := common.DB.Select(&posts, `
		SELECT p.id, p.user_id, u.username, p.title,
			p.created_time
		FROM posts as p
		LEFT JOIN users as u ON p.user_id = u.id
		WHERE p.user_id = $1
		ORDER BY created_time DESC
		LIMIT $2 OFFSET $3
	`, userID, pageSize, offset)

	return posts, total, err
}

// GetRecentPostIDs returns the most recent post IDs, newest first.
func GetRecentPostIDs(count int, cursor time.Time) ([]int, error) {
	query := `SELECT p.id FROM posts as p`
	var args []interface{}

	if !cursor.IsZero() {
		query += ` WHERE p.created_time < $1`
		args = append(args, cursor)
	}
	query += fmt.Sprintf(` ORDER BY p.created_time DESC LIMIT $%d`, len(args)+1)
	args = append(args, count)

	var ids []int
	if err := common.DB.Select(&ids, common.DB.Rebind(query), args...); err != nil {
		return nil, err
	}
	return ids, nil
}

// GetRecentPostIDsCursor returns recent post IDs using id cursor.
func GetRecentPostIDsCursor(count int, cursorID int) ([]int, error) {
	query := `
		SELECT p.id FROM posts p
		WHERE ($1 = 0 OR p.id < $1)
		ORDER BY p.created_time DESC
		LIMIT $2
	`
	var ids []int
	if err := common.DB.Select(&ids, query, cursorID, count); err != nil {
		return nil, err
	}
	return ids, nil
}

// GetPostsByFollowingIDs returns post IDs from the given following user ID list.
func GetPostsByFollowingIDs(count int, followingIDs []int, cursor time.Time) ([]int, error) {
	if len(followingIDs) == 0 {
		return []int{}, nil
	}

	query := `SELECT p.id FROM posts as p`
	var args []interface{}

	q, idsArgs, err := sqlx.In(` WHERE p.user_id IN (?)`, followingIDs)
	if err != nil {
		return nil, err
	}
	query += q
	args = append(args, idsArgs...)

	if !cursor.IsZero() {
		query += fmt.Sprintf(` AND p.created_time < $%d`, len(args)+1)
		args = append(args, cursor)
	}
	query += fmt.Sprintf(` ORDER BY p.created_time DESC LIMIT $%d`, len(args)+1)
	args = append(args, count)

	var ids []int
	if err := common.DB.Select(&ids, common.DB.Rebind(query), args...); err != nil {
		return nil, err
	}
	return ids, nil
}

// GetPostsByFollowingIDsCursor returns post IDs from following users using id cursor.
func GetPostsByFollowingIDsCursor(count int, followingIDs []int, cursorID int) ([]int, error) {
	if len(followingIDs) == 0 {
		return []int{}, nil
	}

	query := `SELECT p.id FROM posts as p`
	var args []interface{}

	q, idsArgs, err := sqlx.In(` WHERE p.user_id IN (?)`, followingIDs)
	if err != nil {
		return nil, err
	}
	query += q
	args = append(args, idsArgs...)

	if cursorID != 0 {
		query += fmt.Sprintf(` AND p.id < $%d`, len(args)+1)
		args = append(args, cursorID)
	}
	query += fmt.Sprintf(` ORDER BY p.created_time DESC LIMIT $%d`, len(args)+1)
	args = append(args, count)

	var ids []int
	if err := common.DB.Select(&ids, common.DB.Rebind(query), args...); err != nil {
		return nil, err
	}
	return ids, nil
}
