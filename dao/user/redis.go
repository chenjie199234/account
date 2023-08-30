package user

import (
	"context"

	"github.com/chenjie199234/Corelib/redis"
)

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
func (d *Dao) RedisSetCode(ctx context.Context, userid, action string, code string) (bool, error) {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return false, e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "SET", rediskey(userid, action), code, "EX", 300, "NX"))
	if e == nil {
		return true, nil
	}
	if e == redis.ErrNil {
		return false, nil
	}
	return false, e
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
