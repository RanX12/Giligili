package api

import (
	"giligili/service"

	"github.com/gin-gonic/gin"
)

// ShowVideoComments 视频评论详情接口
func ShowVideoComments(c *gin.Context) {
	service := service.ShowVideoCommentService{}

	res := service.Show(c.Param("id"))
	c.JSON(200, res)
}
