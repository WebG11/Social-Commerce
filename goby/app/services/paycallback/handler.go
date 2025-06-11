package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/bitdance-panic/gobuy/app/consts"
	rpc_order "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order"
	"github.com/bitdance-panic/gobuy/app/services/paycallback/biz/clients"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
)

// 与用户响应同步
func handleCallback(c *gin.Context) {
	// 解析请求参数
	if err := c.Request.ParseForm(); err != nil {
		log.Println("解析请求参数发生错误", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "解析请求参数发生错误"})
		return
	}

	// 获取通知参数
	outTradeNo := c.Request.Form.Get("out_trade_no")

	// 获取系统订单ID (PassbackParams)
	passbackParams := c.Request.Form.Get("passback_params")
	systemOrderID, err := url.QueryUnescape(passbackParams)
	if err != nil {
		log.Println("解析系统订单ID发生错误", err)
		systemOrderID = ""
	}

	// 验证签名
	if err := clients.AlipayClient.VerifySign(c.Request.Form); err != nil {
		log.Println("回调验证签名发生错误", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "回调验证签名发生错误"})
		return
	}
	// 查询订单状态
	var p = alipay.TradeQuery{}
	p.OutTradeNo = outTradeNo

	rsp, err := clients.AlipayClient.TradeQuery(c, p)
	if err != nil {
		log.Printf("验证订单 %s 信息发生错误: %s", outTradeNo, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("验证订单 %s 信息发生错误: %s", outTradeNo, err.Error())})
		return
	}

	if rsp.IsFailure() {
		log.Printf("验证订单 %s 信息发生错误: %s-%s", outTradeNo, rsp.Msg, rsp.SubMsg)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("验证订单 %s 信息发生错误: %s-%s", outTradeNo, rsp.Msg, rsp.SubMsg)})
		return
	}

	if systemOrderID != "" {
		log.Printf("系统订单 %s 支付成功，支付宝交易号: %s", systemOrderID, outTradeNo)
	}
	redirectUrl := "http://localhost:3000/"
	fmt.Println("重定向到", redirectUrl)
	c.Redirect(http.StatusFound, redirectUrl)
}

