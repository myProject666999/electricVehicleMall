package models

type Notice struct {
	BaseModel
	Title   string `json:"title" gorm:"size:255;not null;comment:标题"`
	Content string `json:"content" gorm:"type:text;not null;comment:内容"`
	Author  string `json:"author" gorm:"size:50;comment:作者"`
	Sort    int    `json:"sort" gorm:"default:0;comment:排序"`
	Status  int    `json:"status" gorm:"default:1;comment:状态 1:显示 0:隐藏"`
	Views   int    `json:"views" gorm:"default:0;comment:浏览量"`
}

func (Notice) TableName() string {
	return "notices"
}
