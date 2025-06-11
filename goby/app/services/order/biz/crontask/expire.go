package crontask

import (
	"log"
	"time"

	"github.com/bitdance-panic/gobuy/app/consts"

	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/services/cart/biz/dal/tidb"
	"github.com/bitdance-panic/gobuy/app/services/order/biz/dao"
)

func CheckAndUpdateExpiredOrders() {
	// ctx := context.Background()
	// 获取所有状态为“待支付”的订单
	orders, err := dao.ListPendingOrder(tidb.DB)
	if err != nil {
		log.Printf("Failed to get pending orders:%v", err)
		return
	}
	// 遍历订单，检查是否过期
	for _, order := range *orders {
		if isOrderExpired(order) {
			// 更新订单状态为"已取消"
			err := dao.UpdateOrderStatus(tidb.DB, &order, consts.OrderStatusCancelled)
			if err != nil {
				log.Printf("Failed to update order:%v", err)
				continue
			}
			log.Printf("Order %d has been cancelled due to expiration", order.ID)
		}
	}
}

// 检查订单是否过期
func isOrderExpired(order models.Order) bool {
	// 假设订单创建后30分钟未支付则过期
	expirationTime := order.CreatedAt.Add(30 * time.Minute)
	return time.Now().After(expirationTime)
}
