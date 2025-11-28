package models

import "gorm.io/gorm"

// 用户体重记录表
type UserWeightLog struct {
	gorm.Model
	UserID uint    `gorm:"not null;index;comment:关联用户ID" json:"userId"`
	Weight float64 `gorm:"type:decimal(5,2);not null;comment:记录时的体重 (kg)" json:"weight"`
	BMI    float64 `gorm:"type:decimal(4,2);comment:记录时的BMI" json:"bmi"`
}

func (UserWeightLog) TableName() string {
	return "user_weight_logs"
}
