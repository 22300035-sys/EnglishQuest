package repositories

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/englishquest/backend/internal/models"
	apperr "github.com/englishquest/backend/internal/pkg/errors"
	"gorm.io/gorm"
	"math/big"
	"time"
)

type Learning struct{ DB *gorm.DB }

func (r Learning) Levels(ctx context.Context) ([]map[string]interface{}, error) {
	rows := []map[string]interface{}{}
	err := r.DB.WithContext(ctx).Raw(`SELECT c.*,
 (SELECT count(DISTINCT u.id) FROM units u JOIN lessons l ON l.unit_id=u.id JOIN activities a ON a.lesson_id=l.id WHERE u.level_code=c.code AND u.published AND a.published) unit_count,
 (SELECT count(*) FROM units u JOIN lessons l ON l.unit_id=u.id JOIN activities a ON a.lesson_id=l.id WHERE u.level_code=c.code AND u.published AND a.published) lesson_count
 FROM cefr_levels c ORDER BY c.position`).Scan(&rows).Error
	for _, row := range rows {
		code, _ := row["code"].(string)
		row["unlocked"] = code == "A1" || code == "A2" || code == "B1"
	}
	return rows, err
}
func lessonList(db *gorm.DB, userID string) ([]models.LessonSummary, error) {
	lessons := []models.LessonSummary{}
	err := db.Raw(`SELECT l.id,l.unit_id,u.title unit_title,u.description unit_description,EXISTS(SELECT 1 FROM unit_projects pr WHERE pr.unit_id=u.id) has_project,u.level_code,l.title,l.objective,l.content,
 a.id activity_id,a.kind,a.assessment,(p.user_id IS NOT NULL) completed,COALESCE(p.best_score,0) best_score
  FROM lessons l JOIN units u ON u.id=l.unit_id AND u.published JOIN cefr_levels c ON c.code=u.level_code AND c.available JOIN activities a ON a.lesson_id=l.id AND a.published
 LEFT JOIN user_progress p ON p.activity_id=a.id AND p.user_id=?
 ORDER BY u.level_code,u.position,l.position,a.id`, userID).Scan(&lessons).Error
	previousLevel := ""
	unlocked := true
	for i := range lessons {
		if lessons[i].LevelCode != previousLevel {
			unlocked = true
			previousLevel = lessons[i].LevelCode
		}
		lessons[i].Locked = !unlocked
		unlocked = unlocked && lessons[i].Completed
	}
	if err != nil {
		return nil, err
	}
	if err := applyLevelLocks(db, userID, lessons); err != nil {
		return nil, err
	}
	return lessons, nil
}
func (r Learning) Lessons(ctx context.Context, userID string) ([]models.LessonSummary, error) {
	return lessonList(r.DB.WithContext(ctx), userID)
}

