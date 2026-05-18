package common

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the standard JSON response envelope for all API endpoints.
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Success sends a 200 OK response with the provided data.
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

// Error logs the original error (if any) to the backend terminal and sends
// a user-friendly JSON error response to the frontend.
func Error(c *gin.Context, httpCode int, err *AppError) {
	log.Printf("[ERROR %d] code=%d message=%q", httpCode, int(err.Code), err.Message)
	if err.Err != nil {
		log.Printf("[ERROR %d] cause: %v", httpCode, err.Err)
	}
	c.JSON(httpCode, Response{
		Code:    int(err.Code),
		Message: err.Message,
	})
}

// PaginatedResponse wraps a paginated list with total count, page and page size.
type PaginatedResponse struct {
	Items    any `json:"items"`
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// SuccessPaginated sends a 200 OK paginated response.
func SuccessPaginated(c *gin.Context, items any, total, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "success",
		Data: PaginatedResponse{
			Items:    items,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}
