package dal

import (
	"github.com/bitdance-panic/gobuy/app/services/order/biz/dal/tidb"
)

func Init() {
	tidb.Init()
}
