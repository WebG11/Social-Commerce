package clients

import (
	"github.com/bitdance-panic/gobuy/app/clients"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order/orderservice"
	"github.com/bitdance-panic/gobuy/app/services/paycallback/conf"
	"github.com/smartwalle/alipay/v3"
)

var (
	AlipayClient *alipay.Client
	OrderClient  orderservice.Client
)

func init() {
	payConf := conf.GetConf().Alipay
	AlipayClient = clients.NewAlipayClient(payConf.APPID, payConf.PrivateKey)
	addr := conf.GetConf().Registry.RegistryAddress[0]
	OrderClient = clients.NewOrderClient(addr)
}
