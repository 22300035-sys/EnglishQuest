package handler

import (
	"github.com/englishquest/backend/internal/models"
	"github.com/englishquest/backend/internal/services"
	"github.com/gin-gonic/gin"
)

func (a *App) LevelAccess(c *gin.Context) {
	out, err := a.Learning.LevelAccess(c.Request.Context(), currentUser(c).ID)
	if err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 200, out)
}
func (a *App) StartAccessExam(c *gin.Context) {
	out, err := a.Learning.StartAccessExam(c.Request.Context(), currentUser(c).ID, c.Param("code"), services.BuildAccessExam)
	if err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 201, out)
}
func (a *App) AccessExam(c *gin.Context) {
	if !validID(c) {
		return
	}
	out, err := a.Learning.AccessExam(c.Request.Context(), currentUser(c).ID, c.Param("id"))
	if err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 200, out)
}
func (a *App) SubmitAccessExam(c *gin.Context) {
	if !validID(c) {
		return
	}
	var in models.Submission
	if !decode(c, &in) {
		return
	}
	out, err := a.Learning.SubmitAccessExam(c.Request.Context(), currentUser(c).ID, c.Param("id"), in, services.Grade)
	if err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 200, out)
}
func (a *App) DiplomaStatus(c *gin.Context) {
	out, err := a.Learning.DiplomaStatus(c.Request.Context(), currentUser(c).ID)
	if err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 200, out)
}
func (a *App) IssueDiploma(c *gin.Context) {
	var in struct {
		Name string `json:"name"`
	}
	if !decode(c, &in) {
		return
	}
	name, err := services.DiplomaName(in.Name)
	if err != nil {
		RespondError(c, err)
		return
	}
	out, err := a.Learning.IssueDiploma(c.Request.Context(), currentUser(c).ID, name)
	if err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 200, out)
}
