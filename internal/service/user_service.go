package service

import (
	"NutriPlan/internal/repository/dao"
	"NutriPlan/internal/repository/models"
	"NutriPlan/pkg/xerr"
	"fmt"
	"time"
)

// UserService 定义了用户业务流程的接口
type UserService interface {
	// RegisterUser 处理用户注册逻辑，包括密码加密和创建数据库记录
	RegisterUser(user *models.User) error

	// LoginUser 处理用户登录逻辑，验证密码等
	LoginUser(username, password string) (*models.User, error)

	// UpdateUserProfile 处理用户健康档案更新，包括触发 TDEE 计算
	UpdateUserProfile(id uint, profile *models.User) error

	// GetUserByID 获取用户信息
	GetUserByID(id uint) (*models.User, error)

	// GetNutritionRequirements 获取用户营养需求
	GetNutritionRequirements(user *models.User) (targetCalorie, proteinGram, carbGram, fatGram, proteinRatio, carbRatio, fatRatio float64, err error)

	// GetDietMode 获取当前饮食模式
	GetDietMode(userID uint) (*models.User, error)

	// SetDietMode 设置饮食模式（手动）
	SetDietMode(userID uint, mode models.DietMode, days int, reason string) (*models.User, error)
}

// UserServiceImpl 是 UserService 接口的具体实现
type UserServiceImpl struct {
	userRepo     dao.UserRepository
	nutriService NutriService
}

// NewUserService 创建 UserService 实例，注入依赖
func NewUserService(userRepo dao.UserRepository, nutriService NutriService) UserService {
	return &UserServiceImpl{
		userRepo:     userRepo,
		nutriService: nutriService,
	}
}

func (s *UserServiceImpl) RegisterUser(user *models.User) error {
	// 检查用户是否已存在
	existingUser, _ := s.userRepo.GetUserByUsername(user.Username)
	if existingUser != nil {
		return xerr.ErrUsernameExists
	}

	return s.userRepo.CreateUser(user)
}

func (s *UserServiceImpl) LoginUser(username, password string) (*models.User, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}
	if user == nil {
		return nil, xerr.ErrUserNotFound
	}

	// TODO 可以加密密码验证
	// 验证密码
	// if !util.CheckPasswordHash(password, user.Password) {
	// 	return nil, errors.New("密码错误")
	// }
	return user, nil
}

func (s *UserServiceImpl) UpdateUserProfile(id uint, profile *models.User) error {
	// 确保用户存在
	existingUser, err := s.userRepo.GetUserByID(id)
	if err != nil || existingUser == nil {
		return xerr.ErrUserNotFound
	}

	// 更新用户健康档案
	bmr, tdee, bmi := s.nutriService.CalculateTDEE(profile)

	// 确定目标热量（这个步骤通常在食谱推荐时做，但可以在此预先计算并存储）
	// targetCalorie := s.nutriService.DetermineTargetCalorie(tdee, profile.HealthGoal)

	profile.ID = id
	profile.BMR = bmr
	profile.TDEE = tdee
	profile.BMI = bmi

	return s.userRepo.UpdateHealthProfile(id, profile)
}

func (s *UserServiceImpl) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}

// GetNutritionRequirements 获取用户营养需求
func (s *UserServiceImpl) GetNutritionRequirements(user *models.User) (targetCalorie, proteinGram, carbGram, fatGram, proteinRatio, carbRatio, fatRatio float64, err error) {
	// 检查档案是否完整
	if user.TDEE == 0 || user.HealthGoal == "" {
		return 0, 0, 0, 0, 0, 0, 0, xerr.ErrUserProfileIncomplete
	}

	// 计算目标热量
	targetCalorie = s.nutriService.DetermineTargetCalorie(user.TDEE, user.HealthGoal)

	// 计算宏量营养素分配
	proteinGram, carbGram, fatGram = s.nutriService.AllocateMacros(targetCalorie, user.HealthGoal)

	// 计算比例
	proteinRatio = (proteinGram * 4.0 / targetCalorie) * 100
	carbRatio = (carbGram * 4.0 / targetCalorie) * 100
	fatRatio = (fatGram * 9.0 / targetCalorie) * 100

	return targetCalorie, proteinGram, carbGram, fatGram, proteinRatio, carbRatio, fatRatio, nil
}

// GetDietMode 获取当前饮食模式
func (s *UserServiceImpl) GetDietMode(userID uint) (*models.User, error) {
	return s.userRepo.GetUserByID(userID)
}

// SetDietMode 设置饮食模式（手动）
func (s *UserServiceImpl) SetDietMode(userID uint, mode models.DietMode, days int, reason string) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}
	if user == nil {
		return nil, xerr.ErrUserNotFound
	}

	if days <= 0 {
		days = 1
	}
	if days > 30 {
		days = 30
	}

	user.DietMode = mode
	user.DietModeSource = "manual"
	user.DietModeReason = reason

	now := time.Now()
	if mode == models.DietModeNormal {
		user.DietModeSource = "auto"
		user.DietModeReason = ""
		user.DietModeUntil = nil
	} else {
		until := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, days)
		user.DietModeUntil = &until
	}

	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, fmt.Errorf("保存饮食模式失败: %w", err)
	}

	return user, nil
}
