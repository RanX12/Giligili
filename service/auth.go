package service

import (
	"giligili/model"

	"github.com/gin-gonic/gin"
)

func CurrentUser(c *gin.Context) *model.User {
	if user, _ := c.Get("user"); user != nil {
		if u, ok := user.(*model.User); ok {
			return u
		}
	}
	return nil
}
