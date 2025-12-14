package dao

import (
	"NutriPlan/internal/repository/models"

	"gorm.io/gorm"
)

type ShoppingListRepository interface {
	Create(list *models.ShoppingList) error
	GetByID(id uint) (*models.ShoppingList, error)
	GetByUserID(userID uint, page, pageSize int) ([]models.ShoppingList, int64, error)
	Update(list *models.ShoppingList) error
	Delete(id uint) error
	UpdateStatus(id uint, status string) error
}

type GormShoppingListRepository struct {
	db *gorm.DB
}

func NewShoppingListRepository(db *gorm.DB) ShoppingListRepository {
	return &GormShoppingListRepository{db: db}
}

func (r *GormShoppingListRepository) Create(list *models.ShoppingList) error {
	return r.db.Create(list).Error
}

func (r *GormShoppingListRepository) GetByID(id uint) (*models.ShoppingList, error) {
	var list models.ShoppingList
	err := r.db.First(&list, id).Error
	return &list, err
}

func (r *GormShoppingListRepository) GetByUserID(userID uint, page, pageSize int) ([]models.ShoppingList, int64, error) {
	var lists []models.ShoppingList
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.ShoppingList{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at desc").Limit(pageSize).Offset(offset).Find(&lists).Error
	return lists, total, err
}

func (r *GormShoppingListRepository) Update(list *models.ShoppingList) error {
	return r.db.Save(list).Error
}

func (r *GormShoppingListRepository) Delete(id uint) error {
	return r.db.Delete(&models.ShoppingList{}, id).Error
}

func (r *GormShoppingListRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.ShoppingList{}).Where("id = ?", id).Update("status", status).Error
}
