package service

import (
	"giligili/model"
	"giligili/serializer"

	"github.com/gin-gonic/gin"
)

// UserLikeVideoService 某个用户点赞视频的服务
type UserLikeVideoService struct {
}

// Like 点赞视频详情
func (service *UserLikeVideoService) Like(c *gin.Context) serializer.Response {
	user := CurrentUser(c)
	id := c.Param("id")

	var video model.Video
	if err := model.DB.First(&video, "id = ?", id).Error; err != nil {
		return serializer.Response{
			Code:  50001,
			Msg:   "视频 ID 验证失败",
			Error: err.Error(),
		}
	}

	like := model.Like{
		UserId:  user.ID,
		VideoId: video.ID,
	}

	if err := model.DB.Create(&like).Error; err != nil {
		return serializer.Response{
			Code:  50001,
			Msg:   "视频点赞失败",
			Error: err.Error(),
		}
	}

	return serializer.Response{}
}

// UnLike 点赞视频详情
func (service *UserLikeVideoService) UnLike(c *gin.Context) serializer.Response {
	user := CurrentUser(c)
	id := c.Param("id")

	var video model.Video
	if err := model.DB.First(&video, "id = ?", id).Error; err != nil {
		return serializer.Response{
			Code:  50001,
			Msg:   "视频 ID 验证失败",
			Error: err.Error(),
		}
	}

	like := model.Like{
		UserId:  user.ID,
		VideoId: video.ID,
	}

	if err := model.DB.Delete(&like).Error; err != nil {
		return serializer.Response{
			Code:  50001,
			Msg:   "视频取消点赞失败",
			Error: err.Error(),
		}
	}

	return serializer.Response{}
}
