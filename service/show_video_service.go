package service

import (
	"giligili/model"
	"giligili/serializer"
)

// ShowVideoService 视频详情的服务
type ShowVideoService struct {
}

// Show 展示视频详情
func (service *ShowVideoService) Show(id string) serializer.Response {
	var video model.Video
	if err := model.DB.Preload("User").First(&video, "id = ?", id).Error; err != nil {
		return serializer.Response{
			Code:  404,
			Msg:   "视频不存在",
			Error: err.Error(),
		}
	}

	return serializer.BuildVideoResponse(video)
}
