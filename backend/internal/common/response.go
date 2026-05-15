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
	if err.Err != nil {
		log.Printf("[ERROR %d] %s | Original error: %v", err.Code, err.Message, err.Err)
	} else {
		log.Printf("[ERROR %d] %s", err.Code, err.Message)
	}

	c.JSON(httpCode, Response{
		Code:    int(err.Code),
		Message: err.Message,
		Data:    nil,
	})
}
