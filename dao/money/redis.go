package money

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/chenjie199234/account/ecode"
	"github.com/chenjie199234/account/model"

	"github.com/chenjie199234/Corelib/redis"
)

func init() {
	h := sha1.Sum([]byte(setMoneyLogs))
	hSetMoneyLogs = hex.EncodeToString(h[:])

	h = sha1.Sum([]byte(appendMoneyLogs))
	hAppendMoneyLogs = hex.EncodeToString(h[:])

	h = sha1.Sum([]byte(getMoneyLogs))
	hGetMoneyLogs = hex.EncodeToString(h[:])
}

// argv 1 = expire time(uint second)
// argv rest = data
const setMoneyLogs = `redis.call("DEL",KEYS[1])
redis.call("ZADD",KEYS[1],0,"")
for i=2,#ARGV,2 do
	redis.call("ZADD",KEYS[1],ARGV[i],ARGV[i+1])
end
redis.call("EXPIRE",KEYS[1],ARGV[1])
return "OK"`

var hSetMoneyLogs = ""

func (d *Dao) RedisSetMoneyLogs(ctx context.Context, userid, opaction string, logs []*model.MoneyLog) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	args := make([]interface{}, 0, len(logs)*2+4)
	args = append(args, hSetMoneyLogs, 1, opaction+"_money_logs_{"+userid+"}", 604800)
	for _, log := range logs {
		if log == nil || log.LogID.IsZero() {
			continue
		}
		data, _ := json.Marshal(log)
		args = append(args, log.LogID.Timestamp().Unix(), data)
	}
	_, e = redis.String(c.DoContext(ctx, "EVALSHA", args...))
	if e != nil && strings.Contains(e.Error(), "NOSCRIPT") {
		args[0] = setMoneyLogs
		_, e = redis.String(c.DoContext(ctx, "EVAL", args...))
	}
	return e
}

// argv 1 = expire time(unit second)
// argv 2 = score
// argc 3 = data
const appendMoneyLogs = `if(redis.call("EXPIRE",KEYS[1],ARGV[1])==0)
then
	return nil
end
redis.call("ZADD",KEYS[1],ARGV[2],ARGV[3])
return "OK"`

var hAppendMoneyLogs = ""

// log will be appended only when the user's money logs' redis key exist.if key not exist,ecode.ErrMoneyLogsNotExist will return
func (d *Dao) RedisAppendMoneyLogs(ctx context.Context, userid, opaction string, log *model.MoneyLog) error {
	if log == nil || log.LogID.IsZero() {
		return nil
	}
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	data, _ := json.Marshal(log)
	_, e = redis.String(c.DoContext(ctx, "EVALSHA", hAppendMoneyLogs, 1, opaction+"_money_logs_{"+userid+"}", 604800, log.LogID.Timestamp().Unix(), data))
	if e != nil && strings.Contains(e.Error(), "NOSCRIPT") {
		_, e = redis.String(c.DoContext(ctx, "EVAL", appendMoneyLogs, 1, opaction+"_money_logs_{"+userid+"}", 604800, log.LogID.Timestamp().Unix(), data))
	}
	if e != nil {
		if e == redis.ErrNil {
			e = ecode.ErrMoneyLogsNotExist
		}
	}
	return e
}
func (d *Dao) RedisDelMoneyLogs(ctx context.Context, userid, opaction string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = c.DoContext(ctx, "DEL", opaction+"_money_logs_{"+userid+"}")
	return e
}

// argv 1 = starttime
// argv 2 = endtime
// argv 3 = pagesize
// argv 4 = page
const getMoneyLogs = `if(redis.call("EXISTS",KEYS[1])==0)
then
	return nil
end`

var hGetMoneyLogs = ""

// if page == 0,return all money logs
// if page != 0,return the required page's logs
// if page != 0 and overflow,return the last page's logs
func (d *Dao) RedisGetMoneyLogs(ctx context.Context, userid, action string, starttime, endtime, pagesize, page int64) ([]*model.MoneyLog, int64, int64, error) {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return nil, 0, 0, e
	}
	defer c.Close()
	var totalsize int64
	var curpage int64
	var datas [][]byte

	if page == 0 {
		// if datas, e = redis.ByteSlices(c.DoContext(ctx, "ZRANGE", action+"_money_logs_{"+userid+"}", endtime, starttime, "BYSCORE", "REV")); e != nil {
		// 	return nil, 0, 0, e
		// }
		// if len(datas) == 0 {
		// 	return nil, 0, 0, ecode.ErrMoneyLogsNotExist
		// }
		// totalsize = int64(len(datas)) - 1 //delete the score 0(empty key,placeholder)
		// curpage = 0
	} else {
		// totalsize, e = redis.Int64(c.DoContext(ctx, "ZCOUNT", action+"_money_logs_{"+userid+"}", starttime, endtime))
		// if e != nil {
		// 	return nil, 0, 0, e
		// }
		// if totalsize == 0 {
		// 	return nil, 0, 0, nil
		// }
		// curpage = page
		// skip := (curpage - 1) * pagesize
		// if skip >= totalsize {
		// 	if totalsize%pagesize > 0 {
		// 		curpage = totalsize/page + 1
		// 	} else {
		// 		curpage = totalsize / page
		// 	}
		// 	skip = (curpage - 1) * pagesize
		// }
		// if datas, e = redis.ByteSlices(c.DoContext(ctx, "LRANGE", action+"_money_logs_{"+userid+"}", skip, skip+pagesize)); e != nil {
		// 	return nil, 0, 0, e
		// }
	}
}
