package service

import (
	"NutriPlan/internal/repository/models"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
)

type AdminService interface {
	ListRecipes(query AdminRecipeQuery) ([]models.Recipe, int64, error)
	GetRecipe(id uint) (*models.Recipe, error)
	CreateRecipe(recipe *models.Recipe) error
	UpdateRecipe(id uint, recipe *models.Recipe) (*models.Recipe, error)
	DeleteRecipe(id uint) error
	GetHealthStats(query AdminStatsQuery) (*AdminHealthStats, error)
	GetRecipeCompletionStats(query AdminStatsQuery) (*AdminRecipeCompletionStats, error)
}

type AdminServiceImpl struct {
	db           *gorm.DB
	nutriService NutriService
}

type AdminRecipeQuery struct {
	Page     int
	PageSize int
	Keyword  string
	MealType string
}

type AdminStatsQuery struct {
	Limit int
	Days  int
}

type AdminHealthStats struct {
	Summary AdminHealthSummary    `json:"summary"`
	Users   []AdminHealthUserStat `json:"users"`
}

type AdminHealthSummary struct {
	TotalUsers        int64   `json:"totalUsers"`
	ProfileComplete   int64   `json:"profileComplete"`
	AttentionUsers    int64   `json:"attentionUsers"`
	AverageBMI        float64 `json:"averageBmi"`
	AverageCompliance float64 `json:"averageCompliance"`
}

type AdminHealthUserStat struct {
	ID                uint       `json:"id"`
	Username          string     `json:"username"`
	Role              string     `json:"role"`
	Gender            string     `json:"gender"`
	Age               int        `json:"age"`
	HealthGoal        string     `json:"healthGoal"`
	Weight            float64    `json:"weight"`
	TargetWeight      float64    `json:"targetWeight"`
	BMI               float64    `json:"bmi"`
	BMR               float64    `json:"bmr"`
	TDEE              float64    `json:"tdee"`
	TargetCalorie     float64    `json:"targetCalorie"`
	AverageIntake     float64    `json:"averageIntake"`
	Compliance        float64    `json:"compliance"`
	Risk              string     `json:"risk"`
	LastProfileUpdate *time.Time `json:"lastProfileUpdate"`
}

type AdminRecipeCompletionStats struct {
	Summary AdminRecipeCompletionSummary `json:"summary"`
	Plans   []AdminRecipeCompletionRow   `json:"plans"`
}

type AdminRecipeCompletionSummary struct {
	TotalPlans        int64   `json:"totalPlans"`
	CompletedPlans    int64   `json:"completedPlans"`
	ActivePlans       int64   `json:"activePlans"`
	MissedPlans       int64   `json:"missedPlans"`
	AverageCompletion float64 `json:"averageCompletion"`
}

type AdminRecipeCompletionRow struct {
	ID             uint      `json:"id"`
	UserID         uint      `json:"userId"`
	Username       string    `json:"username"`
	PlanDate       time.Time `json:"planDate"`
	PlanName       string    `json:"planName"`
	CompletedMeals int       `json:"completedMeals"`
	TotalMeals     int       `json:"totalMeals"`
	CompletionRate float64   `json:"completionRate"`
	Calories       float64   `json:"calories"`
	TargetCalories float64   `json:"targetCalories"`
	Status         string    `json:"status"`
}

func NewAdminService(db *gorm.DB, nutriService NutriService) AdminService {
	return &AdminServiceImpl{db: db, nutriService: nutriService}
}

