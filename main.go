package main

import (
	"giligili/conf"
	"giligili/server"
	"giligili/service"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 从配置文件读取配置
	conf.Init()

	// 装载路由
	gin.SetMode(os.Getenv("GIN_MODE"))
	go service.Manager.Start()
	r := server.NewRouter()
	r.Run(":3001")
}
