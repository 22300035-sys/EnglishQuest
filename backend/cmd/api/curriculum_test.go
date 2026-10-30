package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/englishquest/backend/curriculum"
	"github.com/englishquest/backend/internal/db"
	"github.com/englishquest/backend/internal/models"
	"github.com/englishquest/backend/internal/repositories"
	"github.com/englishquest/backend/internal/services"
)

func TestGuidedCurriculumAllSixRoutes(t *testing.T) {
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
	learning := repositories.Learning{DB: database}
	var units []struct{ Key, Level, Title string }
	if err := json.Unmarshal(curriculum.Data, &units); err != nil {
		t.Fatal(err)
	}
	expected := map[string]int{}
	for _, u := range units {
		expected[u.Level]++
	}
	if len(expected) != 6 || len(units) != 36 {
		t.Fatal("incomplete curriculum manifest")
	}
	levels, err := learning.Levels(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range levels {
		if l["available"] != true {
			t.Fatalf("level unavailable: %+v", l)
		}
	}
	for _, unit := range units {
		var ids []string
		if err := database.Raw(`SELECT a.id FROM activities a JOIN lessons l ON l.id=a.lesson_id JOIN units u ON u.id=l.unit_id WHERE u.level_code=? AND u.title=? AND u.published AND a.published ORDER BY l.position`, unit.Level, unit.Title).Scan(&ids).Error; err != nil {
			t.Fatal(err)
		}
		if len(ids) != 5 {
			t.Fatalf("%s: %d lessons", unit.Key, len(ids))
		}
		for _, id := range ids {
			activity, err := learning.AdminActivity(ctx, id)
			if err != nil {
				t.Fatal(err)
			}
			if err := services.ValidateActivity(activity); err != nil {
				t.Fatalf("%s %s: %v", unit.Key, activity.Title, err)
			}
			if len(activity.Questions) < 1 {
				t.Fatal("empty lesson")
			}
			for _, q := range activity.Questions {
				if q.AudioKey != "" {
					if !services.BundledAudio[q.AudioURL] {
						t.Fatalf("audio not allowlisted: %s", q.AudioURL)
					}
					path := filepath.Join("..", "..", "..", "frontend", "public", strings.TrimPrefix(q.AudioURL, "/"))
					if stat, err := os.Stat(path); err != nil || stat.Size() < 1000 {
						t.Fatalf("missing real audio: %s", path)
					}
				}
			}
			result, err := services.Grade(activity, correctSubmission(t, activity))
			if err != nil || result.Score != 100 {
				t.Fatalf("invalid solution %s: %+v %v", unit.Key, result, err)
			}
		}
	}
	// Traverse every published route through normal snapshot/start/submit logic.
	auth := services.Auth{Accounts: repositories.Accounts{DB: database}, TTL: time.Hour}
	var suffix string
	database.Raw("SELECT gen_random_uuid()").Scan(&suffix)
	user, err := auth.Register(ctx, services.Registration{Name: "Curriculum traversal", Email: suffix + "@example.org", Password: "Curriculum-test-password-2026", Timezone: "UTC"})
	if err != nil {
		t.Fatal(err)
	}
	lessons, err := learning.Lessons(ctx, user.ID)
	if err != nil {
		t.Fatal(err)
	}
	counts := map[string]int{}
	for _, lesson := range lessons {
		attempt, err := learning.Start(ctx, user.ID, lesson.ActivityID)
		if err != nil {
			t.Fatalf("start %s %s: %v", lesson.LevelCode, lesson.Title, err)
		}
		for _, q := range attempt.Activity.Questions {
			if len(q.Solution) > 0 || q.Explanation != "" {
				t.Fatal("solutions exposed before submission")
			}
		}
		var snapshot string
		if err := database.Raw("SELECT snapshot::text FROM activity_attempts WHERE id=?", attempt.ID).Scan(&snapshot).Error; err != nil {
			t.Fatal(err)
		}
		var activity models.Activity
		if err := json.Unmarshal([]byte(snapshot), &activity); err != nil {
			t.Fatal(err)
		}
		result, err := learning.Submit(ctx, user.ID, attempt.ID, correctSubmission(t, activity), services.Grade, services.NextStreak)
		if err != nil || !result.Passed || result.Score != 100 {
			t.Fatalf("submit %s %s: %+v %v", lesson.LevelCode, lesson.Title, result, err)
		}
		counts[lesson.LevelCode]++
	}
	for level := range expected {
		if counts[level] < 30 {
			t.Fatalf("%s route incomplete: %d", level, counts[level])
		}
	}
	after, err := learning.Lessons(ctx, user.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, lesson := range after {
		if !lesson.Completed || lesson.Locked {
			t.Fatalf("completion not durable: %+v", lesson)
		}
	}
	var earned int
	database.Raw("SELECT count(*) FROM user_achievements WHERE user_id=? AND achievement_code='level_path'", user.ID).Scan(&earned)
	if earned != 1 {
		t.Fatal("route completion achievement missing or duplicated")
	}
	// Productive writing is a private portfolio, not fabricated automatic marking.
	unitID := ""
	for _, lesson := range lessons {
		if lesson.HasProject && lesson.LevelCode == "A1" {
			unitID = lesson.UnitID
			break
		}
	}
	if unitID == "" {
		t.Fatal("initial-level project missing")
	}
	project, err := learning.UnitProject(ctx, user.ID, unitID)
	if err != nil || project.Prompt == "" || project.Rubric == "" || project.Body != "" {
		t.Fatalf("missing production project: %+v %v", project, err)
	}
	var beforeXP int
	database.Raw("SELECT COALESCE(sum(amount),0) FROM xp_transactions WHERE user_id=?", user.ID).Scan(&beforeXP)
	if err := learning.SaveUnitProject(ctx, user.ID, unitID, "My original argument preserves the limits of the evidence.", "I revised the certainty of my conclusion."); err != nil {
		t.Fatal(err)
	}
	project, err = learning.UnitProject(ctx, user.ID, unitID)
	if err != nil || project.Body != "My original argument preserves the limits of the evidence." || project.UpdatedAt == nil {
		t.Fatal("project not saved durably")
	}
	other, err := auth.Register(ctx, services.Registration{Name: "Other writer", Email: "other-" + suffix + "@example.org", Password: "Curriculum-test-password-2026", Timezone: "UTC"})
	if err != nil {
		t.Fatal(err)
	}
	private, err := learning.UnitProject(ctx, other.ID, unitID)
	if err != nil || private.Body != "" || private.Reflection != "" {
		t.Fatal("another writer's project was exposed")
	}
	if err := learning.SaveUnitProject(ctx, user.ID, unitID, "", ""); err == nil {
		t.Fatal("empty draft accepted")
	}
	if err := learning.SaveUnitProject(ctx, user.ID, unitID, strings.Repeat("x", 12001), ""); err == nil {
		t.Fatal("oversized draft accepted")
	}
	var afterXP int
	database.Raw("SELECT COALESCE(sum(amount),0) FROM xp_transactions WHERE user_id=?", user.ID).Scan(&afterXP)
	if afterXP != beforeXP {
		t.Fatal("self-review invented XP")
	}
	// Progress, not access-exam shortcuts, completes the course and earns a symbolic diploma.
	states, err := learning.LevelAccess(ctx, user.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, state := range states {
		if !state.Unlocked {
			t.Fatal("full previous-level completion did not unlock next level")
		}
	}
	status, err := learning.DiplomaStatus(ctx, user.ID)
	if err != nil || !status.Eligible || status.Completed != len(lessons) {
		t.Fatalf("completed course not eligible: %+v %v", status, err)
	}
	var wg sync.WaitGroup
	ids := make(chan string, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			record, err := learning.IssueDiploma(ctx, user.ID, "María del Carmen O'Neill")
			if err != nil {
				t.Error(err)
				return
			}
			ids <- record.ID
		}()
	}
	wg.Wait()
	close(ids)
	id := ""
	for next := range ids {
		if id != "" && id != next {
			t.Fatal("duplicate diploma issuance")
		}
		id = next
	}
	if id == "" {
		t.Fatal("diploma not issued")
	}
	status, err = learning.DiplomaStatus(ctx, user.ID)
	if err != nil || status.Diploma == nil || status.Diploma.ActivityCount != len(lessons) || len(status.Diploma.Levels) != 6 {
		t.Fatal("diploma snapshot not durable")
	}
	otherStatus, err := learning.DiplomaStatus(ctx, other.ID)
	if err != nil || otherStatus.Diploma != nil || otherStatus.Eligible {
		t.Fatal("another student's diploma exposed")
	}
	if _, err := learning.IssueDiploma(ctx, user.ID, "Changed Name"); err == nil {
		t.Fatal("issued name silently overwritten")
	}
}
