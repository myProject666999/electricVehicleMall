package database

import (
	"fmt"
	"log"

	"electricVehicleMall/models"

	"golang.org/x/crypto/bcrypt"
)

func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.Category{},
		&models.ElectricVehicle{},
		&models.Banner{},
		&models.Notice{},
		&models.Cart{},
		&models.Order{},
		&models.OrderItem{},
		&models.Address{},
		&models.Favorite{},
		&models.Comment{},
		&models.ChatMessage{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Database migration completed successfully")

	// 创建默认管理员
	seedData()
}

func seedData() {
	// 创建默认分类
	var categoryCount int64
	DB.Model(&models.Category{}).Count(&categoryCount)
	if categoryCount == 0 {
		categories := []models.Category{
			{Name: "电动自行车", Sort: 1, Status: 1},
			{Name: "电动摩托车", Sort: 2, Status: 1},
			{Name: "电动三轮车", Sort: 3, Status: 1},
			{Name: "电动四轮车", Sort: 4, Status: 1},
		}
		DB.Create(&categories)
		fmt.Println("Default categories created")
	}

	// 创建默认管理员
	var adminCount int64
	DB.Model(&models.Admin{}).Count(&adminCount)
	if adminCount == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := models.Admin{
			Username: "admin",
			Password: string(hashedPassword),
			RealName: "超级管理员",
			Role:     1,
			Status:   1,
		}
		DB.Create(&admin)
		fmt.Println("Default admin created: username=admin, password=admin123")
	}

	// 创建默认轮播图
	var bannerCount int64
	DB.Model(&models.Banner{}).Count(&bannerCount)
	if bannerCount == 0 {
		banners := []models.Banner{
			{Title: "新品上市", Image: "", Sort: 1, Status: 1},
			{Title: "限时优惠", Image: "", Sort: 2, Status: 1},
			{Title: "热销推荐", Image: "", Sort: 3, Status: 1},
		}
		DB.Create(&banners)
		fmt.Println("Default banners created")
	}

	// 创建默认公告
	var noticeCount int64
	DB.Model(&models.Notice{}).Count(&noticeCount)
	if noticeCount == 0 {
		notices := []models.Notice{
			{Title: "欢迎使用电动车销售系统", Content: "欢迎使用电动车销售管理系统，祝您购物愉快！", Author: "管理员", Sort: 1, Status: 1},
			{Title: "新用户注册优惠", Content: "新用户注册即送100元优惠券，满5000可用！", Author: "管理员", Sort: 2, Status: 1},
		}
		DB.Create(&notices)
		fmt.Println("Default notices created")
	}
}
