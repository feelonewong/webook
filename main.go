package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
	"webook/internal/web"
)

func main() {
	hdl := &web.UserHandler{}
	server := gin.Default()
	// resolve cors
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))
	hdl.RegisterRoutes(server)
	server.Run(":8801")
}
