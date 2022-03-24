package loginLogModel

import (
	"fmt"
	"github.com/wuchunfu/nginx-web/middleware/databasex"
	"github.com/wuchunfu/nginx-web/middleware/logx"
)

// LoginLog 用户 model
type LoginLog struct {
	Id         int64  `json:"id" gorm:"primary_key; auto_increment; not null"` // 设置字段 id 自增类型。
	Username   string `json:"username" gorm:"type:varchar(50); not null"`      // 用户名
	Ip         string `json:"ip" gorm:"type:varchar(64); not null "`           // ip 地址
	CreateTime string `json:"createTime" gorm:"type:varchar(50); not null"`    // 创建时间
}

func (loginLog *LoginLog) List(page int, pageSize int, username string) ([]LoginLog, int64) {
	db := databasex.GetDB()
	if username != "" {
		db = db.Where("username like ?", fmt.Sprintf("%%%s%%", username))
	}
	list := make([]LoginLog, 0)
	findErr := db.Model(&loginLog).Offset((page - 1) * pageSize).Limit(pageSize).Order("create_time desc").Find(&list)
	if findErr.Error != nil {
		logx.GetLogger().Sugar().Error(findErr.Error)
		return nil, -1
	}
	var count int64
	countErr := db.Model(&loginLog).Count(&count)
	if countErr.Error != nil {
		logx.GetLogger().Sugar().Error(countErr.Error)
		return nil, -1
	}
	return list, count
}

// Save 新增
func (loginLog *LoginLog) Save() {
	db := databasex.GetDB().Begin()
	err := db.Model(&loginLog).Create(&loginLog)
	if err.Error != nil {
		db.Rollback()
		logx.GetLogger().Sugar().Error(err.Error)
	}
	db.Commit()
}
