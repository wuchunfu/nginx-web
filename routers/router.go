package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/wuchunfu/nginx-web/api"
	"github.com/wuchunfu/nginx-web/api/configApi"
	"github.com/wuchunfu/nginx-web/api/loginApi"
	"github.com/wuchunfu/nginx-web/api/loginLogApi"
	"github.com/wuchunfu/nginx-web/api/userApi"
	"github.com/wuchunfu/nginx-web/middleware/cors"
	"github.com/wuchunfu/nginx-web/middleware/logx"
)

func InitRouter() *gin.Engine {
	// 全局中间件
	// 路由设置
	router := gin.New()
	router.Use(logx.ZapLogger(), logx.ZapRecovery(true))
	// 设置 Recovery 中间件，主要用于拦截 panic 错误，不至于导致进程崩掉
	router.Use(gin.Recovery())
	// 允许使用跨域请求  全局中间件
	router.Use(cors.Cors())

	gin.SetMode(gin.ReleaseMode)

	router.GET("/ping", api.PingHandler)
	router.GET("/", api.RedirectIndex)
	//router.StaticFS("/ui", web.GetFS())

	loginGroup := router.Group("/sys")
	{
		loginGroup.POST("/login", loginApi.Login)
		loginGroup.POST("/logout", loginApi.Logout)
	}

	router.GET("/sys/login/log/list", loginLogApi.List)

	userGroup := router.Group("/sys/user")
	{
		userGroup.GET("/list", userApi.List)
		userGroup.GET("/detail/:userId", userApi.Detail)
		userGroup.POST("/save", userApi.Save)
		userGroup.PUT("/update", userApi.Update)
		userGroup.POST("/changePassword", userApi.ChangePassword)
		userGroup.POST("/changeLoginPassword", userApi.ChangeLoginPassword)
		userGroup.POST("/delete", userApi.Delete)
	}

	configGroup := router.Group("/sys/config")
	{
		configGroup.GET("/list", configApi.List)
		configGroup.GET("/detail", configApi.Detail)
		configGroup.PUT("/update", configApi.Update)
		configGroup.GET("/changeFolder", configApi.ChangeFolder)
	}

	router.NoRoute(api.NoRouteHandler)
	return router
}
