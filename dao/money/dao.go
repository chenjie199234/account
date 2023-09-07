package money

import (
	"context"
	csql "database/sql"
	"unsafe"

	"github.com/chenjie199234/Corelib/log"
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

type getLogsResult struct {
	Logs      []*model.MoneyLog
	Totalsize uint32
	CurPage   uint32
}

func (d *Dao) GetMoneyLogs(ctx context.Context, userid primitive.ObjectID, opaction string, starttime, endtime, pagesize, page uint32) ([]*model.MoneyLog, uint32, uint32, error) {
	if moneylogs, totalsize, curpage, e := d.RedisGetMoneyLogs(ctx, userid.Hex(), opaction, starttime, endtime, pagesize, page); e != nil {
		log.Error(nil, "[dao.GetMoneyLogs] redis op failed", map[string]interface{}{"user_id": userid.Hex(), "error": e})
	} else {
		return moneylogs, totalsize, curpage, nil
	}
	//redis error or redis not exist,we need to query db
	unsafeAll, e := oneshot.Do("GetMoneyLogs_"+opaction+"_"+userid.Hex(), func() (unsafe.Pointer, error) {
		all, e := d.MongoGetMoneyLogs(ctx, userid, opaction)
		if e != nil {
			return nil, e
		}
		//update redis
		go func() {
			if e := d.RedisSetMoneyLogs(context.Background(), userid.Hex(), opaction, all); e != nil {
				log.Error(nil, "[dao.GetMoneyLogs] update redis failed", map[string]interface{}{"user_id": userid.Hex(), "error": e})
			}
		}()
		return unsafe.Pointer(&all), nil
	})
	if e != nil {
		return nil, 0, 0, e
	}
	all := *(*[]*model.MoneyLog)(unsafeAll)
	need := make([]*model.MoneyLog, 0, len(all))
	for _, moneylog := range all {
		logtime := moneylog.LogID.Timestamp().Unix()
		if logtime >= int64(starttime) && logtime <= int64(endtime) {
			need = append(need, moneylog)
		}
	}
	totalsize := uint32(len(need))
	if totalsize == 0 {
		return need, 0, 0, nil
	}
	if page == 0 {
		return need, totalsize, 0, nil
	}
	curpage := page
	skip := (curpage - 1) * pagesize
	if skip >= totalsize {
		if totalsize%pagesize > 0 {
			curpage = totalsize/pagesize + 1
		} else {
			curpage = totalsize / pagesize
		}
		skip = (curpage - 1) * pagesize
	}
	if skip+pagesize > totalsize {
		return need[skip:], totalsize, curpage, nil
	}
	return need[skip : skip+pagesize], totalsize, curpage, nil
}
