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

type Recipe struct {
	gorm.Model
	Name        string   `gorm:"type:varchar(200);not null;index;comment:食谱名称" json:"name"`
	Description string   `gorm:"type:text;comment:食谱简介" json:"description"`
	ImageURL    string   `gorm:"type:varchar(500);comment:封面图片链接" json:"imageUrl"`
	MealType    MealType `gorm:"type:varchar(20);not null;index;comment:适合餐点类型" json:"mealType"`
	Difficulty  string   `gorm:"type:varchar(20);comment:烹饪难度" json:"difficulty"`
	CookingTime int      `gorm:"type:bigint;default:0;comment:烹饪耗时(分)" json:"cookingTime"`

	// 核心营养素
	Energy       float64 `gorm:"type:decimal(10,2);not null;comment:总热量 (kcal)" json:"energy"`
	Protein      float64 `gorm:"type:decimal(10,2);not null;comment:蛋白质 (g)" json:"protein"`
	Carbohydrate float64 `gorm:"type:decimal(10,2);not null;comment:碳水化合物 (g)" json:"carbohydrate"`
	Fat          float64 `gorm:"type:decimal(10,2);not null;comment:脂肪 (g)" json:"fat"`
	DietaryFiber float64 `gorm:"type:decimal(10,2);comment:膳食纤维 (g)" json:"dietaryFiber"`

	// 长文本字段
	Ingredients  string `gorm:"type:text;not null;comment:所需食材清单(JSON)" json:"ingredients"`
	CookingSteps string `gorm:"type:text;comment:烹饪步骤" json:"cookingSteps"`
}

func (Recipe) TableName() string {
	return "recipes"
}
