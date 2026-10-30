package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger writes structured-ish access logs (method, path, status, latency).
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		attrs := []any{
			"method", c.Request.Method,
			"route", c.FullPath(),
			"status", c.Writer.Status(),
			"latency_ms", float64(time.Since(start).Microseconds()) / 1000.0,
		}
		slog.Info("http_request", attrs...)
	}
}

// Recovery converts panics into JSON 500 responses instead of crashing the
// engine or returning an HTML error page.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered", "method", c.Request.Method, "route", c.FullPath())
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
		}()
		c.Next()
	}
}
