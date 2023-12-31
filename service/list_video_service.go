package service

import (
	"giligili/model"
	"giligili/serializer"
)

// ListVideoService 视频列表的服务
type ListVideoService struct {
}

// List 视频列表
func (service *ListVideoService) List() serializer.Response {
	var videos []model.Video
	if err := model.DB.Find(&videos).Error; err != nil {
		return serializer.Response{
			Code:  400,
			Msg:   "视频列表获取失败",
			Error: err.Error(),
		}
	}

	return serializer.BuildVideosResponse(videos)
}
