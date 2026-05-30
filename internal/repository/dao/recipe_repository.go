package dao

import (
	"NutriPlan/internal/repository/models"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RecipeRepository interface {
	FindByMealType(mealType models.MealType, target models.HealthGoal, forbidden []string) ([]models.Recipe, error)

	FindLowCalorieRecipes(mealType models.MealType, target models.HealthGoal, limit int) ([]models.Recipe, error)
	FindHighCalorieRecipes(mealType models.MealType, target models.HealthGoal, limit int) ([]models.Recipe, error)
	FindByIngredientTags(mealType models.MealType, target models.HealthGoal, tags []string, limit int) ([]models.Recipe, error)
	FindRandomRecipes(mealType models.MealType, target models.HealthGoal, limit int) ([]models.Recipe, error)

	FindByID(id uint) (*models.Recipe, error)
	FindByIDsFromRuntimeTable(ids []uint) ([]models.Recipe, error)
	ResolveRecipeIDForPlan(recipeID uint) (uint, error)

	CreateDailyPlan(plan *models.DailyRecipePlan) error
	FindPlansByUserID(userID uint) ([]models.DailyRecipePlan, error)
	UpdatePlanSelection(planID uint, selected bool) error
	FindSelectedPlan(userID uint) (*models.DailyRecipePlan, error)
	FindRecentPlanRecipeIDs(userID uint, days int) ([]uint, error)

	AddFavorite(userID, recipeID uint) error
	RemoveFavorite(userID, recipeID uint) error
	GetUserFavorites(userID uint) ([]models.Recipe, error)
	IsFavorite(userID, recipeID uint) (bool, error)

	ClearAllSelectedPlans(userID uint) error
	DeleteTodayAllPlans(userID uint) error

	FindFuzzyByName(keyword string) (*models.Recipe, error)
}

type GormRecipeRepository struct {
	db          *gorm.DB
	recipeTable string
}

func NewGormRecipeRepository(db *gorm.DB, recipeTable string) RecipeRepository {
	_ = recipeTable
	return &GormRecipeRepository{db: db, recipeTable: "recipes"}
}

func (r *GormRecipeRepository) recipeQuery() *gorm.DB {
	return r.db.Model(&models.Recipe{}).Table(r.recipeTable)
}

func (r *GormRecipeRepository) applyMealTypeFilter(tx *gorm.DB, mealType models.MealType) *gorm.DB {
	allowedMealTypePattern := "%\"" + string(mealType) + "\"%"
	return tx.Where(
		"(allowed_meal_types LIKE ? OR ((allowed_meal_types IS NULL OR allowed_meal_types = '' OR allowed_meal_types = '[]') AND meal_type = ?))",
		allowedMealTypePattern,
		mealType,
	)
}

func (r *GormRecipeRepository) applyHealthGoalFilter(tx *gorm.DB, target models.HealthGoal) *gorm.DB {
	switch target {
	case models.GoalWeightLoss:
		return tx.Where("is_weight_loss_friendly = ?", true)
	case models.GoalMuscleGain:
		return tx.Where("is_muscle_gain_friendly = ?", true)
	case models.GoalSugarControl:
		return tx.Where("is_sugar_control_friendly = ?", true)
	case models.GoalMaintain:
		return tx.Where("is_general_friendly = ?", true)
	default:
		return tx
	}
}

func (r *GormRecipeRepository) FindByMealType(mealType models.MealType, target models.HealthGoal, forbidden []string) ([]models.Recipe, error) {
	var recipes []models.Recipe
	tx := r.applyHealthGoalFilter(r.applyMealTypeFilter(r.recipeQuery(), mealType), target)

	for _, item := range forbidden {
		if item == "" {
			continue
		}
		tx = tx.Where("ingredients NOT LIKE ?", "%"+item+"%")
	}

	err := tx.Find(&recipes).Error
	if err != nil {
		return nil, err
	}
	return recipes, nil
}

func (r *GormRecipeRepository) FindLowCalorieRecipes(mealType models.MealType, target models.HealthGoal, limit int) ([]models.Recipe, error) {
	var recipes []models.Recipe
	err := r.applyHealthGoalFilter(
		r.applyMealTypeFilter(r.recipeQuery().Where("energy > 0"), mealType),
		target,
	).Order("energy ASC").Limit(limit).Find(&recipes).Error
	return recipes, err
}

func (r *GormRecipeRepository) FindHighCalorieRecipes(mealType models.MealType, target models.HealthGoal, limit int) ([]models.Recipe, error) {
	var recipes []models.Recipe
	err := r.applyHealthGoalFilter(
		r.applyMealTypeFilter(r.recipeQuery(), mealType),
		target,
	).Order("energy DESC").Limit(limit).Find(&recipes).Error
	return recipes, err
}

func (r *GormRecipeRepository) FindByIngredientTags(mealType models.MealType, target models.HealthGoal, tags []string, limit int) ([]models.Recipe, error) {
	var recipes []models.Recipe
	if len(tags) == 0 {
		return r.FindRandomRecipes(mealType, target, limit)
	}

	tx := r.applyHealthGoalFilter(r.applyMealTypeFilter(r.recipeQuery(), mealType), target)
	var orConditions []string
	var args []interface{}
	for _, tag := range tags {
		if tag == "" {
			continue
		}
		orConditions = append(orConditions, "ingredients LIKE ?")
		args = append(args, "%"+tag+"%")
	}

	if len(orConditions) > 0 {
		tx = tx.Where(strings.Join(orConditions, " OR "), args...)
	}

	err := tx.Limit(limit).Find(&recipes).Error
	return recipes, err
}

func (r *GormRecipeRepository) FindRandomRecipes(mealType models.MealType, target models.HealthGoal, limit int) ([]models.Recipe, error) {
	var recipes []models.Recipe
	err := r.applyHealthGoalFilter(
		r.applyMealTypeFilter(r.recipeQuery(), mealType),
		target,
	).Order("RAND()").Limit(limit).Find(&recipes).Error
	return recipes, err
}

func (r *GormRecipeRepository) FindByID(id uint) (*models.Recipe, error) {
	var recipe models.Recipe
	err := r.recipeQuery().Where("id = ?", id).First(&recipe).Error
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

func (r *GormRecipeRepository) FindByIDsFromRuntimeTable(ids []uint) ([]models.Recipe, error) {
	if len(ids) == 0 {
		return []models.Recipe{}, nil
	}

	var recipes []models.Recipe
	err := r.db.Model(&models.Recipe{}).
		Table("recipes").
		Where("id IN ?", ids).
		Find(&recipes).Error
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func (r *GormRecipeRepository) ResolveRecipeIDForPlan(recipeID uint) (uint, error) {
	if recipeID == 0 {
		return 0, nil
	}
	return recipeID, nil
}

func (r *GormRecipeRepository) CreateDailyPlan(plan *models.DailyRecipePlan) error {
	return r.db.Omit(clause.Associations).Create(plan).Error
}

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

func (r *GormRecipeRepository) UpdatePlanSelection(planID uint, selected bool) error {
	updates := map[string]interface{}{
		"is_selected": selected,
	}
	if selected {
		updates["plan_date"] = time.Now()
	}
	return r.db.Model(&models.DailyRecipePlan{}).
		Where("id = ?", planID).
		Updates(updates).Error
}

func (r *GormRecipeRepository) FindSelectedPlan(userID uint) (*models.DailyRecipePlan, error) {
	var plan models.DailyRecipePlan
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

func (r *GormRecipeRepository) AddFavorite(userID, recipeID uint) error {
	var existing models.UserFavoriteRecipe
	err := r.db.Unscoped().Where("user_id = ? AND recipe_id = ?", userID, recipeID).First(&existing).Error

	if err == nil {
		return r.db.Unscoped().Model(&existing).Update("deleted_at", nil).Error
	}

	favorite := &models.UserFavoriteRecipe{
		UserID:   userID,
		RecipeID: recipeID,
	}
	return r.db.Create(favorite).Error
}

func (r *GormRecipeRepository) RemoveFavorite(userID, recipeID uint) error {
	return r.db.Where("user_id = ? AND recipe_id = ?", userID, recipeID).
		Delete(&models.UserFavoriteRecipe{}).Error
}

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

func (r *GormRecipeRepository) IsFavorite(userID, recipeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.UserFavoriteRecipe{}).
		Where("user_id = ? AND recipe_id = ?", userID, recipeID).
		Count(&count).Error
	return count > 0, err
}

func (r *GormRecipeRepository) ClearAllSelectedPlans(userID uint) error {
	return r.db.Model(&models.DailyRecipePlan{}).
		Where("user_id = ? AND is_selected = ?", userID, true).
		Update("is_selected", false).Error
}

func (r *GormRecipeRepository) DeleteTodayAllPlans(userID uint) error {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return r.db.Unscoped().
		Where("user_id = ? AND plan_date = ?", userID, today).
		Delete(&models.DailyRecipePlan{}).Error
}

func (r *GormRecipeRepository) FindFuzzyByName(keyword string) (*models.Recipe, error) {
	var recipe models.Recipe

	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil, nil
	}

	err := r.recipeQuery().Where("name LIKE ?", "%"+keyword+"%").
		Order("LENGTH(name) ASC").
		First(&recipe).Error
	if err != nil {
		return nil, err
	}

	return &recipe, nil
}
