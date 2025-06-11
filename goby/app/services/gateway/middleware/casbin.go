package middleware

import (
	"context"
	"strconv"

	"github.com/bitdance-panic/gobuy/app/services/gateway/casbin"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
)

func CasbinMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 检查是否跳过认证
		if skip, exists := c.Get("skip_auth"); exists && skip.(bool) {
			c.Next(ctx)
			return
		}

		// 如果请求已被终止（如JWT认证失败），直接跳过
		if c.IsAborted() {
			c.Next(ctx)
			return
		}

		// 提取用户ID
		claims := jwt.ExtractClaims(ctx, c)
		userID := claims[IdentityKey].(float64)
		userIDStr := strconv.Itoa(int(userID))

		// 提取请求路径和方法
		path := string(c.URI().Path())
		method := string(c.Request.Method())

		// Casbin权限检查
		ok, err := casbin.Enforcer.Enforce(userIDStr, path, method)
		if err != nil {
			utils.FailFull(c, consts.StatusInternalServerError, "权限校验错误", nil)
			c.Abort()
			return
		}
		if !ok {
			utils.FailFull(c, consts.StatusForbidden, "无访问权限", nil)
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}
