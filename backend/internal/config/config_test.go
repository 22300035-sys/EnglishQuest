package config

import (
	"strings"
	"testing"
)

func baseEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{"APP_ENV", "PORT", "DATABASE_URL", "JWT_SECRET", "JWT_ACCESS_TTL", "JWT_REFRESH_TTL", "CORS_ALLOWED_ORIGINS"} {
		t.Setenv(key, "")
	}
	t.Setenv("DATABASE_URL", "postgres://test:password@localhost:5432/test?sslmode=disable")
	t.Setenv("JWT_SECRET", strings.Repeat("x", 32))
}

func TestDevelopmentDefaults(t *testing.T) {
	baseEnv(t)
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.Port != "8080" || c.Env != "development" || len(c.AllowedOrigins) != 2 {
		t.Fatalf("unexpected defaults: port=%s env=%s", c.Port, c.Env)
	}
}

func TestInvalidConfiguration(t *testing.T) {
	for _, tc := range []struct{ name, key, value string }{
		{"missing database", "DATABASE_URL", ""},
		{"invalid database", "DATABASE_URL", "not-a-url"},
		{"wrong driver", "DATABASE_URL", "mysql://localhost/db"},
		{"duplicate connection parameters", "DATABASE_URL", "postgres://localhost/test?sslmode=require&sslmode=disable"},
		{"short secret", "JWT_SECRET", "short"},
		{"invalid port", "PORT", "70000"},
		{"invalid environment", "APP_ENV", "prod"},
		{"invalid access duration", "JWT_ACCESS_TTL", "bad"},
		{"negative access duration", "JWT_ACCESS_TTL", "-1m"},
		{"oversized access duration", "JWT_ACCESS_TTL", "24h"},
		{"invalid refresh duration", "JWT_REFRESH_TTL", "bad"},
		{"short refresh duration", "JWT_REFRESH_TTL", "1m"},
		{"wildcard origin", "CORS_ALLOWED_ORIGINS", "https://*.example.org"},
		{"origin path", "CORS_ALLOWED_ORIGINS", "https://example.org/path"},
		{"origin credentials", "CORS_ALLOWED_ORIGINS", "https://user:pass@example.org"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			baseEnv(t)
			t.Setenv(tc.key, tc.value)
			if _, err := Load(); err == nil {
				t.Fatal("expected validation failure")
			}
		})
	}
}

func TestProductionConfiguration(t *testing.T) {
	baseEnv(t)
	t.Setenv("APP_ENV", "production")
	if _, err := Load(); err == nil {
		t.Fatal("expected SSL rejection")
	}
	t.Setenv("DATABASE_URL", "postgres://test:password@db.example.org:5432/test?sslmode=require")
	if _, err := Load(); err == nil {
		t.Fatal("expected explicit CORS origins")
	}
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")
	if _, err := Load(); err == nil {
		t.Fatal("expected production origin rejection")
	}
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://englishquest.example.org")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(c.AllowedOrigins) != 1 {
		t.Fatal("development origins must not leak into production")
	}
}
