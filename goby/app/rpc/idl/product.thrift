namespace go product


struct Product {
    1: i32 id;
    2: string name;
    3: double price;
    4: i32 stock;
    5: string image;
    6: string description;
    7: string created_at;
    8: bool is_deleted;
    9: i64 seller_id
}

struct ProductReview {
    1: i32 id;
    2: i32 product_id;
    3: i32 user_id;
    4: string user_name;
    5: i32 rating;
    6: string comment;
    7: string created_at;
}

struct Promotion {
    1: string id
    2: string type
    3: string product_id
    4: string discount_rate
    5: string coupon_code
    6: string coupon_type
    7: string coupon_value
    8: string min_purchase
    9: string usage_limit
    10: string start_date
    11: string end_date
}

struct GetActivePromotionsResp {
    1: list<Promotion> promotions
}

service ProductService {
    CreateProductResp createProduct(1: CreateProductReq req);
    UpdateProductResp updateProduct(1: UpdateProductReq req);
    RemoveProductResp removeProduct(1: RemoveProductReq req);
    GetProductByIDResp getProductByID(1: GetProductByIDReq req);
    ListProductResp listProduct();
    ListProductResp adminListProduct(1: ListProductReq req);
    SearchProductsResp searchProducts(1: SearchProductsReq req);
    CreateProductReviewResp createProductReview(1: CreateProductReviewReq req);
    GetProductReviewsResp GetProductReviews(1: GetProductReviewsReq req);
    CreatePromotionResp CreatePromotion(1: CreatePromotionReq req);
    GetActivePromotionsResp GetActivePromotions();
    DeletePromotionResp DeletePromotion(1: DeletePromotionReq req);
    CheckStockResp CheckStock(1: CheckStockReq req);
}

struct DeletePromotionReq {
    1: i32 id;
}

struct DeletePromotionResp {
    1: bool success;
}

struct GetProductReviewsReq {
    1: i32 product_id;
}

struct GetProductReviewsResp {
    1: list<ProductReview> Reviews;
}

struct CreateProductReq {
    1: string name;
    2: string description;
    3: double price;
    4: i32 stock;
    5: string image;
    6: i64 seller_id
}

struct CreateProductResp {
    1: Product product;
}

struct UpdateProductReq {
    1: i32 id;
    2: string name;
    3: string description;
    4: double price;
    5: i32 stock;
    6: string image;
}

struct UpdateProductResp {
    1: Product product;
}

struct RemoveProductReq {
    1: i32 id;
}

struct RemoveProductResp {
    1: bool success;
}

struct GetProductByIDReq {
    1: i32 id;
}

struct GetProductByIDResp {
    1: Product product;
}

struct SearchProductsReq {
    1: string query;
    2: string category;
    3: i32 page_num;
    4: i32 page_size;
    5: i32 min_price;
    6: i32 max_price;
    7: string brand;
}

struct SearchProductsResp {
    1: list<Product> products;
    2: i64 total_count;
}

struct ListProductReq {
    1: i32 page_num;
    2: i32 page_size;
    3: i32 SellerId;
}

struct ListProductResp {
    1: list<Product> products;
    2: bool success;
    3: i64 total_count;
}

struct CreateProductReviewReq {
    1: i32 product_id;
    2: i32 user_id;
    3: string user_name;
    4: i32 rating;
    5: string comment;
}

struct CreateProductReviewResp {
    1: bool success;
}

struct CreatePromotionReq {
    1: string type; 
    2: i32 product_id;
    3: double discount_rate;
    4: string coupon_code;
    5: string coupon_type; 
    6: double coupon_value;
    7: double min_purchase;
    8: i32 usage_limit;
    9: string start_date;
    10: string end_date;
}

struct CreatePromotionResp {
    1:  bool success;
}

struct CheckStockReq {
    1: optional i32 min_stock;
    2: optional i32 max_stock;       
}

struct CheckStockResp {
    1: list<Product> products;
    2: required bool low_stock_warning;
}