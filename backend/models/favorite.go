package models

type Favorite struct {
	BaseModel
	UserID    uint `json:"user_id" gorm:"not null;index;comment:用户ID"`
	VehicleID uint `json:"vehicle_id" gorm:"not null;index;comment:电动车ID"`
}

func (Favorite) TableName() string {
	return "favorites"
}
