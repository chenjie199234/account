package money

import (
	"context"
	"log/slog"
	"unsafe"

	"github.com/chenjie199234/account/model"

	"github.com/chenjie199234/Corelib/cotel"
	cmongo "github.com/chenjie199234/Corelib/mongo"
	cmysql "github.com/chenjie199234/Corelib/mysql"
	credis "github.com/chenjie199234/Corelib/redis"
	"github.com/chenjie199234/Corelib/util/oneshot"
	gredis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const DefaultMoneyLogsPageSize = 10

// Dao this is a data operation layer to operate money service's data
type Dao struct {
	mysql *cmysql.Client
	redis *credis.Client
	mongo *cmongo.Client
}

// NewDao Dao is only a data operation layer
// don't write business logic in this package
// business logic should be written in service package
func NewDao(mysql *cmysql.Client, redis *credis.Client, mongo *cmongo.Client) *Dao {
	return &Dao{
		mysql: mysql,
		redis: redis,
		mongo: mongo,
	}
}

func (d *Dao) GetMoneyLogs(ctx context.Context, userid bson.ObjectID, opaction string, starttime, endtime, pagesize, page uint32) ([]*model.MoneyLog, uint32, uint32, error) {
	if moneylogs, totalsize, curpage, e := d.RedisGetMoneyLogs(ctx, userid.Hex(), opaction, starttime, endtime, pagesize, page); e != nil {
		if e != gredis.Nil {
			slog.ErrorContext(ctx, "[dao.GetMoneyLogs] redis op failed", slog.String("user_id", userid.Hex()), slog.String("opaction", opaction), slog.String("error", e.Error()))
		}
	} else {
		return moneylogs, totalsize, curpage, nil
	}
	//redis error or redis not exist,we need to query db
	unsafeAll, e := oneshot.Do("GetMoneyLogs_"+opaction+"_"+userid.Hex(), func() (unsafe.Pointer, error) {
		all, e := d.MongoGetMoneyLogs(ctx, userid, opaction)
		if e != nil {
			slog.ErrorContext(nil, "[dao.GetMoneyLogs] db op failed", slog.String("user_id", userid.Hex()), slog.String("opaction", opaction), slog.String("error", e.Error()))
			return nil, e
		}
		//update redis
		go func() {
			ctx := cotel.CloneTrace(ctx)
			if e := d.RedisSetMoneyLogs(ctx, userid.Hex(), opaction, all); e != nil {
				slog.ErrorContext(ctx, "[dao.GetMoneyLogs] update redis failed", slog.String("user_id", userid.Hex()), slog.String("opaction", opaction), slog.String("error", e.Error()))
			}
		}()
		return unsafe.Pointer(&all), nil
	})
	if e != nil {
		return nil, 0, 0, e
	}
	//all is sorted by DESC by db's query
	all := *(*[]*model.MoneyLog)(unsafeAll)
	start := 0
	startfind := false
	end := len(all)
	endfind := false
	for i, moneylog := range all {
		if moneylog.LogID.Timestamp().Unix() <= int64(endtime) && !startfind {
			start = i
			startfind = true
		}
		if moneylog.LogID.Timestamp().Unix() < int64(starttime) && !endfind {
			end = i
			endfind = true
		}
		if startfind && endfind {
			break
		}
	}
	if !startfind && !endfind {
		return nil, 0, 0, nil
	}
	need := all[start:end]
	totalsize := uint32(len(need))
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
