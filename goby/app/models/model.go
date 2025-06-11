package models

import (
	"log"
	"os"
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        int `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDeleted bool
}

// 贫血模型
type User struct {
	Base
	Email          string `gorm:"unique"`
	PasswordHashed string
	Username       string
	RefreshToken   string
	IsSeller       bool

	Roles         []Role `gorm:"many2many:user_role;"` // 用户和角色的多对多关系
	Orders        []Order
	UserAddresses []UserAddress
}

type Role struct {
	Base
	Name        string       `gorm:"unique;not null"`
	Permissions []Permission `gorm:"many2many:role_permission;"` // 角色和权限的多对多关系
}

type Permission struct {
	Base
	Name   string `gorm:"unique;not null"`            // 权限名称，例如 "create_user", "delete_user"
	Path   string `gorm:"not null"`                   // 资源路径，支持通配符如 /product/*
	Method string `gorm:"not null"`                   // 请求方法，支持正则如 GET|POST
	Roles  []Role `gorm:"many2many:role_permission;"` // 权限和角色的多对多关系
}

type UserRole struct {
	UserID int
	RoleID int
}

type RolePermission struct {
	RoleID       int
	PermissionID int
}

type ProductReview struct {
	Base
	ProductID int
	UserID    int
	UserName  string
	Rating    int
	Comment   string
}

type Product struct {
	Base
	SellerID    int
	Name        string
	Price       float64
	Stock       int
	Image       string
	Category    string
	Description string
}

type CartItem struct {
	Base
	UserID    int
	ProductID int
	Quantity  int     `gorm:"not null"`
	Product   Product `gorm:"foreignKey:ProductID"`
}

type Order struct {
	Base
	UserID     int
	Number     string  `gorm:"unique;not null"`
	TotalPrice float64 `gorm:"not null"`
	Discount   float64
	// OrderStatus
	Status         int `gorm:"type:varchar(20);not null"`
	TrackingNumber string
	Items          []OrderItem
	PayTime        *time.Time
	Name           string
	Phone          string
	Address        string
	SellerID       int
}

type OrderItem struct {
	Base
	OrderID      int     // 订单 ID
	ProductID    int     // 商品 ID
	ItemStatus   bool    // 商品状态
	Quantity     int     `gorm:"not null"`             // 商品数量
	Price        float64 `gorm:"not null"`             // 商品单价
	Product      Product `gorm:"foreignKey:ProductID"` // 关联商品
	Order        Order   `gorm:"foreignKey:OrderID"`   // 关联订单
	ProductName  string
	ProductImage string
}

type UserAddress struct {
	Base
	UserID      int
	UserName    string
	Phone       string
	UserAddress string
}

// 黑名单条目模型
type Blacklist struct {
	Base
	Identifier string    `gorm:"type:varchar(255);uniqueIndex;not null"` // 封禁标识（用户ID）
	Reason     string    `gorm:"type:text;not null"`                     // 封禁原因
	ExpiresAt  time.Time `gorm:"index"`                                  // 过期时间（为空表示永久封禁）
}

// 操作类型（用于日志）
type BlacklistOperation string

const (
	AddToBlacklist      BlacklistOperation = "ADD"
	RemoveFromBlacklist BlacklistOperation = "REMOVE"
)

type Promotion struct {
	Base
	SellerID     int
	Type         string
	ProductID    int
	DiscountRate float64
	CouponCode   string
	CouponType   string
	CouponValue  float64
	MinPurchase  float64
	UsageLimit   int
	StartDate    time.Time
	EndDate      time.Time
}

type PromotionProduct struct {
	Base
	PromotionID int
	ProductID   int
	Product     Product `gorm:"foreignKey:ProductID"`
}

type Coupon struct {
	Base
	SellerID     int
	Code         string `gorm:"unique;not null"`
	Name         string `gorm:"not null"`
	Description  string
	DiscountType int       `gorm:"not null"` // 1: 固定金额, 2: 折扣
	Discount     float64   `gorm:"not null"`
	MinAmount    float64   // 最低使用金额
	StartTime    time.Time `gorm:"not null"`
	EndTime      time.Time `gorm:"not null"`
	Status       int       `gorm:"not null"` // 0: 未开始, 1: 可用, 2: 已过期
	UsageLimit   int       // 使用次数限制
	UsedCount    int       // 已使用次数
}

type ChatMessage struct {
	ID         int64     `json:"id"`
	SenderID   int64     `json:"sender_id"`
	ReceiverID int64     `json:"receiver_id"`
	Content    string    `json:"content"`
	Timestamp  time.Time `json:"timestamp"`
}

func (User) TableName() string {
	return "user"
}

func (Role) TableName() string {
	return "role"
}

func (Permission) TableName() string {
	return "permission"
}

func (UserRole) TableName() string {
	return "user_role"
}

func (RolePermission) TableName() string {
	return "role_permission"
}

func (CartItem) TableName() string {
	return "cart_item"
}

func (Order) TableName() string {
	return "order"
}

func (OrderItem) TableName() string {
	return "order_item"
}

func (Product) TableName() string {
	return "product"
}

func (Blacklist) TableName() string {
	return "blacklist"
}

func (UserAddress) TableName() string { return "user_address" }

func (ProductReview) TableName() string { return "product_review" }

func (Promotion) TableName() string {
	return "promotion"
}

func (PromotionProduct) TableName() string {
	return "promotion_product"
}

func (Coupon) TableName() string {
	return "coupon"
}

// 调用这个来自动调整表结构
func AutoMigrate(db *gorm.DB) {
	if os.Getenv("GO_ENV") != "production" {
		log.Println("进行数据表的Migrate")
		//nolint:errcheck
		db.AutoMigrate(
			&User{},
			&Role{},
			&Permission{},
			&UserRole{},
			&RolePermission{},
			&CartItem{},
			&Order{},
			&OrderItem{},
			&Product{},
			&Blacklist{},
			&UserAddress{},
			&ProductReview{},
			&Promotion{},
			&PromotionProduct{},
			&Coupon{},
		)
	}
}
