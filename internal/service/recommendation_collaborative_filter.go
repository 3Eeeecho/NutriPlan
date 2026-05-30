package service

import (
	"NutriPlan/internal/repository/dao"
	"NutriPlan/internal/repository/models"
	"context"
	"math"
	"sort"
	"sync"
	"time"
)

const (
	defaultCollaborativeFilterDays = 90
	defaultSimilarRecipeLimit      = 200
	defaultCollaborativeCacheTTL   = 2 * time.Minute
)

type RecipeScore struct {
	RecipeID uint          `json:"recipeId"`
	Recipe   models.Recipe `json:"recipe"`
	Score    float64       `json:"score"`
}

type CollaborativeFilter interface {
	RecommendRecipes(ctx context.Context, userID uint, mealType models.MealType, limit int) ([]RecipeScore, error)
}

type ItemCollaborativeFilter struct {
	behaviorRepo BehaviorRepository
	recipeRepo   dao.RecipeRepository
	days         int
	cacheTTL     time.Duration
	cacheMu      sync.RWMutex
	cache        map[uint]collaborativeRankCacheEntry
}

type collaborativeRecipeScore struct {
	RecipeID uint
	Score    float64
}

type collaborativeRankCacheEntry struct {
	limit     int
	ranked    []collaborativeRecipeScore
	expiresAt time.Time
}

func NewItemCollaborativeFilter(behaviorRepo BehaviorRepository, recipeRepo dao.RecipeRepository) CollaborativeFilter {
	return &ItemCollaborativeFilter{
		behaviorRepo: behaviorRepo,
		recipeRepo:   recipeRepo,
		days:         defaultCollaborativeFilterDays,
		cacheTTL:     defaultCollaborativeCacheTTL,
		cache:        make(map[uint]collaborativeRankCacheEntry),
	}
}

func (f *ItemCollaborativeFilter) RecommendRecipes(ctx context.Context, userID uint, mealType models.MealType, limit int) ([]RecipeScore, error) {
	if f == nil || f.behaviorRepo == nil || f.recipeRepo == nil || userID == 0 || limit <= 0 {
		return []RecipeScore{}, nil
	}

	candidateLimit := limit * 8
	if candidateLimit < defaultSimilarRecipeLimit {
		candidateLimit = defaultSimilarRecipeLimit
	}
	ranked, err := f.rankRecipeCandidates(ctx, userID, candidateLimit)
	if err != nil {
		return nil, err
	}
	if len(ranked) == 0 {
		return []RecipeScore{}, nil
	}

	recipeIDs := make([]uint, 0, len(ranked))
	for _, candidate := range ranked {
		recipeIDs = append(recipeIDs, candidate.RecipeID)
	}
	recipes, err := f.recipeRepo.FindByIDsFromRuntimeTable(recipeIDs)
	if err != nil {
		return nil, err
	}

	recipeByID := make(map[uint]models.Recipe, len(recipes))
	for _, recipe := range recipes {
		recipeByID[recipe.ID] = recipe
	}

	results := make([]RecipeScore, 0, limit)
	for _, candidate := range ranked {
		recipe, ok := recipeByID[candidate.RecipeID]
		if !ok || !recipeSupportsMealType(recipe, mealType) {
			continue
		}
		results = append(results, RecipeScore{
			RecipeID: candidate.RecipeID,
			Recipe:   recipe,
			Score:    candidate.Score,
		})
		if len(results) >= limit {
			break
		}
	}

	return results, nil
}

func (f *ItemCollaborativeFilter) rankRecipeCandidates(ctx context.Context, userID uint, limit int) ([]collaborativeRecipeScore, error) {
	if f == nil || f.behaviorRepo == nil || userID == 0 || limit <= 0 {
		return nil, nil
	}

	if cached, ok := f.getCachedRankedCandidates(userID, limit); ok {
		return cached, nil
	}

	userPrefs, err := f.behaviorRepo.GetUserRecipePreferences(ctx, userID, f.days)
	if err != nil {
		return nil, err
	}
	if len(userPrefs) == 0 {
		return []collaborativeRecipeScore{}, nil
	}

	allPrefs, err := f.behaviorRepo.GetAllRecipePreferences(ctx, f.days)
	if err != nil {
		return nil, err
	}

	ranked := rankCollaborativeRecipeScores(userPrefs, allPrefs, limit)
	f.setCachedRankedCandidates(userID, limit, ranked)
	return ranked, nil
}

