package main

import (
	"context"
	"encoding/json"
	"github.com/englishquest/backend/internal/db"
	"github.com/englishquest/backend/internal/models"
	"github.com/englishquest/backend/internal/repositories"
	"github.com/englishquest/backend/internal/services"
	"os"
	"strings"
	"testing"
	"time"
)

func correctSubmission(t *testing.T, activity models.Activity) models.Submission {
	t.Helper()
	out := models.Submission{Answers: map[string]json.RawMessage{}}
	for _, q := range activity.Questions {
		if q.AnswerType(activity.Kind) == "letters" {
			var word string
			if err := json.Unmarshal(q.Solution, &word); err != nil {
				t.Fatal(err)
			}
			ids := []string{}
			used := map[string]bool{}
			for _, letter := range word {
				found := false
				for i := len(q.Options) - 1; i >= 0; i-- {
					o := q.Options[i]
					if !used[o.ID] && strings.EqualFold(o.Label, string(letter)) {
						ids = append(ids, o.ID)
						used[o.ID] = true
						found = true
						break
					}
				}
				if !found {
					t.Fatal("letter missing from bag")
				}
			}
			raw, _ := json.Marshal(ids)
			out.Answers[q.ID] = raw
		} else {
			out.Answers[q.ID] = q.Solution
		}
	}
	return out
}

func TestAllTenGameTypesAndServerDeadlines(t *testing.T) {
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
	var suffix string
	database.Raw("SELECT gen_random_uuid()").Scan(&suffix)
	auth := services.Auth{Accounts: repositories.Accounts{DB: database}, TTL: time.Hour}
	user, err := auth.Register(context.Background(), services.Registration{Name: "All games test", Email: suffix + "@example.org", Password: "Safe-test-password-2026", Timezone: "UTC"})
	if err != nil {
		t.Fatal(err)
	}
	learning := repositories.Learning{DB: database}
	ctx := context.Background()
	lessons, err := learning.Lessons(ctx, user.ID)
	if err != nil {
		t.Fatal(err)
	}
	kinds := map[string]int{}
	var timedLesson models.LessonSummary
	for _, lesson := range lessons {
		started, err := learning.Start(ctx, user.ID, lesson.ActivityID)
		if err != nil {
			t.Fatalf("%s start: %v", lesson.Kind, err)
		}
		public, _ := json.Marshal(started)
		if strings.Contains(string(public), `"solution"`) || strings.Contains(string(public), `"explanation"`) {
			t.Fatal("private key leaked")
		}
		var snapshot string
		database.Raw("SELECT snapshot::text FROM activity_attempts WHERE id=?", started.ID).Scan(&snapshot)
		var activity models.Activity
		if err := json.Unmarshal([]byte(snapshot), &activity); err != nil {
			t.Fatal(err)
		}
		if err := services.ValidateActivity(activity); err != nil {
			t.Fatalf("invalid seeded %s: %v", activity.Kind, err)
		}
		if activity.Kind == "sound_detective" && activity.Questions[0].AudioURL == "" {
			t.Fatal("listening has no real resource URL")
		}
		if activity.TimeLimitSeconds > 0 && started.ExpiresAt == nil {
			t.Fatal("missing server deadline")
		}
		result, err := learning.Submit(ctx, user.ID, started.ID, correctSubmission(t, activity), services.Grade, services.NextStreak)
		if err != nil || result.Score != 100 || !result.Passed {
			t.Fatalf("%s grade %+v: %v", activity.Kind, result, err)
		}
		kinds[activity.Kind]++
		if activity.Kind == "english_rush" {
			timedLesson = lesson
		}
	}
	if len(kinds) != 10 {
		t.Fatalf("only %d game types tested", len(kinds))
	}
	for _, kind := range []string{"word_rain", "english_rush"} {
		for _, l := range lessons {
			if l.Kind == kind {
				timedLesson = l
				break
			}
		}
		started, err := learning.Start(ctx, user.ID, timedLesson.ActivityID)
		if err != nil {
			t.Fatal(err)
		}
		database.Exec("UPDATE activity_attempts SET started_at=now()-interval '10 minutes' WHERE id=?", started.ID)
		result, err := learning.Submit(ctx, user.ID, started.ID, models.Submission{}, services.Grade, services.NextStreak)
		if err != nil || !result.Expired || result.Passed || result.XP != 0 || result.Score != 0 {
			t.Fatalf("late %s rewarded %+v %v", kind, result, err)
		}
		resumed, err := learning.Attempt(ctx, user.ID, started.ID)
		if err != nil || !resumed.Result.Expired {
			t.Fatal("expired result not saved")
		}
	}
}
