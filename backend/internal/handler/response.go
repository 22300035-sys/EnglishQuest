package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/englishquest/backend/internal/pkg/errors"
)

// WriteJSON writes a JSON response using the conventional envelope
// {"data": ...} for successful responses and {"error": "..."} for failures.
func WriteJSON(c *gin.Context, status int, data interface{}) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(status, gin.H{"data": data})
}

// WriteError writes a JSON error response using the envelope {"error": "..."}.
func WriteError(c *gin.Context, status int, message string) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(status, gin.H{"error": message})
}

// RespondError writes the HTTP response mapped from the given error.
func RespondError(c *gin.Context, err error) {
	if c.Request != nil && c.Request.Context().Err() == context.Canceled {
		c.AbortWithStatus(499) // Client closed the request; not an application 500.
		return
	}
	if c.Request != nil && c.Request.Context().Err() == context.DeadlineExceeded {
		c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{"error": "La operación tardó demasiado. Inténtalo de nuevo."})
		return
	}
	appErr, ok := errors.IsAppError(err)
	if !ok {
		appErr = &errors.AppError{StatusCode: http.StatusInternalServerError, Message: "internal server error"}
	}
	c.AbortWithStatusJSON(appErr.StatusCode, gin.H{"error": appErr.Message})
}
