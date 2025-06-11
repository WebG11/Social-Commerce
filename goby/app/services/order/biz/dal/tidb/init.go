package tidb

import (
	"crypto/tls"
	"fmt"

	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/services/order/conf"
	driver "github.com/go-sql-driver/mysql"

	"gorm.io/driver/mysql"
	"gorm.io/plugin/opentelemetry/tracing"

	// "gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "gorm.io/plugin/opentelemetry/tracing"
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
	// log.Printf("%s", dsn)
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

	models.AutoMigrate(DB)
	// 第一次连接后就关闭

	// type TimeZoneInfo struct {
	// 	GlobalTimeZone  string `gorm:"column:global_time_zone"`
	// 	SessionTimeZone string `gorm:"column:session_time_zone"`
	// 	SystemTimeZone  string `gorm:"column:system_time_zone"`
	// 	Noww            string `gorm:"column:now"`
	// }
	// var timeZoneInfo TimeZoneInfo
	// err = DB.Raw("SELECT Now() as now,  @@global.time_zone AS global_time_zone, @@session.time_zone AS session_time_zone, @@global.system_time_zone AS system_time_zone").Scan(&timeZoneInfo).Error
	// if err != nil {
	// 	panic("failed to execute query")
	// }

	// // 打印查询结果
	// fmt.Println(timeZoneInfo.Noww)
	// fmt.Printf("Global Time Zone: %s\n", timeZoneInfo.GlobalTimeZone)
	// fmt.Printf("Session Time Zone: %s\n", timeZoneInfo.SessionTimeZone)
	// fmt.Printf("System Time Zone: %s\n", timeZoneInfo.SystemTimeZone)
}
