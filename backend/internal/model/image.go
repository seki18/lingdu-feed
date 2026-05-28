package model

// PostImage represents a single image attached to a post.
type PostImage struct {
	ID        int    `db:"id" json:"id"`
	PostID    int    `db:"post_id" json:"post_id"`
	ImageURL  string `db:"image_url" json:"image_url"`
	SortOrder int    `db:"sort_order" json:"sort_order"`
}

// AddPostImagesRequest is the JSON body for POST /posts/:id/images.
type AddPostImagesRequest struct {
	Images []string `json:"images"`
}
