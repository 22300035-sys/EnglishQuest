package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/englishquest/backend/internal/config"
	"github.com/englishquest/backend/internal/db"
	"github.com/englishquest/backend/internal/models"
	"github.com/englishquest/backend/internal/routes"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestPrototypeFlow(t *testing.T) {
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
	gin.SetMode(gin.TestMode)
	router := routes.NewRouter(&config.Config{Env: "test", AllowedOrigins: []string{"http://localhost:5173"}, JWTRefreshTTL: time.Hour}, database)
	request := func(method, path string, payload interface{}, cookie *http.Cookie, origin bool) *httptest.ResponseRecorder {
		raw, _ := json.Marshal(payload)
		req := httptest.NewRequest(method, "/api/v1"+path, bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		if origin {
			req.Header.Set("Origin", "http://localhost:5173")
		}
		if cookie != nil {
			req.AddCookie(cookie)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		return w
	}
	expect := func(w *httptest.ResponseRecorder, status int) {
		t.Helper()
		if w.Code != status {
			t.Fatalf("status %d wanted %d: %s", w.Code, status, w.Body.String())
		}
	}
	var suffix string
	database.Raw("SELECT gen_random_uuid()").Scan(&suffix)
	email := "prototype-" + suffix + "@example.org"
	password := "Test-only-password-42"
	registration := map[string]interface{}{"email": email, "password": password, "name": "Explorer", "timezone": "America/Mexico_City"}
	expect(request("POST", "/auth/register", registration, nil, false), 403)
	injected := map[string]interface{}{"email": email, "password": password, "name": "Explorer", "role": "ADMIN"}
	expect(request("POST", "/auth/register", injected, nil, true), 400)
	expect(request("POST", "/auth/register", registration, nil, true), 201)
	expect(request("POST", "/auth/register", registration, nil, true), 409)
	var stored models.User
	database.Where("email=?", email).First(&stored)
	if stored.PasswordHash == password || !strings.HasPrefix(stored.PasswordHash, "$2") {
		t.Fatal("password not hashed")
	}
	expect(request("POST", "/auth/login", map[string]string{"email": email, "password": "wrong"}, nil, true), 401)
	login := request("POST", "/auth/login", map[string]string{"email": email, "password": password}, nil, true)
	expect(login, 200)
	cookies := login.Result().Cookies()
	if len(cookies) != 1 || !cookies[0].HttpOnly || len(cookies[0].Value) != 64 {
		t.Fatal("unsafe session cookie")
	}
	cookie := cookies[0]
	if strings.Contains(login.Body.String(), "password_hash") || strings.Contains(login.Body.String(), cookie.Value) {
		t.Fatal("credential leaked")
	}
	expect(request("GET", "/auth/me", nil, cookie, true), 200)
	expect(request("GET", "/dashboard", nil, nil, true), 401)
	expect(request("GET", "/admin/overview", nil, cookie, true), 403)
	expect(request("POST", "/activities/bad-id/start", nil, cookie, true), 400)
	var catalog struct{ Data []models.LessonSummary }
	w := request("GET", "/lessons", nil, cookie, true)
	expect(w, 200)
	json.Unmarshal(w.Body.Bytes(), &catalog)
	legacy := []models.LessonSummary{}
	for _, lesson := range catalog.Data {
		if lesson.Kind == "word_hunter" || lesson.Kind == "sentence_builder" || lesson.Kind == "memory_match" {
			legacy = append(legacy, lesson)
		}
	}
	catalog.Data = legacy
	if len(catalog.Data) != 9 {
		t.Fatalf("expected 9 starter lessons; got %d", len(catalog.Data))
	}
	expect(request("POST", "/activities/"+catalog.Data[1].ActivityID+"/start", nil, cookie, true), 403)
	var firstID string
	for index, l := range catalog.Data {
		started := request("POST", "/activities/"+l.ActivityID+"/start", nil, cookie, true)
		expect(started, 201)
		if strings.Contains(started.Body.String(), `"solution"`) || strings.Contains(started.Body.String(), `"explanation"`) {
			t.Fatal("answer key leaked before completion")
		}
		var attempt struct {
			Data struct {
				ID       string
				Activity models.Activity
			}
		}
		json.Unmarshal(started.Body.Bytes(), &attempt)
		resumed := request("POST", "/activities/"+l.ActivityID+"/start", nil, cookie, true)
		expect(resumed, 201)
		if !strings.Contains(resumed.Body.String(), attempt.Data.ID) {
			t.Fatal("start is not resumable")
		}
		var snapshot string
		database.Raw("SELECT snapshot::text FROM activity_attempts WHERE id=?", attempt.Data.ID).Scan(&snapshot)
		var private models.Activity
		json.Unmarshal([]byte(snapshot), &private)
		answers := models.Submission{Answers: map[string]json.RawMessage{}}
		for _, q := range private.Questions {
			answers.Answers[q.ID] = q.Solution
			for _, o := range q.Options {
				if len(o.ID) != 24 {
					t.Fatal("option identifier reveals answer position")
				}
			}
		}
		if index == 0 {
			firstID = attempt.Data.ID
			expect(request("POST", "/attempts/"+firstID+"/submit", models.Submission{}, cookie, true), 400)
			expect(request("POST", "/attempts/"+firstID+"/submit", map[string]interface{}{"answers": answers.Answers, "score": 100, "xp": 999}, cookie, true), 400)
		}
		var result struct{ Data models.Result }
		// Concurrent identical submissions must only commit one completion/reward.
		var group sync.WaitGroup
		responses := make(chan *httptest.ResponseRecorder, 2)
		for i := 0; i < 2; i++ {
			group.Add(1)
			go func() {
				defer group.Done()
				responses <- request("POST", "/attempts/"+attempt.Data.ID+"/submit", answers, cookie, true)
			}()
		}
		group.Wait()
		close(responses)
		for response := range responses {
			expect(response, 200)
			json.Unmarshal(response.Body.Bytes(), &result)
			if result.Data.Score != 100 || !result.Data.Passed {
				t.Fatalf("incorrect result: %+v", result.Data)
			}
		}
	}
	var total int
	database.Raw("SELECT sum(amount) FROM xp_transactions WHERE user_id=?", stored.ID).Scan(&total)
	if total != 410 {
		t.Fatalf("XP was duplicated or missing: %d", total)
	}
	dashboard := request("GET", "/dashboard", nil, cookie, true)
	expect(dashboard, 200)
	var dashboardBody struct {
		Data struct {
			XP     int
			Streak struct {
				CurrentDays int `json:"current_days"`
			}
			Lessons  []models.LessonSummary
			Missions []struct {
				Code     string
				Rewarded bool
			}
		}
	}
	json.Unmarshal(dashboard.Body.Bytes(), &dashboardBody)
	missionRewarded := false
	for _, mission := range dashboardBody.Data.Missions {
		if mission.Code == "three_activities" {
			missionRewarded = mission.Rewarded
		}
	}
	if dashboardBody.Data.XP != 410 || dashboardBody.Data.Streak.CurrentDays != 1 || !missionRewarded {
		t.Fatalf("dashboard did not persist rewards: %s", dashboard.Body.String())
	}
	// A fresh repeat attempt has no new XP, even when answered perfectly.
	repeated := request("POST", "/activities/"+catalog.Data[0].ActivityID+"/start", nil, cookie, true)
	expect(repeated, 201)
	var again struct{ Data struct{ ID string } }
	json.Unmarshal(repeated.Body.Bytes(), &again)
	var snapshot string
	database.Raw("SELECT snapshot::text FROM activity_attempts WHERE id=?", again.Data.ID).Scan(&snapshot)
	var private models.Activity
	json.Unmarshal([]byte(snapshot), &private)
	answers := models.Submission{Answers: map[string]json.RawMessage{}}
	for _, q := range private.Questions {
		answers.Answers[q.ID] = q.Solution
	}
	repeat := request("POST", "/attempts/"+again.Data.ID+"/submit", answers, cookie, true)
	expect(repeat, 200)
	var result struct{ Data models.Result }
	json.Unmarshal(repeat.Body.Bytes(), &result)
	if result.Data.XP != 0 || result.Data.NewCompletion {
		t.Fatal("repeat farming rewarded")
	}
	// Another account cannot read someone else's attempt.
	otherEmail := "other-" + email
	registration["email"] = otherEmail
	expect(request("POST", "/auth/register", registration, nil, true), 201)
	other := request("POST", "/auth/login", map[string]string{"email": otherEmail, "password": password}, nil, true)
	expect(other, 200)
	expect(request("GET", "/attempts/"+firstID, nil, other.Result().Cookies()[0], true), 404)
	// Promote only within the disposable test DB; authorization reads the new DB role.
	database.Exec("UPDATE users SET role='ADMIN' WHERE id=?", stored.ID)
	expect(request("GET", "/admin/overview", nil, cookie, true), 200)
	edit := request("GET", "/admin/activities/"+catalog.Data[0].ActivityID, nil, cookie, true)
	expect(edit, 200)
	var edited struct{ Data models.Activity }
	json.Unmarshal(edit.Body.Bytes(), &edited)
	var originalSnapshot string
	database.Raw("SELECT snapshot::text FROM activity_attempts WHERE id=?", firstID).Scan(&originalSnapshot)
	edited.Data.Title = "Edited safely"
	payload := map[string]interface{}{"activity": edited.Data, "lesson_title": "Edited lesson", "objective": "Recognize fruit words.", "content": "An apple is a fruit. Practice these English words."}
	expect(request("PUT", "/admin/activities/"+catalog.Data[0].ActivityID, payload, cookie, true), 200)
	// Create hidden content without altering the published sample used by other tests.
	edited.Data.Published = false
	payload["activity"] = edited.Data
	payload["unit_id"] = catalog.Data[0].UnitID
	expect(request("POST", "/admin/lessons", payload, cookie, true), 200)
	var otherUser models.User
	database.Where("email=?", otherEmail).First(&otherUser)
	expect(request("PATCH", "/admin/users/"+otherUser.ID, map[string]bool{"active": false}, cookie, true), 200)
	expect(request("GET", "/auth/me", nil, other.Result().Cookies()[0], true), 401)
	expect(request("POST", "/auth/login", map[string]string{"email": otherEmail, "password": password}, nil, true), 401)
	expect(request("PATCH", "/admin/users/"+otherUser.ID, map[string]bool{"active": true}, cookie, true), 200)
	// Reactivating does not resurrect a revoked session.
	expect(request("GET", "/auth/me", nil, other.Result().Cookies()[0], true), 401)
	// Completed snapshot remains intact despite editing content.
	expiring := request("POST", "/auth/login", map[string]string{"email": otherEmail, "password": password}, nil, true)
	expect(expiring, 200)
	database.Exec("UPDATE refresh_tokens SET expires_at=now()-interval '1 second' WHERE user_id=?", otherUser.ID)
	expect(request("GET", "/auth/me", nil, expiring.Result().Cookies()[0], true), 401)
	saved := request("GET", "/attempts/"+firstID, nil, cookie, true)
	expect(saved, 200)
	var currentSnapshot string
	database.Raw("SELECT snapshot::text FROM activity_attempts WHERE id=?", firstID).Scan(&currentSnapshot)
	if originalSnapshot != currentSnapshot {
		t.Fatal("attempt snapshot changed")
	}
	expect(request("POST", "/auth/logout", nil, cookie, true), 200)
	expect(request("GET", "/auth/me", nil, cookie, true), 401)
	t.Log(fmt.Sprintf("Full flow passed for %d seeded lessons and all three games", len(catalog.Data)))
}
