package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", h.SignUp)
	ug.POST("/login", h.Login)
	ug.GET("/profile", h.Profile)
	ug.POST("/edit", h.Edit)
}
func (h *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	// bind最常用的用于接受请求的方法
	// 用于解析传递过来的参数类型 一旦格式不对 返回报错.
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "SignUp",
	})
}
func (h *UserHandler) Login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login",
	})
}
func (h *UserHandler) Profile(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Profile",
	})
}
func (h *UserHandler) Edit(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Edit",
	})
}
