package tidb

import (
	"crypto/tls"
	"fmt"

	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/services/product/conf"
	driver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/plugin/opentelemetry/tracing"

	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	conf_ := conf.GetConf()
	err := driver.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: conf_.Tidb.Host,
	})
	if err != nil {
		panic(err)
	}
	dsn := fmt.Sprintf(conf_.Tidb.DSN, conf_.Tidb.User, conf_.Tidb.Password, conf_.Tidb.Host, conf_.Tidb.Port)
	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}

	if err := DB.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
		panic(err)
	}
	models.AutoMigrate(DB) // 第一次连接后就关闭
}
