namespace go cart

struct CartItem {
    1: i32 id;
    2: string name;
    3: double price;
    4: i32 quantity;
    5: string image;
    6: bool valid;
    7: i32 product_id;
    8: i32 seller_id
}

service CartService {
    CreateItemResp createItem(1: CreateItemReq req);
    UpdateQuantityResp updateQuantity(1: UpdateQuantityReq req);
    DeleteItemResp deleteItem(1: DeleteItemReq req);
    GetItemResp getItem(1: GetItemReq req);
    ListItemResp listItem(1: ListItemReq req);
}

struct CreateItemReq {
    1: i32 product_id;
    2: i32 user_id;
}

struct CreateItemResp {
    1: bool success;
}

struct UpdateQuantityReq {
    1: i32 item_id;
    2: i32 new_quantity;
}

struct UpdateQuantityResp {
    1: bool success;
}

struct DeleteItemReq {
    1: i32 item_id;
}

struct DeleteItemResp {
    1: bool success;
}

struct ListItemReq {
    1: i32 user_id;
    //2: i32 page_num;
    //3: i32 page_size;
}

struct ListItemResp {
    1: list<CartItem> items;
}

struct GetItemReq {
    1: i32 item_id;
}

struct GetItemResp {
    1: CartItem item;
}