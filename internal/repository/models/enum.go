package models

// HealthGoal 定义用户的健康目标
type HealthGoal string

// 定义健康目标常量，用于提高代码可读性和类型安全
const (
	GoalWeightLoss   HealthGoal = "减脂"
	GoalMuscleGain   HealthGoal = "增肌"
	GoalSugarControl HealthGoal = "控糖"
	GoalMaintain     HealthGoal = "维持健康"
)

// MealType 定义枚举，防止字符串拼写错误
type MealType string

const (
	MealTypeBreakfast MealType = "早餐"
	MealTypeLunch     MealType = "午餐"
	MealTypeDinner    MealType = "晚餐"
	MealTypeSnack     MealType = "加餐"
)

// DietMode 饮食调节模式
type DietMode string

const (
	DietModeNormal      DietMode = "normal"
	DietModeLightAdjust DietMode = "light_adjust"
	DietModeBland       DietMode = "bland"
	DietModeHeavyAdjust DietMode = "heavy_adjust"
)

// --- 2. 禁忌/过敏标签 (对应 Recipe.ForbiddenUsers) ---
const (
	// 常见过敏源
	TagSeafood = "海鲜"  // 对应食材：虾、蟹、贝类
	TagNut     = "坚果"  // 对应食材：花生、核桃、腰果
	TagDairy   = "乳制品" // 对应食材：牛奶、芝士、奶油
	TagGluten  = "麸质"  // 对应食材：面粉、面包 (针对乳糜泻)
	TagEgg     = "蛋类"  // 对应食材：鸡蛋

	// 饮食偏好/禁忌
	TagPork       = "猪肉"  // 针对不吃猪肉人群
	TagBeef       = "牛肉"  // 针对不吃牛肉人群
	TagSpicy      = "辛辣"  // 针对不吃辣人群
	TagHighSugar  = "高糖"  // 针对糖尿病/抗糖人群
	TagHighPurine = "高嘌呤" // 针对痛风人群
)
