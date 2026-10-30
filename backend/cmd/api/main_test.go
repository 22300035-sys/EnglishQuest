package main

import (
	"context"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/englishquest/backend/internal/config"
	"github.com/englishquest/backend/internal/db"
	"github.com/englishquest/backend/internal/models"
	"github.com/englishquest/backend/internal/repositories"
	"github.com/englishquest/backend/internal/routes"
)

// Only point TEST_DATABASE_URL at a disposable development database. This test
// applies additive migrations and inserts test users. It never drops data.
func TestPostgresIntegration(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not configured; PostgreSQL integration not executed")
	}
	database, err := db.Open(dsn)
	if err != nil {
		t.Fatal("cannot connect to integration database")
	}
	defer db.Close(database)
	if err := runMigrations(database); err != nil {
		t.Fatal(err)
	}
	if err := runMigrations(database); err != nil {
		t.Fatalf("migration rerun: %v", err)
	}
	users := repositories.NewUserRepository(database)
	// UUID-based email permits repeated non-destructive test runs.
	var id string
	if err := database.Raw("SELECT gen_random_uuid()").Scan(&id).Error; err != nil {
		t.Fatal(err)
	}
	user := &models.User{Email: "integration-" + id + "@example.org", PasswordHash: "not-a-usable-password-hash", FirstName: "Test", CEFRLevel: "A1", Role: "STUDENT", IsActive: true}
	if err := users.Create(context.Background(), user); err != nil {
		t.Fatal(err)
	}
	found, err := users.FindByEmail(context.Background(), user.Email)
	if err != nil || found.ID != user.ID {
		t.Fatal("user did not persist")
	}
	duplicate := *user
	duplicate.ID = ""
	if err := users.Create(context.Background(), &duplicate); err == nil {
		t.Fatal("duplicate email allowed")
	}
	invalid := *user
	invalid.ID = ""
	invalid.Email = "invalid-" + user.Email
	invalid.Role = "SUPERADMIN"
	if err := users.Create(context.Background(), &invalid); err == nil {
		t.Fatal("invalid role allowed")
	}
	r := routes.NewRouter(&config.Config{AllowedOrigins: []string{"http://localhost:5173"}}, database)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/api/v1/ready", nil))
	if w.Code != 200 || !strings.Contains(w.Body.String(), "ready") {
		t.Fatalf("readiness failed: %s", w.Body.String())
	}
	sqlDB, _ := database.DB()
	_ = sqlDB.Close()
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/api/v1/ready", nil))
	if w.Code != 503 {
		t.Fatalf("closed DB should not report ready: %d", w.Code)
	}
}
