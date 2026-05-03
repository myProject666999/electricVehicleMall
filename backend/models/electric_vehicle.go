package models

type ElectricVehicle struct {
	BaseModel
	CategoryID     uint    `json:"category_id" gorm:"not null;comment:分类ID"`
	Name           string  `json:"name" gorm:"size:255;not null;comment:名称"`
	Brand          string  `json:"brand" gorm:"size:100;comment:品牌"`
	Price          float64 `json:"price" gorm:"type:decimal(10,2);not null;comment:价格"`
	OriginalPrice  float64 `json:"original_price" gorm:"type:decimal(10,2);comment:原价"`
	Stock          int     `json:"stock" gorm:"default:0;comment:库存"`
	Sales          int     `json:"sales" gorm:"default:0;comment:销量"`
	Image          string  `json:"image" gorm:"size:255;comment:主图"`
	Images         string  `json:"images" gorm:"type:text;comment:轮播图，JSON格式"`
	Description    string  `json:"description" gorm:"type:text;comment:详情描述"`
	Specifications string  `json:"specifications" gorm:"type:text;comment:规格参数，JSON格式"`
	Status         int     `json:"status" gorm:"default:1;comment:状态 1:上架 0:下架"`
	Sort           int     `json:"sort" gorm:"default:0;comment:排序"`
	IsRecommend    int     `json:"is_recommend" gorm:"default:0;comment:是否推荐 1:推荐 0:不推荐"`
}

func (ElectricVehicle) TableName() string {
	return "electric_vehicles"
}
