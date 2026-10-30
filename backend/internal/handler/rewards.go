package handler

import "github.com/gin-gonic/gin"

func (a *App) ActivateBooster(c *gin.Context) {
	if !validID(c) {
		return
	}
	out, err := a.Learning.ActivateBooster(c.Request.Context(), currentUser(c).ID, c.Param("id"))
	if err != nil {
		RespondError(c, err)
		return
	}
	WriteJSON(c, 200, out)
}
