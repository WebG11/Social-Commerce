package handlers

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/bitdance-panic/gobuy/app/consts"
	rpc_user "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/clients"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dal/redis"
	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client/callopt"
)

// 获取用户信息
func HandleGetUser(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt(consts.CONTEXT_UID_KEY)
	req := rpc_user.GetUserReq{UserId: int32(userID)}
	resp, err := clients.UserClient.GetUser(context.Background(), &req, callopt.WithRPCTimeout(10*time.Second))
	if err != nil || !resp.Success {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{
		"name":      resp.Username,
		"email":     resp.Email,
		"password":  resp.Password,
		"addresses": resp.Addresses,
	})
}

// 更新用户信息
func HandleUpdateUser(ctx context.Context, c *app.RequestContext) {
	req := rpc_user.UpdateUserReq{UserId: int32(c.GetInt(consts.CONTEXT_UID_KEY))}

	if err := c.BindAndValidate(&req); err != nil {
		hlog.Warnf("User update failed for user id: %s, validation error: %v", req.UserId, err)
		utils.Fail(c, err.Error())
		return
	}
	resp, err := clients.UserClient.UpdateUser(context.Background(), &req)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	if !resp.Success {
		utils.Fail(c, "User update failed")
		return
	}
	utils.Success(c, utils.H{"userID": req.UserId})
}

// 获取用户列表
func HandleAdminListUser(ctx context.Context, c *app.RequestContext) {
	pageNum, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	pageSize, err := strconv.Atoi(c.Query("size"))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	req := rpc_user.AdminListUserReq{
		PageNum:  int32(pageNum),
		PageSize: int32(pageSize),
	}
	resp, err := clients.UserClient.AdminListUser(context.Background(), &req, callopt.WithRPCTimeout(10*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{
		"users":       resp.Users,
		"total_count": resp.TotalCount,
	})
}

// 移除用户
func HandleRemoveUser(ctx context.Context, c *app.RequestContext) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	req := rpc_user.RemoveUserReq{
		UserId: int32(userID),
	}
	resp, err := clients.UserClient.RemoveUser(context.Background(), &req, callopt.WithRPCTimeout(10*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{
		"success": resp.Success,
	})
}

// 封禁用户
// @Summary 添加黑名单条目
// @Description 封禁指定用户/IP
// @Accept json
// @Produce json
// @Router /admin/block_user [post]
func HandleBlockUser(ctx context.Context, c *app.RequestContext) {
	req := rpc_user.BlockUserReq{}

	if err := c.BindAndValidate(&req); err != nil {
		hlog.Warnf("Block user failed for: %s, validation error: %v", req.Identifier, err)
		utils.Fail(c, err.Error())
		return
	}

	resp, err := clients.UserClient.BlockUser(context.Background(), &req, callopt.WithRPCTimeout(310*time.Second))

	if err != nil || !resp.Success {
		utils.Fail(c, "封禁失败: "+err.Error())
		return
	}

	// 更新Redis
	data, _ := json.Marshal(map[string]interface{}{
		"reason":     req.Reason,
		"expires_at": req.ExpiresAt,
	})
	pipe := redis.RedisClient.Pipeline()
	pipe.HSet(ctx, "blacklist:entries", req.Identifier, data)
	if req.ExpiresAt != 0 {
		expirationTime := time.Unix(req.ExpiresAt, 0)
		pipe.ZAdd(ctx, "blacklist:expiry", &redis.Z{
			Score:  float64(expirationTime.UnixNano()),
			Member: req.Identifier,
		})
	}
	if _, err := pipe.Exec(ctx); err != nil {
		utils.Fail(c, "Redis更新失败: "+err.Error())
		return
	}

	utils.Success(c, utils.H{"blockID": resp.BlockId})
}

// 解封用户
// @Summary 移除黑名单条目
// @Description 解除封禁
// @Accept json
// @Produce json
// @Router /admin/unblock [delete]
func HandleUnblockUser(ctx context.Context, c *app.RequestContext) {
	req := rpc_user.UnblockUserReq{
		Identifier: c.Param("identifier"),
	}

	resp, err := clients.UserClient.UnblockUser(context.Background(), &req, callopt.WithRPCTimeout(10*time.Second))

	if err != nil || !resp.Success {
		utils.Fail(c, "解禁失败: "+err.Error())
		return
	}

	// 更新Redis
	pipe := redis.RedisClient.Pipeline()
	pipe.HDel(ctx, "blacklist:entries", req.Identifier)
	pipe.ZRem(ctx, "blacklist:expiry", req.Identifier)
	if _, err := pipe.Exec(ctx); err != nil {
		utils.Fail(c, "Redis删除失败: "+err.Error())
		return
	}

	utils.Success(c, nil)
}
