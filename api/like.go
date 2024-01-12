package api

import (
	"giligili/service"

	"github.com/gin-gonic/gin"
)

// LikeVideo 点赞视频
func LikeVideo(c *gin.Context) {
	service := service.UserLikeVideoService{}

	res := service.Like(c)
	c.JSON(200, res)
}

// UnLikeVideo 取消点赞
func UnLikeVideo(c *gin.Context) {
	service := service.UserLikeVideoService{}

	res := service.UnLike(c)
	c.JSON(200, res)
}

// LikeList 点赞列表
func LikeList(c *gin.Context) {
	service := service.VideoLikeListService{}
	res := service.List(c)
	c.JSON(200, res)
}
