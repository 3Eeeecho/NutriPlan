package router

import (
	"NutriPlan/internal/api/handler"
	"NutriPlan/internal/service"
	"NutriPlan/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// RouterDeps 结构体用于接收所有需要的依赖服务
type RouterDeps struct {
	UserService   service.UserService
	RecipeService service.RecipeService
	IntakeService service.IntakeService
}

// NewRouter 初始化并配置 Gin 路由
func NewRouter(deps RouterDeps) *gin.Engine {
	// 生产环境应使用 gin.ReleaseMode
	r := gin.Default()

	// 跨域中间件 (CORS) - 生产环境需严格配置
	r.Use(gin.Recovery())

	// 实例化 Handler
	userHandler := handler.NewUserHandler(deps.UserService)
	recipeHandler := handler.NewRecipeHandler(deps.RecipeService, deps.UserService)
	intakeHandler := handler.NewIntakeHandler(deps.IntakeService)

	// 基础路由组
	v1 := r.Group("/api/v1")
	{
		// --- 用户认证/档案路由 ---

		// 注册接口 (POST /api/v1/register)
		v1.POST("/register", userHandler.Register)

		// 登录接口 (POST /api/v1/login)
		v1.POST("/login", userHandler.Login)

		// 受保护的路由：需要 JWT 认证
		auth := v1.Group("/user")
		auth.Use(jwt.AuthMiddleware()) // 认证中间件
		{
			// 更新健康档案 (PUT /api/v1/user/profile)
			auth.PUT("/profile", userHandler.UpdateProfile)

			// 获取用户档案 (GET /api/v1/user/profile)
			auth.GET("/profile", userHandler.GetProfile)

			// 获取营养需求 (GET /api/v1/user/nutrition)
			auth.GET("/nutrition", userHandler.GetNutritionRequirements)
		}

		// --- 食谱推荐路由 ---
		recipes := v1.Group("/recipes")
		recipes.Use(jwt.AuthMiddleware()) // 需要认证
		{
			// 获取食谱推荐 (GET /api/v1/recipes/recommend)
			recipes.GET("/recommend", recipeHandler.GetRecommendations)

			// 选择食谱计划 (POST /api/v1/recipes/plan/:planId/select)
			recipes.POST("/plan/:planId/select", recipeHandler.SelectPlan)

			// 保存推荐结果 (POST /api/v1/recipes/recommend)
			recipes.POST("/recommend", recipeHandler.SaveRecommendations)

			// 获取已选食谱计划 (GET /api/v1/recipes/selected)
			recipes.GET("/selected", recipeHandler.GetSelectedPlan)

			// 获取食谱详情 (GET /api/v1/recipes/:id)
			recipes.GET("/:id", recipeHandler.GetRecipeDetail)

			// 添加收藏 (POST /api/v1/recipes/:id/favorite)
			recipes.POST("/:id/favorite", recipeHandler.AddFavorite)

			// 取消收藏 (DELETE /api/v1/recipes/:id/favorite)
			recipes.DELETE("/:id/favorite", recipeHandler.RemoveFavorite)

			// 获取收藏列表 (GET /api/v1/recipes/favorites)
			recipes.GET("/favorites", recipeHandler.GetFavoriteList)
		}

		// --- 饮食记录路由 ---
		intake := v1.Group("/intake")
		intake.Use(jwt.AuthMiddleware()) // 需要认证
		{
			// 添加饮食记录 (POST /api/v1/intake/records)
			intake.POST("/records", intakeHandler.AddIntakeRecord)

			// 删除饮食记录 (DELETE /api/v1/intake/records/:id)
			intake.DELETE("/records/:id", intakeHandler.DeleteIntakeRecord)

			// 获取当日营养状态 (GET /api/v1/intake/today)
			intake.GET("/today", intakeHandler.GetTodayStatus)

			// 获取周报告 (GET /api/v1/intake/weekly)
			intake.GET("/weekly", intakeHandler.GetWeeklyReport)
		}
	}

	return r
}
