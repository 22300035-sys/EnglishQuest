package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/englishquest/backend/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Health does not query the DB. Database readiness has a real integration test.
func testRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return NewRouter(&config.Config{AllowedOrigins: []string{"http://localhost:5173"}}, &gorm.DB{})
}

func TestHealth(t *testing.T) {
	w := httptest.NewRecorder()
	testRouter().ServeHTTP(w, httptest.NewRequest("GET", "/api/v1/health", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	var body struct {
		Data struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil || body.Data.Status != "ok" {
		t.Fatal("invalid health response")
	}
	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatal("missing security header")
	}
}

func TestRoutingErrors(t *testing.T) {
	for _, tc := range []struct {
		method, path string
		status       int
	}{
		{"GET", "/api/v1/missing", 404}, {"POST", "/api/v1/health", 405},
	} {
		w := httptest.NewRecorder()
		testRouter().ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
		if w.Code != tc.status || !json.Valid(w.Body.Bytes()) {
			t.Fatalf("unexpected response: %d %s", w.Code, w.Body.String())
		}
	}
}

func TestCORS(t *testing.T) {
	for _, tc := range []struct {
		origin string
		status int
	}{
		{"http://localhost:5173", 204}, {"https://untrusted.example.org", 403},
	} {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("OPTIONS", "/api/v1/health", nil)
		r.Header.Set("Origin", tc.origin)
		r.Header.Set("Access-Control-Request-Method", "GET")
		testRouter().ServeHTTP(w, r)
		if w.Code != tc.status {
			t.Fatalf("origin %s: status %d", tc.origin, w.Code)
		}
		if tc.status == 403 && w.Header().Get("Access-Control-Allow-Origin") != "" {
			t.Fatal("untrusted origin allowed")
		}
	}
}
