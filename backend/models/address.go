package models

type Address struct {
	BaseModel
	UserID    uint   `json:"user_id" gorm:"not null;index;comment:用户ID"`
	Consignee string `json:"consignee" gorm:"size:50;not null;comment:收货人"`
	Phone     string `json:"phone" gorm:"size:20;not null;comment:联系电话"`
	Province  string `json:"province" gorm:"size:50;comment:省份"`
	City      string `json:"city" gorm:"size:50;comment:城市"`
	District  string `json:"district" gorm:"size:50;comment:区县"`
	Address   string `json:"address" gorm:"size:500;not null;comment:详细地址"`
	IsDefault int    `json:"is_default" gorm:"default:0;comment:是否默认地址 1:默认 0:非默认"`
}

func (Address) TableName() string {
	return "addresses"
}
