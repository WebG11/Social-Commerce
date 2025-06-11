package dal

import (
	"github.com/bitdance-panic/gobuy/app/services/agent/biz/dal/tidb"
)

func Init() {
	tidb.Init()
}
