package main

import (
	"context"
	"log"
	"net"
	"strings"

	"github.com/bitdance-panic/gobuy/app/common/mtl"
	"github.com/bitdance-panic/gobuy/app/common/serversuite"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order/orderservice"
	"github.com/bitdance-panic/gobuy/app/services/order/biz/clients"
	"github.com/bitdance-panic/gobuy/app/services/order/biz/crontask"
	"github.com/bitdance-panic/gobuy/app/services/order/biz/dal"
	"github.com/bitdance-panic/gobuy/app/services/order/conf"
	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/robfig/cron/v3"

	"github.com/cloudwego/kitex/server"
)

var (
	ServiceName  = "order" // conf.GetConf().Kitex.Service
	RegistryAddr = conf.GetConf().Registry.RegistryAddress[0]
)

func kitexInit() (opts []server.Option) {
	address := conf.GetConf().Kitex.Address
	if strings.HasPrefix(address, ":") {
		localIp := utils.MustGetLocalIPv4()
		address = localIp + address
	}
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		panic(err)
	}
	// opts = append(opts, server.WithServiceAddr(addr))
	opts = append(opts, server.WithServiceAddr(addr), server.WithSuite(serversuite.CommonServerSuite{
		CurrentServiceName: ServiceName,
		RegistryAddr:       RegistryAddr,
	}))
	return
}

func main() {
	// 初始化指标监控
	mtl.InitMetric(ServiceName, conf.GetConf().Kitex.MetricsPort, RegistryAddr)
	p := mtl.InitTracing(ServiceName)
	defer p.Shutdown(context.Background())

	dal.Init()

	clients.NewClients()

	opts := kitexInit()

	svr := orderservice.NewServer(new(OrderServiceImpl), opts...)

	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
	startExpireCron()
}

func startExpireCron() {
	// 创建一个新的cron实例
	c := cron.New()
	// 添加定时任务，每分钟执行一次
	_, err := c.AddFunc("@every 1m", crontask.CheckAndUpdateExpiredOrders)
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}
	// 启动cron
	c.Start()
	// 保持主程序运行
	select {}
}
