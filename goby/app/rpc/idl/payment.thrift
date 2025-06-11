namespace go payment

// 这里根据实际需要引入相关 Thrift 文件，
// 以便在需要时可以使用其中的类型或结构。
// 如果仅用到 ID，就可以不在结构中强依赖它们。
include "user.thrift"
include "product.thrift"
include "order.thrift"

/**
 * 定义支付状态枚举
 */
enum PaymentStatus {
    PENDING   = 0,  // 待支付
    SUCCESS   = 1,  // 支付成功
    FAILED    = 2,  // 支付失败
    CANCELLED = 3,  // 支付已取消
    REFUNDED  = 4   // 已退款
}

/**
 * 支付记录结构体
 * - 只存储了 user_id / order_id 用于外部关联
 * - 如果需要更详细的第三方支付回调信息、交易流水号等，可自行补充
 */
struct Payment {
    1: i32 id;              // 支付记录主键ID
    2: i32 user_id;         // 对应的用户ID (外键自 user.UserService)
    3: i32 order_id;        // 对应的订单ID (外键自 order.OrderService)
    4: double amount;       // 支付金额
    5: PaymentStatus status; // 支付状态
    6: string created_at;   // 创建时间
    7: string updated_at;   // 更新时间
}

/**
 * 创建支付请求 & 响应
 */
struct CreatePaymentRequest {
    1: i32 user_id;   // 关联的用户ID
    2: i32 order_id;  // 关联的订单ID
    3: double amount; // 需要支付的金额
}

struct CreatePaymentResponse {
    1: Payment payment; // 创建后的支付记录
}

/**
 * 获取支付请求 & 响应
 */
struct GetPaymentRequest {
    1: i32 payment_id; // 需要获取的支付记录ID
}

struct GetPaymentResponse {
    1: Payment payment;
}

/**
 * 更新支付请求 & 响应
 * - 可能只更新状态，也可能在业务中需要更新其他字段
 */
struct UpdatePaymentRequest {
    1: i32 payment_id;       // 待更新的支付记录ID
    2: PaymentStatus status; // 更新后的支付状态（示例：从PENDING更新为SUCCESS等）
}

struct UpdatePaymentResponse {
    1: Payment payment;
}

/**
 * 删除支付请求 & 响应
 * - 一般支付记录不建议物理删除，此处仅作演示
 */
struct DeletePaymentRequest {
    1: i32 payment_id;
}

struct DeletePaymentResponse {
    1: bool success;
}

/**
 * （可选）查询用户所有支付记录
 */
// struct ListUserPaymentsRequest {
//     1: i32 user_id;
// }

// struct ListUserPaymentsResponse {
//     1: list<Payment> payments;
// }

/**
 * PaymentService
 * - 提供最常见的增删改查方法
 * - 可根据需要增加“退款接口”、“回调接口”等
 */
service PaymentService {
    // 创建支付记录
    CreatePaymentResponse createPayment(1: CreatePaymentRequest req);

    // 获取单条支付记录
    GetPaymentResponse getPayment(1: GetPaymentRequest req);

    // 更新支付记录
    UpdatePaymentResponse updatePayment(1: UpdatePaymentRequest req);

    // 删除支付记录（如无此需求可去掉）
    DeletePaymentResponse deletePayment(1: DeletePaymentRequest req);

    // 获取某个用户的所有支付记录（可选）
    // ListUserPaymentsResponse listUserPayments(1: ListUserPaymentsRequest req); 
}
