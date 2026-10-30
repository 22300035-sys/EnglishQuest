package handler

import (
	"github.com/englishquest/backend/internal/models"
	"github.com/englishquest/backend/internal/services"
	"github.com/gin-gonic/gin"
	"strings"
)

func (a *App) SaveLevel(c *gin.Context) {
	var in models.LevelEdit
	if !decode(c, &in) {
		return
	}
	in.Title = strings.TrimSpace(in.Title)
	if in.Title == "" || len(in.Title) > 80 || len(in.Description) > 2000 {
		WriteError(c, 400, "Revisa el título y la descripción del nivel.")
		return
	}
	if err := a.Learning.SaveLevel(c.Request.Context(), currentUser(c).ID, c.Param("code"), in); err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 200, gin.H{"updated": true})
}
func (a *App) SaveUnit(c *gin.Context) {
	if c.Param("id") != "" && !validID(c) {
		return
	}
	var in models.UnitEdit
	if !decode(c, &in) {
		return
	}
	in.Title = strings.TrimSpace(in.Title)
	if in.Title == "" || len(in.Title) > 150 || len(in.Description) > 2000 {
		WriteError(c, 400, "Revisa el título y la descripción de la unidad.")
		return
	}
	id, err := a.Learning.SaveUnit(c.Request.Context(), currentUser(c).ID, c.Param("id"), in)
	if err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 200, gin.H{"id": id})
}
func (a *App) SaveResource(c *gin.Context) {
	var in models.Resource
	if !decode(c, &in) {
		return
	}
	if err := services.ValidateResource(in); err != nil {
		RespondError(c, err)
		return
	}
	if err := a.Learning.SaveResource(c.Request.Context(), currentUser(c).ID, in); err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 200, gin.H{"updated": true})
}
func (a *App) ChangePassword(c *gin.Context) {
	var in struct {
		Current string `json:"current_password"`
		New     string `json:"new_password"`
	}
	if !decode(c, &in) {
		return
	}
	if err := a.Auth.ChangePassword(c.Request.Context(), currentUser(c).ID, in.Current, in.New); err != nil {
		RespondError(c, err)
		return
	}
	a.cookie(c, "", -1)
	WriteJSON(c, 200, gin.H{"changed": true, "sessions_revoked": true})
}
