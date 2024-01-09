package service

import (
	"giligili/model"
	"giligili/serializer"
)

// ShowVideoCommentService 视频评论详情的服务
type ShowVideoCommentService struct {
}

// Show 展示视频评论详情
func (service *ShowVideoCommentService) Show(VideoId string) serializer.Response {
	var comments []model.Comment
	if err := model.DB.Preload("User").Preload("Video").Where("video_id = ?", VideoId).Find(&comments).Error; err != nil {
		return serializer.Response{
			Code:  404,
			Msg:   "视频不存在",
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Data: serializer.BuildComments(comments),
	}
}
