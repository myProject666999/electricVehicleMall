package controllers

import (
	"electricVehicleMall/database"
	"electricVehicleMall/models"
	"electricVehicleMall/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

type NoticeController struct{}

type CreateNoticeRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Author  string `json:"author"`
	Sort    int    `json:"sort"`
	Status  int    `json:"status"`
}

func (nc *NoticeController) List(c *gin.Context) {
	var notices []models.Notice
	var total int64

	query := database.GetDB().Model(&models.Notice{})

	// 搜索
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	// 状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 统计总数
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
	query.Offset(offset).Limit(pageSize).Order("sort DESC, created_at DESC").Find(&notices)

	utils.Success(c, gin.H{
		"list":      notices,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (nc *NoticeController) Detail(c *gin.Context) {
	id := c.Param("id")

	var notice models.Notice
	if err := database.GetDB().First(&notice, id).Error; err != nil {
		utils.NotFound(c, "公告不存在")
		return
	}

	// 增加浏览量
	database.GetDB().Model(&notice).Update("views", notice.Views+1)

	utils.Success(c, notice)
}

func (nc *NoticeController) Create(c *gin.Context) {
	var req CreateNoticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	notice := models.Notice{
		Title:   req.Title,
		Content: req.Content,
		Author:  req.Author,
		Sort:    req.Sort,
		Status:  req.Status,
	}

	if notice.Status == 0 {
		notice.Status = 1
	}

	if notice.Author == "" {
		notice.Author = "管理员"
	}

	if err := database.GetDB().Create(&notice).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{
		"id": notice.ID,
	})
}

func (nc *NoticeController) Update(c *gin.Context) {
	id := c.Param("id")

	var req CreateNoticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var notice models.Notice
	if err := database.GetDB().First(&notice, id).Error; err != nil {
		utils.NotFound(c, "公告不存在")
		return
	}

	updates := map[string]interface{}{
		"title":   req.Title,
		"content": req.Content,
		"author":  req.Author,
		"sort":    req.Sort,
		"status":  req.Status,
	}

	if err := database.GetDB().Model(&notice).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (nc *NoticeController) Delete(c *gin.Context) {
	id := c.Param("id")

	var notice models.Notice
	if err := database.GetDB().First(&notice, id).Error; err != nil {
		utils.NotFound(c, "公告不存在")
		return
	}

	if err := database.GetDB().Delete(&notice).Error; err != nil {
		utils.InternalServerError(c, "删除失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}
