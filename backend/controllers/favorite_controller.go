package controllers

import (
	"electricVehicleMall/database"
	"electricVehicleMall/models"
	"electricVehicleMall/utils"

	"github.com/gin-gonic/gin"
)

type FavoriteController struct{}

func (fc *FavoriteController) List(c *gin.Context) {
	userID := c.GetUint("user_id")

	var favorites []struct {
		models.Favorite
		VehicleName  string  `json:"vehicle_name"`
		VehicleImage string `json:"vehicle_image"`
		Price        float64 `json:"price"`
		OriginalPrice float64 `json:"original_price"`
		Status       int     `json:"vehicle_status"`
	}

	query := database.GetDB().Table("favorites").
		Select("favorites.*, electric_vehicles.name as vehicle_name, electric_vehicles.image as vehicle_image, electric_vehicles.price, electric_vehicles.original_price, electric_vehicles.status as vehicle_status").
		Joins("LEFT JOIN electric_vehicles ON favorites.vehicle_id = electric_vehicles.id").
		Where("favorites.user_id = ?", userID).
		Order("favorites.created_at DESC")

	query.Find(&favorites)

	utils.Success(c, favorites)
}

func (fc *FavoriteController) Toggle(c *gin.Context) {
	userID := c.GetUint("user_id")
	vehicleID := c.Param("vehicle_id")

	// 检查电动车是否存在
	var vehicle models.ElectricVehicle
	if err := database.GetDB().First(&vehicle, vehicleID).Error; err != nil {
		utils.NotFound(c, "电动车不存在")
		return
	}

	// 检查是否已收藏
	var favorite models.Favorite
	if err := database.GetDB().Where("user_id = ? AND vehicle_id = ?", userID, vehicleID).First(&favorite).Error; err == nil {
		// 已收藏，取消收藏
		database.GetDB().Delete(&favorite)
		utils.SuccessWithMessage(c, "已取消收藏", gin.H{
			"is_favorite": false,
		})
		return
	}

	// 未收藏，添加收藏
	newFavorite := models.Favorite{
		UserID:    userID,
		VehicleID: vehicle.ID,
	}

	if err := database.GetDB().Create(&newFavorite).Error; err != nil {
		utils.InternalServerError(c, "收藏失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "已收藏", gin.H{
		"is_favorite": true,
	})
}

func (fc *FavoriteController) Check(c *gin.Context) {
	userID := c.GetUint("user_id")
	vehicleID := c.Param("vehicle_id")

	var favorite models.Favorite
	if err := database.GetDB().Where("user_id = ? AND vehicle_id = ?", userID, vehicleID).First(&favorite).Error; err == nil {
		utils.Success(c, gin.H{
			"is_favorite": true,
		})
		return
	}

	utils.Success(c, gin.H{
		"is_favorite": false,
	})
}

func (fc *FavoriteController) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	var favorite models.Favorite
	if err := database.GetDB().Where("id = ? AND user_id = ?", id, userID).First(&favorite).Error; err != nil {
		utils.NotFound(c, "收藏不存在")
		return
	}

	if err := database.GetDB().Delete(&favorite).Error; err != nil {
		utils.InternalServerError(c, "删除失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}
