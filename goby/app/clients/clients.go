package clients

import (
	"log"

	"github.com/bitdance-panic/gobuy/app/common/clientsuite"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/agent/agentservice"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/cart/cartservice"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order/orderservice"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product/productservice"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user/userservice"
	"github.com/cloudwego/kitex/client"
	"github.com/smartwalle/alipay/v3"
)

func NewUserClient(registryAddress string) userservice.Client {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: "user",
			RegistryAddr:       registryAddress,
		}),
	}
	client, err := userservice.NewClient(
		"user",
		opts...,
	)
	if err != nil {
		log.Fatalf("faild to new user client: %v", err)
	}
	return client
}

func NewCartClient(registryAddress string) cartservice.Client {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: "cart",
			RegistryAddr:       registryAddress,
		}),
	}
	client, err := cartservice.NewClient(
		"cart",
		opts...,
	)
	if err != nil {
		log.Fatalf("faild to new cart client: %v", err)
	}
	return client
}

func NewAgentClient(registryAddress string) agentservice.Client {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: "agent",
			RegistryAddr:       registryAddress,
		}),
	}
	client, err := agentservice.NewClient(
		"agent",
		opts...,
	)
	if err != nil {
		log.Fatalf("faild to new agent client: %v", err)
	}
	return client
}

func NewOrderClient(registryAddress string) orderservice.Client {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: "order",
			RegistryAddr:       registryAddress,
		}),
	}
	client, err := orderservice.NewClient(
		"order",
		opts...,
	)
	if err != nil {
		log.Fatalf("faild to new order client: %v", err)
	}
	return client
}

func NewProductClient(registryAddress string) productservice.Client {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: "product",
			RegistryAddr:       registryAddress,
		}),
	}
	client, err := productservice.NewClient(
		"product",
		opts...,
	)
	if err != nil {
		log.Fatalf("faild to new product client: %v", err)
	}
	return client
}

func NewAlipayClient(appID string, privateKey string) *alipay.Client {
	var client *alipay.Client
	var err error
	// 支付宝提供了用于开发时测试的 sandbox 环境，对接的时候需要注意相关的 app id 和密钥是 sandbox 环境还是 production 环境的。如果是 sandbox 环境，本参数应该传 false，否则为 true。
	if client, err = alipay.New(appID, privateKey, false); err != nil {
		log.Panic(err)
	}
	// 加载证书
	// 加载应用公钥证书
	if err = client.LoadAppCertPublicKeyFromFile("conf/appPublicCert.crt"); err != nil {
		log.Panic(err)
	}
	// 加载支付宝根证书
	if err = client.LoadAliPayRootCertFromFile("conf/alipayRootCert.crt"); err != nil {
		log.Panic(err)
	}
	// 加载支付宝公钥证书
	if err = client.LoadAlipayCertPublicKeyFromFile("conf/alipayPublicCert.crt"); err != nil {
		log.Panic(err)
	}
	return client
}
