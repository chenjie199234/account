package user

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/chenjie199234/account/ecode"

	"github.com/chenjie199234/Corelib/redis"
	"github.com/chenjie199234/Corelib/util/common"
)

const DefaultExpireSeconds = 300
const DefaultCheckTimes = 5

func init() {
	h := sha1.Sum([]byte(setcode))
	hsetcode = hex.EncodeToString(h[:])
	h = sha1.Sum([]byte(checkcode))
	hcheckcode = hex.EncodeToString(h[:])
}

// argv 1 = code
// argv 2 = max check times
// argv 3 = expire time(unit second)
// return nil = set success
// return <=0 = already setted before and all check times failed,ban some time
// return >0 = already setted before,data is the rest check times
const setcode = `local used=redis.call("HGET",KEYS[1],"check")
if(used)
then
	return ARGV[2]-used
end
redis.call("HMSET",KEYS[1],"code",ARGV[1],"check",0)
redis.call("EXPIRE",KEYS[1],ARGV[3])
return nil`

var hsetcode string

func (d *Dao) RedisSetCode(ctx context.Context, target, action, code string) (int, error) {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return 0, e
	}
	defer c.Close()
	var rest int
	rest, e = redis.Int(c.DoContext(ctx, "EVALSHA", hsetcode, 1, "code_{"+target+"}_"+action, code, DefaultCheckTimes, DefaultExpireSeconds))
	if e != nil && strings.Contains(e.Error(), "NOSCRIPT") {
		rest, e = redis.Int(c.DoContext(ctx, "EVAL", setcode, 1, "code_{"+target+"}_"+action, code, DefaultCheckTimes, DefaultExpireSeconds))
	}
	if e == redis.ErrNil {
		return DefaultCheckTimes, nil
	}
	if e != nil {
		return 0, e
	}
	e = ecode.ErrCodeAlreadySend
	if rest < 0 {
		rest = 0
	}
	return rest, e
}

// argv 1 = code
// argv 2 = max check times
// argv 3 = expire time(unit second)
// return nil = key already expired
// return number = rest check times
// -1 means code same,check success
// 0 means all check times used,this key will be expired after expire time
const checkcode = `local data=redis.call("HMGET",KEYS[1],"code","check")
if(not data[1] or not data[2])
then
	return nil
end
if(data[1]==ARGV[1] and data[2]<ARGV[2])
then
	redis.call("DEL",KEYS[1])
	return -1
end
if(data[2]>=ARGV[2])
then
	return 0
end
if(redis.call("EXPIRE",KEYS[1],ARGV[3])==0)
then
	return nil
end
data[2]=redis.call("HINCRBY",KEYS[1],"check",1)
return ARGV[2]-data[2]`

var hcheckcode string

func (d *Dao) RedisCheckCode(ctx context.Context, target, action, code string) (int, error) {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return 0, e
	}
	defer c.Close()
	var rest int
	rest, e = redis.Int(c.DoContext(ctx, "EVALSHA", hcheckcode, 1, "code_{"+target+"}_"+action, code, DefaultCheckTimes, DefaultExpireSeconds))
	if e != nil && strings.Contains(e.Error(), "NOSCRIPT") {
		rest, e = redis.Int(c.DoContext(ctx, "EVAL", checkcode, 1, "code_{"+target+"}_"+action, code, DefaultCheckTimes, DefaultExpireSeconds))
	}
	if e == redis.ErrNil {
		e = ecode.ErrCodeNotExist
	}
	return rest, e
}
func (d *Dao) RedisGetCode(ctx context.Context, target, action string) (string, int, error) {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return "", 0, e
	}
	defer c.Close()
	values, e := redis.Values(c.DoContext(ctx, "HMGET", "code_{"+target+"}_"+action, "code", "check"))
	if e != nil {
		if e == redis.ErrNil {
			e = ecode.ErrCodeNotExist
		}
		return "", 0, e
	}
	if values[0] == nil || values[1] == nil {
		return "", 0, ecode.ErrCodeNotExist
	}
	code := common.Byte2str(values[0].([]byte))
	check, e := strconv.Atoi(common.Byte2str(values[1].([]byte)))
	if e != nil {
		e = ecode.ErrRedisDataBroken
	}
	return code, DefaultCheckTimes - check, e
}

func (d *Dao) RedisDelCode(ctx context.Context, target, action string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.Int64(c.DoContext(ctx, "DEL", "code_{"+target+"}_"+action))
	return e
}
