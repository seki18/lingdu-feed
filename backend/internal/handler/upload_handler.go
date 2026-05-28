package handler

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/storage"

	"github.com/gin-gonic/gin"
)

// UploadImage handles POST /api/upload (auth required).
// Accepts multipart form with field "file" and optional "post_id".
func UploadImage(c *gin.Context) {
	// Validate file presence
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Printf("[UploadImage] No file in request: %v", err)
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}

	// Validate file size (10 MB max)
	if fileHeader.Size > storage.ImageMaxSize {
		log.Printf("[UploadImage] File too large: %d bytes", fileHeader.Size)
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowed[ext] {
		log.Printf("[UploadImage] Unsupported file type: %s", ext)
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}

	// Post ID (optional — use "general" if not specified)
	postID := c.PostForm("post_id")
	if postID == "" {
		postID = "general"
	}

	// Upload to S3 with compression
	url, err := storage.UploadImage(
		c.Request.Context(),
		postID,
		fileHeader,
		1920, // maxWidth — could read from config
		85,   // jpegQuality
	)
	if err != nil {
		log.Printf("[UploadImage] Upload failed: %v", err)
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

	log.Printf("[UploadImage] Success: post=%s file=%s url=%s", postID, fileHeader.Filename, url)
	common.Success(c, gin.H{"url": url})
}
