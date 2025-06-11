package main

import (
	"context"

	"github.com/bitdance-panic/gobuy/app/services/cart/biz/bll"

	rpc_cart "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/cart"
)

type CartServiceImpl struct{}

func (*CartServiceImpl) CreateItem(ctx context.Context, req *rpc_cart.CreateItemReq) (resp *rpc_cart.CreateItemResp, err error) {
	return bll.CreateItem(ctx, req)
}

func (*CartServiceImpl) DeleteItem(ctx context.Context, req *rpc_cart.DeleteItemReq) (resp *rpc_cart.DeleteItemResp, err error) {
	return bll.DeleteItem(ctx, req)
}

func (*CartServiceImpl) UpdateQuantity(ctx context.Context, req *rpc_cart.UpdateQuantityReq) (resp *rpc_cart.UpdateQuantityResp, err error) {
	return bll.UpdateQuantity(ctx, req)
}

func (*CartServiceImpl) ListItem(ctx context.Context, req *rpc_cart.ListItemReq) (resp *rpc_cart.ListItemResp, err error) {
	return bll.ListItem(ctx, req)
}

func (*CartServiceImpl) GetItem(ctx context.Context, req *rpc_cart.GetItemReq) (resp *rpc_cart.GetItemResp, err error) {
	return bll.GetItem(ctx, req)
}
