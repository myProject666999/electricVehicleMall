package models

type Admin struct {
	BaseModel
	Username string `json:"username" gorm:"uniqueIndex;size:50;not null;comment:用户名"`
	Password string `json:"-" gorm:"size:255;not null;comment:密码"`
	RealName string `json:"real_name" gorm:"size:50;comment:真实姓名"`
	Role     int    `json:"role" gorm:"default:2;comment:角色 1:超级管理员 2:普通管理员"`
	Status   int    `json:"status" gorm:"default:1;comment:状态 1:正常 0:禁用"`
}

func (Admin) TableName() string {
	return "admins"
}
