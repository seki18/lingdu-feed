package middleware

import (
	"community-backend/internal/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware returns a Gin middleware that validates JWT tokens.
// It extracts the user ID from the token and stores it in the context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no token"})
			c.Abort()
			log.Printf("AuthMiddleware: missing Authorization header")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (any, error) {
			return []byte("your-secret-key"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			log.Printf("AuthMiddleware: invalid token (%v)", err)
			return
		}

		claims := token.Claims.(*utils.Claims)
		// Store user ID in context for downstream handlers
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}
