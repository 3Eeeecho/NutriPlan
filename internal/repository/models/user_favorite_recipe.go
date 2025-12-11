package models

import (
	"time"

	"gorm.io/gorm"
)

// UserFavoriteRecipe 用户收藏食谱关联表
type UserFavoriteRecipe struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"not null;index:idx_user_recipe,unique" json:"user_id"`
	RecipeID  uint           `gorm:"not null;index:idx_user_recipe,unique" json:"recipe_id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	User   User   `gorm:"foreignKey:UserID" json:"-"`
	Recipe Recipe `gorm:"foreignKey:RecipeID" json:"recipe"`
}

func (UserFavoriteRecipe) TableName() string {
	return "user_favorite_recipes"
}
