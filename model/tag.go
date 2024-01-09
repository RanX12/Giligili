package model

import "gorm.io/gorm"

// Tag 标签模型
type Tag struct {
	gorm.Model
	Name   string
	Videos []Video `gorm:"many2many:video_tags;"` // 自动创建 video_tags 表
}
