package controllers

import (
	"electricVehicleMall/database"
	"electricVehicleMall/models"
	"electricVehicleMall/utils"

	"github.com/gin-gonic/gin"
)

type CategoryController struct{}

type CreateCategoryRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID uint   `json:"parent_id"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status"`
}

func (cc *CategoryController) List(c *gin.Context) {
	var categories []models.Category

	query := database.GetDB().Model(&models.Category{})

	// 状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Order("sort ASC, created_at ASC").Find(&categories)

	utils.Success(c, categories)
}

func (cc *CategoryController) Detail(c *gin.Context) {
	id := c.Param("id")

	var category models.Category
	if err := database.GetDB().First(&category, id).Error; err != nil {
		utils.NotFound(c, "分类不存在")
		return
	}

	utils.Success(c, category)
}

func (cc *CategoryController) Create(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	category := models.Category{
		Name:     req.Name,
		ParentID: req.ParentID,
		Sort:     req.Sort,
		Status:   req.Status,
	}

	if category.Status == 0 {
		category.Status = 1
	}

	if err := database.GetDB().Create(&category).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{
		"id": category.ID,
	})
}

func (cc *CategoryController) Update(c *gin.Context) {
	id := c.Param("id")

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var category models.Category
	if err := database.GetDB().First(&category, id).Error; err != nil {
		utils.NotFound(c, "分类不存在")
		return
	}

	updates := map[string]interface{}{
		"name":      req.Name,
		"parent_id": req.ParentID,
		"sort":      req.Sort,
		"status":    req.Status,
	}

	if err := database.GetDB().Model(&category).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (cc *CategoryController) Delete(c *gin.Context) {
	id := c.Param("id")

	var category models.Category
	if err := database.GetDB().First(&category, id).Error; err != nil {
		utils.NotFound(c, "分类不存在")
		return
	}

	// 检查是否有子分类
	var childCount int64
	database.GetDB().Model(&models.Category{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		utils.BadRequest(c, "该分类下有子分类，无法删除")
		return
	}

	// 检查是否有关联的电动车
	var vehicleCount int64
	database.GetDB().Model(&models.ElectricVehicle{}).Where("category_id = ?", id).Count(&vehicleCount)
	if vehicleCount > 0 {
		utils.BadRequest(c, "该分类下有电动车，无法删除")
		return
	}

	if err := database.GetDB().Delete(&category).Error; err != nil {
		utils.InternalServerError(c, "删除失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}
