package clients

import (
	"github.com/bitdance-panic/gobuy/app/clients"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/cart/cartservice"
	"github.com/bitdance-panic/gobuy/app/services/order/conf"
)

var CartClient cartservice.Client

func NewClients() {
	registryConf := conf.GetConf().Registry
	CartClient = clients.NewCartClient(registryConf.RegistryAddress[0])
}
