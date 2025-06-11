package main

import (
	"context"

	rpc_product "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product"
	"github.com/bitdance-panic/gobuy/app/services/product/biz/bll"
)

type ProductServiceImpl struct{}

func (*ProductServiceImpl) ListProduct(ctx context.Context) (*rpc_product.ListProductResp, error) {
	return bll.ListProduct(ctx)
}

func (*ProductServiceImpl) GetProductByID(ctx context.Context, req *rpc_product.GetProductByIDReq) (*rpc_product.GetProductByIDResp, error) {
	return bll.GetProductByID(ctx, req)
}

func (*ProductServiceImpl) CreateProduct(ctx context.Context, req *rpc_product.CreateProductReq) (*rpc_product.CreateProductResp, error) {
	return bll.CreateProduct(ctx, req)
}

func (*ProductServiceImpl) UpdateProduct(ctx context.Context, req *rpc_product.UpdateProductReq) (*rpc_product.UpdateProductResp, error) {
	return bll.UpdateProduct(ctx, req)
}

func (*ProductServiceImpl) RemoveProduct(ctx context.Context, req *rpc_product.RemoveProductReq) (*rpc_product.RemoveProductResp, error) {
	return bll.RemoveProduct(ctx, req)
}

func (*ProductServiceImpl) AdminListProduct(ctx context.Context, req *rpc_product.ListProductReq) (resp *rpc_product.ListProductResp, err error) {
	return bll.AdminListProduct(ctx, req)
}

func (*ProductServiceImpl) SearchProducts(ctx context.Context, req *rpc_product.SearchProductsReq) (*rpc_product.SearchProductsResp, error) {
	return bll.SearchProducts(ctx, req)
}

func (*ProductServiceImpl) CreateProductReview(ctx context.Context, req *rpc_product.CreateProductReviewReq) (*rpc_product.CreateProductReviewResp, error) {
	return bll.CreateProductReview(ctx, req)
}

func (*ProductServiceImpl) GetProductReviews(ctx context.Context, req *rpc_product.GetProductReviewsReq) (*rpc_product.GetProductReviewsResp, error) {
	return bll.GetProductReviews(ctx, req)
}

func (*ProductServiceImpl) CreatePromotion(ctx context.Context, req *rpc_product.CreatePromotionReq) (*rpc_product.CreatePromotionResp, error) {
	return bll.CreatePromotion(ctx, req)
}

func (*ProductServiceImpl) GetActivePromotions(ctx context.Context) (*rpc_product.GetActivePromotionsResp, error) {
	return bll.GetActivePromotions(ctx)
}

func (*ProductServiceImpl) DeletePromotion(ctx context.Context, req *rpc_product.DeletePromotionReq) (*rpc_product.DeletePromotionResp, error) {
	return bll.DeletePromotion(ctx, req)
}

func (*ProductServiceImpl) CheckStock(ctx context.Context, req *rpc_product.CheckStockReq) (*rpc_product.CheckStockResp, error) {
	return bll.CheckStock(ctx, req)
}
