package tokenModel

import (
	"github.com/wuchunfu/nginx-web/middleware/databasex"
	"github.com/wuchunfu/nginx-web/middleware/logx"
	"github.com/wuchunfu/nginx-web/utils/datetimex"
)

// UserToken 用户model
type UserToken struct {
	UserId     int64  `json:"userId" gorm:"primary_key; not null"`
	Token      string `json:"token" gorm:"type:varchar(200); not null; unique"`
	ExpireTime string `json:"expireTime" gorm:"type:varchar(50); not null"`
	UpdateTime string `json:"updateTime" gorm:"type:varchar(50);"`
}

// Save 新增
func (token *UserToken) Save(userId int64, tokenInfo map[string]string) {
	db := databasex.GetDB()
	idErr := db.Model(&token).Where("user_id = ?", userId).Find(&token)
	if idErr.Error != nil {
		logx.GetLogger().Sugar().Error(idErr.Error)
		token.UserId = userId
		token.Token = tokenInfo["token"]
		token.ExpireTime = tokenInfo["expireTime"]
		token.UpdateTime = datetimex.FormatNowDateTime()
		err := db.Model(&token).Create(&token)
		if err.Error != nil {
			logx.GetLogger().Sugar().Error(err.Error)
		}
	} else {
		token.Token = tokenInfo["token"]
		token.ExpireTime = tokenInfo["expireTime"]
		token.UpdateTime = datetimex.FormatNowDateTime()
		err := db.Model(&token).Where("user_id = ?", userId).Updates(&token)
		if err.Error != nil {
			logx.GetLogger().Sugar().Error(err.Error)
		}
	}
}

func (token *UserToken) Update(userId int64, tokenInfo map[string]string) {
	db := databasex.GetDB()
	token.Token = tokenInfo["token"]
	token.ExpireTime = tokenInfo["expireTime"]
	token.UpdateTime = datetimex.FormatNowDateTime()
	err := db.Model(&token).Where("user_id = ?", userId).Updates(&token)
	if err.Error != nil {
		logx.GetLogger().Sugar().Error(err.Error)
	}
}
