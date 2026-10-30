package repositories

import (
	"context"
	"github.com/englishquest/backend/internal/models"
	apperr "github.com/englishquest/backend/internal/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type Accounts struct{ DB *gorm.DB }

func (r Accounts) Create(ctx context.Context, user *models.User) error {
	return r.DB.WithContext(ctx).Create(user).Error
}
func (r Accounts) ByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return &u, err
}
func (r Accounts) Session(ctx context.Context, userID, hash string, expiry time.Time) error {
	return r.DB.WithContext(ctx).Exec("INSERT INTO refresh_tokens(user_id, token_hash, expires_at) VALUES (?, ?, ?)", userID, hash, expiry).Error
}
func (r Accounts) Authenticate(ctx context.Context, hash string) (*models.User, error) {
	var u models.User
	err := r.DB.WithContext(ctx).Raw(`SELECT u.* FROM users u JOIN refresh_tokens s ON s.user_id=u.id
 WHERE s.token_hash=? AND NOT s.revoked AND s.expires_at>now() AND u.is_active AND u.deleted_at IS NULL`, hash).Scan(&u).Error
	if err == nil && u.ID == "" {
		err = gorm.ErrRecordNotFound
	}
	return &u, err
}
func (r Accounts) Logout(ctx context.Context, hash string) error {
	return r.DB.WithContext(ctx).Exec("UPDATE refresh_tokens SET revoked=true WHERE token_hash=?", hash).Error
}
func (r Accounts) UpdateProfile(ctx context.Context, id, name, level string) error {
	res := r.DB.WithContext(ctx).Exec("UPDATE users SET first_name=?,cefr_level=?,updated_at=now() WHERE id=? AND EXISTS(SELECT 1 FROM cefr_levels WHERE code=? AND available)", name, level, id, level)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 1 {
		return apperr.NewAppError(400, "El nivel seleccionado no está disponible.")
	}
	return nil
}

func (r Accounts) ChangePassword(ctx context.Context, id string, replace func(string) (string, error)) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var old string
		if err := tx.Raw("SELECT password_hash FROM users WHERE id=? AND is_active AND deleted_at IS NULL FOR UPDATE", id).Scan(&old).Error; err != nil {
			return err
		}
		if old == "" {
			return gorm.ErrRecordNotFound
		}
		hash, err := replace(old)
		if err != nil {
			return err
		}
		if err := tx.Exec("UPDATE users SET password_hash=?,updated_at=now() WHERE id=?", hash, id).Error; err != nil {
			return err
		}
		return tx.Exec("UPDATE refresh_tokens SET revoked=true WHERE user_id=?", id).Error
	})
}
