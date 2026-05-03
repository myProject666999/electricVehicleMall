package models

type Banner struct {
	BaseModel
	Title   string `json:"title" gorm:"size:255;not null;comment:标题"`
	Image   string `json:"image" gorm:"size:255;not null;comment:图片"`
	Link    string `json:"link" gorm:"size:255;comment:链接地址"`
	Sort    int    `json:"sort" gorm:"default:0;comment:排序"`
	Status  int    `json:"status" gorm:"default:1;comment:状态 1:启用 0:禁用"`
	VehicleID uint  `json:"vehicle_id" gorm:"comment:关联电动车ID"`
}

func (Banner) TableName() string {
	return "banners"
}
