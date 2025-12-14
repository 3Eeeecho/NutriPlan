package service

import (
	"NutriPlan/internal/repository/dao"
	"NutriPlan/internal/repository/models"
	"encoding/json"
	"regexp"
	"sort"
	"strings"
	"time"
)

type ShoppingListService interface {
	GenerateFromRecipes(userID uint, recipeIDs []uint, listName string) (*models.ShoppingList, error)
	GetShoppingLists(userID uint, page, pageSize int) ([]models.ShoppingList, int64, error)
	GetShoppingListDetail(id uint) (*models.ShoppingList, error)
	UpdateShoppingList(list *models.ShoppingList) error
	DeleteShoppingList(id uint) error
	CompleteShoppingList(id uint) error
}

type shoppingListService struct {
	repo       dao.ShoppingListRepository
	recipeRepo dao.RecipeRepository
}

func NewShoppingListService(repo dao.ShoppingListRepository, recipeRepo dao.RecipeRepository) ShoppingListService {
	return &shoppingListService{
		repo:       repo,
		recipeRepo: recipeRepo,
	}
}

// ShoppingItem 清单项结构
type ShoppingItem struct {
	Name     string `json:"name"`
	Amount   string `json:"amount"`
	Category string `json:"category"`
	Checked  bool   `json:"checked"`
}

// CategoryGroup 分组结构
type CategoryGroup struct {
	Category string         `json:"category"`
	Items    []ShoppingItem `json:"items"`
}

func (s *shoppingListService) GenerateFromRecipes(userID uint, recipeIDs []uint, listName string) (*models.ShoppingList, error) {
	// 1. 获取所有食谱
	var recipes []models.Recipe
	for _, id := range recipeIDs {
		recipe, err := s.recipeRepo.FindByID(id)
		if err != nil {
			continue
		}
		recipes = append(recipes, *recipe)
	}

	// 2. 提取并合并食材
	ingredientMap := make(map[string][]string) // name -> amounts

	for _, recipe := range recipes {
		var ingredients []string
		// 尝试解析 JSON
		if err := json.Unmarshal([]byte(recipe.Ingredients), &ingredients); err != nil {
			// 如果不是 JSON，尝试按逗号分隔
			ingredients = strings.Split(recipe.Ingredients, ",")
		}

		for _, raw := range ingredients {
			name, amount := parseIngredient(raw)
			if name != "" {
				ingredientMap[name] = append(ingredientMap[name], amount)
			}
		}
	}

	// 3. 分类并构建列表
	var items []ShoppingItem
	for name, amounts := range ingredientMap {
		// 合并数量 (简单拼接)
		totalAmount := mergeAmounts(amounts)
		category := categorizeIngredient(name)

		items = append(items, ShoppingItem{
			Name:     name,
			Amount:   totalAmount,
			Category: category,
			Checked:  false,
		})
	}

	// 4. 排序 (按分类顺序)
	sort.Slice(items, func(i, j int) bool {
		orderI := getCategoryOrder(items[i].Category)
		orderJ := getCategoryOrder(items[j].Category)
		if orderI != orderJ {
			return orderI < orderJ
		}
		return items[i].Name < items[j].Name
	})

	// 5. 序列化并保存
	itemsJSON, _ := json.Marshal(items)
	recipeIDsJSON, _ := json.Marshal(recipeIDs)

	list := &models.ShoppingList{
		UserID:     userID,
		Name:       listName,
		RecipeIDs:  string(recipeIDsJSON),
		Items:      string(itemsJSON),
		Status:     "pending",
		ShopDate:   time.Now(),
		TotalItems: len(items),
	}

	if err := s.repo.Create(list); err != nil {
		return nil, err
	}

	return list, nil
}

func (s *shoppingListService) GetShoppingLists(userID uint, page, pageSize int) ([]models.ShoppingList, int64, error) {
	return s.repo.GetByUserID(userID, page, pageSize)
}

func (s *shoppingListService) GetShoppingListDetail(id uint) (*models.ShoppingList, error) {
	return s.repo.GetByID(id)
}

func (s *shoppingListService) UpdateShoppingList(list *models.ShoppingList) error {
	return s.repo.Update(list)
}

func (s *shoppingListService) DeleteShoppingList(id uint) error {
	return s.repo.Delete(id)
}

func (s *shoppingListService) CompleteShoppingList(id uint) error {
	return s.repo.UpdateStatus(id, "completed")
}

// 辅助函数

func parseIngredient(raw string) (string, string) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", ""
	}

	// 简单尝试分离名称和数量
	// 假设格式为 "名称 数量" 或 "名称"
	// 使用正则匹配末尾的数字和单位
	re := regexp.MustCompile(`^(.+?)\s*(\d+[a-zA-Z\p{Han}]*)$`)
	matches := re.FindStringSubmatch(raw)

	if len(matches) == 3 {
		return strings.TrimSpace(matches[1]), matches[2]
	}

	return raw, ""
}

func mergeAmounts(amounts []string) string {
	var validAmounts []string
	for _, a := range amounts {
		if a != "" {
			validAmounts = append(validAmounts, a)
		}
	}
	if len(validAmounts) == 0 {
		return ""
	}
	// 这里不做复杂的单位换算，直接拼接，例如 "200g + 100g"
	// 实际项目中可以引入单位换算库
	return strings.Join(validAmounts, " + ")
}

func categorizeIngredient(name string) string {
	name = strings.ToLower(name)

	rules := map[string][]string{
		"蔬菜类":  {"菜", "瓜", "菇", "茄", "椒", "笋", "藕", "豆角", "萝卜", "洋葱", "蒜", "姜", "葱", "西兰花", "菠菜"},
		"肉类海鲜": {"肉", "鸡", "鸭", "鱼", "虾", "蟹", "牛", "羊", "猪", "排骨", "培根", "火腿"},
		"豆制品":  {"豆腐", "豆浆", "腐竹", "豆干"},
		"蛋奶类":  {"蛋", "奶", "芝士", "黄油", "乳"},
		"主食类":  {"米", "面", "粉", "馒头", "饼", "薯", "玉米", "燕麦"},
		"水果类":  {"果", "蕉", "梨", "桃", "橘", "橙", "莓"},
		"调料类":  {"油", "盐", "糖", "酱", "醋", "酒", "粉", "精", "辣", "椒粉"},
	}

	for category, keywords := range rules {
		for _, kw := range keywords {
			if strings.Contains(name, kw) {
				return category
			}
		}
	}

	return "其他"
}

func getCategoryOrder(category string) int {
	order := map[string]int{
		"蔬菜类":  1,
		"肉类海鲜": 2,
		"豆制品":  3,
		"蛋奶类":  4,
		"主食类":  5,
		"水果类":  6,
		"调料类":  7,
		"其他":   8,
	}
	if v, ok := order[category]; ok {
		return v
	}
	return 99
}
