package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	// 账号信息
	Username string   `gorm:"type:varchar(50);not null;uniqueIndex;comment:用户名" json:"username"`
	Password string   `gorm:"type:varchar(255);not null;comment:加密后的密码" json:"-"` // json:"-" 确保密码不返回给前端
	Email    *string  `gorm:"type:varchar(100);uniqueIndex;comment:邮箱地址" json:"email"`
	Role     UserRole `gorm:"type:varchar(20);not null;default:'user';index;comment:用户角色" json:"role"`

	// 基础身体数据
	Gender            string     `gorm:"type:varchar(10);comment:性别 (Male/Female)" json:"gender"`
	Age               int        `gorm:"type:bigint;default:0;comment:年龄" json:"age"`
	Height            float64    `gorm:"type:decimal(5,2);comment:身高 (cm)" json:"height"`
	Weight            float64    `gorm:"type:decimal(5,2);comment:当前体重 (kg)" json:"weight"`
	ActivityLevel     string     `gorm:"type:varchar(20);comment:活动水平" json:"activityLevel"`
	MealTimesPerDay   int        `gorm:"type:bigint;not null;default:3;comment:每日进食次数" json:"mealTimesPerDay"`
	LastProfileUpdate *time.Time `gorm:"type:datetime(3);comment:最后一次更新身体数据的时间" json:"lastProfileUpdate"`

	// 健康目标与计算结果
	HealthGoal   HealthGoal `gorm:"type:varchar(20);comment:健康目标" json:"healthGoal"`
	TargetWeight float64    `gorm:"type:decimal(5,2);comment:目标体重 (kg)" json:"targetWeight"`
	BMI          float64    `gorm:"type:decimal(4,2);comment:当前BMI指数" json:"bmi"`
	TDEE         float64    `gorm:"type:decimal(6,2);comment:每日总能量消耗" json:"tdee"`
	BMR          float64    `gorm:"type:decimal(6,2);comment:基础代谢率" json:"bmr"`

	// 饮食调节模式
	DietMode       DietMode   `gorm:"type:varchar(20);default:'normal';comment:饮食调节模式" json:"dietMode"`
	DietModeSource string     `gorm:"type:varchar(20);default:'auto';comment:模式来源(auto/manual)" json:"dietModeSource"`
	DietModeReason string     `gorm:"type:varchar(255);comment:模式原因" json:"dietModeReason"`
	DietModeUntil  *time.Time `gorm:"type:date;comment:模式截止日期" json:"dietModeUntil"`

	// 标签类数据 (存储为字符串/JSON)
	Allergies        string `gorm:"type:text;comment:过敏源" json:"allergies"`
	DietaryPrefs     string `gorm:"type:text;comment:饮食偏好" json:"dietaryPrefs"`
	HealthConditions string `gorm:"type:text;comment:健康状况" json:"healthConditions"`
}

// TableName 指定数据库表名
func (User) TableName() string {
	return "users"
}
