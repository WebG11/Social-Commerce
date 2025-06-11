package bll

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bitdance-panic/gobuy/app/consts"

	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/services/order/biz/clients"
	"github.com/bitdance-panic/gobuy/app/services/order/biz/dal/tidb"
	"github.com/bitdance-panic/gobuy/app/services/order/biz/dao"
	"github.com/bitdance-panic/gobuy/app/utils"

	rpc_cart "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/cart"
	rpc_order "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order"
)

type OrderBLL struct{}

func NewOrderBLL() *OrderBLL {
	return &OrderBLL{}
}

// TODO 自动机
// 创建一个新订单
func (bll *OrderBLL) CreateOrder(ctx context.Context, req *rpc_order.CreateOrderReq) (*rpc_order.CreateOrderResp, error) {
	orderNumber, err := utils.GenerateID()
	if err != nil {
		return nil, err
	}
	orderItems := make([]models.OrderItem, len(req.CartItemIDs))

	// 用于检查卖家ID
	var sellerID int64 = -1

	for i, cartItemID := range req.CartItemIDs {
		resp, err := clients.CartClient.GetItem(
			ctx,
			&rpc_cart.GetItemReq{ItemId: cartItemID},
		)
		if err != nil {
			return nil, err
		}

		cartItem := resp.Item

		// 获取商品信息以检查卖家ID
		productResp, err := dao.GetByID(tidb.DB, int(cartItem.ProductId))
		if err != nil {
			return nil, err
		}

		// 检查卖家ID
		if sellerID == -1 {
			sellerID = int64(productResp.SellerID)
		} else if sellerID != int64(productResp.SellerID) {
			return nil, errors.New("订单中的商品必须来自同一个卖家")
		}

		orderItem := models.OrderItem{
			ProductID:    int(cartItem.ProductId),
			ProductName:  cartItem.Name,
			ProductImage: cartItem.Image,
			ItemStatus:   false,
			Price:        cartItem.Price,
			Quantity:     int(cartItem.Quantity),
		}
		// 避免外键报错
		orderItem.Product.ID = int(orderItem.ProductID)
		// itemWithID, err := dao.CreateOrderItem(tidb.DB, &orderItem)
		// if err != nil {
		// 	return nil, err
		// }
		orderItems[i] = orderItem
		// p, _ := decimal.NewFromFloat(orderItem.Price).Mul(decimal.NewFromInt(int64(orderItem.Quantity))).Float64()
		// if !exact {
		// 	return nil, errors.New("handle decimal is not exact")
		// }
	}

	for _, item := range req.CartItemIDs {
		_, err := clients.CartClient.DeleteItem(
			ctx,
			&rpc_cart.DeleteItemReq{ItemId: int32(item)},
		)
		if err != nil {
			return nil, err
		}
	}

	// 先保存下订单
	order := models.Order{
		UserID:     int(req.UserId),
		Number:     orderNumber,
		TotalPrice: req.TotalPrice,
		Status:     int(consts.OrderStatusPending),
		PayTime:    nil,
		Phone:      req.Phone,
		Address:    req.OrderAddress,
		SellerID:   int(sellerID), // 设置卖家ID
	}

	err = dao.CreateOrder(tidb.DB, &order)
	if err != nil {
		return nil, err
	}
	// 再存订单项
	for i := range len(orderItems) {
		orderItems[i].OrderID = order.ID
	}
	order.Items = orderItems
	// TODO 事务
	err = dao.SaveOrder(tidb.DB, &order)
	if err != nil {
		return nil, err
	}
	protoOrder := convertOrderToProto(&order)
	return &rpc_order.CreateOrderResp{
		Order: protoOrder,
	}, err
}

func (bll *OrderBLL) GetOrder(ctx context.Context, req *rpc_order.GetOrderReq) (*rpc_order.GetOrderResp, error) {
	order, err := dao.GetOrderByID(tidb.DB, int(req.OrderId))
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	protoOrder := convertOrderToProto(order)
	return &rpc_order.GetOrderResp{
		Order: protoOrder,
	}, nil
}

func (bll *OrderBLL) UpdateOrderDiscount(ctx context.Context, req *rpc_order.UpdateOrderDiscountReq) (*rpc_order.UpdateOrderDiscountResp, error) {
	order, err := dao.GetOrderByID(tidb.DB, int(req.OrderId))
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	err = dao.UpdateOrderDiscount(tidb.DB, order, req.Discount)
	if err != nil {
		return nil, err
	}
	return &rpc_order.UpdateOrderDiscountResp{
		Success: true,
	}, nil
}

func (bll *OrderBLL) UpdateOrderStatus(ctx context.Context, req *rpc_order.UpdateOrderStatusReq) (*rpc_order.UpdateOrderStatusResp, error) {
	order, err := dao.GetOrderByID(tidb.DB, int(req.OrderId))
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	err = dao.UpdateOrderStatus(tidb.DB, order, consts.OrderStatus(req.Status))
	if err != nil {
		return nil, err
	}
	return &rpc_order.UpdateOrderStatusResp{
		Success: true,
	}, nil
}

