package controllers

import (
	"fmt"

	"electricVehicleMall/database"
	"electricVehicleMall/models"
	"electricVehicleMall/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Phone  string `json:"phone"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (uc *UserController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := database.GetDB().Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		utils.BadRequest(c, "用户名已存在")
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.InternalServerError(c, "密码加密失败")
		return
	}

	// 创建用户
	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Phone:    req.Phone,
		Status:   1,
	}

	if err := database.GetDB().Create(&user).Error; err != nil {
		utils.InternalServerError(c, "注册失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "注册成功", gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}

func (uc *UserController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var user models.User
	if err := database.GetDB().Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.BadRequest(c, "用户名或密码错误")
		return
	}

	// 检查密码
	if !utils.CheckPassword(req.Password, user.Password) {
		utils.BadRequest(c, "用户名或密码错误")
		return
	}

	// 检查状态
	if user.Status != 1 {
		utils.Forbidden(c, "账号已被禁用")
		return
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username, "user")
	if err != nil {
		utils.InternalServerError(c, "生成token失败")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"phone":    user.Phone,
			"email":    user.Email,
			"avatar":   user.Avatar,
		},
	})
}

func (uc *UserController) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"phone":    user.Phone,
		"email":    user.Email,
		"avatar":   user.Avatar,
		"status":   user.Status,
	})
}

func (uc *UserController) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if err := database.GetDB().Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (uc *UserController) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var user models.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	// 验证旧密码
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		utils.BadRequest(c, "原密码错误")
		return
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.InternalServerError(c, "密码加密失败")
		return
	}

	// 更新密码
	if err := database.GetDB().Model(&user).Update("password", hashedPassword).Error; err != nil {
		utils.InternalServerError(c, "密码修改失败")
		return
	}

	utils.SuccessWithMessage(c, "密码修改成功", nil)
}

// 管理员获取用户列表
func (uc *UserController) AdminList(c *gin.Context) {
	page := 1
	pageSize := 10

	if p, ok := c.GetQuery("page"); ok {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps, ok := c.GetQuery("page_size"); ok {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	keyword := c.Query("keyword")
	status := c.Query("status")

	query := database.GetDB().Model(&models.User{})

	if keyword != "" {
		query = query.Where("username LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize

	var users []models.User
	query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&users)

	userList := make([]gin.H, 0)
	for _, user := range users {
		userList = append(userList, gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"phone":      user.Phone,
			"email":      user.Email,
			"avatar":     user.Avatar,
			"status":     user.Status,
			"created_at": user.CreatedAt,
		})
	}

	utils.Success(c, gin.H{
		"list":  userList,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

// 管理员更新用户状态
func (uc *UserController) AdminUpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	fmt.Sscanf(idStr, "%d", &id)

	var req struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var user models.User
	if err := database.GetDB().First(&user, id).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	if err := database.GetDB().Model(&user).Update("status", req.Status).Error; err != nil {
		utils.InternalServerError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

// 管理员获取用户详情
func (uc *UserController) AdminDetail(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	fmt.Sscanf(idStr, "%d", &id)

	var user models.User
	if err := database.GetDB().First(&user, id).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"phone":      user.Phone,
		"email":      user.Email,
		"avatar":     user.Avatar,
		"status":     user.Status,
		"created_at": user.CreatedAt,
	})
}
