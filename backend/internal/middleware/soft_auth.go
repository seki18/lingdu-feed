package middleware

import (
	"strings"

	"github.com/seki18/lingdu-feed/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// SoftAuthMiddleware optionally parses a JWT token and stores user_id in the context.
// Unlike AuthMiddleware, missing or invalid tokens are silently ignored (user_id = -1).
func SoftAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		userID := -1 // default: not logged in

		if authHeader != "" {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.ParseWithClaims(tokenString, &utils.Claims{},
				func(token *jwt.Token) (any, error) {
					return []byte("your-secret-key"), nil
				})
			if err == nil && token.Valid {
				claims := token.Claims.(*utils.Claims)
				userID = claims.UserID
			}
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
