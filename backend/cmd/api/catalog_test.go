package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/englishquest/backend/internal/db"
	"github.com/englishquest/backend/internal/models"
	apperr "github.com/englishquest/backend/internal/pkg/errors"
	"github.com/englishquest/backend/internal/repositories"
	"github.com/englishquest/backend/internal/services"
	"gorm.io/gorm"
	"os"
	"testing"
	"time"
)

func TestCatalogueAndPasswordTransactions(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not configured")
	}
	database, err := db.Open(dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close(database)
	if err := runMigrations(database); err != nil {
		t.Fatal(err)
	}
	rollback := errors.New("intentional test rollback")
	err = database.Transaction(func(tx *gorm.DB) error {
		ctx := context.Background()
		auth := services.Auth{Accounts: repositories.Accounts{DB: tx}, TTL: time.Hour}
		learning := repositories.Learning{DB: tx}
		var suffix string
		tx.Raw("SELECT gen_random_uuid()").Scan(&suffix)
		user, err := auth.Register(ctx, services.Registration{Name: "Catalogue test", Email: suffix + "@example.org", Password: "Old-safe-test-password", Timezone: "UTC"})
		if err != nil {
			t.Fatal(err)
		}
		// Simulate an empty published level only within this rolled-back transaction.
		if err := tx.Exec("UPDATE units SET published=false WHERE level_code='B1'").Error; err != nil {
			t.Fatal(err)
		}
		if err := learning.SaveLevel(ctx, user.ID, "B1", models.LevelEdit{Title: "Intermediate", Available: true}); apperr.HTTPStatus(err) != 409 {
			t.Fatalf("empty course enabled: %v", err)
		}
		// The SQL transaction intentionally rolls back, preserving all preexisting catalogue data.
		unitID, err := learning.SaveUnit(ctx, user.ID, "", models.UnitEdit{LevelCode: "B1", Title: "Test unit", Description: "A focused test unit.", Published: true})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := learning.SaveUnit(ctx, user.ID, unitID, models.UnitEdit{LevelCode: "B2", Title: "Invalid move", Published: true}); apperr.HTTPStatus(err) != 409 {
			t.Fatal("unit moved across levels")
		}
		// Expected business errors do not abort the outer transaction.
		resource := models.Resource{Key: "test_image", Title: "Test image", Kind: "image", URL: "/images/apple.svg", Attribution: "EnglishQuest original", Active: true}
		if err := services.ValidateResource(resource); err != nil {
			t.Fatal(err)
		}
		if err := learning.SaveResource(ctx, user.ID, resource); err != nil {
			t.Fatal(err)
		}
		activity := models.Activity{Title: "Choice", Kind: "word_hunter", Instructions: "Choose the matching picture.", Published: true, Questions: []models.Question{{Prompt: "Find: apple", Explanation: "Apple means manzana.", Solution: json.RawMessage(`"a"`), Options: []models.Option{{ID: "a", Label: "Manzana", Image: "test_image"}, {ID: "b", Label: "Plátano", Image: "banana"}}}}}
		id, err := learning.SaveContent(ctx, user.ID, "", repositories.ContentEdit{UnitID: unitID, LessonTitle: "Test lesson", Objective: "Recognize a fruit.", Content: "An apple is a fruit.", Activity: activity})
		if err != nil {
			t.Fatal(err)
		}
		if err := learning.SaveLevel(ctx, user.ID, "B1", models.LevelEdit{Title: "Intermediate", Description: "Introductory test material.", Available: true}); err != nil {
			t.Fatal(err)
		}
		loaded, err := learning.AdminActivity(ctx, id)
		if err != nil || loaded.Questions[0].Options[0].ImageURL != "/images/apple.svg" {
			t.Fatalf("resource not resolved: %+v %v", loaded, err)
		}
		lessons, err := learning.Lessons(ctx, user.ID)
		if err != nil {
			t.Fatal(err)
		}
		found := false
		for _, l := range lessons {
			if l.ActivityID == id {
				found = true
				if l.Locked || l.UnitDescription != "A focused test unit." {
					t.Fatal("invalid unit metadata")
				}
			}
		}
		if !found {
			t.Fatal("published B1 lesson unavailable")
		}
		if _, err := learning.SaveUnit(ctx, user.ID, unitID, models.UnitEdit{LevelCode: "B1", Title: "Hidden", Published: false}); err != nil {
			t.Fatal(err)
		}
		lessons, _ = learning.Lessons(ctx, user.ID)
		for _, l := range lessons {
			if l.ActivityID == id {
				t.Fatal("hidden unit still listed")
			}
		}
		_, raw, err := auth.Login(ctx, user.Email, "Old-safe-test-password")
		if err != nil {
			t.Fatal(err)
		}
		if err := auth.ChangePassword(ctx, user.ID, "wrong", "New-safe-test-password"); apperr.HTTPStatus(err) != 403 {
			t.Fatal("wrong current password accepted")
		}
		if err := auth.ChangePassword(ctx, user.ID, "Old-safe-test-password", "New-safe-test-password"); err != nil {
			t.Fatal(err)
		}
		if _, err := auth.Accounts.Authenticate(ctx, services.HashSession(raw)); err == nil {
			t.Fatal("old session survived password change")
		}
		if _, _, err := auth.Login(ctx, user.Email, "Old-safe-test-password"); apperr.HTTPStatus(err) != 401 {
			t.Fatal("old password still valid")
		}
		if _, _, err := auth.Login(ctx, user.Email, "New-safe-test-password"); err != nil {
			t.Fatal(err)
		}
		return rollback
	})
	if !errors.Is(err, rollback) {
		t.Fatal(err)
	}
}
