package user

import (
	"context"
	"time"
	"unsafe"

	"github.com/chenjie199234/account/ecode"
	"github.com/chenjie199234/account/model"

	"github.com/chenjie199234/Corelib/log"
	cmongo "github.com/chenjie199234/Corelib/mongo"
	cmysql "github.com/chenjie199234/Corelib/mysql"
	credis "github.com/chenjie199234/Corelib/redis"
	"github.com/chenjie199234/Corelib/util/oneshot"
	gredis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Dao this is a data operation layer to operate user service's data
type Dao struct {
	mysql *cmysql.Client
	redis *credis.Client
	mongo *cmongo.Client
}

// NewDao Dao is only a data operation layer
// don't write business logic in this package
// business logic should be written in service package
func NewDao(mysql *cmysql.Client, redis *credis.Client, mongo *cmongo.Client) *Dao {
	return &Dao{
		mysql: mysql,
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
		if e != gredis.Nil {
			log.Error(ctx, "[dao.GetUser] redis op failed", log.String("user_id", userid.Hex()), log.CError(e))
		}
	} else if user != nil {
		return user, nil
	}
	//redis error or redis not exist,we need to query db
	unsafeUser, e := oneshot.Do("GetUser_"+userid.Hex(), func() (unsafe.Pointer, error) {
		user, e := d.MongoGetUser(ctx, userid)
		if e != nil {
			log.Error(nil, "[dao.GetUser] db op failed", log.String("user_id", userid.Hex()), log.CError(e))
			if e == ecode.ErrUserNotExist {
				//set redis empty key
				go func() {
					if e := d.RedisSetUser(context.Background(), userid.Hex(), nil); e != nil {
						log.Error(nil, "[dao.GetUser] update redis failed", log.String("user_id", userid.Hex()), log.CError(e))
					}
				}()
			}
			return nil, e
		}
		//update redis
		go func() {
			if e := d.RedisSetUser(context.Background(), user.UserID.Hex(), user); e != nil {
				log.Error(nil, "[dao.GetUser] update redis failed", log.String("user_id", userid.Hex()), log.CError(e))
			}
		}()
		return unsafe.Pointer(user), e
	})
	return (*model.User)(unsafeUser), e
}
func (d *Dao) GetUserByOAuth(ctx context.Context, oauthservicename, oauthid string) (*model.User, error) {
	if user, e := d.RedisGetUserByOAuth(ctx, oauthservicename, oauthid); e != nil {
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return nil, e
		}
		if e != gredis.Nil && e != ecode.ErrCacheDataConflict {
			log.Error(ctx, "[dao.GetUserByOAuth] redis op failed", log.String(oauthservicename, oauthid), log.CError(e))
		}
	} else if user != nil {
		return user, nil
	}
	//redis error or redis not exist,we need to query db
	unsafeUser, e := oneshot.Do("GetUserByOAuth_"+oauthservicename+"|"+oauthid, func() (unsafe.Pointer, error) {
		var user *model.User
		var e error
		for {
			user, e = d.MongoGetUserByOAuth(ctx, oauthservicename, oauthid)
			if e != nil {
				log.Error(nil, "[dao.GetUserByOAuth] db op failed", log.String(oauthservicename, oauthid), log.CError(e))
				if e == ecode.ErrDBDataConflict {
					time.Sleep(time.Millisecond * 5)
					continue
				}
				if e == ecode.ErrUserNotExist {
					//set redis empty key
					go func() {
						if e := d.RedisSetUserOAuthIndex(context.Background(), oauthservicename, oauthid, ""); e != nil {
							log.Error(nil, "[dao.GetUserByOAuth] update redis failed", log.String(oauthservicename, oauthid), log.CError(e))
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
				log.Error(nil, "[dao.GetUserByOAuth] update redis failed", log.String(oauthservicename, oauthid), log.CError(e))
			}
		}()
		go func() {
			if e := d.RedisSetUserOAuthIndex(context.Background(), oauthservicename, oauthid, user.UserID.Hex()); e != nil {
				log.Error(nil, "[dao.GetUserByOAuth] update redis failed", log.String(oauthservicename, oauthid), log.CError(e))
			}
		}()
		return unsafe.Pointer(user), e
	})
	return (*model.User)(unsafeUser), e
}
func (d *Dao) GetOrCreateUserByOAuth(ctx context.Context, oauthservicename, oauthid string) (*model.User, error) {
	if user, e := d.RedisGetUserByOAuth(ctx, oauthservicename, oauthid); e != nil {
		if e != gredis.Nil && e != ecode.ErrUserNotExist && e != ecode.ErrCacheDataConflict {
			log.Error(ctx, "[dao.GetOrCreateUserByOAuth] redis op failed", log.String(oauthservicename, oauthid), log.CError(e))
		}
	} else if user != nil {
		return user, nil
	}
	//redis error or redis not exist,we need to query db
	unsafeUser, e := oneshot.Do("GetOrCreateUserByOAuth_"+oauthservicename+"|"+oauthid, func() (unsafe.Pointer, error) {
		var user *model.User
		var e error
		for {
			user, e = d.MongoCreateUserByOAuth(ctx, oauthservicename, oauthid)
			if e != nil {
				log.Error(nil, "[dao.GetOrCreateUserByOAuth] db op failed", log.String(oauthservicename, oauthid), log.CError(e))
				if e == ecode.ErrDBDataConflict {
					time.Sleep(time.Millisecond * 5)
					continue
				}
				return nil, e
			}
			break
		}
		//update redis
		go func() {
			if e := d.RedisSetUser(context.Background(), user.UserID.Hex(), user); e != nil {
				log.Error(nil, "[dao.GetOrCreateUserByOAuth] update redis failed", log.String("user_id", user.UserID.Hex()), log.CError(e))
			}
		}()
		go func() {
			if e := d.RedisSetUserOAuthIndex(context.Background(), oauthservicename, oauthid, user.UserID.Hex()); e != nil {
				log.Error(nil, "[dao.GetOrCreateUserByOAuth] update redis failed", log.String(oauthservicename, oauthid), log.CError(e))
			}
		}()
		return unsafe.Pointer(user), nil
	})
	return (*model.User)(unsafeUser), e
}
func (d *Dao) GetUserByTel(ctx context.Context, tel string) (*model.User, error) {
	if user, e := d.RedisGetUserByTel(ctx, tel); e != nil {
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return nil, e
		}
		if e != gredis.Nil && e != ecode.ErrCacheDataConflict {
			log.Error(ctx, "[dao.GetUserByTel] redis op failed", log.String("tel", tel), log.CError(e))
		}
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
				log.Error(nil, "[dao.GetUserByTel] db op failed", log.String("tel", tel), log.CError(e))
				if e == ecode.ErrDBDataConflict {
					time.Sleep(time.Millisecond * 5)
					continue
				}
				if e == ecode.ErrUserNotExist {
					//set redis empty key
					go func() {
						if e := d.RedisSetUserTelIndex(context.Background(), tel, ""); e != nil {
							log.Error(nil, "[dao.GetUserByTel] update redis failed", log.String("tel", tel), log.CError(e))
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
				log.Error(nil, "[dao.GetUserByTel] update redis failed", log.String("user_id", user.UserID.Hex()), log.CError(e))
			}
		}()
		go func() {
			if e := d.RedisSetUserTelIndex(context.Background(), user.Tel, user.UserID.Hex()); e != nil {
				log.Error(nil, "[dao.GetUserByTel] update redis failed", log.String("tel", user.Tel), log.CError(e))
			}
		}()
		return unsafe.Pointer(user), e
	})
	return (*model.User)(unsafeUser), e
}
func (d *Dao) GetOrCreateUserByTel(ctx context.Context, tel string) (*model.User, error) {
	if user, e := d.RedisGetUserByTel(ctx, tel); e != nil {
		if e != gredis.Nil && e != ecode.ErrUserNotExist && e != ecode.ErrCacheDataConflict {
			log.Error(ctx, "[dao.GetOrCreateUserByTel] redis op failed", log.String("tel", tel), log.CError(e))
		}
	} else if user != nil {
		return user, nil
	}
	//redis error or redis not exist,we need to query db
	unsafeUser, e := oneshot.Do("GetOrCreateUserByTel_"+tel, func() (unsafe.Pointer, error) {
		var user *model.User
		var e error
		for {
			user, e = d.MongoCreateUserByTel(ctx, tel)
			if e != nil {
				log.Error(nil, "[dao.GetOrCreateUserByTel] db op failed", log.String("tel", tel), log.CError(e))
				if e == ecode.ErrDBDataConflict {
					time.Sleep(time.Millisecond * 5)
					continue
				}
				return nil, e
			}
			break
		}
		//update redis
		go func() {
			if e := d.RedisSetUser(context.Background(), user.UserID.Hex(), user); e != nil {
				log.Error(nil, "[dao.GetOrCreateUserByTel] update redis failed", log.String("user_id", user.UserID.Hex()), log.CError(e))
			}
		}()
		go func() {
			if e := d.RedisSetUserTelIndex(context.Background(), user.Tel, user.UserID.Hex()); e != nil {
				log.Error(nil, "[dao.GetOrCreateUserByTel] update redis failed", log.String("tel", user.Tel), log.CError(e))
			}
		}()
		return unsafe.Pointer(user), nil
	})
	return (*model.User)(unsafeUser), e
}
func (d *Dao) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	if user, e := d.RedisGetUserByEmail(ctx, email); e != nil {
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return nil, e
		}
		if e != gredis.Nil && e != ecode.ErrCacheDataConflict {
			log.Error(ctx, "[dao.GetUserByEmail] redis op failed", log.String("email", email), log.CError(e))
		}
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
				log.Error(nil, "[dao.GetUserByEmail] db op failed", log.String("email", email), log.CError(e))
				if e == ecode.ErrDBDataConflict {
					time.Sleep(time.Millisecond * 5)
					continue
				}
				if e == ecode.ErrUserNotExist {
					//set redis empty key
					go func() {
						if e := d.RedisSetUserEmailIndex(context.Background(), email, ""); e != nil {
							log.Error(nil, "[dao.GetUserByEmail] update redis failed", log.String("email", email), log.CError(e))
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
				log.Error(nil, "[dao.GetUserByEmail] update redis failed", log.String("user_id", user.UserID.Hex()), log.CError(e))
			}
		}()
		go func() {
			if e := d.RedisSetUserEmailIndex(context.Background(), user.Email, user.UserID.Hex()); e != nil {
				log.Error(nil, "[dao.GetUserByEmail] update redis failed", log.String("email", user.Email), log.CError(e))
			}
		}()
		return unsafe.Pointer(user), e
	})
	return (*model.User)(unsafeUser), e
}
func (d *Dao) GetOrCreateUserByEmail(ctx context.Context, email string) (*model.User, error) {
	if user, e := d.RedisGetUserByEmail(ctx, email); e != nil {
		if e != gredis.Nil && e != ecode.ErrUserNotExist && e != ecode.ErrCacheDataConflict {
			log.Error(ctx, "[dao.GetOrCreateUserByEmail] redis op failed", log.String("email", email), log.CError(e))
		}
	} else if user != nil {
		return user, nil
	}
	//redis error or redis not exist,we need to query db
	unsafeUser, e := oneshot.Do("GetOrCreateUserByEmail_"+email, func() (unsafe.Pointer, error) {
		var user *model.User
		var e error
		for {
			user, e = d.MongoCreateUserByEmail(ctx, email)
			if e != nil {
				log.Error(nil, "[dao.GetOrCreateUserByEmail] db op failed", log.String("email", email), log.CError(e))
				if e == ecode.ErrDBDataConflict {
					time.Sleep(time.Millisecond * 5)
					continue
				}
				return nil, e
			}
			break
		}
		//update redis
		go func() {
			if e := d.RedisSetUser(context.Background(), user.UserID.Hex(), user); e != nil {
				log.Error(nil, "[dao.GetOrCreateUserByEmail] update redis failed", log.String("user_id", user.UserID.Hex()), log.CError(e))
			}
		}()
		go func() {
			if e := d.RedisSetUserEmailIndex(context.Background(), user.Email, user.UserID.Hex()); e != nil {
				log.Error(nil, "[dao.GetOrCreateUserByEmail] update redis failed", log.String("email", user.Email), log.CError(e))
			}
		}()
		return unsafe.Pointer(user), nil
	})
	return (*model.User)(unsafeUser), e
}
func (d *Dao) GetUserByIDCard(ctx context.Context, idcard string) (*model.User, error) {
	if user, e := d.RedisGetUserByIDCard(ctx, idcard); e != nil {
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return nil, e
		}
		if e != gredis.Nil && e != ecode.ErrCacheDataConflict {
			log.Error(ctx, "[dao.GetUserByIDCard] redis op failed", log.String("idcard", idcard), log.CError(e))
		}
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
				log.Error(nil, "[dao.GetUserByIDCard] db op failed", log.String("idcard", idcard), log.CError(e))
				if e == ecode.ErrDBDataConflict {
					time.Sleep(time.Millisecond * 5)
					continue
				}
				if e == ecode.ErrUserNotExist {
					//set redis empty key
					go func() {
						if e := d.RedisSetUserIDCardIndex(context.Background(), idcard, ""); e != nil {
							log.Error(nil, "[dao.GetUserByIDCard] update redis failed", log.String("idcard", idcard), log.CError(e))
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
				log.Error(nil, "[dao.GetUserByIDCard] update redis failed", log.String("user_id", user.UserID.Hex()), log.CError(e))
			}
		}()
		go func() {
			if e := d.RedisSetUserIDCardIndex(context.Background(), user.IDCard, user.UserID.Hex()); e != nil {
				log.Error(nil, "[dao.GetUserByIDCard] update redis failed", log.String("idcard", user.IDCard), log.CError(e))
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
		if e != gredis.Nil && e != ecode.ErrCacheDataConflict {
			log.Error(ctx, "[dao.GetUserByNickName] redis op failed", log.String("nick_name", nickname), log.CError(e))
		}
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
				log.Error(nil, "[dao.GetUserByNickName] db op failed", log.String("nick_name", nickname), log.CError(e))
				if e == ecode.ErrDBDataConflict {
					time.Sleep(time.Millisecond * 5)
					continue
				}
				if e == ecode.ErrUserNotExist {
					//set redis empty key
					go func() {
						if e := d.RedisSetUserNickNameIndex(context.Background(), nickname, ""); e != nil {
							log.Error(nil, "[dao.GetUserByNickName] update redis failed", log.String("nick_name", nickname), log.CError(e))
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
				log.Error(nil, "[dao.GetUserByNickName] update redis failed", log.String("user_id", user.UserID.Hex()), log.CError(e))
			}
		}()
		go func() {
			if e := d.RedisSetUserNickNameIndex(context.Background(), user.NickName, user.UserID.Hex()); e != nil {
				log.Error(nil, "[dao.GetUserByNickName] update redis failed", log.String("nick_name", user.NickName), log.CError(e))
			}
		}()
		return unsafe.Pointer(user), e
	})
	return (*model.User)(unsafeUser), e
}

// oauth index
func (d *Dao) GetUserOAuthIndex(ctx context.Context, oauthservicename, oauthid string) (string, error) {
	if userid, e := d.RedisGetUserOAuthIndex(ctx, oauthservicename, oauthid); e != nil {
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return userid, e
		}
		if e != gredis.Nil {
			log.Error(ctx, "[dao.GetUserOAuthIndex] redis op failed", log.String(oauthservicename, oauthid), log.CError(e))
		}
	} else if userid != "" {
		return userid, nil
	}
	unsafeUserid, e := oneshot.Do("GetUserOAuthIndex_"+oauthservicename+"|"+oauthid, func() (unsafe.Pointer, error) {
		index, e := d.MongoGetUserOAuthIndex(ctx, oauthservicename, oauthid)
		if e != nil {
			log.Error(nil, "[dao.GetUserOAuthIndex] db op failed", log.String(oauthservicename, oauthid), log.CError(e))
			if e == ecode.ErrUserNotExist {
				//set redis empty key
				go func() {
					if e := d.RedisSetUserOAuthIndex(context.Background(), oauthservicename, oauthid, ""); e != nil {
						log.Error(nil, "[dao.GetUserOAuthIndex] update redis failed", log.String(oauthservicename, oauthid), log.CError(e))
					}
				}()
			}
			return nil, e
		}
		userid := index.UserID.Hex()
		//update redis
		go func() {
			if e := d.RedisSetUserOAuthIndex(context.Background(), oauthservicename, oauthid, userid); e != nil {
				log.Error(ctx, "[dao.GetUserOAuthIndex] update redis failed", log.String(oauthservicename, oauthid), log.CError(e))
			}
		}()
		return unsafe.Pointer(&userid), nil
	})
	if e != nil {
		return "", e
	}
	return *(*string)(unsafeUserid), e
}

// tel index
func (d *Dao) GetUserTelIndex(ctx context.Context, tel string) (string, error) {
	if userid, e := d.RedisGetUserTelIndex(ctx, tel); e != nil {
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return userid, e
		}
		if e != gredis.Nil {
			log.Error(ctx, "[dao.GetUserTelIndex] redis op failed", log.String("tel", tel), log.CError(e))
		}
	} else if userid != "" {
		return userid, nil
	}
	//redis error or redis not exist,we need to query db
	unsafeUserid, e := oneshot.Do("GetUserTelIndex_"+tel, func() (unsafe.Pointer, error) {
		index, e := d.MongoGetUserTelIndex(ctx, tel)
		if e != nil {
			log.Error(nil, "[dao.GetUserTelIndex] db op failed", log.String("tel", tel), log.CError(e))
			if e == ecode.ErrUserNotExist {
				//set redis empty key
				go func() {
					if e := d.RedisSetUserTelIndex(context.Background(), tel, ""); e != nil {
						log.Error(nil, "[dao.GetUserTelIndex] update redis failed", log.String("tel", tel), log.CError(e))
					}
				}()
			}
			return nil, e
		}
		userid := index.UserID.Hex()
		//update redis
		go func() {
			if e := d.RedisSetUserTelIndex(context.Background(), tel, userid); e != nil {
				log.Error(ctx, "[dao.GetUserTelIndex] update redis failed", log.String("tel", tel), log.CError(e))
			}
		}()
		return unsafe.Pointer(&userid), nil
	})
	if e != nil {
		return "", e
	}
	return *(*string)(unsafeUserid), e
}

