package dao

import (
	"context"
	"time"

	"github.com/bitdance-panic/gobuy/app/models"

	"gorm.io/gorm"
)

type User = models.User

func RegisterUser(db *gorm.DB, ctx context.Context, username, password, email string, isSeller bool) (*User, error) {
	hashedPassword := password

	user := &User{
		Username:       username,
		PasswordHashed: hashedPassword,
		Email:          email,
		IsSeller:       isSeller,
	}

	// 插入新用户
	err := db.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func Seller(db *gorm.DB, ctx context.Context, userID int32) (bool, error) {
	var user User
	err := db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		return false, err
	}
	return user.IsSeller, nil
}

// 查询所有用户信息，并进行分页
func AdminListUser(db *gorm.DB, ctx context.Context, page int, pageSize int) ([]User, int64, error) {
	var users []User
	offset := (page - 1) * pageSize
	err := db.WithContext(ctx).Limit(pageSize).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	var count int64
	db.Model(&User{}).Count(&count)
	return users, count, nil
}

func GetUserByEmailAndPass(db *gorm.DB, ctx context.Context, email string, password string) (*User, error) {
	var userPO User
	err := db.WithContext(ctx).Model(&User{}).Where(&User{Email: email, PasswordHashed: password}).First(&userPO).Error
	return &userPO, err
}

func CreateUser(db *gorm.DB, ctx context.Context, user *User) error {
	return db.WithContext(ctx).Create(user).Error
}

// GetUserByID 根据用户 ID 查询用户信息
func GetUserByID(db *gorm.DB, ctx context.Context, userID int) (*User, error) {
	user := &User{}
	err := db.WithContext(ctx).
		Where("id = ?", userID).
		First(&user).
		Error
	return user, err
}

// UpdateUserAddressByID 更新用户地址信息
func UpdateUserAddressByID(db *gorm.DB, ctx context.Context, userID int, userName, phone, userAddress string) error {
	return db.Model(&models.UserAddress{}).Where("user_id = ?", userID).Updates(models.UserAddress{
		UserName:    userName,
		Phone:       phone,
		UserAddress: userAddress,
	}).Error
}

// 更新
func UpdateUserByID(db *gorm.DB, ctx context.Context, userID int, username, email string, passwordHashed string) error {
	// if db.DB == nil {
	// 	return errors.New("database connection is nil")
	// }
	return db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"username":        username,
			"email":           email,
			"password_hashed": passwordHashed,
		}).Error
}

func GetUserAddressesByUserID(db *gorm.DB, ctx context.Context, userID int) ([]*models.UserAddress, error) {
	var addresses []*models.UserAddress
	err := db.WithContext(ctx).Where("user_id = ?", userID).Find(&addresses).Error
	return addresses, err
}

func UpdateUserAddress(db *gorm.DB, ctx context.Context, userID int, username, email string, passwordHashed string) error {
	// if db.DB == nil {
	// 	return errors.New("database connection is nil")
	// }
	return db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"username":        username,
			"email":           email,
			"password_hashed": passwordHashed,
		}).Error
}

func DeleteUserAddressesByUserID(db *gorm.DB, ctx context.Context, userID int) error {
	return db.Where("user_id = ?", userID).Delete(&models.UserAddress{}).Error
}

// InsertUserAddress 插入新的用户地址
func InsertUserAddress(db *gorm.DB, ctx context.Context, userID int, userName, userAddress, phone string) error {
	return db.Create(&models.UserAddress{
		UserID:      userID,
		UserName:    userName,
		UserAddress: userAddress,
		Phone:       phone,
	}).Error
}

// 删除
func DeleteUserByID(db *gorm.DB, ctx context.Context, userID int) error {
	return db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", userID).
		Update("is_deleted", true).Error
}

func BlockUser(db *gorm.DB, ctx context.Context, identifier string, reason string, expires_at int64) (*models.Blacklist, error) {
	entry := &models.Blacklist{
		Identifier: identifier,
		Reason:     reason,
		ExpiresAt:  time.Unix(expires_at, 0),
	}

	if err := db.WithContext(ctx).Create(entry).Error; err != nil {
		return nil, err
	}

	return entry, nil
}

func UnblockUser(db *gorm.DB, ctx context.Context, identifier string) error {
	return db.WithContext(ctx).
		Model(&models.Blacklist{}).
		Where("identifier = ?", identifier).
		Update("is_deleted", true).Error

	// 删除记录
	// if err := db.Where("identifier = ?", identifier).Delete(&models.Blacklist{}).Error; err != nil {
	// 	return err
	// }
	// return nil
}
