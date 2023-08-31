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
			if rest, e := s.userDao.RedisSetCode(ctx, req.Src, util.LoginTel, code); e != nil {
				s.stop.DoneOne()
				if e != ecode.ErrCodeAlreadySend {
					log.Error(ctx, "[Login] redis op failed", map[string]interface{}{"tel": req.Src, "error": e})
					return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
				}
				if rest != 0 {
					return &api.LoginResp{Step: "verify"}, nil
				}
				log.Error(ctx, "[Login] all check times failed", map[string]interface{}{"tel": req.Src, "ban_seconds": userdao.DefaultExpireSeconds})
				return nil, ecode.ErrBan
			}
			if e := util.SendTelCode(ctx, req.Src, code, util.LoginTel); e != nil {
				log.Error(ctx, "[Login] send tel failed", map[string]interface{}{"tel": req.Src, "error": e})
				//clean redis code
				if e := s.userDao.RedisDelCode(ctx, req.Src, util.LoginTel); e != nil {
					log.Error(ctx, "[Login] del redis code failed", map[string]interface{}{"tel": req.Src, "error": e})
					go func() {
						if e := s.userDao.RedisDelCode(context.Background(), req.Src, util.LoginTel); e != nil {
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
			rest, e := s.userDao.RedisCheckCode(ctx, req.Src, util.LoginTel, req.Password)
			if e != nil {
				log.Error(ctx, "[Login] redis op failed", map[string]interface{}{"tel": req.Src, "code": req.Password, "error": e})
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			if rest > 0 {
				log.Error(ctx, "[Login] check failed", map[string]interface{}{"tel": req.Src, "code": req.Password, "rest": rest})
				return nil, ecode.ErrPasswordWrong
			} else if rest == 0 {
				log.Error(ctx, "[Login] all check times failed", map[string]interface{}{"tel": req.Src, "ban_seconds": userdao.DefaultExpireSeconds})
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
			if rest, e := s.userDao.RedisSetCode(ctx, req.Src, util.LoginEmail, code); e != nil {
				s.stop.DoneOne()
				if e != ecode.ErrCodeAlreadySend {
					log.Error(ctx, "[Login] redis op failed", map[string]interface{}{"email": req.Src, "error": e})
					return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
				}
				if rest != 0 {
					return &api.LoginResp{Step: "verify"}, nil
				}
				log.Error(ctx, "[Login] all check times failed", map[string]interface{}{"email": req.Src, "ban_seconds": userdao.DefaultExpireSeconds})
				return nil, ecode.ErrBan
			}
			if e := util.SendEmailCode(ctx, req.Src, code, util.LoginEmail); e != nil {
				log.Error(ctx, "[Login] send email failed", map[string]interface{}{"email": req.Src, "error": e})
				//clean redis code
				if e := s.userDao.RedisDelCode(ctx, req.Src, util.LoginEmail); e != nil {
					log.Error(ctx, "[Login] del redis code failed", map[string]interface{}{"email": req.Src, "error": e})
					go func() {
						if e := s.userDao.RedisDelCode(context.Background(), req.Src, util.LoginEmail); e != nil {
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
			rest, e := s.userDao.RedisCheckCode(ctx, req.Src, util.LoginEmail, req.Password)
			if e != nil {
				log.Error(ctx, "[Login] redis op failed", map[string]interface{}{"email": req.Src, "code": req.Password, "error": e})
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			if rest > 0 {
				log.Error(ctx, "[Login] check failed", map[string]interface{}{"email": req.Src, "code": req.Password, "rest": rest})
				return nil, ecode.ErrPasswordWrong
			} else if rest == 0 {
				log.Error(ctx, "[Login] all check times failed", map[string]interface{}{"email": req.Src, "ban_seconds": userdao.DefaultExpireSeconds})
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
	log.Info(ctx, "[Login] success", map[string]interface{}{"operator": user.UserID.Hex()})
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
	log.Info(ctx, "[UpdateStaticPassword] success", map[string]interface{}{"operator": md["Token-Data"]})
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
	log.Info(ctx, "[UpdateNickName] success", map[string]interface{}{"operator": md["Token-Data"], "new_nick_name": req.NewNickName})
	return &api.UpdateNickNameResp{}, nil
}

// UpdateTel Step 1:send dynamic password to old email or tel
// UpdateTel Step 2:verify old email's or tel's dynamic password and send dynamic password to new email
// UpdateTel Step 3:verify new email's dynamic and update
func (s *Service) UpdateEmail(ctx context.Context, req *api.UpdateEmailReq) (*api.UpdateEmailResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[UpdateEmail] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}
	if req.NewEmailDynamicPassword != "" {
		//step 3
		rest, e := s.userDao.RedisCheckCode(ctx, md["Token-Data"]+"_"+req.NewEmail, util.NewEmail, req.NewEmailDynamicPassword)
		if e != nil {
			log.Error(ctx, "[UpdateEmail] redis op failed", map[string]interface{}{"operator": md["Token-Data"], "code": req.NewEmailDynamicPassword, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if rest > 0 {
			log.Error(ctx, "[UpdateEmail] check failed", map[string]interface{}{"operator": md["Token-Data"], "code": req.NewEmailDynamicPassword, "rest": rest})
			return nil, ecode.ErrPasswordWrong
		} else if rest == 0 {
			log.Error(ctx, "[UpdateEmail] all check times failed", map[string]interface{}{"operator": md["Token-Data"], "ban_seconds": userdao.DefaultExpireSeconds})
			return nil, ecode.ErrBan
		}
		//verify success

		if e := s.userDao.MongoUpdateUserEmail(ctx, operator, req.NewEmail); e != nil {
			log.Error(ctx, "[UpdateEmail] db op failed", map[string]interface{}{"operator": md["Token-Data"], "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		log.Info(ctx, "[UpdateEmail] success", map[string]interface{}{"operator": md["Token-Data"], "new_email": req.NewEmail})
		return &api.UpdateEmailResp{Step: "success"}, nil
	} else if req.OldDynamicPassword != "" {
		//step 2
		var rest int
		var e error
		switch req.OldReceiverType {
		case "tel":
			rest, e = s.userDao.RedisCheckCode(ctx, md["Token-Data"], util.OldTel, req.OldDynamicPassword)
		case "email":
			rest, e = s.userDao.RedisCheckCode(ctx, md["Token-Data"], util.OldEmail, req.OldDynamicPassword)
		}
		if e != nil {
			log.Error(ctx, "[UpdateEmail] redis op failed", map[string]interface{}{"operator": md["Token-Data"], "code": req.OldDynamicPassword, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if rest > 0 {
			log.Error(ctx, "[UpdateEmail] check failed", map[string]interface{}{"operator": md["Token-Data"], "code": req.OldDynamicPassword, "rest": rest})
			return nil, ecode.ErrPasswordWrong
		} else if rest == 0 {
			log.Error(ctx, "[UpdateEmail] all check times failed", map[string]interface{}{"operator": md["Token-Data"], "ban_seconds": userdao.DefaultExpireSeconds})
			return nil, ecode.ErrBan
		}
		//verify success

		//send code to new email
		//set redis and send email is async
		//if set redis success and send email failed
		//we need to clean the redis
		if !s.stop.AddOne() {
			return nil, cerror.ErrServerClosing
		}
		code := util.MakeRandCode()
		if rest, e := s.userDao.RedisSetCode(ctx, md["Token-Data"]+"_"+req.NewEmail, util.NewEmail, code); e != nil {
			s.stop.DoneOne()
			if e != ecode.ErrCodeAlreadySend {
				log.Error(ctx, "[UpdateEmail] redis op failed", map[string]interface{}{"operator": md["Token-Data"], "error": e})
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			if rest != 0 {
				return &api.UpdateEmailResp{Step: "newverify"}, nil
			}
			log.Error(ctx, "[UpdateEmail] all check times failed", map[string]interface{}{"operator": md["Token-Data"], "ban_seconds": userdao.DefaultExpireSeconds})
			return nil, ecode.ErrBan
		}
		if e := util.SendEmailCode(ctx, req.NewEmail, code, util.NewEmail); e != nil {
			log.Error(ctx, "[UpdateEmail] send email failed", map[string]interface{}{"operator": md["Token-Data"], "error": e})
			if e := s.userDao.RedisDelCode(ctx, md["Token-Data"]+"_"+req.NewEmail, util.NewEmail); e != nil {
				log.Error(ctx, "[UpdateEmail] del redis code failed", map[string]interface{}{"operator": md["Token-Data"], "error": e})
				go func() {
					if e := s.userDao.RedisDelCode(context.Background(), md["Token-Data"]+"_"+req.NewEmail, util.NewEmail); e != nil {
						log.Error(ctx, "[UpdateEmail] del redis code failed in goroutine", map[string]interface{}{"operator": md["Token-Data"], "error": e})
					}
					s.stop.DoneOne()
				}()
			} else {
				s.stop.DoneOne()
			}
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		log.Info(ctx, "[UpdateEmail] send dynamic password to new email success", map[string]interface{}{"operator": md["Token-Data"], "new_email": req.NewEmail, "code": code})
		s.stop.DoneOne()
		return &api.UpdateEmailResp{Step: "newverify"}, nil
	}
	//step 1
	if _, rest, e := s.userDao.RedisGetCode(ctx, md["Token-Data"]+"_"+req.NewEmail, util.NewEmail); e != nil {
		if e != ecode.ErrCodeAlreadyExpire {
			log.Error(ctx, "[UpdateEmail] redis op failed", map[string]interface{}{"operator": md["Token-Data"], "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if rest == 0 {
			log.Error(ctx, "[UpdateEmail] all check times failed", map[string]interface{}{"operator": md["Token-Data"], "max_checktimes": userdao.DefaultCheckTimes, "ban_seconds": userdao.DefaultExpireSeconds})
		}
	}
	user, e := s.userDao.MongoGetUserByUserID(ctx, operator)
	if e != nil {
		log.Error(ctx, "[UpdateEmail] db op failed", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if req.OldReceiverType == "tel" && user.Tel == "" {
		log.Error(ctx, "[UpdateEmail] missing tel,can't use tel to receive dynamic password", map[string]interface{}{"operator": md["Token-Data"]})
		return nil, ecode.ErrReq
	}
	if req.OldReceiverType == "email" && user.Email == "" {
		log.Error(ctx, "[UpdateEmail] missing email,can't use email to receive dynamic password", map[string]interface{}{"operator": md["Token-Data"]})
		return nil, ecode.ErrReq
	}

	//send code
	//set redis and send is async
	//if set redis success and send failed
	//we need to clean the redis
	if !s.stop.AddOne() {
		return nil, cerror.ErrServerClosing
	}
	code := util.MakeRandCode()
	var rest int
	switch req.OldReceiverType {
	case "tel":
		rest, e = s.userDao.RedisSetCode(ctx, md["Token-Data"], util.OldTel, code)
	case "email":
		rest, e = s.userDao.RedisSetCode(ctx, md["Token-Data"], util.OldEmail, code)
	}
	if e != nil {
		s.stop.DoneOne()
		if e != ecode.ErrCodeAlreadySend {
			log.Error(ctx, "[UpdateEmail] redis op failed", map[string]interface{}{"operator": md["Token-Data"], "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if rest != 0 {
			return &api.UpdateEmailResp{Step: "oldverify"}, nil
		}
		log.Error(ctx, "[UpdateEmail] all check times failed", map[string]interface{}{"operator": md["Token-Data"], "ban_seconds": userdao.DefaultExpireSeconds})
		return nil, ecode.ErrBan
	}
	switch req.OldReceiverType {
	case "tel":
		if e = util.SendTelCode(ctx, user.Tel, code, util.OldTel); e != nil {
			log.Error(ctx, "[UpdateEmail] send tel failed", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		}
	case "email":
		if e = util.SendEmailCode(ctx, user.Email, code, util.OldEmail); e != nil {
			log.Error(ctx, "[UpdateEmail] send email failed", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		}
	}
	if e != nil {
		var ee error
		switch req.OldReceiverType {
		case "tel":
			ee = s.userDao.RedisDelCode(ctx, md["Token-Data"], util.OldTel)
		case "email":
			ee = s.userDao.RedisDelCode(ctx, md["Token-Data"], util.OldEmail)
		}
		if ee != nil {
			log.Error(ctx, "[UpdateEmail] del redis code failed", map[string]interface{}{"operator": md["Token-Data"], "error": ee})
			go func() {
				var e error
				switch req.OldReceiverType {
				case "tel":
					e = s.userDao.RedisDelCode(context.Background(), md["Token-Data"], util.OldTel)
				case "email":
					e = s.userDao.RedisDelCode(context.Background(), md["Token-Data"], util.OldEmail)
				}
				if e != nil {
					log.Error(ctx, "[UpdateEmail] del redis code failed in goroutine", map[string]interface{}{"operator": md["Token-Data"], "error": e})
				}
				s.stop.DoneOne()
			}()
		} else {
			s.stop.DoneOne()
		}
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	switch req.OldReceiverType {
	case "tel":
		log.Info(ctx, "[UpdateEmail] send dynamic password to old tel success", map[string]interface{}{"operator": md["Token-Data"], "old_tel": user.Tel, "code": code})
	case "email":
		log.Info(ctx, "[UpdateEmail] send dynamic password to old email success", map[string]interface{}{"operator": md["Token-Data"], "old_email": user.Email, "code": code})
	}
	s.stop.DoneOne()
	return &api.UpdateEmailResp{Step: "oldverify"}, nil
}

// UpdateTel Step 1:send dynamic password to old email or tel
// UpdateTel Step 2:verify old email's or tel's dynamic password and send dynamic password to new tel
// UpdateTel Step 3:verify new tel's dynamic and update
func (s *Service) UpdateTel(ctx context.Context, req *api.UpdateTelReq) (*api.UpdateTelResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[UpdateTel] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}
	if req.NewTelDynamicPassword != "" {
		//step 3
		rest, e := s.userDao.RedisCheckCode(ctx, md["Token-Data"]+"_"+req.NewTel, util.NewTel, req.NewTelDynamicPassword)
		if e != nil {
			log.Error(ctx, "[UpdateTel] redis op failed", map[string]interface{}{"operator": md["Token-Data"], "code": req.NewTelDynamicPassword, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if rest > 0 {
			log.Error(ctx, "[UpdateTel] check failed", map[string]interface{}{"operator": md["Token-Data"], "code": req.NewTelDynamicPassword, "rest": rest})
			return nil, ecode.ErrPasswordWrong
		} else if rest == 0 {
			log.Error(ctx, "[UpdateTel] all check times failed", map[string]interface{}{"operator": md["Token-Data"], "ban_seconds": userdao.DefaultExpireSeconds})
			return nil, ecode.ErrBan
		}
		//verify success

		if e := s.userDao.MongoUpdateUserTel(ctx, operator, req.NewTel); e != nil {
			log.Error(ctx, "[UpdateTel] db op failed", map[string]interface{}{"operator": md["Token-Data"], "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		log.Info(ctx, "[UpdateTel] success", map[string]interface{}{"operator": md["Token-Data"], "new_tel": req.NewTel})
		return &api.UpdateTelResp{Step: "success"}, nil
	} else if req.OldDynamicPassword != "" {
		//step 2
		var rest int
		var e error
		switch req.OldReceiverType {
		case "tel":
			rest, e = s.userDao.RedisCheckCode(ctx, md["Token-Data"], util.OldTel, req.OldDynamicPassword)
		case "email":
			rest, e = s.userDao.RedisCheckCode(ctx, md["Token-Data"], util.OldEmail, req.OldDynamicPassword)
		}
		if e != nil {
			log.Error(ctx, "[UpdateTel] redis op failed", map[string]interface{}{"operator": md["Token-Data"], "code": req.OldDynamicPassword, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if rest > 0 {
			log.Error(ctx, "[UpdateTel] check failed", map[string]interface{}{"operator": md["Token-Data"], "code": req.OldDynamicPassword, "rest": rest})
			return nil, ecode.ErrPasswordWrong
		} else if rest == 0 {
			log.Error(ctx, "[UpdateTel] all check times failed", map[string]interface{}{"operator": md["Token-Data"], "ban_seconds": userdao.DefaultExpireSeconds})
			return nil, ecode.ErrBan
		}
		//verify success

		//send code
		//set redis and send is async
		//if set redis success and send failed
		//we need to clean the redis
		if !s.stop.AddOne() {
			return nil, cerror.ErrServerClosing
		}
		code := util.MakeRandCode()
		if rest, e := s.userDao.RedisSetCode(ctx, md["Token-Data"]+"_"+req.NewTel, util.NewTel, code); e != nil {
			s.stop.DoneOne()
			if e != ecode.ErrCodeAlreadySend {
				log.Error(ctx, "[UpdateTel] redis op failed", map[string]interface{}{"operator": md["Token-Data"], "error": e})
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			if rest != 0 {
				return &api.UpdateTelResp{Step: "newverify"}, nil
			}
			log.Error(ctx, "[UpdateTel] all check times failed", map[string]interface{}{"operator": md["Token-Data"], "ban_seconds": userdao.DefaultExpireSeconds})
			return nil, ecode.ErrBan
		}
		if e := util.SendTelCode(ctx, req.NewTel, code, util.NewTel); e != nil {
			log.Error(ctx, "[UpdateTel] send tel failed", map[string]interface{}{"operator": md["Token-Data"], "error": e})
			if e := s.userDao.RedisDelCode(ctx, md["Token-Data"]+"_"+req.NewTel, util.NewTel); e != nil {
				log.Error(ctx, "[UpdateTel] del redis code failed", map[string]interface{}{"operator": md["Token-Data"], "error": e})
				go func() {
					if e := s.userDao.RedisDelCode(context.Background(), md["Token-Data"]+"_"+req.NewTel, util.NewTel); e != nil {
						log.Error(ctx, "[UpdateTel] del redis code failed in goroutine", map[string]interface{}{"operator": md["Token-Data"], "error": e})
					}
					s.stop.DoneOne()
				}()
			} else {
				s.stop.DoneOne()
			}
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		log.Info(ctx, "[UpdateTel] send dynamic password to new tel success", map[string]interface{}{"operator": md["Token-Data"], "new_tel": req.NewTel, "code": code})
		s.stop.DoneOne()
		return &api.UpdateTelResp{Step: "newverify"}, nil
	}
	//step 1
	if _, rest, e := s.userDao.RedisGetCode(ctx, md["Token-Data"]+"_"+req.NewTel, util.NewTel); e != nil {
		if e != ecode.ErrCodeAlreadyExpire {
			log.Error(ctx, "[UpdateTel] redis op failed", map[string]interface{}{"operator": md["Token-Data"], "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if rest == 0 {
			log.Error(ctx, "[UpdateTel] all check times failed", map[string]interface{}{"operator": md["Token-Data"], "max_checktimes": userdao.DefaultCheckTimes, "ban_seconds": userdao.DefaultExpireSeconds})
		}
	}
	user, e := s.userDao.MongoGetUserByUserID(ctx, operator)
	if e != nil {
		log.Error(ctx, "[UpdateTel] db op failed", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if req.OldReceiverType == "tel" && user.Tel == "" {
		log.Error(ctx, "[UpdateTel] missing tel,can't use tel to receive dynamic password", map[string]interface{}{"operator": md["Token-Data"]})
		return nil, ecode.ErrReq
	}
	if req.OldReceiverType == "email" && user.Email == "" {
		log.Error(ctx, "[UpdateTel] missing email,can't use email to receive dynamic password", map[string]interface{}{"operator": md["Token-Data"]})
		return nil, ecode.ErrReq
	}

	//send code
	//set redis and send is async
	//if set redis success and send failed
	//we need to clean the redis
	if !s.stop.AddOne() {
		return nil, cerror.ErrServerClosing
	}
	code := util.MakeRandCode()
	var rest int
	switch req.OldReceiverType {
	case "tel":
		rest, e = s.userDao.RedisSetCode(ctx, md["Token-Data"], util.OldTel, code)
	case "email":
		rest, e = s.userDao.RedisSetCode(ctx, md["Token-Data"], util.OldEmail, code)
	}
	if e != nil {
		s.stop.DoneOne()
		if e != ecode.ErrCodeAlreadySend {
			log.Error(ctx, "[UpdateTel] redis op failed", map[string]interface{}{"operator": md["Token-Data"], "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if rest != 0 {
			return &api.UpdateTelResp{Step: "oldverify"}, nil
		}
		log.Error(ctx, "[UpdateTel] all check times failed", map[string]interface{}{"operator": md["Token-Data"], "ban_seconds": userdao.DefaultExpireSeconds})
		return nil, ecode.ErrBan
	}
	switch req.OldReceiverType {
	case "tel":
		if e = util.SendTelCode(ctx, user.Tel, code, util.OldTel); e != nil {
			log.Error(ctx, "[UpdateTel] send tel failed", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		}
	case "email":
		if e = util.SendEmailCode(ctx, user.Email, code, util.OldEmail); e != nil {
			log.Error(ctx, "[UpdateTel] send email failed", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		}
	}
	if e != nil {
		var ee error
		switch req.OldReceiverType {
		case "tel":
			ee = s.userDao.RedisDelCode(ctx, md["Token-Data"], util.OldTel)
		case "email":
			ee = s.userDao.RedisDelCode(ctx, md["Token-Data"], util.OldEmail)
		}
		if ee != nil {
			log.Error(ctx, "[UpdateTel] del redis code failed", map[string]interface{}{"operator": md["Token-Data"], "error": ee})
			go func() {
				var e error
				switch req.OldReceiverType {
				case "tel":
					e = s.userDao.RedisDelCode(context.Background(), md["Token-Data"], util.OldTel)
				case "email":
					e = s.userDao.RedisDelCode(context.Background(), md["Token-Data"], util.OldEmail)
				}
				if e != nil {
					log.Error(ctx, "[UpdateTel] del redis code failed in goroutine", map[string]interface{}{"operator": md["Token-Data"], "error": e})
				}
				s.stop.DoneOne()
			}()
		} else {
			s.stop.DoneOne()
		}
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	switch req.OldReceiverType {
	case "tel":
		log.Info(ctx, "[UpdateTel] send dynamic password to old tel success", map[string]interface{}{"operator": md["Token-Data"], "old_tel": user.Tel, "code": code})
	case "email":
		log.Info(ctx, "[UpdateTel] send dynamic password to old email success", map[string]interface{}{"operator": md["Token-Data"], "old_email": user.Email, "code": code})
	}
	s.stop.DoneOne()
	return &api.UpdateTelResp{Step: "oldverify"}, nil
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
