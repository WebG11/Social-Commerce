package handlers

import (
	"context"
	"time"

	gconsts "github.com/bitdance-panic/gobuy/app/consts"
	rpc_user "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/clients"
	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client/callopt"
)

// 用户注册
func HandleRegister(ctx context.Context, c *app.RequestContext) {
	req := rpc_user.RegisterReq{}

	// 从请求体中绑定参数并验证
	if err := c.BindAndValidate(&req); err != nil {
		hlog.Warnf("User register failed, validation error: %v", err)
		utils.Fail(c, err.Error())
		return
	}

	resp, err := clients.UserClient.Register(context.Background(), &req, callopt.WithRPCTimeout(5*time.Second))
	if err != nil || !resp.Success {
		utils.FailFull(c, consts.StatusCreated, err.Error(), map[string]interface{}{"status": 201})
		return
	}
	utils.Success(c, utils.H{"userID": resp.UserId})
}

func HandleSeller(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt(gconsts.CONTEXT_UID_KEY)
	req := rpc_user.SellerReq{
		UserId: int32(userID),
	}

	resp, err := clients.UserClient.Seller(context.Background(), &req, callopt.WithRPCTimeout(5*time.Second))
	if err != nil {
		utils.FailFull(c, consts.StatusCreated, err.Error(), map[string]interface{}{"status": 201})
		return
	}
	utils.Success(c, utils.H{"isSeller": resp.IsSeller})
}
