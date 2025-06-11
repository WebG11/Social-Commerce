package main

import (
	"github.com/bitdance-panic/gobuy/app/services/paycallback/conf"
	"github.com/gin-gonic/gin"
)

func main() {
	// 路由函数
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		// 支付宝沙箱环境的网关域名
		// allowedOrigin := "https://openapi.alipaydev.com"
		// // 获取请求的来源
		// origin := c.Request.Header.Get("Origin")
		// fmt.Println(origin)
		// // 仅允许支付宝沙箱环境的域名访问
		// if origin == allowedOrigin {
		// 	c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		// } else {
		// 	c.AbortWithStatusJSON(403, gin.H{"error": "Forbidden"})
		// 	return
		// }
		// 其他CORS配置
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
	})
	// http://127.0.0.1:8280/alipay/pay?order_id=150003
	r.GET("/alipay/callback", handleCallback)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/alipay/notify", handleNotify)
	hertzConf := conf.GetConf().Hertz
	r.Run(hertzConf.Address)
}
