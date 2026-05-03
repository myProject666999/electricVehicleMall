package models

type User struct {
	BaseModel
	Username string `json:"username" gorm:"uniqueIndex;size:50;not null;comment:用户名"`
	Password string `json:"-" gorm:"size:255;not null;comment:密码"`
	Email    string `json:"email" gorm:"size:100;comment:邮箱"`
	Phone    string `json:"phone" gorm:"size:20;comment:手机号"`
	Avatar   string `json:"avatar" gorm:"size:255;comment:头像"`
	Status   int    `json:"status" gorm:"default:1;comment:状态 1:正常 0:禁用"`
}

func (User) TableName() string {
	return "users"
}
