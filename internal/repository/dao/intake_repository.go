package dao

import (
	"NutriPlan/internal/repository/models"
	"time"

	"gorm.io/gorm"
)

// IntakeRepository 饮食记录仓储接口
type IntakeRepository interface {
	// 创建饮食记录
	Create(record *models.DailyIntakeRecord) error
	// 删除饮食记录
	Delete(id uint) error
	// 获取某用户某一天的所有记录
	GetByUserAndDate(userID uint, date time.Time) ([]models.DailyIntakeRecord, error)
	// 获取某用户某一周的记录
	GetByUserAndWeek(userID uint, startDate, endDate time.Time) ([]models.DailyIntakeRecord, error)
	// 计算某天的营养总和
	GetDailySummary(userID uint, date time.Time) (*DailyNutritionSummary, error)
	// 获取一周的每日汇总
	GetWeeklySummary(userID uint, startDate, endDate time.Time) ([]DailyNutritionSummary, error)
}

// DailyNutritionSummary 每日营养汇总
type DailyNutritionSummary struct {
	Date              time.Time `json:"date"`
	TotalEnergy       float64   `json:"total_energy"`
	TotalProtein      float64   `json:"total_protein"`
	TotalCarbohydrate float64   `json:"total_carbohydrate"`
	TotalFat          float64   `json:"total_fat"`
	RecordCount       int       `json:"record_count"`
}

type intakeRepositoryImpl struct {
	db *gorm.DB
}

// NewIntakeRepository 创建饮食记录仓储实例
func NewIntakeRepository(db *gorm.DB) IntakeRepository {
	return &intakeRepositoryImpl{db: db}
}

// Create 创建饮食记录
func (r *intakeRepositoryImpl) Create(record *models.DailyIntakeRecord) error {
	// 如果没有设置记录日期，使用当前日期
	if record.RecordDate.IsZero() {
		now := time.Now()
		record.RecordDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	}
	return r.db.Create(record).Error
}

// Delete 删除饮食记录
func (r *intakeRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.DailyIntakeRecord{}, id).Error
}

// GetByUserAndDate 获取某用户某一天的所有记录
func (r *intakeRepositoryImpl) GetByUserAndDate(userID uint, date time.Time) ([]models.DailyIntakeRecord, error) {
	var records []models.DailyIntakeRecord
	// 只比较日期部分
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := r.db.Where("user_id = ? AND record_date >= ? AND record_date < ?", userID, startOfDay, endOfDay).
		Order("created_at DESC").
		Find(&records).Error

	return records, err
}

// GetByUserAndWeek 获取某用户某一周的记录
func (r *intakeRepositoryImpl) GetByUserAndWeek(userID uint, startDate, endDate time.Time) ([]models.DailyIntakeRecord, error) {
	var records []models.DailyIntakeRecord
	err := r.db.Where("user_id = ? AND record_date >= ? AND record_date <= ?", userID, startDate, endDate).
		Order("record_date ASC, created_at DESC").
		Find(&records).Error

	return records, err
}

// GetDailySummary 计算某天的营养总和
func (r *intakeRepositoryImpl) GetDailySummary(userID uint, date time.Time) (*DailyNutritionSummary, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var summary DailyNutritionSummary
	summary.Date = startOfDay

	err := r.db.Model(&models.DailyIntakeRecord{}).
		Select(`
			COALESCE(SUM(calculated_energy), 0) as total_energy,
			COALESCE(SUM(calculated_protein), 0) as total_protein,
			COALESCE(SUM(calculated_carb), 0) as total_carbohydrate,
			COALESCE(SUM(calculated_fat), 0) as total_fat,
			COUNT(*) as record_count
		`).
		Where("user_id = ? AND record_date >= ? AND record_date < ?", userID, startOfDay, endOfDay).
		Scan(&summary).Error

	return &summary, err
}

// GetWeeklySummary 获取一周的每日汇总
func (r *intakeRepositoryImpl) GetWeeklySummary(userID uint, startDate, endDate time.Time) ([]DailyNutritionSummary, error) {
	var summaries []DailyNutritionSummary

	// 按日期分组统计
	err := r.db.Model(&models.DailyIntakeRecord{}).
		Select(`
			DATE(record_date) as date,
			COALESCE(SUM(calculated_energy), 0) as total_energy,
			COALESCE(SUM(calculated_protein), 0) as total_protein,
			COALESCE(SUM(calculated_carb), 0) as total_carbohydrate,
			COALESCE(SUM(calculated_fat), 0) as total_fat,
			COUNT(*) as record_count
		`).
		Where("user_id = ? AND record_date >= ? AND record_date <= ?", userID, startDate, endDate).
		Group("DATE(record_date)").
		Order("date ASC").
		Scan(&summaries).Error

	return summaries, err
}
