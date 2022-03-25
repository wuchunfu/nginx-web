package loginLogApi

import (
	"github.com/gin-gonic/gin"
	"github.com/wuchunfu/nginx-web/common"
	"github.com/wuchunfu/nginx-web/model/loginLogModel"
	"net/http"
)

func List(ctx *gin.Context) {
	username := ctx.Query("username")

	loginLog := new(loginLogModel.LoginLog)
	page, pageSize := common.ParseQueryParams(ctx)
	dataList, count := loginLog.List(page, pageSize, username)
	ctx.JSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"data":     &dataList,
		"msg":      "获取数据成功！",
		"page":     page,
		"pageSize": pageSize,
		"total":    count,
	})
}
