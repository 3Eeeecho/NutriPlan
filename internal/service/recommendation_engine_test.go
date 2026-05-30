package service

import (
	"NutriPlan/internal/repository/models"
	"context"
	"testing"
)

type stubRecommendationEngine struct {
	called bool
	plans  []*models.DailyRecipePlan
}

func (e *stubRecommendationEngine) Recommend(ctx context.Context, user *models.User, count int) ([]*models.DailyRecipePlan, error) {
	e.called = true
	return e.plans, nil
}

func TestRecipeServiceDelegatesRecommendationsToEngine(t *testing.T) {
	expectedPlans := []*models.DailyRecipePlan{{UserID: 7, MatchScore: 88}}
	engine := &stubRecommendationEngine{plans: expectedPlans}
	service := &RecipeServiceImpl{recommendationEngine: engine}

	plans, err := service.RecommendRecipes(&models.User{}, 3)
	if err != nil {
		t.Fatalf("RecommendRecipes returned error: %v", err)
	}
	if !engine.called {
		t.Fatal("expected recommendation engine to be called")
	}
	if len(plans) != 1 || plans[0] != expectedPlans[0] {
		t.Fatalf("unexpected plans: %#v", plans)
	}
}
