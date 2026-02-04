// internal/api/middleware.go
package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupMiddleware configures global middleware for the router
func SetupMiddleware(engine *gin.Engine) {
	// Recovery middleware
	engine.Use(gin.Recovery())

	// Logger middleware
	engine.Use(gin.Logger())

	// CORS middleware
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

// RequestLogger is a custom logger middleware
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)

		if raw != "" {
			path = path + "?" + raw
		}

		statusCode := c.Writer.Status()

		// Log format: status | latency | method | path
		gin.DefaultWriter.Write([]byte(
			c.ClientIP() + " | " +
				time.Now().Format(time.RFC3339) + " | " +
				string(rune(statusCode)) + " | " +
				latency.String() + " | " +
				c.Request.Method + " | " +
				path + "\\n",
		))
	}
}
