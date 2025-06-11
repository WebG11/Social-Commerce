package tidb

import (
	"crypto/tls"
	"fmt"

	"github.com/bitdance-panic/gobuy/app/services/paycallback/conf"
	driver "github.com/go-sql-driver/mysql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

var DB *gorm.DB

// Init initializes the database connection for the payment service
func Init() {
	// Retrieve the configuration for the payment service
	conf_ := conf.GetConf()

	// Register TLS configuration for secure connection
	err := driver.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: conf_.Tidb.Host,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to register TLS config: %v", err))
	}

	// Construct the Data Source Name (DSN) for connecting to the database
	dsn := fmt.Sprintf(conf_.Tidb.DSN, conf_.Tidb.User, conf_.Tidb.Password, conf_.Tidb.Host, conf_.Tidb.Port)

	// Open a connection to the database using GORM
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	if err := DB.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
		panic(err)
	}

	// Models.AutoMigrate(DB) // Uncomment this line to automatically migrate models if needed
}
