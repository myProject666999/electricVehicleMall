package controllers

import (
	"electricVehicleMall/database"
	"electricVehicleMall/models"
	"electricVehicleMall/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

type ChatController struct{}

type SendMessageRequest struct {
	Content     string `json:"content" binding:"required"`
	MessageType int    `json:"message_type"`
}

func (cc *ChatController) UserList(c *gin.Context) {
	userID := c.GetUint("user_id")

	var messages []models.ChatMessage

	query := database.GetDB().Where("user_id = ?", userID)

	// 分页
	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := c.Query("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at ASC").Find(&messages)

	// 标记为已读
	database.GetDB().Model(&models.ChatMessage{}).
		Where("user_id = ? AND is_from_admin = 1 AND is_read = 0", userID).
		Update("is_read", 1)

	utils.Success(c, messages)
}

func (cc *ChatController) UserSend(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	messageType := req.MessageType
	if messageType == 0 {
		messageType = 1
	}

	message := models.ChatMessage{
		UserID:      userID,
		Content:     req.Content,
		MessageType: messageType,
		IsFromAdmin: 0,
		IsRead:      1,
	}

	if err := database.GetDB().Create(&message).Error; err != nil {
		utils.InternalServerError(c, "发送失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "发送成功", message)
}

func (cc *ChatController) AdminListUsers(c *gin.Context) {
	// 获取所有有聊天记录的用户列表
	var userChats []struct {
		UserID       uint   `json:"user_id"`
		Username     string `json:"username"`
		LastMessage  string `json:"last_message"`
		LastTime     string `json:"last_time"`
		UnreadCount  int64  `json:"unread_count"`
	}

	// 复杂查询，先获取每个用户的最新消息和未读数量
	query := `
		SELECT 
			cm.user_id,
			u.username,
			cm.content as last_message,
			cm.created_at as last_time,
			(SELECT COUNT(*) FROM chat_messages WHERE user_id = cm.user_id AND is_from_admin = 0 AND is_read = 0) as unread_count
		FROM chat_messages cm
		LEFT JOIN users u ON cm.user_id = u.id
		WHERE cm.created_at = (
			SELECT MAX(created_at) FROM chat_messages WHERE user_id = cm.user_id
		)
		GROUP BY cm.user_id
		ORDER BY cm.created_at DESC
	`

	database.GetDB().Raw(query).Scan(&userChats)

	utils.Success(c, userChats)
}

func (cc *ChatController) AdminListMessages(c *gin.Context) {
	userID := c.Param("user_id")

	if userID == "" {
		userID = c.Query("user_id")
	}

	if userID == "" {
		utils.BadRequest(c, "请指定用户ID")
		return
	}

	var messages []struct {
		models.ChatMessage
		AdminName string `json:"admin_name"`
		UserName  string `json:"user_name"`
	}

	query := database.GetDB().Table("chat_messages").
		Select("chat_messages.*, admins.real_name as admin_name, users.username as user_name").
		Joins("LEFT JOIN admins ON chat_messages.admin_id = admins.id").
		Joins("LEFT JOIN users ON chat_messages.user_id = users.id").
		Where("chat_messages.user_id = ?", userID)

	// 分页
	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := c.Query("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("chat_messages.created_at ASC").Find(&messages)

	// 标记用户消息为已读
	database.GetDB().Model(&models.ChatMessage{}).
		Where("user_id = ? AND is_from_admin = 0 AND is_read = 0", userID).
		Update("is_read", 1)

	utils.Success(c, messages)
}

func (cc *ChatController) AdminSend(c *gin.Context) {
	adminID := c.GetUint("user_id")

	var req struct {
		UserID      uint   `json:"user_id" binding:"required"`
		Content     string `json:"content" binding:"required"`
		MessageType int    `json:"message_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	messageType := req.MessageType
	if messageType == 0 {
		messageType = 1
	}

	message := models.ChatMessage{
		UserID:      req.UserID,
		AdminID:     &adminID,
		Content:     req.Content,
		MessageType: messageType,
		IsFromAdmin: 1,
		IsRead:      0,
	}

	if err := database.GetDB().Create(&message).Error; err != nil {
		utils.InternalServerError(c, "发送失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "发送成功", message)
}