func loadActivity(db *gorm.DB, id string, lock ...bool) (models.Activity, error) {
	a := models.Activity{}
	query := "SELECT * FROM activities WHERE id=?"
	if len(lock) > 0 && lock[0] {
		query += " FOR SHARE"
	}
	res := db.Raw(query, id).Scan(&a)
	if res.Error != nil {
		return a, res.Error
	}
	if a.ID == "" {
		return a, apperr.NewAppError(404, "Actividad no encontrada.")
	}
	if len(lock) > 0 && lock[0] {
		var unavailable int64
		if err := db.Raw(`SELECT count(*) FROM questions q WHERE q.activity_id=? AND
 ((q.audio_key IS NOT NULL AND NOT EXISTS(SELECT 1 FROM educational_resources r WHERE r.key=q.audio_key AND r.active AND r.kind='audio')) OR
 EXISTS(SELECT 1 FROM answer_options o WHERE o.question_id=q.id AND o.image_resource_key IS NOT NULL AND NOT EXISTS(SELECT 1 FROM educational_resources r WHERE r.key=o.image_resource_key AND r.active AND r.kind='image')))`, id).Scan(&unavailable).Error; err != nil {
			return a, err
		}
		if unavailable > 0 {
			return a, apperr.NewAppError(409, "La actividad contiene recursos inactivos. Pide al administrador que los revise.")
		}
	}
	var questions []struct{ ID, Prompt, Explanation, Key, Response, Passage, AudioKey, AudioURL string }
	if err := db.Raw(`SELECT q.id,q.prompt,q.explanation,q.solution::text key,q.response,q.passage,COALESCE(q.audio_key,'') audio_key,COALESCE(r.url,'') audio_url FROM questions q LEFT JOIN educational_resources r ON r.key=q.audio_key WHERE q.activity_id=? ORDER BY q.position`, id).Scan(&questions).Error; err != nil {
		return a, err
	}
	a.Questions = []models.Question{}
	for _, q := range questions {
		options := []models.Option{}
		if err := db.Raw(`SELECT o.id,o.label,o.image,COALESCE(r.url,'') image_url FROM answer_options o LEFT JOIN educational_resources r ON r.key=o.image_resource_key WHERE o.question_id=? ORDER BY o.id`, q.ID).Scan(&options).Error; err != nil {
			return a, err
		}
		a.Questions = append(a.Questions, models.Question{ID: q.ID, Prompt: q.Prompt, Explanation: q.Explanation, Solution: json.RawMessage(q.Key), Options: options, Response: q.Response, Passage: q.Passage, AudioKey: q.AudioKey, AudioURL: q.AudioURL})
	}
	return a, nil
}
func shuffle(options []models.Option) {
	for i := len(options) - 1; i > 0; i-- {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err == nil {
			j := int(n.Int64())
			options[i], options[j] = options[j], options[i]
		}
	}
}

// Opaque per-attempt option identifiers prevent ordering/pair IDs from encoding
// the answer key (e.g. token-1 or English-1 / Spanish-1).
func anonymize(a *models.Activity) error {
	for i := range a.Questions {
		q := &a.Questions[i]
		mapping := map[string]string{}
		for j := range q.Options {
			bytes := make([]byte, 12)
			if _, err := rand.Read(bytes); err != nil {
				return err
			}
			opaque := hex.EncodeToString(bytes)
			mapping[q.Options[j].ID] = opaque
			q.Options[j].ID = opaque
		}
		if q.AnswerType(a.Kind) == "text" || q.AnswerType(a.Kind) == "letters" {
			shuffle(q.Options)
			continue
		}
		var key interface{}
		if err := json.Unmarshal(q.Solution, &key); err != nil {
			return err
		}
		var rewrite func(interface{}) interface{}
		rewrite = func(value interface{}) interface{} {
			switch v := value.(type) {
			case string:
				return mapping[v]
			case []interface{}:
				for j := range v {
					v[j] = rewrite(v[j])
				}
				return v
			default:
				return v
			}
		}
		raw, err := json.Marshal(rewrite(key))
		if err != nil {
			return err
		}
		q.Solution = raw
		shuffle(q.Options)
	}
	return nil
}

type AttemptView struct {
	ID         string          `json:"id"`
	Activity   models.Activity `json:"activity"`
	Result     *models.Result  `json:"result,omitempty"`
	StartedAt  time.Time       `json:"started_at"`
	ServerTime time.Time       `json:"server_time"`
	ExpiresAt  *time.Time      `json:"expires_at,omitempty"`
}

