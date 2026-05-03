package controllers

import (
	"electricVehicleMall/database"
	"electricVehicleMall/models"
	"electricVehicleMall/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

type CommentController struct{}

type CreateCommentRequest struct {
	VehicleID uint   `json:"vehicle_id" binding:"required"`
	OrderID   uint   `json:"order_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Rating    int    `json:"rating"`
	Images    string `json:"images"`
}

func (cc *CommentController) List(c *gin.Context) {
	vehicleID := c.Query("vehicle_id")

	var comments []struct {
		models.Comment
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
	}

	query := database.GetDB().Table("comments").
		Select("comments.*, users.username, users.avatar").
		Joins("LEFT JOIN users ON comments.user_id = users.id").
		Where("comments.status = 1")

	if vehicleID != "" {
		query = query.Where("comments.vehicle_id = ?", vehicleID)
	}

	// 统计总数
	var total int64
	query.Count(&total)

	// 分页
	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := c.Query("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("comments.created_at DESC").Find(&comments)

	utils.Success(c, gin.H{
		"list":      comments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (cc *CommentController) AdminList(c *gin.Context) {
	var comments []struct {
		models.Comment
		Username    string `json:"username"`
		VehicleName string `json:"vehicle_name"`
	}

	query := database.GetDB().Table("comments").
		Select("comments.*, users.username, electric_vehicles.name as vehicle_name").
		Joins("LEFT JOIN users ON comments.user_id = users.id").
		Joins("LEFT JOIN electric_vehicles ON comments.vehicle_id = electric_vehicles.id")

	// 搜索
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("users.username LIKE ? OR electric_vehicles.name LIKE ? OR comments.content LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("comments.status = ?", status)
	}

	// 统计总数
	var total int64
	query.Count(&total)

	// 分页
	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := c.Query("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("comments.created_at DESC").Find(&comments)

	utils.Success(c, gin.H{
		"list":      comments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (cc *CommentController) Create(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 检查订单是否存在且属于当前用户
	var order models.Order
	if err := database.GetDB().Where("id = ? AND user_id = ?", req.OrderID, userID).First(&order).Error; err != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	// 检查订单状态是否为已完成
	if order.Status != int(OrderStatusCompleted) {
		utils.BadRequest(c, "只有已完成的订单才能评价")
		return
	}

	// 检查是否已评价过
	var existingComment models.Comment
	if err := database.GetDB().Where("order_id = ? AND user_id = ?", req.OrderID, userID).First(&existingComment).Error; err == nil {
		utils.BadRequest(c, "该订单已评价")
		return
	}

	// 检查电动车是否在订单中
	var orderItem models.OrderItem
	if err := database.GetDB().Where("order_id = ? AND vehicle_id = ?", req.OrderID, req.VehicleID).First(&orderItem).Error; err != nil {
		utils.BadRequest(c, "该电动车不在此订单中")
		return
	}

	rating := req.Rating
	if rating < 1 || rating > 5 {
		rating = 5
	}

	comment := models.Comment{
		UserID:    userID,
		VehicleID: req.VehicleID,
		OrderID:   req.OrderID,
		Content:   req.Content,
		Rating:    rating,
		Images:    req.Images,
		Status:    1,
	}

	if err := database.GetDB().Create(&comment).Error; err != nil {
		utils.InternalServerError(c, "评价失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "评价成功", gin.H{
		"id": comment.ID,
	})
}

func (cc *CommentController) UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var comment models.Comment
	if err := database.GetDB().First(&comment, id).Error; err != nil {
		utils.NotFound(c, "评论不存在")
		return
	}

	database.GetDB().Model(&comment).Update("status", req.Status)

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (cc *CommentController) Delete(c *gin.Context) {
	id := c.Param("id")

	var comment models.Comment
	if err := database.GetDB().First(&comment, id).Error; err != nil {
		utils.NotFound(c, "评论不存在")
		return
	}

	if err := database.GetDB().Delete(&comment).Error; err != nil {
		utils.InternalServerError(c, "删除失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}
