package dao

import (
	"NutriPlan/internal/repository/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RecipeRepository 食谱仓储接口
type RecipeRepository interface {
	// FindByMealType 根据餐次类型以及目标用户查询食谱,并过滤掉禁忌食材
	FindByMealType(mealType models.MealType, target models.HealthGoal, forbidden []string) ([]models.Recipe, error)

	// FindByID 根据ID查询食谱
	FindByID(id uint) (*models.Recipe, error)

	// CreateDailyPlan 创建每日食谱计划
	CreateDailyPlan(plan *models.DailyRecipePlan) error

	// FindPlansByUserID 根据用户ID查询食谱计划
	FindPlansByUserID(userID uint) ([]models.DailyRecipePlan, error)

	// UpdatePlanSelection 更新计划选择状态
	UpdatePlanSelection(planID uint, selected bool) error

	// FindSelectedPlan 查询用户当前选中的计划
	FindSelectedPlan(userID uint) (*models.DailyRecipePlan, error)

	// FindRecentPlanRecipeIDs 查询用户最近几天选择的食谱ID
	FindRecentPlanRecipeIDs(userID uint, days int) ([]uint, error)

	// AddFavorite 添加收藏
	AddFavorite(userID, recipeID uint) error

	// RemoveFavorite 取消收藏
	RemoveFavorite(userID, recipeID uint) error

	// GetUserFavorites 获取用户收藏的食谱列表
	GetUserFavorites(userID uint) ([]models.Recipe, error)

	// IsFavorite 检查是否已收藏
	IsFavorite(userID, recipeID uint) (bool, error)

	// ClearAllSelectedPlans 取消用户所有已选中的计划(is_selected=false)
	ClearAllSelectedPlans(userID uint) error

	// DeleteTodayAllPlans 删除用户当天的所有计划
	DeleteTodayAllPlans(userID uint) error
}

// GormRecipeRepository GORM 实现
type GormRecipeRepository struct {
	db *gorm.DB
}

// NewGormRecipeRepository 创建 GORM 食谱仓储实例
func NewGormRecipeRepository(db *gorm.DB) RecipeRepository {
	return &GormRecipeRepository{db: db}
}

// FindByMealType 根据餐次类型以及目标用户查询食谱,并过滤掉禁忌食材
func (r *GormRecipeRepository) FindByMealType(mealType models.MealType, target models.HealthGoal, forbidden []string) ([]models.Recipe, error) {
	var recipes []models.Recipe
	tx := r.db.Model(&models.Recipe{}).Where("meal_type = ?", mealType)

	//添加禁忌食材过滤
	for _, item := range forbidden {
		if item == "" {
			continue
		}
		tx = tx.Where("ingredients NOT LIKE ?", "%"+item+"%")
	}
	err := tx.Find(&recipes).Limit(100).Error
	if err != nil {
		return nil, err
	}
	return recipes, err
}

// FindByID 根据ID查询食谱
func (r *GormRecipeRepository) FindByID(id uint) (*models.Recipe, error) {
	var recipe models.Recipe
	err := r.db.First(&recipe, id).Error
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

// CreateDailyPlan 创建每日食谱计划
func (r *GormRecipeRepository) CreateDailyPlan(plan *models.DailyRecipePlan) error {
	return r.db.Omit(clause.Associations).Create(plan).Error
}

// FindPlansByUserID 根据用户ID查询食谱计划（预加载关联的食谱）
func (r *GormRecipeRepository) FindPlansByUserID(userID uint) ([]models.DailyRecipePlan, error) {
	var plans []models.DailyRecipePlan
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
	updates := map[string]interface{}{
		"is_selected": selected,
	}
	// 如果是选中操作，同时更新plan_date为当前日期
	if selected {
		updates["plan_date"] = time.Now()
	}
	return r.db.Model(&models.DailyRecipePlan{}).
		Where("id = ?", planID).
		Updates(updates).Error
}

// FindSelectedPlan 查询用户当前选中的计划(只查询今天的)
func (r *GormRecipeRepository) FindSelectedPlan(userID uint) (*models.DailyRecipePlan, error) {
	var plan models.DailyRecipePlan
	// 构造今天的日期(只包含年月日,不包含时间)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	err := r.db.Preload("BreakfastRecipe").
		Preload("LunchRecipe").
		Preload("DinnerRecipe").
		Preload("SnackRecipe").
		Where("user_id = ? AND is_selected = ? AND plan_date = ?", userID, true, today).
		Order("created_at DESC").
		First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// FindRecentPlanRecipeIDs 查询用户最近几天选择的食谱ID
func (r *GormRecipeRepository) FindRecentPlanRecipeIDs(userID uint, days int) ([]uint, error) {
	var plans []models.DailyRecipePlan
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
		recipeIDs = append(recipeIDs, plan.SnackRecipeID)
	}
	return recipeIDs, nil
}

// AddFavorite 添加收藏（支持恢复软删除的记录）
func (r *GormRecipeRepository) AddFavorite(userID, recipeID uint) error {
	// 先检查是否存在软删除的记录
	var existing models.UserFavoriteRecipe
	err := r.db.Unscoped().Where("user_id = ? AND recipe_id = ?", userID, recipeID).First(&existing).Error

	if err == nil {
		// 记录存在（可能软删除），恢复它
		return r.db.Unscoped().Model(&existing).Update("deleted_at", nil).Error
	}

	// 不存在，创建新记录
	favorite := &models.UserFavoriteRecipe{
		UserID:   userID,
		RecipeID: recipeID,
	}
	return r.db.Create(favorite).Error
}

// RemoveFavorite 取消收藏
func (r *GormRecipeRepository) RemoveFavorite(userID, recipeID uint) error {
	return r.db.Where("user_id = ? AND recipe_id = ?", userID, recipeID).
		Delete(&models.UserFavoriteRecipe{}).Error
}

// GetUserFavorites 获取用户收藏的食谱列表
func (r *GormRecipeRepository) GetUserFavorites(userID uint) ([]models.Recipe, error) {
	var favorites []models.UserFavoriteRecipe
	err := r.db.Preload("Recipe").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&favorites).Error
	if err != nil {
		return nil, err
	}

	recipes := make([]models.Recipe, 0, len(favorites))
	for _, fav := range favorites {
		recipes = append(recipes, fav.Recipe)
	}
	return recipes, nil
}

// IsFavorite 检查是否已收藏
func (r *GormRecipeRepository) IsFavorite(userID, recipeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.UserFavoriteRecipe{}).
		Where("user_id = ? AND recipe_id = ?", userID, recipeID).
		Count(&count).Error
	return count > 0, err
}

// ClearAllSelectedPlans 取消用户所有已选中的计划
func (r *GormRecipeRepository) ClearAllSelectedPlans(userID uint) error {
	return r.db.Model(&models.DailyRecipePlan{}).
		Where("user_id = ? AND is_selected = ?", userID, true).
		Update("is_selected", false).Error
}

// DeleteTodayAllPlans 删除用户当天的所有计划
func (r *GormRecipeRepository) DeleteTodayAllPlans(userID uint) error {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return r.db.Unscoped().
		Where("user_id = ? AND plan_date = ?", userID, today).
		Delete(&models.DailyRecipePlan{}).Error
}