func publicAttempt(id string, a models.Activity, started ...time.Time) AttemptView {
	for i := range a.Questions {
		a.Questions[i].Solution = nil
		a.Questions[i].Explanation = ""
	}
	view := AttemptView{ID: id, Activity: a, ServerTime: time.Now().UTC()}
	if len(started) > 0 {
		view.StartedAt = started[0]
		if a.TimeLimitSeconds > 0 {
			end := started[0].Add(time.Duration(a.TimeLimitSeconds) * time.Second)
			view.ExpiresAt = &end
		}
	}
	return view
}
func (r Learning) Start(ctx context.Context, userID, activityID string) (AttemptView, error) {
	var out AttemptView
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT id FROM users WHERE id=? FOR NO KEY UPDATE", userID).Error; err != nil {
			return err
		}
		lessons, err := lessonList(tx, userID)
		if err != nil {
			return err
		}
		found := false
		for _, l := range lessons {
			if l.ActivityID == activityID {
				found = true
				if l.Locked {
					return apperr.NewAppError(403, "Completa la lección anterior para desbloquear esta actividad.")
				}
			}
		}
		if !found {
			return apperr.NewAppError(404, "Actividad no disponible.")
		}
		var old struct {
			ID, Snapshot string
			StartedAt    time.Time
		}
		if err := tx.Raw("SELECT id,snapshot::text,started_at FROM activity_attempts WHERE user_id=? AND activity_id=? AND completed_at IS NULL", userID, activityID).Scan(&old).Error; err != nil {
			return err
		}
		if old.ID != "" {
			var a models.Activity
			if err := json.Unmarshal([]byte(old.Snapshot), &a); err != nil {
				return err
			}
			out = publicAttempt(old.ID, a, old.StartedAt)
			return nil
		}
		a, err := loadActivity(tx, activityID, true)
		if err != nil {
			return err
		}
		if len(a.Questions) == 0 {
			return apperr.NewAppError(409, "La actividad no tiene preguntas publicadas.")
		}
		if err := anonymize(&a); err != nil {
			return err
		}
		snapshot, err := json.Marshal(a)
		if err != nil {
			return err
		}
		var created struct {
			ID        string
			StartedAt time.Time
		}
		if err := tx.Raw("INSERT INTO activity_attempts(user_id,activity_id,snapshot) VALUES (?,?,?::jsonb) RETURNING id,started_at", userID, activityID, string(snapshot)).Scan(&created).Error; err != nil {
			return err
		}
		out = publicAttempt(created.ID, a, created.StartedAt)
		return nil
	})
	return out, err
}
func (r Learning) Attempt(ctx context.Context, userID, id string) (AttemptView, error) {
	var row struct {
		ID, Snapshot, Result string
		StartedAt            time.Time
	}
	if err := r.DB.WithContext(ctx).Raw("SELECT id,snapshot::text,COALESCE(result::text,'') result,started_at FROM activity_attempts WHERE id=? AND user_id=?", id, userID).Scan(&row).Error; err != nil {
		return AttemptView{}, err
	}
	if row.ID == "" {
		return AttemptView{}, apperr.NewAppError(404, "Intento no encontrado.")
	}
	var a models.Activity
	if err := json.Unmarshal([]byte(row.Snapshot), &a); err != nil {
		return AttemptView{}, err
	}
	out := publicAttempt(row.ID, a, row.StartedAt)
	if row.Result != "" {
		out.Result = &models.Result{}
		if err := json.Unmarshal([]byte(row.Result), out.Result); err != nil {
			return out, err
		}
	}
	return out, nil
}

func award(tx *gorm.DB, userID, source string, amount int) (int, error) {
	res := tx.Exec("INSERT INTO xp_transactions(user_id,source,amount) VALUES (?,?,?) ON CONFLICT DO NOTHING", userID, source, amount)
	if res.Error != nil {
		return 0, res.Error
	}
	if res.RowsAffected == 0 {
		return 0, nil
	}
	return amount, nil
}

