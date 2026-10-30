package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"

	"github.com/englishquest/backend/internal/config"
	"github.com/englishquest/backend/internal/handler"
	"github.com/englishquest/backend/internal/middleware"
)

// NewRouter builds the versioned REST API router (base path /api/v1).
func NewRouter(cfg *config.Config, database *gorm.DB) *gin.Engine {
	if database == nil {
		panic("db is required to create router")
	}
	if cfg == nil {
		panic("config is required to create router")
	}

	r := gin.New()
	_ = r.SetTrustedProxies(nil)
	r.HandleMethodNotAllowed = true
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORSMiddleware(cfg.AllowedOrigins))
	r.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Cache-Control", "no-store")
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
		c.Next()
	})
	r.NoRoute(func(c *gin.Context) { handler.WriteError(c, http.StatusNotFound, "route not found") })
	r.NoMethod(func(c *gin.Context) { handler.WriteError(c, http.StatusMethodNotAllowed, "method not allowed") })

	api := r.Group("/api")
	v1 := api.Group("/v1")

	h := handler.NewHandler()
	h.RegisterRoutes(v1)
	v1.GET("/ready", handler.Readiness(database))
	handler.NewApp(cfg, database).Routes(v1)

	return r
}
