package clients

import (
	"github.com/bitdance-panic/gobuy/app/clients"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/agent/agentservice"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/cart/cartservice"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order/orderservice"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product/productservice"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user/userservice"
	"github.com/bitdance-panic/gobuy/app/services/gateway/conf"
	"github.com/smartwalle/alipay/v3"
	// payclient "github.com/bitdance-panic/gobuy/app/services/payment/biz/clients"
)

var (
	UserClient    userservice.Client
	ProductClient productservice.Client
	// PaymentClient paymentservice.Client  // 暂未完成
	OrderClient  orderservice.Client
	CartClient   cartservice.Client
	AgentClient  agentservice.Client
	AlipayClient *alipay.Client
)

func init() {
	addr := conf.GetConf().Registry.RegistryAddress[0]
	ProductClient = clients.NewProductClient(addr)
	// PaymentClient = payclient.PaymentClient  // 暂未完成
	OrderClient = clients.NewOrderClient(addr)
	CartClient = clients.NewCartClient(addr)
	UserClient = clients.NewUserClient(addr)
	AgentClient = clients.NewAgentClient(addr)
	payConf := conf.GetConf().Alipay
	AlipayClient = clients.NewAlipayClient(payConf.APPID, payConf.PrivateKey)
}
