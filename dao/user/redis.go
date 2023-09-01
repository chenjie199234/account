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
		_, e = redis.String(c.DoContext(ctx, "SET", "user_{"+userid+"}_info", "", "EX", 604800))
	} else {
		data, _ := json.Marshal(user)
		_, e = redis.String(c.DoContext(ctx, "SET", "user_{"+userid+"}_info", data, "EX", 604800))
	}
	return e
}
func (d *Dao) RedisGetUser(ctx context.Context, userid string) (*model.User, error) {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return nil, e
	}
	defer c.Close()
	data, e := redis.Bytes(c.DoContext(ctx, "GET", "user_{"+userid+"}_info"))
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
	_, e = redis.String(c.DoContext(ctx, "DEL", "user_{"+userid+"}_info"))
	return e
}
func (d *Dao) RedisSetUserIndexTel(ctx context.Context, tel string, userid string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "SET", "tel_{"+tel+"}_index", userid, "EX", 604800))
	return e
}
func (d *Dao) RedisGetUserIndexTel(ctx context.Context, tel string) (string, error) {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return "", e
	}
	defer c.Close()
	userid, e := redis.String(c.DoContext(ctx, "GET", "tel_{"+tel+"}_index"))
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
func (d *Dao) RedisGetUserByTel(ctx context.Context, tel string) (*model.User, error) {
	userid, e := d.RedisGetUserIndexTel(ctx, tel)
	if e != nil {
		return nil, e
	}
	if userid == "" {
		return nil, nil
	}
	user, e := d.RedisGetUser(ctx, userid)
	if e != nil {
		if e == ecode.ErrUserNotExist {
			e = ecode.ErrRedisConflict
		}
		return nil, e
	}
	if user.Tel != tel {
		return nil, ecode.ErrRedisConflict
	}
	return user, nil
}
func (d *Dao) RedisDelUserIndexTel(ctx context.Context, tel string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "DEL", "tel_{"+tel+"}_index"))
	return e
}
func (d *Dao) RedisSetUserIndexEmail(ctx context.Context, email string, userid string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "SET", "email_{"+email+"}_index", userid, "EX", 604800))
	return e
}
func (d *Dao) RedisGetUserIndexEmail(ctx context.Context, email string) (string, error) {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return "", e
	}
	defer c.Close()
	userid, e := redis.String(c.DoContext(ctx, "GET", "email_{"+email+"}_index"))
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
func (d *Dao) RedisGetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	userid, e := d.RedisGetUserIndexEmail(ctx, email)
	if e != nil {
		return nil, e
	}
	if userid == "" {
		return nil, nil
	}
	user, e := d.RedisGetUser(ctx, userid)
	if e != nil {
		if e == ecode.ErrUserNotExist {
			e = ecode.ErrRedisConflict
		}
		return nil, e
	}
	if user.Email != email {
		return nil, ecode.ErrRedisConflict
	}
	return user, nil
}
func (d *Dao) RedisDelUserIndexEmail(ctx context.Context, email string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "DEL", "email_{"+email+"}_index"))
	return e
}
func (d *Dao) RedisSetUserIndexIDCard(ctx context.Context, idcard string, userid string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "SET", "idcard_{"+idcard+"}_index", userid, "EX", 604800))
	return e
}
func (d *Dao) RedisGetUserIndexIDCard(ctx context.Context, idcard string) (string, error) {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return "", e
	}
	defer c.Close()
	userid, e := redis.String(c.DoContext(ctx, "GET", "idcard_{"+idcard+"}_index"))
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
func (d *Dao) RedisGetUserByIDCard(ctx context.Context, idcard string) (*model.User, error) {
	userid, e := d.RedisGetUserIndexIDCard(ctx, idcard)
	if e != nil {
		return nil, e
	}
	if userid == "" {
		return nil, nil
	}
	user, e := d.RedisGetUser(ctx, userid)
	if e != nil {
		if e == ecode.ErrUserNotExist {
			e = ecode.ErrRedisConflict
		}
		return nil, e
	}
	if user.IDCard != idcard {
		return nil, ecode.ErrRedisConflict
	}
	return user, nil
}
func (d *Dao) RedisDelUserIndexIDCard(ctx context.Context, idcard string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "DEL", "idcard_{"+idcard+"}_index"))
	return e
}
func (d *Dao) RedisSetUserIndexNickName(ctx context.Context, nickname string, userid string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "SET", "nickname_{"+nickname+"}_index", userid, "EX", 604800))
	return e
}
func (d *Dao) RedisGetUserIndexNickName(ctx context.Context, nickname string) (string, error) {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return "", e
	}
	defer c.Close()
	userid, e := redis.String(c.DoContext(ctx, "GET", "nickname_{"+nickname+"}_index"))
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
func (d *Dao) RedisGetUserByNickName(ctx context.Context, nickname string) (*model.User, error) {
	userid, e := d.RedisGetUserIndexNickName(ctx, nickname)
	if e != nil {
		return nil, e
	}
	if userid == "" {
		return nil, nil
	}
	user, e := d.RedisGetUser(ctx, userid)
	if e != nil {
		if e == ecode.ErrUserNotExist {
			e = ecode.ErrRedisConflict
		}
		return nil, e
	}
	if user.NickName != nickname {
		return nil, ecode.ErrRedisConflict
	}
	return user, nil
}
func (d *Dao) RedisDelUserIndexNickName(ctx context.Context, nickname string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "DEL", "nickname_{"+nickname+"}_index"))
	return e
}
