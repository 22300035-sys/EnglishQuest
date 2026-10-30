package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"gorm.io/gorm"

	"github.com/englishquest/backend/internal/config"
	"github.com/englishquest/backend/internal/db"
	"github.com/englishquest/backend/internal/routes"

	"github.com/englishquest/backend/migrations"
)

func main() {
	applyMigrations := flag.Bool("migrate", false, "Apply reviewed forward migrations and exit (explicit operation)")
	promoteAdmin := flag.String("promote-admin", "", "Explicitly grant ADMIN to an existing account email and exit; owner operation")
	flag.Parse()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	gormDB, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("database connection failed; verify credentials, SSL and availability")
	}
	defer db.Close(gormDB)
	if *promoteAdmin != "" {
		res := gormDB.Exec("UPDATE users SET role='ADMIN',updated_at=now() WHERE email=? AND is_active AND deleted_at IS NULL", strings.ToLower(strings.TrimSpace(*promoteAdmin)))
		if res.Error != nil || res.RowsAffected != 1 {
			log.Fatal("admin promotion failed: account must exist and be active")
		}
		log.Println("existing account granted ADMIN role")
		return
	}

	if *applyMigrations {
		if err := runMigrations(gormDB); err != nil {
			log.Fatal("migration failed; inspect schema version/dirty state before retrying")
		}
		log.Println("migrations applied (or already current)")
		return
	}

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := routes.NewRouter(cfg, gormDB)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}
	log.Println("server stopped")
}

func runMigrations(gormDB *gorm.DB) error {
	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}
	source, err := iofs.New(migrations.FS, "sql")
	if err != nil {
		return err
	}
	databaseDriver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return err
	}
	instance, err := migrate.NewWithInstance("iofs", source, "postgres", databaseDriver)
	if err != nil {
		return err
	}
	if err := instance.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
