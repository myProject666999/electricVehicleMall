package controllers

import (
	"electricVehicleMall/database"
	"electricVehicleMall/models"
	"electricVehicleMall/utils"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type OrderController struct{}

type CreateOrderRequest struct {
	AddressID uint   `json:"address_id" binding:"required"`
	VehicleIDs []uint `json:"vehicle_ids" binding:"required"`
	Remark     string `json:"remark"`
}

type ShipOrderRequest struct {
	LogisticsNo  string `json:"logistics_no" binding:"required"`
	LogisticsCom string `json:"logistics_com" binding:"required"`
}

type OrderStatus int

const (
	OrderStatusPendingPayment OrderStatus = iota // 0: 待支付
	OrderStatusPendingShip                         // 1: 待发货
	OrderStatusPendingReceive                      // 2: 待收货
	OrderStatusCompleted                           // 3: 已完成
	OrderStatusCancelled                           // 4: 已取消
	OrderStatusRefunded                            // 5: 已退款
)

func (oc *OrderController) GenerateOrderNo() string {
	now := time.Now()
	return fmt.Sprintf("EV%s%d", now.Format("20060102150405"), now.Nanosecond()/1000)
}

func (oc *OrderController) List(c *gin.Context) {
	userID := c.GetUint("user_id")

	var orders []models.Order
	var total int64

	query := database.GetDB().Model(&models.Order{}).Where("user_id = ?", userID)

	// 状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
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
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders)

	// 获取订单详情
	var orderList []gin.H
	for _, order := range orders {
		var orderItems []models.OrderItem
		database.GetDB().Where("order_id = ?", order.ID).Find(&orderItems)

		orderList = append(orderList, gin.H{
			"order":       order,
			"order_items": orderItems,
		})
	}

	utils.Success(c, gin.H{
		"list":      orderList,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (oc *OrderController) AdminList(c *gin.Context) {
	var orders []models.Order
	var total int64

	query := database.GetDB().Model(&models.Order{})

	// 搜索
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("order_no LIKE ? OR consignee LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
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
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders)

	// 获取订单详情
	var orderList []gin.H
	for _, order := range orders {
		var orderItems []models.OrderItem
		database.GetDB().Where("order_id = ?", order.ID).Find(&orderItems)

		var user models.User
		database.GetDB().Select("id, username, phone").First(&user, order.UserID)

		orderList = append(orderList, gin.H{
			"order":       order,
			"order_items": orderItems,
			"user":        user,
		})
	}

	utils.Success(c, gin.H{
		"list":      orderList,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (oc *OrderController) Detail(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	var order models.Order
	if err := database.GetDB().Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	var orderItems []models.OrderItem
	database.GetDB().Where("order_id = ?", order.ID).Find(&orderItems)

	utils.Success(c, gin.H{
		"order":       order,
		"order_items": orderItems,
	})
}

func (oc *OrderController) AdminDetail(c *gin.Context) {
	id := c.Param("id")

	var order models.Order
	if err := database.GetDB().First(&order, id).Error; err != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	var orderItems []models.OrderItem
	database.GetDB().Where("order_id = ?", order.ID).Find(&orderItems)

	var user models.User
	database.GetDB().Select("id, username, phone, email").First(&user, order.UserID)

	utils.Success(c, gin.H{
		"order":       order,
		"order_items": orderItems,
		"user":        user,
	})
}

func (oc *OrderController) Create(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 检查地址
	var address models.Address
	if err := database.GetDB().Where("id = ? AND user_id = ?", req.AddressID, userID).First(&address).Error; err != nil {
		utils.NotFound(c, "地址不存在")
		return
	}

	// 检查购物车中的商品
	var cartItems []models.Cart
	database.GetDB().Where("user_id = ? AND vehicle_id IN ?", userID, req.VehicleIDs).Find(&cartItems)

	if len(cartItems) == 0 {
		utils.BadRequest(c, "购物车中没有选中的商品")
		return
	}

	// 开始事务
	tx := database.GetDB().Begin()

	// 生成订单号
	orderNo := oc.GenerateOrderNo()

	// 计算订单金额
	var totalAmount float64
	var orderItems []models.OrderItem

	for _, cartItem := range cartItems {
		var vehicle models.ElectricVehicle
		if err := tx.First(&vehicle, cartItem.VehicleID).Error; err != nil {
			tx.Rollback()
			utils.NotFound(c, "电动车不存在")
			return
		}

		if vehicle.Status != 1 {
			tx.Rollback()
			utils.BadRequest(c, "电动车 "+vehicle.Name+" 已下架")
			return
		}

		if vehicle.Stock < cartItem.Quantity {
			tx.Rollback()
			utils.BadRequest(c, "电动车 "+vehicle.Name+" 库存不足")
			return
		}

		itemTotal := vehicle.Price * float64(cartItem.Quantity)
		totalAmount += itemTotal

		orderItem := models.OrderItem{
			OrderNo:     orderNo,
			VehicleID:   vehicle.ID,
			VehicleName: vehicle.Name,
			Image:       vehicle.Image,
			Price:       vehicle.Price,
			Quantity:    cartItem.Quantity,
			TotalPrice:  itemTotal,
		}
		orderItems = append(orderItems, orderItem)

		// 扣减库存
		tx.Model(&vehicle).Update("stock", vehicle.Stock-cartItem.Quantity)
		// 增加销量
		tx.Model(&vehicle).Update("sales", vehicle.Sales+cartItem.Quantity)
	}

	// 创建订单
	order := models.Order{
		OrderNo:     orderNo,
		UserID:      userID,
		TotalAmount: totalAmount,
		Status:      int(OrderStatusPendingPayment),
		Consignee:   address.Consignee,
		Phone:       address.Phone,
		Address:     address.Province + " " + address.City + " " + address.District + " " + address.Address,
		Remark:      req.Remark,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "创建订单失败: "+err.Error())
		return
	}

	// 创建订单项
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
	}
	if err := tx.Create(&orderItems).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "创建订单项失败: "+err.Error())
		return
	}

	// 删除购物车中的商品
	tx.Where("user_id = ? AND vehicle_id IN ?", userID, req.VehicleIDs).Delete(&models.Cart{})

	// 提交事务
	tx.Commit()

	utils.SuccessWithMessage(c, "订单创建成功", gin.H{
		"order_id":   order.ID,
		"order_no":   order.OrderNo,
		"total_amount": order.TotalAmount,
	})
}

func (oc *OrderController) Cancel(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	var order models.Order
	if err := database.GetDB().Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	if order.Status != int(OrderStatusPendingPayment) {
		utils.BadRequest(c, "只有待支付的订单可以取消")
		return
	}

	// 开始事务
	tx := database.GetDB().Begin()

	// 恢复库存
	var orderItems []models.OrderItem
	tx.Where("order_id = ?", order.ID).Find(&orderItems)

	for _, item := range orderItems {
		var vehicle models.ElectricVehicle
		tx.First(&vehicle, item.VehicleID)
		tx.Model(&vehicle).Update("stock", vehicle.Stock+item.Quantity)
		tx.Model(&vehicle).Update("sales", vehicle.Sales-item.Quantity)
	}

	// 更新订单状态
	tx.Model(&order).Update("status", int(OrderStatusCancelled))

	// 提交事务
	tx.Commit()

	utils.SuccessWithMessage(c, "订单已取消", nil)
}

func (oc *OrderController) Pay(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	var order models.Order
	if err := database.GetDB().Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	if order.Status != int(OrderStatusPendingPayment) {
		utils.BadRequest(c, "订单状态不正确")
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	order.Status = int(OrderStatusPendingShip)
	order.PayAmount = order.TotalAmount
	order.PayTime = &now

	database.GetDB().Save(&order)

	utils.SuccessWithMessage(c, "支付成功", nil)
}

func (oc *OrderController) Ship(c *gin.Context) {
	id := c.Param("id")

	var req ShipOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var order models.Order
	if err := database.GetDB().First(&order, id).Error; err != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	if order.Status != int(OrderStatusPendingShip) {
		utils.BadRequest(c, "只有待发货的订单可以发货")
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	order.Status = int(OrderStatusPendingReceive)
	order.LogisticsNo = req.LogisticsNo
	order.LogisticsCom = req.LogisticsCom
	order.ShipTime = &now

	database.GetDB().Save(&order)

	utils.SuccessWithMessage(c, "发货成功", nil)
}

func (oc *OrderController) Receive(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	var order models.Order
	if err := database.GetDB().Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	if order.Status != int(OrderStatusPendingReceive) {
		utils.BadRequest(c, "只有待收货的订单可以确认收货")
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	order.Status = int(OrderStatusCompleted)
	order.ReceiveTime = &now

	database.GetDB().Save(&order)

	utils.SuccessWithMessage(c, "确认收货成功", nil)
}

func (oc *OrderController) Statistics(c *gin.Context) {
	// 订单统计
	var totalOrders int64
	var pendingPayment int64
	var pendingShip int64
	var pendingReceive int64
	var completed int64
	var todaySales float64

	// 总订单数
	database.GetDB().Model(&models.Order{}).Count(&totalOrders)

	// 待支付
	database.GetDB().Model(&models.Order{}).Where("status = ?", OrderStatusPendingPayment).Count(&pendingPayment)

	// 待发货
	database.GetDB().Model(&models.Order{}).Where("status = ?", OrderStatusPendingShip).Count(&pendingShip)

	// 待收货
	database.GetDB().Model(&models.Order{}).Where("status = ?", OrderStatusPendingReceive).Count(&pendingReceive)

	// 已完成
	database.GetDB().Model(&models.Order{}).Where("status = ?", OrderStatusCompleted).Count(&completed)

	// 今日销售额
	today := time.Now().Format("2006-01-02")
	var todayOrders []models.Order
	database.GetDB().Model(&models.Order{}).
		Where("DATE(created_at) = ? AND status >= ?", today, OrderStatusPendingShip).
		Find(&todayOrders)

	for _, order := range todayOrders {
		if order.PayAmount > 0 {
			todaySales += order.PayAmount
		} else {
			todaySales += order.TotalAmount
		}
	}

	utils.Success(c, gin.H{
		"total_orders":     totalOrders,
		"pending_payment":  pendingPayment,
		"pending_ship":     pendingShip,
		"pending_receive":  pendingReceive,
		"completed":        completed,
		"today_sales":      todaySales,
	})
}
