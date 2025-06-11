package bll

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/bitdance-panic/gobuy/app/models"
	rpc_product "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product"
	"github.com/bitdance-panic/gobuy/app/services/product/biz/dal/tidb"
	"github.com/bitdance-panic/gobuy/app/services/product/biz/dao"
)

func ListProduct(ctx context.Context) (*rpc_product.ListProductResp, error) {
	products, err := dao.ListAll(tidb.DB)
	if err != nil {
		return nil, err
	}
	protoProducts := make([]*rpc_product.Product, len(products))
	for i, v := range products {
		protoProducts[i] = convertToProtoProduct(&v)
	}
	return &rpc_product.ListProductResp{
		Products:   protoProducts,
		TotalCount: int64(len(products)),
	}, nil
}

func AdminListProduct(ctx context.Context, req *rpc_product.ListProductReq) (*rpc_product.ListProductResp, error) {
	p, total, err := dao.AdminList(tidb.DB, int(req.PageNum), int(req.PageSize), int(req.SellerId))
	if err != nil {
		return nil, err
	}
	protoProducts := make([]*rpc_product.Product, len(*p))
	for i, v := range *p {
		protoProducts[i] = convertToProtoProduct(&v)
	}
	return &rpc_product.ListProductResp{
		Products:   protoProducts,
		TotalCount: total,
	}, nil
}

func GetProductByID(ctx context.Context, req *rpc_product.GetProductByIDReq) (*rpc_product.GetProductByIDResp, error) {
	p, err := dao.GetByID(tidb.DB, int(req.Id))
	if err != nil {
		return nil, err
	}
	if p == nil {
		return &rpc_product.GetProductByIDResp{
			Product: nil,
		}, nil
	}
	protoProduct := convertToProtoProduct(p)
	return &rpc_product.GetProductByIDResp{
		Product: protoProduct,
	}, nil
}

func CreateProduct(ctx context.Context, req *rpc_product.CreateProductReq) (*rpc_product.CreateProductResp, error) {
	p := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int(req.Stock),
		Image:       req.Image,
		SellerID:    int(req.SellerId),
	}
	if err := dao.Create(tidb.DB, &p); err != nil {
		return nil, err
	}
	protoProduct := convertToProtoProduct(&p)
	return &rpc_product.CreateProductResp{
		Product: protoProduct,
	}, nil
}

func UpdateProduct(ctx context.Context, req *rpc_product.UpdateProductReq) (*rpc_product.UpdateProductResp, error) {
	p, err := dao.GetByID(tidb.DB, int(req.Id))
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		p.Name = req.Name
	}
	if req.Description != "" {
		p.Description = req.Description
	}
	if req.Price != 0 {
		p.Price = req.Price
	}
	if req.Stock != 0 {
		p.Stock = int(req.Stock)
	}
	if req.Image != "" {
		p.Image = req.Image
	}
	if err := dao.Update(tidb.DB, p); err != nil {
		return nil, err
	}
	protoProduct := convertToProtoProduct(p)
	return &rpc_product.UpdateProductResp{
		Product: protoProduct,
	}, nil
}

