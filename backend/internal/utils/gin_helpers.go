package utils

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ParseExcludeIDs extracts and parses the "current_ids" comma-separated query
// parameter into a slice of ints. Invalid or non-positive values are silently skipped.
func ParseExcludeIDs(c *gin.Context) []int {
	var ids []int
	raw := c.Query("current_ids")
	if raw == "" {
		return ids
	}
	for _, s := range strings.Split(raw, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		id, err := strconv.Atoi(s)
		if err != nil || id <= 0 {
			continue
		}
		ids = append(ids, id)
	}
	return ids
}

// GetSoftUserID returns the user ID set by SoftAuthMiddleware, or -1 if
// the user is not logged in.
func GetSoftUserID(c *gin.Context) int {
	uid, _ := c.Get("user_id")
	if v, ok := uid.(int); ok {
		return v
	}
	return -1
}

// GetAuthUserID returns the user ID set by AuthMiddleware. It panics if
// user_id is missing, which should never happen when the middleware is active.
func GetAuthUserID(c *gin.Context) int {
	uid, _ := c.Get("user_id")
	return uid.(int)
}

// ParsePagination extracts page and page_size from query parameters with
// defaults of 1 and 10 respectively.
func ParsePagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return page, pageSize
}

// GetQueryInt reads an integer query parameter with a default value.
func GetQueryInt(c *gin.Context, key string, defaultVal int) int {
	val, err := strconv.Atoi(c.DefaultQuery(key, strconv.Itoa(defaultVal)))
	if err != nil {
		return defaultVal
	}
	return val
}
