package common

import "github.com/gin-gonic/gin"

const (
	HTTPOK                  = 200
	HTTPBadRequest          = 400
	HTTPUnauthorized        = 401
	HTTPForbidden           = 403
	HTTPNotFound            = 404
	HTTPConflict            = 409
	HTTPInternalServerError = 500
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Success(c *gin.Context, data any) {
	c.JSON(HTTPOK, Response{
		Code:    HTTPOK,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, httpCode int, err *AppError) {
	c.JSON(httpCode, gin.H{
		"code":    err.Code,
		"message": err.Message,
	})
}
