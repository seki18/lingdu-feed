package common

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Success(c *gin.Context, data any) {
	fmt.Println(data)
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, httpCode int, err *AppError) {
	c.JSON(httpCode, Response{
		Code:    int(err.Code),
		Message: err.Message,
		Data:    nil,
	})
}
