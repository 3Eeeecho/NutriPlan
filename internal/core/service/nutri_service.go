package service

import (
	"NutriPlan/internal/config"
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

	// AllocateMacros 根据目标热量和健康目标分配宏量营养素的克数
	// 返回 P_g, C_g, F_g (蛋白质克数, 碳水克数, 脂肪克数)
	AllocateMacros(targetCalorie float64, goal domain.HealthGoal) (pGram float64, cGram float64, fGram float64)
}

type NutriServiceImpl struct {
	nutriCfg config.NutrientsConfig
}

func NewNutriService() NutriService {
	return &NutriServiceImpl{
		nutriCfg: config.AppConfig.Nutrients,
	}
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

	// 根据健康目标调整 TDEE（使用配置中的值）
	switch goal {
	case domain.GoalWeightLoss:
		// 减脂：TDEE - 热量赤字
		targetCalorie = tdee - s.nutriCfg.LossCalorieDeficit
	case domain.GoalMuscleGain:
		// 增肌：TDEE + 热量盈余
		targetCalorie = tdee + s.nutriCfg.GainCalorieSurplus
	case domain.GoalSugarControl:
		// 控糖：TDEE - 轻微热量赤字
		targetCalorie = tdee - s.nutriCfg.ControlCalorieDeficit
	case domain.GoalMaintain:
		// 维持健康：保持 TDEE 不变
		targetCalorie = tdee
	default:
		targetCalorie = tdee
	}

	// 确保不低于最低安全热量
	if targetCalorie < s.nutriCfg.MinCalorie {
		return s.nutriCfg.MinCalorie
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

// getMacroRatios 根据健康目标返回 [蛋白质%, 碳水%, 脂肪%] 比例
func (s *NutriServiceImpl) getMacroRatios(goal domain.HealthGoal) (pRatio, cRatio, fRatio float64) {
	switch goal {
	case domain.GoalWeightLoss:
		// 减脂：25% 蛋白质, 40% 碳水, 35% 脂肪
		return 0.25, 0.40, 0.35
	case domain.GoalMuscleGain:
		// 增肌：30% 蛋白质, 50% 碳水, 20% 脂肪
		return 0.30, 0.50, 0.20
	case domain.GoalSugarControl:
		// 控糖：25% 蛋白质, 35% 碳水, 40% 脂肪
		return 0.25, 0.35, 0.40
	case domain.GoalMaintain:
		// 维持健康：20% 蛋白质, 55% 碳水, 25% 脂肪
		return 0.20, 0.55, 0.25
	default:
		// 默认使用维持健康的比例
		return 0.20, 0.55, 0.25
	}
}

// AllocateMacros 根据目标热量和健康目标分配宏量营养素的克数
// 返回 P_g, C_g, F_g (蛋白质克数, 碳水克数, 脂肪克数)
func (s *NutriServiceImpl) AllocateMacros(targetCalorie float64, goal domain.HealthGoal) (pGram float64, cGram float64, fGram float64) {
	// 1. 获取宏量营养素比例
	pRatio, cRatio, fRatio := s.getMacroRatios(goal)

	// 2. 计算每种营养素的总热量 (kcal)
	pCalorie := targetCalorie * pRatio
	cCalorie := targetCalorie * cRatio
	fCalorie := targetCalorie * fRatio

	// 3. 转换为克数 (g)
	// 使用配置中的热量因子
	pGram = pCalorie / s.nutriCfg.ProteinFactor
	cGram = cCalorie / s.nutriCfg.CarbFactor
	fGram = fCalorie / s.nutriCfg.FatFactor

	// 返回四舍五入的克数
	return math.Round(pGram), math.Round(cGram), math.Round(fGram)
}
