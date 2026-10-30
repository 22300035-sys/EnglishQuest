package repositories

import (
	"context"
	"encoding/json"
	"github.com/englishquest/backend/internal/models"
	apperr "github.com/englishquest/backend/internal/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type examRow struct {
	ID, TargetLevel, Snapshot, Result string
	StartedAt, ExpiresAt              time.Time
	CompletedAt                       *time.Time
}
type AccessExamView struct {
	AttemptView
	TargetLevel string     `json:"target_level"`
	RetryAt     *time.Time `json:"retry_at,omitempty"`
}

func examView(row examRow) (AccessExamView, error) {
	var activity models.Activity
	if err := json.Unmarshal([]byte(row.Snapshot), &activity); err != nil {
		return AccessExamView{}, err
	}
	out := AccessExamView{AttemptView: publicAttempt(row.ID, activity, row.StartedAt), TargetLevel: row.TargetLevel}
	out.ExpiresAt = &row.ExpiresAt
	if row.Result != "" {
		out.Result = &models.Result{}
		if err := json.Unmarshal([]byte(row.Result), out.Result); err != nil {
			return out, err
		}
		if !out.Result.Passed && row.CompletedAt != nil {
			retry := row.CompletedAt.Add(10 * time.Minute)
			out.RetryAt = &retry
		}
	}
	return out, nil
}

func (r Learning) StartAccessExam(ctx context.Context, userID, level string, build func(string) (models.Activity, error)) (AccessExamView, error) {
	var row examRow
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT id FROM users WHERE id=? FOR NO KEY UPDATE", userID).Error; err != nil {
			return err
		}
		levels, err := levelAccess(tx, userID)
		if err != nil {
			return err
		}
		found := false
		for _, state := range levels {
			if state.Code == level {
				found = true
				if state.Unlocked {
					return apperr.NewAppError(409, "Ya tienes este nivel desbloqueado.")
				}
				if !state.CanTakeExam {
					return apperr.NewAppError(403, "Primero debes tener abierto el nivel anterior.")
				}
			}
		}
		if !found {
			return apperr.NewAppError(404, "Nivel no encontrado.")
		}
		if err := tx.Raw("SELECT id,target_level,snapshot::text,COALESCE(result::text,'') result,started_at,expires_at,completed_at FROM access_exam_attempts WHERE user_id=? AND target_level=? AND completed_at IS NULL", userID, level).Scan(&row).Error; err != nil {
			return err
		}
		if row.ID != "" {
			return nil
		} // Resume, including expiry; submit to close an expired attempt.
		var last struct{ At *time.Time }
		if err := tx.Raw("SELECT max(completed_at) at FROM access_exam_attempts WHERE user_id=? AND target_level=?", userID, level).Scan(&last).Error; err != nil {
			return err
		}
		if last.At != nil && time.Now().Before(last.At.Add(10*time.Minute)) {
			return apperr.NewAppError(429, "Espera 10 minutos desde el último examen antes de reintentarlo.")
		}
		activity, err := build(level)
		if err != nil {
			return err
		}
		if err := anonymize(&activity); err != nil {
			return err
		}
		snapshot, err := json.Marshal(activity)
		if err != nil {
			return err
		}
		return tx.Raw(`INSERT INTO access_exam_attempts(user_id,target_level,snapshot,expires_at) VALUES (?,?,?::jsonb,now()+interval '20 minutes')
  RETURNING id,target_level,snapshot::text,'' result,started_at,expires_at,completed_at`, userID, level, string(snapshot)).Scan(&row).Error
	})
	if err != nil {
		return AccessExamView{}, err
	}
	return examView(row)
}
func (r Learning) AccessExam(ctx context.Context, userID, id string) (AccessExamView, error) {
	var row examRow
	if err := r.DB.WithContext(ctx).Raw("SELECT id,target_level,snapshot::text,COALESCE(result::text,'') result,started_at,expires_at,completed_at FROM access_exam_attempts WHERE id=? AND user_id=?", id, userID).Scan(&row).Error; err != nil {
		return AccessExamView{}, err
	}
	if row.ID == "" {
		return AccessExamView{}, apperr.NewAppError(404, "Examen no encontrado.")
	}
	return examView(row)
}
func (r Learning) SubmitAccessExam(ctx context.Context, userID, id string, in models.Submission, grade func(models.Activity, models.Submission) (models.Result, error)) (AccessExamView, error) {
	var row examRow
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT id FROM users WHERE id=? FOR NO KEY UPDATE", userID).Error; err != nil {
			return err
		}
		if err := tx.Raw("SELECT id,target_level,snapshot::text,COALESCE(result::text,'') result,started_at,expires_at,completed_at FROM access_exam_attempts WHERE id=? AND user_id=? FOR UPDATE", id, userID).Scan(&row).Error; err != nil {
			return err
		}
		if row.ID == "" {
			return apperr.NewAppError(404, "Examen no encontrado.")
		}
		if row.CompletedAt != nil {
			return nil
		}
		var activity models.Activity
		if err := json.Unmarshal([]byte(row.Snapshot), &activity); err != nil {
			return err
		}
		result, err := grade(activity, in)
		if err != nil {
			return err
		}
		result.AttemptID = id
		result.XP = 0
		result.NewCompletion = false
		result.ElapsedSeconds = int(time.Since(row.StartedAt).Seconds())
		if time.Now().After(row.ExpiresAt) {
			result.Expired = true
			result.Score = 0
			result.Passed = false
			result.Stars = 0
		}
		if result.Passed {
			if err := tx.Exec("INSERT INTO user_level_unlocks(user_id,level_code,source) VALUES (?,?,'exam') ON CONFLICT DO NOTHING", userID, row.TargetLevel).Error; err != nil {
				return err
			}
		}
		answers, err := json.Marshal(in)
		if err != nil {
			return err
		}
		raw, err := json.Marshal(result)
		if err != nil {
			return err
		}
		now := time.Now().UTC()
		row.Result = string(raw)
		row.CompletedAt = &now
		return tx.Exec("UPDATE access_exam_attempts SET answers=?::jsonb,result=?::jsonb,completed_at=? WHERE id=?", string(answers), row.Result, now, id).Error
	})
	if err != nil {
		return AccessExamView{}, err
	}
	return examView(row)
}
