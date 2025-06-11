package tidb

import (
	"fmt"
	"log"

	"github.com/bitdance-panic/gobuy/app/services/gateway/conf"

	"gorm.io/driver/mysql"
	"gorm.io/plugin/opentelemetry/tracing"

	"gorm.io/gorm"
	// "gorm.io/plugin/opentelemetry/tracing"
)

var DB *gorm.DB

func Init() {
	conf_ := conf.GetConf()
	// err := driver.RegisterTLSConfig("tidb", &tls.Config{
	// 	MinVersion: tls.VersionTLS12,
	// 	ServerName: conf_.Tidb.Host,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	dsn := fmt.Sprintf(conf_.Tidb.DSN, conf_.Tidb.User, conf_.Tidb.Password, conf_.Tidb.Host, conf_.Tidb.Port)
	// log.Printf("%s", dsn)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}
	log.Println("数据库连接成功")
	if err := DB.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
		panic(err)
	}
	// models.AutoMigrate(DB) // 第一次连接后就关闭
}
