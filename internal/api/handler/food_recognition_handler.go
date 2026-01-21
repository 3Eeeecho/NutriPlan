package handler

import (
	"NutriPlan/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FoodRecognitionHandler 负责处理与食物识别相关的HTTP请求。
type FoodRecognitionHandler struct {
	svc service.FoodRecognitionService
}

// NewFoodRecognitionHandler 创建一个新的食物识别处理器实例。
func NewFoodRecognitionHandler(svc service.FoodRecognitionService) *FoodRecognitionHandler {
	return &FoodRecognitionHandler{
		svc: svc,
	}
}

// RegisterRoutes 在Gin路由中注册食物识别相关的端点。
func (h *FoodRecognitionHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/recognize", h.RecognizeFood)
}

// RecognizeFood 处理菜品图片上传和识别的请求。
func (h *FoodRecognitionHandler) RecognizeFood(c *gin.Context) {
	// 从HTTP请求中获取上传的图片文件 (字段名为 "image")。
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取图片失败，请确保上传字段为 'image'"})
		return
	}

	// 检查文件类型
	contentType := file.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的图片格式，请上传 JPEG 或 PNG 格式的图片"})
		return
	}

	// 打开文件获取 io.Reader。
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法打开上传的文件"})
		return
	}
	defer f.Close()

	nutritionInfo, err := h.svc.RecognizeAndGetNutrition(c.Request.Context(), f)
	if err != nil {
		// TODO 根据 service 返回的错误类型，可以返回更具体的HTTP状态码
		c.JSON(http.StatusInternalServerError, gin.H{"error": "食物识别失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, nutritionInfo)
}
