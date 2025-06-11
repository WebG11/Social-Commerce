package main

import (
	"context"

	cart "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/cart"
	order "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order"
	product "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product"
	user "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user"
)

// OrderServiceImpl implements the last service interface defined in the IDL.
type OrderServiceImpl struct{}

// CreateOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) CreateOrder(ctx context.Context, req *order.CreateOrderReq) (resp *order.CreateOrderResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateOrderStatus implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) UpdateOrderStatus(ctx context.Context, req *order.UpdateOrderStatusReq) (resp *order.UpdateOrderStatusResp, err error) {
	// TODO: Your code here...
	return
}

// GetOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) GetOrder(ctx context.Context, req *order.GetOrderReq) (resp *order.GetOrderResp, err error) {
	// TODO: Your code here...
	return
}

// ListUserOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) ListUserOrder(ctx context.Context, req *order.ListOrderReq) (resp *order.ListOrderResp, err error) {
	// TODO: Your code here...
	return
}

// AdminListOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) AdminListOrder(ctx context.Context, req *order.ListOrderReq) (resp *order.ListOrderResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateOrderAddressResp implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) UpdateOrderAddressResp(ctx context.Context, req *order.UpdateOrderAddressReq) (resp *order.UpdateOrderAddressResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateOrderAddress implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) UpdateOrderAddress(ctx context.Context, req *order.UpdateOrderAddressReq) (resp *order.UpdateOrderAddressResp, err error) {
	// TODO: Your code here...
	return
}

// CreateUserAddress implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) CreateUserAddress(ctx context.Context, req *order.CreateUserAddressReq) (resp *order.CreateUserAddressResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteUserAddress implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) DeleteUserAddress(ctx context.Context, req *order.DeleteUserAddressReq) (resp *order.DeleteUserAddressResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateUserAddress implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) UpdateUserAddress(ctx context.Context, req *order.UpdateUserAddressReq) (resp *order.UpdateUserAddressResp, err error) {
	// TODO: Your code here...
	return
}

// GetUserAddress implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) GetUserAddress(ctx context.Context, req *order.GetUserAddressReq) (resp *order.GetUserAddressResp, err error) {
	// TODO: Your code here...
	return
}

type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	// TODO: Your code here...
	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {
	// TODO: Your code here...
	return
}

// GetUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.GetUserReq) (resp *user.GetUserResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *user.UpdateUserReq) (resp *user.UpdateUserResp, err error) {
	// TODO: Your code here...
	return
}

// RemoveUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) RemoveUser(ctx context.Context, req *user.RemoveUserReq) (resp *user.RemoveUserResp, err error) {
	// TODO: Your code here...
	return
}

// BlockUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) BlockUser(ctx context.Context, req *user.BlockUserReq) (resp *user.BlockUserResp, err error) {
	// TODO: Your code here...
	return
}

// UnblockUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) UnblockUser(ctx context.Context, req *user.UnblockUserReq) (resp *user.UnblockUserResp, err error) {
	// TODO: Your code here...
	return
}

// AdminListUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) AdminListUser(ctx context.Context, req *user.AdminListUserReq) (resp *user.AdminListUserResp, err error) {
	// TODO: Your code here...
	return
}

type ProductServiceImpl struct{}

// CreateProduct implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) CreateProduct(ctx context.Context, req *product.CreateProductReq) (resp *product.CreateProductResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateProduct implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) UpdateProduct(ctx context.Context, req *product.UpdateProductReq) (resp *product.UpdateProductResp, err error) {
	// TODO: Your code here...
	return
}

// RemoveProduct implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) RemoveProduct(ctx context.Context, req *product.RemoveProductReq) (resp *product.RemoveProductResp, err error) {
	// TODO: Your code here...
	return
}

// GetProductByID implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) GetProductByID(ctx context.Context, req *product.GetProductByIDReq) (resp *product.GetProductByIDResp, err error) {
	// TODO: Your code here...
	return
}

