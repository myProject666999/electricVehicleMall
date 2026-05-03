package models

type OrderItem struct {
	BaseModel
	OrderID     uint    `json:"order_id" gorm:"not null;index;comment:订单ID"`
	OrderNo     string  `json:"order_no" gorm:"size:50;not null;comment:订单号"`
	VehicleID   uint    `json:"vehicle_id" gorm:"not null;comment:电动车ID"`
	VehicleName string  `json:"vehicle_name" gorm:"size:255;not null;comment:电动车名称"`
	Image       string  `json:"image" gorm:"size:255;comment:商品图片"`
	Price       float64 `json:"price" gorm:"type:decimal(10,2);not null;comment:单价"`
	Quantity    int     `json:"quantity" gorm:"not null;comment:数量"`
	TotalPrice  float64 `json:"total_price" gorm:"type:decimal(10,2);not null;comment:小计"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
