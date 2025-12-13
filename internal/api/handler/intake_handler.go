package handler

import (
	"NutriPlan/internal/repository/models"
	"NutriPlan/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// IntakeHandler 饮食记录处理器
type IntakeHandler struct {
	intakeService service.IntakeService
}

// NewIntakeHandler 创建饮食记录处理器实例
func NewIntakeHandler(intakeService service.IntakeService) *IntakeHandler {
	return &IntakeHandler{
		intakeService: intakeService,
	}
}

// AddIntakeRecordRequest 添加饮食记录请求
type AddIntakeRecordRequest struct {
	MealType          string  `json:"meal_type" binding:"required"`          // 餐点类型: breakfast/lunch/dinner/snack
	FoodSource        int     `json:"food_source"`                           // 1=食谱, 2=基础食物，默认2
	SourceID          uint    `json:"source_id"`                             // 食谱ID或食物ID（可选）
	FoodName          string  `json:"food_name" binding:"required"`          // 食物名称
	IntakeAmount      float64 `json:"intake_amount" binding:"required,gt=0"` // 摄入量
	CalculatedEnergy  float64 `json:"calculated_energy" `                    // 计算后的热量
	CalculatedProtein float64 `json:"calculated_protein"`                    // 计算后的蛋白质
	CalculatedCarb    float64 `json:"calculated_carb"`                       // 计算后的碳水
	CalculatedFat     float64 `json:"calculated_fat"`                        // 计算后的脂肪
}

// AddIntakeRecord 添加饮食记录
// @Summary 添加饮食记录
// @Description 用户手动记录当日已吃食物
// @Tags Intake
// @Accept json
// @Produce json
// @Param record body AddIntakeRecordRequest true "饮食记录"
// @Success 200 {object} map[string]interface{}
// @Router /api/intake/records [post]
func (h *IntakeHandler) AddIntakeRecord(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var req AddIntakeRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 如果未指定food_source，默认为2（基础食物）
	if req.FoodSource == 0 {
		req.FoodSource = 2
	}

	// 构建记录模型
	record := &models.DailyIntakeRecord{
		UserID:            userID.(uint),
		MealType:          req.MealType,
		FoodSource:        req.FoodSource,
		SourceID:          req.SourceID,
		FoodName:          req.FoodName,
		IntakeAmount:      req.IntakeAmount,
		CalculatedEnergy:  req.CalculatedEnergy,
		CalculatedProtein: req.CalculatedProtein,
		CalculatedCarb:    req.CalculatedCarb,
		CalculatedFat:     req.CalculatedFat,
	}

	// 添加记录
	if err := h.intakeService.AddIntakeRecord(record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "记录添加成功",
		"record":  record,
	})
}

// DeleteIntakeRecord 删除饮食记录
// @Summary 删除饮食记录
// @Description 删除某条饮食记录
// @Tags Intake
// @Accept json
// @Produce json
// @Param id path int true "记录ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/intake/records/{id} [delete]
func (h *IntakeHandler) DeleteIntakeRecord(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 解析记录ID
	recordIDStr := c.Param("id")
	recordID, err := strconv.ParseUint(recordIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录ID"})
		return
	}

	// 删除记录
	if err := h.intakeService.DeleteIntakeRecord(userID.(uint), uint(recordID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "记录删除成功"})
}

// GetTodayStatus 获取当日营养状态
// @Summary 获取当日营养状态
// @Description 获取当日饮食记录和营养达标率
// @Tags Intake
// @Accept json
// @Produce json
// @Success 200 {object} service.NutritionStatus
// @Router /api/intake/today [get]
func (h *IntakeHandler) GetTodayStatus(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取当日营养状态
	status, err := h.intakeService.GetTodayNutritionStatus(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取营养状态失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// GetWeeklyReport 获取周报告
// @Summary 获取周报告
// @Description 获取本周饮食趋势报告
// @Tags Intake
// @Accept json
// @Produce json
// @Success 200 {object} service.WeeklyReport
// @Router /api/intake/weekly [get]
func (h *IntakeHandler) GetWeeklyReport(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取周报告
	report, err := h.intakeService.GetWeeklyReport(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取周报告失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}
