package repositories

import (
	"context"
	"github.com/englishquest/backend/internal/models"
	apperr "github.com/englishquest/backend/internal/pkg/errors"
	"gorm.io/gorm"
)

type LevelAccess struct {
	Code          string `json:"code"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Position      int    `json:"position"`
	Available     bool   `json:"available"`
	UnitCount     int    `json:"unit_count"`
	LessonCount   int    `json:"lesson_count"`
	Completed     int    `json:"completed"`
	Unlocked      bool   `json:"unlocked"`
	Source        string `json:"source"`
	PreviousLevel string `json:"previous_level"`
	CanTakeExam   bool   `json:"can_take_exam"`
}

func levelAccess(db *gorm.DB, userID string) ([]LevelAccess, error) {
	levels := []LevelAccess{}
	err := db.Raw(`SELECT c.code,c.title,c.description,c.position,c.available,
 count(DISTINCT u.id) FILTER(WHERE a.id IS NOT NULL) unit_count,count(a.id) lesson_count,count(p.activity_id) completed,
 COALESCE(g.source,'') source FROM cefr_levels c
 LEFT JOIN units u ON u.level_code=c.code AND u.published
 LEFT JOIN lessons l ON l.unit_id=u.id LEFT JOIN activities a ON a.lesson_id=l.id AND a.published
 LEFT JOIN user_progress p ON p.activity_id=a.id AND p.user_id=?
 LEFT JOIN user_level_unlocks g ON g.level_code=c.code AND g.user_id=?
 GROUP BY c.code,c.title,c.description,c.position,c.available,g.source ORDER BY c.position`, userID, userID).Scan(&levels).Error
	if err != nil {
		return nil, err
	}
	for i := range levels {
		l := &levels[i]
		l.Unlocked = l.Available && (l.Position <= 3 || l.Source != "")
		if l.Position <= 3 {
			l.Source = "initial"
		}
		if i > 0 {
			previous := levels[i-1]
			l.PreviousLevel = previous.Code
			if l.Available && previous.Available && previous.Unlocked && previous.LessonCount > 0 && previous.Completed == previous.LessonCount {
				l.Unlocked = true
				if l.Source == "" {
					l.Source = "progress"
				}
			}
			l.CanTakeExam = l.Available && !l.Unlocked && previous.Available && previous.Unlocked
		}
	}
	return levels, nil
}
func (r Learning) LevelAccess(ctx context.Context, userID string) ([]LevelAccess, error) {
	return levelAccess(r.DB.WithContext(ctx), userID)
}

func syncProgressUnlocks(tx *gorm.DB, userID string) error {
	levels, err := levelAccess(tx, userID)
	if err != nil {
		return err
	}
	for _, level := range levels {
		if level.Position > 3 && level.Unlocked && level.Source == "progress" {
			if err := tx.Exec("INSERT INTO user_level_unlocks(user_id,level_code,source) VALUES (?,?,'progress') ON CONFLICT DO NOTHING", userID, level.Code).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func requireUnitAccess(db *gorm.DB, userID, unitID string) error {
	var code string
	if err := db.Raw("SELECT level_code FROM units WHERE id=? AND published", unitID).Scan(&code).Error; err != nil {
		return err
	}
	if code == "" {
		return apperr.NewAppError(404, "Unidad no disponible.")
	}
	levels, err := levelAccess(db, userID)
	if err != nil {
		return err
	}
	for _, level := range levels {
		if level.Code == code && level.Unlocked {
			return nil
		}
	}
	return apperr.NewAppError(403, "Desbloquea este nivel completando el anterior o aprobando el examen de acceso.")
}

func applyLevelLocks(db *gorm.DB, userID string, lessons []models.LessonSummary) error {
	levels, err := levelAccess(db, userID)
	if err != nil {
		return err
	}
	access := map[string]bool{}
	for _, l := range levels {
		access[l.Code] = l.Unlocked
	}
	for i := range lessons {
		lessons[i].LevelLocked = !access[lessons[i].LevelCode]
		if lessons[i].LevelLocked {
			lessons[i].Locked = true
			lessons[i].Content = ""
		}
	}
	return nil
}
