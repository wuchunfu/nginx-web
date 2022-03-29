package websiteApi

import (
	"github.com/gin-gonic/gin"
	"github.com/wuchunfu/nginx-web/middleware/logx"
	"github.com/wuchunfu/nginx-web/utils/bytex"
	"github.com/wuchunfu/nginx-web/utils/datetimex"
	"github.com/wuchunfu/nginx-web/utils/filex"
	"github.com/wuchunfu/nginx-web/utils/nginxx"
	"net/http"
	"os"
	"path/filepath"
)

type FileList struct {
	FileName  string `json:"fileName"`  // 文件名
	FileSize  string `json:"fileSize"`  // 文件大小
	IsEnabled bool   `json:"isEnabled"` // 是否已启用
	DateTime  string `json:"dateTime"`  // 文件修改时间
}

func List(ctx *gin.Context) {
	availableConfPath, availablePathErr := nginxx.GetConfPath("sites-available")
	if availablePathErr != nil {
		logx.GetLogger().Sugar().Errorf("Get nginx conf path failed: %s", availablePathErr.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Get nginx conf path failed!",
			"data": nil,
		})
		return
	}
	availableBaseAbsPath, _ := filepath.Abs(availableConfPath)

	availableIsExistPath := filex.FilePathExists(availableBaseAbsPath)
	if !availableIsExistPath {
		logx.GetLogger().Sugar().Errorf("File or directory does not exist!: %s", availableBaseAbsPath)
		ok := filex.MkdirAll(availableBaseAbsPath)
		if !ok {
			logx.GetLogger().Sugar().Errorf("Directory create failed!: %s", availableBaseAbsPath)
			return
		}
	}

	enabledConfPath, enabledPathErr := nginxx.GetConfPath("sites-enabled")
	if enabledPathErr != nil {
		logx.GetLogger().Sugar().Errorf("Get nginx conf path failed: %s", enabledPathErr.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Get nginx conf path failed!",
			"data": nil,
		})
		return
	}
	enabledBaseAbsPath, _ := filepath.Abs(enabledConfPath)

	enabledIsExistPath := filex.FilePathExists(enabledBaseAbsPath)
	if !enabledIsExistPath {
		logx.GetLogger().Sugar().Errorf("File or directory does not exist!: %s", enabledBaseAbsPath)
		ok := filex.MkdirAll(enabledBaseAbsPath)
		if !ok {
			logx.GetLogger().Sugar().Errorf("Directory create failed!: %s", enabledBaseAbsPath)
			return
		}
	}

	availableConfigFileList, availableReadErr := os.ReadDir(availableBaseAbsPath)
	if availableReadErr != nil {
		logx.GetLogger().Sugar().Errorf("Read sites-available path failed: %s", availableReadErr.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Read sites-available path failed!",
			"data": nil,
		})
		return
	}

	enabledConfigFileList, enabledReadErr := os.ReadDir(enabledBaseAbsPath)
	if enabledReadErr != nil {
		logx.GetLogger().Sugar().Errorf("Read sites-enabled path failed: %s", enabledReadErr.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Read sites-enabled path failed!",
			"data": nil,
		})
		return
	}

	enabledConfigMap := make(map[string]bool)
	for i := range enabledConfigFileList {
		enabledConfigMap[enabledConfigFileList[i].Name()] = true
	}

	fileList := make([]FileList, 0)
	for i := range availableConfigFileList {
		file := availableConfigFileList[i]
		fileInfo, _ := file.Info()
		if !fileInfo.IsDir() {
			fileName := fileInfo.Name()
			fileSize := bytex.FormatFileSize(fileInfo.Size())
			dateTime := datetimex.FormatDateTime(fileInfo.ModTime())
			isEnabled := enabledConfigMap[fileName]
			list := FileList{
				FileName:  fileName,
				FileSize:  fileSize,
				IsEnabled: isEnabled,
				DateTime:  dateTime,
			}
			fileList = append(fileList, list)
		}
	}

	// 返回目录json数据
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Get data successfully!",
		"data": fileList,
	})
}
