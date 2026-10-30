package repositories

import (
	"context"
	"encoding/json"
	apperr "github.com/englishquest/backend/internal/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type CourseDiploma struct {
	ID            string         `json:"id"`
	RecipientName string         `json:"recipient_name"`
	IssuedAt      time.Time      `json:"issued_at"`
	CourseTitle   string         `json:"course_title"`
	ActivityCount int            `json:"activity_count"`
	LevelCounts   string         `json:"-"`
	Levels        map[string]int `json:"levels" gorm:"-"`
}
type DiplomaStatus struct {
	Eligible  bool           `json:"eligible"`
	Total     int            `json:"total"`
	Completed int            `json:"completed"`
	Diploma   *CourseDiploma `json:"diploma"`
}
type diplomaActivity struct {
	ID, LevelCode        string
	Completed, Available bool
}

func diplomaRequirements(db *gorm.DB, userID string) (DiplomaStatus, []string, map[string]int, error) {
	rows := []diplomaActivity{}
	err := db.Raw(`SELECT a.id,u.level_code,c.available,(p.activity_id IS NOT NULL) completed FROM activities a
 JOIN lessons l ON l.id=a.lesson_id JOIN units u ON u.id=l.unit_id AND u.published JOIN cefr_levels c ON c.code=u.level_code
 LEFT JOIN user_progress p ON p.activity_id=a.id AND p.user_id=? WHERE a.published ORDER BY c.position,u.position,l.position,a.id`, userID).Scan(&rows).Error
	out := DiplomaStatus{Total: len(rows), Eligible: len(rows) > 0}
	ids := []string{}
	counts := map[string]int{}
	if err != nil {
		return out, ids, counts, err
	}
	for _, row := range rows {
		ids = append(ids, row.ID)
		counts[row.LevelCode]++
		if row.Completed {
			out.Completed++
		}
		if !row.Available {
			out.Eligible = false
		}
	}
	out.Eligible = out.Eligible && len(counts) == 6 && out.Completed == out.Total
	return out, ids, counts, nil
}
func readDiploma(db *gorm.DB, userID string) (*CourseDiploma, error) {
	var out CourseDiploma
	if err := db.Raw("SELECT id,recipient_name,issued_at,course_title,activity_count,level_counts::text FROM course_diplomas WHERE user_id=?", userID).Scan(&out).Error; err != nil {
		return nil, err
	}
	if out.ID == "" {
		return nil, nil
	}
	if err := json.Unmarshal([]byte(out.LevelCounts), &out.Levels); err != nil {
		return nil, err
	}
	return &out, nil
}
func (r Learning) DiplomaStatus(ctx context.Context, userID string) (DiplomaStatus, error) {
	db := r.DB.WithContext(ctx)
	out, _, _, err := diplomaRequirements(db, userID)
	if err != nil {
		return out, err
	}
	out.Diploma, err = readDiploma(db, userID)
	return out, err
}

// name must first pass services.DiplomaName at the handler boundary.
func (r Learning) IssueDiploma(ctx context.Context, userID, name string) (*CourseDiploma, error) {
	var out *CourseDiploma
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT id FROM users WHERE id=? FOR NO KEY UPDATE", userID).Error; err != nil {
			return err
		}
		old, err := readDiploma(tx, userID)
		if err != nil {
			return err
		}
		if old != nil {
			if old.RecipientName != name {
				return apperr.NewAppError(409, "Tu diploma ya fue emitido con otro nombre. No se cambia automáticamente.")
			}
			out = old
			return nil
		}
		status, ids, counts, err := diplomaRequirements(tx, userID)
		if err != nil {
			return err
		}
		if !status.Eligible {
			return apperr.NewAppError(409, "Primero aprueba todas las actividades publicadas de los seis niveles. Un examen de acceso no completa el curso.")
		}
		idJSON, _ := json.Marshal(ids)
		countJSON, _ := json.Marshal(counts)
		if err := tx.Exec(`INSERT INTO course_diplomas(user_id,recipient_name,course_title,activity_count,activity_ids,level_counts)
   VALUES (?,?,'Ruta guiada de práctica A1–C2',?,?::jsonb,?::jsonb)`, userID, name, status.Total, string(idJSON), string(countJSON)).Error; err != nil {
			return err
		}
		out, err = readDiploma(tx, userID)
		return err
	})
	return out, err
}
