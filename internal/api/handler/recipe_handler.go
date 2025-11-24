package handler

import (
	"NutriPlan/internal/core/domain"
	"NutriPlan/internal/core/service"
	"NutriPlan/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RecipeHandler 食谱处理器
type RecipeHandler struct {
	recipeService service.RecipeService
	userRepo      repository.UserRepository
}

// NewRecipeHandler 创建食谱处理器实例
func NewRecipeHandler(recipeService service.RecipeService, userRepo repository.UserRepository) *RecipeHandler {
	return &RecipeHandler{
		recipeService: recipeService,
		userRepo:      userRepo,
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
	ID                 uint       `json:"id"`
	Breakfast          RecipeDTO  `json:"breakfast"`
	Lunch              RecipeDTO  `json:"lunch"`
	Dinner             RecipeDTO  `json:"dinner"`
	Snack              *RecipeDTO `json:"snack,omitempty"`
	TotalEnergy        float64    `json:"total_energy"`
	TotalProtein       float64    `json:"total_protein"`
	TotalCarbohydrate  float64    `json:"total_carbohydrate"`
	TotalFat           float64    `json:"total_fat"`
	TargetEnergy       float64    `json:"target_energy"`
	TargetProtein      float64    `json:"target_protein"`
	TargetCarbohydrate float64    `json:"target_carbohydrate"`
	TargetFat          float64    `json:"target_fat"`
	MatchScore         float64    `json:"match_score"`
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
	Ingredients  []string `json:"ingredients"`
	CookingTime  int      `json:"cooking_time"`
	Difficulty   string   `json:"difficulty"`
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

	// 获取用户信息
	user, err := h.userRepo.GetUserByID(userID.(uint))
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
		planDTO, err := h.convertPlanToDTO(plan)
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
	planDTO, err := h.convertPlanToDTO(plan)
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

	var plans []*domain.DailyRecipePlan
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

// convertPlanToDTO 转换计划为DTO（需要加载关联的食谱）
func (h *RecipeHandler) convertPlanToDTO(plan *domain.DailyRecipePlan) (DailyPlanDTO, error) {
	// 注意：这里需要预加载食谱数据
	// 实际实现中应该在repository层使用Preload
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

	// 转换食谱（这里简化处理，实际应从数据库加载）
	dto.Breakfast = convertRecipeToDTO(&plan.BreakfastRecipe)
	dto.Lunch = convertRecipeToDTO(&plan.LunchRecipe)
	dto.Dinner = convertRecipeToDTO(&plan.DinnerRecipe)

	if plan.SnackRecipe != nil {
		snackDTO := convertRecipeToDTO(plan.SnackRecipe)
		dto.Snack = &snackDTO
	}

	return dto, nil
}

// convertRecipeToDTO 转换食谱为DTO
func convertRecipeToDTO(recipe *domain.Recipe) RecipeDTO {
	// 解析食材列表
	var ingredients []string
	// 这里简化处理，实际应该解析JSON
	ingredients = []string{recipe.Ingredients}

	return RecipeDTO{
		ID:           recipe.ID,
		Name:         recipe.Name,
		Description:  recipe.Description,
		ImageURL:     recipe.ImageURL,
		MealType:     string(recipe.MealType),
		Energy:       recipe.Energy,
		Protein:      recipe.Protein,
		Carbohydrate: recipe.Carbohydrate,
		Fat:          recipe.Fat,
		Ingredients:  ingredients,
		CookingTime:  recipe.CookingTime,
		Difficulty:   recipe.Difficulty,
	}
}
