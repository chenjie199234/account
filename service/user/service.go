package user

import (
	"context"

	"github.com/chenjie199234/account/api"
	"github.com/chenjie199234/account/config"
	userdao "github.com/chenjie199234/account/dao/user"
	"github.com/chenjie199234/account/ecode"
	"github.com/chenjie199234/account/model"
	"github.com/chenjie199234/account/util"

	// "github.com/chenjie199234/Corelib/cgrpc"
	// "github.com/chenjie199234/Corelib/crpc"
	// "github.com/chenjie199234/Corelib/log"
	// "github.com/chenjie199234/Corelib/web"
	"github.com/chenjie199234/Corelib/cerror"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/metadata"
	publicmids "github.com/chenjie199234/Corelib/mids"
	"github.com/chenjie199234/Corelib/util/graceful"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service subservice for user business
type Service struct {
	stop *graceful.Graceful

	userDao *userdao.Dao
}

// Start -
func Start() *Service {
	return &Service{
		stop: graceful.New(),

		userDao: userdao.NewDao(config.GetSql("user_sql"), config.GetRedis("user_redis"), config.GetMongo("user_mongo")),
	}
}
func (s *Service) GetUserInfo(ctx context.Context, req *api.GetUserInfoReq) (*api.GetUserInfoResp, error) {
	var user *model.User
	var e error
	switch req.SrcType {
	case "user_id":
		userid, e := primitive.ObjectIDFromHex(req.Src)
		if e != nil {
			log.Error(ctx, "[GetUserInfo] user_id format wrong", map[string]interface{}{"user_id": req.Src, "error": e})
			return nil, ecode.ErrReq
		}
		user, e = s.userDao.MongoGetUserByUserID(ctx, userid)
		if e != nil {
			log.Error(ctx, "[GetUserInfo] db op failed", map[string]interface{}{"user_id": req.Src, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	case "tel":
		user, e = s.userDao.MongoGetUserByTel(ctx, req.Src)
		if e != nil {
			log.Error(ctx, "[GetUserInfo] db op failed", map[string]interface{}{"tel": req.Src, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	case "email":
		user, e = s.userDao.MongoGetUserByEmail(ctx, req.Src)
		if e != nil {
			log.Error(ctx, "[GetUserInfo] db op failed", map[string]interface{}{"email": req.Src, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	case "idcard":
		user, e = s.userDao.MongoGetUserByIDCard(ctx, req.Src)
		if e != nil {
			log.Error(ctx, "[GetUserInfo] db op failed", map[string]interface{}{"idcard": req.Src, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	case "nickname":
		user, e = s.userDao.MongoGetUserByNickName(ctx, req.Src)
		if e != nil {
			log.Error(ctx, "[GetUserInfo] db op failed", map[string]interface{}{"nick_name": req.Src, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	}
	return &api.GetUserInfoResp{
		Info: &api.UserInfo{
			UserId:   req.Src,
			Idcard:   user.IDCard,
			Tel:      user.Tel,
			Email:    user.Email,
			NickName: user.NickName,
			Money:    user.Money,
			Ctime:    user.UserID.Timestamp().Unix(),
		},
	}, nil
}
func (s *Service) Login(ctx context.Context, req *api.LoginReq) (*api.LoginResp, error) {
	if req.PasswordType == "static" && req.Password == "" {
		log.Error(ctx, "[Login] empty static password", map[string]interface{}{"src_type": req.SrcType, "src": req.Src})
		return nil, ecode.ErrReq
	}
	var user *model.User
	var e error
	switch req.SrcType {
	case "idcard":
		if req.PasswordType == "dynamic" {
			log.Error(ctx, "[Login] idcard can't use dynamic password", nil)
			return nil, ecode.ErrReq
		}
		//static
		if user, e = s.userDao.MongoGetUserByIDCard(ctx, req.Src); e != nil {
			log.Error(ctx, "[Login] db op failed", map[string]interface{}{"idcard": req.Src, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	case "nickname":
		if req.PasswordType == "dynamic" {
			log.Error(ctx, "[Login] nickname can't use dynamic password", nil)
			return nil, ecode.ErrReq
		}
		//static
		if user, e = s.userDao.MongoGetUserByNickName(ctx, req.Src); e != nil {
			log.Error(ctx, "[Login] db op failed", map[string]interface{}{"nickname": req.Src, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	case "tel":
		if req.PasswordType == "dynamic" && req.Password == "" {
			//send code
			//set redis and send tel is async
			//if set redis success and send tel failed
			//we need to clean the redis
			if !s.stop.AddOne() {
				return nil, cerror.ErrServerClosing
			}
			code := util.MakeRandCode()
			if e := s.userDao.RedisSetCode(ctx, req.Src, userdao.LoginTel, code); e != nil {
				log.Error(ctx, "[Login] redis op failed", map[string]interface{}{"tel": req.Src, "error": e})
				s.stop.DoneOne()
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			if e := util.SendTelCode(ctx, code); e != nil {
				log.Error(ctx, "[Login] send tel failed", map[string]interface{}{"tel": req.Src, "error": e})
				//clean redis code
				if e := s.userDao.RedisDelCode(ctx, req.Src, userdao.LoginTel); e != nil {
					log.Error(ctx, "[Login] del redis code failed", map[string]interface{}{"tel": req.Src, "error": e})
					go func() {
						if e := s.userDao.RedisDelCode(context.Background(), req.Src, userdao.LoginTel); e != nil {
							log.Error(ctx, "[Login] del redis code failed in goroutine", map[string]interface{}{"tel": req.Src, "error": e})
						}
						s.stop.DoneOne()
					}()
				} else {
					s.stop.DoneOne()
				}
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			s.stop.DoneOne()
			return &api.LoginResp{Step: "verify"}, nil
		} else if req.PasswordType == "dynamic" {
			//verify code
			rest, e := s.userDao.RedisCheckCode(ctx, req.Src, userdao.LoginTel, req.Password)
			if e != nil {
				log.Error(ctx, "[Login] redis op failed", map[string]interface{}{"tel": req.Src, "code": req.Password, "error": e})
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			if rest > 0 {
				log.Error(ctx, "[Login] check failed", map[string]interface{}{"tel": req.Src, "code": req.Password, "rest": rest})
				return nil, ecode.ErrPasswordWrong
			} else if rest == 0 {
				log.Error(ctx, "[Login] all check times failed", map[string]interface{}{"tel": req.Src, "max_checktimes": userdao.DefaultCheckTimes, "ban_seconds": userdao.DefaultExpireSeconds})
				return nil, ecode.ErrBan
			}
			//verify success
		}
		//static or dynamic's verify success
		user, e = s.userDao.MongoGetUserByTel(ctx, req.Src)
		if e != nil {
			if req.PasswordType == "static" || e != ecode.ErrUserNotExist {
				log.Error(ctx, "[Login] db op failed", map[string]interface{}{"tel": req.Src, "error": e})
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			if user, e = s.userDao.MongoCreateUserByTel(ctx, req.Src); e != nil {
				log.Error(ctx, "[Login] db op failed", map[string]interface{}{"tel": req.Src, "error": e})
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
		}
	case "email":
		if req.PasswordType == "dynamic" && req.Password == "" {
			//send code
			//set redis and send email is async
			//if set redis success and send email failed
			//we need to clean the redis
			if !s.stop.AddOne() {
				return nil, cerror.ErrServerClosing
			}
			code := util.MakeRandCode()
			if e := s.userDao.RedisSetCode(ctx, req.Src, userdao.LoginEmail, code); e != nil {
				log.Error(ctx, "[Login] redis op failed", map[string]interface{}{"email": req.Src, "error": e})
				s.stop.DoneOne()
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			if e := util.SendEmailCode(ctx, code); e != nil {
				log.Error(ctx, "[Login] send email failed", map[string]interface{}{"email": req.Src, "error": e})
				//clean redis code
				if e := s.userDao.RedisDelCode(ctx, req.Src, userdao.LoginEmail); e != nil {
					log.Error(ctx, "[Login] del redis code failed", map[string]interface{}{"email": req.Src, "error": e})
					go func() {
						if e := s.userDao.RedisDelCode(context.Background(), req.Src, userdao.LoginEmail); e != nil {
							log.Error(ctx, "[Login] del redis code failed in goroutine", map[string]interface{}{"email": req.Src, "error": e})
						}
						s.stop.DoneOne()
					}()
				} else {
					s.stop.DoneOne()
				}
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			s.stop.DoneOne()
			return &api.LoginResp{Step: "verify"}, nil
		} else if req.PasswordType == "dynamic" {
			//verify code
			rest, e := s.userDao.RedisCheckCode(ctx, req.Src, userdao.LoginEmail, req.Password)
			if e != nil {
				log.Error(ctx, "[Login] redis op failed", map[string]interface{}{"email": req.Src, "code": req.Password, "error": e})
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			if rest > 0 {
				log.Error(ctx, "[Login] check failed", map[string]interface{}{"email": req.Src, "code": req.Password, "rest": rest})
				return nil, ecode.ErrPasswordWrong
			} else if rest == 0 {
				log.Error(ctx, "[Login] all check times failed", map[string]interface{}{"email": req.Src, "max_checktimes": userdao.DefaultCheckTimes, "ban_seconds": userdao.DefaultExpireSeconds})
				return nil, ecode.ErrBan
			}
			//verify success
		}
		//static or dynamic's verify success
		user, e = s.userDao.MongoGetUserByEmail(ctx, req.Src)
		if e != nil {
			if req.PasswordType == "static" || e != ecode.ErrUserNotExist {
				log.Error(ctx, "[Login] db op failed", map[string]interface{}{"email": req.Src, "error": e})
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			if user, e = s.userDao.MongoCreateUserByEmail(ctx, req.Src); e != nil {
				log.Error(ctx, "[Login] db op failed", map[string]interface{}{"email": req.Src, "error": e})
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
		}
	}
	needSetPassword := false
	if req.PasswordType == "dynamic" {
		needSetPassword = util.SignCheck("", user.Password) == nil
	} else if e := util.SignCheck(req.Password, user.Password); e != nil {
		if e == ecode.ErrSignCheckFailed {
			e = ecode.ErrPasswordWrong
		}
		log.Error(ctx, "[Login] sign check failed", map[string]interface{}{"src_type": req.SrcType, "src": req.Src, "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	//TODO set the puber
	token := publicmids.MakeToken(ctx, "", *config.EC.DeployEnv, *config.EC.RunEnv, user.UserID.Hex())
	resp := &api.LoginResp{
		Token: token,
		Info: &api.UserInfo{
			UserId:   user.UserID.Hex(),
			Idcard:   util.MaskIDCard(user.IDCard),
			Tel:      util.MaskTel(user.Tel),
			Email:    util.MaskEmail(user.Email),
			NickName: user.NickName,
			Ctime:    user.UserID.Timestamp().Unix(),
			Money:    user.Money,
		},
		Step: "success",
	}
	if needSetPassword {
		resp.Step = "password"
	}
	return resp, nil
}
func (s *Service) UpdateStaticPassword(ctx context.Context, req *api.UpdateStaticPasswordReq) (*api.UpdateStaticPasswordResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[UpdateStaticPassword] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}
	if e := s.userDao.MongoUpdateUserPassword(ctx, operator, req.OldStaticPassword, req.NewStaticPassword); e != nil {
		log.Error(ctx, "[UpdateStaticPassword] db op failed", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.UpdateStaticPasswordResp{}, nil
}
func (s *Service) UpdateNickName(ctx context.Context, req *api.UpdateNickNameReq) (*api.UpdateNickNameResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[UpdateNickName] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}
	if e := s.userDao.MongoUpdateUserNickName(ctx, operator, req.NewNickName); e != nil {
		log.Error(ctx, "[UpdateNickName] db op failed", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.UpdateNickNameResp{}, nil
}

// UpdateTel Step 1:send dynamic password to old email or tail
// UpdateTel Step 2:verify old email's or tel's dynamic password and send dynamic password to new email
// UpdateTel Step 3:verify new email's dynamic and update
func (s *Service) UpdateEmail(ctx context.Context, req *api.UpdateEmailReq) (*api.UpdateEmailResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[UpdateEmail] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}
}

// UpdateTel Step 1:send dynamic password to old email or tail
// UpdateTel Step 2:verify old email's or tel's dynamic password and send dynamic password to new tel
// UpdateTel Step 3:verify new tel's dynamic and update
func (s *Service) UpdateTel(ctx context.Context, req *api.UpdateTelReq) (*api.UpdateTelResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[UpdateTel] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}
	if req.NewEmailDynamicPassword != "" {

	}
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
