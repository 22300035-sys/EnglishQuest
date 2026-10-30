package repositories

import (
	"context"
	apperr "github.com/englishquest/backend/internal/pkg/errors"
	"strings"
	"time"
	"unicode/utf8"
)

type UnitProject struct {
	UnitID     string     `json:"unit_id"`
	Prompt     string     `json:"prompt"`
	Rubric     string     `json:"rubric"`
	Body       string     `json:"body"`
	Reflection string     `json:"reflection"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

func (r Learning) UnitProject(ctx context.Context, userID, unitID string) (UnitProject, error) {
	var out UnitProject
	if err := requireUnitAccess(r.DB.WithContext(ctx), userID, unitID); err != nil {
		return out, err
	}
	err := r.DB.WithContext(ctx).Raw(`SELECT p.unit_id,p.prompt,p.rubric,COALESCE(d.body,'') body,COALESCE(d.reflection,'') reflection,d.updated_at
 FROM unit_projects p JOIN units u ON u.id=p.unit_id AND u.published JOIN cefr_levels c ON c.code=u.level_code AND c.available
 LEFT JOIN user_unit_projects d ON d.unit_id=p.unit_id AND d.user_id=? WHERE p.unit_id=?`, userID, unitID).Scan(&out).Error
	if err != nil {
		return out, err
	}
	if out.UnitID == "" {
		return out, apperr.NewAppError(404, "Esta unidad no tiene un proyecto disponible.")
	}
	return out, nil
}
func (r Learning) SaveUnitProject(ctx context.Context, userID, unitID, body, reflection string) error {
	if err := requireUnitAccess(r.DB.WithContext(ctx), userID, unitID); err != nil {
		return err
	}
	if utf8.RuneCountInString(body) > 12000 || utf8.RuneCountInString(reflection) > 2000 || strings.TrimSpace(body) == "" {
		return apperr.NewAppError(400, "Escribe tu borrador (máximo 12000 caracteres) y una reflexión opcional de hasta 2000 caracteres.")
	}
	res := r.DB.WithContext(ctx).Exec(`INSERT INTO user_unit_projects(user_id,unit_id,body,reflection)
 SELECT ?,p.unit_id,?,? FROM unit_projects p JOIN units u ON u.id=p.unit_id AND u.published JOIN cefr_levels c ON c.code=u.level_code AND c.available WHERE p.unit_id=?
 ON CONFLICT(user_id,unit_id) DO UPDATE SET body=EXCLUDED.body,reflection=EXCLUDED.reflection,updated_at=now()`, userID, body, reflection, unitID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return apperr.NewAppError(404, "Esta unidad no tiene un proyecto disponible.")
	}
	return nil
}
