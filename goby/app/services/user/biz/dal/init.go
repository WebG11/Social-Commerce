package dal

import (
	"github.com/bitdance-panic/gobuy/app/services/user/biz/dal/tidb"
)

func Init() {
	tidb.Init()
}
