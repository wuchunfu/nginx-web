package websiteApi

import (
	"github.com/gin-gonic/gin"
	"github.com/wuchunfu/nginx-web/middleware/logx"
	"github.com/wuchunfu/nginx-web/model/websiteModel"
	"github.com/wuchunfu/nginx-web/utils/bytex"
	"github.com/wuchunfu/nginx-web/utils/datetimex"
	"github.com/wuchunfu/nginx-web/utils/filex"
	"github.com/wuchunfu/nginx-web/utils/nginxx"
	"github.com/wuchunfu/nginx-web/utils/templatex"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// WebsiteForm 网站表单
type WebsiteForm struct {
	WebsiteId         int64  // 配置文件 id
	FileName          string // 配置文件名称
	ServerName        string // 网站域名
	RootDirectory     string // 网站根目录
	HomePage          string // 网站主页
	HttpPort          int16  // http 监听端口
	SupportSsl        int8   // 启用 TLS 状态 1: 启用 0: 禁用
	HttpsPort         int16  // https 监听端口
	SslCertificate    string // TLS 证书路径
	SslCertificateKey string // 私钥路径
	Status            int8   // 启用状态 1:启用 0:禁用
	CreateTime        string // 创建时间
	UpdateTime        string // 修改时间
}

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
			"code": http.StatusBadRequest,
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
			"code": http.StatusBadRequest,
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
			"code": http.StatusBadRequest,
			"msg":  "Read sites-available path failed!",
			"data": nil,
		})
		return
	}

	enabledConfigFileList, enabledReadErr := os.ReadDir(enabledBaseAbsPath)
	if enabledReadErr != nil {
		logx.GetLogger().Sugar().Errorf("Read sites-enabled path failed: %s", enabledReadErr.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
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

func Save(ctx *gin.Context) {
	form := new(WebsiteForm)
	ctx.Bind(form)

	fileName := strings.TrimSpace(form.FileName)
	serverName := strings.TrimSpace(form.ServerName)
	rootDirectory := strings.TrimSpace(form.RootDirectory)
	homePage := strings.TrimSpace(form.HomePage)
	httpPort := form.HttpPort
	supportSsl := form.SupportSsl
	httpsPort := form.HttpsPort
	sslCertificate := strings.TrimSpace(form.SslCertificate)
	sslCertificateKey := strings.TrimSpace(form.SslCertificateKey)
	status := form.Status

	website := new(websiteModel.Website)
	fileNameCount := website.IsExistsFileName(fileName)
	if fileNameCount > 0 {
		logx.GetLogger().Sugar().Warn("Config name already exists!")
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": nil,
			"msg":  "Config name already exists!",
		})
		return
	}

	website.FileName = fileName
	website.ServerName = serverName
	website.RootDirectory = rootDirectory
	website.HomePage = homePage
	website.HttpPort = httpPort
	if supportSsl == 1 {
		website.SupportSsl = supportSsl
		website.HttpsPort = httpsPort
		website.SslCertificate = sslCertificate
		website.SslCertificateKey = sslCertificateKey
	}
	website.Status = status
	website.CreateTime = datetimex.FormatNowDateTime()
	website.UpdateTime = datetimex.FormatNowDateTime()

	availableConfPath, availablePathErr := nginxx.GetConfPath("sites-available")
	if availablePathErr != nil {
		logx.GetLogger().Sugar().Errorf("Get nginx conf path failed: %s", availablePathErr.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
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
			"code": http.StatusBadRequest,
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

	availableFileAbsPath := filepath.Join(availableBaseAbsPath, fileName)
	availableFileIsExists := filex.FilePathExists(availableFileAbsPath)
	if availableFileIsExists {
		logx.GetLogger().Sugar().Errorf("File already exists!: %s", availableFileAbsPath)
		return
	}

	content := templatex.TemplateReplace(website)
	err := os.WriteFile(availableFileAbsPath, content, 0644)
	if err != nil {
		logx.GetLogger().Sugar().Errorf("Write file failed: %s", err.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Write file failed!",
			"data": nil,
		})
		return
	}

	website.Save()

	enabledFileAbsPath := filepath.Join(availableBaseAbsPath, fileName)
	enabledFileIsExists := filex.FilePathExists(enabledFileAbsPath)
	if enabledFileIsExists {
		logx.GetLogger().Sugar().Infof("File already exists!: %s", enabledFileAbsPath)
		err := nginxx.Test()
		if err != nil {
			logx.GetLogger().Sugar().Errorf("There is a problem with the configuration file, please check: %s", err.Error())
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "There is a problem with the configuration file, please check!",
				"data": nil,
			})
			return
		}

		//output, err := nginxx.Reload()
		//if err != nil {
		//	logx.GetLogger().Sugar().Errorf("Reload failed: %s", output)
		//	logx.GetLogger().Sugar().Errorf("Reload failed: %s", err.Error())
		//	ctx.JSON(http.StatusOK, gin.H{
		//		"code": http.StatusInternalServerError,
		//		"msg":  "Reload failed!",
		//		"data": nil,
		//	})
		//	return
		//}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": nil,
		"msg":  "保存成功！",
	})
}