func (bll *OrderBLL) UpdateOrderAddress(ctx context.Context, req *rpc_order.UpdateOrderAddressReq) (*rpc_order.UpdateOrderAddressResp, error) {
	order, err := dao.GetOrderByID(tidb.DB, int(req.OrderId))
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	err = dao.UpdateOrderAddress(tidb.DB, order, req.AddressName, req.AddressPhone, req.Address)
	if err != nil {
		return nil, err
	}
	return &rpc_order.UpdateOrderAddressResp{
		OrderAddress: req.Address,
		Success:      true,
	}, nil
}

func (bll *OrderBLL) ListUserOrder(ctx context.Context, req *rpc_order.ListOrderReq) (*rpc_order.ListOrderResp, error) {
	orders, err := dao.ListUserOrder(tidb.DB, int(req.UserId), int(req.PageNum), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	protoOrders := make([]*rpc_order.Order, len(*orders))
	for i, order := range *orders {
		protoOrder := convertOrderToProto(&order)
		protoOrders[i] = protoOrder
	}
	return &rpc_order.ListOrderResp{
		Orders: protoOrders,
	}, nil
}

func (bll *OrderBLL) UpdateOrderTracking(ctx context.Context, req *rpc_order.UpdateOrderTrackingReq) (*rpc_order.UpdateOrderTrackingResp, error) {
	order, err := dao.GetOrderByID(tidb.DB, int(req.OrderId))
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	err = dao.UpdateOrderTracking(tidb.DB, order, req.TrackingNumber)
	if err != nil {
		return nil, err
	}
	return &rpc_order.UpdateOrderTrackingResp{
		Success: true,
	}, nil
}

func (bll *OrderBLL) AdminListOrder(ctx context.Context, req *rpc_order.ListSellerOrderReq) (*rpc_order.ListOrderResp, error) {
	orders, total, err := dao.AdminListOrder(tidb.DB, int(req.PageNum), int(req.PageSize), int(req.SellerId))
	if err != nil {
		return nil, err
	}
	protoOrders := make([]*rpc_order.Order, 0, len(*orders))
	for _, order := range *orders {
		protoOrder := convertOrderToProto(&order)
		protoOrders = append(protoOrders, protoOrder)
	}
	return &rpc_order.ListOrderResp{
		Orders:     protoOrders,
		TotalCount: total,
	}, nil
}

func convertItemFromProto(items rpc_order.OrderItem) *models.OrderItem {
	return &models.OrderItem{
		OrderID:      int(items.OrderId),
		ProductID:    int(items.ProductId),
		Quantity:     int(items.Quantity),
		Price:        items.Price,
		ProductName:  items.ProductName,
		ProductImage: items.ProductImage,
	}
}

func convertItemToProto(items models.OrderItem) *rpc_order.OrderItem {
	return &rpc_order.OrderItem{
		OrderId:      int32(items.OrderID),
		ProductId:    int32(items.ProductID),
		ItemStatus:   items.ItemStatus,
		Quantity:     int32(items.Quantity),
		Price:        items.Price,
		ProductName:  items.ProductName,
		ProductImage: items.ProductImage,
	}
}

func convertOrderToProto(order *models.Order) *rpc_order.Order {
	protoItems := make([]*rpc_order.OrderItem, len(order.Items))
	for i, item := range order.Items {
		orderItem := convertItemToProto(item)
		protoItems[i] = orderItem
	}
	return &rpc_order.Order{
		Id:             int32(order.ID),
		UserId:         int32(order.UserID),
		Number:         order.Number,
		TotalPrice:     order.TotalPrice,
		Discount:       order.Discount,
		Status:         int32(order.Status),
		Items:          protoItems,
		CreatedAt:      utils.FormatTime(order.CreatedAt),
		Name:           order.Name,
		Phone:          order.Phone,
		Address:        order.Address,
		TrackingNumber: order.TrackingNumber,
		SellerId:       int32(order.SellerID),
	}
}

// 创建订单地址
func (bll *OrderBLL) CreateUserAddress(ctx context.Context, req *rpc_order.CreateUserAddressReq) (*rpc_order.CreateUserAddressResp, error) {
	userAddress := models.UserAddress{
		UserID:      int(req.UserId),
		Phone:       req.Phone,
		UserAddress: req.UserAddress,
	}

	// 调用 DAO 层创建订单地址
	err := dao.CreateUserAddress(tidb.DB, &userAddress)
	if err != nil {
		return nil, err
	}

	return &rpc_order.CreateUserAddressResp{
		UserId:  int32(userAddress.UserID),
		Success: true,
	}, nil
}

// 删除订单地址
func (bll *OrderBLL) DeleteUserAddress(ctx context.Context, req *rpc_order.DeleteUserAddressReq) (*rpc_order.DeleteUserAddressResp, error) {
	// 调用 DAO 层删除订单地址
	err := dao.DeleteUserAddress(tidb.DB, req.UserId)
	if err != nil {
		return nil, err
	}

	return &rpc_order.DeleteUserAddressResp{
		UserId:  req.UserId,
		Success: true,
	}, nil
}

// 更新订单地址
func (bll *OrderBLL) UpdateUserAddress(ctx context.Context, req *rpc_order.UpdateUserAddressReq) (*rpc_order.UpdateUserAddressResp, error) {
	// 调用 DAO 层更新订单地址
	err := dao.UpdateUserAddress(tidb.DB, req.UserId, req.UserAddress)
	if err != nil {
		return nil, err
	}

	return &rpc_order.UpdateUserAddressResp{
		UserAddress: req.UserAddress,
		Success:     true,
	}, nil
}

// 获取订单地址
func (bll *OrderBLL) GetUserAddress(ctx context.Context, req *rpc_order.GetUserAddressReq) (*rpc_order.GetUserAddressResp, error) {
	// 调用 DAO 层获取订单地址
	userAddresses, err := dao.GetUserAddress(tidb.DB, req.UserId)
	if err != nil {
		return nil, err
	}

	// 将地址列表转换为 RPC 格式
	var rpcAddresses []*rpc_order.UserAddress
	for _, address := range userAddresses {
		rpcAddresses = append(rpcAddresses, &rpc_order.UserAddress{
			UserId:  int32(address.UserID),
			Name:    address.UserName,
			Phone:   address.Phone,
			Address: address.UserAddress,
		})
	}

	return &rpc_order.GetUserAddressResp{
		UserAddresses: rpcAddresses,
	}, nil
}

func (bll *OrderBLL) GenerateSalesReport(ctx context.Context, req *rpc_order.SalesReportReq) (*rpc_order.SalesReportResp, error) {
	var startTime, endTime time.Time
	var err error

	if req.IsSetStartDate() {
		startTime, err = time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return nil, errors.New("invalid start date format")
		}
	} else {
		startTime, err = time.Parse("2006-01-02", "2000-01-01")
	}

	if req.IsSetEndDate() {
		endTime, err = time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, errors.New("invalid end date format")
		}
	} else {
		endTime = time.Now()
	}

	// 调用 DAO 层获取订单数据
	orders, err := dao.ListOrderByDateRange(tidb.DB, startTime, endTime)
	if err != nil {
		return nil, err
	}
	filteredOrders := make([]models.Order, 0)
	for _, order := range *orders {
		if (order.Status == 1 || order.Status == 2 || order.Status == 3) && order.SellerID == int(req.SellerID) {
			filteredOrders = append(filteredOrders, order)
		}
	}
	orders = &filteredOrders

	// 分析订单数据
	totalRevenue := 0.0
	totalProductCount := 0
	orderCount := len(*orders)
	productCount := make(map[string]int)
	productSales := make(map[string]int) // 产品销量统计
	dailyRevenue := make(map[string]float64)

	for _, order := range *orders {
		// 累加总销售额
		totalRevenue += order.TotalPrice
		orderDate := order.CreatedAt.Format("2006-01-02")
		dailyRevenue[orderDate] += order.TotalPrice

		// 统计每个产品的销量
		for _, item := range order.Items {
			productName := item.ProductName
			productSales[productName] += item.Quantity
			productCount[productName] += 1
			totalProductCount += item.Quantity
		}
	}

	// 计算平均订单金额
	averageOrderAmt := 0.0
	if orderCount > 0 {
		averageOrderAmt = totalRevenue / float64(orderCount)
	}

	// 找出最畅销的三件产品及其销量

	topProducts := make(map[string]int32)
	for productName, quantity := range productSales {
		topProducts[productName] = int32(quantity)
	}

	// 构造结果
	resp := &rpc_order.SalesReportResp{
		TotalRevenue:      totalRevenue,
		OrderCount:        int32(orderCount),
		TopProducts:       topProducts,
		AverageOrderAmt:   averageOrderAmt,
		TotalProductCount: int32(totalProductCount),
		DailyRevenue:      dailyRevenue,
	}

	return resp, nil
}

func (bll *OrderBLL) GenerateSalesReportByDate(ctx context.Context, req *rpc_order.SalesReportByDateReq) (*rpc_order.SalesReportByDateResp, error) {
	var startTime, endTime time.Time
	var err error

	endTime = time.Now().AddDate(0, 0, 1)
	startTime = endTime.AddDate(0, 0, -14)

	orders, err := dao.ListOrderByDateRange(tidb.DB, startTime, endTime)
	if err != nil {
		return nil, err
	}

	dailyRevenue := make(map[string]float64)
	currentTime := startTime
	for currentTime.Before(time.Now()) {
		dailyRevenue[currentTime.Format("2006-01-02")] = 0
		currentTime = currentTime.AddDate(0, 0, 1)
	}

	for _, order := range *orders {
		// 获取订单的日期部分
		orderDate := order.CreatedAt.Format("2006-01-02")
		fmt.Println(orderDate)
		// 累加当天的营业额
		dailyRevenue[orderDate] += order.TotalPrice
	}

	// 构造结果
	resp := &rpc_order.SalesReportByDateResp{
		DateRevenue: dailyRevenue,
	}

	return resp, nil
}
