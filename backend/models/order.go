package models

type Order struct {
	BaseModel
	OrderNo      string  `json:"order_no" gorm:"uniqueIndex;size:50;not null;comment:订单号"`
	UserID       uint    `json:"user_id" gorm:"not null;index;comment:用户ID"`
	TotalAmount  float64 `json:"total_amount" gorm:"type:decimal(10,2);not null;comment:订单总金额"`
	PayAmount    float64 `json:"pay_amount" gorm:"type:decimal(10,2);comment:实际支付金额"`
	Status       int     `json:"status" gorm:"default:0;comment:订单状态 0:待支付 1:待发货 2:待收货 3:已完成 4:已取消 5:已退款"`
	PayType      int     `json:"pay_type" gorm:"comment:支付方式 1:微信 2:支付宝 3:银行卡"`
	PayTime      *string `json:"pay_time" gorm:"comment:支付时间"`
	Consignee    string  `json:"consignee" gorm:"size:50;not null;comment:收货人"`
	Phone        string  `json:"phone" gorm:"size:20;not null;comment:联系电话"`
	Address      string  `json:"address" gorm:"size:500;not null;comment:收货地址"`
	LogisticsNo  string  `json:"logistics_no" gorm:"size:50;comment:物流单号"`
	LogisticsCom string  `json:"logistics_com" gorm:"size:100;comment:物流公司"`
	ShipTime     *string `json:"ship_time" gorm:"comment:发货时间"`
	ReceiveTime  *string `json:"receive_time" gorm:"comment:收货时间"`
	Remark       string  `json:"remark" gorm:"size:500;comment:备注"`
}

func (Order) TableName() string {
	return "orders"
}
