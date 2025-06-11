package utils

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Response 是统一的响应结构体
type Response struct {
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

type H map[string]interface{}

func Success(c *app.RequestContext, data map[string]interface{}) {
	response := Response{
		Message: "success",
		Data:    data,
	}
	c.JSON(consts.StatusOK, response)
}

func Fail(c *app.RequestContext, message string) {
	response := Response{
		Message: message,
	}
	c.JSON(consts.StatusInternalServerError, response)
}

// 比如鉴权的就不一样
func FailFull(c *app.RequestContext, code int, message string, data map[string]interface{}) {
	response := Response{
		Message: message,
		Data:    data,
	}
	c.JSON(code, response)
}
