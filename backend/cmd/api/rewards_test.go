package main

import (
	"context"
	"encoding/json"
	"github.com/englishquest/backend/internal/db"
	"github.com/englishquest/backend/internal/models"
	apperr "github.com/englishquest/backend/internal/pkg/errors"
	"github.com/englishquest/backend/internal/repositories"
	"github.com/englishquest/backend/internal/services"
	"os"
	"sync"
	"testing"
	"time"
)

func TestBoosterOwnershipConcurrencyExpiryAndNoFarming(t *testing.T) {
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
	ctx := context.Background()
	auth := services.Auth{Accounts: repositories.Accounts{DB: database}, TTL: time.Hour}
	learning := repositories.Learning{DB: database}
	var suffix string
	database.Raw("SELECT gen_random_uuid()").Scan(&suffix)
	user, err := auth.Register(ctx, services.Registration{Name: "Rewards test", Email: suffix + "@example.org", Password: "Safe-rewards-password", Timezone: "UTC"})
	if err != nil {
		t.Fatal(err)
	}
	lessons, err := learning.Lessons(ctx, user.ID)
	if err != nil {
		t.Fatal(err)
	}
	complete := func(lesson models.LessonSummary) models.Result {
		t.Helper()
		attempt, err := learning.Start(ctx, user.ID, lesson.ActivityID)
		if err != nil {
			t.Fatal(err)
		}
		var snapshot string
		database.Raw("SELECT snapshot::text FROM activity_attempts WHERE id=?", attempt.ID).Scan(&snapshot)
		var activity models.Activity
		json.Unmarshal([]byte(snapshot), &activity)
		result, err := learning.Submit(ctx, user.ID, attempt.ID, correctSubmission(t, activity), services.Grade, services.NextStreak)
		if err != nil {
			t.Fatal(err)
		}
		return result
	}
	for i := 0; i < 3; i++ {
		complete(lessons[i])
	}
	var inventory []struct{ ID string }
	database.Raw("SELECT id FROM user_boosters WHERE user_id=?", user.ID).Scan(&inventory)
	if len(inventory) != 1 {
		t.Fatalf("mission granted %d boosters", len(inventory))
	}
	var group sync.WaitGroup
	results := make(chan error, 2)
	for i := 0; i < 2; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			_, err := learning.ActivateBooster(ctx, user.ID, inventory[0].ID)
			results <- err
		}()
	}
	group.Wait()
	close(results)
	success, conflict := 0, 0
	for err := range results {
		if err == nil {
			success++
		} else if apperr.HTTPStatus(err) == 409 {
			conflict++
		} else {
			t.Fatal(err)
		}
	}
	if success != 1 || conflict != 1 {
		t.Fatal("concurrent activation was not serialized")
	}
	// Simulate a counter near the achievement boundary, not fabricated public progress.
	database.Exec("UPDATE user_skill_stats SET words_practiced=99 WHERE user_id=?", user.ID)
	result := complete(lessons[3])
	if result.Multiplier != 2 || result.XP != 50 {
		t.Fatalf("base/perfect multiplication wrong: %+v", result)
	}
	var earned int
	database.Raw("SELECT count(*) FROM user_achievements WHERE user_id=? AND achievement_code='hundred_words'", user.ID).Scan(&earned)
	if earned != 1 {
		t.Fatal("100-word threshold not awarded")
	}
	repeated := complete(lessons[0])
	if repeated.XP != 0 || repeated.NewCompletion {
		t.Fatal("booster rewarded a repeated activity")
	}
	database.Exec("UPDATE user_boosters SET activated_at=now()-interval '16 minutes',expires_at=now()-interval '1 minute' WHERE id=?", inventory[0].ID)
	result = complete(lessons[4])
	if result.Multiplier != 1 {
		t.Fatal("expired booster still applied")
	}
	var activityXP int
	database.Raw("SELECT amount FROM xp_transactions WHERE user_id=? AND source=?", user.ID, "activity:"+lessons[4].ActivityID).Scan(&activityXP)
	if activityXP != 30 {
		t.Fatal("expired factor applied to XP ledger")
	}
	if _, err := learning.ActivateBooster(ctx, user.ID, inventory[0].ID); apperr.HTTPStatus(err) != 409 {
		t.Fatal("used booster reactivated")
	}
	other, err := auth.Register(ctx, services.Registration{Name: "Other rewards test", Email: "other-" + suffix + "@example.org", Password: "Safe-rewards-password", Timezone: "UTC"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := learning.ActivateBooster(ctx, other.ID, inventory[0].ID); apperr.HTTPStatus(err) != 404 {
		t.Fatal("another user activated private inventory")
	}
	var boosters int
	database.Raw("SELECT count(*) FROM user_boosters WHERE user_id=?", user.ID).Scan(&boosters)
	if boosters != 1 {
		t.Fatal("repeat mission inventory farming")
	}
}