// Submit serializes a user's rewards and attempt completion in ONE transaction.
func (r Learning) Submit(ctx context.Context, userID, id string, answers models.Submission,
	grade func(models.Activity, models.Submission) (models.Result, error),
	streak func(string, string, int, int) (int, int)) (models.Result, error) {
	var result models.Result
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.Raw("SELECT * FROM users WHERE id=? FOR NO KEY UPDATE", userID).Scan(&user).Error; err != nil {
			return err
		}
		var row struct {
			ID, ActivityID, Snapshot, Result string
			StartedAt                        time.Time
		}
		if err := tx.Raw("SELECT id,activity_id,snapshot::text,COALESCE(result::text,'') result,started_at FROM activity_attempts WHERE id=? AND user_id=? FOR UPDATE", id, userID).Scan(&row).Error; err != nil {
			return err
		}
		if row.ID == "" {
			return apperr.NewAppError(404, "Intento no encontrado.")
		}
		if row.Result != "" {
			return json.Unmarshal([]byte(row.Result), &result)
		}
		var a models.Activity
		if err := json.Unmarshal([]byte(row.Snapshot), &a); err != nil {
			return err
		}
		var err error
		result, err = grade(a, answers)
		if err != nil {
			return err
		}
		result.AttemptID = id
		result.Multiplier = 1
		result.ElapsedSeconds = int(time.Since(row.StartedAt).Seconds())
		if result.ElapsedSeconds < 0 {
			result.ElapsedSeconds = 0
		}
		if a.TimeLimitSeconds > 0 && time.Now().After(row.StartedAt.Add(time.Duration(a.TimeLimitSeconds)*time.Second)) {
			result.Expired = true
			result.Score = 0
			result.Passed = false
			result.Stars = 0
		}
		loc, err := time.LoadLocation(user.Timezone)
		if err != nil {
			return err
		}
		today := time.Now().In(loc).Format("2006-01-02")
		if result.Passed {
			inserted := tx.Exec("INSERT INTO user_progress(user_id,activity_id,best_score) VALUES (?,?,?) ON CONFLICT DO NOTHING", userID, a.ID, result.Score)
			if inserted.Error != nil {
				return inserted.Error
			}
			result.NewCompletion = inserted.RowsAffected == 1
			if err := tx.Exec("UPDATE user_progress SET best_score=GREATEST(best_score,?) WHERE user_id=? AND activity_id=?", result.Score, userID, a.ID).Error; err != nil {
				return err
			}
			if result.NewCompletion {
				factor, err := activeMultiplier(tx, userID)
				if err != nil {
					return err
				}
				result.Multiplier = factor
				xp := a.RewardXP(result.Score == 100, factor)
				gained, err := award(tx, userID, "activity:"+a.ID, xp)
				if err != nil {
					return err
				}
				result.XP += gained
				bonus, err := applyMissions(tx, userID, today, a, result.Score == 100)
				if err != nil {
					return err
				}
				result.XP += bonus
				if err := tx.Exec(`INSERT INTO user_skill_stats(user_id,words_practiced) VALUES (?,?) ON CONFLICT(user_id) DO UPDATE SET words_practiced=user_skill_stats.words_practiced+EXCLUDED.words_practiced`, userID, a.VocabularyCount).Error; err != nil {
					return err
				}
			}
			var old struct {
				CurrentDays, BestDays int
				LastDay               string
			}
			if err := tx.Raw("SELECT current_days,best_days,COALESCE(to_char(last_day,'YYYY-MM-DD'),'') last_day FROM streaks WHERE user_id=?", userID).Scan(&old).Error; err != nil {
				return err
			}
			current, best := streak(old.LastDay, today, old.CurrentDays, old.BestDays)
			if err := tx.Exec(`INSERT INTO streaks(user_id,current_days,best_days,last_day) VALUES (?,?,?,?::date)
    ON CONFLICT(user_id) DO UPDATE SET current_days=EXCLUDED.current_days,best_days=EXCLUDED.best_days,last_day=EXCLUDED.last_day`, userID, current, best, today).Error; err != nil {
				return err
			}
			if err := tx.Exec("INSERT INTO practice_days VALUES (?,?::date) ON CONFLICT DO NOTHING", userID, today).Error; err != nil {
				return err
			}
			var total int64
			if err := tx.Raw("SELECT count(*) FROM user_progress WHERE user_id=?", userID).Scan(&total).Error; err != nil {
				return err
			}
			codes := []string{"first_activity"}
			if total >= 10 {
				codes = append(codes, "ten_activities")
			}
			if a.Assessment {
				codes = append(codes, "first_assessment")
			}
			if current >= 7 {
				codes = append(codes, "week_streak")
			}
			var words int
			if err := tx.Raw("SELECT words_practiced FROM user_skill_stats WHERE user_id=?", userID).Scan(&words).Error; err != nil {
				return err
			}
			if words >= 100 {
				codes = append(codes, "hundred_words")
			}
			var levelComplete bool
			if err := tx.Raw(`SELECT count(*)>0 AND bool_and(p.user_id IS NOT NULL) FROM activities a JOIN lessons l ON l.id=a.lesson_id JOIN units u ON u.id=l.unit_id
 LEFT JOIN user_progress p ON p.activity_id=a.id AND p.user_id=? WHERE a.published AND u.published AND u.level_code=(SELECT u2.level_code FROM units u2 JOIN lessons l2 ON l2.unit_id=u2.id WHERE l2.id=?)`, userID, a.LessonID).Scan(&levelComplete).Error; err != nil {
				return err
			}
			if levelComplete {
				codes = append(codes, "level_path")
			}
			for _, code := range codes {
				if err := tx.Exec("INSERT INTO user_achievements VALUES (?,?,now()) ON CONFLICT DO NOTHING", userID, code).Error; err != nil {
					return err
				}
			}
		}
		if result.Passed {
			if err := syncProgressUnlocks(tx, userID); err != nil {
				return err
			}
		}
		raw, err := json.Marshal(result)
		if err != nil {
			return err
		}
		submitted, err := json.Marshal(answers)
		if err != nil {
			return err
		}
		return tx.Exec("UPDATE activity_attempts SET score=?,answers=?::jsonb,result=?::jsonb,completed_at=now() WHERE id=?", result.Score, string(submitted), string(raw), id).Error
	})
	return result, err
}