// email index
func (d *Dao) GetUserEmailIndex(ctx context.Context, email string) (string, error) {
	if userid, e := d.RedisGetUserEmailIndex(ctx, email); e != nil {
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return userid, e
		}
		if e != gredis.Nil {
			log.Error(ctx, "[dao.GetUserEmailIndex] redis op failed", log.String("email", email), log.CError(e))
		}
	} else if userid != "" {
		return userid, nil
	}
	//redis error or redis not exist,we need to query db
	unsafeUserid, e := oneshot.Do("GetUserEmailIndex_"+email, func() (unsafe.Pointer, error) {
		index, e := d.MongoGetUserEmailIndex(ctx, email)
		if e != nil {
			log.Error(nil, "[dao.GetUserEmailIndex] db op failed", log.String("email", email), log.CError(e))
			if e == ecode.ErrUserNotExist {
				//set redis empty key
				go func() {
					if e := d.RedisSetUserEmailIndex(context.Background(), email, ""); e != nil {
						log.Error(nil, "[dao.GetUserEmailIndex] update redis failed", log.String("email", email), log.CError(e))
					}
				}()
			}
			return nil, e
		}
		userid := index.UserID.Hex()
		//update redis
		go func() {
			if e := d.RedisSetUserEmailIndex(context.Background(), email, userid); e != nil {
				log.Error(ctx, "[dao.GetUserEmailIndex] update redis failed", log.String("email", email), log.CError(e))
			}
		}()
		return unsafe.Pointer(&userid), nil
	})
	if e != nil {
		return "", e
	}
	return *(*string)(unsafeUserid), e
}

