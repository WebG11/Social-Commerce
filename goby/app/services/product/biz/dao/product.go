package dao

import (
	"errors"

	"gorm.io/gorm"

	"github.com/bitdance-panic/gobuy/app/models"
)

type Product = models.Product

func Create(db *gorm.DB, product *Product) error {
	return db.Create(product).Error
}

func CreateReview(db *gorm.DB, review *models.ProductReview) error {
	return db.Create(review).Error
}

func GetReviews(db *gorm.DB, productID int) (*[]models.ProductReview, error) {
	var reviews []models.ProductReview
	if err := db.Where("product_id = ?", productID).Find(&reviews).Error; err != nil {
		return nil, err
	}
	return &reviews, nil
}

// ListAll 获取所有商品
func ListAll(db *gorm.DB) ([]models.Product, error) {
	var products []models.Product
	err := db.Find(&products).Error
	return products, err
}

func List(db *gorm.DB, pageNum int, pageSize int) (*[]Product, int64, error) {
	var products []Product
	if err := db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Where("is_deleted = false").Find(&products).Error; err != nil {
		return nil, 0, err
	}
	var count int64
	db.Model(&Product{}).Where("is_deleted = false").Count(&count)
	return &products, count, nil
}

func AdminList(db *gorm.DB, pageNum int, pageSize int, sellerID int) (*[]Product, int64, error) {
	var products []Product
	if err := db.Where("seller_id = ?", sellerID).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}
	var count int64
	db.Model(&Product{}).Where("seller_id = ?", sellerID).Count(&count)
	return &products, count, nil
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

func Update(db *gorm.DB, product *Product) error {
	if product == nil {
		return errors.New("product is nil")
	}
	return db.Save(product).Error
}

func Remove(db *gorm.DB, id int) error {
	result := db.Model(&Product{}).Where("id = ? AND is_deleted = false", id).Update("is_deleted", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no product found with the given ID")
	}
	return nil
}

func Search(db *gorm.DB, query string, category string, pageNum int, pageSize int, minPrice int, maxPrice int, brand string) ([]Product, int64, error) {
	var products []Product
	searchQuery := "%" + query + "%"

	queryBuilder := db.Limit(pageSize).Offset((pageNum-1)*pageSize).
		Where("is_deleted = ? AND (name LIKE ? OR description LIKE ?)", false, searchQuery, searchQuery)

	if category != "" {
		queryBuilder = queryBuilder.Where("category = ?", category)
	}

	queryBuilder = queryBuilder.Where("price BETWEEN ? AND ?", minPrice, maxPrice)

	if brand != "" {
		brandQuery := "%" + brand + "%"
		queryBuilder = queryBuilder.Where("name LIKE ?", brandQuery)
	}

	if err := queryBuilder.Find(&products).Error; err != nil {
		return nil, 0, err
	}

	countQuery := db.Model(&Product{}).Where("is_deleted = ? AND (name LIKE ? OR description LIKE ?)", false, searchQuery, searchQuery)
	if category != "" {
		countQuery = countQuery.Where("category = ?", category)
	}

	countQuery = countQuery.Where("price BETWEEN ? AND ?", minPrice, maxPrice)

	if brand != "" {
		brandQuery := "%" + brand + "%"
		countQuery = countQuery.Where("name LIKE ?", brandQuery)
	}
	var count int64
	countQuery.Count(&count)

	return products, count, nil
}

func CreatePromotion(db *gorm.DB, promotion *models.Promotion) error {
	return db.Create(promotion).Error
}

func GetActivePromotions(db *gorm.DB) ([]models.Promotion, error) {
	var promotions []models.Promotion
	err := db.Find(&promotions).Error
	return promotions, err
}

func DeletePromotion(db *gorm.DB, id int) error {
	return db.Delete(&models.Promotion{}, id).Error
}

func GetProductsByStockRange(db *gorm.DB, minStock *int, maxStock *int) ([]Product, error) {
	var products []Product

	// 构建查询条件
	query := db.Where("is_deleted = ?", false)

	if minStock != nil {
		query = query.Where("stock >= ?", *minStock)
	}

	if maxStock != nil {
		query = query.Where("stock <= ?", *maxStock)
	}

	// 查询符合条件的商品
	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}