func (r Learning) Dashboard(ctx context.Context, user *models.User) (map[string]interface{}, error) {
	db := r.DB.WithContext(ctx)
	out := map[string]interface{}{}
	lessons, err := lessonList(db, user.ID)
	if err != nil {
		return nil, err
	}
	out["lessons"] = lessons
	var xp int64
	if err := db.Raw("SELECT COALESCE(sum(amount),0) FROM xp_transactions WHERE user_id=?", user.ID).Scan(&xp).Error; err != nil {
		return nil, err
	}
	out["xp"] = xp
	var streak struct {
		CurrentDays int    `json:"current_days"`
		BestDays    int    `json:"best_days"`
		LastDay     string `json:"last_day"`
	}
	if err := db.Raw("SELECT current_days,best_days,COALESCE(to_char(last_day,'YYYY-MM-DD'),'') last_day FROM streaks WHERE user_id=?", user.ID).Scan(&streak).Error; err != nil {
		return nil, err
	}
	loc, err := time.LoadLocation(user.Timezone)
	if err != nil {
		return nil, err
	}
	today := time.Now().In(loc)
	day := today.Format("2006-01-02")
	if streak.LastDay != day && streak.LastDay != today.AddDate(0, 0, -1).Format("2006-01-02") {
		streak.CurrentDays = 0
	}
	out["streak"] = streak
	days := []string{}
	if err := db.Raw("SELECT to_char(day,'YYYY-MM-DD') FROM practice_days WHERE user_id=? AND day>=?::date ORDER BY day", user.ID, today.AddDate(0, 0, -6).Format("2006-01-02")).Scan(&days).Error; err != nil {
		return nil, err
	}
	out["practice_days"] = days
	out["today"] = day
	missions := []map[string]interface{}{}
	if err := db.Raw(`SELECT d.*,COALESCE(m.progress,0) progress,COALESCE(m.rewarded,false) rewarded FROM daily_missions d LEFT JOIN user_daily_missions m ON m.mission_code=d.code AND m.user_id=? AND m.day=?::date`, user.ID, day).Scan(&missions).Error; err != nil {
		return nil, err
	}
	out["missions"] = missions
	achievements := []map[string]interface{}{}
	if err := db.Raw(`SELECT a.*,u.earned_at,(u.user_id IS NOT NULL) unlocked FROM achievements a LEFT JOIN user_achievements u ON u.achievement_code=a.code AND u.user_id=? ORDER BY a.threshold,a.code`, user.ID).Scan(&achievements).Error; err != nil {
		return nil, err
	}
	out["achievements"] = achievements
	history := []map[string]interface{}{}
	if err := db.Raw(`SELECT t.id,a.title,t.score,t.completed_at FROM activity_attempts t JOIN activities a ON a.id=t.activity_id WHERE t.user_id=? AND t.completed_at IS NOT NULL ORDER BY t.completed_at DESC LIMIT 20`, user.ID).Scan(&history).Error; err != nil {
		return nil, err
	}
	out["history"] = history
	if err := addRewardDashboard(db, user.ID, out); err != nil {
		return nil, err
	}
	return out, nil
}
