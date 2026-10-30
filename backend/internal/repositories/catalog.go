package repositories

import (
	"context"
	"github.com/englishquest/backend/internal/models"
	apperr "github.com/englishquest/backend/internal/pkg/errors"
	"gorm.io/gorm"
)

func audit(tx *gorm.DB, actor, action, target string) error {
	// Existing UUID target retained for compatibility; target_key records non-UUID catalogue keys.
	return tx.Exec("INSERT INTO admin_audit(actor_id,action,resource_id,target_key) VALUES (?,?,?,?)", actor, action, actor, target).Error
}
func (r Learning) SaveLevel(ctx context.Context, actor, code string, in models.LevelEdit) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing string
		if err := tx.Raw("SELECT code FROM cefr_levels WHERE code=? FOR UPDATE", code).Scan(&existing).Error; err != nil {
			return err
		}
		if existing == "" {
			return apperr.NewAppError(404, "Nivel no encontrado.")
		}
		if in.Available {
			var count int64
			if err := tx.Raw(`SELECT count(*) FROM units u JOIN lessons l ON l.unit_id=u.id JOIN activities a ON a.lesson_id=l.id AND a.published
    WHERE u.level_code=? AND u.published AND EXISTS(SELECT 1 FROM questions q WHERE q.activity_id=a.id)`, code).Scan(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				return apperr.NewAppError(409, "Añade contenido revisado y publicado antes de habilitar el nivel.")
			}
		}
		if err := tx.Exec("UPDATE cefr_levels SET title=?,description=?,available=? WHERE code=?", in.Title, in.Description, in.Available, code).Error; err != nil {
			return err
		}
		return audit(tx, actor, "save_level", code)
	})
}
func (r Learning) SaveUnit(ctx context.Context, actor, id string, in models.UnitEdit) (string, error) {
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var level string
		if err := tx.Raw("SELECT code FROM cefr_levels WHERE code=? FOR UPDATE", in.LevelCode).Scan(&level).Error; err != nil {
			return err
		}
		if level == "" {
			return apperr.NewAppError(404, "Nivel no encontrado.")
		}
		if id == "" {
			if err := tx.Raw(`INSERT INTO units(level_code,title,description,published,position) SELECT ?,?,?,?,COALESCE(max(position),0)+1 FROM units WHERE level_code=? RETURNING id`, in.LevelCode, in.Title, in.Description, in.Published, in.LevelCode).Scan(&id).Error; err != nil {
				return err
			}
		} else {
			res := tx.Exec("UPDATE units SET title=?,description=?,published=? WHERE id=? AND level_code=?", in.Title, in.Description, in.Published, id, in.LevelCode)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected != 1 {
				return apperr.NewAppError(409, "Unidad inexistente o intento de cambiarla de nivel. Conserva su nivel original.")
			}
		}
		return audit(tx, actor, "save_unit", id)
	})
	return id, err
}
func (r Learning) SaveResource(ctx context.Context, actor string, in models.Resource) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var old models.Resource
		if err := tx.Raw("SELECT * FROM educational_resources WHERE key=? FOR UPDATE", in.Key).Scan(&old).Error; err != nil {
			return err
		}
		if old.Key != "" && old.Kind != in.Kind {
			return apperr.NewAppError(409, "No puedes cambiar el tipo de un recurso existente; crea otra clave.")
		}
		if err := tx.Exec(`INSERT INTO educational_resources(key,title,kind,url,attribution,active) VALUES (?,?,?,?,?,?) ON CONFLICT(key) DO UPDATE SET
   title=EXCLUDED.title,url=EXCLUDED.url,attribution=EXCLUDED.attribution,active=EXCLUDED.active`, in.Key, in.Title, in.Kind, in.URL, in.Attribution, in.Active).Error; err != nil {
			return err
		}
		return audit(tx, actor, "save_resource", in.Key)
	})
}
