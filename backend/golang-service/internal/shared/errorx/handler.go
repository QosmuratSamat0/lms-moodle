package errorx

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool       `json:"success"`
	Error   *ErrorInfo `json:"error"`
}

// ErrorInfo contains error details
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// HandleError is a Gin middleware-friendly error handler
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// Check for AppError
	var appErr *AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.StatusCode, ErrorResponse{
			Success: false,
			Error: &ErrorInfo{
				Code:    appErr.Code,
				Message: appErr.Message,
			},
		})
		return
	}

	// Check for pgx.ErrNoRows
	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Success: false,
			Error: &ErrorInfo{
				Code:    "NOT_FOUND",
				Message: "resource not found",
			},
		})
		return
	}

	// Log unexpected errors
	log.Printf("Unexpected error: %v", err)

	// Default to internal server error
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Success: false,
		Error: &ErrorInfo{
			Code:    "INTERNAL_ERROR",
			Message: "an unexpected error occurred",
		},
	})
}

// AbortWithError aborts the request with an error response
func AbortWithError(c *gin.Context, err error) {
	HandleError(c, err)
	c.Abort()
}
