package routers

import (
	"electricVehicleMall/controllers"
	"electricVehicleMall/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 跨域中间件
	r.Use(middleware.CORS())

	// 控制器实例
	userController := &controllers.UserController{}
	adminController := &controllers.AdminController{}
	vehicleController := &controllers.VehicleController{}
	categoryController := &controllers.CategoryController{}
	noticeController := &controllers.NoticeController{}
	bannerController := &controllers.BannerController{}
	cartController := &controllers.CartController{}
	orderController := &controllers.OrderController{}
	addressController := &controllers.AddressController{}
	favoriteController := &controllers.FavoriteController{}
	commentController := &controllers.CommentController{}
	chatController := &controllers.ChatController{}

	// 公开路由 - 用户相关
	userPublic := r.Group("/api/user")
	{
		userPublic.POST("/register", userController.Register)
		userPublic.POST("/login", userController.Login)
	}

	// 公开路由 - 电动车相关（不需要登录）
	vehiclePublic := r.Group("/api/vehicle")
	{
		vehiclePublic.GET("/list", vehicleController.List)
		vehiclePublic.GET("/detail/:id", vehicleController.Detail)
	}

	// 公开路由 - 分类相关
	categoryPublic := r.Group("/api/category")
	{
		categoryPublic.GET("/list", categoryController.List)
	}

	// 公开路由 - 公告相关
	noticePublic := r.Group("/api/notice")
	{
		noticePublic.GET("/list", noticeController.List)
		noticePublic.GET("/detail/:id", noticeController.Detail)
	}

	// 公开路由 - 轮播图相关
	bannerPublic := r.Group("/api/banner")
	{
		bannerPublic.GET("/list", bannerController.List)
	}

	// 公开路由 - 评论相关
	commentPublic := r.Group("/api/comment")
	{
		commentPublic.GET("/list", commentController.List)
	}

	// 需要用户登录的路由
	userAuth := r.Group("/api/user")
	userAuth.Use(middleware.JWTAuth(), middleware.UserAuth())
	{
		// 用户个人信息
		userAuth.GET("/profile", userController.GetProfile)
		userAuth.PUT("/profile", userController.UpdateProfile)
		userAuth.PUT("/password", userController.ChangePassword)

		// 购物车
		userAuth.GET("/cart/list", cartController.List)
		userAuth.POST("/cart/add", cartController.Add)
		userAuth.PUT("/cart/update/:id", cartController.Update)
		userAuth.DELETE("/cart/delete/:id", cartController.Delete)
		userAuth.POST("/cart/clear", cartController.Clear)

		// 地址
		userAuth.GET("/address/list", addressController.List)
		userAuth.GET("/address/detail/:id", addressController.Detail)
		userAuth.POST("/address/create", addressController.Create)
		userAuth.PUT("/address/update/:id", addressController.Update)
		userAuth.DELETE("/address/delete/:id", addressController.Delete)
		userAuth.POST("/address/default/:id", addressController.SetDefault)

		// 订单
		userAuth.GET("/order/list", orderController.List)
		userAuth.GET("/order/detail/:id", orderController.Detail)
		userAuth.POST("/order/create", orderController.Create)
		userAuth.POST("/order/cancel/:id", orderController.Cancel)
		userAuth.POST("/order/pay/:id", orderController.Pay)
		userAuth.POST("/order/receive/:id", orderController.Receive)

		// 收藏
		userAuth.GET("/favorite/list", favoriteController.List)
		userAuth.POST("/favorite/toggle/:vehicle_id", favoriteController.Toggle)
		userAuth.GET("/favorite/check/:vehicle_id", favoriteController.Check)
		userAuth.DELETE("/favorite/delete/:id", favoriteController.Delete)

		// 评论
		userAuth.POST("/comment/create", commentController.Create)

		// 客服
		userAuth.GET("/chat/list", chatController.UserList)
		userAuth.POST("/chat/send", chatController.UserSend)
	}

	// 公开路由 - 管理员登录
	adminPublic := r.Group("/api/admin")
	{
		adminPublic.POST("/login", adminController.Login)
	}

	// 需要管理员登录的路由
	adminAuth := r.Group("/api/admin")
	adminAuth.Use(middleware.JWTAuth(), middleware.AdminAuth())
	{
		// 管理员个人信息
		adminAuth.GET("/profile", adminController.GetProfile)

		// 管理员管理
		adminAuth.GET("/list", adminController.List)
		adminAuth.POST("/create", adminController.Create)
		adminAuth.PUT("/update/:id", adminController.Update)
		adminAuth.DELETE("/delete/:id", adminController.Delete)
		adminAuth.POST("/reset-password/:id", adminController.ResetPassword)

		// 电动车管理
		adminAuth.POST("/vehicle/create", vehicleController.Create)
		adminAuth.PUT("/vehicle/update/:id", vehicleController.Update)
		adminAuth.DELETE("/vehicle/delete/:id", vehicleController.Delete)

		// 分类管理
		adminAuth.POST("/category/create", categoryController.Create)
		adminAuth.PUT("/category/update/:id", categoryController.Update)
		adminAuth.DELETE("/category/delete/:id", categoryController.Delete)

		// 公告管理
		adminAuth.GET("/notice/list", noticeController.List)
		adminAuth.GET("/notice/detail/:id", noticeController.Detail)
		adminAuth.POST("/notice/create", noticeController.Create)
		adminAuth.PUT("/notice/update/:id", noticeController.Update)
		adminAuth.DELETE("/notice/delete/:id", noticeController.Delete)

		// 轮播图管理
		adminAuth.GET("/banner/list", bannerController.List)
		adminAuth.GET("/banner/detail/:id", bannerController.Detail)
		adminAuth.POST("/banner/create", bannerController.Create)
		adminAuth.PUT("/banner/update/:id", bannerController.Update)
		adminAuth.DELETE("/banner/delete/:id", bannerController.Delete)

		// 订单管理
		adminAuth.GET("/order/list", orderController.AdminList)
		adminAuth.GET("/order/detail/:id", orderController.AdminDetail)
		adminAuth.POST("/order/ship/:id", orderController.Ship)

		// 评论管理
		adminAuth.GET("/comment/list", commentController.AdminList)
		adminAuth.PUT("/comment/status/:id", commentController.UpdateStatus)
		adminAuth.DELETE("/comment/delete/:id", commentController.Delete)

		// 订单统计
		adminAuth.GET("/statistics", orderController.Statistics)

		// 用户管理
		adminAuth.GET("/user/list", userController.AdminList)
		adminAuth.GET("/user/detail/:id", userController.AdminDetail)
		adminAuth.PUT("/user/status/:id", userController.AdminUpdateStatus)

		// 客服管理
		adminAuth.GET("/chat/users", chatController.AdminListUsers)
		adminAuth.GET("/chat/messages/:user_id", chatController.AdminListMessages)
		adminAuth.POST("/chat/reply", chatController.AdminSend)
	}

	return r
}
