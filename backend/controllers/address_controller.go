package controllers

import (
	"electricVehicleMall/database"
	"electricVehicleMall/models"
	"electricVehicleMall/utils"

	"github.com/gin-gonic/gin"
)

type AddressController struct{}

type CreateAddressRequest struct {
	Consignee string `json:"consignee" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Province  string `json:"province" binding:"required"`
	City      string `json:"city" binding:"required"`
	District  string `json:"district" binding:"required"`
	Address   string `json:"address" binding:"required"`
	IsDefault int    `json:"is_default"`
}

func (ac *AddressController) List(c *gin.Context) {
	userID := c.GetUint("user_id")

	var addresses []models.Address
	database.GetDB().Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&addresses)

	utils.Success(c, addresses)
}

func (ac *AddressController) Detail(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	var address models.Address
	if err := database.GetDB().Where("id = ? AND user_id = ?", id, userID).First(&address).Error; err != nil {
		utils.NotFound(c, "地址不存在")
		return
	}

	utils.Success(c, address)
}

func (ac *AddressController) Create(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 开始事务
	tx := database.GetDB().Begin()

	// 如果设置为默认地址，先取消其他地址的默认状态
	if req.IsDefault == 1 {
		tx.Model(&models.Address{}).Where("user_id = ?", userID).Update("is_default", 0)
	}

	// 检查是否是第一个地址，如果是则设为默认
	var count int64
	tx.Model(&models.Address{}).Where("user_id = ?", userID).Count(&count)
	if count == 0 {
		req.IsDefault = 1
	}

	address := models.Address{
		UserID:    userID,
		Consignee: req.Consignee,
		Phone:     req.Phone,
		Province:  req.Province,
		City:      req.City,
		District:  req.District,
		Address:   req.Address,
		IsDefault: req.IsDefault,
	}

	if err := tx.Create(&address).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	// 提交事务
	tx.Commit()

	utils.SuccessWithMessage(c, "创建成功", gin.H{
		"id": address.ID,
	})
}

func (ac *AddressController) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	var req CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var address models.Address
	if err := database.GetDB().Where("id = ? AND user_id = ?", id, userID).First(&address).Error; err != nil {
		utils.NotFound(c, "地址不存在")
		return
	}

	// 开始事务
	tx := database.GetDB().Begin()

	// 如果设置为默认地址，先取消其他地址的默认状态
	if req.IsDefault == 1 {
		tx.Model(&models.Address{}).Where("user_id = ? AND id != ?", userID, id).Update("is_default", 0)
	}

	updates := map[string]interface{}{
		"consignee":  req.Consignee,
		"phone":      req.Phone,
		"province":   req.Province,
		"city":       req.City,
		"district":   req.District,
		"address":    req.Address,
		"is_default": req.IsDefault,
	}

	if err := tx.Model(&address).Updates(updates).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	// 提交事务
	tx.Commit()

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (ac *AddressController) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	var address models.Address
	if err := database.GetDB().Where("id = ? AND user_id = ?", id, userID).First(&address).Error; err != nil {
		utils.NotFound(c, "地址不存在")
		return
	}

	// 开始事务
	tx := database.GetDB().Begin()

	// 如果删除的是默认地址，需要设置其他地址为默认
	if address.IsDefault == 1 {
		var otherAddress models.Address
		if err := tx.Where("user_id = ? AND id != ?", userID, id).Order("created_at ASC").First(&otherAddress).Error; err == nil {
			tx.Model(&otherAddress).Update("is_default", 1)
		}
	}

	if err := tx.Delete(&address).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "删除失败: "+err.Error())
		return
	}

	// 提交事务
	tx.Commit()

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (ac *AddressController) SetDefault(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	var address models.Address
	if err := database.GetDB().Where("id = ? AND user_id = ?", id, userID).First(&address).Error; err != nil {
		utils.NotFound(c, "地址不存在")
		return
	}

	// 开始事务
	tx := database.GetDB().Begin()

	// 取消其他地址的默认状态
	tx.Model(&models.Address{}).Where("user_id = ?", userID).Update("is_default", 0)

	// 设置当前地址为默认
	tx.Model(&address).Update("is_default", 1)

	// 提交事务
	tx.Commit()

	utils.SuccessWithMessage(c, "设置成功", nil)
}
