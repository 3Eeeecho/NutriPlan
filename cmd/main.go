package main

import (
	"NutriPlan/internal/config"
	"NutriPlan/internal/core/service"
	"NutriPlan/internal/repository"
	"NutriPlan/internal/router"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}
	log.Println("配置加载成功")

	// 初始化数据库连接
	if err := repository.InitDatabase(config.AppConfig.Database); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	log.Println("数据库初始化成功")

	// 创建 Repository 实例
	userRepo := repository.NewGormUserRepository(repository.DB)

	// 创建 Service 实例
	nutriService := service.NewNutriService()
	userService := service.NewUserService(userRepo, nutriService)

	// 创建 Router 并注入依赖
	r := router.NewRouter(router.RouterDeps{
		UserService: userService,
	})

	// 根据环境设置 Gin 模式
	if config.AppConfig.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 启动服务器
	port := config.AppConfig.Server.Port
	if port == 0 {
		port = 8080
	}
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	log.Printf("服务器启动在端口 %d", port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
