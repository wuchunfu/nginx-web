package configApi

import (
	"github.com/gin-gonic/gin"
	"github.com/wuchunfu/nginx-web/middleware/configx"
	"github.com/wuchunfu/nginx-web/middleware/logx"
	"github.com/wuchunfu/nginx-web/utils/bytex"
	"github.com/wuchunfu/nginx-web/utils/datetimex"
	"github.com/wuchunfu/nginx-web/utils/filetypex"
	"github.com/wuchunfu/nginx-web/utils/filex"
	"github.com/wuchunfu/nginx-web/utils/nginxx"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type FileList struct {
	BasePath   string `json:"basePath"`   // 基础目录
	ParentPath string `json:"parentPath"` // 文件父目录
	FilePath   string `json:"filePath"`   // 文件路径
	FileName   string `json:"fileName"`   // 文件名
	IsFile     bool   `json:"isFile"`     // 是否是文件
	FileType   string `json:"fileType"`   // 文件类型
	FileSize   string `json:"fileSize"`   // 文件大小
	SuffixName string `json:"suffixName"` // 文件后缀名
	DateTime   string `json:"dateTime"`   // 文件修改时间
}

type ConfigForm struct {
	ParentPath string `json:"parentPath" binding:"required"` // 文件路径
	Content    string `json:"content" binding:"required"`    // 文件内容
}

func List(ctx *gin.Context) {
	setting := configx.ServerSetting
	storagePath := setting.System.StoragePath
	confPath, pathErr := nginxx.GetConfPath(storagePath)
	if pathErr != nil {
		logx.GetLogger().Sugar().Errorf("Get nginx conf path failed: %s", pathErr.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Get nginx conf path failed!",
			"data": nil,
		})
		return
	}
	baseAbsPath, _ := filepath.Abs(confPath)

	isExistPath := filex.FilePathExists(baseAbsPath)
	if !isExistPath {
		logx.GetLogger().Sugar().Errorf("File or directory does not exist!: %s", baseAbsPath)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "File or directory does not exist!",
			"data": nil,
		})
		return
	}

	parentPath := ctx.Query("parentPath")
	logx.GetLogger().Sugar().Infof("parentPath: %s", parentPath)

	fileList := &[]FileList{}
	if parentPath == "" {
		isDir := filex.IsDir(baseAbsPath)
		if !isDir {
			logx.GetLogger().Sugar().Errorf("File or directory does not exist!: %s", baseAbsPath)
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "File or directory does not exist!",
				"data": nil,
			})
			return
		}
		fileList = ListFolder(baseAbsPath)
	} else {
		parentAbsPath := filepath.Join(baseAbsPath, parentPath)
		logx.GetLogger().Sugar().Infof("parentAbsPath: %s", parentAbsPath)

		isDir := filex.IsDir(parentAbsPath)
		if !isDir {
			logx.GetLogger().Sugar().Errorf("File or directory does not exist!: %s", parentAbsPath)
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "File or directory does not exist!",
				"data": nil,
			})
			return
		}
		fileList = ListFolder(parentAbsPath)
	}
	// 返回目录json数据
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Get data successfully!",
		"data": fileList,
	})
}

// Detail 文件详情
func Detail(ctx *gin.Context) {
	setting := configx.ServerSetting
	storagePath := setting.System.StoragePath
	confPath, pathErr := nginxx.GetConfPath(storagePath)
	if pathErr != nil {
		logx.GetLogger().Sugar().Errorf("Get nginx conf path failed: %s", pathErr.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Get nginx conf path failed!",
			"data": nil,
		})
		return
	}
	baseAbsPath, _ := filepath.Abs(confPath)

	isExistPath := filex.FilePathExists(baseAbsPath)
	if !isExistPath {
		logx.GetLogger().Sugar().Errorf("File or directory does not exist!: %s", baseAbsPath)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "File or directory does not exist!",
			"data": nil,
		})
		return
	}

	parentPath := ctx.Query("parentPath")
	logx.GetLogger().Sugar().Infof("parentPath: %s", parentPath)

	if parentPath != "" {
		parentAbsPath := filepath.Join(baseAbsPath, parentPath)
		logx.GetLogger().Sugar().Infof("parentAbsPath: %s", parentAbsPath)

		exists := filex.FilePathExists(parentAbsPath)
		if !exists {
			logx.GetLogger().Sugar().Errorf("File or directory does not exist!: %s", parentAbsPath)
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "File or directory does not exist!",
				"data": nil,
			})
			return
		}

		isFile := filex.IsFile(parentAbsPath)
		if !isFile {
			logx.GetLogger().Sugar().Errorf("File or directory does not exist!: %s", parentAbsPath)
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "File or directory does not exist!",
				"data": nil,
			})
			return
		}

		originalContent, err := os.ReadFile(parentAbsPath)
		if err != nil {
			logx.GetLogger().Sugar().Errorf("Read file failed: %s", err.Error())
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "Read file failed!",
				"data": nil,
			})
			return
		}
		// 返回目录json数据
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "Get data successfully!",
			"data": string(originalContent),
		})
	}
}

