package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole middleware checks if the user has one of the allowed roles
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		role, ok := v.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid role type"})
			c.Abort()
			return
		}

		for _, r := range allowedRoles {
			if role == r {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		c.Abort()
	}
}

// RequirePermission middleware checks if the user has the required permission
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, exists := c.Get("permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "permissions not found"})
			c.Abort()
			return
		}

		perms, ok := v.([]string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid permissions format"})
			c.Abort()
			return
		}

		for _, p := range perms {
			if p == permission || p == "*" {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": fmt.Sprintf("permission '%s' required", permission),
		})
		c.Abort()
	}
}

// RequireAccessLevel checks if the user has the required access level
func RequireAccessLevel(requiredLevel string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, exists := c.Get("accessLevel")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "access level not found"})
			c.Abort()
			return
		}

		accessLevel, ok := v.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid access level type"})
			c.Abort()
			return
		}

		// Super admin has full access
		if accessLevel == "super_admin" {
			c.Next()
			return
		}

		if accessLevel != requiredLevel {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient access level"})
			c.Abort()
			return
		}

		c.Next()
	}
}
