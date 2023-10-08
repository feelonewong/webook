package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
}

func (a *AuthHandler) RegisterRoutes(server *gin.Engine) {
	ag := server.Group("/auths")
	ag.POST("/signout", a.SignOut)
	ag.POST("/back", a.Back)
}

func (a *AuthHandler) SignOut(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "SignOut",
	})
}

func (a *AuthHandler) Back(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Back",
	})
}
