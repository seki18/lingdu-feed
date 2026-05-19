package repository

import (
	"errors"

	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/model"

	"github.com/jmoiron/sqlx"
)

// GetPostContentByID retrieves post content fields by primary key.
func GetPostContentByID(id int) (model.Post, error) {
	var post model.Post

	err := common.DB.Get(&post, `
		SELECT p.id, p.user_id, u.username, p.title, p.content,
			p.praise_count, p.comment_count, p.collection_count, p.view_count,
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

// GetPostsByIDs returns post summary rows for the given IDs, preserving the
// order of the supplied slice by re-sorting in Go after the SQL query.
func GetPostsByIDs(ids []int) ([]model.Posts, error) {
	if len(ids) == 0 {
		return []model.Posts{}, nil
	}

	query := `
		SELECT p.id, p.user_id, u.username, p.title,
			p.praise_count, p.comment_count, p.collection_count, p.view_count,
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

	var rows []model.Posts
	if err := common.DB.Select(&rows, query, args...); err != nil {
		return nil, err
	}

	// Re-sort rows to match input ID order (deduplicated: only first occurrence gets the row)
	idIndex := make(map[int]int, len(ids))
	seen := make(map[int]bool, len(ids))
	for i, id := range ids {
		if !seen[id] {
			idIndex[id] = i
			seen[id] = true
		}
	}
	posts := make([]model.Posts, len(ids))
	written := 0
	for _, row := range rows {
		if idx, ok := idIndex[row.ID]; ok {
			posts[idx] = row
			written++
		}
	}
	return posts[:written], nil
}

// GetPostStatsByIDs returns lightweight stat records for the given post IDs,
// preserving input order. Only fetches count columns (no content).
func GetPostStatsByIDs(ids []int) ([]model.Posts, error) {
	if len(ids) == 0 {
		return []model.Posts{}, nil
	}

	query := `
		SELECT p.id, p.user_id, u.username, p.title,
			p.praise_count, p.comment_count, p.collection_count, p.view_count,
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

	var rows []model.Posts
	if err := common.DB.Select(&rows, query, args...); err != nil {
		return nil, err
	}

	// Re-sort rows to match input ID order (deduplicated: only first occurrence gets the row)
	idIndex := make(map[int]int, len(ids))
	seen := make(map[int]bool, len(ids))
	for i, id := range ids {
		if !seen[id] {
			idIndex[id] = i
			seen[id] = true
		}
	}
	posts := make([]model.Posts, len(ids))
	written := 0
	for _, row := range rows {
		if idx, ok := idIndex[row.ID]; ok {
			posts[idx] = row
			written++
		}
	}
	return posts[:written], nil
}

// GetPostsByUserID returns posts authored by the given user, newest first, with pagination.
func GetPostsByUserID(userID int, page, pageSize int) ([]model.Posts, int, error) {
	var posts []model.Posts
	var total int

	// Count total
	if err := common.DB.Get(&total, `SELECT COUNT(1) FROM posts WHERE user_id = $1`, userID); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := common.DB.Select(&posts, `
		SELECT p.id, p.user_id, u.username, p.title,
			p.praise_count, p.comment_count, p.collection_count, p.view_count,
			p.created_time
		FROM posts as p
		LEFT JOIN users as u ON p.user_id = u.id
		WHERE p.user_id = $1
		ORDER BY created_time DESC
		LIMIT $2 OFFSET $3
	`, userID, pageSize, offset)

	return posts, total, err
}

// IncrPraiseCount atomically increments the praise_count for a post.
func IncrPraiseCount(postID int) error {
	_, err := common.DB.Exec(`UPDATE posts SET praise_count = praise_count + 1 WHERE id = $1`, postID)
	return err
}

// DecrPraiseCount atomically decrements the praise_count for a post (floor 0).
func DecrPraiseCount(postID int) error {
	_, err := common.DB.Exec(`UPDATE posts SET praise_count = GREATEST(praise_count - 1, 0) WHERE id = $1`, postID)
	return err
}

// IncrCommentCount atomically increments the comment_count for a post.
func IncrCommentCount(postID int) error {
	_, err := common.DB.Exec(`UPDATE posts SET comment_count = comment_count + 1 WHERE id = $1`, postID)
	return err
}

// DecrCommentCount atomically decrements the comment_count for a post (floor 0).
func DecrCommentCount(postID int) error {
	_, err := common.DB.Exec(`UPDATE posts SET comment_count = GREATEST(comment_count - 1, 0) WHERE id = $1`, postID)
	return err
}

// IncrCollectionCount atomically increments the collection_count for a post.
func IncrCollectionCount(postID int) error {
	_, err := common.DB.Exec(`UPDATE posts SET collection_count = collection_count + 1 WHERE id = $1`, postID)
	return err
}

// DecrCollectionCount atomically decrements the collection_count for a post (floor 0).
func DecrCollectionCount(postID int) error {
	_, err := common.DB.Exec(`UPDATE posts SET collection_count = GREATEST(collection_count - 1, 0) WHERE id = $1`, postID)
	return err
}

// IncrViewCount atomically increments the view_count for a post.
func IncrViewCount(postID int) error {
	_, err := common.DB.Exec(`UPDATE posts SET view_count = view_count + 1 WHERE id = $1`, postID)
	return err
}
