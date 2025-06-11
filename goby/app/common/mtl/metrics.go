package mtl

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/server"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Registry *prometheus.Registry // 全局 Prometheus 注册表

func InitMetric(serviceName, metricsPort, registryAddr string) (registry.Registry, *registry.Info) {
	// 初始化 prometheus 注册表
	Registry = prometheus.NewRegistry()
	Registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	// 注册到 Consul
	r, _ := consul.NewConsulRegister(registryAddr)

	if strings.HasPrefix(metricsPort, ":") {
		localIp := utils.MustGetLocalIPv4()
		// localIp := "0.0.0.0"
		metricsPort = localIp + metricsPort
	}
	addr, _ := net.ResolveTCPAddr("tcp", metricsPort)
	registryInfo := &registry.Info{
		ServiceName: "prometheus",
		Addr:        addr,
		Weight:      1,
		Tags:        map[string]string{"service": serviceName},
	}

	_ = r.Register(registryInfo)

	server.RegisterShutdownHook(func() {
		if err := r.Deregister(registryInfo); err != nil {
			panic(err)
		}
	})

	// 暴露指标端点
	http.Handle("/metrics", promhttp.HandlerFor(Registry, promhttp.HandlerOpts{}))
	go func() {
		log.Printf("[Metrics] Serving metrics at %s/metrics", metricsPort)
		if err := http.ListenAndServe(metricsPort, nil); err != nil {
			panic(err)
		}
	}()
	return r, registryInfo
}
