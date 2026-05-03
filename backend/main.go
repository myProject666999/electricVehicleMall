package main

import (
	"fmt"

	"electricVehicleMall/config"
	"electricVehicleMall/database"
	"electricVehicleMall/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 设置Gin模式
	if config.AppConfig.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化数据库
	database.InitDB()

	// 数据库迁移
	database.Migrate()

	// 设置路由
	r := routers.SetupRouter()

	// 启动服务器
	port := config.AppConfig.Server.Port
	fmt.Printf("Server starting on port %d...\n", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
