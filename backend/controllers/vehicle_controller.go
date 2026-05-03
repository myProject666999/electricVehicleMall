package controllers

import (
	"electricVehicleMall/database"
	"electricVehicleMall/models"
	"electricVehicleMall/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

type VehicleController struct{}

type CreateVehicleRequest struct {
	CategoryID     uint    `json:"category_id" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	Brand          string  `json:"brand"`
	Price          float64 `json:"price" binding:"required"`
	OriginalPrice  float64 `json:"original_price"`
	Stock          int     `json:"stock"`
	Image          string  `json:"image"`
	Images         string  `json:"images"`
	Description    string  `json:"description"`
	Specifications string  `json:"specifications"`
	Status         int     `json:"status"`
	Sort           int     `json:"sort"`
	IsRecommend    int     `json:"is_recommend"`
}

func (vc *VehicleController) List(c *gin.Context) {
	var vehicles []models.ElectricVehicle
	var total int64

	query := database.GetDB().Model(&models.ElectricVehicle{})

	// 搜索
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR brand LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 分类筛选
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// 状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 是否推荐
	if isRecommend := c.Query("is_recommend"); isRecommend != "" {
		query = query.Where("is_recommend = ?", isRecommend)
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
	query.Offset(offset).Limit(pageSize).Order("sort DESC, created_at DESC").Find(&vehicles)

	utils.Success(c, gin.H{
		"list":      vehicles,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (vc *VehicleController) Detail(c *gin.Context) {
	id := c.Param("id")

	var vehicle models.ElectricVehicle
	if err := database.GetDB().First(&vehicle, id).Error; err != nil {
		utils.NotFound(c, "电动车不存在")
		return
	}

	// 增加销量
	database.GetDB().Model(&vehicle).Update("sales", vehicle.Sales+1)

	utils.Success(c, vehicle)
}

func (vc *VehicleController) Create(c *gin.Context) {
	var req CreateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	vehicle := models.ElectricVehicle{
		CategoryID:     req.CategoryID,
		Name:           req.Name,
		Brand:          req.Brand,
		Price:          req.Price,
		OriginalPrice:  req.OriginalPrice,
		Stock:          req.Stock,
		Image:          req.Image,
		Images:         req.Images,
		Description:    req.Description,
		Specifications: req.Specifications,
		Status:         req.Status,
		Sort:           req.Sort,
		IsRecommend:    req.IsRecommend,
	}

	if vehicle.Status == 0 {
		vehicle.Status = 1
	}

	if err := database.GetDB().Create(&vehicle).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{
		"id": vehicle.ID,
	})
}

func (vc *VehicleController) Update(c *gin.Context) {
	id := c.Param("id")

	var req CreateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var vehicle models.ElectricVehicle
	if err := database.GetDB().First(&vehicle, id).Error; err != nil {
		utils.NotFound(c, "电动车不存在")
		return
	}

	updates := map[string]interface{}{
		"category_id":     req.CategoryID,
		"name":            req.Name,
		"brand":           req.Brand,
		"price":           req.Price,
		"original_price":  req.OriginalPrice,
		"stock":           req.Stock,
		"image":           req.Image,
		"images":          req.Images,
		"description":     req.Description,
		"specifications":  req.Specifications,
		"status":          req.Status,
		"sort":            req.Sort,
		"is_recommend":    req.IsRecommend,
	}

	if err := database.GetDB().Model(&vehicle).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (vc *VehicleController) Delete(c *gin.Context) {
	id := c.Param("id")

	var vehicle models.ElectricVehicle
	if err := database.GetDB().First(&vehicle, id).Error; err != nil {
		utils.NotFound(c, "电动车不存在")
		return
	}

	if err := database.GetDB().Delete(&vehicle).Error; err != nil {
		utils.InternalServerError(c, "删除失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}
