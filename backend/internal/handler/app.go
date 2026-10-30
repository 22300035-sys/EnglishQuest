package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/englishquest/backend/internal/config"
	"github.com/englishquest/backend/internal/models"
	"github.com/englishquest/backend/internal/repositories"
	"github.com/englishquest/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

type window struct {
	Count int
	Until time.Time
}
type App struct {
	Config   *config.Config
	Auth     services.Auth
	Learning repositories.Learning
	mu       sync.Mutex
	limits   map[string]window
}

func NewApp(cfg *config.Config, db *gorm.DB) *App {
	ttl := cfg.JWTRefreshTTL
	if ttl <= 0 {
		ttl = 72 * time.Hour
	}
	return &App{Config: cfg, Auth: services.Auth{Accounts: repositories.Accounts{DB: db}, TTL: ttl}, Learning: repositories.Learning{DB: db}, limits: map[string]window{}}
}
func (a *App) cookieName() string {
	if a.Config.Env == "production" {
		return "__Host-eq_session"
	}
	return "eq_session"
}
func (a *App) cookie(c *gin.Context, token string, age int) {
	sameSite := http.SameSiteLaxMode
	if a.Config.Env == "production" {
		sameSite = http.SameSiteNoneMode
	}
	http.SetCookie(c.Writer, &http.Cookie{Name: a.cookieName(), Value: token, Path: "/", HttpOnly: true, Secure: a.Config.Env == "production", SameSite: sameSite, Partitioned: a.Config.Env == "production", MaxAge: age})
}
func decode(c *gin.Context, out interface{}) bool {
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(out); err != nil {
		WriteError(c, 400, "Solicitud JSON inválida.")
		return false
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		WriteError(c, 400, "Envía un único documento JSON.")
		return false
	}
	return true
}

var uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func validID(c *gin.Context) bool {
	if !uuidPattern.MatchString(c.Param("id")) {
		WriteError(c, 400, "Identificador inválido.")
		return false
	}
	return true
}
func currentUser(c *gin.Context) *models.User { return c.MustGet("user").(*models.User) }

