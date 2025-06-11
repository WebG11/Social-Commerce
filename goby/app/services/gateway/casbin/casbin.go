package casbin

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

var Enforcer *casbin.Enforcer

func InitCasbin(db *gorm.DB) error {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return err
	}

	enforcer, err := casbin.NewEnforcer("./casbin/rbac_model.conf", adapter)
	if err != nil {
		return err
	}

	Enforcer = enforcer
	enforcer.EnableAutoSave(true) // 自动保存策略到数据库

	if err := Enforcer.LoadPolicy(); err != nil {
		return err
	}

	return nil
}