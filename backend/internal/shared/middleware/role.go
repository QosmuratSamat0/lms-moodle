// internal/shared/middleware/role.go
package middleware

import (
	"github.com/gin-gonic/gin"
)

// RequireAdmin is a convenience middleware for admin-only routes
func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin")
}

// RequireTeacherOrAdmin is a convenience middleware for teacher or admin routes
func RequireTeacherOrAdmin() gin.HandlerFunc {
	return RequireRole("teacher", "admin")
}

// RequireManagerOrAdmin is a convenience middleware for manager or admin routes
func RequireManagerOrAdmin() gin.HandlerFunc {
	return RequireRole("manager", "admin")
}

// RequireStudent is a convenience middleware for student-only routes
func RequireStudent() gin.HandlerFunc {
	return RequireRole("student")
}

// RequireTeacher is a convenience middleware for teacher-only routes
func RequireTeacher() gin.HandlerFunc {
	return RequireRole("teacher")
}
