package handler

import (
	"context"
	"gorm.io/gorm"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles infrastructure and liveness probes used by Render
// and local development.
type HealthHandler struct{}

func NewHealthHandler() *HealthHandler { return &HealthHandler{} }

func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"status": "ok"}})
}

// RegisterHealth registers the health check endpoint under the given
// versioned API group (e.g. /api/v1).
func (h *HealthHandler) RegisterHealth(r *gin.RouterGroup) {
	r.GET("/health", h.Health)
}

// Readiness checks the actual database without disclosing connection details.
func Readiness(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		sqlDB, err := database.DB()
		if err != nil {
			WriteError(c, http.StatusServiceUnavailable, "database unavailable")
			return
		}
		if err := sqlDB.PingContext(ctx); err != nil {
			WriteError(c, http.StatusServiceUnavailable, "database unavailable")
			return
		}
		WriteJSON(c, http.StatusOK, gin.H{"status": "ready"})
	}
}
