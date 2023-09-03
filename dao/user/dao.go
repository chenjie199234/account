package user

import (
	"context"
	csql "database/sql"

	"github.com/chenjie199234/account/ecode"
	"github.com/chenjie199234/account/model"

	"github.com/chenjie199234/Corelib/log"
	credis "github.com/chenjie199234/Corelib/redis"
	"go.mongodb.org/mongo-driver/bson/primitive"
	cmongo "go.mongodb.org/mongo-driver/mongo"
)

// Dao this is a data operation layer to operate user service's data
type Dao struct {
	sql   *csql.DB
	redis *credis.Pool
	mongo *cmongo.Client
}

// NewDao Dao is only a data operation layer
// don't write business logic in this package
// business logic should be written in service package
func NewDao(sql *csql.DB, redis *credis.Pool, mongo *cmongo.Client) *Dao {
	return &Dao{
		sql:   sql,
		redis: redis,
		mongo: mongo,
	}
}
func (d *Dao) GetUser(ctx context.Context, callerName string, userid primitive.ObjectID) (*model.User, error) {
	user, e := d.RedisGetUser(ctx, userid.Hex())
	if e != nil {
		log.Error(ctx, "["+callerName+"] redis op failed", map[string]interface{}{"user_id": userid.Hex(), "error": e})
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return nil, e
		}
	}
	if user != nil {
		return user, nil
	}
	if user, e = d.MongoGetUser(ctx, userid); e != nil {
		log.Error(ctx, "["+callerName+"] db op failed", map[string]interface{}{"user_id": userid.Hex(), "error": e})
		if e == ecode.ErrUserNotExist {
			//set redis empty key
			go func() {
				if e := d.RedisSetUser(context.Background(), userid.Hex(), nil); e != nil {
					log.Error(ctx, "["+callerName+"] update redis failed", map[string]interface{}{"user_id": userid.Hex(), "error": e})
				}
			}()
		}
		return nil, e
	}
	//update redis
	go func() {
		if e := d.RedisSetUser(context.Background(), user.UserID.Hex(), user); e != nil {
			log.Error(ctx, "["+callerName+"] update redis failed", map[string]interface{}{"user_id": user.UserID.Hex(), "error": e})
		}
	}()
	if user.Tel != "" {
		go func() {
			if e := d.RedisSetUserIndexTel(context.Background(), user.Tel, user.UserID.Hex()); e != nil {
				log.Error(ctx, "["+callerName+"] update redis failed", map[string]interface{}{"tel": user.Tel, "error": e})
			}
		}()
	}
	if user.Email != "" {
		go func() {
			if e := d.RedisSetUserIndexEmail(context.Background(), user.Email, user.UserID.Hex()); e != nil {
				log.Error(ctx, "["+callerName+"] update redis failed", map[string]interface{}{"email": user.Email, "error": e})
			}
		}()
	}
	if user.IDCard != "" {
		go func() {
			if e := d.RedisSetUserIndexIDCard(context.Background(), user.IDCard, user.UserID.Hex()); e != nil {
				log.Error(ctx, "["+callerName+"] update redis failed", map[string]interface{}{"idcard": user.IDCard, "error": e})
			}
		}()
	}
	if user.NickName != "" {
		go func() {
			if e := d.RedisSetUserIndexNickName(context.Background(), user.NickName, user.UserID.Hex()); e != nil {
				log.Error(ctx, "["+callerName+"] update redis failed", map[string]interface{}{"nick_name": user.NickName, "error": e})
			}
		}()
	}
	return user, nil
}
func (d *Dao) GetUserByTel(ctx context.Context, callerName, tel string) (*model.User, error) {
	user, e := d.RedisGetUserByTel(ctx, tel)
	if e != nil {
		log.Error(ctx, "["+callerName+"] redis op failed", map[string]interface{}{"tel": tel, "error": e})
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return nil, e
		}
	}
	if user != nil {
		return user, nil
	}
	for {
		user, e = d.MongoGetUserByTel(ctx, tel)
		if e != nil {
			log.Error(ctx, "["+callerName+"] db op failed", map[string]interface{}{"tel": tel, "error": e})
			if e == ecode.ErrDBConflict {
				continue
			}
			if e == ecode.ErrUserNotExist {
				//set redis empty key
				go func() {
					if e := d.RedisSetUserIndexTel(context.Background(), tel, ""); e != nil {
						log.Error(ctx, "["+callerName+"] update redis failed", map[string]interface{}{"tel": tel, "error": e})
					}
				}()
			}
			return nil, e
		}
		break
	}
	//update redis
	go func() {
		if e := d.RedisSetUser(context.Background(), user.UserID.Hex(), user); e != nil {
			log.Error(ctx, "["+callerName+"] update redis failed", map[string]interface{}{"user_id": user.UserID.Hex(), "error": e})
		}
	}()
	go func() {
		if e := d.RedisSetUserIndexTel(context.Background(), user.Tel, user.UserID.Hex()); e != nil {
			log.Error(ctx, "["+callerName+"] update redis failed", map[string]interface{}{"tel": user.Tel, "error": e})
		}
	}()
	return user, nil
}
func (d *Dao) GetUserByEmail(ctx context.Context, callerName, email string) (*model.User, error) {
	user, e := d.RedisGetUserByEmail(ctx, email)
	if e != nil {
		log.Error(ctx, "["+callerName+"] redis op failed", map[string]interface{}{"email": email, "error": e})
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return nil, e
		}
	}
	if user != nil {
		return user, nil
	}
	for {
		user, e = d.MongoGetUserByEmail(ctx, email)
		if e != nil {
			log.Error(ctx, "["+callerName+"] db op failed", map[string]interface{}{"email": email, "error": e})
			if e == ecode.ErrDBConflict {
				continue
			}
			if e == ecode.ErrUserNotExist {
				//set redis empty key
				go func() {
					if e := d.RedisSetUserIndexEmail(context.Background(), email, ""); e != nil {
						log.Error(ctx, "["+callerName+"] update redis failed", map[string]interface{}{"email": email, "error": e})
					}
				}()
			}
			return nil, e
		}
		break
	}
	//update redis
	go func() {
		if e := d.RedisSetUser(context.Background(), user.UserID.Hex(), user); e != nil {
			log.Error(ctx, "["+callerName+"] update redis failed", map[string]interface{}{"user_id": user.UserID.Hex(), "error": e})
		}
	}()
	go func() {
		if e := d.RedisSetUserIndexEmail(context.Background(), user.Email, user.UserID.Hex()); e != nil {
			log.Error(ctx, "["+callerName+"] update redis failed", map[string]interface{}{"email": user.Email, "error": e})
		}
	}()
	return user, nil
}
func (d *Dao) GetUserByIDCard(ctx context.Context, callerName, idcard string) (*model.User, error) {
	user, e := d.RedisGetUserByIDCard(ctx, idcard)
	if e != nil {
		log.Error(ctx, "["+callerName+"] redis op failed", map[string]interface{}{"idcard": idcard, "error": e})
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return nil, e
		}
	}
	if user != nil {
		return user, nil
	}
	for {
		user, e = d.MongoGetUserByIDCard(ctx, idcard)
		if e != nil {
			log.Error(ctx, "["+callerName+"] db op failed", map[string]interface{}{"idcard": idcard, "error": e})
			if e == ecode.ErrDBConflict {
				continue
			}
			if e == ecode.ErrUserNotExist {
				//set redis empty key
				go func() {
					if e := d.RedisSetUserIndexIDCard(context.Background(), idcard, ""); e != nil {
						log.Error(ctx, "["+callerName+"] update redis failed", map[string]interface{}{"idcard": idcard, "error": e})
					}
				}()
			}
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		break
	}
	//update redis
	go func() {
		if e := d.RedisSetUser(context.Background(), user.UserID.Hex(), user); e != nil {
			log.Error(ctx, "["+callerName+"] update redis failed", map[string]interface{}{"user_id": user.UserID.Hex(), "error": e})
		}
	}()
	go func() {
		if e := d.RedisSetUserIndexIDCard(context.Background(), user.IDCard, user.UserID.Hex()); e != nil {
			log.Error(ctx, "["+callerName+"] update redis failed", map[string]interface{}{"idcard": user.IDCard, "error": e})
		}
	}()
	return user, nil
}
func (d *Dao) GetUserByNickName(ctx context.Context, callerName, nickname string) (*model.User, error) {
	user, e := d.RedisGetUserByNickName(ctx, nickname)
	if e != nil {
		log.Error(ctx, "["+callerName+"] redis op failed", map[string]interface{}{"nickname": nickname, "error": e})
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return nil, e
		}
	}
	if user != nil {
		return user, nil
	}
	for {
		user, e = d.MongoGetUserByNickName(ctx, nickname)
		if e != nil {
			log.Error(ctx, "["+callerName+"] db op failed", map[string]interface{}{"nick_name": nickname, "error": e})
			if e == ecode.ErrDBConflict {
				continue
			}
			if e == ecode.ErrUserNotExist {
				//set redis empty key
				go func() {
					if e := d.RedisSetUserIndexNickName(context.Background(), nickname, ""); e != nil {
						log.Error(ctx, "["+callerName+"] update redis failed", map[string]interface{}{"nick_name": nickname, "error": e})
					}
				}()
			}
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		break
	}
	//update redis
	go func() {
		if e := d.RedisSetUser(context.Background(), user.UserID.Hex(), user); e != nil {
			log.Error(ctx, "["+callerName+"] update redis failed", map[string]interface{}{"user_id": user.UserID.Hex(), "error": e})
		}
	}()
	go func() {
		if e := d.RedisSetUserIndexNickName(context.Background(), user.NickName, user.UserID.Hex()); e != nil {
			log.Error(ctx, "["+callerName+"] update redis failed", map[string]interface{}{"nick_name": user.NickName, "error": e})
		}
	}()
	return user, nil
}
