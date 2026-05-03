package controllers

import (
	"electricVehicleMall/database"
	"electricVehicleMall/models"
	"electricVehicleMall/utils"

	"github.com/gin-gonic/gin"
)

type CartController struct{}

type AddCartRequest struct {
	VehicleID uint `json:"vehicle_id" binding:"required"`
	Quantity  int  `json:"quantity"`
}

type UpdateCartRequest struct {
	Quantity int `json:"quantity"`
	Selected *int `json:"selected"`
}

func (cc *CartController) List(c *gin.Context) {
	userID := c.GetUint("user_id")

	var cartItems []struct {
		models.Cart
		VehicleName  string  `json:"vehicle_name"`
		VehicleImage string `json:"vehicle_image"`
		Price        float64 `json:"price"`
		Stock        int     `json:"stock"`
		Status       int     `json:"vehicle_status"`
	}

	query := database.GetDB().Table("carts").
		Select("carts.*, electric_vehicles.name as vehicle_name, electric_vehicles.image as vehicle_image, electric_vehicles.price, electric_vehicles.stock, electric_vehicles.status as vehicle_status").
		Joins("LEFT JOIN electric_vehicles ON carts.vehicle_id = electric_vehicles.id").
		Where("carts.user_id = ?", userID).
		Order("carts.created_at DESC")

	query.Find(&cartItems)

	utils.Success(c, cartItems)
}

func (cc *CartController) Add(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req AddCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 检查电动车是否存在
	var vehicle models.ElectricVehicle
	if err := database.GetDB().First(&vehicle, req.VehicleID).Error; err != nil {
		utils.NotFound(c, "电动车不存在")
		return
	}

	// 检查电动车是否上架
	if vehicle.Status != 1 {
		utils.BadRequest(c, "该电动车已下架")
		return
	}

	// 检查购物车中是否已存在
	var existingCart models.Cart
	if err := database.GetDB().Where("user_id = ? AND vehicle_id = ?", userID, req.VehicleID).First(&existingCart).Error; err == nil {
		// 已存在，更新数量
		newQuantity := existingCart.Quantity + req.Quantity
		if newQuantity <= 0 {
			newQuantity = 1
		}
		database.GetDB().Model(&existingCart).Update("quantity", newQuantity)
		utils.SuccessWithMessage(c, "添加成功", gin.H{
			"id": existingCart.ID,
		})
		return
	}

	// 不存在，创建新的购物车项
	quantity := req.Quantity
	if quantity <= 0 {
		quantity = 1
	}

	cart := models.Cart{
		UserID:    userID,
		VehicleID: req.VehicleID,
		Quantity:  quantity,
		Selected:  1,
	}

	if err := database.GetDB().Create(&cart).Error; err != nil {
		utils.InternalServerError(c, "添加失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "添加成功", gin.H{
		"id": cart.ID,
	})
}

func (cc *CartController) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	var req UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var cart models.Cart
	if err := database.GetDB().Where("id = ? AND user_id = ?", id, userID).First(&cart).Error; err != nil {
		utils.NotFound(c, "购物车项不存在")
		return
	}

	updates := make(map[string]interface{})
	if req.Quantity > 0 {
		updates["quantity"] = req.Quantity
	}
	if req.Selected != nil {
		updates["selected"] = *req.Selected
	}

	if err := database.GetDB().Model(&cart).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (cc *CartController) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	var cart models.Cart
	if err := database.GetDB().Where("id = ? AND user_id = ?", id, userID).First(&cart).Error; err != nil {
		utils.NotFound(c, "购物车项不存在")
		return
	}

	if err := database.GetDB().Delete(&cart).Error; err != nil {
		utils.InternalServerError(c, "删除失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (cc *CartController) Clear(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := database.GetDB().Where("user_id = ?", userID).Delete(&models.Cart{}).Error; err != nil {
		utils.InternalServerError(c, "清空失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "清空成功", nil)
}
