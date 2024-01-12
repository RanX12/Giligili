package service

import (
	"giligili/model"
	"giligili/serializer"

	"github.com/gin-gonic/gin"
)

// VideoLikeListService 某个视频点赞列表的服务
type VideoLikeListService struct {
}

func (service *VideoLikeListService) List(c *gin.Context) serializer.Response {
	var video model.Video
	if err := model.DB.First(&video, "id = ?", c.Param("id")).Error; err != nil {
		return serializer.Response{
			Code:  404,
			Msg:   "视频不存在",
			Error: err.Error(),
		}
	}

	var likes []model.Like

	if err := model.DB.Preload("Video").Where("video_id = ?", video.ID).Find(&likes).Error; err != nil {
		return serializer.Response{
			Code:  50001,
			Msg:   "视频点赞列表获取失败",
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Data: likes,
	}
}
