package service

import (
	"giligili/model"
	"giligili/serializer"
)

// CreateVideoService 创建视频投稿的服务
type CreateVideoService struct {
	Title string `form:"title" json:"title" binding:"required,min=3,max=30"`
	Info  string `form:"info" json:"info" binding:"required,max=40"`
}

// Create 创建视频
func (service *CreateVideoService) Create() serializer.Response {
	// user := sessions.Default(c).Get("user").(model.User)
	video := model.Video{
		Title: service.Title,
		Info:  service.Info,
	}
	if err := model.DB.Create(&video).Error; err != nil {
		return serializer.Response{
			Code:  50001,
			Msg:   "视频创建失败",
			Error: err.Error(),
		}
	}

	return serializer.BuildVideoResponse(video)
}
