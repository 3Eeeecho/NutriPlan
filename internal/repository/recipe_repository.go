package repository

import (
	"NutriPlan/internal/core/domain"
	"time"

	"gorm.io/gorm"
)

// RecipeRepository 食谱仓储接口
type RecipeRepository interface {
	// FindByMealType 根据餐次类型查询食谱
	FindByMealType(mealType domain.MealType) ([]domain.Recipe, error)

	// FindByID 根据ID查询食谱
	FindByID(id uint) (*domain.Recipe, error)

	// FindAll 查询所有食谱
	FindAll() ([]domain.Recipe, error)

	// Create 创建食谱
	Create(recipe *domain.Recipe) error

	// CreateDailyPlan 创建每日食谱计划
	CreateDailyPlan(plan *domain.DailyRecipePlan) error

	// FindPlansByUserID 根据用户ID查询食谱计划
	FindPlansByUserID(userID uint) ([]domain.DailyRecipePlan, error)

	// UpdatePlanSelection 更新计划选择状态
	UpdatePlanSelection(planID uint, selected bool) error

	// FindSelectedPlan 查询用户当前选中的计划
	FindSelectedPlan(userID uint) (*domain.DailyRecipePlan, error)

	// FindRecentPlanRecipeIDs 查询用户最近几天选择的食谱ID
	FindRecentPlanRecipeIDs(userID uint, days int) ([]uint, error)
}

// GormRecipeRepository GORM 实现
type GormRecipeRepository struct {
	db *gorm.DB
}

// NewGormRecipeRepository 创建 GORM 食谱仓储实例
func NewGormRecipeRepository(db *gorm.DB) RecipeRepository {
	return &GormRecipeRepository{db: db}
}

// FindByMealType 根据餐次类型查询食谱
func (r *GormRecipeRepository) FindByMealType(mealType domain.MealType) ([]domain.Recipe, error) {
	var recipes []domain.Recipe
	err := r.db.Where("meal_type = ?", mealType).Find(&recipes).Error
	return recipes, err
}

// FindByID 根据ID查询食谱
func (r *GormRecipeRepository) FindByID(id uint) (*domain.Recipe, error) {
	var recipe domain.Recipe
	err := r.db.First(&recipe, id).Error
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

// FindAll 查询所有食谱
func (r *GormRecipeRepository) FindAll() ([]domain.Recipe, error) {
	var recipes []domain.Recipe
	err := r.db.Find(&recipes).Error
	return recipes, err
}

// Create 创建食谱
func (r *GormRecipeRepository) Create(recipe *domain.Recipe) error {
	return r.db.Create(recipe).Error
}

// CreateDailyPlan 创建每日食谱计划
func (r *GormRecipeRepository) CreateDailyPlan(plan *domain.DailyRecipePlan) error {
	return r.db.Create(plan).Error
}

// FindPlansByUserID 根据用户ID查询食谱计划（预加载关联的食谱）
func (r *GormRecipeRepository) FindPlansByUserID(userID uint) ([]domain.DailyRecipePlan, error) {
	var plans []domain.DailyRecipePlan
	err := r.db.Preload("BreakfastRecipe").
		Preload("LunchRecipe").
		Preload("DinnerRecipe").
		Preload("SnackRecipe").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&plans).Error
	return plans, err
}

// UpdatePlanSelection 更新计划选择状态
func (r *GormRecipeRepository) UpdatePlanSelection(planID uint, selected bool) error {
	return r.db.Model(&domain.DailyRecipePlan{}).
		Where("id = ?", planID).
		Update("is_selected", selected).Error
}

// FindSelectedPlan 查询用户当前选中的计划
func (r *GormRecipeRepository) FindSelectedPlan(userID uint) (*domain.DailyRecipePlan, error) {
	var plan domain.DailyRecipePlan
	err := r.db.Preload("BreakfastRecipe").
		Preload("LunchRecipe").
		Preload("DinnerRecipe").
		Preload("SnackRecipe").
		Where("user_id = ? AND is_selected = ?", userID, true).
		Order("created_at DESC").
		First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// FindRecentPlanRecipeIDs 查询用户最近几天选择的食谱ID
func (r *GormRecipeRepository) FindRecentPlanRecipeIDs(userID uint, days int) ([]uint, error) {
	var plans []domain.DailyRecipePlan
	startTime := time.Now().AddDate(0, 0, -days)

	err := r.db.Where("user_id = ? AND is_selected = ? AND created_at >= ?", userID, true, startTime).
		Find(&plans).Error
	if err != nil {
		return nil, err
	}

	recipeIDs := make([]uint, 0)
	for _, plan := range plans {
		recipeIDs = append(recipeIDs, plan.BreakfastRecipeID)
		recipeIDs = append(recipeIDs, plan.LunchRecipeID)
		recipeIDs = append(recipeIDs, plan.DinnerRecipeID)
		if plan.SnackRecipeID != nil {
			recipeIDs = append(recipeIDs, *plan.SnackRecipeID)
		}
	}
	return recipeIDs, nil
}
