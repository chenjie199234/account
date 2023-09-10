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
			//key not exist
			e = ecode.ErrRedisKeyMissing
		}
		return nil, e
	}
	if len(data) == 0 {
		//key exist but value is empty
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
	_, e = redis.Int64(c.DoContext(ctx, "DEL", "user_{"+userid+"}_info"))
	return e
}
func (d *Dao) RedisSetUserTelIndex(ctx context.Context, tel string, userid string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "SET", "tel_{"+tel+"}_index", userid, "EX", 604800))
	return e
}
func (d *Dao) RedisGetUserTelIndex(ctx context.Context, tel string) (string, error) {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return "", e
	}
	defer c.Close()
	userid, e := redis.String(c.DoContext(ctx, "GET", "tel_{"+tel+"}_index"))
	if e != nil {
		if e == redis.ErrNil {
			//key not exist
			e = ecode.ErrRedisKeyMissing
		}
		return "", e
	}
	if userid == "" {
		//key exist but value is empty
		return "", ecode.ErrUserNotExist
	}
	return userid, nil
}
func (d *Dao) RedisGetUserByTel(ctx context.Context, tel string) (*model.User, error) {
	//tel -> userid -> user
	if userid, e := d.RedisGetUserTelIndex(ctx, tel); e != nil {
		return nil, e
	} else if userid == "" {
		return nil, nil
	} else if user, e := d.RedisGetUser(ctx, userid); e != nil {
		if e == ecode.ErrUserNotExist {
			e = ecode.ErrRedisConflict
		}
		return nil, e
	} else if user == nil || user.Tel != tel {
		return nil, ecode.ErrRedisConflict
	} else {
		return user, nil
	}
}
func (d *Dao) RedisDelUserTelIndex(ctx context.Context, tel string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.Int64(c.DoContext(ctx, "DEL", "tel_{"+tel+"}_index"))
	return e
}
func (d *Dao) RedisSetUserEmailIndex(ctx context.Context, email string, userid string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "SET", "email_{"+email+"}_index", userid, "EX", 604800))
	return e
}
func (d *Dao) RedisGetUserEmailIndex(ctx context.Context, email string) (string, error) {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return "", e
	}
	defer c.Close()
	userid, e := redis.String(c.DoContext(ctx, "GET", "email_{"+email+"}_index"))
	if e != nil {
		if e == redis.ErrNil {
			//key not exist
			e = ecode.ErrRedisKeyMissing
		}
		return "", e
	}
	if userid == "" {
		//key exist but value is empty
		return "", ecode.ErrUserNotExist
	}
	return userid, nil
}
func (d *Dao) RedisGetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	//email -> userid -> user
	if userid, e := d.RedisGetUserEmailIndex(ctx, email); e != nil {
		return nil, e
	} else if userid == "" {
		return nil, nil
	} else if user, e := d.RedisGetUser(ctx, userid); e != nil {
		if e == ecode.ErrUserNotExist {
			e = ecode.ErrRedisConflict
		}
		return nil, e
	} else if user == nil || user.Email != email {
		return nil, ecode.ErrRedisConflict
	} else {
		return user, nil
	}
}
func (d *Dao) RedisDelUserEmailIndex(ctx context.Context, email string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.Int64(c.DoContext(ctx, "DEL", "email_{"+email+"}_index"))
	return e
}
func (d *Dao) RedisSetUserIDCardIndex(ctx context.Context, idcard string, userid string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "SET", "idcard_{"+idcard+"}_index", userid, "EX", 604800))
	return e
}
func (d *Dao) RedisGetUserIDCardIndex(ctx context.Context, idcard string) (string, error) {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return "", e
	}
	defer c.Close()
	userid, e := redis.String(c.DoContext(ctx, "GET", "idcard_{"+idcard+"}_index"))
	if e != nil {
		if e == redis.ErrNil {
			//key not exist
			e = ecode.ErrRedisKeyMissing
		}
		return "", e
	}
	if userid == "" {
		//key exist but value is empty
		return "", ecode.ErrUserNotExist
	}
	return userid, nil
}
func (d *Dao) RedisGetUserByIDCard(ctx context.Context, idcard string) (*model.User, error) {
	//idcard -> userid -> user
	if userid, e := d.RedisGetUserIDCardIndex(ctx, idcard); e != nil {
		return nil, e
	} else if userid == "" {
		return nil, nil
	} else if user, e := d.RedisGetUser(ctx, userid); e != nil {
		if e == ecode.ErrUserNotExist {
			e = ecode.ErrRedisConflict
		}
		return nil, e
	} else if user == nil || user.IDCard != idcard {
		return nil, ecode.ErrRedisConflict
	} else {
		return user, nil
	}
}
func (d *Dao) RedisDelUserIDCardIndex(ctx context.Context, idcard string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.Int64(c.DoContext(ctx, "DEL", "idcard_{"+idcard+"}_index"))
	return e
}
func (d *Dao) RedisSetUserNickNameIndex(ctx context.Context, nickname string, userid string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.String(c.DoContext(ctx, "SET", "nickname_{"+nickname+"}_index", userid, "EX", 604800))
	return e
}
func (d *Dao) RedisGetUserNickNameIndex(ctx context.Context, nickname string) (string, error) {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return "", e
	}
	defer c.Close()
	userid, e := redis.String(c.DoContext(ctx, "GET", "nickname_{"+nickname+"}_index"))
	if e != nil {
		if e == redis.ErrNil {
			//key not exist
			e = ecode.ErrRedisKeyMissing
		}
		return "", e
	}
	if userid == "" {
		//key exist but value is empty
		return "", ecode.ErrUserNotExist
	}
	return userid, nil
}
func (d *Dao) RedisGetUserByNickName(ctx context.Context, nickname string) (*model.User, error) {
	//nickname -> userid -> user
	if userid, e := d.RedisGetUserNickNameIndex(ctx, nickname); e != nil {
		return nil, e
	} else if userid == "" {
		return nil, nil
	} else if user, e := d.RedisGetUser(ctx, userid); e != nil {
		if e == ecode.ErrUserNotExist {
			e = ecode.ErrRedisConflict
		}
		return nil, e
	} else if user == nil || user.NickName != nickname {
		return nil, ecode.ErrRedisConflict
	} else {
		return user, nil
	}
}
func (d *Dao) RedisDelUserNickNameIndex(ctx context.Context, nickname string) error {
	c, e := d.redis.GetContext(ctx)
	if e != nil {
		return e
	}
	defer c.Close()
	_, e = redis.Int64(c.DoContext(ctx, "DEL", "nickname_{"+nickname+"}_index"))
	return e
}
