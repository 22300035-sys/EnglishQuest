package models

import "time"

import "gorm.io/gorm"

// User represents an EnglishQuest student or administrator.
// The role field stores the application role: STUDENT or ADMIN.
type User struct {
	ID           string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email        string         `gorm:"type:varchar(255);not null;uniqueIndex;size:255" json:"email"`
	PasswordHash string         `gorm:"type:varchar(255);not null;size:255" json:"-"`
	FirstName    string         `gorm:"type:varchar(100)" json:"first_name"`
	LastName     string         `gorm:"type:varchar(100)" json:"last_name"`
	CEFRLevel    string         `gorm:"type:varchar(4)" json:"cefr_level"`
	Role         string         `gorm:"type:varchar(20)" json:"role"`
	Timezone     string         `gorm:"type:varchar(100);default:UTC" json:"timezone"`
	IsActive     bool           `gorm:"not null" json:"is_active"`
	CreatedAt    time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-"`
}

// PublicUser is the safe user representation returned by the API.
type PublicUser struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CEFRLevel string    `json:"cefr_level"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToPublicUser(u *User) *PublicUser {
	if u == nil {
		return nil
	}
	return &PublicUser{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CEFRLevel: u.CEFRLevel,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
