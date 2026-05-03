package models

type Category struct {
	BaseModel
	Name     string `json:"name" gorm:"size:100;not null;comment:分类名称"`
	ParentID uint   `json:"parent_id" gorm:"default:0;comment:父分类ID"`
	Sort     int    `json:"sort" gorm:"default:0;comment:排序"`
	Status   int    `json:"status" gorm:"default:1;comment:状态 1:启用 0:禁用"`
}

func (Category) TableName() string {
	return "categories"
}
