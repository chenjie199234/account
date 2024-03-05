package user

import (
	"context"
	"encoding/json"
	"time"

	"github.com/chenjie199234/account/ecode"
	"github.com/chenjie199234/account/model"
)

func (d *Dao) RedisSetUser(ctx context.Context, userid string, user *model.User) error {
	var e error
	if user == nil || user.UserID.IsZero() {
		//set empty key
		_, e = d.redis.SetEx(ctx, "user_{"+userid+"}_info", "", 7*24*time.Hour).Result()
	} else {
		data, _ := json.Marshal(user)
		_, e = d.redis.SetEx(ctx, "user_{"+userid+"}_info", data, 7*24*time.Hour).Result()
	}
	return e
}
func (d *Dao) RedisGetUser(ctx context.Context, userid string) (*model.User, error) {
	data, e := d.redis.GetEx(ctx, "user_{"+userid+"}_info", 7*24*time.Hour).Bytes()
	if e != nil {
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
	_, e := d.redis.Del(ctx, "user_{"+userid+"}_info").Result()
	return e
}

// oauth
func (d *Dao) RedisSetUserOAuthIndex(ctx context.Context, oauthservicename, oauthid, userid string) error {
	_, e := d.redis.SetEx(ctx, "oauth_{"+oauthservicename+"|"+oauthid+"}_index", userid, 7*24*time.Hour).Result()
	return e
}
func (d *Dao) RedisGetUserOAuthIndex(ctx context.Context, oauthservicename, oauthid string) (string, error) {
	userid, e := d.redis.GetEx(ctx, "oauth_{"+oauthservicename+"|"+oauthid+"}_index", 7*24*time.Hour).Result()
	if e != nil {
		return "", e
	}
	if userid == "" {
		//key exist but value is empty
		return "", ecode.ErrUserNotExist
	}
	return userid, nil
}
func (d *Dao) RedisGetUserByOAuth(ctx context.Context, oauthservicename, oauthid string) (*model.User, error) {
	//oauth -> userid -> user
	if userid, e := d.RedisGetUserOAuthIndex(ctx, oauthservicename, oauthid); e != nil {
		return nil, e
	} else if user, e := d.RedisGetUser(ctx, userid); e != nil {
		if e == ecode.ErrUserNotExist {
			e = ecode.ErrCacheDataConflict
		}
		return nil, e
	} else if user.OAuths[oauthservicename] != oauthid {
		return nil, ecode.ErrCacheDataConflict
	} else {
		return user, nil
	}
}
func (d *Dao) RedisDelUserOAuthIndex(ctx context.Context, oauthservicename, oauthid string) error {
	_, e := d.redis.Del(ctx, "oauth_{"+oauthservicename+"|"+oauthid+"}_index").Result()
	return e
}

// tel
func (d *Dao) RedisSetUserTelIndex(ctx context.Context, tel string, userid string) error {
	_, e := d.redis.SetEx(ctx, "tel_{"+tel+"}_index", userid, 7*24*time.Hour).Result()
	return e
}
func (d *Dao) RedisGetUserTelIndex(ctx context.Context, tel string) (string, error) {
	userid, e := d.redis.GetEx(ctx, "tel_{"+tel+"}_index", 7*24*time.Hour).Result()
	if e != nil {
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
	} else if user, e := d.RedisGetUser(ctx, userid); e != nil {
		if e == ecode.ErrUserNotExist {
			e = ecode.ErrCacheDataConflict
		}
		return nil, e
	} else if user.Tel != tel {
		return nil, ecode.ErrCacheDataConflict
	} else {
		return user, nil
	}
}
func (d *Dao) RedisDelUserTelIndex(ctx context.Context, tel string) error {
	_, e := d.redis.Del(ctx, "tel_{"+tel+"}_index").Result()
	return e
}

// email
func (d *Dao) RedisSetUserEmailIndex(ctx context.Context, email string, userid string) error {
	_, e := d.redis.SetEx(ctx, "email_{"+email+"}_index", userid, 7*24*time.Hour).Result()
	return e
}
func (d *Dao) RedisGetUserEmailIndex(ctx context.Context, email string) (string, error) {
	userid, e := d.redis.GetEx(ctx, "email_{"+email+"}_index", 7*24*time.Hour).Result()
	if e != nil {
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
	} else if user, e := d.RedisGetUser(ctx, userid); e != nil {
		if e == ecode.ErrUserNotExist {
			e = ecode.ErrCacheDataConflict
		}
		return nil, e
	} else if user.Email != email {
		return nil, ecode.ErrCacheDataConflict
	} else {
		return user, nil
	}
}
func (d *Dao) RedisDelUserEmailIndex(ctx context.Context, email string) error {
	_, e := d.redis.Del(ctx, "email_{"+email+"}_index").Result()
	return e
}

// idcard
func (d *Dao) RedisSetUserIDCardIndex(ctx context.Context, idcard string, userid string) error {
	_, e := d.redis.SetEx(ctx, "idcard_{"+idcard+"}_index", userid, 7*24*time.Hour).Result()
	return e
}
func (d *Dao) RedisGetUserIDCardIndex(ctx context.Context, idcard string) (string, error) {
	userid, e := d.redis.GetEx(ctx, "idcard_{"+idcard+"}_index", 7*24*time.Hour).Result()
	if e != nil {
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
	} else if user, e := d.RedisGetUser(ctx, userid); e != nil {
		if e == ecode.ErrUserNotExist {
			e = ecode.ErrCacheDataConflict
		}
		return nil, e
	} else if user.IDCard != idcard {
		return nil, ecode.ErrCacheDataConflict
	} else {
		return user, nil
	}
}
func (d *Dao) RedisDelUserIDCardIndex(ctx context.Context, idcard string) error {
	_, e := d.redis.Del(ctx, "idcard_{"+idcard+"}_index").Result()
	return e
}
