package dal

import (
	"github.com/bitdance-panic/gobuy/app/services/paycallback/biz/dal/tidb"
)

func Init() {
	tidb.Init()
}
