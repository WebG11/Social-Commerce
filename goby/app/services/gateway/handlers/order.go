package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/bitdance-panic/gobuy/app/consts"
	gconsts "github.com/bitdance-panic/gobuy/app/consts"
	rpc_order "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order"
	clients "github.com/bitdance-panic/gobuy/app/services/gateway/biz/clients"
	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client/callopt"
)

func HandleCreateOrder(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt(consts.CONTEXT_UID_KEY)
	var body struct {
		ItemIDs      []int   `json:"itemIDs"`
		Phone        string  `json:"phone"`
		OrderAddress string  `json:"order_address"`
		TotalPrice   float64 `json:"total_price"`
	}
	if err := c.Bind(&body); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	itemIDs := make([]int32, len(body.ItemIDs))
	for i, itemID := range body.ItemIDs {
		itemIDs[i] = int32(itemID)
	}
	req := rpc_order.CreateOrderReq{
		UserId:       int32(userID),
		CartItemIDs:  itemIDs,
		Phone:        body.Phone,
		OrderAddress: body.OrderAddress,
		TotalPrice:   body.TotalPrice,
	}
	resp, err := clients.OrderClient.CreateOrder(context.Background(), &req, callopt.WithRPCTimeout(5*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"order": resp.Order})
}

func HandleUpdateOrderStatus(ctx context.Context, c *app.RequestContext) {
	// 1. 绑定请求参数
	var body struct {
		OrderID int `json:"order_id" binding:"required"`
		Status  int `json:"status" binding:"required"`
	}
	if err := c.Bind(&body); err != nil {
		utils.Fail(c, err.Error())
		return
	}

	// 2. 参数验证
	if body.OrderID <= 0 {
		utils.Fail(c, "无效的订单ID")
		return
	}

	// 3. 构建RPC请求
	req := rpc_order.UpdateOrderStatusReq{
		OrderId: int32(body.OrderID),
		Status:  int32(body.Status),
	}

	// 4. 调用服务
	resp, err := clients.OrderClient.UpdateOrderStatus(
		context.Background(),
		&req,
		callopt.WithRPCTimeout(5*time.Second),
	)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	// 5. 返回响应
	utils.Success(c, utils.H{
		"success": resp.Success,
		"status":  resp.Status,
		"message": "订单状态更新成功",
	})
}

func HandleUpdateOrderAddress(ctx context.Context, c *app.RequestContext) {
	// 1. 绑定请求参数
	var body struct {
		OrderID      int    `json:"order_id" binding:"required"`
		AddressName  string `json:"address_name"`
		AddressPhone string `json:"address_phone"`
		Address      string `json:"address"`
	}
	if err := c.Bind(&body); err != nil {
		utils.Fail(c, err.Error())
		return
	}

	// 2. 参数验证
	if body.OrderID <= 0 {
		utils.Fail(c, "无效的订单ID")
		return
	}

	// 3. 构建RPC请求
	req := rpc_order.UpdateOrderAddressReq{
		OrderId:      int32(body.OrderID),
		AddressName:  body.AddressName,
		AddressPhone: body.AddressPhone,
		Address:      body.Address,
	}

	// 4. 调用服务
	resp, err := clients.OrderClient.UpdateOrderAddress(
		context.Background(),
		&req,
		callopt.WithRPCTimeout(5*time.Second),
	)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	// 5. 返回响应
	utils.Success(c, utils.H{
		"success": resp.Success,
		"message": "地址更新成功",
	})
}

func HandleGetOrder(ctx context.Context, c *app.RequestContext) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	req := rpc_order.GetOrderReq{
		OrderId: int32(orderID),
	}
	resp, err := clients.OrderClient.GetOrder(context.Background(), &req, callopt.WithRPCTimeout(5*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"order": resp.Order})
}

func HandleUpdateOrderDiscount(ctx context.Context, c *app.RequestContext) {
	var body struct {
		OrderId  int    `json:"order_id"`
		Discount string `json:"discount"`
	}
	if err := c.Bind(&body); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	discount, err := strconv.ParseFloat(body.Discount, 64)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	req := rpc_order.UpdateOrderDiscountReq{
		OrderId:  int32(body.OrderId),
		Discount: discount,
	}
	resp, err := clients.OrderClient.UpdateOrderDiscount(context.Background(), &req, callopt.WithRPCTimeout(5*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"success": resp.Success})
}

func HandlePayOrder(ctx context.Context, c *app.RequestContext) {
	// 从路径参数中获取 orderId
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil || orderID <= 0 {
		utils.Fail(c, "无效的订单ID")
		return
	}

	// 构建 RPC 请求
	req := rpc_order.UpdateOrderStatusReq{
		OrderId: int32(orderID),
		Status:  1, // 设置订单状态为已支付
	}

	// 调用 OrderClient 的 UpdateOrderStatus 方法
	resp, err := clients.OrderClient.UpdateOrderStatus(
		context.Background(),
		&req,
		callopt.WithRPCTimeout(5*time.Second),
	)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	// 检查更新是否成功
	if !resp.Success {
		utils.Fail(c, "订单支付失败")
		return
	}

	// 返回成功响应
	utils.Success(c, utils.H{
		"success": true,
		"message": "订单支付成功",
	})
}

