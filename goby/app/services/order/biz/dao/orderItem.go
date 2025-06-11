package dao

import (
	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type OrderItem = models.OrderItem

func CreateOrderItem(db *gorm.DB, item *OrderItem) (*OrderItem, error) {
	if item == nil {
		return nil, errors.New("item is nil")
	}
	if err := db.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}
