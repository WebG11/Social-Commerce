package middleware

import (
	"context"

	app_consts "github.com/bitdance-panic/gobuy/app/consts"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
)

func AddUidMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if skip, exists := c.Get("skip_auth"); exists && skip.(bool) {
			c.Next(ctx) // 白名单跳过认证
			return
		}

		if claims := jwt.ExtractClaims(ctx, c); claims != nil {
			userID := claims[IdentityKey]
			c.Set(app_consts.CONTEXT_UID_KEY, int(userID.(float64)))
		}
		// c.Set(app_consts.CONTEXT_UID_KEY, 1)
		c.Next(ctx)
	}
}
