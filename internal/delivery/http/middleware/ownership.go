package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireOwnership middleware checks if the user owns the resource
func RequireOwnership(resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// super admin bypass
		if isSuperAdmin(c) {
			c.Next()
			return
		}

		ownerID := getOwnerIDFromContext(c, resourceType)
		if ownerID == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "owner not resolved"})
			c.Abort()
			return
		}

		if ownerID != userID.(string) {
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not own this resource"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// isSuperAdmin helper function to check if user is super admin
func isSuperAdmin(c *gin.Context) bool {
	v, exists := c.Get("accessLevel")
	if !exists {
		return false
	}
	level, ok := v.(string)
	return ok && level == "super_admin"
}

// getOwnerIDFromContext helper to extract owner ID based on resource type
func getOwnerIDFromContext(c *gin.Context, resourceType string) string {
	switch resourceType {
	case "course":
		return c.GetString("courseOwnerID")
	case "submission":
		return c.GetString("submissionOwnerID")
	case "group":
		return c.GetString("groupOwnerID")
	default:
		return ""
	}
}
