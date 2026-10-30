package handler

import (
	"github.com/gin-gonic/gin"
)

// Handler groups the HTTP handlers of the application and takes care of
// registering them in the versioned API router.
type Handler struct {
	Health *HealthHandler
}

func NewHandler() *Handler {
	return &Handler{Health: NewHealthHandler()}
}

func (h *Handler) RegisterHealth(r *gin.RouterGroup) {
	r.GET("/health", h.Health.Health)
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	h.RegisterHealth(r)
}
