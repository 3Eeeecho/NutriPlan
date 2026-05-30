package service

import (
	"NutriPlan/internal/repository/models"
	"context"
	"math"
	"time"

	"gorm.io/gorm"
)

const (
	RecipePreferenceSourceFavorite     = "favorite"
	RecipePreferenceSourceSelectedPlan = "selected_plan"
	RecipePreferenceSourceIntake       = "intake"
)

const (
	recipeFavoriteWeight     = 2.0
	recipeSelectedPlanWeight = 3.0
	recipeIntakeWeight       = 5.0
	recipePreferenceHalfLife = 30.0
)

type UserRecipePreference struct {
	UserID      uint      `json:"userId"`
	RecipeID    uint      `json:"recipeId"`
	Score       float64   `json:"score"`
	EventCount  int       `json:"eventCount"`
	LastEventAt time.Time `json:"lastEventAt"`
}

type recipePreferenceEvent struct {
	UserID    uint
	RecipeID  uint
	Weight    float64
	EventTime time.Time
	Source    string
}

type BehaviorRepository interface {
	GetUserRecipePreferences(ctx context.Context, userID uint, days int) ([]UserRecipePreference, error)
	GetAllRecipePreferences(ctx context.Context, days int) ([]UserRecipePreference, error)
}

type GormBehaviorRepository struct {
	db  *gorm.DB
	now func() time.Time
}

func NewGormBehaviorRepository(db *gorm.DB) BehaviorRepository {
	return &GormBehaviorRepository{
		db:  db,
		now: time.Now,
	}
}

func (r *GormBehaviorRepository) GetUserRecipePreferences(ctx context.Context, userID uint, days int) ([]UserRecipePreference, error) {
	return r.loadRecipePreferences(ctx, userID, days)
}

func (r *GormBehaviorRepository) GetAllRecipePreferences(ctx context.Context, days int) ([]UserRecipePreference, error) {
	return r.loadRecipePreferences(ctx, 0, days)
}

func (r *GormBehaviorRepository) loadRecipePreferences(ctx context.Context, userID uint, days int) ([]UserRecipePreference, error) {
	if r == nil || r.db == nil {
		return []UserRecipePreference{}, nil
	}
	days = normalizeBehaviorWindowDays(days)
	start := r.now().AddDate(0, 0, -days)

	events := make([]recipePreferenceEvent, 0)
	favoriteEvents, err := r.loadFavoriteEvents(ctx, userID, start)
	if err != nil {
		return nil, err
	}
	events = append(events, favoriteEvents...)

	selectedPlanEvents, err := r.loadSelectedPlanEvents(ctx, userID, start)
	if err != nil {
		return nil, err
	}
	events = append(events, selectedPlanEvents...)

	intakeEvents, err := r.loadIntakeEvents(ctx, userID, start)
	if err != nil {
		return nil, err
	}
	events = append(events, intakeEvents...)

	return aggregateRecipePreferenceEvents(events, r.now()), nil
}

func (r *GormBehaviorRepository) loadFavoriteEvents(ctx context.Context, userID uint, start time.Time) ([]recipePreferenceEvent, error) {
	var favorites []models.UserFavoriteRecipe
	tx := r.db.WithContext(ctx).Where("created_at >= ?", start)
	if userID > 0 {
		tx = tx.Where("user_id = ?", userID)
	}
	if err := tx.Find(&favorites).Error; err != nil {
		return nil, err
	}

	events := make([]recipePreferenceEvent, 0, len(favorites))
	for _, favorite := range favorites {
		if favorite.UserID == 0 || favorite.RecipeID == 0 {
			continue
		}
		events = append(events, recipePreferenceEvent{
			UserID:    favorite.UserID,
			RecipeID:  favorite.RecipeID,
			Weight:    recipeFavoriteWeight,
			EventTime: favorite.CreatedAt,
			Source:    RecipePreferenceSourceFavorite,
		})
	}
	return events, nil
}

