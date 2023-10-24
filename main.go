package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
	"webook/internal/repository"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/middleware"
)

func main() {
	// 初始化数据库
	db := initDB()

	// 初始化server服务
	server := initWebServer()

	// 初始化User
	initUser(db, server)

	server.Run(":8801")
}

func initDB() *gorm.DB {
	dsn := "root:root@tcp(127.0.0.1:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	user := dao.User{}
	err = dao.InitTables(db, &user)
	if err != nil {
		panic(err)
	}
	// 打印数据库日志
	db = db.Debug()
	return db
}
func initWebServer() *gin.Engine {
	server := gin.Default()
	// resolve cors
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		// 允许前端添加的请求头
		ExposeHeaders: []string{"x-jwt-token"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	// 限流阈值

	useJWT(server)
	return server
}

func useSession(server *gin.Engine) {
	login := &middleware.LoginMiddlewareBuilder{}
	store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
		[]byte("yLbakqA10vl62ADPax5ScvE69B0Ph43W"),
		[]byte("bqD05B9Ze6UDkwX2OSk5AA1sFp19KFxO"),
	)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
}

func useJWT(server *gin.Engine) {
	login := &middleware.LoginJWTMiddlewareBuilder{}
	server.Use(login.CheckLoginJwt())
}

func redisSaveCookie(server *gin.Engine) {
	// 基于redis存储cookie

}

func initUser(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	hdl := web.NewUserHandler(us)
	hdl.RegisterRoutes(server)
}
