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
redis.call("ZADD",KEYS[1],-1,"")
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
	if e != nil && e == redis.ErrNil {
		e = ecode.ErrRedisKeyMissing
	}
	return e
}
func (d *Dao) RedisDelMoneyLogs(ctx context.Context, userid, opaction string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.Int64(c.DoContext(ctx, "DEL", opaction+"_money_logs_{"+userid+"}"))
	return e
}

// argv 1 = starttime(unit second)
// argv 2 = endtime(unit second)
// argv 3 = page //if page == 0,return all data //if page != 0,return the required page's data //if page != 0 and page overflow,return the last page's data
// argv 4 = pagesize
// return data is a list
// return data's last element = curpage
// return data's last second element = totalsize
// return data's rest element is the data
const getMoneyLogs = `if(redis.call("EXISTS",KEYS[1])==0)
then
	return nil
end
local result={}
local totalsize=0
local curpage=0
if(tonumber(ARGV[3])==0)
then
	local tmp=redis.call("ZRANGE",KEYS[1],ARGV[2],ARGV[1],"BYSCORE","REV")
	if(tmp and #tmp>0)
	then
		result=tmp
		totalsize=#result
	end
else
	totalsize=redis.call("ZCOUNT",KEYS[1],ARGV[1],ARGV[2])
	if(totalsize~=0)
	then
		curpage=tonumber(ARGV[3])
		local skip=(curpage-1)*ARGV[4]
		if(skip>=totalsize)
		then
			if(totalsize%ARGV[4]>0)
			then
				curpage=(totalsize-totalsize%ARGV[4])/ARGV[4]+1
			else
				curpage=totalsize/ARGV[4]
			end
			skip=(curpage-1)*ARGV[4]
		end
		result=redis.call("ZRANGE",KEYS[1],ARGV[2],ARGV[1],"BYSCORE","REV","LIMIT",skip,ARGV[4])
	end
end
result[#result+1]=totalsize
result[#result+1]=curpage
return result`

var hGetMoneyLogs = ""

// if page == 0,return all money logs
// if page != 0,return the required page's logs
// if page != 0 and page overflow,return the last page's logs
func (d *Dao) RedisGetMoneyLogs(ctx context.Context, userid, opaction string, starttime, endtime, pagesize, page uint32) ([]*model.MoneyLog, uint32, uint32, error) {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return nil, 0, 0, e
	}
	defer c.Close()
	values, e := redis.Values(c.DoContext(ctx, "EVALSHA", hGetMoneyLogs, 1, opaction+"_money_logs_{"+userid+"}", starttime, endtime, page, pagesize))
	if e != nil && strings.Contains(e.Error(), "NOSCRIPT") {
		values, e = redis.Values(c.DoContext(ctx, "EVAL", getMoneyLogs, 1, opaction+"_money_logs_{"+userid+"}", starttime, endtime, page, pagesize))
	}
	if e != nil {
		if e == redis.ErrNil {
			e = ecode.ErrRedisKeyMissing
		}
		return nil, 0, 0, e
	}
	curpage := values[len(values)-1].(int64)
	totalsize := values[len(values)-2].(int64)
	values = values[:len(values)-2]
	r := make([]*model.MoneyLog, 0, len(values))
	for i := range values {
		tmp := &model.MoneyLog{}
		if e := json.Unmarshal(values[i].([]byte), tmp); e != nil {
			return nil, 0, 0, ecode.ErrRedisDataBroken
		}
		r = append(r, tmp)
	}
	return r, uint32(totalsize), uint32(curpage), nil
}
