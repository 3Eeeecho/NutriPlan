package service

import (
	"NutriPlan/internal/client"
	"context"
	"io"
)

// FoodRecognitionService 定义了食物识别与记录服务的接口。
type FoodRecognitionService interface {
	// RecognizeAndGetNutrition 通过图片识别菜品并获取其营养信息。
	// imageReader: 包含待识别菜品图片数据。
	// 返回一个包含识别结果和营养信息的结构体。
	RecognizeAndGetNutrition(ctx context.Context, imageReader io.Reader) (*FoodNutritionInfo, error)

	// AnalyzeFoodText 根据用户提供的文本描述分析食物并获取其营养信息。
	// text: 用户输入的食物描述。
	// 返回一个包含分析结果和营养信息的结构体。
	AnalyzeFoodText(ctx context.Context, text string) (*FoodNutritionInfo, error)
}

// foodRecognitionServiceImpl 是 FoodRecognitionService 的具体实现。
type foodRecognitionServiceImpl struct {
	zhipuClient client.ZhipuAIClient
}

// NewFoodRecognitionService 创建一个新的食物识别服务实例。
func NewFoodRecognitionService(zhipuClient client.ZhipuAIClient) FoodRecognitionService {
	return &foodRecognitionServiceImpl{
		zhipuClient: zhipuClient,
	}
}

// FoodNutritionInfo 封装了从AI模型返回并经过处理的菜品营养信息。
// 这个结构体旨在直接服务于前端展示。
type FoodNutritionInfo struct {
	DishName        string `json:"dish_name"`        // 菜品名称
	Reasoning       string `json:"reasoning"`        // AI分析理由
	EstimatedWeight int    `json:"estimated_weight"` // 整份食物的估算重量 (克)

	// --- 每100克标准营养值 ---
	CaloriesPer100g float64 `json:"calories_per_100g"` // 每100g热量 (千卡)
	ProteinPer100g  float64 `json:"protein_per_100g"`  // 每100g蛋白质 (克)
	FatPer100g      float64 `json:"fat_per_100g"`      // 每100g脂肪 (克)
	CarbsPer100g    float64 `json:"carbs_per_100g"`    // 每100g碳水化合物 (克)

	// --- 根据估算重量计算的总营养值 ---
	TotalCalories float64 `json:"total_calories"` // 整份总热量 (千卡)
	TotalProtein  float64 `json:"total_protein"`  // 整份总蛋白质 (克)
	TotalFat      float64 `json:"total_fat"`      // 整份总脂肪 (克)
	TotalCarbs    float64 `json:"total_carbs"`    // 整份总碳水化合物 (克)
}

// RecognizeAndGetNutrition 通过图片调用智谱模型识别菜品，处理后返回详细营养信息。
func (s *foodRecognitionServiceImpl) RecognizeAndGetNutrition(ctx context.Context, imageReader io.Reader) (*FoodNutritionInfo, error) {
	// 1. 调用智谱客户端识别食物，获取基础营养数据。
	zhipuResp, err := s.zhipuClient.RecognizeFood(ctx, imageReader)
	if err != nil {
		return nil, err
	}

	// 2. 计算总营养成分。
	// weightFactor = 估算重量 / 100g
	weightFactor := float64(zhipuResp.EstimatedWeight) / 100.0

	totalCalories := zhipuResp.Calories100g * weightFactor
	totalProtein := zhipuResp.Protein100g * weightFactor
	totalFat := zhipuResp.Fat100g * weightFactor
	totalCarbs := zhipuResp.Carbs100g * weightFactor

	// 3. 将数据映射到最终的服务层结构体中。
	info := &FoodNutritionInfo{
		DishName:        zhipuResp.Name,
		Reasoning:       zhipuResp.Reasoning,
		EstimatedWeight: zhipuResp.EstimatedWeight,
		CaloriesPer100g: zhipuResp.Calories100g,
		ProteinPer100g:  zhipuResp.Protein100g,
		FatPer100g:      zhipuResp.Fat100g,
		CarbsPer100g:    zhipuResp.Carbs100g,
		TotalCalories:   totalCalories,
		TotalProtein:    totalProtein,
		TotalFat:        totalFat,
		TotalCarbs:      totalCarbs,
	}

	return info, nil
}

// AnalyzeFoodText 根据用户提供的文本描述分析食物，处理后返回详细营养信息。
func (s *foodRecognitionServiceImpl) AnalyzeFoodText(ctx context.Context, text string) (*FoodNutritionInfo, error) {
	// 1. 调用智谱客户端分析食物文本，获取基础营养数据。
	zhipuResp, err := s.zhipuClient.AnalyzeFoodText(ctx, text)
	if err != nil {
		return nil, err
	}

	// 2. 计算总营养成分。
	// weightFactor = 估算重量 / 100g
	weightFactor := float64(zhipuResp.EstimatedWeight) / 100.0

	totalCalories := zhipuResp.Calories100g * weightFactor
	totalProtein := zhipuResp.Protein100g * weightFactor
	totalFat := zhipuResp.Fat100g * weightFactor
	totalCarbs := zhipuResp.Carbs100g * weightFactor

	// 3. 将数据映射到最终的服务层结构体中。
	info := &FoodNutritionInfo{
		DishName:        zhipuResp.Name,
		Reasoning:       zhipuResp.Reasoning,
		EstimatedWeight: zhipuResp.EstimatedWeight,
		CaloriesPer100g: zhipuResp.Calories100g,
		ProteinPer100g:  zhipuResp.Protein100g,
		FatPer100g:      zhipuResp.Fat100g,
		CarbsPer100g:    zhipuResp.Carbs100g,
		TotalCalories:   totalCalories,
		TotalProtein:    totalProtein,
		TotalFat:        totalFat,
		TotalCarbs:      totalCarbs,
	}

	return info, nil
}
