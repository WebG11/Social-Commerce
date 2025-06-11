package handlers

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	rpc_order "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/clients"
	"github.com/bitdance-panic/gobuy/app/services/gateway/conf"
	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/smartwalle/alipay/v3"
	"github.com/smartwalle/xid"
)

var (
	notifyURL   = conf.GetConf().Alipay.ServerDomain + "/alipay/notify"
	callbackURL = conf.GetConf().Alipay.ServerDomain + "/alipay/callback"
)

func HandleGetPayUrl(ctx context.Context, c *app.RequestContext) {
	// 从请求中获取系统订单ID
	systemOrderID := c.Query("order_id")
	orderID, err := strconv.Atoi(systemOrderID)
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

	// 生成支付宝交易号
	tradeNo := fmt.Sprintf("alipay_%d", xid.Next())
	p := alipay.TradePagePay{}
	p.NotifyURL = notifyURL
	p.ReturnURL = callbackURL
	p.Subject = "订单支付:" + tradeNo
	p.OutTradeNo = tradeNo
	p.TotalAmount = strconv.FormatFloat(resp.Order.TotalPrice-resp.Order.Discount, 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	// 使用 PassbackParams 传递系统订单ID
	// 需要进行URL编码，因为这个参数会通过URL传递
	p.PassbackParams = url.QueryEscape(systemOrderID)
	url, err := clients.AlipayClient.TradePagePay(p)
	if err != nil {
		utils.Fail(c, "")
		return
	}
	payURL := url.String()
	log.Printf("创建支付链接成功，系统订单ID: %s, 支付宝交易号: %s", systemOrderID, tradeNo)
	// 直接重定向
	utils.Success(c, utils.H{"url": payURL})
}
