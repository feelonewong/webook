package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
	"webook/internal/domain"
	"webook/internal/service"
)

const (
	emailRegexPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$`
	// 密码规则： 至少8位、至少包含一个数字和字母
)

type UserHandler struct {
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
	svc              *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRegexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:              svc,
	}
}
func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", h.SignUp)
	ug.POST("/login", h.LoginJwt)
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
	// 规则校验
	isEmail, err := h.emailRegexExp.MatchString(req.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "系统错误",
			"code":    200,
		})
		return
	}
	if !isEmail {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "邮箱格式不正确",
		})
		return
	}
	isPassword, err := h.passwordRegexExp.MatchString(req.Password)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "密码-系统错误",
		})
		return
	}

	if !isPassword {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "密码格式不正确:至少八位，包含一个数组和一个字母",
		})
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "密码输入不一致",
		})
		return
	}

	// service方法调用, 这里不会有ConfirmPassword 前端校验即可
	err = h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "SignUp登录校验成功",
		})
	case service.ErrDuplicateEmail:
		ctx.JSON(http.StatusOK, gin.H{
			"message": "邮箱冲突，请更换",
		})
	default:
		ctx.JSON(http.StatusOK, gin.H{
			"message": "系统错误",
		})
	}

}
func (h *UserHandler) Login(ctx *gin.Context) {
	type loginReq struct {
		Email    string
		Password string
	}

	var req loginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 调用service的业务逻辑
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		sess := sessions.Default(ctx)
		sess.Set("userId", u.Id)
		sess.Options(sessions.Options{
			// 过期时间15分钟
			MaxAge: 900,
		})
		if errorSess := sess.Save(); errorSess != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "系统错误-session",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "登录成功!",
		})
	case service.ErrInvalidUserOrPassword:
		ctx.JSON(http.StatusOK, gin.H{
			"message": "邮箱或者密码不正确",
		})

	default:
		ctx.JSON(http.StatusOK, gin.H{
			"message": "系统错误",
		})
	}

}

func (h *UserHandler) LoginJwt(ctx *gin.Context) {
	type loginReq struct {
		Email    string
		Password string
	}

	var req loginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 调用service的业务逻辑
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		uc := UseClaims{
			Uid: u.Id,
			RegisteredClaims: jwt.RegisteredClaims{
				// 设置过期时间为30分钟
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
		// 生成秘钥
		tokenStr, err := token.SignedString(JWTKey)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": "系统错误",
			})
		}
		ctx.Header("x-jwt-token", tokenStr)
		//sess := sessions.Default(ctx)
		//sess.Set("userId", u.Id)
		//sess.Options(sessions.Options{
		//	// 过期时间15分钟
		//	MaxAge: 900,
		//})
		//if errorSess := sess.Save(); errorSess != nil {
		//	ctx.JSON(http.StatusOK, gin.H{
		//		"message": "系统错误-session",
		//	})
		//	return
		//}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "登录成功!",
		})
	case service.ErrInvalidUserOrPassword:
		ctx.JSON(http.StatusOK, gin.H{
			"message": "邮箱或者密码不正确",
		})

	default:
		ctx.JSON(http.StatusOK, gin.H{
			"message": "系统错误",
		})
	}

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

var JWTKey = []byte("72wNs6ENZnUsmXUC37ro")

type UseClaims struct {
	Uid int64
	jwt.RegisteredClaims
}