func (f *ItemCollaborativeFilter) getCachedRankedCandidates(userID uint, limit int) ([]collaborativeRecipeScore, bool) {
	if f == nil || f.cacheTTL <= 0 {
		return nil, false
	}

	f.cacheMu.RLock()
	entry, ok := f.cache[userID]
	f.cacheMu.RUnlock()
	if !ok || entry.limit < limit || time.Now().After(entry.expiresAt) {
		return nil, false
	}

	ranked := append([]collaborativeRecipeScore(nil), entry.ranked...)
	if len(ranked) > limit {
		ranked = ranked[:limit]
	}
	return ranked, true
}

func (f *ItemCollaborativeFilter) setCachedRankedCandidates(userID uint, limit int, ranked []collaborativeRecipeScore) {
	if f == nil || f.cacheTTL <= 0 || userID == 0 {
		return
	}

	f.cacheMu.Lock()
	if f.cache == nil {
		f.cache = make(map[uint]collaborativeRankCacheEntry)
	}
	f.cache[userID] = collaborativeRankCacheEntry{
		limit:     limit,
		ranked:    append([]collaborativeRecipeScore(nil), ranked...),
		expiresAt: time.Now().Add(f.cacheTTL),
	}
	f.cacheMu.Unlock()
}

func rankCollaborativeRecipeScores(userPrefs, allPrefs []UserRecipePreference, limit int) []collaborativeRecipeScore {
	if limit <= 0 || len(userPrefs) == 0 || len(allPrefs) == 0 {
		return nil
	}

	history := make(map[uint]float64)
	for _, preference := range userPrefs {
		if preference.UserID == 0 || preference.RecipeID == 0 || preference.Score <= 0 {
			continue
		}
		history[preference.RecipeID] += preference.Score
	}
	if len(history) == 0 {
		return nil
	}

	itemVectors := buildItemUserPreferenceVectors(allPrefs)
	if len(itemVectors) == 0 {
		return nil
	}

	scoreSums := make(map[uint]float64)
	weightSums := make(map[uint]float64)
	for historyRecipeID, historyScore := range history {
		historyVector := itemVectors[historyRecipeID]
		if len(historyVector) == 0 {
			continue
		}

		for candidateRecipeID, candidateVector := range itemVectors {
			if candidateRecipeID == historyRecipeID {
				continue
			}
			if _, alreadySeen := history[candidateRecipeID]; alreadySeen {
				continue
			}

			similarity := cosineSimilarity(historyVector, candidateVector)
			if similarity <= 0 {
				continue
			}
			scoreSums[candidateRecipeID] += historyScore * similarity
			weightSums[candidateRecipeID] += math.Abs(similarity)
		}
	}

	ranked := make([]collaborativeRecipeScore, 0, len(scoreSums))
	for recipeID, scoreSum := range scoreSums {
		weightSum := weightSums[recipeID]
		if weightSum <= 0 {
			continue
		}
		ranked = append(ranked, collaborativeRecipeScore{
			RecipeID: recipeID,
			Score:    round4(scoreSum / weightSum),
		})
	}

	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].Score == ranked[j].Score {
			return ranked[i].RecipeID < ranked[j].RecipeID
		}
		return ranked[i].Score > ranked[j].Score
	})
	if len(ranked) > limit {
		ranked = ranked[:limit]
	}
	return ranked
}

func buildItemUserPreferenceVectors(preferences []UserRecipePreference) map[uint]map[uint]float64 {
	vectors := make(map[uint]map[uint]float64)
	for _, preference := range preferences {
		if preference.UserID == 0 || preference.RecipeID == 0 || preference.Score <= 0 {
			continue
		}
		if vectors[preference.RecipeID] == nil {
			vectors[preference.RecipeID] = make(map[uint]float64)
		}
		vectors[preference.RecipeID][preference.UserID] += preference.Score
	}
	return vectors
}

func cosineSimilarity(left, right map[uint]float64) float64 {
	if len(left) == 0 || len(right) == 0 {
		return 0
	}

	var dot, leftNorm, rightNorm float64
	for userID, leftScore := range left {
		leftNorm += leftScore * leftScore
		if rightScore, ok := right[userID]; ok {
			dot += leftScore * rightScore
		}
	}
	for _, rightScore := range right {
		rightNorm += rightScore * rightScore
	}
	if leftNorm == 0 || rightNorm == 0 || dot == 0 {
		return 0
	}
	return dot / (math.Sqrt(leftNorm) * math.Sqrt(rightNorm))
}
