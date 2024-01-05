package service

import (
	"giligili/model"
	"giligili/serializer"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// CheckUserService 创建视频投稿的服务
type CheckUserService struct{}

// Check 检查用户是否登录
func (service *CheckUserService) Check(c *gin.Context) serializer.Response {
	s := sessions.Default(c)
	user_id := s.Get("user_id")

	if user_id == nil {
		return serializer.Err(
			401,
			"未登录",
			nil,
		)
	} else {
		var user model.User

		if err := model.DB.First(&user, "id = ?", user_id); err != nil {
			return serializer.Err(
				401,
				"未登录",
				nil,
			)
		}

		return serializer.BuildUserResponse(user)
	}
}
