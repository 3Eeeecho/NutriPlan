package models

import (
	"gorm.io/gorm"
)

// MealType 定义餐次类型
type MealType string

const (
	MealTypeBreakfast MealType = "早餐"
	MealTypeLunch     MealType = "午餐"
	MealTypeDinner    MealType = "晚餐"
	MealTypeSnack     MealType = "加餐"
)

// Recipe 食谱结构体
type Recipe struct {
	gorm.Model

	// 基本信息
	Name        string   `gorm:"type:varchar(200);not null;index" json:"name"`     // 食谱名称
	Description string   `gorm:"type:text" json:"description"`                     // 食谱描述
	ImageURL    string   `gorm:"type:varchar(500)" json:"image_url"`               // 食谱图片URL
	MealType    MealType `gorm:"type:varchar(20);not null;index" json:"meal_type"` // 餐次类型

	// 营养成分 (每份)
	Energy       float64 `gorm:"type:decimal(10,2);not null" json:"energy"`       // 能量 (kcal)
	Protein      float64 `gorm:"type:decimal(10,2);not null" json:"protein"`      // 蛋白质 (g)
	Carbohydrate float64 `gorm:"type:decimal(10,2);not null" json:"carbohydrate"` // 碳水化合物 (g)
	Fat          float64 `gorm:"type:decimal(10,2);not null" json:"fat"`          // 脂肪 (g)
	DietaryFiber float64 `gorm:"type:decimal(10,2)" json:"dietary_fiber"`         // 膳食纤维 (g)

	// 食材列表 (JSON格式存储，例如: ["鸡胸肉", "西兰花", "糙米"])
	Ingredients string `gorm:"type:text;not null" json:"ingredients"`

	// 烹饪步骤 (可选)
	CookingSteps string `gorm:"type:text" json:"cooking_steps"`

	// 烹饪时间（分钟）
	CookingTime int `gorm:"default:0" json:"cooking_time"`

	// 难度等级 (简单/中等/困难)
	Difficulty string `gorm:"type:varchar(20)" json:"difficulty"`
}

// TableName 指定表名
func (Recipe) TableName() string {
	return "recipes"
}
