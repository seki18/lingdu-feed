package storage

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/seki18/lingdu-feed/internal/common"
	"golang.org/x/image/draw"
)

// ImageMaxSize is the maximum allowed file size (10 MB).
const ImageMaxSize = 10 << 20

// UploadImage reads an uploaded file, compresses it, and uploads to S3.
// Returns the public URL of the uploaded image.
func UploadImage(
	ctx context.Context,
	postID string,
	fileHeader *multipart.FileHeader,
	maxWidth int,
	jpegQuality int,
) (string, error) {
	if common.S3Client == nil {
		return "", fmt.Errorf("S3 client not initialized")
	}

	// Determine filename with timestamp prefix to avoid collisions
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// Read uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	// Compress image
	compressed, contentType, err := compressImage(file, ext, maxWidth, jpegQuality)
	if err != nil {
		return "", err
	}

	// Build S3 key: posts/{postId}/{filename}
	s3Key := fmt.Sprintf("posts/%s/%s", postID, filename)

	// Upload to S3
	if _, err := common.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(common.S3Bucket),
		Key:         aws.String(s3Key),
		Body:        bytes.NewReader(compressed),
		ContentType: aws.String(contentType),
	}); err != nil {
		return "", fmt.Errorf("S3 upload failed: %w", err)
	}

	// Build the public URL: https://{bucket}.s3.{region}.amazonaws.com/{key}
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", common.S3Bucket, common.S3Region, s3Key)
	return url, nil
}

// compressImage reads an image, resizes it to at most maxWidth, and encodes as JPEG or PNG.
func compressImage(src io.Reader, ext string, maxWidth int, quality int) ([]byte, string, error) {
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()

	// Resize only if wider than maxWidth
	if width > maxWidth {
		newHeight := bounds.Dy() * maxWidth / width
		dst := image.NewRGBA(image.Rect(0, 0, maxWidth, newHeight))
		draw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)
		img = dst
	}

	var buf bytes.Buffer
	switch ext {
	case ".png":
		if err := png.Encode(&buf, img); err != nil {
			return nil, "", fmt.Errorf("failed to encode PNG: %w", err)
		}
		return buf.Bytes(), "image/png", nil
	default:
		// JPEG for .jpg, .jpeg, and fallback
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
			return nil, "", fmt.Errorf("failed to encode JPEG: %w", err)
		}
		return buf.Bytes(), "image/jpeg", nil
	}
}
