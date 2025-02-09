package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware ensures the user has at least one of the allowed roles to access that route.
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No user role found in token"})
			c.Abort()
			return
		}
		role := userRole.(string)

		for _, r := range allowedRoles {
			if r == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}