func (r *GormBehaviorRepository) loadSelectedPlanEvents(ctx context.Context, userID uint, start time.Time) ([]recipePreferenceEvent, error) {
	var plans []models.DailyRecipePlan
	tx := r.db.WithContext(ctx).Where("is_selected = ? AND plan_date >= ?", true, start)
	if userID > 0 {
		tx = tx.Where("user_id = ?", userID)
	}
	if err := tx.Find(&plans).Error; err != nil {
		return nil, err
	}

	events := make([]recipePreferenceEvent, 0, len(plans)*3)
	for _, plan := range plans {
		if plan.UserID == 0 {
			continue
		}
		eventTime := plan.PlanDate
		if eventTime.IsZero() {
			eventTime = plan.UpdatedAt
		}
		for _, recipeID := range planRecipeIDs(&plan) {
			events = append(events, recipePreferenceEvent{
				UserID:    plan.UserID,
				RecipeID:  recipeID,
				Weight:    recipeSelectedPlanWeight,
				EventTime: eventTime,
				Source:    RecipePreferenceSourceSelectedPlan,
			})
		}
	}
	return events, nil
}

func (r *GormBehaviorRepository) loadIntakeEvents(ctx context.Context, userID uint, start time.Time) ([]recipePreferenceEvent, error) {
	var records []models.DailyIntakeRecord
	tx := r.db.WithContext(ctx).Where("food_source = ? AND source_id > 0 AND record_date >= ?", 1, start)
	if userID > 0 {
		tx = tx.Where("user_id = ?", userID)
	}
	if err := tx.Find(&records).Error; err != nil {
		return nil, err
	}

	events := make([]recipePreferenceEvent, 0, len(records))
	for _, record := range records {
		if record.UserID == 0 || record.SourceID == 0 {
			continue
		}
		events = append(events, recipePreferenceEvent{
			UserID:    record.UserID,
			RecipeID:  record.SourceID,
			Weight:    recipeIntakeWeight,
			EventTime: record.RecordDate,
			Source:    RecipePreferenceSourceIntake,
		})
	}
	return events, nil
}

func aggregateRecipePreferenceEvents(events []recipePreferenceEvent, now time.Time) []UserRecipePreference {
	type aggregate struct {
		userID      uint
		recipeID    uint
		score       float64
		eventCount  int
		lastEventAt time.Time
	}

	aggregates := make(map[[2]uint]*aggregate)
	for _, event := range events {
		if event.UserID == 0 || event.RecipeID == 0 || event.Weight == 0 {
			continue
		}
		key := [2]uint{event.UserID, event.RecipeID}
		if aggregates[key] == nil {
			aggregates[key] = &aggregate{userID: event.UserID, recipeID: event.RecipeID}
		}
		current := aggregates[key]
		current.score += event.Weight * recipePreferenceTimeDecay(event.EventTime, now)
		current.eventCount++
		if event.EventTime.After(current.lastEventAt) {
			current.lastEventAt = event.EventTime
		}
	}

	preferences := make([]UserRecipePreference, 0, len(aggregates))
	for _, current := range aggregates {
		preferences = append(preferences, UserRecipePreference{
			UserID:      current.userID,
			RecipeID:    current.recipeID,
			Score:       round4(current.score),
			EventCount:  current.eventCount,
			LastEventAt: current.lastEventAt,
		})
	}
	return preferences
}

func recipePreferenceTimeDecay(eventTime, now time.Time) float64 {
	if eventTime.IsZero() || now.IsZero() {
		return 1
	}
	daysAgo := now.Sub(eventTime).Hours() / 24
	if daysAgo <= 0 {
		return 1
	}
	return math.Exp(-daysAgo / recipePreferenceHalfLife)
}

func planRecipeIDs(plan *models.DailyRecipePlan) []uint {
	seen := make(map[uint]struct{})
	ids := make([]uint, 0, 8)
	add := func(recipeID uint) {
		if recipeID == 0 {
			return
		}
		if _, ok := seen[recipeID]; ok {
			return
		}
		seen[recipeID] = struct{}{}
		ids = append(ids, recipeID)
	}

	for _, recipeID := range plan.BreakfastItemIDs {
		add(recipeID)
	}
	for _, recipeID := range plan.LunchItemIDs {
		add(recipeID)
	}
	for _, recipeID := range plan.DinnerItemIDs {
		add(recipeID)
	}
	for _, recipeID := range plan.SnackItemIDs {
		add(recipeID)
	}

	add(plan.BreakfastRecipeID)
	add(plan.LunchRecipeID)
	add(plan.DinnerRecipeID)
	add(plan.SnackRecipeID)
	return ids
}

func normalizeBehaviorWindowDays(days int) int {
	if days <= 0 {
		return 90
	}
	if days > 365 {
		return 365
	}
	return days
}

func round4(v float64) float64 {
	return math.Round(v*10000) / 10000
}
