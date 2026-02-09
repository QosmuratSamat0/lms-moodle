package middleware

import (
	"context"
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/domain/categorymanager"
	"github.com/gin-gonic/gin"
)

// RequireCategoryManagerPermission checks if user has required permission level for a category
// Permission levels: "view" (read), "edit" (read+update), "admin" (full access)
func RequireCategoryManagerPermission(repo categorymanager.Repository, requiredLevel string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user id required"})
			c.Abort()
			return
		}

		// Get categoryID from route params
		categoryID := c.Param("id")
		// For some routes, it might be in :categoryID param
		if categoryID == "" {
			categoryID = c.Param("categoryID")
		}
		// For POST, it's in the request body
		if categoryID == "" {
			if err := c.ShouldBindJSON(nil); err == nil {
				if id, exists := c.Get("category_id"); exists {
					if idStr, ok := id.(string); ok {
						categoryID = idStr
					}
				}
			}
		}

		if categoryID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "category id required"})
			c.Abort()
			return
		}

		// Check if user has access to this category with required permission level
		hasPermission, err := checkCategoryManagerPermission(repo, userID, categoryID, requiredLevel)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check permissions"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "insufficient permission level for this category",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkCategoryManagerPermission verifies if user has required permission for category
func checkCategoryManagerPermission(repo categorymanager.Repository, userID, categoryID, requiredLevel string) (bool, error) {
	// Get all category managers for this user in this category
	mgr, err := repo.GetByUserAndCategory(context.Background(), userID, categoryID)
	if err != nil {
		return false, err
	}

	if mgr == nil {
		return false, nil
	}

	return hasRequiredLevel(mgr.PermissionLevel, requiredLevel), nil
}

// hasRequiredLevel checks if userLevel satisfies requiredLevel
func hasRequiredLevel(userLevel, requiredLevel string) bool {
	// Permission hierarchy: view < edit < admin
	levelMap := map[string]int{
		"view":  1,
		"edit":  2,
		"admin": 3,
	}

	userLvl := levelMap[userLevel]
	requiredLvl := levelMap[requiredLevel]

	return userLvl >= requiredLvl
}

// boolPtr helper
func boolPtr(v bool) *bool {
	return &v
}
