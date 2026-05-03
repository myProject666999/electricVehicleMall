package models

type ChatMessage struct {
	BaseModel
	UserID      uint   `json:"user_id" gorm:"not null;index;comment:用户ID"`
	AdminID     *uint  `json:"admin_id" gorm:"index;comment:管理员ID，空表示系统消息"`
	Content     string `json:"content" gorm:"type:text;not null;comment:消息内容"`
	MessageType int    `json:"message_type" gorm:"default:1;comment:消息类型 1:文本 2:图片 3:系统"`
	IsFromAdmin int    `json:"is_from_admin" gorm:"default:0;comment:是否来自管理员 1:是 0:否"`
	IsRead      int    `json:"is_read" gorm:"default:0;comment:是否已读 1:已读 0:未读"`
}

func (ChatMessage) TableName() string {
	return "chat_messages"
}
