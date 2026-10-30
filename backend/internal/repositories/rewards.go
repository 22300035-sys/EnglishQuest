package repositories

import (
	"context"
	"github.com/englishquest/backend/internal/models"
	apperr "github.com/englishquest/backend/internal/pkg/errors"
	"gorm.io/gorm"
	"time"
)

func applyMissions(tx *gorm.DB, userID, day string, a models.Activity, perfect bool) (int, error) {
	var missions []struct {
		Code, Metric   string
		Target, Reward int
	}
	if err := tx.Raw("SELECT code,metric,target,reward FROM daily_missions ORDER BY code").Scan(&missions).Error; err != nil {
		return 0, err
	}
	total := 0
	for _, m := range missions {
		if m.Metric != "new_completion" && m.Metric != a.LearningSkill() && !(m.Metric == "perfect" && perfect) {
			continue
		}
		if err := tx.Exec(`INSERT INTO user_daily_missions(user_id,mission_code,day,progress) VALUES (?,?,?::date,1)
   ON CONFLICT(user_id,mission_code,day) DO UPDATE SET progress=user_daily_missions.progress+1`, userID, m.Code, day).Error; err != nil {
			return 0, err
		}
		var progress int
		if err := tx.Raw("SELECT progress FROM user_daily_missions WHERE user_id=? AND mission_code=? AND day=?::date", userID, m.Code, day).Scan(&progress).Error; err != nil {
			return 0, err
		}
		if progress < m.Target {
			continue
		}
		source := "mission:" + m.Code + ":" + day
		if m.Code == "three_activities" {
			source = "mission:" + day
		} // preserve old idempotency keys
		bonus, err := award(tx, userID, source, m.Reward)
		if err != nil {
			return 0, err
		}
		total += bonus
		if err := tx.Exec("UPDATE user_daily_missions SET rewarded=true WHERE user_id=? AND mission_code=? AND day=?::date", userID, m.Code, day).Error; err != nil {
			return 0, err
		}
		if m.Code == "three_activities" {
			if err := tx.Exec("INSERT INTO user_boosters(user_id,booster_code,source) VALUES (?,'double_xp_15m',?) ON CONFLICT DO NOTHING", userID, "daily:"+day).Error; err != nil {
				return 0, err
			}
		}
	}
	return total, nil
}

func activeMultiplier(tx *gorm.DB, userID string) (int, error) {
	var factor int
	err := tx.Raw(`SELECT d.multiplier FROM user_boosters b JOIN xp_boosters d ON d.code=b.booster_code WHERE b.user_id=? AND b.active AND b.expires_at>now() LIMIT 1`, userID).Scan(&factor).Error
	if factor != 2 {
		factor = 1
	}
	return factor, err
}

func (r Learning) ActivateBooster(ctx context.Context, userID, id string) (map[string]interface{}, error) {
	out := map[string]interface{}{}
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT id FROM users WHERE id=? FOR NO KEY UPDATE", userID).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE user_boosters SET active=false WHERE user_id=? AND active AND expires_at<=now()", userID).Error; err != nil {
			return err
		}
		var active int64
		if err := tx.Raw("SELECT count(*) FROM user_boosters WHERE user_id=? AND active", userID).Scan(&active).Error; err != nil {
			return err
		}
		if active > 0 {
			return apperr.NewAppError(409, "Ya tienes un potenciador activo. No pueden acumularse.")
		}
		var row struct {
			ID              string
			ActivatedAt     *time.Time
			DurationSeconds int
		}
		if err := tx.Raw(`SELECT b.id,b.activated_at,d.duration_seconds FROM user_boosters b JOIN xp_boosters d ON d.code=b.booster_code WHERE b.id=? AND b.user_id=? FOR UPDATE OF b`, id, userID).Scan(&row).Error; err != nil {
			return err
		}
		if row.ID == "" {
			return apperr.NewAppError(404, "Potenciador no encontrado.")
		}
		if row.ActivatedAt != nil {
			return apperr.NewAppError(409, "Este potenciador ya fue utilizado.")
		}
		now := time.Now().UTC()
		end := now.Add(time.Duration(row.DurationSeconds) * time.Second)
		if err := tx.Exec("UPDATE user_boosters SET active=true,activated_at=?,expires_at=? WHERE id=?", now, end, id).Error; err != nil {
			return err
		}
		out["id"] = id
		out["expires_at"] = end
		out["multiplier"] = 2
		return nil
	})
	return out, err
}

func addRewardDashboard(tx *gorm.DB, userID string, out map[string]interface{}) error {
	boosters := []map[string]interface{}{}
	if err := tx.Raw(`SELECT b.id,d.title,d.multiplier,d.duration_seconds,b.earned_at,b.activated_at,b.expires_at,
 (b.active AND b.expires_at>now()) active,(b.activated_at IS NULL) available
 FROM user_boosters b JOIN xp_boosters d ON d.code=b.booster_code WHERE b.user_id=? ORDER BY b.earned_at DESC LIMIT 30`, userID).Scan(&boosters).Error; err != nil {
		return err
	}
	out["boosters"] = boosters
	var words int
	if err := tx.Raw("SELECT words_practiced FROM user_skill_stats WHERE user_id=?", userID).Scan(&words).Error; err != nil {
		return err
	}
	out["words_practiced"] = words
	out["server_time"] = time.Now().UTC()
	return nil
}
