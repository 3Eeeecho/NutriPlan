package domain

import (
	"gorm.io/gorm"
)

// Ingredient 对应 food_nutrition 表
type Ingredient struct {
	gorm.Model

	CateID int `gorm:"column:cate_id;default:0" json:"cate_id"` // 分类ID（对应 food_categories.id）

	// GORM 关系定义：关联到 Category 模型
	Category Category `gorm:"foreignKey:CateID"`

	FoodName    string `gorm:"type:varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;not null;index:idx_food_name" json:"food_name"`
	AliasName   string `gorm:"type:varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci" json:"alias_name"`
	EnglishName string `gorm:"type:varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;index:idx_english_name" json:"english_name"`

	// --- 宏量营养素 (转换为 float64 进行计算) ---
	Water        float64 `gorm:"type:decimal(10,2)" json:"water"`
	Energy       float64 `gorm:"type:decimal(10,2)" json:"energy"`        // 能量 (kcal)
	Protein      float64 `gorm:"type:decimal(10,2)" json:"protein"`       // 蛋白质 (g)
	Fat          float64 `gorm:"type:decimal(10,2)" json:"fat"`           // 脂肪 (g)
	Carbohydrate float64 `gorm:"type:decimal(10,2)" json:"carbohydrate"`  // 碳水化合物 (g)
	DietaryFiber float64 `gorm:"type:decimal(10,2)" json:"dietary_fiber"` // 总膳食纤维 (g)
}

// TableName 用于指定 GORM 映射的表名
func (Ingredient) TableName() string {
	return "food_nutrition"
}