// idcard index
func (d *Dao) GetUserIDCardIndex(ctx context.Context, idcard string) (string, error) {
	if userid, e := d.RedisGetUserIDCardIndex(ctx, idcard); e != nil {
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return userid, e
		}
		if e != gredis.Nil {
			log.Error(ctx, "[dao.GetUserIDCardIndex] redis op failed", log.String("idcard", idcard), log.CError(e))
		}
	} else if userid != "" {
		return userid, nil
	}
	//redis error or redis not exist,we need to query db
	unsafeUserid, e := oneshot.Do("GetUserIDcardIndex_"+idcard, func() (unsafe.Pointer, error) {
		index, e := d.MongoGetUserIDCardIndex(ctx, idcard)
		if e != nil {
			log.Error(nil, "[dao.GetUserIDCardIndex] db op failed", log.String("idcard", idcard), log.CError(e))
			if e == ecode.ErrUserNotExist {
				//set redis empty key
				go func() {
					if e := d.RedisSetUserIDCardIndex(context.Background(), idcard, ""); e != nil {
						log.Error(nil, "[dao.GetUserIDCardIndex] update redis failed", log.String("idcard", idcard))
					}
				}()
			}
			return nil, e
		}
		userid := index.UserID.Hex()
		//update redis
		go func() {
			if e := d.RedisSetUserIDCardIndex(context.Background(), idcard, userid); e != nil {
				log.Error(ctx, "[dao.GetUserIDCardIndex] update redis failed", log.String("idcard", idcard), log.CError(e))
			}
		}()
		return unsafe.Pointer(&userid), nil
	})
	if e != nil {
		return "", e
	}
	return *(*string)(unsafeUserid), e
}

