package middleware

import (
	"context"
	"time"

	app_consts "github.com/bitdance-panic/gobuy/app/consts"
	"github.com/bitdance-panic/gobuy/app/models"
	rpc_user "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	//"github.com/bitdance-panic/gobuy/app/services/gateway/conf"

	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/clients"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dal/tidb"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dao"
	gutils "github.com/bitdance-panic/gobuy/app/services/gateway/utils"
	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/hertz-contrib/jwt"
)

type User = models.User

var IdentityKey = app_consts.CONTEXT_UID_KEY

var secret = "panic-bitdance"

const (
	AccessTokenExpire  = 12 * time.Hour
	RefreshTokenExpire = 7 * 24 * time.Hour
)

var AuthMiddleware *jwt.HertzJWTMiddleware

// BlackList      = gutils.NewSyncSet()

func init() {
	// JWT中间件配置
	initJWTMiddleware()
}

func initJWTMiddleware() {
	var err error
	AuthMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:            "gobuy auth",
		SigningAlgorithm: "HS256",
		Key:              []byte(secret),
		Timeout:          AccessTokenExpire,
		MaxRefresh:       RefreshTokenExpire,
		// 从请求头的Authorization提取token
		TokenLookup:                 "header: Authorization",
		TokenHeadName:               "Bearer",
		SendAuthorization:           true,
		WithoutDefaultTokenHeadName: false,
		TimeFunc:                    time.Now,
		HTTPStatusMessageFunc: func(e error, ctx context.Context, c *app.RequestContext) string {
			return e.Error()
		},
		DisabledAbort: false,

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			userID := data

			return jwt.MapClaims{
				IdentityKey: userID,
			}
		},
		Authenticator: authenticate,
		LoginResponse: loginResponse,
		// Authorizator: authorize,
		Unauthorized: unauthorizedHandler,
		LogoutResponse: func(ctx context.Context, c *app.RequestContext, code int) {
			utils.Success(c, nil)
		},
		RefreshResponse: refreshResponse,
		IdentityKey:     IdentityKey,
		IdentityHandler: identityHandler,
	})

	_ = AuthMiddleware
	if err != nil {
		hlog.Fatalf("auth middleware Error:" + err.Error())
	}
}

// 认证处理，用于用户登录时提取并验证登录凭据
func authenticate(ctx context.Context, c *app.RequestContext) (interface{}, error) {
	loginReq := rpc_user.LoginReq{}

	if err := c.Bind(&loginReq); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}
	loginResp, err := clients.UserClient.Login(context.Background(), &loginReq, callopt.WithRPCTimeout(5*time.Second))
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}
	userID := int(loginResp.UserId)
	c.Set(app_consts.CONTEXT_UID_KEY, userID)
	c.Set("is_seller", loginReq.IsSeller)
	return loginResp.UserId, nil
}

// 登录响应
func loginResponse(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
	userID := c.GetInt(app_consts.CONTEXT_UID_KEY)
	// 生成RefreshToken
	refreshToken, err := gutils.GenerateRefreshToken(userID)
	if err != nil {
		utils.FailFull(c, consts.StatusInternalServerError, "Failed to generate refreshtoken", nil)
		return
	}

	// 保存refreshToken
	if err := dao.UpdateRefreshToken(tidb.DB, userID, refreshToken); err != nil {
		hlog.Errorf("Failed to save refresh token, error: %v", err)
		utils.FailFull(c, consts.StatusInternalServerError, "Failed to store refreshtoken", nil)
		return
	}

	req := rpc_user.GetUserReq{
		UserId:      int32(userID),
		IsSeller:    c.GetBool("is_seller"),
		CheckSeller: true,
	}

	loginResp, err := clients.UserClient.GetUser(context.Background(), &req, callopt.WithRPCTimeout(5*time.Second))
	if err != nil || !loginResp.Success {
		utils.FailFull(c, 401, "Incorrent", nil)
		return
	}

	utils.Success(c, map[string]interface{}{
		"access_token":  token,
		"refresh_token": refreshToken,
		"email":         loginResp.Email,
		"username":      loginResp.Username,
		"userId":        userID,
	})
}

// 授权处理
// func authorize(data interface{}, ctx context.Context, c *app.RequestContext) bool {
// 	uid, ok := data.(int)
// 	if !ok {
// 		return false
// 	}

// refreshResponse 刷新Token响应
func refreshResponse(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
	// 生成新的RefreshToken
	refreshToken, err := gutils.GenerateRefreshToken(c.GetInt(app_consts.CONTEXT_UID_KEY))
	if err != nil {
		utils.FailFull(c, consts.StatusInternalServerError, "Failed to generate refreshtoken", nil)
		return
	}

	// 更新数据库中的RefreshToken
	if err := dao.UpdateRefreshToken(tidb.DB, c.GetInt(app_consts.CONTEXT_UID_KEY), refreshToken); err != nil {
		hlog.Errorf("Failed to save refresh token, error: %v", err)
		utils.FailFull(c, consts.StatusInternalServerError, "Failed to store refreshtoken, error: "+err.Error(), nil)
		return
	}

	// 返回新的双Token
	utils.Success(c, map[string]interface{}{
		"access_token": token,
		"expire":       expire.Unix(),
	})
}

// 身份处理
func identityHandler(ctx context.Context, c *app.RequestContext) interface{} {
	claims := jwt.ExtractClaims(ctx, c)
	return claims[app_consts.CONTEXT_UID_KEY]
}

// 统一错误处理
func unauthorizedHandler(ctx context.Context, c *app.RequestContext, code int, message string) {
	utils.FailFull(c, code, message, nil)
	c.Abort()
}

func ConditionalAuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if skip, exists := c.Get("skip_auth"); exists && skip.(bool) {
			c.Next(ctx) // 跳过认证
			return
		}
		AuthMiddleware.MiddlewareFunc()(ctx, c) // 执行认证
	}
}
