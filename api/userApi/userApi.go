package userApi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wuchunfu/nginx-web/common"
	"github.com/wuchunfu/nginx-web/model/userModel"
	"github.com/wuchunfu/nginx-web/utils"
	"github.com/wuchunfu/nginx-web/utils/datetimex"
	"net/http"
	"strconv"
	"strings"
)

// UserForm 用户表单
type UserForm struct {
	UserId      int64
	UserIds     []int64
	Username    string // 用户名
	Password    string // 密码
	NewPassword string // 确认密码
	Email       string // 邮箱
	IsAdmin     int8   // 是否是管理员 1:管理员 0:普通用户
	Status      int8   // 启用状态 1:启动 0:禁用
	CreateTime  string // 创建时间
	UpdateTime  string // 修改时间
}

// List 用户列表页
func List(ctx *gin.Context) {
	page, pageSize := common.ParseQueryParams(ctx)
	username := ctx.Query("username")
	user := new(userModel.User)
	dataList, count := user.List(page, pageSize, username)
	ctx.JSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"data":     &dataList,
		"msg":      "获取数据成功！",
		"page":     page,
		"pageSize": pageSize,
		"total":    count,
	})
}

// Detail 用户详情
func Detail(ctx *gin.Context) {
	userId, _ := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	user := new(userModel.User)
	detail := user.Detail(userId)
	if detail.UserId >= 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": &detail,
			"msg":  "获取数据成功！",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": &detail,
			"msg":  "获取数据失败！",
		})
	}
}

func Save(ctx *gin.Context) {
	form := new(UserForm)
	ctx.Bind(form)
	username := strings.TrimSpace(form.Username)
	password := strings.TrimSpace(form.Password)
	email := strings.TrimSpace(form.Email)
	isAdmin := form.IsAdmin
	status := form.Status

	user := new(userModel.User)
	usernameCount := user.IsExistsUsername(username)
	emailCount := user.IsExistsEmail(email)
	if usernameCount > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": nil,
			"msg":  "用户名已存在！",
		})
	} else if emailCount > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": nil,
			"msg":  "邮箱已存在！",
		})
	} else {
		user.Username = username
		user.Salt = utils.RandomString(6)
		user.Password = utils.Md5(password + user.Salt)
		user.Email = email
		user.IsAdmin = isAdmin
		user.Status = status
		user.CreateTime = datetimex.FormatNowDateTime()
		user.UpdateTime = datetimex.FormatNowDateTime()

		user.Save()

		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": nil,
			"msg":  "保存成功！",
		})
	}
}

func Update(ctx *gin.Context) {
	form := new(UserForm)
	ctx.Bind(form)
	userId := form.UserId
	username := strings.TrimSpace(form.Username)
	email := strings.TrimSpace(form.Email)
	isAdmin := form.IsAdmin
	status := form.Status

	user := new(userModel.User)
	userMap := make(map[string]interface{})
	userMap["username"] = username
	userMap["email"] = email
	userMap["is_admin"] = isAdmin
	userMap["status"] = status
	userMap["update_time"] = datetimex.FormatNowDateTime()
	user.Update(userId, userMap)

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": nil,
		"msg":  "修改成功！",
	})
}

func ChangePassword(ctx *gin.Context) {
	form := new(UserForm)
	ctx.Bind(form)
	userId := form.UserId
	username := strings.TrimSpace(form.Username)
	password := strings.TrimSpace(form.Password)

	if password == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": nil,
			"msg":  "请输入密码！",
		})
	} else {
		//salt := utils.RandomString(6)
		user := new(userModel.User)
		userMap := make(map[string]interface{})
		userMap["username"] = username
		userMap["password"] = utils.Md5(password + user.Salt)
		userMap["update_time"] = datetimex.FormatNowDateTime()

		user.ChangePassword(userId, userMap)

		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": nil,
			"msg":  "密码修改成功！",
		})
	}
}

func ChangeLoginPassword(ctx *gin.Context) {
	form := new(UserForm)
	ctx.Bind(form)
	userId := form.UserId
	username := strings.TrimSpace(form.Username)
	oldPassword := strings.TrimSpace(form.Password)
	newPassword := strings.TrimSpace(form.NewPassword)

	user := new(userModel.User)
	detail := user.Detail(userId)
	password := utils.Md5(oldPassword + user.Salt)
	if detail.Password != password {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": nil,
			"msg":  "原密码不正确！",
		})
	} else if oldPassword == "" || newPassword == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": nil,
			"msg":  "请输入密码！",
		})
	} else {
		userMap := make(map[string]interface{})
		userMap["username"] = username
		userMap["password"] = utils.Md5(newPassword + user.Salt)
		userMap["update_time"] = datetimex.FormatNowDateTime()

		user.ChangePassword(userId, userMap)

		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": nil,
			"msg":  "密码修改成功！",
		})
	}
}

func Delete(ctx *gin.Context) {
	form := new(UserForm)
	ctx.Bind(form)

	fmt.Println("=-=-=-=")
	fmt.Println(form.UserIds)

	user := new(userModel.User)
	for _, id := range form.UserIds {
		user.Delete(id)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": nil,
		"msg":  "删除成功！",
	})
}
