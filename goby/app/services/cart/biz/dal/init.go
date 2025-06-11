package dal

import (
	"github.com/bitdance-panic/gobuy/app/services/cart/biz/dal/tidb"
)

func Init() {
	tidb.Init()
}
