package main

import (
	"context"

	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dal/tidb"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dao"
	"github.com/bitdance-panic/gobuy/app/services/gateway/middleware"
	gutils "github.com/bitdance-panic/gobuy/app/services/gateway/utils"
	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// RefreshTokenHandler 处理 Token 续期
func RefreshTokenHandler(ctx context.Context, c *app.RequestContext) {
	// 从请求头获取 Access Token
	accessToken := string(c.GetHeader("Authorization"))
	if accessToken == "" {
		utils.FailFull(c, consts.StatusUnauthorized, "Missing Access Token.", nil)
		return
	}

	// 验证 access_token 是否有效
	claims, err := middleware.AuthMiddleware.GetClaimsFromJWT(ctx, c)
	if err == nil {
		// access_token 仍然有效，不需要刷新
		utils.FailFull(c, consts.StatusOK, "Access Token is still valid", nil)
		return
	}

	// access_token 过期，提取 uid
	userID := int(claims["uid"].(float64))

	// 查询数据库中的 refresh_token
	storedRefreshToken, err := dao.GetRefreshTokenByUserID(tidb.DB, userID)
	if err != nil {
		hlog.Warnf("Refresh Token not found for user %d", userID)
		utils.FailFull(c, consts.StatusUnauthorized, "Refresh Token Not Found", nil)
		return
	}

	// 验证 Refresh Token 是否过期
	if gutils.IsRefreshTokenExpired(storedRefreshToken) {
		utils.FailFull(c, consts.StatusUnauthorized, "Refresh Token Expired, Please Login Again", nil)
		return
	}

	// 返回新的 access_token
	tokenString, _, err := middleware.AuthMiddleware.TokenGenerator(userID)
	if err != nil {
		utils.FailFull(c, consts.StatusInternalServerError, "Token Generation Failed", nil)
		return
	}
	utils.Success(c, map[string]interface{}{
		"access_token": tokenString,
	})
}