func RemoveProduct(ctx context.Context, req *rpc_product.RemoveProductReq) (*rpc_product.RemoveProductResp, error) {
	err := dao.Remove(tidb.DB, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &rpc_product.RemoveProductResp{Success: true}, nil
}

func SearchProducts(ctx context.Context, req *rpc_product.SearchProductsReq) (*rpc_product.SearchProductsResp, error) {
	// 处理 minPrice 的默认值
	minPrice := 0
	if req.MinPrice != 0 {
		minPrice = int(req.MinPrice)
	}

	// 处理 maxPrice 的默认值
	maxPrice := math.MaxInt32
	if req.MaxPrice != 0 {
		maxPrice = int(req.MaxPrice)
	}

	products, total, err := dao.Search(tidb.DB, req.Query, req.Category, int(req.PageNum), int(req.PageSize), minPrice, maxPrice, req.Brand)
	if err != nil {
		return nil, err
	}

	protoProducts := make([]*rpc_product.Product, 0, len(products))
	for _, p := range products {
		protoProducts = append(protoProducts, convertToProtoProduct(&p))
	}

	return &rpc_product.SearchProductsResp{
		Products:   protoProducts,
		TotalCount: total,
	}, nil
}

func GetProductReviews(ctx context.Context, req *rpc_product.GetProductReviewsReq) (*rpc_product.GetProductReviewsResp, error) {
	reviews, err := dao.GetReviews(tidb.DB, int(req.ProductId))
	if err != nil {
		return nil, err
	}
	protoReviews := make([]*rpc_product.ProductReview, len(*reviews))
	for i, v := range *reviews {
		protoReviews[i] = convertToProtoProductReview(&v)
	}
	return &rpc_product.GetProductReviewsResp{
		Reviews: protoReviews,
	}, nil
}

func convertToProtoProductReview(review *models.ProductReview) *rpc_product.ProductReview {
	return &rpc_product.ProductReview{
		Id:        int32(review.ID),
		ProductId: int32(review.ProductID),
		UserId:    int32(review.UserID),
		UserName:  review.UserName,
		Rating:    int32(review.Rating),
		Comment:   review.Comment,
		CreatedAt: review.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func CreateProductReview(ctx context.Context, req *rpc_product.CreateProductReviewReq) (*rpc_product.CreateProductReviewResp, error) {
	review := models.ProductReview{
		ProductID: int(req.ProductId),
		UserID:    int(req.UserId),
		UserName:  req.UserName,
		Rating:    int(req.Rating),
		Comment:   req.Comment,
	}
	if err := dao.CreateReview(tidb.DB, &review); err != nil {
		return nil, err
	}
	return &rpc_product.CreateProductReviewResp{
		Success: true,
	}, nil
}

func convertToProtoProduct(p *models.Product) *rpc_product.Product {
	return &rpc_product.Product{
		Id:          int32(p.ID),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       int32(p.Stock),
		Image:       p.Image,
		SellerId:    int64(p.SellerID),
		CreatedAt:   p.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func CreatePromotion(ctx context.Context, req *rpc_product.CreatePromotionReq) (*rpc_product.CreatePromotionResp, error) {
	startTime, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, err
	}
	endTime, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, err
	}
	promotion := models.Promotion{
		Type:         req.Type,
		ProductID:    int(req.ProductId),
		DiscountRate: req.DiscountRate,
		CouponCode:   req.CouponCode,
		CouponType:   req.CouponType,
		CouponValue:  req.CouponValue,
		MinPurchase:  req.MinPurchase,
		UsageLimit:   int(req.UsageLimit),
		StartDate:    startTime,
		EndDate:      endTime,
	}
	if err := dao.CreatePromotion(tidb.DB, &promotion); err != nil {
		return nil, err
	}
	return &rpc_product.CreatePromotionResp{
		Success: true,
	}, nil
}

func GetActivePromotions(ctx context.Context) (*rpc_product.GetActivePromotionsResp, error) {
	promotions, err := dao.GetActivePromotions(tidb.DB)
	if err != nil {
		return nil, err
	}
	protoPromotions := make([]*rpc_product.Promotion, len(promotions))
	for i, v := range promotions {
		protoPromotions[i] = convertToProtoPromotion(&v)
	}
	return &rpc_product.GetActivePromotionsResp{
		Promotions: protoPromotions,
	}, nil
}

func convertToProtoPromotion(p *models.Promotion) *rpc_product.Promotion {
	return &rpc_product.Promotion{
		Id:           strconv.Itoa(p.ID),
		Type:         p.Type,
		ProductId:    strconv.Itoa(p.ProductID),
		DiscountRate: strconv.FormatFloat(p.DiscountRate, 'f', -1, 64),
		CouponCode:   p.CouponCode,
		CouponType:   p.CouponType,
		CouponValue:  strconv.FormatFloat(p.CouponValue, 'f', -1, 64),
		MinPurchase:  strconv.FormatFloat(p.MinPurchase, 'f', -1, 64),
		UsageLimit:   strconv.Itoa(p.UsageLimit),
		StartDate:    p.StartDate.Format("2006-01-02"),
		EndDate:      p.EndDate.Format("2006-01-02"),
	}
}

func DeletePromotion(ctx context.Context, req *rpc_product.DeletePromotionReq) (*rpc_product.DeletePromotionResp, error) {
	err := dao.DeletePromotion(tidb.DB, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &rpc_product.DeletePromotionResp{
		Success: true,
	}, nil
}

func CheckStock(ctx context.Context, req *rpc_product.CheckStockReq) (*rpc_product.CheckStockResp, error) {
	// 处理 minStock 的默认值
	minStock := 0
	if !req.IsSetMinStock() {
		minStock = 0
	} else {
		minStock = int(req.GetMinStock())
	}

	// 处理 maxStock 的默认值
	maxStock := math.MaxInt32 // 如果未设置最大库存，默认为最大整数值
	if !req.IsSetMaxStock() {
		maxStock = math.MaxInt32
	} else {
		maxStock = int(req.GetMaxStock())
	}

	// 调用 DAO 层获取符合条件的商品列表
	products, err := dao.GetProductsByStockRange(tidb.DB, &minStock, &maxStock)
	if err != nil {
		return nil, err
	}

	var protoProducts []*rpc_product.Product
	lowStockWarning := false
	for _, p := range products {
		protoProduct := convertToProtoProduct(&p)
		protoProducts = append(protoProducts, protoProduct)

		// 判断是否存在低库存警告
		if p.Stock <= 10 { // 假设低库存阈值为 10
			lowStockWarning = true
		}
	}

	// 构造返回的响应
	return &rpc_product.CheckStockResp{
		Products:        protoProducts,
		LowStockWarning: lowStockWarning,
	}, nil
}
