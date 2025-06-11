namespace go user

struct RegisterReq {
    1: string email
    2: string password
    3: string username
    4: bool is_seller
}

struct RegisterResp {
    1: i32 user_id
    2: bool success
    3: 
}

struct SellerReq {
    1: i32 user_id
}

struct SellerResp {
    1: bool is_seller
}

struct LoginReq {
    1: string email
    2: string password
    3: bool is_seller
}

struct LoginResp {
    1: bool success
    2: i32 user_id
}

struct GetUserReq {
    1: i32 user_id
    2: bool is_seller
    3: bool CheckSeller = 0
}


struct GetUserResp {
    1: bool success
    2: string username
    3: string email
    4: string password
    5: list<UserAddress> addresses
    6: bool is_seller
}

struct User {
    1: i32 user_id
    2: string username
    3: string email
    4: string refresh_token
}

struct AdminListUserReq {
    1: i32 page_num
    2: i32 page_size
}


struct AdminListUserResp {
    1: list<User> users
    2: string message
    3: i64 total_count
}

struct UserAddress {
  1: string address_name,
  2: string address_phone,
  3: string address,
}

struct UpdateUserReq {
  1: i32 user_id,
  2: string name,
  3: string email,
  4: string password_hashed
  5: list<UserAddress> addresses,
}

struct UpdateUserResp {
    1: bool success
}

//软删除
struct RemoveUserReq {
    1: i32 user_id
}

struct RemoveUserResp {
    1: bool success
}

struct BlockUserReq {
    1: string identifier 
    2: string reason
    3: i64 expires_at     
}

struct BlockUserResp {
    1: string block_id
    2: bool success
}

struct UnblockUserReq {
    1: string identifier     
}

struct UnblockUserResp {
    1: bool success
}


// 定义服务接口
service UserService {
    RegisterResp register(1: RegisterReq req)
    LoginResp login(1: LoginReq req)
    GetUserResp GetUser(1: GetUserReq req)
    UpdateUserResp updateUser(1: UpdateUserReq req)
    RemoveUserResp removeUser(1: RemoveUserReq req)
    BlockUserResp blockUser(1: BlockUserReq req)
    UnblockUserResp unblockUser(1: UnblockUserReq req)
    AdminListUserResp adminListUser(1: AdminListUserReq req)
    SellerResp Seller(1: SellerReq req)
}