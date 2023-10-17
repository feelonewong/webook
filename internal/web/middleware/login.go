package middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
}

func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	// 注册time类型
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/login" || path == "/users/signup" {
			return
		}
		sess := sessions.Default(ctx)
		userId := sess.Get("userId")
		if userId == nil {
			// 缺少userId中断后续操作
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		const updateTimeKey = "update_time"
		now := time.Now()
		// 拿上一次刷新时间
		val := sess.Get(updateTimeKey)
		lastUpdateTime, ok := val.(time.Time)

		// 1.第一次设置进来session的userId
		// 2.不ok的时候重新设置cookie
		// 3.时间大于1分钟设置session的userId
		if val == nil || !ok || now.Sub(lastUpdateTime) > time.Second*10 {
			// 第一次进来存储数据
			sess.Set(updateTimeKey, now)
			sess.Set("userId", userId)
			err := sess.Save()
			if err != nil {
				fmt.Println("session存储异常:", err)
			}
		}

	}
}
