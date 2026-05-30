package service

import (
	"NutriPlan/internal/repository/models"
	"context"
)

type RecommendationEngine interface {
	Recommend(ctx context.Context, user *models.User, count int) ([]*models.DailyRecipePlan, error)
}

type RuleBasedRecommendationEngine struct {
	recipeService *RecipeServiceImpl
}

func NewRuleBasedRecommendationEngine(recipeService *RecipeServiceImpl) RecommendationEngine {
	return &RuleBasedRecommendationEngine{recipeService: recipeService}
}

func (e *RuleBasedRecommendationEngine) Recommend(ctx context.Context, user *models.User, count int) ([]*models.DailyRecipePlan, error) {
	_ = ctx
	return e.recipeService.recommendRecipesRuleBased(user, count)
}
