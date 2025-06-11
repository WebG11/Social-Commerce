package clients

import (
	"github.com/bitdance-panic/gobuy/app/clients"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product/productservice"
	"github.com/bitdance-panic/gobuy/app/services/cart/conf"
)

var ProductClient productservice.Client

func init() {
	registryConf := conf.GetConf().Registry
	ProductClient = clients.NewProductClient(registryConf.RegistryAddress[0])
}
