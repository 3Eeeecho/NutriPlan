package service

import (
	"NutriPlan/internal/core/domain"
	"NutriPlan/internal/repository"
	"errors"
	"fmt"
)

// UserService 定义了用户业务流程的接口
type UserService interface {
	// RegisterUser 处理用户注册逻辑，包括密码加密和创建数据库记录
	RegisterUser(user *domain.User) error

	// LoginUser 处理用户登录逻辑，验证密码等
	LoginUser(username, password string) (*domain.User, error)

	// UpdateUserProfile 处理用户健康档案更新，包括触发 TDEE 计算
	UpdateUserProfile(id uint, profile *domain.User) error

	// GetUserByID 获取用户信息
	GetUserByID(id uint) (*domain.User, error)

	// GetNutritionRequirements 获取用户营养需求
	GetNutritionRequirements(user *domain.User) (targetCalorie, proteinGram, carbGram, fatGram, proteinRatio, carbRatio, fatRatio float64, err error)
}

// UserServiceImpl 是 UserService 接口的具体实现
type UserServiceImpl struct {
	userRepo     repository.UserRepository
	nutriService NutriService
}

// NewUserService 创建 UserService 实例，注入依赖
func NewUserService(userRepo repository.UserRepository, nutriService NutriService) UserService {
	return &UserServiceImpl{
		userRepo:     userRepo,
		nutriService: nutriService,
	}
}

func (s *UserServiceImpl) RegisterUser(user *domain.User) error {
	// 检查用户是否已存在
	existingUser, _ := s.userRepo.GetUserByUsername(user.Username)
	if existingUser != nil {
		return errors.New("用户名已存在")
	}

	return s.userRepo.CreateUser(user)
}

func (s *UserServiceImpl) LoginUser(username, password string) (*domain.User, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 验证密码
	// if !util.CheckPasswordHash(password, user.Password) {
	// 	return nil, errors.New("密码错误")
	// }
	return user, nil
}

func (s *UserServiceImpl) UpdateUserProfile(id uint, profile *domain.User) error {
	// 确保用户存在
	existingUser, err := s.userRepo.GetUserByID(id)
	if err != nil || existingUser == nil {
		return errors.New("用户不存在")
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

func (s *UserServiceImpl) GetUserByID(id uint) (*domain.User, error) {
	return s.userRepo.GetUserByID(id)
}

// GetNutritionRequirements 获取用户营养需求
func (s *UserServiceImpl) GetNutritionRequirements(user *domain.User) (targetCalorie, proteinGram, carbGram, fatGram, proteinRatio, carbRatio, fatRatio float64, err error) {
	// 检查档案是否完整
	if user.TDEE == 0 || user.HealthGoal == "" {
		return 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("用户档案不完整，请先完善健康档案")
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
