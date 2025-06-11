package casbin

import (
	"log"
	"strconv"

	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

// InitRBACData 初始化基础RBAC数据（角色、权限、关联关系）
func InitRBACData(db *gorm.DB, enforcer *casbin.Enforcer) error {
	// 使用事务确保数据一致性
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// ----------------------------
	// 初始化权限
	// ----------------------------
	permissions := []models.Permission{
		{
			Name:   "all",
			Path:   "/*",
			Method: "(GET|POST|PUT|DELETE)",
		},
		{
			Name:   "product_view",
			Path:   "/product/*",
			Method: "GET",
		},
		{
			Name:   "product_edit",
			Path:   "/product/*",
			Method: "(GET|POST|PUT|DELETE)",
		},
		{
			Name:   "payment_view",
			Path:   "/payment/*",
			Method: "GET",
		},
		{
			Name:   "payment_manage",
			Path:   "/payment/*",
			Method: "(POST|PUT|DELETE)",
		},
	}

	for _, perm := range permissions {
		if err := tx.Where(models.Permission{Name: perm.Name}).FirstOrCreate(&perm).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// ----------------------------
	// 初始化角色
	// ----------------------------
	roles := []models.Role{
		{Name: "admin"},
		{Name: "user"},
	}

	for _, role := range roles {
		if err := tx.Where(models.Role{Name: role.Name}).FirstOrCreate(&role).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// ----------------------------
	// 关联角色与权限
	// ----------------------------
	// 管理员拥有所有权限
	// 查询所有权限
	var allPermissions []models.Permission
	if err := tx.Find(&allPermissions).Error; err != nil {
		tx.Rollback()
		return err
	}

	var adminRole models.Role
	tx.Where("name = ?", "admin").First(&adminRole)
	// 将所有权限添加到角色中
	if err := tx.Model(&adminRole).Association("Permissions").Append(allPermissions); err != nil {
		tx.Rollback()
		return err
	}

	// 同步到Casbin策略
	adminIDStr := strconv.Itoa(int(adminRole.ID))
	for _, perm := range allPermissions {
		if _, err := enforcer.AddPolicy(adminIDStr, perm.Path, perm.Method); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 普通用户只有查看权限
	var userRole models.Role
	tx.Where("name = ?", "user").First(&userRole)
	var viewPerm models.Permission
	tx.Where("name = ?", "product_view").First(&viewPerm)
	if err := tx.Model(&userRole).Association("Permissions").Append(&viewPerm); err != nil {
		tx.Rollback()
		return err
	}
	userIDStr := strconv.Itoa(int(userRole.ID))
	if _, err := enforcer.AddPolicy(userIDStr, viewPerm.Path, viewPerm.Method); err != nil {
		tx.Rollback()
		return err
	}

	// ----------------------------
	// 提交事务
	// ----------------------------
	if err := tx.Commit().Error; err != nil {
		return err
	}

	// 重新加载Casbin策略
	if err := enforcer.LoadPolicy(); err != nil {
		log.Printf("Failed to reload Casbin policy: %v", err)
		return err
	}

	log.Println("RBAC data initialized successfully")
	return nil
}
