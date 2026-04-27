package handler

import (
	"NutriPlan/internal/repository/models"
	"NutriPlan/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminService service.AdminService
}

type AdminRecipeRequest struct {
	Name                   string   `json:"name" binding:"required"`
	ImageURL               string   `json:"imageUrl"`
	MealType               string   `json:"mealType" binding:"required"`
	AllowedMealTypes       []string `json:"allowedMealTypes"`
	Difficulty             string   `json:"difficulty"`
	CookingTime            int      `json:"cookingTime"`
	PortionWeightG         float64  `json:"portionWeightG"`
	Energy                 float64  `json:"energy" binding:"required"`
	Protein                float64  `json:"protein"`
	Carbohydrate           float64  `json:"carbohydrate"`
	Fat                    float64  `json:"fat"`
	Ingredients            []string `json:"ingredients"`
	CookingSteps           []string `json:"cookingSteps"`
	TargetUsers            []string `json:"targetUsers"`
	ForbiddenUsers         []string `json:"forbiddenUsers"`
	IsWeightLossFriendly   bool     `json:"isWeightLossFriendly"`
	IsMuscleGainFriendly   bool     `json:"isMuscleGainFriendly"`
	IsSugarControlFriendly bool     `json:"isSugarControlFriendly"`
	IsGeneralFriendly      *bool    `json:"isGeneralFriendly"`
}

func NewAdminHandler(adminService service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func (h *AdminHandler) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":  "admin access granted",
		"username": c.GetString("username"),
		"role":     c.GetString("role"),
	})
}

func (h *AdminHandler) ListRecipes(c *gin.Context) {
	page := parseIntQuery(c, "page", 1)
	pageSize := parseIntQuery(c, "pageSize", parseIntQuery(c, "page_size", 20))
	recipes, total, err := h.adminService.ListRecipes(service.AdminRecipeQuery{
		Page:     page,
		PageSize: pageSize,
		Keyword:  c.Query("keyword"),
		MealType: c.Query("mealType"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load recipes: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recipes":  recipes,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *AdminHandler) CreateRecipe(c *gin.Context) {
	var req AdminRecipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe payload: " + err.Error()})
		return
	}

	recipe := buildRecipeFromAdminRequest(req)
	if err := h.adminService.CreateRecipe(recipe); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create recipe: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "recipe created", "recipe": recipe})
}

func (h *AdminHandler) GetRecipe(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	recipe, err := h.adminService.GetRecipe(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recipe": recipe})
}

func (h *AdminHandler) UpdateRecipe(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	var req AdminRecipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe payload: " + err.Error()})
		return
	}

	recipe, err := h.adminService.UpdateRecipe(id, buildRecipeFromAdminRequest(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update recipe: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "recipe updated", "recipe": recipe})
}

func (h *AdminHandler) DeleteRecipe(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	if err := h.adminService.DeleteRecipe(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete recipe: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "recipe deleted"})
}

func (h *AdminHandler) GetHealthStats(c *gin.Context) {
	stats, err := h.adminService.GetHealthStats(service.AdminStatsQuery{
		Limit: parseIntQuery(c, "limit", 100),
		Days:  parseIntQuery(c, "days", 7),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load health stats: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *AdminHandler) GetRecipeCompletionStats(c *gin.Context) {
	stats, err := h.adminService.GetRecipeCompletionStats(service.AdminStatsQuery{
		Limit: parseIntQuery(c, "limit", 100),
		Days:  parseIntQuery(c, "days", 7),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load recipe completion stats: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func buildRecipeFromAdminRequest(req AdminRecipeRequest) *models.Recipe {
	isGeneralFriendly := true
	if req.IsGeneralFriendly != nil {
		isGeneralFriendly = *req.IsGeneralFriendly
	}
	portionWeight := req.PortionWeightG
	if portionWeight <= 0 {
		portionWeight = 100
	}

	return &models.Recipe{
		Name:                   strings.TrimSpace(req.Name),
		ImageURL:               strings.TrimSpace(req.ImageURL),
		MealType:               models.MealType(strings.TrimSpace(req.MealType)),
		AllowedMealTypes:       cleanStringList(req.AllowedMealTypes),
		Difficulty:             strings.TrimSpace(req.Difficulty),
		CookingTime:            req.CookingTime,
		PortionWeightG:         portionWeight,
		Energy:                 req.Energy,
		Protein:                req.Protein,
		Carbohydrate:           req.Carbohydrate,
		Fat:                    req.Fat,
		Ingredients:            cleanStringList(req.Ingredients),
		CookingSteps:           cleanStringList(req.CookingSteps),
		TargetUsers:            cleanStringList(req.TargetUsers),
		ForbiddenUsers:         cleanStringList(req.ForbiddenUsers),
		IsWeightLossFriendly:   req.IsWeightLossFriendly,
		IsMuscleGainFriendly:   req.IsMuscleGainFriendly,
		IsSugarControlFriendly: req.IsSugarControlFriendly,
		IsGeneralFriendly:      isGeneralFriendly,
	}
}

func cleanStringList(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}

func parseIDParam(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}
	return uint(id), true
}

func parseIntQuery(c *gin.Context, key string, fallback int) int {
	raw := c.Query(key)
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}