// Every state-changing request, including login, needs an exact trusted Origin.
// CORS alone is not a CSRF defense. Non-browser clients must also send Origin.
func (a *App) OriginCheck(c *gin.Context) {
	if c.Request.Method != "GET" && c.Request.Method != "HEAD" && c.Request.Method != "OPTIONS" {
		allowed := false
		for _, origin := range a.Config.AllowedOrigins {
			if c.GetHeader("Origin") == origin {
				allowed = true
			}
		}
		if !allowed {
			c.AbortWithStatusJSON(403, gin.H{"error": "Origen no autorizado."})
			return
		}
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 12*time.Second)
	defer cancel()
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}
func (a *App) RateLimit(c *gin.Context) {
	a.mu.Lock()
	now := time.Now()
	key := c.ClientIP()
	w := a.limits[key]
	if now.After(w.Until) {
		w = window{Until: now.Add(15 * time.Minute)}
	}
	if len(a.limits) >= 5000 {
		for k, v := range a.limits {
			if now.After(v.Until) {
				delete(a.limits, k)
			}
		}
	}
	if w.Count >= 15 || (len(a.limits) >= 5000 && w.Count == 0) {
		a.mu.Unlock()
		c.Header("Retry-After", "900")
		c.AbortWithStatusJSON(429, gin.H{"error": "Demasiados intentos. Espera 15 minutos."})
		return
	}
	w.Count++
	a.limits[key] = w
	a.mu.Unlock()
	c.Next()
}
func (a *App) RequireSession(c *gin.Context) {
	raw, err := c.Cookie(a.cookieName())
	if err != nil || len(raw) != 64 {
		c.AbortWithStatusJSON(401, gin.H{"error": "Inicia sesión para continuar."})
		return
	}
	u, err := a.Auth.Accounts.Authenticate(c.Request.Context(), services.HashSession(raw))
	if err != nil {
		if c.Request.Context().Err() != nil {
			RespondError(c, err)
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(401, gin.H{"error": "La sesión expiró. Inicia sesión de nuevo."})
		} else {
			c.AbortWithStatusJSON(503, gin.H{"error": "No se pudo verificar la sesión."})
		}
		return
	}
	c.Set("user", u)
	c.Next()
}
func (a *App) RequireAdmin(c *gin.Context) {
	if currentUser(c).Role != "ADMIN" {
		c.AbortWithStatusJSON(403, gin.H{"error": "Acceso exclusivo de administración."})
		return
	}
	c.Next()
}
func (a *App) Register(c *gin.Context) {
	var in services.Registration
	if !decode(c, &in) {
		return
	}
	u, err := a.Auth.Register(c.Request.Context(), in)
	if err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 201, u)
}
func (a *App) Login(c *gin.Context) {
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if !decode(c, &in) {
		return
	}
	if len(in.Email) > 254 || len(in.Password) > 72 {
		WriteError(c, 401, "Correo o contraseña incorrectos.")
		return
	}
	u, raw, err := a.Auth.Login(c.Request.Context(), in.Email, in.Password)
	if err != nil {
		RespondError(c, err)
		return
	}
	a.cookie(c, raw, int(a.Auth.TTL.Seconds()))
	WriteJSON(c, 200, u)
}
func (a *App) Logout(c *gin.Context) {
	raw, _ := c.Cookie(a.cookieName())
	if err := a.Auth.Accounts.Logout(c.Request.Context(), services.HashSession(raw)); err != nil {
		RespondError(c, err)
		return
	}
	a.cookie(c, "", -1)
	WriteJSON(c, 200, gin.H{"logged_out": true})
}
func (a *App) Profile(c *gin.Context) { WriteJSON(c, 200, currentUser(c)) }
func (a *App) UpdateProfile(c *gin.Context) {
	var in struct {
		Name  string `json:"name"`
		Level string `json:"cefr_level"`
	}
	if !decode(c, &in) {
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Level = strings.ToUpper(strings.TrimSpace(in.Level))
	if len([]rune(in.Name)) < 1 || len([]rune(in.Name)) > 80 || len(in.Level) != 2 {
		WriteError(c, 400, "Nombre inválido o nivel no disponible.")
		return
	}
	if err := a.Auth.Accounts.UpdateProfile(c.Request.Context(), currentUser(c).ID, in.Name, in.Level); err != nil {
		RespondError(c, err)
		return
	}
	u := *currentUser(c)
	u.FirstName = in.Name
	u.CEFRLevel = in.Level
	WriteJSON(c, 200, &u)
}
func (a *App) Levels(c *gin.Context) {
	out, err := a.Learning.Levels(c.Request.Context())
	if err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 200, out)
}
func (a *App) Lessons(c *gin.Context) {
	out, err := a.Learning.Lessons(c.Request.Context(), currentUser(c).ID)
	if err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 200, out)
}
func (a *App) Dashboard(c *gin.Context) {
	out, err := a.Learning.Dashboard(c.Request.Context(), currentUser(c))
	if err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 200, out)
}
func (a *App) Start(c *gin.Context) {
	if !validID(c) {
		return
	}
	out, err := a.Learning.Start(c.Request.Context(), currentUser(c).ID, c.Param("id"))
	if err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 201, out)
}
func (a *App) Attempt(c *gin.Context) {
	if !validID(c) {
		return
	}
	out, err := a.Learning.Attempt(c.Request.Context(), currentUser(c).ID, c.Param("id"))
	if err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 200, out)
}
func (a *App) Submit(c *gin.Context) {
	if !validID(c) {
		return
	}
	var in models.Submission
	if !decode(c, &in) {
		return
	}
	out, err := a.Learning.Submit(c.Request.Context(), currentUser(c).ID, c.Param("id"), in, services.Grade, services.NextStreak)
	if err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 200, out)
}
func (a *App) AdminOverview(c *gin.Context) {
	out, err := a.Learning.AdminOverview(c.Request.Context())
	if err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 200, out)
}
func (a *App) AdminActivity(c *gin.Context) {
	if !validID(c) {
		return
	}
	out, err := a.Learning.AdminActivity(c.Request.Context(), c.Param("id"))
	if err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 200, out)
}
func (a *App) SaveContent(c *gin.Context) {
	id := c.Param("id")
	if id != "" && !validID(c) {
		return
	}
	var in repositories.ContentEdit
	if !decode(c, &in) {
		return
	}
	if len(strings.TrimSpace(in.LessonTitle)) < 1 || len(in.LessonTitle) > 150 || len(in.Objective) < 5 || len(in.Objective) > 1000 || len(in.Content) < 10 || len(in.Content) > 10000 || (id == "" && !uuidPattern.MatchString(in.UnitID)) {
		WriteError(c, 400, "Revisa los datos de la lección.")
		return
	}
	if err := services.ValidateActivity(in.Activity); err != nil {
		RespondError(c, err)
		return
	}
	saved, err := a.Learning.SaveContent(c.Request.Context(), currentUser(c).ID, id, in)
	if err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 200, gin.H{"id": saved})
}
func (a *App) SetActive(c *gin.Context) {
	if !validID(c) {
		return
	}
	var in struct {
		Active *bool `json:"active"`
	}
	if !decode(c, &in) {
		return
	}
	if in.Active == nil {
		WriteError(c, 400, "Falta el estado de la cuenta.")
		return
	}
	if err := a.Learning.SetActive(c.Request.Context(), currentUser(c).ID, c.Param("id"), *in.Active); err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 200, gin.H{"updated": true})
}

