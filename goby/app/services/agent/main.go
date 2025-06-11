package main

import (
	"context"
	"log"
	"net"
	"strings"

	"github.com/bitdance-panic/gobuy/app/common/mtl"
	"github.com/bitdance-panic/gobuy/app/common/serversuite"
	"github.com/bitdance-panic/gobuy/app/services/agent/biz/bll"
	"github.com/bitdance-panic/gobuy/app/services/agent/biz/dal"
	"github.com/bitdance-panic/gobuy/app/services/agent/conf"
	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/cloudwego/kitex/server"

	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/agent/agentservice"
)

var (
	ServiceName  = "agent" // conf.GetConf().Kitex.Service
	RegistryAddr = conf.GetConf().Registry.RegistryAddress[0]
)

func kitexInit() (opts []server.Option) {
	// address
	address := conf.GetConf().Kitex.Address
	if strings.HasPrefix(address, ":") {
		localIp := utils.MustGetLocalIPv4()
		// localIp := "0.0.0.0"
		address = localIp + address
	}
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		panic(err)
	}
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
	bll.Init()

	opts := kitexInit()
	svr := agentservice.NewServer(new(AgentServiceImpl), opts...)
	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
