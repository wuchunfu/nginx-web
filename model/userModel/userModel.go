package userModel

import (
	"fmt"
	"github.com/wuchunfu/nginx-web/common"
	"github.com/wuchunfu/nginx-web/middleware/databasex"
	"github.com/wuchunfu/nginx-web/middleware/logx"
)

// User 用户model
type User struct {
	UserId     int64  `json:"userId" gorm:"primary_key; auto_increment; not null"`         // 设置字段 userId 自增类型。
	Username   string `json:"username" gorm:"type:varchar(50); not null; unique"`          // 用户名
	Password   string `json:"-" gorm:"type:char(50); not null "`                           // 密码，json 返回忽略本字段
	Salt       string `json:"-" gorm:"type:char(6); not null"`                             // 密码盐值，json 返回忽略本字段
	Email      string `json:"email" gorm:"type:varchar(50); not null; unique; default ''"` // 邮箱
	IsAdmin    int8   `json:"isAdmin" gorm:"tinyint; not null; default 0"`                 // 是否是管理员 1:管理员 0:普通用户
	Status     int8   `json:"status" gorm:"tinyint; not null; default 0"`                  // 1: 正常 0:禁用
	CreateTime string `json:"createTime" gorm:"type:varchar(50); not null"`
	UpdateTime string `json:"updateTime" gorm:"type:varchar(50);"`
}

func (user *User) List(page int, pageSize int, username string) ([]User, int64) {
	db := databasex.GetDB()
	if username != "" {
		db = db.Where("username like ?", fmt.Sprintf("%%%s%%", username))
	}
	list := make([]User, 0)
	findErr := db.Model(&user).Offset((page - 1) * pageSize).Limit(pageSize).Order("update_time desc").Find(&list)
	if findErr.Error != nil {
		logx.GetLogger().Sugar().Error(findErr.Error)
		return nil, -1
	}
	var count int64
	countErr := db.Model(&user).Count(&count)
	if countErr.Error != nil {
		logx.GetLogger().Sugar().Error(countErr.Error)
		return nil, -1
	}
	return list, count
}

func (user *User) Detail(userId int64) *User {
	db := databasex.GetDB()
	err := db.Model(&user).Where("user_id = ?", userId).Find(&user)
	if err.Error != nil {
		logx.GetLogger().Sugar().Error(err.Error)
		return nil
	}
	return user
}

// Save 新增
func (user *User) Save() {
	db := databasex.GetDB().Begin()
	err := db.Model(&user).Create(&user)
	if err.Error != nil {
		db.Rollback()
		logx.GetLogger().Sugar().Error(err.Error)
	}
	db.Commit()
}

// Update 修改
func (user *User) Update(userId int64, fieldMap map[string]interface{}) {
	db := databasex.GetDB().Begin()
	err := db.Model(&user).Where("user_id = ?", userId).Updates(fieldMap)
	if err.Error != nil {
		db.Rollback()
		logx.GetLogger().Sugar().Error(err.Error)
	}
	db.Commit()
}

// ChangePassword 修改密码
func (user *User) ChangePassword(userId int64, fieldMap map[string]interface{}) {
	db := databasex.GetDB().Begin()
	err := db.Model(&user).Where("user_id = ?", userId).Updates(fieldMap)
	if err.Error != nil {
		db.Rollback()
		logx.GetLogger().Sugar().Error(err.Error)
	}
	db.Commit()
}

// Delete 删除
func (user *User) Delete(userId int64) {
	db := databasex.GetDB().Begin()
	err := db.Model(&user).Where("user_id = ?", userId).Delete(&user)
	if err.Error != nil {
		db.Rollback()
		logx.GetLogger().Sugar().Error(err.Error)
	}
	db.Commit()
}

// Enable 激活
func (user *User) Enable(userId int) int64 {
	db := databasex.GetDB().Begin()
	fieldMap := make(map[string]interface{})
	fieldMap["status"] = common.Enabled
	err := db.Model(&user).Where("user_id = ?", userId).Updates(fieldMap)
	if err.Error != nil {
		db.Rollback()
		logx.GetLogger().Sugar().Error(err.Error)
	}
	db.Commit()
	return user.UserId
}

// Disable 激活
func (user *User) Disable(userId int) int64 {
	db := databasex.GetDB().Begin()
	fieldMap := make(map[string]interface{})
	fieldMap["status"] = common.Disabled
	err := db.Model(&user).Where("user_id = ?", userId).Updates(fieldMap)
	if err.Error != nil {
		db.Rollback()
		logx.GetLogger().Sugar().Error(err.Error)
	}
	db.Commit()
	return user.UserId
}

// IsExistsUsername 用户名是否存在
func (user *User) IsExistsUsername(username string) int64 {
	db := databasex.GetDB()
	var count int64
	err := db.Model(&user).Where("username = ?", username).Count(&count)
	if err.Error != nil {
		logx.GetLogger().Sugar().Error(err.Error)
		return -1
	}
	return count
}

// IsExistsEmail 邮箱是否存在
func (user *User) IsExistsEmail(email string) int64 {
	db := databasex.GetDB()
	var count int64
	err := db.Model(&user).Where("email = ?", email).Count(&count)
	if err.Error != nil {
		logx.GetLogger().Sugar().Error(err.Error)
		return -1
	}
	return count
}

// GetInfoByName 根据 用户名获取用户信息
func (user *User) GetInfoByName(username string) *User {
	db := databasex.GetDB()
	err := db.Model(&user).Where("username = ?", username).Find(&user)
	if err.Error != nil {
		logx.GetLogger().Sugar().Error(err.Error)
		return nil
	}
	return user
}