func HandleListUserOrder(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt(consts.CONTEXT_UID_KEY)
	req := rpc_order.ListOrderReq{
		UserId:   int32(userID),
		PageNum:  1,
		PageSize: 1,
	}
	resp, err := clients.OrderClient.ListUserOrder(context.Background(), &req, callopt.WithRPCTimeout(5*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"orders": resp.Orders})
}

func HandleAdminListOrder(ctx context.Context, c *app.RequestContext) {
	pageNum, err := strconv.Atoi(c.Query("page"))
	sellerID := c.GetInt(consts.CONTEXT_UID_KEY)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	pageSize, err := strconv.Atoi(c.Query("size"))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	req := rpc_order.ListSellerOrderReq{
		SellerId: int32(sellerID),
		PageNum:  int32(pageNum),
		PageSize: int32(pageSize),
	}
	resp, err := clients.OrderClient.AdminListOrder(context.Background(), &req, callopt.WithRPCTimeout(5*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"orders": resp.Orders, "total_count": resp.TotalCount})
}

func HandleUpdateOrderTracking(ctx context.Context, c *app.RequestContext) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	req := rpc_order.UpdateOrderTrackingReq{
		OrderId:        int32(orderID),
		TrackingNumber: c.Param("tracking_number"),
	}
	resp, err := clients.OrderClient.UpdateOrderTracking(context.Background(), &req, callopt.WithRPCTimeout(5*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"success": resp.Success})
}

func HandleCreateUserAddress(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt(consts.CONTEXT_UID_KEY)
	var body struct {
		UserID      int32  `json:"user_id"`
		Phone       string `json:"phone"`
		UserAddress string `json:"userAddress"`
	}
	if err := c.Bind(&body); err != nil {
		utils.Fail(c, err.Error())
		return
	}

	req := rpc_order.CreateUserAddressReq{
		UserId:      int32(userID),
		Phone:       body.Phone,
		UserAddress: body.UserAddress,
	}

	// 调用 OrderClient 的 CreateOrderAddress 方法
	resp, err := clients.OrderClient.CreateUserAddress(context.Background(), &req, callopt.WithRPCTimeout(5*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	utils.Success(c, utils.H{"userAddress": resp.UserId, "success": resp.Success})
}

func HandleDeleteUserAddress(ctx context.Context, c *app.RequestContext) {
	var body struct {
		UserID int32 `json:"userID"`
	}
	if err := c.Bind(&body); err != nil {
		utils.Fail(c, err.Error())
		return
	}

	req := rpc_order.DeleteUserAddressReq{
		UserId: body.UserID,
	}

	// 调用 OrderClient 的 DeleteOrderAddress 方法
	resp, err := clients.OrderClient.DeleteUserAddress(context.Background(), &req, callopt.WithRPCTimeout(5*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	utils.Success(c, utils.H{"userID": resp.UserId, "success": resp.Success})
}

func HandleUpdateUserAddress(ctx context.Context, c *app.RequestContext) {
	var body struct {
		UserID      int32  `json:"userID"`
		UserAddress string `json:"userAddress"`
	}
	if err := c.Bind(&body); err != nil {
		utils.Fail(c, err.Error())
		return
	}

	req := rpc_order.UpdateUserAddressReq{
		UserId:      body.UserID,
		UserAddress: body.UserAddress,
	}

	// 调用 OrderClient 的 UpdateOrderAddress 方法
	resp, err := clients.OrderClient.UpdateUserAddress(context.Background(), &req, callopt.WithRPCTimeout(5*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	utils.Success(c, utils.H{"userAddress": resp.UserAddress, "success": resp.Success})
}

func HandleGetUserAddress(ctx context.Context, c *app.RequestContext) {
	// 从字符串参数中获取 UserID
	userIDStr := c.Query("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		utils.Fail(c, "无效的用户ID")
		return
	}

	req := rpc_order.GetUserAddressReq{
		UserId: int32(userID),
	}

	// 调用 OrderClient 的 GetOrderAddress 方法
	resp, err := clients.OrderClient.GetUserAddress(context.Background(), &req, callopt.WithRPCTimeout(10*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	if resp.UserAddresses == nil {
		resp.UserAddresses = []*rpc_order.UserAddress{}
	}
	utils.Success(c, utils.H{"userAddress": resp.UserAddresses})
}

func HandleGenerateSalesReport(ctx context.Context, c *app.RequestContext) {
	var req rpc_order.SalesReportReq

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	userID := c.GetInt(gconsts.CONTEXT_UID_KEY)

	if startDate != "" {
		req.StartDate = &startDate
	}
	if endDate != "" {
		req.EndDate = &endDate
	}

	req.SellerID = int32(userID)

	resp, err := clients.OrderClient.GetSalesReport(context.Background(), &req, callopt.WithRPCTimeout(5*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	utils.Success(c, utils.H{
		"total_revenue":       resp.TotalRevenue,
		"order_count":         resp.OrderCount,
		"top_products":        resp.TopProducts,
		"average_order_amt":   resp.AverageOrderAmt,
		"total_product_count": resp.TotalProductCount,
		"daily_revenue":       resp.DailyRevenue,
	})
}

func HandleGenerateSalesReportByDate(ctx context.Context, c *app.RequestContext) {
	var req rpc_order.SalesReportByDateReq

	resp, err := clients.OrderClient.GetSalesReportByDate(context.Background(), &req, callopt.WithRPCTimeout(5*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	utils.Success(c, utils.H{
		"EveryDateRevenue": resp.GetDateRevenue(),
	})
}