// nickname index
func (d *Dao) GetUserNickNameIndex(ctx context.Context, nickname string) (string, error) {
	if userid, e := d.RedisGetUserNickNameIndex(ctx, nickname); e != nil {
		if e == ecode.ErrUserNotExist {
			//key exist but value is empty
			return userid, e
		}
		if e != gredis.Nil {
			log.Error(ctx, "[dao.GetUserNickNameIndex] redis op failed", log.String("nick_name", nickname), log.CError(e))
		}
	} else if userid != "" {
		return userid, nil
	}
	//redis error or redis not exist,we need to query db
	unsafeUserid, e := oneshot.Do("GetUserNickNameIndex_"+nickname, func() (unsafe.Pointer, error) {
		index, e := d.MongoGetUserNickNameIndex(ctx, nickname)
		if e != nil {
			log.Error(nil, "[dao.GetUserNickNameIndex] db op failed", log.String("nick_name", nickname), log.CError(e))
			if e == ecode.ErrUserNotExist {
				//set redis empty key
				go func() {
					if e := d.RedisSetUserNickNameIndex(context.Background(), nickname, ""); e != nil {
						log.Error(nil, "[dao.GetUserNickNameIndex] update redis failed", log.String("nick_name", nickname), log.CError(e))
					}
				}()
			}
			return nil, e
		}
		userid := index.UserID.Hex()
		//update redis
		go func() {
			if e := d.RedisSetUserNickNameIndex(context.Background(), nickname, userid); e != nil {
				log.Error(ctx, "[dao.GetUserNickNameIndex] update redis failed", log.String("nick_name", nickname), log.CError(e))
			}
		}()
		return unsafe.Pointer(&userid), nil
	})
	if e != nil {
		return "", e
	}
	return *(*string)(unsafeUserid), e
}
