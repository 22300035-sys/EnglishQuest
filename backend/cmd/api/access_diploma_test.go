package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/englishquest/backend/internal/config"
	"github.com/englishquest/backend/internal/db"
	"github.com/englishquest/backend/internal/models"
	apperr "github.com/englishquest/backend/internal/pkg/errors"
	"github.com/englishquest/backend/internal/repositories"
	"github.com/englishquest/backend/internal/routes"
	"github.com/englishquest/backend/internal/services"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"
)

func TestLevelAccessExamsAndDiplomaGuards(t *testing.T) {
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
	user, err := auth.Register(ctx, services.Registration{Name: "Access test", Email: suffix + "@example.org", Password: "Access-safe-test-password", Timezone: "UTC"})
	if err != nil {
		t.Fatal(err)
	}
	other, err := auth.Register(ctx, services.Registration{Name: "Other access test", Email: "other-" + suffix + "@example.org", Password: "Access-safe-test-password", Timezone: "UTC"})
	if err != nil {
		t.Fatal(err)
	}
	levels, err := learning.LevelAccess(ctx, user.ID)
	if err != nil {
		t.Fatal(err)
	}
	for i, level := range levels {
		if level.Unlocked != (i < 3) {
			t.Fatalf("initial access: %+v", levels)
		}
	}
	if !levels[3].CanTakeExam || levels[4].CanTakeExam || levels[5].CanTakeExam {
		t.Fatal("exam prerequisites not progressive")
	}
	lessons, err := learning.Lessons(ctx, user.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, lesson := range lessons {
		if lesson.LevelCode == "B2" {
			if !lesson.LevelLocked || lesson.Content != "" {
				t.Fatal("locked content leaked")
			}
			if _, err := learning.Start(ctx, user.ID, lesson.ActivityID); apperr.HTTPStatus(err) != 403 {
				t.Fatal("direct activity bypass")
			}
			if _, err := learning.UnitProject(ctx, user.ID, lesson.UnitID); apperr.HTTPStatus(err) != 403 {
				t.Fatal("direct project read bypass")
			}
			if err := learning.SaveUnitProject(ctx, user.ID, lesson.UnitID, "Bypass attempt", ""); apperr.HTTPStatus(err) != 403 {
				t.Fatal("direct project write bypass")
			}
			break
		}
	}
	if _, err := learning.StartAccessExam(ctx, user.ID, "C2", services.BuildAccessExam); apperr.HTTPStatus(err) != 403 {
		t.Fatal("skipped exam prerequisite")
	}
	if _, err := learning.IssueDiploma(ctx, user.ID, "Student Name"); apperr.HTTPStatus(err) != 409 {
		t.Fatal("premature diploma")
	}
	// Concurrent starts resume one immutable exam rather than creating extra time.
	var wg sync.WaitGroup
	ids := make(chan string, 2)
	errs := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			exam, err := learning.StartAccessExam(ctx, user.ID, "B2", services.BuildAccessExam)
			ids <- exam.ID
			errs <- err
		}()
	}
	wg.Wait()
	close(ids)
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	examID := ""
	for id := range ids {
		if examID != "" && id != examID {
			t.Fatal("concurrent starts created multiple attempts")
		}
		examID = id
	}
	exam, err := learning.AccessExam(ctx, user.ID, examID)
	if err != nil {
		t.Fatal(err)
	}
	for _, q := range exam.Activity.Questions {
		if len(q.Solution) > 0 || q.Explanation != "" {
			t.Fatal("exam key exposed")
		}
		for _, option := range q.Options {
			if len(option.ID) != 24 {
				t.Fatal("answer IDs encode bank keys")
			}
		}
	}
	resumed, err := learning.StartAccessExam(ctx, user.ID, "B2", services.BuildAccessExam)
	if err != nil || !resumed.ExpiresAt.Equal(*exam.ExpiresAt) {
		t.Fatal("resume reset deadline")
	}
	if _, err := learning.AccessExam(ctx, other.ID, examID); apperr.HTTPStatus(err) != 404 {
		t.Fatal("other user read exam")
	}
	if _, err := learning.SubmitAccessExam(ctx, other.ID, examID, models.Submission{}, services.Grade); apperr.HTTPStatus(err) != 404 {
		t.Fatal("other user submitted exam")
	}
	failed, err := learning.SubmitAccessExam(ctx, user.ID, examID, models.Submission{}, services.Grade)
	if err != nil || failed.Result.Passed || failed.RetryAt == nil {
		t.Fatal("empty exam did not fail")
	}
	if _, err := learning.StartAccessExam(ctx, user.ID, "B2", services.BuildAccessExam); apperr.HTTPStatus(err) != 429 {
		t.Fatalf("cooldown missing: %v", err)
	}
	database.Exec("UPDATE access_exam_attempts SET completed_at=now()-interval '11 minutes' WHERE user_id=?", user.ID)
	expired, err := learning.StartAccessExam(ctx, user.ID, "B2", services.BuildAccessExam)
	if err != nil {
		t.Fatal(err)
	}
	database.Exec("UPDATE access_exam_attempts SET started_at=now()-interval '21 minutes',expires_at=now()-interval '1 minute' WHERE id=?", expired.ID)
	snapshot := func(id string) models.Activity {
		var raw string
		database.Raw("SELECT snapshot::text FROM access_exam_attempts WHERE id=?", id).Scan(&raw)
		var a models.Activity
		if err := json.Unmarshal([]byte(raw), &a); err != nil {
			t.Fatal(err)
		}
		return a
	}
	expired, err = learning.SubmitAccessExam(ctx, user.ID, expired.ID, correctSubmission(t, snapshot(expired.ID)), services.Grade)
	if err != nil || !expired.Result.Expired || expired.Result.Passed {
		t.Fatal("expired exam unlocked content")
	}
	database.Exec("UPDATE access_exam_attempts SET completed_at=now()-interval '11 minutes' WHERE user_id=?", user.ID)
	for _, level := range []string{"B2", "C1", "C2"} {
		opened, err := learning.StartAccessExam(ctx, user.ID, level, services.BuildAccessExam)
		if err != nil {
			t.Fatal(err)
		}
		answers := correctSubmission(t, snapshot(opened.ID))
		results := make(chan error, 2)
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				result, err := learning.SubmitAccessExam(ctx, user.ID, opened.ID, answers, services.Grade)
				if err == nil && (!result.Result.Passed || result.Result.XP != 0) {
					t.Errorf("wrong exam result: %+v", result)
				}
				results <- err
			}()
		}
		wg.Wait()
		close(results)
		for err := range results {
			if err != nil {
				t.Fatal(err)
			}
		}
	}
	levels, _ = learning.LevelAccess(ctx, user.ID)
	for _, level := range levels {
		if !level.Unlocked {
			t.Fatal("passed exam did not persist access")
		}
	}
	var grants, progress, xp int
	database.Raw("SELECT count(*) FROM user_level_unlocks WHERE user_id=? AND source='exam'", user.ID).Scan(&grants)
	database.Raw("SELECT count(*) FROM user_progress WHERE user_id=?", user.ID).Scan(&progress)
	database.Raw("SELECT count(*) FROM xp_transactions WHERE user_id=?", user.ID).Scan(&xp)
	if grants != 3 || progress != 0 || xp != 0 {
		t.Fatal("exam duplicated rewards or manufactured course progress")
	}
	if _, err := learning.IssueDiploma(ctx, user.ID, "Student Name"); apperr.HTTPStatus(err) != 409 {
		t.Fatal("access exams alone issued diploma")
	}
	// HTTP authorization, strict payload and CSRF boundaries remain enforced.
	router := routes.NewRouter(&config.Config{Env: "test", AllowedOrigins: []string{"http://localhost:5173"}, JWTRefreshTTL: time.Hour}, database)
	_, raw, err := auth.Login(ctx, user.Email, "Access-safe-test-password")
	if err != nil {
		t.Fatal(err)
	}
	request := func(method, path string, payload interface{}, authenticated, origin bool) int {
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(method, "/api/v1"+path, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		if authenticated {
			req.AddCookie(&http.Cookie{Name: "eq_session", Value: raw})
		}
		if origin {
			req.Header.Set("Origin", "http://localhost:5173")
		}
		out := httptest.NewRecorder()
		router.ServeHTTP(out, req)
		return out.Code
	}
	for _, path := range []string{"/learning/levels", "/diploma", "/access-exams/" + examID} {
		if request("GET", path, nil, false, true) != 401 {
			t.Fatal("unauthenticated access")
		}
	}
	if request("POST", "/diploma", map[string]interface{}{"name": "Student Name", "eligible": true}, true, true) != 400 {
		t.Fatal("forged eligibility accepted")
	}
	if request("POST", "/diploma", map[string]interface{}{"name": "Student Name"}, true, false) != 403 {
		t.Fatal("Origin check bypassed")
	}
	if request("POST", "/diploma", map[string]interface{}{"name": "<script>bad</script>"}, true, true) != 400 {
		t.Fatal("unsafe diploma name accepted")
	}
}
