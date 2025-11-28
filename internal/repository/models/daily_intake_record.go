package models

import (
	"time"

	"gorm.io/gorm"
)

// 每日饮食打卡表 - 真实记录
type DailyIntakeRecord struct {
	gorm.Model
	UserID     uint      `gorm:"not null;index;comment:关联用户ID" json:"userId"`
	RecordDate time.Time `gorm:"type:date;not null;index;comment:记录日期" json:"recordDate"`
	MealType   string    `gorm:"type:varchar(20);not null;comment:餐点类型" json:"mealType"`

	// 核心逻辑: 区分是食谱还是单一食物
	FoodSource int    `gorm:"type:tinyint;not null;comment:1=食谱, 2=基础食物" json:"foodSource"`
	SourceID   uint   `gorm:"not null;comment:Recipes.ID 或 FoodNutrition.ID" json:"sourceId"`
	FoodName   string `gorm:"type:varchar(200);not null;comment:食物名称备份" json:"foodName"`

	// 实际摄入量
	IntakeAmount      float64 `gorm:"type:decimal(8,2);not null;comment:摄入数量(g/份)" json:"intakeAmount"`
	CalculatedEnergy  float64 `gorm:"type:decimal(10,2);not null;comment:实际热量" json:"calculatedEnergy"`
	CalculatedProtein float64 `gorm:"type:decimal(10,2);default:0;comment:实际蛋白" json:"calculatedProtein"`
	CalculatedCarb    float64 `gorm:"type:decimal(10,2);default:0;comment:实际碳水" json:"calculatedCarb"`
	CalculatedFat     float64 `gorm:"type:decimal(10,2);default:0;comment:实际脂肪" json:"calculatedFat"`
}

func (DailyIntakeRecord) TableName() string {
	return "daily_intake_records"
}