func (a *App) Routes(v1 *gin.RouterGroup) {
	app := v1.Group("")
	app.Use(a.OriginCheck)
	auth := app.Group("/auth")
	auth.POST("/register", a.RateLimit, a.Register)
	auth.POST("/login", a.RateLimit, a.Login)
	app.GET("/levels", a.Levels)
	secure := app.Group("")
	secure.Use(a.RequireSession)
	secure.POST("/auth/logout", a.Logout)
	secure.GET("/auth/me", a.Profile)
	secure.PATCH("/profile", a.UpdateProfile)
	secure.POST("/profile/password", a.ChangePassword)
	secure.GET("/dashboard", a.Dashboard)
	secure.POST("/boosters/:id/activate", a.ActivateBooster)
	secure.GET("/lessons", a.Lessons)
	secure.GET("/learning/levels", a.LevelAccess)
	secure.POST("/levels/:code/access-exams", a.StartAccessExam)
	secure.GET("/access-exams/:id", a.AccessExam)
	secure.POST("/access-exams/:id/submit", a.SubmitAccessExam)
	secure.GET("/diploma", a.DiplomaStatus)
	secure.POST("/diploma", a.IssueDiploma)
	secure.GET("/units/:id/project", a.UnitProject)
	secure.PUT("/units/:id/project", a.SaveUnitProject)
	secure.POST("/activities/:id/start", a.Start)
	secure.GET("/attempts/:id", a.Attempt)
	secure.POST("/attempts/:id/submit", a.Submit)
	admin := secure.Group("/admin")
	admin.Use(a.RequireAdmin)
	admin.GET("/overview", a.AdminOverview)
	admin.PUT("/levels/:code", a.SaveLevel)
	admin.POST("/units", a.SaveUnit)
	admin.PUT("/units/:id", a.SaveUnit)
	admin.PUT("/resources", a.SaveResource)
	admin.GET("/activities/:id", a.AdminActivity)
	admin.POST("/lessons", a.SaveContent)
	admin.PUT("/activities/:id", a.SaveContent)
	admin.PATCH("/users/:id", a.SetActive)
}
