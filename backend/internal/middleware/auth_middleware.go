package middleware

import (
	"FreeConnect/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks for Bearer token, validates it, and sets user info in context.
func AuthMiddleware(jwtService services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
			c.Abort()
			return
		}
		tokenString := parts[1]

		token, err := jwtService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*services.CustomClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Store user info in context for downstream
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.UserRole)

		c.Next()
	}
}
