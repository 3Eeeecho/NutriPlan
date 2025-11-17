package service

import (
	"NutriPlan/internal/core/domain"
	"math"
	"strings"
)

type NutriService interface {
	// CalculateTDEE 基于用户档案计算 BMR, TDEE, 和 BMI
	// 返回计算结果 BMR, TDEE, BMI
	CalculateTDEE(user *domain.User) (float64, float64, float64)

	// DetermineTargetCalorie 根据用户的 TDEE 和 HealthGoal 确定每日目标热量摄入
	DetermineTargetCalorie(tdee float64, goal domain.HealthGoal) float64

	// 其他宏量营养素分配、营养素缺口分析等方法
}

type NutriServiceImpl struct{}

func NewNutriService() NutriService {
	return &NutriServiceImpl{}
}

func (s *NutriServiceImpl) CalculateTDEE(user *domain.User) (bmr float64, tdee float64, bmi float64) {
	// 计算 BMR（基础代谢率）
	heightM := user.Height / 100.0
	if heightM > 0 {
		bmi = user.Weight / (heightM * heightM)
	}

	baseBMR := (10 * user.Weight) + (6.25 * user.Height) - (5 * float64(user.Age))

	if strings.ToLower(user.Gender) == "male" || user.Gender == "男" {
		bmr = baseBMR + 5
	} else if strings.ToLower(user.Gender) == "female" || user.Gender == "女" {
		bmr = baseBMR - 161
	}

	// 确保 BMR 不为负数
	if bmr < 0 {
		bmr = 0
	}

	// 3. 计算 TDEE (每日总消耗热量)
	factor := s.getActivityFactor(user.ActivityLevel)
	tdee = bmr * factor

	// 四舍五入到最近的整数（通常热量值以整数呈现）
	return math.Round(bmr), math.Round(tdee), math.Round(bmi*100) / 100 // BMI 保留两位小数
}

func (s *NutriServiceImpl) DetermineTargetCalorie(tdee float64, goal domain.HealthGoal) float64 {
	var targetCalorie float64

	// 根据健康目标调整 TDEE
	switch goal {
	case domain.GoalWeightLoss:
		// 减脂
		targetCalorie = tdee - 400
	case domain.GoalMuscleGain:
		// 增肌
		targetCalorie = tdee + 300
	case domain.GoalSugarControl:
		// 控糖：通常是轻微的热量赤字或维持，但宏量营养素比例要求严格
		targetCalorie = tdee - 200
	case domain.GoalMaintain:
		// 维持健康：保持 TDEE 不变
		targetCalorie = tdee
	default:
		targetCalorie = tdee
	}

	if targetCalorie < 1200 {
		return 1200 // 设定一个安全底线
	}

	return math.Round(targetCalorie)
}

// getActivityFactor 根据活动水平字符串返回 TDEE 乘数
func (s *NutriServiceImpl) getActivityFactor(level string) float64 {
	standardLevel := strings.TrimSpace(level)
	switch standardLevel {
	case "久坐":
		return 1.2
	case "轻度活动":
		return 1.375
	case "中度活动":
		return 1.55
	case "重度活动":
		return 1.725
	case "极重度活动":
		return 1.9
	default:
		// 默认值，防止因输入错误导致 TDEE 异常
		return 1.2
	}
}
