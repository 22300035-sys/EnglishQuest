package handler

import "github.com/gin-gonic/gin"

func (a *App) UnitProject(c *gin.Context) {
	if !validID(c) {
		return
	}
	out, err := a.Learning.UnitProject(c.Request.Context(), currentUser(c).ID, c.Param("id"))
	if err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 200, out)
}
func (a *App) SaveUnitProject(c *gin.Context) {
	if !validID(c) {
		return
	}
	var in struct {
		Body       string `json:"body"`
		Reflection string `json:"reflection"`
	}
	if !decode(c, &in) {
		return
	}
	if err := a.Learning.SaveUnitProject(c.Request.Context(), currentUser(c).ID, c.Param("id"), in.Body, in.Reflection); err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 200, gin.H{"saved": true})
}
