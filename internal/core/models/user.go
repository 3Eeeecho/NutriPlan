package models

import (
	"time"

	"gorm.io/gorm"
)

// HealthGoal 定义用户的健康目标
type HealthGoal string

// 定义健康目标常量，用于提高代码可读性和类型安全
const (
	GoalWeightLoss   HealthGoal = "减脂"
	GoalMuscleGain   HealthGoal = "增肌"
	GoalSugarControl HealthGoal = "控糖"
	GoalMaintain     HealthGoal = "维持健康"
)

// User 结构体定义了系统中的用户及他们的健康档案
type User struct {
	// GORM 默认字段
	gorm.Model

	// 基础认证信息
	Username string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"` // 用户名（唯一索引）
	Password string `gorm:"type:varchar(255);not null" json:"-"`                   // 密码（存储加密后的哈希值）
	Email    string `gorm:"type:varchar(100);uniqueIndex" json:"email"`            // 邮箱（唯一索引，可选）

	// --- 基础生理数据 (用于 BMR/TDEE 计算) ---
	Gender string  `gorm:"type:varchar(10)" json:"gender"`  // 性别 (Male/Female)
	Age    int     `gorm:"default:0" json:"age"`            // 年龄 (Years)
	Height float64 `gorm:"type:decimal(5,2)" json:"height"` // 身高 (cm)
	Weight float64 `gorm:"type:decimal(5,2)" json:"weight"` // 当前体重 (kg)

	// --- 健康与目标数据 ---
	HealthGoal   HealthGoal `gorm:"type:varchar(20)" json:"health_goal"`    // 健康目标 (减脂/增肌/控糖等)
	TargetWeight float64    `gorm:"type:decimal(5,2)" json:"target_weight"` // 目标体重 (kg)

	// BMI 和 TDEE 为计算结果，可以缓存到数据库，或在Service层实时计算
	BMI  float64 `gorm:"type:decimal(4,2)" json:"bmi"`  // 身体质量指数
	TDEE float64 `gorm:"type:decimal(6,2)" json:"tdee"` // 每日总能量消耗 (kcal/day)
	BMR  float64 `gorm:"type:decimal(6,2)" json:"bmr"`  // 基础代谢率 (kcal/day)

	// 活动水平（用于TDEE计算）
	ActivityLevel string `gorm:"type:varchar(20)" json:"activity_level"` // 活动水平 (如: 久坐, 轻度活动)

	// --- 个性化偏好与禁忌 ---
	// 存储过敏、忌口和特殊健康问题的 JSON 字符串或用分隔符连接的字符串
	Allergies        string `gorm:"type:text" json:"allergies"`         // 过敏源（如: 海鲜、花生）
	DietaryPrefs     string `gorm:"type:text" json:"dietary_prefs"`     // 饮食偏好（如: 喜辣、低脂）
	HealthConditions string `gorm:"type:text" json:"health_conditions"` // 基础疾病（如: 糖尿病、乳糖不耐）

	// 每天用餐次数，用于食谱分配
	MealTimesPerDay int `gorm:"not null;default:3" json:"meal_times_per_day"`

	// 记录最后一次更新健康档案的时间
	LastProfileUpdate time.Time `gorm:"autoUpdateTime" json:"last_profile_update"`
}