// Detail 详情
func Detail(ctx *gin.Context) {
	fileName := ctx.Param("fileName")

	website := new(websiteModel.Website)
	detail := website.Detail(fileName)
	if detail.WebsiteId <= 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": &detail,
			"msg":  "获取数据失败！",
		})
		return
	}

	availableConfPath, availablePathErr := nginxx.GetConfPath("sites-available")
	if availablePathErr != nil {
		logx.GetLogger().Sugar().Errorf("Get nginx conf path failed: %s", availablePathErr.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Get nginx conf path failed!",
			"data": nil,
		})
		return
	}
	availableBaseAbsPath, _ := filepath.Abs(availableConfPath)

	availableIsExistPath := filex.FilePathExists(availableBaseAbsPath)
	if !availableIsExistPath {
		logx.GetLogger().Sugar().Errorf("File or directory does not exist!: %s", availableBaseAbsPath)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "File or directory does not exist!",
			"data": nil,
		})
		return
	}

	availableFileAbsPath := filepath.Join(availableBaseAbsPath, fileName)
	availableFileIsExists := filex.FilePathExists(availableFileAbsPath)
	if !availableFileIsExists {
		logx.GetLogger().Sugar().Errorf("File does not exist!: %s", availableFileAbsPath)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "File does not exist!",
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": &detail,
		"msg":  "获取数据成功！",
	})
}

func Update(ctx *gin.Context) {
	form := new(WebsiteForm)
	ctx.Bind(form)

	websiteId := form.WebsiteId
	fileName := strings.TrimSpace(form.FileName)
	serverName := strings.TrimSpace(form.ServerName)
	rootDirectory := strings.TrimSpace(form.RootDirectory)
	homePage := strings.TrimSpace(form.HomePage)
	httpPort := form.HttpPort
	supportSsl := form.SupportSsl
	httpsPort := form.HttpsPort
	sslCertificate := strings.TrimSpace(form.SslCertificate)
	sslCertificateKey := strings.TrimSpace(form.SslCertificateKey)
	status := form.Status

	website := new(websiteModel.Website)
	fileNameCount := website.IsExistsFileName(fileName)
	if fileNameCount <= 0 {
		logx.GetLogger().Sugar().Error("File does not exist!")
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": nil,
			"msg":  "File does not exist!",
		})
		return
	}

	websiteMap := make(map[string]interface{})
	websiteMap["file_name"] = fileName
	websiteMap["server_name"] = serverName
	websiteMap["root_directory"] = rootDirectory
	websiteMap["home_page"] = homePage
	websiteMap["http_port"] = httpPort
	websiteMap["support_ssl"] = supportSsl
	websiteMap["https_port"] = httpsPort
	websiteMap["ssl_certificate"] = sslCertificate
	websiteMap["ssl_certificate_key"] = sslCertificateKey
	websiteMap["status"] = status
	websiteMap["update_time"] = datetimex.FormatNowDateTime()

	availableConfPath, availablePathErr := nginxx.GetConfPath("sites-available")
	if availablePathErr != nil {
		logx.GetLogger().Sugar().Errorf("Get nginx conf path failed: %s", availablePathErr.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
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
			"code": http.StatusBadRequest,
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

	availableFileAbsPath := filepath.Join(availableBaseAbsPath, fileName)
	availableFileIsExists := filex.FilePathExists(availableFileAbsPath)
	if !availableFileIsExists {
		logx.GetLogger().Sugar().Errorf("File does not exist!: %s", availableFileAbsPath)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "File does not exist!",
			"data": nil,
		})
		return
	}

	content := templatex.TemplateReplace(website)
	err := os.WriteFile(availableFileAbsPath, content, 0644)
	if err != nil {
		logx.GetLogger().Sugar().Errorf("Write file failed: %s", err.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Write file failed!",
			"data": nil,
		})
		return
	}

	website.Update(websiteId, websiteMap)

	enabledFileAbsPath := filepath.Join(enabledBaseAbsPath, fileName)
	enabledFileIsExists := filex.FilePathExists(enabledFileAbsPath)
	if enabledFileIsExists {
		logx.GetLogger().Sugar().Infof("File already exists!: %s", enabledFileAbsPath)
		err := nginxx.Test()
		if err != nil {
			logx.GetLogger().Sugar().Errorf("There is a problem with the configuration file, please check: %s", err.Error())
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "There is a problem with the configuration file, please check!",
				"data": nil,
			})
			return
		}

		//output, err := nginxx.Reload()
		//if err != nil {
		//	logx.GetLogger().Sugar().Errorf("Reload failed: %s", output)
		//	logx.GetLogger().Sugar().Errorf("Reload failed: %s", err.Error())
		//	ctx.JSON(http.StatusOK, gin.H{
		//		"code": http.StatusInternalServerError,
		//		"msg":  "Reload failed!",
		//		"data": nil,
		//	})
		//	return
		//}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": nil,
		"msg":  "修改成功！",
	})
}

