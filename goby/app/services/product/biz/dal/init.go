package dal

import (
	"github.com/bitdance-panic/gobuy/app/services/product/biz/dal/tidb"
)

func Init() {
	tidb.Init()
}