func (s *AdminServiceImpl) ListRecipes(query AdminRecipeQuery) ([]models.Recipe, int64, error) {
	page, pageSize := normalizePage(query.Page, query.PageSize)
	tx := s.db.Model(&models.Recipe{})

	if keyword := strings.TrimSpace(query.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		tx = tx.Where("name LIKE ? OR ingredients LIKE ? OR target_users LIKE ?", like, like, like)
	}
	if mealType := strings.TrimSpace(query.MealType); mealType != "" {
		tx = tx.Where("meal_type = ? OR allowed_meal_types LIKE ?", mealType, "%\""+mealType+"\"%")
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var recipes []models.Recipe
	err := tx.Order("updated_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&recipes).Error
	return recipes, total, err
}

func (s *AdminServiceImpl) CreateRecipe(recipe *models.Recipe) error {
	return s.db.Create(recipe).Error
}

func (s *AdminServiceImpl) GetRecipe(id uint) (*models.Recipe, error) {
	var recipe models.Recipe
	if err := s.db.First(&recipe, id).Error; err != nil {
		return nil, err
	}
	return &recipe, nil
}

func (s *AdminServiceImpl) UpdateRecipe(id uint, recipe *models.Recipe) (*models.Recipe, error) {
	var existing models.Recipe
	if err := s.db.First(&existing, id).Error; err != nil {
		return nil, err
	}

	existing.Name = recipe.Name
	existing.ImageURL = recipe.ImageURL
	existing.MealType = recipe.MealType
	existing.AllowedMealTypes = recipe.AllowedMealTypes
	existing.Difficulty = recipe.Difficulty
	existing.CookingTime = recipe.CookingTime
	existing.PortionWeightG = recipe.PortionWeightG
	existing.Energy = recipe.Energy
	existing.Protein = recipe.Protein
	existing.Carbohydrate = recipe.Carbohydrate
	existing.Fat = recipe.Fat
	existing.Ingredients = recipe.Ingredients
	existing.CookingSteps = recipe.CookingSteps
	existing.TargetUsers = recipe.TargetUsers
	existing.ForbiddenUsers = recipe.ForbiddenUsers
	existing.IsWeightLossFriendly = recipe.IsWeightLossFriendly
	existing.IsMuscleGainFriendly = recipe.IsMuscleGainFriendly
	existing.IsSugarControlFriendly = recipe.IsSugarControlFriendly
	existing.IsGeneralFriendly = recipe.IsGeneralFriendly

	if err := s.db.Save(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

func (s *AdminServiceImpl) DeleteRecipe(id uint) error {
	return s.db.Delete(&models.Recipe{}, id).Error
}

func (s *AdminServiceImpl) GetHealthStats(query AdminStatsQuery) (*AdminHealthStats, error) {
	limit := normalizeLimit(query.Limit, 100)
	days := normalizeDays(query.Days, 7)

	var users []models.User
	if err := s.db.Order("updated_at DESC").Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}

	var totalUsers int64
	if err := s.db.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		return nil, err
	}

	start := startOfDay(time.Now()).AddDate(0, 0, -(days - 1))
	stats := &AdminHealthStats{
		Summary: AdminHealthSummary{TotalUsers: totalUsers},
		Users:   make([]AdminHealthUserStat, 0, len(users)),
	}

	var bmiSum, complianceSum float64
	var bmiCount, complianceCount int64
	for _, user := range users {
		targetCalorie := 0.0
		if user.TDEE > 0 && user.HealthGoal != "" {
			targetCalorie = s.nutriService.DetermineTargetCalorie(user.TDEE, user.HealthGoal)
		}

		averageIntake, err := s.averageDailyIntake(user.ID, start, days)
		if err != nil {
			return nil, err
		}

		compliance := 0.0
		if targetCalorie > 0 {
			compliance = round1(math.Min((averageIntake/targetCalorie)*100, 150))
			complianceSum += math.Min(compliance, 100)
			complianceCount++
		}

		risk := "stable"
		if user.TDEE <= 0 || compliance < 60 || user.BMI >= 28 {
			risk = "attention"
			stats.Summary.AttentionUsers++
		}
		if user.Age > 0 && user.Height > 0 && user.Weight > 0 {
			stats.Summary.ProfileComplete++
		}
		if user.BMI > 0 {
			bmiSum += user.BMI
			bmiCount++
		}

		stats.Users = append(stats.Users, AdminHealthUserStat{
			ID:                user.ID,
			Username:          user.Username,
			Role:              string(user.Role),
			Gender:            user.Gender,
			Age:               user.Age,
			HealthGoal:        string(user.HealthGoal),
			Weight:            user.Weight,
			TargetWeight:      user.TargetWeight,
			BMI:               round1(user.BMI),
			BMR:               round1(user.BMR),
			TDEE:              round1(user.TDEE),
			TargetCalorie:     round1(targetCalorie),
			AverageIntake:     round1(averageIntake),
			Compliance:        compliance,
			Risk:              risk,
			LastProfileUpdate: user.LastProfileUpdate,
		})
	}

	if bmiCount > 0 {
		stats.Summary.AverageBMI = round1(bmiSum / float64(bmiCount))
	}
	if complianceCount > 0 {
		stats.Summary.AverageCompliance = round1(complianceSum / float64(complianceCount))
	}

	return stats, nil
}

func (s *AdminServiceImpl) GetRecipeCompletionStats(query AdminStatsQuery) (*AdminRecipeCompletionStats, error) {
	limit := normalizeLimit(query.Limit, 100)
	days := normalizeDays(query.Days, 7)
	start := startOfDay(time.Now()).AddDate(0, 0, -(days - 1))

	var plans []models.DailyRecipePlan
	if err := s.db.Where("is_selected = ? AND plan_date >= ?", true, start).
		Order("plan_date DESC, updated_at DESC").
		Limit(limit).
		Find(&plans).Error; err != nil {
		return nil, err
	}

	stats := &AdminRecipeCompletionStats{
		Plans: make([]AdminRecipeCompletionRow, 0, len(plans)),
	}

	var completionSum float64
	for _, plan := range plans {
		username, err := s.usernameByID(plan.UserID)
		if err != nil {
			return nil, err
		}

		completedMeals, calories, err := s.planCompletion(plan.UserID, plan.PlanDate)
		if err != nil {
			return nil, err
		}

		totalMeals := planMealCount(&plan)
		if totalMeals == 0 {
			totalMeals = 3
		}
		completionRate := round1(math.Min(float64(completedMeals)/float64(totalMeals)*100, 100))
		status := completionStatus(plan.PlanDate, completedMeals, totalMeals)

		switch status {
		case "done":
			stats.Summary.CompletedPlans++
		case "active":
			stats.Summary.ActivePlans++
		default:
			stats.Summary.MissedPlans++
		}
		completionSum += completionRate

		stats.Plans = append(stats.Plans, AdminRecipeCompletionRow{
			ID:             plan.ID,
			UserID:         plan.UserID,
			Username:       username,
			PlanDate:       plan.PlanDate,
			PlanName:       plan.PlanName,
			CompletedMeals: completedMeals,
			TotalMeals:     totalMeals,
			CompletionRate: completionRate,
			Calories:       round1(calories),
			TargetCalories: round1(plan.TargetEnergy),
			Status:         status,
		})
	}

	stats.Summary.TotalPlans = int64(len(plans))
	if len(plans) > 0 {
		stats.Summary.AverageCompletion = round1(completionSum / float64(len(plans)))
	}

	return stats, nil
}

func (s *AdminServiceImpl) averageDailyIntake(userID uint, start time.Time, days int) (float64, error) {
	var total float64
	err := s.db.Model(&models.DailyIntakeRecord{}).
		Where("user_id = ? AND record_date >= ?", userID, start).
		Select("COALESCE(SUM(calculated_energy), 0)").
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return total / float64(days), nil
}

func (s *AdminServiceImpl) usernameByID(userID uint) (string, error) {
	var user models.User
	if err := s.db.Select("id", "username").First(&user, userID).Error; err != nil {
		return "", err
	}
	return user.Username, nil
}

func (s *AdminServiceImpl) planCompletion(userID uint, date time.Time) (int, float64, error) {
	start := startOfDay(date)
	end := start.AddDate(0, 0, 1)

	var records []models.DailyIntakeRecord
	if err := s.db.Where("user_id = ? AND record_date >= ? AND record_date < ?", userID, start, end).
		Find(&records).Error; err != nil {
		return 0, 0, err
	}

	mealSet := make(map[string]struct{})
	var calories float64
	for _, record := range records {
		calories += record.CalculatedEnergy
		if record.FoodSource == 1 && strings.TrimSpace(record.MealType) != "" {
			mealSet[record.MealType] = struct{}{}
		}
	}

	return len(mealSet), calories, nil
}

func normalizePage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func normalizeLimit(limit, fallback int) int {
	if limit <= 0 {
		limit = fallback
	}
	if limit > 500 {
		limit = 500
	}
	return limit
}

func normalizeDays(days, fallback int) int {
	if days <= 0 {
		days = fallback
	}
	if days > 90 {
		days = 90
	}
	return days
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func planMealCount(plan *models.DailyRecipePlan) int {
	count := 0
	if plan.BreakfastRecipeID > 0 || len(plan.BreakfastItemIDs) > 0 {
		count++
	}
	if plan.LunchRecipeID > 0 || len(plan.LunchItemIDs) > 0 {
		count++
	}
	if plan.DinnerRecipeID > 0 || len(plan.DinnerItemIDs) > 0 {
		count++
	}
	if plan.SnackRecipeID > 0 || len(plan.SnackItemIDs) > 0 {
		count++
	}
	return count
}

func completionStatus(date time.Time, completedMeals, totalMeals int) string {
	if completedMeals >= totalMeals {
		return "done"
	}
	if startOfDay(date).Equal(startOfDay(time.Now())) {
		return "active"
	}
	return "missed"
}

func round1(v float64) float64 {
	return math.Round(v*10) / 10
}
