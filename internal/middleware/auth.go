package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/distroaryan/restaurant-management/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		var userId string
		if id, ok := claims["userId"].(string); ok {
			userId = id
		} else if id, ok := claims["sub"].(string); ok {
			userId = id
		} else if id, ok := claims["id"].(string); ok {
			userId = id
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "userId not found in token"})
			return
		}

		// Add retrieved userId into the request headers and context
		c.Request.Header.Set("X-User-Id", userId)
		c.Set("userId", userId)

		c.Next()
	}
}
