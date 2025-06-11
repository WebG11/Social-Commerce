package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/bitdance-panic/gobuy/app/consts"

	rpc_cart "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/cart"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/clients"
	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client/callopt"
)

// HandleAddToCart 这是更新商品
// @Summary 这是一段Summary
// @Description 这是一段Description
// @Accept application/json
// @Produce application/json
// @Router /product [put]
func HandleCreateCartItem(ctx context.Context, c *app.RequestContext) {
	productID, err := strconv.Atoi(c.Param("productID"))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	userID := c.GetInt(consts.CONTEXT_UID_KEY)
	req := rpc_cart.CreateItemReq{
		UserId:    int32(userID),
		ProductId: int32(productID),
	}
	_, err = clients.CartClient.CreateItem(context.Background(), &req, callopt.WithRPCTimeout(5*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	// 返回是否成功即可
	utils.Success(c, utils.H{})
}

// handleProductPost 这是创建商品
// @Summary 这是一段Summary
// @Description 这是一段Description
// @Accept application/json
// @Produce application/json
// @Router /product [post]
func HandleListCartItem(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt(consts.CONTEXT_UID_KEY)
	// pageNum, err := strconv.Atoi(c.Query("page"))
	// if err != nil {
	// 	utils.Fail(c, err.Error())
	// 	return
	// }
	// pageSize, err := strconv.Atoi(c.Query("size"))
	// if err != nil {
	// 	utils.Fail(c, err.Error())
	// 	return
	// }
	req := rpc_cart.ListItemReq{
		UserId: int32(userID),
		// PageNum:  int32(pageNum),
		// PageSize: int32(pageSize),
	}
	resp, err := clients.CartClient.ListItem(context.Background(), &req, callopt.WithRPCTimeout(5*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"items": resp.Items})
}

// handleProductDELETE 这是删除商品
// @Summary 这是一段Summary
// @Description 这是一段Description
// @Accept application/json
// @Produce application/json
// @Router /product [delete]
func HandleDeleteCartItem(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.Atoi(c.Param("itemID"))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	req := rpc_cart.DeleteItemReq{
		ItemId: int32(id),
	}
	resp, err := clients.CartClient.DeleteItem(context.Background(), &req, callopt.WithRPCTimeout(5*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	if resp.Success {
		utils.Success(c, nil)
	} else {
		utils.Fail(c, "删除失败")
	}
}

// handleProductGet 这是获取一个商品
// @Summary 这是一段Summary
// @Description 这是一段Description
// @Accept application/json
// @Produce application/json
// @Router /product [get]
func HandleUpdateCartItemQuantity(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.Atoi(c.Param("itemID"))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	newQuantity, err := strconv.Atoi(c.Query("quantity"))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	req := rpc_cart.UpdateQuantityReq{
		ItemId:       int32(id),
		NewQuantity_: int32(newQuantity),
	}
	resp, err := clients.CartClient.UpdateQuantity(context.Background(), &req, callopt.WithRPCTimeout(5*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"success": resp.Success})
}
