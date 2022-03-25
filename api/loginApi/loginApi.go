package loginApi

import (
	"github.com/gin-gonic/gin"
	"github.com/wuchunfu/nginx-web/middleware/jwtx"
	"github.com/wuchunfu/nginx-web/model/loginLogModel"
	"github.com/wuchunfu/nginx-web/model/tokenModel"
	"github.com/wuchunfu/nginx-web/model/userModel"
	"github.com/wuchunfu/nginx-web/utils"
	"github.com/wuchunfu/nginx-web/utils/datetimex"
	"net/http"
	"strings"
)

// LoginForm 登录表单
type LoginForm struct {
	Username string // 用户名
	Password string // 密码
}

func Login(ctx *gin.Context) {
	form := new(LoginForm)
	ctx.Bind(form)

	username := strings.TrimSpace(form.Username)
	password := strings.TrimSpace(form.Password)

	user := new(userModel.User)
	usernameCount := user.IsExistsUsername(username)
	if usernameCount <= 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"data": nil,
			"msg":  "用户名不存在！",
		})
		return
	}

	info := user.GetInfoByName(username)
	pwd := utils.Md5(password + user.Salt)
	if info.Password != pwd {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"data": nil,
			"msg":  "用户名或密码不正确！",
		})
		return
	}

	tokenInfo, genErr := jwtx.GenerateToken(user.UserId, user.Username)
	if genErr != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusUnauthorized,
			"data": nil,
			"msg":  "Token 生成失败！",
		})
		return
	}

	// 保存 token
	tokens := new(tokenModel.UserToken)
	tokens.Save(user.UserId, tokenInfo)

	loginLog := new(loginLogModel.LoginLog)
	loginLog.Username = user.Username
	loginLog.Ip = ctx.ClientIP()
	loginLog.CreateTime = datetimex.FormatNowDateTime()
	loginLog.Save()

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": tokenInfo,
		"msg":  "Token 获取成功！",
	})
}

func Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusUnauthorized,
			"data": nil,
			"msg":  "权限不足！",
		})
		return
	}

	parseToken, err := jwtx.ParseToken(token)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusUnauthorized,
			"data": nil,
			"msg":  "权限不足！",
		})
		return
	}

	userId := parseToken.UserId
	username := parseToken.Username
	// 重新生成 token
	tokenInfo, genErr := jwtx.GenerateToken(userId, username)
	if genErr != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusUnauthorized,
			"data": nil,
			"msg":  "Token 生成失败！",
		})
		return
	}

	tokens := new(tokenModel.UserToken)
	tokens.Update(userId, tokenInfo)

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": tokenInfo,
		"msg":  "Token 刷新成功！",
	})
}
