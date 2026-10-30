package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecoveryDoesNotExposePanic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Recovery())
	r.GET("/panic", func(c *gin.Context) { panic("sensitive-database-credentials") })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/panic", nil))
	if w.Code != 500 || strings.Contains(w.Body.String(), "sensitive") {
		t.Fatalf("unsafe response: %s", w.Body.String())
	}
}
