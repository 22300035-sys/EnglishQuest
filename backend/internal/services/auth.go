package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/englishquest/backend/internal/models"
	apperr "github.com/englishquest/backend/internal/pkg/errors"
	"github.com/englishquest/backend/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/mail"
	"strings"
	"time"
	_ "time/tzdata"
)

type Registration struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Timezone string `json:"timezone"`
}
type Auth struct {
	Accounts repositories.Accounts
	TTL      time.Duration
}

var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("invalid-account-timing-placeholder"), bcrypt.DefaultCost)

func HashSession(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
func (s Auth) Register(ctx context.Context, in Registration) (*models.User, error) {
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	in.Name = strings.TrimSpace(in.Name)
	address, err := mail.ParseAddress(in.Email)
	if err != nil || address.Address != in.Email || len(in.Email) > 254 {
		return nil, apperr.NewAppError(400, "Introduce un correo válido.")
	}
	if len(in.Password) < 10 || len(in.Password) > 72 {
		return nil, apperr.NewAppError(400, "La contraseña debe tener entre 10 y 72 bytes.")
	}
	if len([]rune(in.Name)) < 1 || len([]rune(in.Name)) > 80 {
		return nil, apperr.NewAppError(400, "El nombre debe tener entre 1 y 80 caracteres.")
	}
	if in.Timezone == "" {
		in.Timezone = "UTC"
	}
	if in.Timezone == "Local" {
		return nil, apperr.NewAppError(400, "Selecciona una zona horaria IANA explícita.")
	}
	if len(in.Timezone) > 100 {
		return nil, apperr.NewAppError(400, "Zona horaria inválida.")
	}
	if _, err := time.LoadLocation(in.Timezone); err != nil {
		return nil, apperr.NewAppError(400, "Zona horaria inválida.")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &models.User{Email: in.Email, PasswordHash: string(hash), FirstName: in.Name, Timezone: in.Timezone, CEFRLevel: "A1", Role: "STUDENT", IsActive: true}
	if err := s.Accounts.Create(ctx, u); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperr.NewAppError(409, "No se pudo registrar ese correo; prueba iniciar sesión.")
		}
		return nil, err
	}
	return u, nil
}
func (s Auth) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	u, err := s.Accounts.ByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", err
	}
	hash := dummyHash
	if err == nil {
		hash = []byte(u.PasswordHash)
	}
	valid := bcrypt.CompareHashAndPassword(hash, []byte(password)) == nil
	if err != nil || !valid || !u.IsActive {
		return nil, "", apperr.NewAppError(401, "Correo o contraseña incorrectos.")
	}
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return nil, "", err
	}
	raw := hex.EncodeToString(bytes)
	if err := s.Accounts.Session(ctx, u.ID, HashSession(raw), time.Now().Add(s.TTL)); err != nil {
		return nil, "", err
	}
	return u, raw, nil
}

func (s Auth) ChangePassword(ctx context.Context, id, current, next string) error {
	if len(next) < 10 || len(next) > 72 || len(current) > 72 {
		return apperr.NewAppError(400, "La nueva contraseña debe tener entre 10 y 72 bytes.")
	}
	if current == next {
		return apperr.NewAppError(400, "Elige una contraseña diferente.")
	}
	return s.Accounts.ChangePassword(ctx, id, func(hash string) (string, error) {
		if bcrypt.CompareHashAndPassword([]byte(hash), []byte(current)) != nil {
			return "", apperr.NewAppError(403, "La contraseña actual no es correcta.")
		}
		encoded, err := bcrypt.GenerateFromPassword([]byte(next), bcrypt.DefaultCost)
		return string(encoded), err
	})
}
