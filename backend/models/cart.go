package models

type Cart struct {
	BaseModel
	UserID    uint `json:"user_id" gorm:"not null;index;comment:用户ID"`
	VehicleID uint `json:"vehicle_id" gorm:"not null;index;comment:电动车ID"`
	Quantity  int  `json:"quantity" gorm:"default:1;comment:数量"`
	Selected  int  `json:"selected" gorm:"default:1;comment:是否选中 1:选中 0:未选中"`
}

func (Cart) TableName() string {
	return "carts"
}
