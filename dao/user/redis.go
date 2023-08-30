package user

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"strings"

	"github.com/chenjie199234/Corelib/redis"
	"github.com/chenjie199234/account/ecode"
)

const DefaultExpireSeconds = 300

func init() {
	h := sha1.Sum([]byte(setcode))
	hsetcode = hex.EncodeToString(h[:])
	h = sha1.Sum([]byte(checkcode))
	hcheckcode = hex.EncodeToString(h[:])
}

// actions
const (
	LoginEmail = "loginemail"
	LoginTel   = "logintel"
	OldEmail   = "oldemail"
	OldTel     = "oldtel"
	NewEmail   = "newemail"
	NewTel     = "newtel"
)

func rediskey(userid, action string) string {
	return "{" + userid + "}_" + action + "_code"
}

// argv 1 = code
// argc 2 = expire time(unit second)
// return "OK" = set success
// return nil = key already exist
const setcode = `if(redis.call("EXISTS",KEYS[1])==1)
then
	return nil
end
redis.call("HMSET",KEYS[1],"code",ARGV[1],"check",0)
redis.call("EXPIRE",KEYS[1],ARGV[2])
return "OK"`

var hsetcode string

func (d *Dao) RedisSetCode(ctx context.Context, userid, action string, code string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "EVALSHA", hsetcode, 1, rediskey(userid, action), code, DefaultExpireSeconds))
	if e != nil && strings.Contains(e.Error(), "NOSCRIPT") {
		_, e = redis.String(c.DoContext(ctx, "EVAL", setcode, 1, rediskey(userid, action), code, DefaultExpireSeconds))
	}
	if e == redis.ErrNil {
		e = ecode.ErrCodeAlreadySend
	}
	return e
}

// argv 1 = code
// argv 2 = max check times
// argv 3 = expire time(unit second)
// return nil = key already expired
// return number = rest check times
// -1 means code same,check success
// 0 means all check times used,this key will be expired after expire time
const checkcode = `local data=redis.call("HMGET",KEYS[1],"code","check")
if(data[1]==nil or data[2]==nil)
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
if(redis.call("EXPIRE",KEYS[1],ARGV[2])==0)
then
	return nil
end
data[2]=redis.call("HINCRBY",KEYS[1],"check",1)
return ARGV[2]-data[2]`

var hcheckcode string

func (d *Dao) RedisCheckCode(ctx context.Context, userid, action string, code string, delWhenSame bool) (int, error) {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return 0, e
	}
	defer c.Close()
	var rest int
	rest, e = redis.Int(c.DoContext(ctx, "EVALSHA", hcheckcode, 1, rediskey(userid, action), code, 5, DefaultExpireSeconds))
	if e != nil && strings.Contains(e.Error(), "NOSCRIPT") {
		rest, e = redis.Int(c.DoContext(ctx, "EVAL", checkcode, 1, rediskey(userid, action), code, 5, DefaultExpireSeconds))
	}
	if e == redis.ErrNil {
		e = ecode.ErrCodeAlreadyExpire
	}
	return rest, e
}

func (d *Dao) RedisDelCode(ctx context.Context, userid, action string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = c.DoContext(ctx, "DEL", rediskey(userid, action))
	return e
}
