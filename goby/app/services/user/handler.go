package main

// cd app/rpc kitex -module github.com/bitdance-panic/gobuy/app/rpc -service user idl/user.thrift
import (
	"context"

	rpc_user "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user"
	"github.com/bitdance-panic/gobuy/app/services/user/biz/bll"
)

type UserServiceImpl struct{}

func (*UserServiceImpl) Register(ctx context.Context, req *rpc_user.RegisterReq) (resp *rpc_user.RegisterResp, err error) {
	return bll.Register(ctx, req)
}

func (*UserServiceImpl) Login(ctx context.Context, req *rpc_user.LoginReq) (*rpc_user.LoginResp, error) {
	return bll.Login(ctx, req)
}

// 获取用户信息
func (*UserServiceImpl) GetUser(ctx context.Context, req *rpc_user.GetUserReq) (*rpc_user.GetUserResp, error) {
	return bll.GetUser(ctx, req)
}

func (*UserServiceImpl) AdminListUser(ctx context.Context, req *rpc_user.AdminListUserReq) (*rpc_user.AdminListUserResp, error) {
	return bll.AdminListUser(ctx, req)
}

// 更新用户信息
func (*UserServiceImpl) UpdateUser(ctx context.Context, req *rpc_user.UpdateUserReq) (*rpc_user.UpdateUserResp, error) {
	return bll.UpdateUser(ctx, req)
}

// 删除用户
func (*UserServiceImpl) RemoveUser(ctx context.Context, req *rpc_user.RemoveUserReq) (*rpc_user.RemoveUserResp, error) {
	return bll.RemoveUser(ctx, req)
}

// 封禁用户
func (*UserServiceImpl) BlockUser(ctx context.Context, req *rpc_user.BlockUserReq) (*rpc_user.BlockUserResp, error) {
	return bll.BlockUser(ctx, req)
}

// 解禁用户
func (*UserServiceImpl) UnblockUser(ctx context.Context, req *rpc_user.UnblockUserReq) (*rpc_user.UnblockUserResp, error) {
	return bll.UnblockUser(ctx, req)
}

func (*UserServiceImpl) Seller(ctx context.Context, req *rpc_user.SellerReq) (*rpc_user.SellerResp, error) {
	return bll.Seller(ctx, req)
}
