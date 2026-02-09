package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ContextMiddleware sets common context values for authenticated users
func ContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// These values should be set by actual authentication
		// This is a placeholder showing how they would be used
		userID := c.GetString("userID")
		if userID == "" {
			// Try to get from header for testing
			userID = c.GetHeader("X-User-ID")
		}
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user id required"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

// CORSMiddleware for allowing cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("requestID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// ErrorHandlingMiddleware handles panic and unexpected errors
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
			}
		}()
		c.Next()
	}
}

// LoggingMiddleware logs request information
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log request details
		method := c.Request.Method
		path := c.Request.URL.Path
		userID, _ := GetUserIDFromContext(c)

		c.Next()

		// Log response details
		statusCode := c.Writer.Status()
		_ = method
		_ = path
		_ = userID
		_ = statusCode
		// In production, use proper logging library
	}
}
