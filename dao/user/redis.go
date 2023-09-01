package user

import (
	"context"
	"encoding/json"

	"github.com/chenjie199234/account/ecode"
	"github.com/chenjie199234/account/model"

	"github.com/chenjie199234/Corelib/redis"
)

func (d *Dao) RedisSetUser(ctx context.Context, userid string, user *model.User) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	if user == nil || user.UserID.IsZero() {
		//set empty key
		_, e = redis.String(c.DoContext(ctx, "SET", userid, "", "EX", 604800))
	} else {
		data, _ := json.Marshal(user)
		_, e = redis.String(c.DoContext(ctx, "SET", userid, data, "EX", 604800))
	}
	return e
}
func (d *Dao) RedisGetUser(ctx context.Context, userid string) (*model.User, error) {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return nil, e
	}
	defer c.Close()
	data, e := redis.Bytes(c.DoContext(ctx, "GET", userid))
	if e != nil {
		if e == redis.ErrNil {
			return nil, nil
		}
		return nil, e
	}
	if len(data) == 0 {
		//this is empty key
		return nil, ecode.ErrUserNotExist
	}
	user := &model.User{}
	if e = json.Unmarshal(data, user); e != nil {
		return nil, e
	}
	return user, nil
}
func (d *Dao) RedisDelUser(ctx context.Context, userid string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "DEL", userid))
	return e
}
func (d *Dao) RedisSetUserIndexTel(ctx context.Context, tel string, userid string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "SET", tel, userid, "EX", 604800))
	return e
}
func (d *Dao) RedisGetUserIndexTel(ctx context.Context, tel string) (string, error) {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return "", e
	}
	defer c.Close()
	userid, e := redis.String(c.DoContext(ctx, "GET", tel))
	if e != nil {
		if e == redis.ErrNil {
			return "", nil
		}
		return "", e
	}
	if userid == "" {
		//this is empty key
		return "", ecode.ErrUserNotExist
	}
	return userid, nil
}
func (d *Dao) RedisDelUserIndexTel(ctx context.Context, tel string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "DEL", tel))
	return e
}
func (d *Dao) RedisSetUserIndexEmail(ctx context.Context, email string, userid string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "SET", email, userid, "EX", 604800))
	return e
}
func (d *Dao) RedisGetUserIndexEmail(ctx context.Context, email string) (string, error) {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return "", e
	}
	defer c.Close()
	userid, e := redis.String(c.DoContext(ctx, "GET", email))
	if e != nil {
		if e == redis.ErrNil {
			return "", nil
		}
		return "", e
	}
	if userid == "" {
		//this is empty key
		return "", ecode.ErrUserNotExist
	}
	return userid, nil
}
func (d *Dao) RedisDelUserIndexEmail(ctx context.Context, email string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "DEL", email))
	return e
}
func (d *Dao) RedisSetUserIndexIDCard(ctx context.Context, idcard string, userid string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "SET", idcard, userid, "EX", 604800))
	return e
}
func (d *Dao) RedisGetUserIndexIDCard(ctx context.Context, idcard string) (string, error) {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return "", e
	}
	defer c.Close()
	userid, e := redis.String(c.DoContext(ctx, "GET", idcard))
	if e != nil {
		if e == redis.ErrNil {
			return "", nil
		}
		return "", e
	}
	if userid == "" {
		//this is empty key
		return "", ecode.ErrUserNotExist
	}
	return userid, nil
}
func (d *Dao) RedisDelUserIndexIDCard(ctx context.Context, idcard string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "DEL", idcard))
	return e
}
func (d *Dao) RedisSetUserIndexNickName(ctx context.Context, nickname string, userid string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "SET", nickname, userid, "EX", 604800))
	return e
}
func (d *Dao) RedisGetUserIndexNickName(ctx context.Context, nickname string) (string, error) {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return "", e
	}
	defer c.Close()
	userid, e := redis.String(c.DoContext(ctx, "GET", nickname))
	if e != nil {
		if e == redis.ErrNil {
			return "", nil
		}
		return "", e
	}
	if userid == "" {
		//this is empty key
		return "", ecode.ErrUserNotExist
	}
	return userid, nil
}
func (d *Dao) RedisDelUserIndexNickName(ctx context.Context, nickname string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "DEL", nickname))
	return e
}
