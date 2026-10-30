package repositories

import (
	"context"
	"github.com/englishquest/backend/internal/models"
	apperr "github.com/englishquest/backend/internal/pkg/errors"
	"gorm.io/gorm"
)

func (r Learning) AdminOverview(ctx context.Context) (map[string]interface{}, error) {
	db := r.DB.WithContext(ctx)
	out := map[string]interface{}{}
	for key, query := range map[string]string{
		"levels":     "SELECT * FROM cefr_levels ORDER BY position",
		"resources":  "SELECT * FROM educational_resources ORDER BY kind,title",
		"users":      "SELECT id,email,first_name,role,is_active,created_at FROM users ORDER BY created_at DESC LIMIT 100",
		"units":      "SELECT u.*,c.title level_title FROM units u JOIN cefr_levels c ON c.code=u.level_code ORDER BY c.position,u.position",
		"activities": "SELECT a.*,l.title lesson_title,l.objective,l.content,u.title unit_title,u.level_code FROM activities a JOIN lessons l ON l.id=a.lesson_id JOIN units u ON u.id=l.unit_id ORDER BY u.level_code,u.position,l.position",
	} {
		rows := []map[string]interface{}{}
		if err := db.Raw(query).Scan(&rows).Error; err != nil {
			return nil, err
		}
		out[key] = rows
	}
	stats := map[string]interface{}{}
	if err := db.Raw("SELECT (SELECT count(*) FROM users) users,(SELECT count(*) FROM activity_attempts WHERE completed_at IS NOT NULL) attempts,(SELECT count(*) FROM user_progress) completions").Scan(&stats).Error; err != nil {
		return nil, err
	}
	out["stats"] = stats
	return out, nil
}
func (r Learning) AdminActivity(ctx context.Context, id string) (models.Activity, error) {
	return loadActivity(r.DB.WithContext(ctx), id)
}

type ContentEdit struct {
	Activity    models.Activity `json:"activity"`
	UnitID      string          `json:"unit_id"`
	LessonTitle string          `json:"lesson_title"`
	Objective   string          `json:"objective"`
	Content     string          `json:"content"`
}

func saveQuestions(tx *gorm.DB, id string, questions []models.Question) error {
	if err := tx.Exec("DELETE FROM questions WHERE activity_id=?", id).Error; err != nil {
		return err
	}
	for i, q := range questions {
		var qid string
		if q.AudioKey != "" {
			var count int64
			if err := tx.Raw("SELECT count(*) FROM educational_resources WHERE key=? AND kind='audio' AND active", q.AudioKey).Scan(&count).Error; err != nil {
				return err
			}
			if count != 1 {
				return apperr.NewAppError(400, "El audio no existe o está inactivo.")
			}
		}
		if err := tx.Raw("INSERT INTO questions(activity_id,prompt,explanation,solution,position,response,passage,audio_key) VALUES (?,?,?,?::jsonb,?,?,?,NULLIF(?,'')) RETURNING id", id, q.Prompt, q.Explanation, string(q.Solution), i+1, q.Response, q.Passage, q.AudioKey).Scan(&qid).Error; err != nil {
			return err
		}
		for _, o := range q.Options {
			if o.Image != "" {
				var count int64
				if err := tx.Raw("SELECT count(*) FROM educational_resources WHERE key=? AND kind='image' AND active", o.Image).Scan(&count).Error; err != nil {
					return err
				}
				if count != 1 {
					return apperr.NewAppError(400, "La imagen no existe o está inactiva.")
				}
			}
			if err := tx.Exec("INSERT INTO answer_options(id,question_id,label,image,image_resource_key) VALUES (?,?,?,?,NULLIF(?,''))", o.ID, qid, o.Label, o.Image, o.Image).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// Attempts retain immutable snapshots, so editing never changes a running score.
func (r Learning) SaveContent(ctx context.Context, actor, id string, in ContentEdit) (string, error) {
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var lessonID string
		if id == "" {
			var unit string
			if err := tx.Raw("SELECT id FROM units WHERE id=? FOR UPDATE", in.UnitID).Scan(&unit).Error; err != nil {
				return err
			}
			if unit == "" {
				return apperr.NewAppError(404, "Unidad no encontrada.")
			}
			if err := tx.Raw("INSERT INTO lessons(unit_id,title,objective,content,position) SELECT ?,?,?,?,COALESCE(max(position),0)+1 FROM lessons WHERE unit_id=? RETURNING id", in.UnitID, in.LessonTitle, in.Objective, in.Content, in.UnitID).Scan(&lessonID).Error; err != nil {
				return err
			}
			if err := tx.Raw("INSERT INTO activities(lesson_id,title,kind,instructions,published,assessment,time_limit_seconds,vocabulary_count,skill) VALUES (?,?,?,?,?,?,?,?,?) RETURNING id", lessonID, in.Activity.Title, in.Activity.Kind, in.Activity.Instructions, in.Activity.Published, in.Activity.Assessment, in.Activity.TimeLimitSeconds, in.Activity.VocabularyCount, in.Activity.LearningSkill()).Scan(&id).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Raw("SELECT lesson_id FROM activities WHERE id=? FOR UPDATE", id).Scan(&lessonID).Error; err != nil {
				return err
			}
			if lessonID == "" {
				return apperr.NewAppError(404, "Actividad no encontrada.")
			}
			if err := tx.Exec("UPDATE lessons SET title=?,objective=?,content=? WHERE id=?", in.LessonTitle, in.Objective, in.Content, lessonID).Error; err != nil {
				return err
			}
			if err := tx.Exec("UPDATE activities SET title=?,kind=?,instructions=?,published=?,assessment=?,time_limit_seconds=?,vocabulary_count=?,skill=?,version=version+1 WHERE id=?", in.Activity.Title, in.Activity.Kind, in.Activity.Instructions, in.Activity.Published, in.Activity.Assessment, in.Activity.TimeLimitSeconds, in.Activity.VocabularyCount, in.Activity.LearningSkill(), id).Error; err != nil {
				return err
			}
		}
		if err := saveQuestions(tx, id, in.Activity.Questions); err != nil {
			return err
		}
		return tx.Exec("INSERT INTO admin_audit(actor_id,action,resource_id) VALUES (?,'save_content',?)", actor, id).Error
	})
	return id, err
}
func (r Learning) SetActive(ctx context.Context, actor, id string, active bool) error {
	if actor == id {
		return apperr.NewAppError(400, "No puedes desactivar tu propia cuenta.")
	}
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Exec("UPDATE users SET is_active=?,updated_at=now() WHERE id=? AND role='STUDENT'", active, id)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return apperr.NewAppError(404, "Estudiante no encontrado.")
		}
		if !active {
			if err := tx.Exec("UPDATE refresh_tokens SET revoked=true WHERE user_id=?", id).Error; err != nil {
				return err
			}
		}
		return tx.Exec("INSERT INTO admin_audit(actor_id,action,resource_id) VALUES (?,'set_user_active',?)", actor, id).Error
	})
}