func Enable(ctx *gin.Context) {
	fileName := ctx.Param("fileName")

	website := new(websiteModel.Website)
	detail := website.Detail(fileName)
	if detail.WebsiteId <= 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": &detail,
			"msg":  "获取数据失败！",
		})
		return
	}

	availableConfPath, availablePathErr := nginxx.GetConfPath("sites-available")
	if availablePathErr != nil {
		logx.GetLogger().Sugar().Errorf("Get nginx conf path failed: %s", availablePathErr.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Get nginx conf path failed!",
			"data": nil,
		})
		return
	}
	availableBaseAbsPath, _ := filepath.Abs(availableConfPath)

	availableIsExistPath := filex.FilePathExists(availableBaseAbsPath)
	if !availableIsExistPath {
		logx.GetLogger().Sugar().Errorf("File or directory does not exist!: %s", availableBaseAbsPath)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "File or directory does not exist!",
			"data": nil,
		})
		return
	}

	availableFileAbsPath := filepath.Join(availableBaseAbsPath, fileName)
	availableFileIsExists := filex.FilePathExists(availableFileAbsPath)
	if !availableFileIsExists {
		logx.GetLogger().Sugar().Errorf("File does not exist!: %s", availableFileAbsPath)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "File does not exist!",
			"data": nil,
		})
		return
	}

	enabledConfPath, enabledPathErr := nginxx.GetConfPath("sites-enabled")
	if enabledPathErr != nil {
		logx.GetLogger().Sugar().Errorf("Get nginx conf path failed: %s", enabledPathErr.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
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

	enabledFileAbsPath := filepath.Join(enabledBaseAbsPath, fileName)
	enabledFileIsExists := filex.FilePathExists(enabledFileAbsPath)
	if enabledFileIsExists {
		logx.GetLogger().Sugar().Infof("File already exists!: %s", enabledFileAbsPath)
		logx.GetLogger().Sugar().Info("Service is enabled!")
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Service is enabled!",
			"data": nil,
		})
		return
	}

	linkErr := os.Symlink(availableFileAbsPath, enabledFileAbsPath)
	if linkErr != nil {
		logx.GetLogger().Sugar().Errorf("Enable failed: %s", linkErr.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Enable failed!",
			"data": nil,
		})
		return
	}

	err := nginxx.Test()
	if err != nil {
		_ = os.Remove(enabledFileAbsPath)
		logx.GetLogger().Sugar().Errorf("There is a problem with the configuration file, please check: %s", err.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "There is a problem with the configuration file, please check!",
			"data": nil,
		})
		return
	}

	//output, err := nginxx.Reload()
	//if err != nil {
	//	logx.GetLogger().Sugar().Errorf("Reload failed: %s", output)
	//	logx.GetLogger().Sugar().Errorf("Reload failed: %s", err.Error())
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"code": http.StatusInternalServerError,
	//		"msg":  "Reload failed!",
	//		"data": nil,
	//	})
	//	return
	//}

	websiteId := detail.WebsiteId

	websiteMap := make(map[string]interface{})
	websiteMap["file_name"] = fileName
	websiteMap["status"] = 1
	websiteMap["update_time"] = datetimex.FormatNowDateTime()

	website.Update(websiteId, websiteMap)

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Service enabled successfully!",
		"data": nil,
	})
}

func Disable(ctx *gin.Context) {
	fileName := ctx.Param("fileName")

	website := new(websiteModel.Website)
	detail := website.Detail(fileName)
	if detail.WebsiteId <= 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": &detail,
			"msg":  "获取数据失败！",
		})
		return
	}

	enabledConfPath, enabledPathErr := nginxx.GetConfPath("sites-enabled")
	if enabledPathErr != nil {
		logx.GetLogger().Sugar().Errorf("Get nginx conf path failed: %s", enabledPathErr.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
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

	enabledFileAbsPath := filepath.Join(enabledBaseAbsPath, fileName)
	enabledFileIsExists := filex.FilePathExists(enabledFileAbsPath)
	if !enabledFileIsExists {
		logx.GetLogger().Sugar().Errorf("File does not exist!: %s", enabledFileAbsPath)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "File does not exist!",
			"data": nil,
		})
		return
	}

	linkErr := os.Remove(enabledFileAbsPath)
	if linkErr != nil {
		logx.GetLogger().Sugar().Errorf("Disable failed: %s", linkErr.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Disable failed!",
			"data": nil,
		})
		return
	}

	//output, err := nginxx.Reload()
	//if err != nil {
	//	logx.GetLogger().Sugar().Errorf("Reload failed: %s", output)
	//	logx.GetLogger().Sugar().Errorf("Reload failed: %s", err.Error())
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"code": http.StatusInternalServerError,
	//		"msg":  "Reload failed!",
	//		"data": nil,
	//	})
	//	return
	//}

	websiteId := detail.WebsiteId

	websiteMap := make(map[string]interface{})
	websiteMap["file_name"] = fileName
	websiteMap["status"] = 0
	websiteMap["update_time"] = datetimex.FormatNowDateTime()

	website.Update(websiteId, websiteMap)

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Service disabled successfully!",
		"data": nil,
	})
}
