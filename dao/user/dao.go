package user

import (
	"context"
	csql "database/sql"
	"unsafe"

	"github.com/chenjie199234/account/ecode"
	"github.com/chenjie199234/account/model"

	"github.com/chenjie199234/Corelib/log"
	credis "github.com/chenjie199234/Corelib/redis"
	"github.com/chenjie199234/Corelib/util/oneshot"
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

// user

func (d *Dao) GetUser(ctx context.Context, userid primitive.ObjectID) (*model.User, error) {
	if user, e := d.RedisGetUser(ctx, userid.Hex()); e != nil {
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return nil, e
		}
		log.Error(ctx, "[dao.GetUser] redis op failed", map[string]interface{}{"user_id": userid.Hex(), "error": e})
	} else if user != nil {
		return user, nil
	}
	//redis error or redis not exist,we need to query db
	unsafeUser, e := oneshot.Do("GetUser_"+userid.Hex(), func() (unsafe.Pointer, error) {
		user, e := d.MongoGetUser(ctx, userid)
		if e != nil {
			log.Error(nil, "[dao.GetUser] db op failed", map[string]interface{}{"user_id": userid.Hex(), "error": e})
			if e == ecode.ErrUserNotExist {
				//set redis empty key
				go func() {
					if e := d.RedisSetUser(context.Background(), userid.Hex(), nil); e != nil {
						log.Error(nil, "[dao.GetUser] update redis failed", map[string]interface{}{"user_id": userid.Hex(), "error": e})
					}
				}()
			}
			return nil, e
		}
		//update redis
		go func() {
			if e := d.RedisSetUser(context.Background(), user.UserID.Hex(), user); e != nil {
				log.Error(nil, "[dao.GetUser] update redis failed", map[string]interface{}{"user_id": user.UserID.Hex(), "error": e})
			}
		}()
		return unsafe.Pointer(user), e
	})
	return (*model.User)(unsafeUser), e
}
func (d *Dao) GetUserByTel(ctx context.Context, tel string) (*model.User, error) {
	if user, e := d.RedisGetUserByTel(ctx, tel); e != nil {
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return nil, e
		}
		log.Error(nil, "[dao.GetUserByTel] redis op failed", map[string]interface{}{"tel": tel, "error": e})
	} else if user != nil {
		return user, nil
	}
	//redis error or redis not exist,we need to query db
	unsafeUser, e := oneshot.Do("GetUserByTel_"+tel, func() (unsafe.Pointer, error) {
		var user *model.User
		var e error
		for {
			user, e = d.MongoGetUserByTel(ctx, tel)
			if e != nil {
				log.Error(nil, "[dao.GetUserByTel] db op failed", map[string]interface{}{"tel": tel, "error": e})
				if e == ecode.ErrDBConflict {
					continue
				}
				if e == ecode.ErrUserNotExist {
					//set redis empty key
					go func() {
						if e := d.RedisSetUserTelIndex(context.Background(), tel, ""); e != nil {
							log.Error(nil, "[dao.GetUserByTel] update redis failed", map[string]interface{}{"tel": tel, "error": e})
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
				log.Error(nil, "[dao.GetUserByTel] update redis failed", map[string]interface{}{"user_id": user.UserID.Hex(), "error": e})
			}
		}()
		go func() {
			if e := d.RedisSetUserTelIndex(context.Background(), user.Tel, user.UserID.Hex()); e != nil {
				log.Error(nil, "[dao.GetUserByTel] update redis failed", map[string]interface{}{"tel": user.Tel, "error": e})
			}
		}()
		return unsafe.Pointer(user), e
	})
	return (*model.User)(unsafeUser), e
}
func (d *Dao) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	if user, e := d.RedisGetUserByEmail(ctx, email); e != nil {
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return nil, e
		}
		log.Error(nil, "[dao.GetUserByEmail] redis op failed", map[string]interface{}{"email": email, "error": e})
	} else if user != nil {
		return user, nil
	}
	//redis error or redis not exist,we need to query db
	unsafeUser, e := oneshot.Do("GetUserByEmail_"+email, func() (unsafe.Pointer, error) {
		var user *model.User
		var e error
		for {
			user, e = d.MongoGetUserByEmail(ctx, email)
			if e != nil {
				log.Error(nil, "[dao.GetUserByEmail] db op failed", map[string]interface{}{"email": email, "error": e})
				if e == ecode.ErrDBConflict {
					continue
				}
				if e == ecode.ErrUserNotExist {
					//set redis empty key
					go func() {
						if e := d.RedisSetUserEmailIndex(context.Background(), email, ""); e != nil {
							log.Error(nil, "[dao.GetUserByEmail] update redis failed", map[string]interface{}{"email": email, "error": e})
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
				log.Error(nil, "[dao.GetUserByEmail] update redis failed", map[string]interface{}{"user_id": user.UserID.Hex(), "error": e})
			}
		}()
		go func() {
			if e := d.RedisSetUserEmailIndex(context.Background(), user.Email, user.UserID.Hex()); e != nil {
				log.Error(nil, "[dao.GetUserByEmail] update redis failed", map[string]interface{}{"email": user.Email, "error": e})
			}
		}()
		return unsafe.Pointer(user), e
	})
	return (*model.User)(unsafeUser), e
}
func (d *Dao) GetUserByIDCard(ctx context.Context, idcard string) (*model.User, error) {
	if user, e := d.RedisGetUserByIDCard(ctx, idcard); e != nil {
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return nil, e
		}
		log.Error(nil, "[dao.GetUserByIDCard] redis op failed", map[string]interface{}{"idcard": idcard, "error": e})
	} else if user != nil {
		return user, nil
	}
	//redis error or redis not exist,we need to query db
	unsafeUser, e := oneshot.Do("GetUserByIDCard_"+idcard, func() (unsafe.Pointer, error) {
		var user *model.User
		var e error
		for {
			user, e = d.MongoGetUserByIDCard(ctx, idcard)
			if e != nil {
				log.Error(nil, "[dao.GetUserByIDCard] db op failed", map[string]interface{}{"idcard": idcard, "error": e})
				if e == ecode.ErrDBConflict {
					continue
				}
				if e == ecode.ErrUserNotExist {
					//set redis empty key
					go func() {
						if e := d.RedisSetUserIDCardIndex(context.Background(), idcard, ""); e != nil {
							log.Error(nil, "[dao.GetUserByIDCard] update redis failed", map[string]interface{}{"idcard": idcard, "error": e})
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
				log.Error(nil, "[dao.GetUserByIDCard] update redis failed", map[string]interface{}{"user_id": user.UserID.Hex(), "error": e})
			}
		}()
		go func() {
			if e := d.RedisSetUserIDCardIndex(context.Background(), user.IDCard, user.UserID.Hex()); e != nil {
				log.Error(nil, "[dao.GetUserByIDCard] update redis failed", map[string]interface{}{"idcard": user.IDCard, "error": e})
			}
		}()
		return unsafe.Pointer(user), e
	})
	return (*model.User)(unsafeUser), e
}
func (d *Dao) GetUserByNickName(ctx context.Context, nickname string) (*model.User, error) {
	if user, e := d.RedisGetUserByNickName(ctx, nickname); e != nil {
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return nil, e
		}
		log.Error(nil, "[dao.GetUserByNickName] redis op failed", map[string]interface{}{"nickname": nickname, "error": e})
	} else if user != nil {
		return user, nil
	}
	//redis error or redis not exist,we need to query db
	unsafeUser, e := oneshot.Do("GetUserByNickName_"+nickname, func() (unsafe.Pointer, error) {
		var user *model.User
		var e error
		for {
			user, e = d.MongoGetUserByNickName(ctx, nickname)
			if e != nil {
				log.Error(nil, "[dao.GetUserByNickName] db op failed", map[string]interface{}{"nick_name": nickname, "error": e})
				if e == ecode.ErrDBConflict {
					continue
				}
				if e == ecode.ErrUserNotExist {
					//set redis empty key
					go func() {
						if e := d.RedisSetUserNickNameIndex(context.Background(), nickname, ""); e != nil {
							log.Error(nil, "[dao.GetUserByNickName] update redis failed", map[string]interface{}{"nick_name": nickname, "error": e})
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
				log.Error(nil, "[dao.GetUserByNickName] update redis failed", map[string]interface{}{"user_id": user.UserID.Hex(), "error": e})
			}
		}()
		go func() {
			if e := d.RedisSetUserNickNameIndex(context.Background(), user.NickName, user.UserID.Hex()); e != nil {
				log.Error(nil, "[dao.GetUserByNickName] update redis failed", map[string]interface{}{"nick_name": user.NickName, "error": e})
			}
		}()
		return unsafe.Pointer(user), e
	})
	return (*model.User)(unsafeUser), e
}

// tel index

func (d *Dao) GetUserTelIndex(ctx context.Context, tel string) (string, error) {
	if userid, e := d.RedisGetUserTelIndex(ctx, tel); e != nil {
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return userid, e
		}
		log.Error(nil, "[dao.GetUserTelIndex] redis op failed", map[string]interface{}{"tel": tel, "error": e})
	} else if userid != "" {
		return userid, nil
	}
	//redis error or redis not exist,we need to query db
	unsafeUserid, e := oneshot.Do("GetUserTelIndex_"+tel, func() (unsafe.Pointer, error) {
		index, e := d.MongoGetUserTelIndex(ctx, tel)
		if e != nil {
			log.Error(nil, "[dao.GetUserTelIndex] db op failed", map[string]interface{}{"tel": tel, "error": e})
			if e == ecode.ErrUserNotExist {
				//set redis empty key
				go func() {
					if e := d.RedisSetUserTelIndex(context.Background(), tel, ""); e != nil {
						log.Error(nil, "[dao.GetUserTelIndex] update redis failed", map[string]interface{}{"tel": tel, "error": e})
					}
				}()
			}
			return nil, e
		}
		userid := index.UserID.Hex()
		return unsafe.Pointer(&userid), nil
	})
	return *(*string)(unsafeUserid), e
}

// email index

func (d *Dao) GetUserEmailIndex(ctx context.Context, email string) (string, error) {
	if userid, e := d.RedisGetUserEmailIndex(ctx, email); e != nil {
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return userid, e
		}
		log.Error(nil, "[dao.GetUserEmailIndex] redis op failed", map[string]interface{}{"email": email, "error": e})
	} else if userid != "" {
		return userid, nil
	}
	//redis error or redis not exist,we need to query db
	unsafeUserid, e := oneshot.Do("GetUserEmailIndex_"+email, func() (unsafe.Pointer, error) {
		index, e := d.MongoGetUserEmailIndex(ctx, email)
		if e != nil {
			log.Error(nil, "[dao.GetUserEmailIndex] db op failed", map[string]interface{}{"email": email, "error": e})
			if e == ecode.ErrUserNotExist {
				//set redis empty key
				go func() {
					if e := d.RedisSetUserEmailIndex(context.Background(), email, ""); e != nil {
						log.Error(nil, "[dao.GetUserEmailIndex] update redis failed", map[string]interface{}{"email": email, "error": e})
					}
				}()
			}
			return nil, e
		}
		userid := index.UserID.Hex()
		return unsafe.Pointer(&userid), nil
	})
	return *(*string)(unsafeUserid), e
}

// idcard index

func (d *Dao) GetUserIDCardIndex(ctx context.Context, idcard string) (string, error) {
	if userid, e := d.RedisGetUserIDCardIndex(ctx, idcard); e != nil {
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return userid, e
		}
		log.Error(nil, "[dao.GetUserIDCardIndex] redis op failed", map[string]interface{}{"idcard": idcard, "error": e})
	} else if userid != "" {
		return userid, nil
	}
	//redis error or redis not exist,we need to query db
	unsafeUserid, e := oneshot.Do("GetUserIDcardIndex_"+idcard, func() (unsafe.Pointer, error) {
		index, e := d.MongoGetUserIDCardIndex(ctx, idcard)
		if e != nil {
			log.Error(nil, "[dao.GetUserIDCardIndex] db op failed", map[string]interface{}{"idcard": idcard, "error": e})
			if e == ecode.ErrUserNotExist {
				//set redis empty key
				go func() {
					if e := d.RedisSetUserIDCardIndex(context.Background(), idcard, ""); e != nil {
						log.Error(nil, "[dao.GetUserIDCardIndex] update redis failed", map[string]interface{}{"idcard": idcard, "error": e})
					}
				}()
			}
			return nil, e
		}
		userid := index.UserID.Hex()
		return unsafe.Pointer(&userid), nil
	})
	return *(*string)(unsafeUserid), e
}

// nickname index

func (d *Dao) GetUserNickNameIndex(ctx context.Context, nickname string) (string, error) {
	if userid, e := d.RedisGetUserNickNameIndex(ctx, nickname); e != nil {
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return userid, e
		}
		log.Error(nil, "[dao.GetUserNickNameIndex] redis op failed", map[string]interface{}{"nick_name": nickname, "error": e})
	} else if userid != "" {
		return userid, nil
	}
	//redis error or redis not exist,we need to query db
	unsafeUserid, e := oneshot.Do("GetUserNickNameIndex_"+nickname, func() (unsafe.Pointer, error) {
		index, e := d.MongoGetUserNickNameIndex(ctx, nickname)
		if e != nil {
			log.Error(nil, "[dao.GetUserNickNameIndex] db op failed", map[string]interface{}{"nick_name": nickname, "error": e})
			if e == ecode.ErrUserNotExist {
				//set redis empty key
				go func() {
					if e := d.RedisSetUserNickNameIndex(context.Background(), nickname, ""); e != nil {
						log.Error(nil, "[dao.GetUserNickNameIndex] update redis failed", map[string]interface{}{"nick_name": nickname, "error": e})
					}
				}()
			}
			return nil, e
		}
		userid := index.UserID.Hex()
		return unsafe.Pointer(&userid), nil
	})
	return *(*string)(unsafeUserid), e
}
