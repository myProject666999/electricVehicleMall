package models

type Comment struct {
	BaseModel
	UserID    uint   `json:"user_id" gorm:"not null;index;comment:用户ID"`
	VehicleID uint   `json:"vehicle_id" gorm:"not null;index;comment:电动车ID"`
	OrderID   uint   `json:"order_id" gorm:"comment:订单ID"`
	Content   string `json:"content" gorm:"type:text;not null;comment:评论内容"`
	Rating    int    `json:"rating" gorm:"default:5;comment:评分 1-5"`
	Images    string `json:"images" gorm:"type:text;comment:图片，JSON格式"`
	Status    int    `json:"status" gorm:"default:1;comment:状态 1:显示 0:隐藏"`
}

func (Comment) TableName() string {
	return "comments"
}
