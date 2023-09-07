package money

import (
	"context"
	csql "database/sql"
	"unsafe"

	"github.com/chenjie199234/account/ecode"

	credis "github.com/chenjie199234/Corelib/redis"
	"github.com/chenjie199234/Corelib/util/oneshot"
	"github.com/chenjie199234/account/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	cmongo "go.mongodb.org/mongo-driver/mongo"
)

const DefaultMoneyLogsPageSize = 20

// Dao this is a data operation layer to operate money service's data
type Dao struct {
	sql   *csql.DB
	redis *credis.Pool
	mongo *cmongo.Client
}

// NewDao Dao is only a data operation layer
// don't write business logic in this package
// business logic should be written in service package
func NewDao(sql *csql.DB, redis *credis.Pool, mongo *cmongo.Client) *Dao {
	return &Dao{
		sql:   sql,
		redis: redis,
		mongo: mongo,
	}
}

type logsResult struct {
	Logs      []*model.MoneyLog
	Totalsize uint32
	CurPage   uint32
}

func (d *Dao) GetMoneyLogs(ctx context.Context, callerName string, userid primitive.ObjectID, opaction string, starttime, endtime, pagesize, page uint32) ([]*model.MoneyLog, uint32, uint32, error) {
	result, e := oneshot.Do("GetMoneyLogs_"+userid.Hex(), func() (unsafe.Pointer, error) {
		logs, totalsize, curpage, e := d.RedisGetMoneyLogs(ctx, userid.Hex(), opaction, starttime, endtime, pagesize, page)
		if e != nil && e != ecode.ErrRedisKeyMissing {
		}
		return unsafe.Pointer(&logsResult{
			Logs:      logs,
			Totalsize: totalsize,
			CurPage:   curpage,
		}), nil
	})
}
