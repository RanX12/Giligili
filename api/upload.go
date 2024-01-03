package api

import (
	"giligili/service"

	"github.com/gin-gonic/gin"
)

// UploadToken 上传授权
func UploadToken(c *gin.Context) {
	service := service.UploadTokenService{}

	res := service.Post()
	c.JSON(200, res)
}
