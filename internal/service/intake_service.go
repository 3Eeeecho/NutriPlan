package service

import (
	"NutriPlan/internal/repository/dao"
	"NutriPlan/internal/repository/models"
	"errors"
	"time"
)

// IntakeService 饮食记录服务接口
type IntakeService interface {
	// 添加饮食记录
	AddIntakeRecord(record *models.DailyIntakeRecord) error
	// 删除饮食记录
	DeleteIntakeRecord(userID, recordID uint) error
	// 获取当日饮食记录列表
	GetTodayRecords(userID uint) ([]models.DailyIntakeRecord, error)
	// 获取当日营养汇总和达标率
	GetTodayNutritionStatus(userID uint) (*NutritionStatus, error)
	// 获取一周的营养趋势报告
	GetWeeklyReport(userID uint) (*WeeklyReport, error)
}

// NutritionStatus 当日营养状态
type NutritionStatus struct {
	Date               time.Time                  `json:"date"`
	Records            []models.DailyIntakeRecord `json:"records"`
	TotalEnergy        float64                    `json:"total_energy"`
	TotalProtein       float64                    `json:"total_protein"`
	TotalCarbohydrate  float64                    `json:"total_carbohydrate"`
	TotalFat           float64                    `json:"total_fat"`
	TargetEnergy       float64                    `json:"target_energy"`
	TargetProtein      float64                    `json:"target_protein"`
	TargetCarbohydrate float64                    `json:"target_carbohydrate"`
	TargetFat          float64                    `json:"target_fat"`
	EnergyRate         float64                    `json:"energy_rate"`       // 热量达标率
	ProteinRate        float64                    `json:"protein_rate"`      // 蛋白质达标率
	CarbohydrateRate   float64                    `json:"carbohydrate_rate"` // 碳水达标率
	FatRate            float64                    `json:"fat_rate"`          // 脂肪达标率
}

// WeeklyReport 周报告
type WeeklyReport struct {
	StartDate  time.Time            `json:"start_date"`
	EndDate    time.Time            `json:"end_date"`
	DailyData  []DailyNutritionData `json:"daily_data"`
	AvgEnergy  float64              `json:"avg_energy"`
	AvgProtein float64              `json:"avg_protein"`
	AvgCarb    float64              `json:"avg_carb"`
	AvgFat     float64              `json:"avg_fat"`
}

// DailyNutritionData 每日营养数据
type DailyNutritionData struct {
	Date               string  `json:"date"`
	TotalEnergy        float64 `json:"total_energy"`
	TotalProtein       float64 `json:"total_protein"`
	TotalCarbohydrate  float64 `json:"total_carbohydrate"`
	TotalFat           float64 `json:"total_fat"`
	TargetEnergy       float64 `json:"target_energy"`
	TargetProtein      float64 `json:"target_protein"`
	TargetCarbohydrate float64 `json:"target_carbohydrate"`
	TargetFat          float64 `json:"target_fat"`
}

type intakeServiceImpl struct {
	intakeRepo dao.IntakeRepository
	userRepo   dao.UserRepository
	nutriSvc   NutriService
}

// NewIntakeService 创建饮食记录服务实例
func NewIntakeService(intakeRepo dao.IntakeRepository, userRepo dao.UserRepository, nutriSvc NutriService) IntakeService {
	return &intakeServiceImpl{
		intakeRepo: intakeRepo,
		userRepo:   userRepo,
		nutriSvc:   nutriSvc,
	}
}

// AddIntakeRecord 添加饮食记录
func (s *intakeServiceImpl) AddIntakeRecord(record *models.DailyIntakeRecord) error {
	if record.UserID == 0 {
		return errors.New("用户ID不能为空")
	}
	if record.FoodName == "" {
		return errors.New("食物名称不能为空")
	}
	if record.IntakeAmount <= 0 {
		return errors.New("摄入量必须大于0")
	}

	return s.intakeRepo.Create(record)
}

// DeleteIntakeRecord 删除饮食记录
func (s *intakeServiceImpl) DeleteIntakeRecord(userID, recordID uint) error {
	// 可以添加权限检查：确保记录属于该用户
	return s.intakeRepo.Delete(recordID)
}

// GetTodayRecords 获取当日饮食记录列表
func (s *intakeServiceImpl) GetTodayRecords(userID uint) ([]models.DailyIntakeRecord, error) {
	today := time.Now()
	return s.intakeRepo.GetByUserAndDate(userID, today)
}