// Update 文件详情
func Update(ctx *gin.Context) {
	setting := configx.ServerSetting
	storagePath := setting.System.StoragePath
	confPath, pathErr := nginxx.GetConfPath(storagePath)
	if pathErr != nil {
		logx.GetLogger().Sugar().Errorf("Get nginx conf path failed: %s", pathErr.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Get nginx conf path failed!",
			"data": nil,
		})
		return
	}
	baseAbsPath, _ := filepath.Abs(confPath)

	isExistPath := filex.FilePathExists(baseAbsPath)
	if !isExistPath {
		logx.GetLogger().Sugar().Errorf("File or directory does not exist!: %s", baseAbsPath)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "File or directory does not exist!",
			"data": nil,
		})
		return
	}

	form := new(ConfigForm)
	ctx.Bind(form)

	parentPath := form.ParentPath
	content := form.Content

	logx.GetLogger().Sugar().Infof("parentPath: %s", parentPath)

	if parentPath != "" {
		parentAbsPath := filepath.Join(baseAbsPath, parentPath)
		logx.GetLogger().Sugar().Infof("parentAbsPath: %s", parentAbsPath)

		exists := filex.FilePathExists(parentAbsPath)
		if !exists {
			logx.GetLogger().Sugar().Errorf("File or directory does not exist!: %s", parentAbsPath)
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "File or directory does not exist!",
				"data": nil,
			})
			return
		}

		isFile := filex.IsFile(parentAbsPath)
		if !isFile {
			logx.GetLogger().Sugar().Errorf("File or directory does not exist!: %s", parentAbsPath)
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "File or directory does not exist!",
				"data": nil,
			})
			return
		}

		originalContent, err := os.ReadFile(parentAbsPath)
		if err != nil {
			logx.GetLogger().Sugar().Errorf("Read file failed: %s", err.Error())
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "Read file failed!",
				"data": nil,
			})
			return
		}

		if content != "" && content != string(originalContent) {
			err := os.WriteFile(parentAbsPath, []byte(content), 0644)
			if err != nil {
				logx.GetLogger().Sugar().Errorf("Write file failed: %s", err.Error())
				ctx.JSON(http.StatusOK, gin.H{
					"code": http.StatusInternalServerError,
					"msg":  "Write file failed!",
					"data": nil,
				})
				return
			}
		}

		//output, err := nginxx.Reload()
		//if err != nil {
		//	logx.GetLogger().Sugar().Infof("Reload nginx failed: %s", output)
		//	logx.GetLogger().Sugar().Errorf("Reload nginx failed: %s", err.Error())
		//	ctx.JSON(http.StatusOK, gin.H{
		//		"code": http.StatusInternalServerError,
		//		"msg":  "Reload nginx failed!",
		//		"data": nil,
		//	})
		//	return
		//}

		newContent, err := os.ReadFile(parentAbsPath)
		if err != nil {
			logx.GetLogger().Sugar().Errorf("Read file failed: %s", err.Error())
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "Read file failed!",
				"data": nil,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "Get data successfully!",
			"data": string(newContent),
		})
	}
}

