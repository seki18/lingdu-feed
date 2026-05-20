package model

// PostDetailResponse is the full response for a single post detail page,
// including the post content, interaction status, and comments.
type PostDetailResponse struct {
	Post         Post      `json:"post"`
	HasLiked     bool      `json:"has_liked"`
	HasFavorited bool      `json:"has_favorited"`
	Comments     []Comment `json:"comments"`
}
