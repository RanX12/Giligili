package service

import (
	"giligili/model"
	"giligili/serializer"
)

// DeleteVideoService 删除视频的服务
type DeleteVideoService struct {
}

// Delete 删除视频
func (service *DeleteVideoService) Delete(id string) serializer.Response {
	var video model.Video
	if err := model.DB.Delete(&video, "id = ?", id).Error; err != nil {
		return serializer.Response{
			Code:  400,
			Msg:   "视频删除失败",
			Error: err.Error(),
		}
	}
	return serializer.Response{}
}