// ListProduct implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) ListProduct(ctx context.Context, req *product.ListProductReq) (resp *product.ListProductResp, err error) {
	// TODO: Your code here...
	return
}

// AdminListProduct implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) AdminListProduct(ctx context.Context, req *product.ListProductReq) (resp *product.ListProductResp, err error) {
	// TODO: Your code here...
	return
}

// SearchProducts implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) SearchProducts(ctx context.Context, req *product.SearchProductsReq) (resp *product.SearchProductsResp, err error) {
	// TODO: Your code here...
	return
}

// Seller implements the UserServiceImpl interface.
func (s *UserServiceImpl) Seller(ctx context.Context, req *user.SellerReq) (resp *user.SellerResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateOrderTracking implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) UpdateOrderTracking(ctx context.Context, req *order.UpdateOrderTrackingReq) (resp *order.UpdateOrderTrackingResp, err error) {
	// TODO: Your code here...
	return
}

// CreateProductReview implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) CreateProductReview(ctx context.Context, req *product.CreateProductReviewReq) (resp *product.CreateProductReviewResp, err error) {
	// TODO: Your code here...
	return
}

// GetProductReviews implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) GetProductReviews(ctx context.Context, req *product.GetProductReviewsReq) (resp *product.GetProductReviewsResp, err error) {
	// TODO: Your code here...
	return
}

type CartServiceImpl struct{}

// CreateItem implements the CartServiceImpl interface.
func (s *CartServiceImpl) CreateItem(ctx context.Context, req *cart.CreateItemReq) (resp *cart.CreateItemResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateQuantity implements the CartServiceImpl interface.
func (s *CartServiceImpl) UpdateQuantity(ctx context.Context, req *cart.UpdateQuantityReq) (resp *cart.UpdateQuantityResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteItem implements the CartServiceImpl interface.
func (s *CartServiceImpl) DeleteItem(ctx context.Context, req *cart.DeleteItemReq) (resp *cart.DeleteItemResp, err error) {
	// TODO: Your code here...
	return
}

// GetItem implements the CartServiceImpl interface.
func (s *CartServiceImpl) GetItem(ctx context.Context, req *cart.GetItemReq) (resp *cart.GetItemResp, err error) {
	// TODO: Your code here...
	return
}

// ListItem implements the CartServiceImpl interface.
func (s *CartServiceImpl) ListItem(ctx context.Context, req *cart.ListItemReq) (resp *cart.ListItemResp, err error) {
	// TODO: Your code here...
	return
}

// CreatePromotion implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) CreatePromotion(ctx context.Context, req *product.CreatePromotionReq) (resp *product.CreatePromotionResp, err error) {
	// TODO: Your code here...
	return
}

// GetActivePromotions implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) GetActivePromotions(ctx context.Context) (resp *product.GetActivePromotionsResp, err error) {
	// TODO: Your code here...
	return
}

// DeletePromotion implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) DeletePromotion(ctx context.Context, req *product.DeletePromotionReq) (resp *product.DeletePromotionResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateOrderDiscount implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) UpdateOrderDiscount(ctx context.Context, req *order.UpdateOrderDiscountReq) (resp *order.UpdateOrderDiscountResp, err error) {
	// TODO: Your code here...
	return
}

// GetSalesReport implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) GetSalesReport(ctx context.Context, req *order.SalesReportReq) (resp *order.SalesReportResp, err error) {
	// TODO: Your code here...
	return
}

// CheckStock implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) CheckStock(ctx context.Context, req *product.CheckStockReq) (resp *product.CheckStockResp, err error) {
	// TODO: Your code here...
	return
}

// GetSalesReportByDate implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) GetSalesReportByDate(ctx context.Context, req *order.SalesReportByDateReq) (resp *order.SalesReportByDateResp, err error) {
	// TODO: Your code here...
	return
}
