package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginMiddlewareBuilder struct {
}

func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/login" || path == "/users/signup" {
			return
		}
		sess := sessions.Default(ctx)
		if sess.Get("userId") == nil {
			// 缺少userId中断后续操作
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