// GetTodayNutritionStatus 获取当日营养汇总和达标率
func (s *intakeServiceImpl) GetTodayNutritionStatus(userID uint) (*NutritionStatus, error) {
	today := time.Now()

	// 获取用户信息和目标营养
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	targets := s.nutriSvc.CalculateNutritionTargets(user)

	// 获取今日记录
	records, err := s.intakeRepo.GetByUserAndDate(userID, today)
	if err != nil {
		return nil, err
	}

	// 获取今日汇总
	summary, err := s.intakeRepo.GetDailySummary(userID, today)
	if err != nil {
		return nil, err
	}

	// 计算达标率
	status := &NutritionStatus{
		Date:               today,
		Records:            records,
		TotalEnergy:        summary.TotalEnergy,
		TotalProtein:       summary.TotalProtein,
		TotalCarbohydrate:  summary.TotalCarbohydrate,
		TotalFat:           summary.TotalFat,
		TargetEnergy:       targets.DailyEnergy,
		TargetProtein:      targets.DailyProtein,
		TargetCarbohydrate: targets.DailyCarbohydrate,
		TargetFat:          targets.DailyFat,
	}

	// 计算达标率（百分比）
	if targets.DailyEnergy > 0 {
		status.EnergyRate = (summary.TotalEnergy / targets.DailyEnergy) * 100
	}
	if targets.DailyProtein > 0 {
		status.ProteinRate = (summary.TotalProtein / targets.DailyProtein) * 100
	}
	if targets.DailyCarbohydrate > 0 {
		status.CarbohydrateRate = (summary.TotalCarbohydrate / targets.DailyCarbohydrate) * 100
	}
	if targets.DailyFat > 0 {
		status.FatRate = (summary.TotalFat / targets.DailyFat) * 100
	}

	return status, nil
}

// GetWeeklyReport 获取一周的营养趋势报告
func (s *intakeServiceImpl) GetWeeklyReport(userID uint) (*WeeklyReport, error) {
	// 计算本周的起止日期（周一到周日）
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // 将周日从0改为7
	}
	startDate := now.AddDate(0, 0, -(weekday - 1))
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDate := startDate.AddDate(0, 0, 6)

	// 获取用户信息和目标营养
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	targets := s.nutriSvc.CalculateNutritionTargets(user)

	// 获取一周的汇总数据
	summaries, err := s.intakeRepo.GetWeeklySummary(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 构建每日数据
	dailyData := make([]DailyNutritionData, 0)
	var totalEnergy, totalProtein, totalCarb, totalFat float64
	dayCount := 0

	// 创建日期映射
	summaryMap := make(map[string]*dao.DailyNutritionSummary)
	for i := range summaries {
		dateKey := summaries[i].Date.Format("2006-01-02")
		summaryMap[dateKey] = &summaries[i]
	}

	// 遍历7天，填充数据
	for i := 0; i < 7; i++ {
		currentDate := startDate.AddDate(0, 0, i)
		dateKey := currentDate.Format("2006-01-02")

		data := DailyNutritionData{
			Date:               dateKey,
			TargetEnergy:       targets.DailyEnergy,
			TargetProtein:      targets.DailyProtein,
			TargetCarbohydrate: targets.DailyCarbohydrate,
			TargetFat:          targets.DailyFat,
		}

		if summary, exists := summaryMap[dateKey]; exists {
			data.TotalEnergy = summary.TotalEnergy
			data.TotalProtein = summary.TotalProtein
			data.TotalCarbohydrate = summary.TotalCarbohydrate
			data.TotalFat = summary.TotalFat

			totalEnergy += summary.TotalEnergy
			totalProtein += summary.TotalProtein
			totalCarb += summary.TotalCarbohydrate
			totalFat += summary.TotalFat
			dayCount++
		}

		dailyData = append(dailyData, data)
	}

	// 计算平均值
	report := &WeeklyReport{
		StartDate: startDate,
		EndDate:   endDate,
		DailyData: dailyData,
	}

	if dayCount > 0 {
		report.AvgEnergy = totalEnergy / float64(dayCount)
		report.AvgProtein = totalProtein / float64(dayCount)
		report.AvgCarb = totalCarb / float64(dayCount)
		report.AvgFat = totalFat / float64(dayCount)
	}

	return report, nil
}
