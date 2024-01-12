package model

import "gorm.io/gorm"

type Like struct {
	gorm.Model
	UserId  uint `gorm:"index"` // 用户 ID
	User    User
	VideoId uint `gorm:"index"` // 视频 ID
	Video   Video
}
