package money

import (
	"context"
	"encoding/json"
	"time"

	"github.com/chenjie199234/account/ecode"
	"github.com/chenjie199234/account/model"

	"github.com/chenjie199234/Corelib/util/common"
	gredis "github.com/redis/go-redis/v9"
)

var setMoneyLogsScript *gredis.Script
var addMoneyLogsScript *gredis.Script
var getMoneyLogsScript *gredis.Script

func init() {
	// argv 1 = expire time(uint second)
	// argv rest = data
	setMoneyLogsScript = gredis.NewScript(`redis.call("DEL",KEYS[1])
redis.call("ZADD",KEYS[1],-1,"")
if(#ARGV>1)
then
	redis.call("ZADD",KEYS[1],unpack(ARGV,2))
end
redis.call("EXPIRE",KEYS[1],ARGV[1])
return "OK"`)

	// argv 1 = expire time(unit second)
	// argv 2 = score
	// argc 3 = data
	addMoneyLogsScript = gredis.NewScript(`if(redis.call("EXPIRE",KEYS[1],ARGV[1])==0)
then
	return nil
end
redis.call("ZADD",KEYS[1],ARGV[2],ARGV[3])
return "OK"`)

	// argv 1 = starttime(unit second)
	// argv 2 = endtime(unit second)
	// argv 3 = page //if page == 0,return all data //if page != 0,return the required page's data //if page != 0 and page overflow,return the last page's data
	// argv 4 = pagesize
	// return data is a list
	// return data's last element = curpage
	// return data's last second element = totalsize
	// return data's rest element is the data
	getMoneyLogsScript = gredis.NewScript(`if(redis.call("EXISTS",KEYS[1])==0)
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
return result`)
}

func (d *Dao) RedisSetMoneyLogs(ctx context.Context, userid, opaction string, logs []*model.MoneyLog) error {
	args := make([]interface{}, 0, len(logs)*2+1)
	args = append(args, int64(30*24*time.Hour.Seconds())) //first arg is the expire time
	for _, log := range logs {
		if log == nil || log.LogID.IsZero() {
			continue
		}
		data, _ := json.Marshal(log)
		args = append(args, log.LogID.Timestamp().Unix(), data)
	}
	_, e := setMoneyLogsScript.Run(ctx, d.redis, []string{opaction + "_money_logs_{" + userid + "}"}, args...).Result()
	return e
}

// log will be added only when the user's money logs' redis key exist.if key not exist,ecode.ErrMoneyLogsNotExist will return
func (d *Dao) RedisAddMoneyLogs(ctx context.Context, userid, opaction string, log *model.MoneyLog) error {
	data, _ := json.Marshal(log)
	_, e := addMoneyLogsScript.Run(ctx, d.redis, []string{opaction + "_money_logs_{" + userid + "}"}, int64(30*24*time.Hour.Seconds()), log.LogID.Timestamp().Unix(), data).Result()
	return e
}
func (d *Dao) RedisDelMoneyLogs(ctx context.Context, userid, opaction string) error {
	_, e := d.redis.Del(ctx, opaction+"_money_logs_{"+userid+"}").Result()
	return e
}

// if page == 0,return all money logs
// if page != 0,return the required page's logs
// if page != 0 and page overflow,return the last page's logs
func (d *Dao) RedisGetMoneyLogs(ctx context.Context, userid, opaction string, starttime, endtime, pagesize, page uint32) ([]*model.MoneyLog, uint32, uint32, error) {
	values, e := getMoneyLogsScript.Run(ctx, d.redis, []string{opaction + "_money_logs_{" + userid + "}"}, starttime, endtime, page, pagesize).Slice()
	if e != nil {
		return nil, 0, 0, e
	}
	curpage := values[len(values)-1].(int64)
	totalsize := values[len(values)-2].(int64)
	values = values[:len(values)-2]
	r := make([]*model.MoneyLog, 0, len(values))
	for i := range values {
		tmp := &model.MoneyLog{}
		if e := json.Unmarshal(common.STB(values[i].(string)), tmp); e != nil {
			return nil, 0, 0, ecode.ErrCacheDataBroken
		}
		r = append(r, tmp)
	}
	return r, uint32(totalsize), uint32(curpage), nil
}
