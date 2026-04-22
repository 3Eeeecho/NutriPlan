package models

import (
	"gorm.io/gorm"
)

type Recipe struct {
	gorm.Model

	Name             string   `gorm:"type:varchar(200);not null;index;comment:椋熻氨鍚嶇О" json:"name"`
	ImageURL         string   `gorm:"type:varchar(500);comment:灏侀潰鍥剧墖閾炬帴" json:"imageUrl"`
	MealType         MealType `gorm:"type:varchar(20);not null;index;comment:閫傚悎椁愭绫诲瀷" json:"mealType"`
	AllowedMealTypes []string `gorm:"serializer:json;type:text;comment:allowed meal types" json:"allowedMealTypes"`
	Difficulty       string   `gorm:"type:varchar(20);comment:鐑归オ闅惧害" json:"difficulty"`
	CookingTime      int      `gorm:"default:0;comment:鐑归ऑ鑰楁椂(鍒嗛挓)" json:"cookingTime"`
	PortionWeightG   float64  `gorm:"type:decimal(8,2);default:100;comment:鎺ㄨ崘椋熺敤閲嶉噺(g)" json:"portionWeightG"`

	Energy       float64 `gorm:"type:decimal(10,2);not null;index;comment:姣?00g鐑噺(kcal)" json:"energy"`
	Protein      float64 `gorm:"type:decimal(10,2);not null;comment:姣?00g铔嬬櫧璐?g)" json:"protein"`
	Carbohydrate float64 `gorm:"type:decimal(10,2);not null;comment:姣?00g纰虫按鍖栧悎鐗?g)" json:"carbohydrate"`
	Fat          float64 `gorm:"type:decimal(10,2);not null;comment:姣?00g鑴傝偑(g)" json:"fat"`

	Ingredients  []string `gorm:"serializer:json;not null;comment:鎵€闇€椋熸潗娓呭崟" json:"ingredients"`
	CookingSteps []string `gorm:"serializer:json;comment:鐑归ऑ姝ラ" json:"cookingSteps"`

	TargetUsers    []string `gorm:"serializer:json;index;comment:閫傜敤浜虹兢" json:"targetUsers"`
	ForbiddenUsers []string `gorm:"serializer:json;index;comment:绂佸繉浜虹兢" json:"forbiddenUsers"`

	IsWeightLossFriendly   bool `gorm:"not null;default:false;index;comment:鏄惁閫傚悎鍑忚剛" json:"isWeightLossFriendly"`
	IsMuscleGainFriendly   bool `gorm:"not null;default:false;index;comment:鏄惁閫傚悎澧炶倢" json:"isMuscleGainFriendly"`
	IsSugarControlFriendly bool `gorm:"not null;default:false;index;comment:鏄惁閫傚悎鎺х硸" json:"isSugarControlFriendly"`
	IsGeneralFriendly      bool `gorm:"not null;default:true;index;comment:鏄惁灞炰簬鏅€氬彲鎺ㄨ崘椋熷搧" json:"isGeneralFriendly"`
}

func (Recipe) TableName() string {
	return "recipes"
}
