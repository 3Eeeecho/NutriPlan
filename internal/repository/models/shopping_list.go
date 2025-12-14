package models

import (
	"time"

	"gorm.io/gorm"
)

// ShoppingList 购物清单模型
type ShoppingList struct {
	gorm.Model
	UserID     uint      `gorm:"not null;index;comment:用户ID" json:"userId"`
	Name       string    `gorm:"type:varchar(100);not null;comment:清单名称" json:"name"`
	RecipeIDs  string    `gorm:"type:text;comment:关联的食谱ID列表(JSON)" json:"recipeIds"`
	Items      string    `gorm:"type:text;comment:清单项(JSON)" json:"items"`
	Status     string    `gorm:"type:varchar(20);default:'pending';comment:状态(pending/completed)" json:"status"`
	ShopDate   time.Time `gorm:"comment:计划采购日期" json:"shopDate"`
	TotalItems int       `gorm:"comment:总项数" json:"totalItems"`
}

func (ShoppingList) TableName() string {
	return "shopping_lists"
}
