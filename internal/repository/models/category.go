package models

import (
	"gorm.io/gorm"
)

// Category 对应 food_categories 表
type Category struct {
	gorm.Model

	Title         string `gorm:"type:varchar(255);not null" json:"title"`         // 食物分类名称
	FatherID      string `gorm:"type:varchar(50)" json:"father_id"`               // 父级分类ID (原表中的cate_id体系)
	CateID        string `gorm:"type:varchar(50)" json:"cate_id"`                 // 当前分类ID (原表中的cate_id体系)
	IsSubcategory bool   `gorm:"type:tinyint(1);default:0" json:"is_subcategory"` // 是否为子分类

	// 用于 GORM 关联的字段
	ParentCategoryID *uint `gorm:"column:parent_category_id" json:"parent_category_id"` // 父分类在数据库中的ID

	// GORM 关系定义：自关联，用于获取子分类
	SubCategories []Category `gorm:"foreignKey:ParentCategoryID"`
}

func (Category) TableName() string {
	return "food_categories"
}
