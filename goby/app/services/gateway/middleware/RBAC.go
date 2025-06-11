package middleware

import (
	"context"
	"net/http"

	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dal/tidb"
	"github.com/cloudwego/hertz/pkg/app"
)

func RBACMiddleware(requiredRole string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 获取用户信息
		claims, _ := AuthMiddleware.GetClaimsFromJWT(ctx, c)

		userID := int(claims["uid"].(float64))

		// 查询用户角色
		var roles []string
		tidb.DB.Model(&models.UserRole{}).
			Joins("JOIN role ON user_role.role_id = role.id").
			Where("user_role.user_id = ?", userID).
			Pluck("role.name", &roles)

		// 检查是否包含所需角色
		hasRole := false
		for _, role := range roles {
			if role == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, map[string]interface{}{
				"code":    http.StatusForbidden,
				"message": "需要 " + requiredRole + " 权限",
			})
			return
		}

		c.Next(ctx)
	}
}
