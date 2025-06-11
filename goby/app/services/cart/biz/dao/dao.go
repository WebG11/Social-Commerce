package dao

// 数据访问层

import (
	"errors"

	"github.com/bitdance-panic/gobuy/app/models"
	"gorm.io/gorm"

	"github.com/bitdance-panic/gobuy/app/consts"
)

type CartItem = models.CartItem

func ListItemsByUserID(db *gorm.DB, userID int) (*[]CartItem, error) {
	var items []CartItem
	err := db.Preload("Product").Where("user_id =?", userID).Find(&items).Error
	// .Limit(pageSize).Offset((pageNum-1)*pageSize)
	if err != nil {
		return nil, err
	}
	return &items, nil
}

func Delete(db *gorm.DB, itemID int) error {
	result := db.Delete(&CartItem{}, itemID)
	if result.RowsAffected == 0 {
		return errors.New("cartItem not found")
	}
	return result.Error
}

func GetItemByID(db *gorm.DB, itemID int) (*CartItem, error) {
	var item CartItem
	if err := db.Preload("Product").First(&item, itemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("cartItem not found")
		}
		return nil, err
	}
	return &item, nil
}

func GetItemByUserID(db *gorm.DB, userID int, productID int) (*CartItem, error) {
	var item CartItem
	if err := db.Preload("Product").Where("user_id = ? AND product_id = ? ", userID, productID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没找到没关系
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func UpdateQuantity(db *gorm.DB, item *CartItem, newQuantity int) error {
	if item == nil {
		return errors.New("item is nil")
	}
	item.Quantity = newQuantity
	return db.Save(&item).Error
}

func Create(db *gorm.DB, userID int, productID int) (bool, error) {
	// 尝试找到现有的购物车项
	item, err := GetItemByUserID(db, int(userID), int(productID))
	if err != nil {
		return false, err
	}
	// 已存在，增加数量
	if item != nil {
		err = UpdateQuantity(db, item, item.Quantity+1)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	// 如果没有找到记录，则创建一个新的购物车项
	item = &CartItem{UserID: userID, ProductID: productID, Quantity: consts.ITEM_INITIAL_QUANTITY}
	if err := db.Create(&item).Error; err != nil {
		return false, err
	}
	return true, nil
}
