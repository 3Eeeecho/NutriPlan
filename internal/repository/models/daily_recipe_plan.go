package models

import (
	"time"

	"gorm.io/gorm"
)

type DailyRecipePlan struct {
	gorm.Model
	UserID     uint      `gorm:"not null;index;comment:关联用户ID" json:"userId"`
	PlanDate   time.Time `gorm:"type:date;not null;comment:计划日期" json:"planDate"` // 注意这里的 type:date
	PlanName   string    `gorm:"type:varchar(200);comment:计划名称" json:"planName"`
	IsSelected bool      `gorm:"type:tinyint(1);default:0;comment:用户是否采纳" json:"isSelected"`

	// 关联的食谱 ID
	BreakfastRecipeID uint `gorm:"not null;comment:早餐食谱ID" json:"breakfastRecipeId"`
	LunchRecipeID     uint `gorm:"not null;comment:午餐食谱ID" json:"lunchRecipeId"`
	DinnerRecipeID    uint `gorm:"not null;comment:晚餐食谱ID" json:"dinnerRecipeId"`
	SnackRecipeID     uint `gorm:"default:null;comment:加餐食谱ID" json:"snackRecipeId"`

	// 每餐组合明细（持久化为JSON数组，兼容“每餐多食物”）
	BreakfastItemIDs []uint `gorm:"serializer:json;type:longtext;comment:早餐组合食谱ID列表" json:"breakfastItemIds"`
	LunchItemIDs     []uint `gorm:"serializer:json;type:longtext;comment:午餐组合食谱ID列表" json:"lunchItemIds"`
	DinnerItemIDs    []uint `gorm:"serializer:json;type:longtext;comment:晚餐组合食谱ID列表" json:"dinnerItemIds"`
	SnackItemIDs     []uint `gorm:"serializer:json;type:longtext;comment:加餐组合食谱ID列表" json:"snackItemIds"`

	// 计划总营养数据
	TotalEnergy       float64 `gorm:"type:decimal(10,2);comment:计划总热量" json:"totalEnergy"`
	TotalProtein      float64 `gorm:"type:decimal(10,2);comment:计划总蛋白" json:"totalProtein"`
	TotalCarbohydrate float64 `gorm:"type:decimal(10,2);comment:计划总碳水" json:"totalCarbohydrate"`
	TotalFat          float64 `gorm:"type:decimal(10,2);comment:计划总脂肪" json:"totalFat"`

	// 用户目标对比
	TargetEnergy       float64 `gorm:"type:decimal(10,2);comment:目标热量" json:"targetEnergy"`
	TargetProtein      float64 `gorm:"type:decimal(10,2);comment:目标蛋白" json:"targetProtein"`
	TargetCarbohydrate float64 `gorm:"type:decimal(10,2);comment:目标碳水" json:"targetCarbohydrate"`
	TargetFat          float64 `gorm:"type:decimal(10,2);comment:目标脂肪" json:"targetFat"`

	MatchScore float64 `gorm:"type:decimal(5,2);comment:匹配度评分" json:"matchScore"`

	// 关联关系 (Preload 用)
	BreakfastRecipe Recipe `gorm:"foreignKey:BreakfastRecipeID" json:"breakfastRecipe"`
	LunchRecipe     Recipe `gorm:"foreignKey:LunchRecipeID" json:"lunchRecipe"`
	DinnerRecipe    Recipe `gorm:"foreignKey:DinnerRecipeID" json:"dinnerRecipe"`
	SnackRecipe     Recipe `gorm:"foreignKey:SnackRecipeID" json:"snackRecipe"`

	// 推荐阶段扩展字段（不落库）
	BreakfastItems []Recipe `gorm:"-" json:"breakfastItems,omitempty"`
	LunchItems     []Recipe `gorm:"-" json:"lunchItems,omitempty"`
	DinnerItems    []Recipe `gorm:"-" json:"dinnerItems,omitempty"`
	SnackItems     []Recipe `gorm:"-" json:"snackItems,omitempty"`
}

func (DailyRecipePlan) TableName() string {
	return "daily_recipe_plans"
}
