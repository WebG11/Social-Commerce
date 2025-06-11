package dao

import (
	"fmt"
	"time"

	"github.com/bitdance-panic/gobuy/app/consts"
	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type (
	Order       = models.Order
	UserAddress = models.UserAddress
	Product     = models.Product
)

func CreateOrder(db *gorm.DB, order *Order) error {
	if err := db.Create(order).Error; err != nil {
		return err
	}
	return nil
}

func GetByID(db *gorm.DB, id int) (*Product, error) {
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 处理记录未找到的情况
			return nil, nil
		}
		// 处理其他错误
		return nil, err
	}
	return &product, nil
}

func SaveOrder(db *gorm.DB, order *Order) error {
	if err := db.Save(order).Error; err != nil {
		return err
	}
	return nil
}

func GetOrderByID(db *gorm.DB, orderID int) (*Order, error) {
	var order Order
	result := db.Preload("Items").First(&order, "id = ?", orderID)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &order, nil
}

func UpdateOrderStatus(db *gorm.DB, order *Order, newStatus consts.OrderStatus) error {
	if order == nil {
		return errors.New("order is nil")
	}
	order.Status = int(newStatus)
	return db.Save(order).Error
}

func UpdateOrderAddress(db *gorm.DB, order *Order, newName, newPhone, newAddress string) error {
	if order == nil {
		return errors.New("order is nil")
	}
	order.Address = newAddress
	order.Name = newName
	order.Phone = newPhone
	return db.Save(order).Error
}

func ListUserOrder(db *gorm.DB, userID int, pageNum int, pageSize int) (*[]Order, error) {
	var orders []Order
	err := db.Preload("Items").Where("user_id = ?", userID).Find(&orders).Error
	// .Limit(pageSize).Offset((pageNum-1)*pageSize).
	if err != nil {
		return nil, err
	}
	return &orders, nil
}

func AdminListOrder(db *gorm.DB, pageNum int, pageSize int, sellerID int) (*[]Order, int64, error) {
	var orders []Order
	err := db.Where("seller_id = ?", sellerID).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}
	var count int64
	db.Model(&Order{}).Where("seller_id = ?", sellerID).Count(&count)
	return &orders, count, nil
}

func UpdateOrderDiscount(db *gorm.DB, order *Order, discount float64) error {
	order.Discount = discount
	return db.Save(order).Error
}

func UpdateOrderTracking(db *gorm.DB, order *Order, trackingNumber string) error {
	order.TrackingNumber = trackingNumber
	order.Status = int(consts.OrderStatusShipped)
	return db.Save(order).Error
}

func ListPendingOrder(db *gorm.DB) (*[]Order, error) {
	var orders []Order
	err := db.Where("status = ?", consts.OrderStatusPending).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return &orders, nil
}

// 创建订单地址
func CreateUserAddress(db *gorm.DB, userAddress *UserAddress) error {
	if err := db.Create(userAddress).Error; err != nil {
		return err
	}
	return nil
}

func DeleteUserAddress(db *gorm.DB, userID int32) error {
	result := db.Where("user_id = ?", userID).Delete(&UserAddress{})
	if result.RowsAffected == 0 {
		return fmt.Errorf("no records found with user_id %d", userID)
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 更新订单地址
func UpdateUserAddress(db *gorm.DB, userID int32, userAddress string) error {
	var address *UserAddress
	if err := db.Model(&address).Where("user_id = ?", userID).Update("user_address", userAddress).Error; err != nil {
		return err
	}
	return nil
}

// 获取订单地址
func GetUserAddress(db *gorm.DB, userID int32) ([]models.UserAddress, error) {
	var addresses []models.UserAddress
	if err := db.Where("user_id = ?", userID).Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}

// 根据时间范围返回order数组
func ListOrderByDateRange(db *gorm.DB, startTime, endTime time.Time) (*[]Order, error) {
	var orders []Order

	err := db.Where("created_at >= ? AND created_at <= ?", startTime, endTime).
		Preload("Items").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return &orders, nil
}
