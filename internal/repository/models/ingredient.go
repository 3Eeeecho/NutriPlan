package models

import "gorm.io/gorm"

// Ingredient 食材基础信息
type Ingredient struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100);not null;uniqueIndex;comment:食材名称" json:"name"`
	Category string `gorm:"type:varchar(50);comment:食材分类(蔬菜/肉类/谷物等)" json:"category"`

	// 营养素 (每100g)
	Energy       float64 `gorm:"type:decimal(10,2);default:0;comment:热量(kcal/100g)" json:"energy"`
	Protein      float64 `gorm:"type:decimal(10,2);default:0;comment:蛋白质(g/100g)" json:"protein"`
	Carbohydrate float64 `gorm:"type:decimal(10,2);default:0;comment:碳水(g/100g)" json:"carbohydrate"`
	Fat          float64 `gorm:"type:decimal(10,2);default:0;comment:脂肪(g/100g)" json:"fat"`
}

func (Ingredient) TableName() string {
	return "ingredients"
}
