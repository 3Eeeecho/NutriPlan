package handler

import (
	"NutriPlan/internal/repository/models"
	"NutriPlan/internal/service"
	"NutriPlan/pkg/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler 处理用户相关的HTTP请求
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 创建UserHandler实例
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterRequest 注册请求结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email"`
}

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应结构体
type LoginResponse struct {
	Token   string       `json:"token"`
	User    *models.User `json:"user"`
	Message string       `json:"message"`
}

// UpdateProfileRequest 更新档案请求结构体
type UpdateProfileRequest struct {
	Gender           string  `json:"gender"`
	Age              int     `json:"age" binding:"min=1,max=150"`
	Height           float64 `json:"height" binding:"min=50,max=300"`
	Weight           float64 `json:"weight" binding:"min=20,max=500"`
	HealthGoal       string  `json:"health_goal"`
	TargetWeight     float64 `json:"target_weight"`
	ActivityLevel    string  `json:"activity_level"`
	Allergies        string  `json:"allergies"`
	DietaryPrefs     string  `json:"dietary_prefs"`
	HealthConditions string  `json:"health_conditions"`
	MealTimesPerDay  int     `json:"meal_times_per_day" binding:"min=1,max=6"`
}

// Register 处理用户注册
// @Summary 用户注册
// @Description 注册新用户账号
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册信息"
// @Success 201 {object} map[string]interface{} "注册成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 409 {object} map[string]interface{} "用户名已存在"
// @Router /api/v1/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 创建用户对象，只设置基本注册信息
	// 身体数据等健康档案信息后续通过更新档案接口填写
	user := &models.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		// 设置默认值，避免数据库约束问题
		MealTimesPerDay: 3,
	}

	// 调用服务层注册用户
	if err := h.userService.RegisterUser(user); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "用户名已存在" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"user_id": user.ID,
	})
}

// Login 处理用户登录
// @Summary 用户登录
// @Description 用户登录并获取JWT token
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} LoginResponse "登录成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "用户名或密码错误"
// @Router /api/v1/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 调用服务层验证登录
	user, err := h.userService.LoginUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 生成JWT token
	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成token失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token:   token,
		User:    user,
		Message: "登录成功",
	})
}

// UpdateProfile 处理用户档案更新
// @Summary 更新用户健康档案
// @Description 更新当前登录用户的健康档案信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProfileRequest true "更新信息"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 404 {object} map[string]interface{} "用户不存在"
// @Router /api/v1/user/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// 从JWT中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授权访问",
		})
		return
	}

	// 类型断言
	id, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用户ID格式错误",
		})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 创建更新对象
	profile := &models.User{
		Gender:           req.Gender,
		Age:              req.Age,
		Height:           req.Height,
		Weight:           req.Weight,
		HealthGoal:       models.HealthGoal(req.HealthGoal),
		TargetWeight:     req.TargetWeight,
		ActivityLevel:    req.ActivityLevel,
		Allergies:        req.Allergies,
		DietaryPrefs:     req.DietaryPrefs,
		HealthConditions: req.HealthConditions,
		MealTimesPerDay:  req.MealTimesPerDay,
	}

	// 调用服务层更新档案
	if err := h.userService.UpdateUserProfile(id, profile); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "用户不存在" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 获取更新后的用户信息
	updatedUser, err := h.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取用户信息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"user":    updatedUser,
	})
}

// GetProfile 获取用户档案
// @Summary 获取用户健康档案
// @Description 获取当前登录用户的健康档案信息
// @Tags 用户
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User "用户档案"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 404 {object} map[string]interface{} "用户不存在"
// @Router /api/v1/user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	// 从JWT中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授权访问",
		})
		return
	}

	// 类型断言
	id, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用户ID格式错误",
		})
		return
	}

	// 调用服务层获取用户信息
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取用户信息失败: " + err.Error(),
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// NutritionRequirementsResponse 营养需求响应结构体
type NutritionRequirementsResponse struct {
	TargetCalorie float64 `json:"target_calorie"` // 目标热量 (kcal)
	ProteinGram   float64 `json:"protein_gram"`   // 蛋白质 (g)
	CarbGram      float64 `json:"carb_gram"`      // 碳水化合物 (g)
	FatGram       float64 `json:"fat_gram"`       // 脂肪 (g)
	ProteinRatio  float64 `json:"protein_ratio"`  // 蛋白质比例 (%)
	CarbRatio     float64 `json:"carb_ratio"`     // 碳水化合物比例 (%)
	FatRatio      float64 `json:"fat_ratio"`      // 脂肪比例 (%)
}

// GetNutritionRequirements 获取用户营养需求
// @Summary 获取营养需求
// @Description 根据用户档案计算目标热量和宏量营养素分配
// @Tags 用户
// @Produce json
// @Security BearerAuth
// @Success 200 {object} NutritionRequirementsResponse "营养需求"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 404 {object} map[string]interface{} "用户不存在或档案不完整"
// @Router /api/v1/user/nutrition [get]
func (h *UserHandler) GetNutritionRequirements(c *gin.Context) {
	// 从JWT中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授权访问",
		})
		return
	}

	// 类型断言
	id, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用户ID格式错误",
		})
		return
	}

	// 获取用户信息
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取用户信息失败: " + err.Error(),
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "用户不存在",
		})
		return
	}

	// 获取营养需求
	targetCalorie, proteinGram, carbGram, fatGram, proteinRatio, carbRatio, fatRatio, err := h.userService.GetNutritionRequirements(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, NutritionRequirementsResponse{
		TargetCalorie: targetCalorie,
		ProteinGram:   proteinGram,
		CarbGram:      carbGram,
		FatGram:       fatGram,
		ProteinRatio:  proteinRatio,
		CarbRatio:     carbRatio,
		FatRatio:      fatRatio,
	})
}
