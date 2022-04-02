package websiteModel

import (
	"github.com/wuchunfu/nginx-web/middleware/databasex"
	"github.com/wuchunfu/nginx-web/middleware/logx"
)

// Website 网站 model
type Website struct {
	WebsiteId         int64  `json:"websiteId" gorm:"primary_key; auto_increment; not null"`          // 设置字段 websiteId 自增类型。
	FileName          string `json:"fileName" gorm:"type:varchar(50); not null; default ''; unique"`  // 配置文件名称
	ServerName        string `json:"serverName" gorm:"type:varchar(50); not null; default ''"`        // 网站域名
	RootDirectory     string `json:"rootDirectory" gorm:"type:varchar(50); not null; default ''"`     // 网站根目录
	HomePage          string `json:"homePage" gorm:"type:varchar(50); not null; default ''"`          // 网站主页
	HttpPort          int16  `json:"httpPort" gorm:"tinyint; not null; default 80"`                   // http 监听端口
	SupportSsl        int8   `json:"supportSsl" gorm:"tinyint; not null; default 0"`                  // 启用 TLS 状态 1: 启用 0: 禁用
	HttpsPort         int16  `json:"httpsPort" gorm:"tinyint; not null; default 443"`                 // https 监听端口
	SslCertificate    string `json:"sslCertificate" gorm:"type:varchar(50); not null; default ''"`    // TLS 证书路径
	SslCertificateKey string `json:"sslCertificateKey" gorm:"type:varchar(50); not null; default ''"` // 私钥路径
	Status            int8   `json:"status" gorm:"tinyint; not null; default 0"`                      // 启用状态 1: 正常 0:禁用
	CreateTime        string `json:"createTime" gorm:"type:varchar(50); not null"`
	UpdateTime        string `json:"updateTime" gorm:"type:varchar(50);"`
}

func (website *Website) Detail(fileName string) *Website {
	db := databasex.GetDB()
	err := db.Model(&website).Where("file_name = ?", fileName).Find(&website)
	if err.Error != nil {
		logx.GetLogger().Sugar().Error(err.Error)
		return nil
	}
	return website
}

// Save 新增
func (website *Website) Save() {
	db := databasex.GetDB().Begin()
	err := db.Model(&website).Create(&website)
	if err.Error != nil {
		db.Rollback()
		logx.GetLogger().Sugar().Error(err.Error)
	}
	db.Commit()
}

// Update 修改
func (website *Website) Update(websiteId int64, fieldMap map[string]interface{}) {
	db := databasex.GetDB().Begin()
	err := db.Model(&website).Where("website_id = ?", websiteId).Updates(fieldMap)
	if err.Error != nil {
		db.Rollback()
		logx.GetLogger().Sugar().Error(err.Error)
	}
	db.Commit()
}

// Delete 删除
func (website *Website) Delete(websiteId int64) {
	db := databasex.GetDB().Begin()
	err := db.Model(&website).Where("website_id = ?", websiteId).Delete(&website)
	if err.Error != nil {
		db.Rollback()
		logx.GetLogger().Sugar().Error(err.Error)
	}
	db.Commit()
}

// IsExistsFileName 用户名是否存在
func (website *Website) IsExistsFileName(fileName string) int64 {
	db := databasex.GetDB()
	var count int64
	err := db.Model(&website).Where("file_name = ?", fileName).Count(&count)
	if err.Error != nil {
		logx.GetLogger().Sugar().Error(err.Error)
		return -1
	}
	return count
}
