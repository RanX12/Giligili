package model

import (
	"gorm.io/gorm"
)

// Comment 评论模型
type Comment struct {
	gorm.Model
	UserId   uint `gorm:"index"`
	User     User
	VideoId  uint `gorm:"index"`
	Video    Video
	Content  string
	ParentID *uint     // 为空则为顶级评论
	Parent   *Comment  `gorm:"foreignKey:ParentID"`
	Replies  []Comment `gorm:"foreignKey:ParentID"`
}
