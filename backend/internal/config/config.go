package config

import (
	"errors"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Env            string
	Port           string
	DatabaseURL    string
	JWTSecret      string
	JWTAccessTTL   time.Duration
	JWTRefreshTTL  time.Duration
	AllowedOrigins []string
}

func Load() (*Config, error) {
	env := strings.ToLower(os.Getenv("APP_ENV"))
	if env == "" {
		env = "development"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if n, err := strconv.Atoi(port); err != nil || n < 1 || n > 65535 {
		return nil, errors.New("PORT must be between 1 and 65535")
	}
	if env != "development" && env != "test" && env != "production" {
		return nil, errors.New("APP_ENV must be development, test or production")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	dbURL, err := url.Parse(databaseURL)
	if err != nil || (dbURL.Scheme != "postgres" && dbURL.Scheme != "postgresql") || dbURL.Hostname() == "" || dbURL.Path == "" || dbURL.Path == "/" {
		return nil, errors.New("DATABASE_URL must be a PostgreSQL URL with host and database")
	}
	query, err := url.ParseQuery(dbURL.RawQuery)
	if err != nil {
		return nil, errors.New("DATABASE_URL contains an invalid query")
	}
	for _, values := range query {
		if len(values) != 1 {
			return nil, errors.New("DATABASE_URL must not contain duplicate connection parameters")
		}
	}
	if env == "production" {
		ssl := query.Get("sslmode")
		if ssl != "require" && ssl != "verify-ca" && ssl != "verify-full" {
			return nil, errors.New("production DATABASE_URL requires sslmode=require or certificate verification")
		}
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if len(strings.TrimSpace(jwtSecret)) < 32 {
		return nil, errors.New("JWT_SECRET must contain at least 32 characters")
	}

	accessTTL := 15 * time.Minute
	if v := os.Getenv("JWT_ACCESS_TTL"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil || d <= 0 || d > time.Hour {
			return nil, errors.New("JWT_ACCESS_TTL must be positive and at most 1h")
		}
		accessTTL = d
	}

	refreshTTL := 72 * time.Hour
	if v := os.Getenv("JWT_REFRESH_TTL"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil || d <= accessTTL || d > 30*24*time.Hour {
			return nil, errors.New("JWT_REFRESH_TTL must exceed access TTL and be at most 720h")
		}
		refreshTTL = d
	}

	origins := []string{}
	for _, o := range strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",") {
		if o = strings.TrimSpace(o); o != "" {
			u, err := url.Parse(o)
			if err != nil || u.Hostname() == "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" || strings.Contains(o, "*") || (u.Scheme != "http" && u.Scheme != "https") {
				return nil, errors.New("CORS_ALLOWED_ORIGINS must contain exact HTTP origins without paths")
			}
			if env == "production" && (u.Scheme != "https" || u.Hostname() == "localhost" || net.ParseIP(u.Hostname()).IsLoopback()) {
				return nil, errors.New("production CORS origins must use HTTPS and must not be localhost")
			}
			origins = append(origins, o)
		}
	}
	if len(origins) == 0 {
		if env == "production" {
			return nil, errors.New("CORS_ALLOWED_ORIGINS is required in production")
		}
		origins = []string{"http://localhost:5173", "http://localhost:4173"}
	}

	return &Config{
		Env:            env,
		Port:           port,
		DatabaseURL:    databaseURL,
		JWTSecret:      jwtSecret,
		JWTAccessTTL:   accessTTL,
		JWTRefreshTTL:  refreshTTL,
		AllowedOrigins: origins,
	}, nil
}
