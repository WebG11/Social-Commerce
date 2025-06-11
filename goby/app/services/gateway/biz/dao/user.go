package dao

import (
	"errors"
	"fmt"
	"time"

	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/services/gateway/casbin"
	"gorm.io/gorm"
)

type User = models.User

// UpdateRefreshToken 更新刷新令牌
func UpdateRefreshToken(db *gorm.DB, userID int, refreshToken string) error {
	return db.Model(&User{}).Where("id = ?", userID).
		Update("refresh_token", refreshToken).
		Error
}

// 根据用户ID查询 Refresh Token
func GetRefreshTokenByUserID(db *gorm.DB, userID int) (string, error) {
	var user User
	result := db.Model(&User{}).Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return "", errors.New("refresh token not found")
	}
	return user.RefreshToken, nil
}

// 添加用户角色并同步Casbin策略
func AddUserRole(db *gorm.DB, userID int, roleID int) error {
	// 查询角色信息
	var role models.Role
	if err := db.First(&role, roleID).Error; err != nil {
		return err
	}

	// 添加用户角色关联
	userRole := models.UserRole{UserID: userID, RoleID: int(role.Base.ID)}
	if err := db.Create(&userRole).Error; err != nil {
		return err
	}

	// 同步到Casbin策略
	_, err := casbin.Enforcer.AddGroupingPolicy(fmt.Sprintf("%d", userID), role.Name)
	return err
}

// 删除用户角色并同步Casbin策略
func DeleteUserRole(db *gorm.DB, userID int, roleID int) error {
	// 查询角色信息
	var role models.Role
	if err := db.First(&role, roleID).Error; err != nil {
		return err
	}

	// 删除用户角色关联
	if err := db.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&models.UserRole{}).Error; err != nil {
		return err
	}

	// 同步到Casbin策略
	_, err := casbin.Enforcer.RemoveGroupingPolicy(fmt.Sprintf("%d", userID), role.Name)
	return err
}

func ClearExpBlacklist(db *gorm.DB) error {
	// 删除过期条目
	return db.Where("expires_at < ?", time.Now()).Delete(&models.Blacklist{}).Error
}

func GetUserName(db *gorm.DB, userID int) (string, error) {
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return "", err
	}
	return user.Username, nil
}
