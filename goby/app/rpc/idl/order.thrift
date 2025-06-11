namespace go order

service OrderService{
    CreateOrderResp createOrder(1: CreateOrderReq req);
    UpdateOrderStatusResp updateOrderStatus(1: UpdateOrderStatusReq req);
    UpdateOrderAddressResp updateOrderAddress(1:UpdateOrderAddressReq req);
    GetOrderResp getOrder(1: GetOrderReq req);
    ListOrderResp listUserOrder(1: ListOrderReq req);
    ListOrderResp adminListOrder(1: ListSellerOrderReq req);
    CreateUserAddressResp createUserAddress(1: CreateUserAddressReq req);
    DeleteUserAddressResp deleteUserAddress(1: DeleteUserAddressReq req);
    UpdateUserAddressResp updateUserAddress(1: UpdateUserAddressReq req);
    GetUserAddressResp getUserAddress(1: GetUserAddressReq req);
    UpdateOrderTrackingResp UpdateOrderTracking(1: UpdateOrderTrackingReq req);
    UpdateOrderDiscountResp UpdateOrderDiscount(1: UpdateOrderDiscountReq req);
    SalesReportResp GetSalesReport(1: SalesReportReq req);
    SalesReportByDateResp GetSalesReportByDate(1: SalesReportByDateReq req);
}

struct OrderItem{
    1: i32 order_id;
    2: i32 product_id;
    3: bool item_status
    4: string shipped_at
    5: string tracking_number
    6: string product_name;
    7: i32 quantity;
    8: double price;
    9: string product_image;
}

struct Order{
    1: i32 id;
    2: i32 user_id;
    // 订单号
    3: string number;
    4: double total_price;
    5: double discount;
    6: i32 status;
    7: list<OrderItem> items;
    8: string created_at;
    9: string pay_time;
    10: string name;
    11: string phone;
    12: string address;
    13: string tracking_number;
    14: i32 seller_id;
}

struct UserAddress{
    1: i32 id;
    2: i32 user_id;
    3: string name;
    4: string phone;
    5: string address;
}

struct CreateUserAddressReq{
    1: i32 user_id;
    2: string phone;
    3: string user_address;
}

struct CreateUserAddressResp{
    1: i32 user_id;
    2: bool success;
}

struct DeleteUserAddressReq{
    1: i32 user_id;
}

struct DeleteUserAddressResp{
    1: i32 user_id;
    2: bool success;
}

struct UpdateUserAddressReq{
    1: i32 user_id;
    2: string user_address;
}

struct UpdateUserAddressResp{
    1: string user_address;
    2: bool success;
}

struct GetUserAddressReq{
    1: i32 user_id;
}

struct GetUserAddressResp{
    1: list<UserAddress> user_addresses;
}

struct UpdateOrderTrackingReq{
    1: i32 order_id;
    2: string tracking_number;
}

struct UpdateOrderTrackingResp{
    1: bool success;
}

struct UpdateOrderDiscountReq {
    1: i32 order_id;
    2: double discount;
}

struct UpdateOrderDiscountResp {
    1: bool success;
}
struct CreateOrderReq{
    1: i32 user_id;
    2: list<i32> cartItemIDs;
    3: string phone;
    4: string order_address;
    5: double TotalPrice;
}
struct CreateOrderResp{
    1: Order order;
}

struct UpdateOrderAddressReq{
    1: i32 order_id;
    2: string address_name;
    3: string address_phone;
    4: string address;
}
struct UpdateOrderAddressResp{
    1: string order_address;
    2: bool success;
}

struct UpdateOrderStatusReq{
    1: i32 order_id;
    2: i32 status;
}
struct UpdateOrderStatusResp{
    1: i32 status;
    2: bool success;
}
struct ListOrderReq{
    1: i32 user_id;
    2: i32 page_num;
    3: i32 page_size;
}
struct ListSellerOrderReq{
    1: i32 seller_id;
    2: i32 page_num;
    3: i32 page_size;
}
struct ListOrderResp{
    1: list<Order> orders;
    2: i64 total_count
}
struct GetOrderReq{
    1: i32 order_id;
}
struct GetOrderResp{
    1: Order order;
}

struct SalesReportReq {
    1: optional string start_date, // 可选参数，用于过滤开始日期
    2: optional string end_date,    // 可选参数，用于过滤结束日期
    3: i32 SellerID,
}

struct SalesReportResp {
    1: required double total_revenue,
    2: required i32 order_count,
    3: optional map<string, i32> top_products // 最畅销的产品及其销量
    4: required double average_orderAmt // 平均订单金额
    5: i32 TotalProductCount,
    6: map<string, double> DailyRevenue,
}

struct SalesReportByDateReq {
}

struct SalesReportByDateResp {
    1: optional map<string, double> date_revenue
}