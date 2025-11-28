package dao

import (
	"NutriPlan/internal/repository/models"
	"fmt"

	"gorm.io/gorm"
)

// UserRepository 定义了对 User 数据的操作接口
type UserRepository interface {
	// CreateUser 创建新用户，并保存其健康档案
	CreateUser(user *models.User) error

	// GetUserByID 通过 ID 获取用户信息
	GetUserByID(id uint) (*models.User, error)

	// GetUserByUsername 通过用户名获取用户信息（用于登录）
	GetUserByUsername(username string) (*models.User, error)

	// UpdateUser 更新用户的基本信息或健康档案
	UpdateUser(user *models.User) error

	// DeleteUser 通过 ID 删除用户
	DeleteUser(id uint) error

	// UpdateHealthProfile 专门用于更新健康档案数据
	UpdateHealthProfile(id uint, profile *models.User) error
}

type GormUserRepository struct {
	DB *gorm.DB
}

// NewGormUserRepository 创建一个新的 GormUserRepository 实例
func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{
		DB: db,
	}
}

func (r *GormUserRepository) CreateUser(user *models.User) error {
	if err := r.DB.Create(user).Error; err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}
	return nil
}

func (r *GormUserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询用户(ID:%d)失败: %w", id, err)
	}
	return &user, nil
}

func (r *GormUserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询用户(Username:%s)失败: %w", username, err)
	}
	return &user, nil
}

func (r *GormUserRepository) UpdateUser(user *models.User) error {
	if err := r.DB.Save(user).Error; err != nil {
		return fmt.Errorf("更新用户失败: %w", err)
	}
	return nil
}

// 注意：这里仅更新模型中非零值（Zero Value）的字段
func (r *GormUserRepository) UpdateHealthProfile(id uint, profile *models.User) error {
	result := r.DB.Model(&models.User{}).Where("id = ?", id).Updates(map[string]any{
		"gender":             profile.Gender,
		"age":                profile.Age,
		"height":             profile.Height,
		"weight":             profile.Weight,
		"health_goal":        profile.HealthGoal,
		"target_weight":      profile.TargetWeight,
		"bmi":                profile.BMI,
		"tdee":               profile.TDEE,
		"bmr":                profile.BMR,
		"activity_level":     profile.ActivityLevel,
		"allergies":          profile.Allergies,
		"dietary_prefs":      profile.DietaryPrefs,
		"health_conditions":  profile.HealthConditions,
		"meal_times_per_day": profile.MealTimesPerDay,
	})

	if result.Error != nil {
		return fmt.Errorf("更新用户健康档案失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到 ID 为 %d 的用户或数据无变化", id)
	}

	return nil
}

func (r *GormUserRepository) DeleteUser(id uint) error {
	if err := r.DB.Delete(&models.User{}, id).Error; err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}
	return nil
}
