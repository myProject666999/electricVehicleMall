package controllers

import (
	"electricVehicleMall/database"
	"electricVehicleMall/models"
	"electricVehicleMall/utils"

	"github.com/gin-gonic/gin"
)

type BannerController struct{}

type CreateBannerRequest struct {
	Title     string `json:"title" binding:"required"`
	Image     string `json:"image" binding:"required"`
	Link      string `json:"link"`
	Sort      int    `json:"sort"`
	Status    int    `json:"status"`
	VehicleID uint   `json:"vehicle_id"`
}

func (bc *BannerController) List(c *gin.Context) {
	var banners []models.Banner

	query := database.GetDB().Model(&models.Banner{})

	// 状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Order("sort ASC, created_at DESC").Find(&banners)

	utils.Success(c, banners)
}

func (bc *BannerController) Detail(c *gin.Context) {
	id := c.Param("id")

	var banner models.Banner
	if err := database.GetDB().First(&banner, id).Error; err != nil {
		utils.NotFound(c, "轮播图不存在")
		return
	}

	utils.Success(c, banner)
}

func (bc *BannerController) Create(c *gin.Context) {
	var req CreateBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	banner := models.Banner{
		Title:     req.Title,
		Image:     req.Image,
		Link:      req.Link,
		Sort:      req.Sort,
		Status:    req.Status,
		VehicleID: req.VehicleID,
	}

	if banner.Status == 0 {
		banner.Status = 1
	}

	if err := database.GetDB().Create(&banner).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{
		"id": banner.ID,
	})
}

func (bc *BannerController) Update(c *gin.Context) {
	id := c.Param("id")

	var req CreateBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var banner models.Banner
	if err := database.GetDB().First(&banner, id).Error; err != nil {
		utils.NotFound(c, "轮播图不存在")
		return
	}

	updates := map[string]interface{}{
		"title":      req.Title,
		"image":      req.Image,
		"link":       req.Link,
		"sort":       req.Sort,
		"status":     req.Status,
		"vehicle_id": req.VehicleID,
	}

	if err := database.GetDB().Model(&banner).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (bc *BannerController) Delete(c *gin.Context) {
	id := c.Param("id")

	var banner models.Banner
	if err := database.GetDB().First(&banner, id).Error; err != nil {
		utils.NotFound(c, "轮播图不存在")
		return
	}

	if err := database.GetDB().Delete(&banner).Error; err != nil {
		utils.InternalServerError(c, "删除失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}
