package handler

import (
	"NutriPlan/internal/repository/models"
	"NutriPlan/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RecipeHandler 食谱处理器
type RecipeHandler struct {
	recipeService service.RecipeService
	userService   service.UserService
}

// NewRecipeHandler 创建食谱处理器实例
func NewRecipeHandler(recipeService service.RecipeService, userService service.UserService) *RecipeHandler {
	return &RecipeHandler{
		recipeService: recipeService,
		userService:   userService,
	}
}

// RecommendRequest 推荐请求
type RecommendRequest struct {
	Count int `form:"count" json:"count"` // 推荐方案数量
}

// RecommendResponse 推荐响应
type RecommendResponse struct {
	Plans []DailyPlanDTO `json:"plans"`
}

// DailyPlanDTO 每日食谱计划DTO
type DailyPlanDTO struct {
	ID                 uint      `json:"id"`
	Breakfast          RecipeDTO `json:"breakfast"`
	Lunch              RecipeDTO `json:"lunch"`
	Dinner             RecipeDTO `json:"dinner"`
	Snack              RecipeDTO `json:"snack,omitempty"`
	TotalEnergy        float64   `json:"total_energy"`
	TotalProtein       float64   `json:"total_protein"`
	TotalCarbohydrate  float64   `json:"total_carbohydrate"`
	TotalFat           float64   `json:"total_fat"`
	TargetEnergy       float64   `json:"target_energy"`
	TargetProtein      float64   `json:"target_protein"`
	TargetCarbohydrate float64   `json:"target_carbohydrate"`
	TargetFat          float64   `json:"target_fat"`
	MatchScore         float64   `json:"match_score"`
}

// RecipeDTO 食谱DTO
type RecipeDTO struct {
	ID           uint     `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	ImageURL     string   `json:"image_url"`
	MealType     string   `json:"meal_type"`
	Energy       float64  `json:"energy"`
	Protein      float64  `json:"protein"`
	Carbohydrate float64  `json:"carbohydrate"`
	Fat          float64  `json:"fat"`
	DietaryFiber float64  `json:"dietary_fiber"`
	Ingredients  []string `json:"ingredients"`
	CookingSteps []string `json:"cooking_steps"`
	CookingTime  int      `json:"cooking_time"`
	Difficulty   string   `json:"difficulty"`
	IsFavorite   bool     `json:"is_favorite,omitempty"`
}

// GetRecommendations 获取食谱推荐
// @Summary 获取食谱推荐
// @Description 根据用户营养需求推荐每日食谱计划
// @Tags Recipe
// @Accept json
// @Produce json
// @Param count query int false "推荐方案数量" default(3)
// @Success 200 {object} RecommendResponse
// @Router /api/recipes/recommend [get]
func (h *RecipeHandler) GetRecommendations(c *gin.Context) {
	// 从上下文获取用户ID（需要认证中间件）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 通过 UserService 获取用户信息
	user, err := h.userService.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 解析请求参数
	var req RecommendRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Count = 3 // 默认推荐3套
	}

	// 限制推荐数量在3-5之间
	if req.Count < 3 {
		req.Count = 3
	}
	if req.Count > 5 {
		req.Count = 5
	}

	// 调用服务生成推荐
	plans, err := h.recipeService.RecommendRecipes(user, req.Count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成推荐失败: " + err.Error()})
		return
	}

	// 加载完整的食谱信息
	planDTOs := make([]DailyPlanDTO, 0, len(plans))
	for _, plan := range plans {
		planDTO, err := h.convertPlanToDTO(plan, userID.(uint))
		if err != nil {
			continue // 跳过转换失败的计划
		}
		planDTOs = append(planDTOs, planDTO)
	}

	c.JSON(http.StatusOK, RecommendResponse{
		Plans: planDTOs,
	})
}

// SelectPlan 选择食谱计划
// @Summary 选择食谱计划
// @Description 用户选择某套食谱方案
// @Tags Recipe
// @Accept json
// @Produce json
// @Param planId path int true "计划ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/recipes/plan/{planId}/select [post]
func (h *RecipeHandler) SelectPlan(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 解析计划ID
	planIDStr := c.Param("planId")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的计划ID"})
		return
	}

	// 调用服务选择计划
	if err := h.recipeService.SelectDailyPlan(userID.(uint), uint(planID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "选择计划失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "计划选择成功",
		"plan_id": planID,
	})
}

// GetSelectedPlan 获取已选食谱计划
// @Summary 获取已选食谱计划
// @Description 获取用户当前选中的每日食谱计划
// @Tags Recipe
// @Accept json
// @Produce json
// @Success 200 {object} DailyPlanDTO
// @Router /api/recipes/selected [get]
func (h *RecipeHandler) GetSelectedPlan(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取已选计划
	plan, err := h.recipeService.GetSelectedPlan(userID.(uint))
	if err != nil {
		// 如果没有找到已选计划，返回404
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到已选计划"})
		return
	}

	// 转换为DTO
	planDTO, err := h.convertPlanToDTO(plan, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据转换失败"})
		return
	}

	c.JSON(http.StatusOK, planDTO)
}

// SaveRecommendations 保存推荐结果
// @Summary 保存推荐结果
// @Description 保存推荐的食谱计划到数据库
// @Tags Recipe
// @Accept json
// @Produce json
// @Param plans body []DailyPlanDTO true "食谱计划列表"
// @Success 200 {object} map[string]interface{}
// @Router /api/recipes/recommend [post]
func (h *RecipeHandler) SaveRecommendations(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var plans []*models.DailyRecipePlan
	if err := c.ShouldBindJSON(&plans); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	// 设置用户ID
	for _, plan := range plans {
		plan.UserID = userID.(uint)
		if err := h.recipeService.SaveDailyPlan(plan); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存计划失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "保存成功",
		"count":   len(plans),
	})
}

// convertPlanToDTO 转换计划为DTO
func (h *RecipeHandler) convertPlanToDTO(plan *models.DailyRecipePlan, userID uint) (DailyPlanDTO, error) {
	// Repository 层已使用 Preload 预加载了关联的食谱数据
	dto := DailyPlanDTO{
		ID:                 plan.ID,
		TotalEnergy:        plan.TotalEnergy,
		TotalProtein:       plan.TotalProtein,
		TotalCarbohydrate:  plan.TotalCarbohydrate,
		TotalFat:           plan.TotalFat,
		TargetEnergy:       plan.TargetEnergy,
		TargetProtein:      plan.TargetProtein,
		TargetCarbohydrate: plan.TargetCarbohydrate,
		TargetFat:          plan.TargetFat,
		MatchScore:         plan.MatchScore,
	}

	// 转换预加载的食谱数据为 DTO 并标记收藏状态
	dto.Breakfast = convertRecipeToDTO(&plan.BreakfastRecipe)
	dto.Lunch = convertRecipeToDTO(&plan.LunchRecipe)
	dto.Dinner = convertRecipeToDTO(&plan.DinnerRecipe)
	dto.Snack = convertRecipeToDTO(&plan.SnackRecipe)

	if userID > 0 {
		if dto.Breakfast.ID > 0 {
			if ok, _ := h.recipeService.IsFavorite(userID, dto.Breakfast.ID); ok {
				dto.Breakfast.IsFavorite = true
			}
		}
		if dto.Lunch.ID > 0 {
			if ok, _ := h.recipeService.IsFavorite(userID, dto.Lunch.ID); ok {
				dto.Lunch.IsFavorite = true
			}
		}
		if dto.Dinner.ID > 0 {
			if ok, _ := h.recipeService.IsFavorite(userID, dto.Dinner.ID); ok {
				dto.Dinner.IsFavorite = true
			}
		}
		if dto.Snack.ID > 0 {
			if ok, _ := h.recipeService.IsFavorite(userID, dto.Snack.ID); ok {
				dto.Snack.IsFavorite = true
			}
		}
	}

	return dto, nil
}

// convertRecipeToDTO 转换食谱为DTO
func convertRecipeToDTO(recipe *models.Recipe) RecipeDTO {
	return RecipeDTO{
		ID:           recipe.ID,
		Name:         recipe.Name,
		ImageURL:     recipe.ImageURL,
		MealType:     string(recipe.MealType),
		Energy:       recipe.Energy,
		Protein:      recipe.Protein,
		Carbohydrate: recipe.Carbohydrate,
		Fat:          recipe.Fat,
		Ingredients:  recipe.Ingredients,
		CookingSteps: recipe.CookingSteps,
		CookingTime:  recipe.CookingTime,
		Difficulty:   recipe.Difficulty,
	}
}

// GetRecipeDetail 获取食谱详情
// @Summary 获取食谱详情
// @Description 根据ID获取单个食谱的详细信息
// @Tags Recipe
// @Accept json
// @Produce json
// @Param id path int true "食谱ID"
// @Success 200 {object} RecipeDTO
// @Router /api/recipes/{id} [get]
func (h *RecipeHandler) GetRecipeDetail(c *gin.Context) {
	// 解析食谱ID
	recipeIDStr := c.Param("id")
	recipeID, err := strconv.ParseUint(recipeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的食谱ID"})
		return
	}

	// 获取用户ID（可选，用于检查收藏状态）
	var userID uint
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(uint)
	}

	// 查询食谱和收藏状态
	recipe, isFavorite, err := h.recipeService.GetRecipeDetail(uint(recipeID), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "食谱不存在"})
		return
	}

	// 转换为DTO
	dto := convertRecipeToDTO(recipe)
	dto.IsFavorite = isFavorite

	c.JSON(http.StatusOK, dto)
}

// AddFavorite 添加收藏
// @Summary 添加收藏
// @Description 收藏某个食谱
// @Tags Recipe
// @Accept json
// @Produce json
// @Param id path int true "食谱ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/recipes/{id}/favorite [post]
func (h *RecipeHandler) AddFavorite(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 解析食谱ID
	recipeIDStr := c.Param("id")
	recipeID, err := strconv.ParseUint(recipeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的食谱ID"})
		return
	}

	// 添加收藏
	err = h.recipeService.AddFavorite(userID.(uint), uint(recipeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "收藏失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "收藏成功"})
}

// RemoveFavorite 取消收藏
// @Summary 取消收藏
// @Description 取消收藏某个食谱
// @Tags Recipe
// @Accept json
// @Produce json
// @Param id path int true "食谱ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/recipes/{id}/favorite [delete]
func (h *RecipeHandler) RemoveFavorite(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 解析食谱ID
	recipeIDStr := c.Param("id")
	recipeID, err := strconv.ParseUint(recipeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的食谱ID"})
		return
	}

	// 取消收藏
	err = h.recipeService.RemoveFavorite(userID.(uint), uint(recipeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取消收藏失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "已取消收藏"})
}

// GetFavoriteList 获取收藏列表
// @Summary 获取收藏列表
// @Description 获取用户收藏的食谱列表
// @Tags Recipe
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/recipes/favorites [get]
func (h *RecipeHandler) GetFavoriteList(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取收藏列表
	recipes, err := h.recipeService.GetUserFavorites(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取收藏列表失败"})
		return
	}

	// 转换为DTO
	recipeDTOs := make([]RecipeDTO, 0, len(recipes))
	for _, recipe := range recipes {
		dto := convertRecipeToDTO(&recipe)
		dto.IsFavorite = true // 收藏列表中的都是已收藏
		recipeDTOs = append(recipeDTOs, dto)
	}

	c.JSON(http.StatusOK, gin.H{
		"recipes": recipeDTOs,
		"count":   len(recipeDTOs),
	})
}