// 只在服务端
func handleNotify(c *gin.Context) {
	// 解析请求参数
	if err := c.Request.ParseForm(); err != nil {
		log.Println("解析请求参数发生错误", err)
		c.String(http.StatusBadRequest, "fail")
		return
	}
	// 解析异步通知
	notification, err := clients.AlipayClient.DecodeNotification(c.Request.Form)
	if err != nil {
		log.Println("解析异步通知发生错误", err)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	// 获取系统订单ID
	systemOrderID, err := url.QueryUnescape(notification.PassbackParams)
	if err != nil {
		log.Println("解析系统订单ID发生错误", err)
		systemOrderID = ""
	}

	// 查询订单状态
	var p = alipay.NewPayload("alipay.trade.query")
	p.AddBizField("out_trade_no", notification.OutTradeNo)

	var rsp *alipay.TradeQueryRsp
	if err := clients.AlipayClient.Request(c, p, &rsp); err != nil {
		log.Printf("异步通知验证订单 %s 信息发生错误: %s", notification.OutTradeNo, err.Error())
		c.String(http.StatusBadRequest, "fail")
		return
	}

	if rsp.IsFailure() {
		log.Printf("异步通知验证订单 %s 信息发生错误: %s-%s", notification.OutTradeNo, rsp.Msg, rsp.SubMsg)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	// 在这里处理系统订单状态更新
	if systemOrderID != "" {
		log.Printf("【异步通知】系统订单 %s 支付成功，支付宝交易号: %s", systemOrderID, notification.OutTradeNo)
		orderIDInt, err := strconv.ParseInt(systemOrderID, 10, 32)
		if err != nil {
			log.Panicln(err.Error())
		}
		req := &rpc_order.UpdateOrderStatusReq{
			OrderId: int32(orderIDInt),
			Status:  int32(consts.OrderStatusPaid),
		}
		_, err = clients.OrderClient.UpdateOrderStatus(c.Request.Context(), req)
		if err != nil {
			log.Panicln(err.Error())
		}
		log.Println("订单状态修改成功")
	} else {
		log.Printf("【异步通知】支付成功，支付宝交易号: %s，但未获取到系统订单ID", notification.OutTradeNo)
	}

	// 返回成功结果给支付宝
	clients.AlipayClient.ACKNotification(c.Writer)
}

// func configServer() *server.Hertz {
// 	address := conf.GetConf().Hertz.Address
// 	s := fmt.Sprintf("localhost%s", address)
// 	return server.New(server.WithHostPorts(s))
// }
// func handleGetUrl(ctx context.Context, c *app.RequestContext) {
// 	// 从请求中获取系统订单ID
// 	orderID := c.Query("order_id")
// 	// 生成支付宝交易号
// 	tradeNo := fmt.Sprintf("alipay_%d", xid.Next())
// 	// TODO 根据ID查询订单详情
// 	p := alipay.TradePagePay{}
// 	p.NotifyURL = notifyURL
// 	p.ReturnURL = callbackURL
// 	// TODO 和订单一致
// 	p.Subject = "支付订单:" + "qwerqwer"
// 	p.OutTradeNo = tradeNo
// 	// TODO 和订单一致
// 	p.TotalAmount = "20.00"
// 	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
// 	// 使用 PassbackParams 传递系统订单ID
// 	// 需要进行URL编码，因为这个参数会通过URL传递
// 	p.PassbackParams = url.QueryEscape(orderID)
// 	url, err := clients.AliPayClient.TradePagePay(p)
// 	if err != nil {
// 		utils.Fail(c, "创建支付链接失败")
// 		return
// 	}
// 	payURL := url.String()
// 	log.Printf("创建支付链接成功, 系统订单ID: %s, 支付宝交易号: %s", orderID, tradeNo)
// 	fmt.Println(payURL)
// 	// 重定向到支付页面
// 	// utils.Success(c, utils.H{"payURL": payURL})
// 	// 重定向到支付页面
// 	c.Redirect(http.StatusFound, []byte(payURL))
// }

// // 先notify，成功了再会调用callback
// func handleCallback(ctx context.Context, c *app.RequestContext) {
// 	log.Println("callback")
// 	var form struct {
// 		AppID          string `form:"app_id"`
// 		AuthAppID      string `form:"auth_app_id"`
// 		Charset        string `form:"charset"`
// 		Method         string `form:"method"`
// 		OutTradeNo     string `form:"out_trade_no"`
// 		SellerID       string `form:"seller_id"`
// 		Sign           string `form:"sign"`
// 		PassbackParams string `form:"passback_params"`
// 		SignType       string `form:"sign_type"`
// 		Timestamp      string `form:"timestamp"`
// 		TotalAmount    string `form:"total_amount"`
// 		TradeNo        string `form:"trade_no"`
// 		Version        string `form:"version"`
// 	}
// 	// 解析请求参数
// 	if err := c.Bind(&form); err != nil {
// 		log.Println("解析请求参数发生错误", err)
// 		c.String(http.StatusBadRequest, "fail")
// 		return
// 	}
// 	// 获取系统订单ID (PassbackParams)
// 	systemOrderID, err := url.QueryUnescape(form.PassbackParams)
// 	if err != nil {
// 		log.Println("解析系统订单ID发生错误", err.Error())
// 		systemOrderID = ""
// 		return
// 	}
// 	log.Println("准备验证签名")
// 	var values url.Values = url.Values{}
// 	t := reflect.TypeOf(form)
// 	v := reflect.ValueOf(form)
// 	for i := range t.NumField() {
// 		fieldType := t.Field(i)
// 		fieldValue := v.Field(i)
// 		// 获取字段的标签值（即form标签的值）
// 		fieldTag := fieldType.Tag.Get("form")
// 		// 打印字段的名称和值
// 		// fmt.Printf("字段名: %s, 字段值: %v\n", fieldTag, fieldValue)
// 		// 将字段名和值添加到 url.Values 中
// 		values.Add(fieldTag, fieldValue.String())
// 	}
// 	// 验证签名
// 	if err := clients.AliPayClient.VerifySign(values); err != nil {
// 		log.Println("回调验证签名发生错误", err.Error())
// 		utils.Fail(c, err.Error())
// 		return
// 	}
// 	log.Println("回调验证签名通过")
// 	// 查询订单状态
// 	p := alipay.TradeQuery{}
// 	p.OutTradeNo = form.OutTradeNo

// 	rsp, err := clients.AliPayClient.TradeQuery(ctx, p)
// 	if err != nil {
// 		log.Printf("验证订单 %s 信息发生错误: %s", form.OutTradeNo, err.Error())
// 		return
// 	}
// 	if rsp.IsFailure() {
// 		log.Printf("验证订单 %s 信息发生错误: %s-%s", form.OutTradeNo, rsp.Msg, rsp.SubMsg)
// 		return
// 	}
// 	// 在这里处理系统订单状态更新
// 	if systemOrderID != "" {
// 		log.Printf("系统订单 %s 支付成功，支付宝交易号: %s", systemOrderID, form.OutTradeNo)
// 		// TODO: 在这里更新你系统中的订单状态
// 		// updateOrderStatus(systemOrderID, "PAID")
// 	}
// }

// func handleNotify(ctx context.Context, c *app.RequestContext) {
// 	log.Println("notify")
// 	var form struct {
// 		AppID          string              `form:"app_id"`
// 		AuthAppID      string              `form:"auth_app_id"`
// 		BuyerID        string              `form:"buyer_id"`
// 		BuyerPayAmount string              `form:"buyer_pay_amount"`
// 		Charset        string              `form:"charset"`
// 		FundBillList   []map[string]string `form:"fund_bill_list"`
// 		GMTCreate      string              `form:"gmt_create"`
// 		GMTPayMent     string              `form:"gmt_payment"`
// 		InvoiceAmount  string              `form:"invoice_amount"`
// 		NotifyID       string              `form:"notify_id"`
// 		NotifyTime     string              `form:"notify_time"`
// 		NotifyType     string              `form:"notify_type"`
// 		OutTradeNo     string              `form:"out_trade_no"`
// 		PassbackParams string              `form:"passback_params"`
// 		PointAmount    string              `form:"point_amount"`
// 		ReceiptAmount  string              `form:"receipt_amount"`
// 		SellerID       string              `form:"seller_id"`
// 		Sign           string              `form:"sign"`
// 		SignType       string              `form:"sign_type"`
// 		Subject        string              `form:"subject"`
// 		TotalAmount    string              `form:"total_amount"`
// 		TradeNo        string              `form:"trade_no"`
// 		TradeStatus    string              `form:"trade_status"`
// 		Version        string              `form:"version"`
// 	}
// 	// 解析请求参数
// 	if err := c.Bind(&form); err != nil {
// 		log.Println("解析请求参数发生错误", err)
// 		c.String(http.StatusBadRequest, "fail")
// 		return
// 	}
// 	var values url.Values = url.Values{}
// 	t := reflect.TypeOf(form)
// 	v := reflect.ValueOf(form)
// 	for i := range t.NumField() {
// 		fieldType := t.Field(i)
// 		fieldValue := v.Field(i)
// 		// 获取字段的标签值（即form标签的值）
// 		fieldTag := fieldType.Tag.Get("form")
// 		// 将字段名和值添加到 url.Values 中
// 		values.Add(fieldTag, fieldValue.String())
// 		// 打印字段的名称和值
// 		// fmt.Printf("字段名: %s, 字段值: %v\n", fieldTag, fieldValue)
// 	}
// 	fmt.Printf("form in notify: %+v", values)
// 	// 解析异步通知
// 	notification, err := clients.AliPayClient.DecodeNotification(values)
// 	if err != nil {
// 		log.Println("解析异步通知发生错误", err)
// 		c.String(http.StatusBadRequest, "fail")
// 		return
// 	}

// 	log.Println("解析异步通知成功:", notification.NotifyId)

// 	// 获取系统订单ID
// 	systemOrderID, err := url.QueryUnescape(notification.PassbackParams)
// 	if err != nil {
// 		log.Println("解析系统订单ID发生错误", err)
// 		systemOrderID = ""
// 	}

// 	// 查询订单状态
// 	var p = alipay.NewPayload("alipay.trade.query")
// 	p.AddBizField("out_trade_no", notification.OutTradeNo)

// 	var rsp *alipay.TradeQueryRsp
// 	if err := clients.AliPayClient.Request(ctx, p, &rsp); err != nil {
// 		log.Printf("异步通知验证订单 %s 信息发生错误: %s", notification.OutTradeNo, err.Error())
// 		c.String(http.StatusBadRequest, "fail")
// 		return
// 	}

// 	if rsp.IsFailure() {
// 		log.Printf("异步通知验证订单 %s 信息发生错误: %s-%s", notification.OutTradeNo, rsp.Msg, rsp.SubMsg)
// 		c.String(http.StatusBadRequest, "fail")
// 		return
// 	}

// 	// 在这里处理系统订单状态更新
// 	if systemOrderID != "" {
// 		log.Printf("【异步通知】系统订单 %s 支付成功，支付宝交易号: %s", systemOrderID, notification.OutTradeNo)
// 		// TODO: 在这里更新你系统中的订单状态
// 		// updateOrderStatus(systemOrderID, "PAID")
// 	} else {
// 		log.Printf("【异步通知】支付成功，支付宝交易号: %s，但未获取到系统订单ID", notification.OutTradeNo)
// 	}

// 	// 返回成功结果给支付宝
// 	// 这里只能需要根据源码逻辑自行添加了
// 	// clients.AliPayClient.ACKNotification()
// 	c.Response.SetStatusCode(http.StatusOK)
// 	c.Response.BodyBuffer().Write([]byte("success"))
// }

// func main() {
// 	h := configServer()
// 	{
// 		// http://127.0.0.1:8080/alipay/pay?order_id=12
// 		h.GET("/alipay/pay", handleGetUrl)
// 		h.GET("/alipay/callback", handleCallback)
// 		h.POST("/alipay/notify", handleNotify)
// 	}

// 	h.Spin()
// }
