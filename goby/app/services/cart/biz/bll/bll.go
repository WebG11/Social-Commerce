package bll

import (
	"context"
	"errors"

	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/cart"
	rpc_cart "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/cart"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product"
	rpc_product "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product"
	"github.com/bitdance-panic/gobuy/app/services/cart/biz/clients"
	"github.com/bitdance-panic/gobuy/app/services/cart/biz/dal/tidb"
	"github.com/bitdance-panic/gobuy/app/services/cart/biz/dao"
)

type CartBLL struct{}

func CreateItem(ctx context.Context, req *rpc_cart.CreateItemReq) (*rpc_cart.CreateItemResp, error) {
	// 先检查商品是否存在
	getProductResp, err := clients.ProductClient.GetProductByID(ctx, &rpc_product.GetProductByIDReq{Id: int32(req.ProductId)})
	if err != nil {
		return nil, err
	}
	if getProductResp.Product == nil {
		return nil, errors.New("product not found")
	}
	// 再添加
	ok, err := dao.Create(tidb.DB, int(req.UserId), int(req.ProductId))
	if err != nil {
		return nil, err
	}
	return &rpc_cart.CreateItemResp{Success: ok}, nil
}

func DeleteItem(ctx context.Context, req *rpc_cart.DeleteItemReq) (*rpc_cart.DeleteItemResp, error) {
	err := dao.Delete(tidb.DB, int(req.ItemId))
	if err != nil {
		return &cart.DeleteItemResp{Success: false}, err
	}
	return &cart.DeleteItemResp{Success: true}, nil
}

func GetItem(ctx context.Context, req *rpc_cart.GetItemReq) (*rpc_cart.GetItemResp, error) {
	item, err := dao.GetItemByID(tidb.DB, int(req.ItemId))
	if err != nil {
		return nil, err
	}
	protoItem := convertToProtoCartItem(item)
	return &cart.GetItemResp{Item: protoItem}, nil
}

func UpdateQuantity(ctx context.Context, req *rpc_cart.UpdateQuantityReq) (*rpc_cart.UpdateQuantityResp, error) {
	item, err := dao.GetItemByID(tidb.DB, int(req.ItemId))
	if err != nil {
		return nil, err
	}
	getProductResp, err := clients.ProductClient.GetProductByID(ctx, &product.GetProductByIDReq{Id: int32(item.ProductID)})
	if err != nil {
		return nil, err
	}
	// 1. 商品已被删除
	if getProductResp.Product == nil {
		return nil, errors.New("product not found")
	}
	// 否则直接改，在订单结算时再考虑
	err = dao.UpdateQuantity(tidb.DB, item, int(req.NewQuantity_))
	if err != nil {
		return nil, err
	}
	return &rpc_cart.UpdateQuantityResp{
		Success: true,
	}, nil
}

func ListItem(ctx context.Context, req *rpc_cart.ListItemReq) (*rpc_cart.ListItemResp, error) {
	items, err := dao.ListItemsByUserID(tidb.DB, int(req.UserId))
	if err != nil {
		return nil, err
	}
	protoItems := make([]*rpc_cart.CartItem, 0, len(*items))
	for _, item := range *items {
		protoItem := convertToProtoCartItem(&item)
		protoItems = append(protoItems, protoItem)
	}
	return &rpc_cart.ListItemResp{
		Items: protoItems,
	}, nil
}

func convertToProtoCartItem(item *models.CartItem) *rpc_cart.CartItem {
	valid := item.Product.IsDeleted
	return &rpc_cart.CartItem{
		Id:        int32(item.ID),
		Name:      item.Product.Name,
		Price:     item.Product.Price,
		Quantity:  int32(item.Quantity),
		Image:     item.Product.Image,
		Valid:     valid,
		ProductId: int32(item.ProductID),
		SellerId:  int32(item.Product.SellerID),
	}
}
