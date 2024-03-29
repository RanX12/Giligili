package server

import (
	"giligili/api"
	"giligili/middleware"
	"giligili/service"
	"os"

	"github.com/gin-gonic/gin"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()

	// 中间件, 顺序不能改
	r.Use(middleware.Session(os.Getenv("SESSION_SECRET")))
	r.Use(middleware.Cors())
	r.Use(middleware.CurrentUser())

	// 路由
	v1 := r.Group("/api/v1")
	{
		v1.POST("ping", api.Ping)

		// 用户登录
		v1.POST("user/register", api.UserRegister)

		// 用户登录
		v1.POST("user/login", api.UserLogin)

		// 用户是否登录
		// v1.GET("user/check_status", api.CheckUserLoginStatus)

		// 需要登录保护的
		auth := v1.Group("")
		auth.Use(middleware.AuthRequired())
		{
			// User Routing
			auth.GET("user/me", api.UserMe)
			auth.DELETE("user/logout", api.UserLogout)

			// 投稿视频
			auth.GET("videos", api.ListVideo)
			auth.PUT("videos/:id", api.UpdateVideo)
			auth.DELETE("videos/:id", api.DeleteVideo)

			// 点赞视频
			auth.GET("videos/:id/likes", api.LikeList)
			auth.POST("videos/:id/like", api.LikeVideo)
			auth.DELETE("videos/:id/like", api.UnLikeVideo)

			auth.POST("upload/token", api.UploadToken)
		}

		v1.POST("videos", api.CreateVideo)
		v1.GET("videos/:id", api.ShowVideo)
		v1.GET("videos/:id/comments", api.ShowVideoComments)
		// 使用 WebSocket 中间件
		v1.GET("ws", service.WsHandler)
	}
	return r
}
