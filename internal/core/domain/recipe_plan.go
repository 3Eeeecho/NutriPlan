package domain

import (
	"gorm.io/gorm"
)

// DailyRecipePlan 每日食谱计划
type DailyRecipePlan struct {
	gorm.Model

	// 用户关联
	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID"`

	// 计划名称
	PlanName string `gorm:"type:varchar(200)" json:"plan_name"`

	// 各餐次的食谱ID
	BreakfastRecipeID uint   `gorm:"not null" json:"breakfast_recipe_id"`
	BreakfastRecipe   Recipe `gorm:"foreignKey:BreakfastRecipeID"`

	LunchRecipeID uint   `gorm:"not null" json:"lunch_recipe_id"`
	LunchRecipe   Recipe `gorm:"foreignKey:LunchRecipeID"`

	DinnerRecipeID uint   `gorm:"not null" json:"dinner_recipe_id"`
	DinnerRecipe   Recipe `gorm:"foreignKey:DinnerRecipeID"`

	SnackRecipeID *uint   `gorm:"" json:"snack_recipe_id"` // 加餐可选
	SnackRecipe   *Recipe `gorm:"foreignKey:SnackRecipeID"`

	// 总营养成分汇总
	TotalEnergy       float64 `gorm:"type:decimal(10,2)" json:"total_energy"`
	TotalProtein      float64 `gorm:"type:decimal(10,2)" json:"total_protein"`
	TotalCarbohydrate float64 `gorm:"type:decimal(10,2)" json:"total_carbohydrate"`
	TotalFat          float64 `gorm:"type:decimal(10,2)" json:"total_fat"`

	// 用户目标营养
	TargetEnergy       float64 `gorm:"type:decimal(10,2)" json:"target_energy"`
	TargetProtein      float64 `gorm:"type:decimal(10,2)" json:"target_protein"`
	TargetCarbohydrate float64 `gorm:"type:decimal(10,2)" json:"target_carbohydrate"`
	TargetFat          float64 `gorm:"type:decimal(10,2)" json:"target_fat"`

	// 匹配度评分 (0-100)
	MatchScore float64 `gorm:"type:decimal(5,2)" json:"match_score"`

	// 是否被用户选择
	IsSelected bool `gorm:"default:false" json:"is_selected"`
}

// TableName 指定表名
func (DailyRecipePlan) TableName() string {
	return "daily_recipe_plans"
}
