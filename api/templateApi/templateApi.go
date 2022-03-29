package templateApi

import (
	"github.com/gin-gonic/gin"
	"github.com/wuchunfu/nginx-web/middleware/logx"
	"github.com/wuchunfu/nginx-web/template"
	"net/http"
	"strings"
)

func GetTemplate(ctx *gin.Context) {
	fileName := ctx.Param("fileName")
	contentByte, err := template.TemplateFS.ReadFile(fileName)
	if err != nil {
		logx.GetLogger().Sugar().Errorf("Get template file failed: %s", err.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Get template file failed!",
			"data": nil,
		})
		return
	}

	content := string(contentByte)
	content = strings.ReplaceAll(content, "{{ HTTP01PORT }}", "9180")

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Get data successfully!",
		"data": content,
	})
}
