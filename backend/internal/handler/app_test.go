package handler

import (
	"github.com/englishquest/backend/internal/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProductionCookieFlags(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	app := &App{Config: &config.Config{Env: "production"}}
	app.cookie(c, "opaque-test-token", 3600)
	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatal("missing cookie")
	}
	cookie := cookies[0]
	if cookie.Name != "__Host-eq_session" || !cookie.HttpOnly || !cookie.Secure || !cookie.Partitioned || cookie.SameSite != http.SameSiteNoneMode || cookie.Path != "/" || cookie.Domain != "" {
		t.Fatalf("unsafe production cookie flags: %+v", cookie)
	}
}
func TestAuthRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	app := NewApp(&config.Config{}, nil)
	r := gin.New()
	r.POST("/login", app.RateLimit, func(c *gin.Context) { c.Status(204) })
	for i := 0; i < 16; i++ {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("POST", "/login", nil))
		expected := 204
		if i == 15 {
			expected = 429
		}
		if w.Code != expected {
			t.Fatalf("request %d status %d", i, w.Code)
		}
	}
}
