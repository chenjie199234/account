package user

import (
	"context"
	"strconv"

	"github.com/chenjie199234/account/ecode"

	gredis "github.com/redis/go-redis/v9"
)

const DefaultExpireSeconds = 300
const DefaultCheckTimes = 5

var setCodeScript *gredis.Script
var checkCodeScript *gredis.Script

func init() {
	// argv 1 = code
	// argv 2 = max check times
	// argv 3 = expire time(unit second)
	// return nil = set success
	// return <=0 = already setted before and all check times failed,ban some time
	// return >0 = already setted before,data is the rest check times
	setCodeScript = gredis.NewScript(`local used=redis.call("HGET",KEYS[1],"check")
if(used)
then
	return ARGV[2]-used
end
redis.call("HMSET",KEYS[1],"code",ARGV[1],"check",0)
redis.call("EXPIRE",KEYS[1],ARGV[3])
return nil`)

	// argv 1 = code
	// argv 2 = max check times
	// argv 3 = expire time(unit second)
	// return nil = key already expired
	// return number = rest check times
	// -1 means code same,check success
	// 0 means all check times used,this key will be expired after expire time
	checkCodeScript = gredis.NewScript(`local data=redis.call("HMGET",KEYS[1],"code","check")
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
return ARGV[2]-data[2]`)
}

// return rest check times and the error
// if error is not nil and rest check times is 0,means all check attempts failed,should block some time
func (d *Dao) RedisSetCode(ctx context.Context, target, action, code string) (int, error) {
	rest, e := setCodeScript.Run(ctx, d.redis, []string{"code_{" + target + "}_" + action}, code, DefaultCheckTimes, DefaultExpireSeconds).Int()
	if e == gredis.Nil {
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

// return rest check times and the error
// if rest check times is -1,means check success
// if rest check times is 0,means all check attempts failed,should block some time
func (d *Dao) RedisCheckCode(ctx context.Context, target, action, code string) (int, error) {
	rest, e := checkCodeScript.Run(ctx, d.redis, []string{"code_{" + target + "}_" + action}, code, DefaultCheckTimes, DefaultExpireSeconds).Int()
	if e == gredis.Nil {
		e = ecode.ErrCodeNotExist
	}
	return rest, e
}

// return code and rest check times and the error
// if rest check times is 0,means all check attempts failed,should block some time
func (d *Dao) RedisGetCode(ctx context.Context, target, action string) (string, int, error) {
	values, e := d.redis.HMGet(ctx, "code_{"+target+"}_"+action, "code", "check").Result()
	if e != nil {
		if e == gredis.Nil {
			e = ecode.ErrCodeNotExist
		}
		return "", 0, e
	}
	if values[0] == nil || len(values[0].(string)) == 0 || values[1] == nil || len(values[1].(string)) == 0 {
		return "", 0, ecode.ErrCodeNotExist
	}
	code := values[0].(string)
	check, e := strconv.Atoi(values[1].(string))
	if e != nil {
		e = ecode.ErrRedisDataBroken
	}
	return code, DefaultCheckTimes - check, e
}

func (d *Dao) RedisDelCode(ctx context.Context, target, action string) error {
	_, e := d.redis.Del(ctx, "code_{"+target+"}_"+action).Result()
	return e
}
