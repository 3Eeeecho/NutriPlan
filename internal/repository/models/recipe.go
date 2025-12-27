package models

import (
	"gorm.io/gorm"
)

type Recipe struct {
	gorm.Model // 包含了 ID, CreatedAt, UpdatedAt, DeletedAt

	Name        string   `gorm:"type:varchar(200);not null;index;comment:食谱名称" json:"name"`
	ImageURL    string   `gorm:"type:varchar(500);comment:封面图片链接" json:"imageUrl"`
	MealType    MealType `gorm:"type:varchar(20);not null;index;comment:适合餐点类型" json:"mealType"`
	Difficulty  string   `gorm:"type:varchar(20);comment:烹饪难度" json:"difficulty"`
	CookingTime int      `gorm:"default:0;comment:烹饪耗时(分)" json:"cookingTime"`

	// --- 核心营养素
	Energy       float64 `gorm:"type:decimal(10,2);not null;index;comment:总热量 (kcal)" json:"energy"` // 加 index 方便按热量排序推荐
	Protein      float64 `gorm:"type:decimal(10,2);not null;comment:蛋白质 (g)" json:"protein"`
	Carbohydrate float64 `gorm:"type:decimal(10,2);not null;comment:碳水化合物 (g)" json:"carbohydrate"`
	Fat          float64 `gorm:"type:decimal(10,2);not null;comment:脂肪 (g)" json:"fat"`

	// --- 复杂结构 (直接使用 Slice，让 GORM 自动序列化) ---
	Ingredients  []string `gorm:"serializer:json;not null;comment:所需食材清单" json:"ingredients"`
	CookingSteps []string `gorm:"serializer:json;comment:烹饪步骤" json:"cookingSteps"`

	// TargetUsers: 推荐给谁吃 (白名单) -> 例如 ["减脂", "高血压"]
	// ForbiddenUsers: 谁绝对不能吃 (黑名单) -> 例如 ["糖尿病", "痛风"] -> 过滤user的HealthConditions(健康状况)
	TargetUsers    []string `gorm:"serializer:json;index;comment:适用人群" json:"targetUsers"`
	ForbiddenUsers []string `gorm:"serializer:json;index;comment:禁忌人群" json:"forbiddenUsers"`

	// TODO 暂时没有用于逻辑判断
	HealthScore int `gorm:"default:5;comment:推荐评分1-10" json:"healthScore"`
}

func (Recipe) TableName() string {
	return "recipes"
}
