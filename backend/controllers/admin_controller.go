package controllers

import (
	"electricVehicleMall/database"
	"electricVehicleMall/models"
	"electricVehicleMall/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

type AdminController struct{}

type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateAdminRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	RealName string `json:"real_name"`
	Role     int    `json:"role"`
}

type UpdateAdminRequest struct {
	RealName string `json:"real_name"`
	Role     *int   `json:"role"`
	Status   *int   `json:"status"`
}

func (ac *AdminController) Login(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var admin models.Admin
	if err := database.GetDB().Where("username = ?", req.Username).First(&admin).Error; err != nil {
		utils.BadRequest(c, "用户名或密码错误")
		return
	}

	// 检查密码
	if !utils.CheckPassword(req.Password, admin.Password) {
		utils.BadRequest(c, "用户名或密码错误")
		return
	}

	// 检查状态
	if admin.Status != 1 {
		utils.Forbidden(c, "账号已被禁用")
		return
	}

	// 生成token
	token, err := utils.GenerateToken(admin.ID, admin.Username, "admin")
	if err != nil {
		utils.InternalServerError(c, "生成token失败")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"admin": gin.H{
			"id":        admin.ID,
			"username":  admin.Username,
			"real_name": admin.RealName,
			"role":      admin.Role,
		},
	})
}

func (ac *AdminController) GetProfile(c *gin.Context) {
	adminID := c.GetUint("user_id")

	var admin models.Admin
	if err := database.GetDB().First(&admin, adminID).Error; err != nil {
		utils.NotFound(c, "管理员不存在")
		return
	}

	utils.Success(c, gin.H{
		"id":        admin.ID,
		"username":  admin.Username,
		"real_name": admin.RealName,
		"role":      admin.Role,
		"status":    admin.Status,
	})
}

func (ac *AdminController) List(c *gin.Context) {
	var admins []models.Admin
	var total int64

	query := database.GetDB().Model(&models.Admin{})

	// 搜索
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("username LIKE ? OR real_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
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
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&admins)

	utils.Success(c, gin.H{
		"list":      admins,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (ac *AdminController) Create(c *gin.Context) {
	var req CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 检查用户名是否已存在
	var existingAdmin models.Admin
	if err := database.GetDB().Where("username = ?", req.Username).First(&existingAdmin).Error; err == nil {
		utils.BadRequest(c, "用户名已存在")
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.InternalServerError(c, "密码加密失败")
		return
	}

	admin := models.Admin{
		Username: req.Username,
		Password: hashedPassword,
		RealName: req.RealName,
		Role:     req.Role,
		Status:   1,
	}

	if admin.Role == 0 {
		admin.Role = 2
	}

	if err := database.GetDB().Create(&admin).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{
		"id": admin.ID,
	})
}

func (ac *AdminController) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var admin models.Admin
	if err := database.GetDB().First(&admin, id).Error; err != nil {
		utils.NotFound(c, "管理员不存在")
		return
	}

	updates := make(map[string]interface{})
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.Role != nil {
		updates["role"] = *req.Role
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := database.GetDB().Model(&admin).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (ac *AdminController) Delete(c *gin.Context) {
	id := c.Param("id")

	var admin models.Admin
	if err := database.GetDB().First(&admin, id).Error; err != nil {
		utils.NotFound(c, "管理员不存在")
		return
	}

	// 超级管理员不能删除
	if admin.Role == 1 {
		utils.BadRequest(c, "超级管理员不能删除")
		return
	}

	if err := database.GetDB().Delete(&admin).Error; err != nil {
		utils.InternalServerError(c, "删除失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (ac *AdminController) ResetPassword(c *gin.Context) {
	id := c.Param("id")

	var admin models.Admin
	if err := database.GetDB().First(&admin, id).Error; err != nil {
		utils.NotFound(c, "管理员不存在")
		return
	}

	// 重置密码为 123456
	hashedPassword, err := utils.HashPassword("123456")
	if err != nil {
		utils.InternalServerError(c, "密码加密失败")
		return
	}

	if err := database.GetDB().Model(&admin).Update("password", hashedPassword).Error; err != nil {
		utils.InternalServerError(c, "重置密码失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "密码已重置为 123456", nil)
}
