package service

import (
	"giligili/model"
	"giligili/serializer"
)

// UserInfoService 某个用户信息的服务
type UserInfoService struct {
}

// Show 展示视频详情
func (service *UserInfoService) Show(id string) serializer.Response {
	var user model.User
	if err := model.DB.Preload("Video").First(&user, "id = ?", id).Error; err != nil {
		return serializer.Response{
			Code:  404,
			Msg:   "用户不存在",
			Error: err.Error(),
		}
	}

	return serializer.BuildUserResponse(user)
}
