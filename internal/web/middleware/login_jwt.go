package middleware

import (
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
	"webook/internal/web"
)

type LoginJWTMiddlewareBuilder struct {
}

func (m *LoginJWTMiddlewareBuilder) CheckLoginJwt() gin.HandlerFunc {
	// 注册time类型
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/login" || path == "/users/signup" {
			return
		}
		// token 存放在 Authorization
		// format: Bearer Xxxx
		authCode := ctx.GetHeader("Authorization")
		if authCode == "" {
			// empty token
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(authCode, " ")
		// 切割成一个长度为2的数组
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		var uc web.UseClaims
		// 校验token
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return web.JWTKey, nil
		})
		if err != nil {
			// 有token但是token不正确
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		if token == nil || !token.Valid {
			// token 非法的或者过期了
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		expireTime := uc.ExpiresAt
		if expireTime.Sub(time.Now()) < time.Second*50 {
			// 小于50s刷新token
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString(web.JWTKey)
			ctx.Header("x-jwt-token", tokenStr)
			if err != nil {
				// 过期时间没有刷新，用户登录不了
				log.Println(err)
			}
		}
		ctx.Set("user", uc)
	}
}