// ChangeFolder 切换文件夹
func ChangeFolder(ctx *gin.Context) {
	setting := configx.ServerSetting
	storagePath := setting.System.StoragePath
	confPath, pathErr := nginxx.GetConfPath(storagePath)
	if pathErr != nil {
		logx.GetLogger().Sugar().Errorf("Get nginx conf path failed: %s", pathErr.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Get nginx conf path failed!",
			"data": nil,
		})
		return
	}
	baseAbsPath, _ := filepath.Abs(confPath)

	isExistPath := filex.FilePathExists(baseAbsPath)
	if !isExistPath {
		logx.GetLogger().Sugar().Errorf("File or directory does not exist!: %s", baseAbsPath)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "File or directory does not exist!",
			"data": nil,
		})
		return
	}

	parentPath := ctx.Query("parentPath")
	logx.GetLogger().Sugar().Infof("parentPath: %s", parentPath)

	fileList := &[]FileList{}
	if parentPath != "" {
		parentAbsPath := filepath.Join(baseAbsPath, parentPath)
		logx.GetLogger().Sugar().Infof("parentAbsPath: %s", parentAbsPath)

		exists := filex.FilePathExists(parentAbsPath)
		if !exists {
			logx.GetLogger().Sugar().Errorf("File or directory does not exist!: %s", parentAbsPath)
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "File or directory does not exist!",
				"data": nil,
			})
			return
		}

		isDir := filex.IsDir(parentAbsPath)
		if !isDir {
			logx.GetLogger().Sugar().Errorf("File or directory does not exist!: %s", parentAbsPath)
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "File or directory does not exist!",
				"data": nil,
			})
			return
		}
		parentPath = filepath.Dir(parentAbsPath)
		fileList = ListFolder(parentPath)
	}
	// 返回目录json数据
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Folder switch successfully!",
		"data": fileList,
	})
}

// ListFolder 列出文件夹中的文件夹及文件
func ListFolder(parentAbsPath string) *[]FileList {
	setting := configx.ServerSetting
	storagePath := setting.System.StoragePath
	confPath, pathErr := nginxx.GetConfPath(storagePath)
	if pathErr != nil {
		logx.GetLogger().Sugar().Errorf("Get nginx conf path failed: %s", pathErr.Error())
		return nil
	}
	basePath, _ := filepath.Abs(confPath)

	exists := filex.FilePathExists(parentAbsPath)
	if !exists {
		logx.GetLogger().Sugar().Errorf("File or directory does not exist!: %s", parentAbsPath)
		return nil
	}

	fileList := make([]FileList, 0)
	// 遍历目录，读出文件名、大小
	err := filepath.WalkDir(parentAbsPath, func(fileAbsPath string, info fs.DirEntry, err error) error {
		fileInfo, err := info.Info()
		if nil == fileInfo {
			return err
		}

		if fileAbsPath == parentAbsPath {
			return nil
		}

		fileName := info.Name()
		fullPath := filepath.Join(parentAbsPath, fileName)
		parentPath, _ := filepath.Rel(basePath, parentAbsPath)
		filePath, _ := filepath.Rel(basePath, fullPath)

		if filePath == fileName {
			parentPath = ""
			filePath = ""
		}

		suffixName := strings.ToLower(path.Ext(fileName))
		fileType := filetypex.FileType(suffixName)
		fileSize := bytex.FormatFileSize(fileInfo.Size())
		dateTime := datetimex.FormatDateTime(fileInfo.ModTime())

		list := &FileList{}
		if info.IsDir() {
			list = &FileList{
				BasePath:   basePath,
				ParentPath: parentPath,
				FilePath:   filePath,
				FileName:   fileName,
				IsFile:     false,
				FileSize:   "-",
				DateTime:   dateTime,
			}
			fileList = append(fileList, *list)
			return filepath.SkipDir
		}

		if fileAbsPath != parentAbsPath {
			list = &FileList{
				BasePath:   basePath,
				ParentPath: parentPath,
				FilePath:   filePath,
				FileName:   fileName,
				IsFile:     true,
				FileType:   fileType,
				SuffixName: suffixName,
				FileSize:   fileSize,
				DateTime:   dateTime,
			}
			fileList = append(fileList, *list)
		}
		return nil
	})

	if err != nil {
		logx.GetLogger().Sugar().Errorf("Traversal file directory failed!: %s", err.Error())
		panic(err.Error())
	}
	return &fileList
}
