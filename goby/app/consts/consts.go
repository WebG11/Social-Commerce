package consts

// 定义订单状态常量
type OrderStatus int

const (
	OrderStatusPending   OrderStatus = 0 // 待支付
	OrderStatusPaid      OrderStatus = 1 // 已支付
	OrderStatusShipped   OrderStatus = 2 // 已发货
	OrderStatusDelivered OrderStatus = 3 // 已收货
	OrderStatusCancelled OrderStatus = 4 // 已取消
	OrderStatusRefunded  OrderStatus = 5 // 已退款
)

const (
	CONTEXT_UID_KEY       string = "uid"
	ITEM_INITIAL_QUANTITY int    = 1
)
