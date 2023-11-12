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

		userDao: userdao.NewDao(nil, config.GetRedis("account_redis"), config.GetMongo("account_mongo")),
	}
}

// srctype: email/tel
// src: email addr/tel number
func (s *Service) sendcode(ctx context.Context, callerName, srctype, src, operator, action string) error {
	//send code
	//set redis and send is async
	//if set redis success and send failed
	//we need to clean the redis
	e := s.stop.Add(1)
	if e != nil {
		if e == graceful.ErrClosing {
			return cerror.ErrServerClosing
		}
		return ecode.ErrBusy
	}
	code, dup, e := s.userDao.RedisSetCode(ctx, operator, action, src)
	if e != nil {
		s.stop.DoneOne()
		log.Error(ctx, "["+callerName+"] redis op failed", log.String("operator", operator), log.String(srctype, src), log.CError(e))
		return e
	}
	if dup {
		//if tel's or email's code already send,we jump to verify step
		return nil
	}
	switch action {
	case util.Login:
		e = s.userDao.RedisLockLoginDynamic(ctx, operator)
	case util.DelEmail:
		fallthrough
	case util.UpdateEmailStep1:
		e = s.userDao.RedisLockEmailOP(ctx, operator)
	case util.UpdateEmailStep2:
		//this is controled by step1
	case util.DelTel:
		fallthrough
	case util.UpdateTelStep1:
		e = s.userDao.RedisLockTelOP(ctx, operator)
	case util.UpdateTelStep2:
		//this is controled by step2
	case util.DelIDCard:
		fallthrough
	case util.UpdateIDCard:
		e = s.userDao.RedisLockIDCardOP(ctx, operator)
	case util.DelNickName:
		fallthrough
	case util.UpdateNickName:
		e = s.userDao.RedisLockNickNameOP(ctx, operator)
	default:
		return ecode.ErrUnknownAction
	}
	if e != nil {
		log.Error(ctx, "["+callerName+"] rate check failed", log.String("operator", operator), log.String(srctype, src), log.CError(e))
	} else if srctype == "email" {
		if e = util.SendEmailCode(ctx, src, code, action); e != nil {
			log.Error(ctx, "["+callerName+"] send dynamic password failed", log.String("operator", operator), log.String(srctype, src), log.CError(e))
		}
	} else if e = util.SendTelCode(ctx, src, code, action); e != nil {
		log.Error(ctx, "["+callerName+"] send dynamic password failed", log.String("operator", operator), log.String(srctype, src), log.CError(e))
	}
	if e == nil {
		log.Info(ctx, "["+callerName+"] send dynamic password success", log.String("operator", operator), log.String(srctype, src), log.String("code", code))
		s.stop.DoneOne()
		return nil
	}
	//if rate check failed or send failed,clean redis code
	if ee := s.userDao.RedisDelCode(ctx, operator, action); ee != nil {
		go func() {
			if ee := s.userDao.RedisDelCode(context.Background(), operator, action); ee != nil {
				log.Error(ctx, "["+callerName+"] del redis code failed", log.String("operator", operator), log.String(srctype, src), log.CError(ee))
			}
			s.stop.DoneOne()
		}()
	} else {
		s.stop.DoneOne()
	}
	//this e is SendEmailCode's or SendTelCode's or RateCheck's
	//not the RedisDelCode's
	return e
}

func (s *Service) GetUserInfo(ctx context.Context, req *api.GetUserInfoReq) (*api.GetUserInfoResp, error) {
	var user *model.User
	switch req.SrcType {
	case "user_id":
		userid, e := primitive.ObjectIDFromHex(req.Src)
		if e != nil {
			log.Error(ctx, "[GetUserInfo] user_id format wrong", log.String("user_id", req.Src), log.CError(e))
			return nil, ecode.ErrReq
		}
		if user, e = s.userDao.GetUser(ctx, userid); e != nil {
			log.Error(ctx, "[GetUserInfo] dao op failed", log.String("user_id", req.Src), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	case "tel":
		var e error
		if user, e = s.userDao.GetUserByTel(ctx, req.Src); e != nil {
			log.Error(ctx, "[GetUserInfo] dao op failed", log.String("tel", req.Src), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	case "email":
		var e error
		if user, e = s.userDao.GetUserByEmail(ctx, req.Src); e != nil {
			log.Error(ctx, "[GetUserInfo] dao op failed", log.String("email", req.Src), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	case "idcard":
		var e error
		if user, e = s.userDao.GetUserByIDCard(ctx, req.Src); e != nil {
			log.Error(ctx, "[GetUserInfo] dao op failed", log.String("idcard", req.Src), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	case "nick_name":
		var e error
		if user, e = s.userDao.GetUserByNickName(ctx, req.Src); e != nil {
			log.Error(ctx, "[GetUserInfo] dao op failed", log.String("nick_name", req.Src), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	}
	return &api.GetUserInfoResp{
		Info: &api.UserInfo{
			UserId:   user.UserID.Hex(),
			Idcard:   user.IDCard,
			Tel:      user.Tel,
			Email:    user.Email,
			NickName: user.NickName,
			Money:    user.Money,
			Ctime:    uint32(user.UserID.Timestamp().Unix()),
		},
	}, nil
}
func (s *Service) Login(ctx context.Context, req *api.LoginReq) (*api.LoginResp, error) {
	if req.PasswordType == "static" && req.Password == "" {
		log.Error(ctx, "[Login] empty static password", log.String(req.SrcType, req.Src))
		return nil, ecode.ErrReq
	}
	var user *model.User
	var e error
	switch req.SrcType {
	case "idcard":
		if req.PasswordType == "dynamic" {
			log.Error(ctx, "[Login] idcard can't use dynamic password")
			return nil, ecode.ErrReq
		}
		//static
		if user, e = s.userDao.GetUserByIDCard(ctx, req.Src); e != nil {
			log.Error(ctx, "[Login] dao op failed", log.String("idcard", req.Src), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	case "nick_name":
		if req.PasswordType == "dynamic" {
			log.Error(ctx, "[Login] nick_name can't use dynamic password")
			return nil, ecode.ErrReq
		}
		//static
		if user, e = s.userDao.GetUserByNickName(ctx, req.Src); e != nil {
			log.Error(ctx, "[Login] dao op failed", log.String("nick_name", req.Src), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	case "tel":
		if req.PasswordType == "static" {
			if user, e = s.userDao.GetUserByTel(ctx, req.Src); e != nil {
				log.Error(ctx, "[Login] dao op failed", log.String("tel", req.Src), log.CError(e))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
		} else if req.Password == "" {
			if e = s.sendcode(ctx, "Login", req.SrcType, req.Src, req.Src, util.Login); e != nil {
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			return &api.LoginResp{Step: "verify"}, nil
		} else {
			if e = s.userDao.RedisCheckCode(ctx, req.Src, util.Login, req.Password, ""); e != nil {
				log.Error(ctx, "[Login] redis op failed", log.String("tel", req.Src), log.String("code", req.Password), log.CError(e))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			if user, e = s.userDao.GetOrCreateUserByTel(ctx, req.Src); e != nil {
				log.Error(ctx, "[Login] dao op failed", log.String("tel", req.Src), log.CError(e))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
		}
	case "email":
		if req.PasswordType == "static" {
			if user, e = s.userDao.GetUserByEmail(ctx, req.Src); e != nil {
				log.Error(ctx, "[Login] dao op failed", log.String("email", req.Src), log.CError(e))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
		} else if req.Password == "" {
			if e := s.sendcode(ctx, "Login", req.SrcType, req.Src, req.Src, util.Login); e != nil {
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			return &api.LoginResp{Step: "verify"}, nil
		} else {
			if e = s.userDao.RedisCheckCode(ctx, req.Src, util.Login, req.Password, ""); e != nil {
				log.Error(ctx, "[Login] redis op failed", log.String("email", req.Src), log.String("code", req.Password), log.CError(e))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			if user, e = s.userDao.GetOrCreateUserByEmail(ctx, req.Src); e != nil {
				log.Error(ctx, "[Login] dao op failed", log.String("email", req.Src), log.CError(e))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
		}
	}
	if req.PasswordType == "static" {
		if e := util.SignCheck(req.Password, user.Password); e != nil {
			if e == ecode.ErrSignCheckFailed {
				e = ecode.ErrPasswordWrong
			}
			log.Error(ctx, "[Login] sign check failed", log.String(req.SrcType, req.Src), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	}
	//TODO set the puber
	token := publicmids.MakeToken(ctx, "", *config.EC.DeployEnv, *config.EC.RunEnv, user.UserID.Hex(), "")
	resp := &api.LoginResp{
		Token: token,
		Info: &api.UserInfo{
			UserId:   user.UserID.Hex(),
			Idcard:   util.MaskIDCard(user.IDCard),
			Tel:      util.MaskTel(user.Tel),
			Email:    util.MaskEmail(user.Email),
			NickName: user.NickName,
			Ctime:    uint32(user.UserID.Timestamp().Unix()),
			Money:    user.Money,
		},
		Step: "success",
	}
	if req.PasswordType == "dynamic" && util.SignCheck("", user.Password) == nil {
		//this is a new account
		resp.Step = "password"
	}
	log.Info(ctx, "[Login] success", log.String("operator", user.UserID.Hex()))
	return resp, nil
}
func (s *Service) SelfUserInfo(ctx context.Context, req *api.SelfUserInfoReq) (*api.SelfUserInfoResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[SelfUserInfo] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}
	user, e := s.userDao.GetUser(ctx, operator)
	if e != nil {
		log.Error(ctx, "[SelfUserInfo] dao op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.SelfUserInfoResp{
		Info: &api.UserInfo{
			UserId:   user.UserID.Hex(),
			Idcard:   util.MaskIDCard(user.IDCard),
			Tel:      util.MaskTel(user.Tel),
			Email:    util.MaskEmail(user.Email),
			NickName: user.NickName,
			Ctime:    uint32(user.UserID.Timestamp().Unix()),
			Money:    user.Money,
		},
	}, nil
}
func (s *Service) UpdateStaticPassword(ctx context.Context, req *api.UpdateStaticPasswordReq) (*api.UpdateStaticPasswordResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[UpdateStaticPassword] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}
	//update db and clean redis is async
	//the service's rolling update may happened between update db and clean redis
	//so we need to make this not happened
	if e := s.stop.Add(1); e != nil {
		if e == graceful.ErrClosing {
			return nil, cerror.ErrServerClosing
		}
		return nil, ecode.ErrBusy
	}
	//redis lock
	if e := s.userDao.RedisLockUpdatePassword(ctx, md["Token-User"]); e != nil {
		s.stop.DoneOne()
		log.Error(ctx, "[UpdateStaticPassword] redis op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}

	if e := s.userDao.MongoUpdateUserPassword(ctx, operator, req.OldStaticPassword, req.NewStaticPassword); e != nil {
		s.stop.DoneOne()
		log.Error(ctx, "[UpdateStaticPassword] db op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[UpdateStaticPassword] success", log.String("operator", md["Token-User"]))
	go func() {
		if e := s.userDao.RedisDelUser(context.Background(), md["Token-User"]); e != nil {
			log.Error(ctx, "[UpdateStaticPassword] clean redis failed", log.String("operator", md["Token-User"]), log.CError(e))
		}
		s.stop.DoneOne()
	}()
	return &api.UpdateStaticPasswordResp{}, nil
}

func (s *Service) NickNameDuplicateCheck(ctx context.Context, req *api.NickNameDuplicateCheckReq) (*api.NickNameDuplicateCheckResp, error) {
	md := metadata.GetMetadata(ctx)
	//redis lock
	if e := s.userDao.RedisLockDuplicateCheck(ctx, "nick_name", md["Token-User"]); e != nil {
		log.Error(ctx, "[NickNameDuplicateCheck] redis op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	userid, e := s.userDao.GetUserNickNameIndex(ctx, req.NickName)
	if e != nil && e != ecode.ErrUserNotExist {
		log.Error(ctx, "[NickNameDuplicateCheck] dao op failed", log.String("operator", md["Token-User"]), log.String("nick_name", req.NickName), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.NickNameDuplicateCheckResp{Duplicate: userid != ""}, nil
}

// UpdateNickName By OAuth
//
//	Step1:verify oauth belong to this account
//
// UpdateNickName By Dynamic Password
//
//	Step1:send dynamic password to email or tel
//	Step2:verify email's or tel's dynamic password
func (s *Service) UpdateNickName(ctx context.Context, req *api.UpdateNickNameReq) (*api.UpdateNickNameResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[UpdateNickName] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}
	update := func() error {
		//update db and clean redis is async
		//the service's rolling update may happened between update db and clean redis
		//so we need to make this not happened
		if e := s.stop.Add(3); e != nil {
			if e == graceful.ErrClosing {
				return cerror.ErrServerClosing
			}
			return ecode.ErrBusy
		}

		var olduser *model.User
		if olduser, e = s.userDao.MongoUpdateUserNickName(ctx, operator, req.NewNickName); e != nil {
			s.stop.DoneOne()
			s.stop.DoneOne()
			s.stop.DoneOne()
			log.Error(ctx, "[UpdateNickName] db op failed", log.String("operator", md["Token-user"]), log.CError(e))
			return e
		}
		log.Info(ctx, "[UpdateNickName] success", log.String("operator", md["Token-User"]), log.String("new_nick_name", req.NewNickName))
		if olduser.NickName != req.NewNickName {
			go func() {
				if e := s.userDao.RedisDelUser(context.Background(), md["Token-User"]); e != nil {
					log.Error(ctx, "[UpdateNickName] clean redis failed", log.String("operator", md["Token-User"]), log.CError(e))
				}
				s.stop.DoneOne()
			}()
			go func() {
				if olduser.NickName != "" {
					if e := s.userDao.RedisDelUserNickNameIndex(context.Background(), olduser.NickName); e != nil {
						log.Error(ctx, "[UpdateNickName] clean redis failed", log.String("nick_name", olduser.NickName), log.CError(e))
					}
				}
				s.stop.DoneOne()
			}()
			go func() {
				if req.NewNickName != "" {
					if e := s.userDao.RedisDelUserNickNameIndex(context.Background(), req.NewNickName); e != nil {
						log.Error(ctx, "[UpdateNickName] clean redis failed", log.String("nick_name", req.NewNickName), log.CError(e))
					}
				}
				s.stop.DoneOne()
			}()
		} else {
			s.stop.DoneOne()
			s.stop.DoneOne()
			s.stop.DoneOne()
		}
		return nil
	}
	if req.SrcType == "oauth" {
		if req.SrcTypeExtra == "" || req.DynamicPassword == "" {
			return nil, ecode.ErrReq
		}
		if e := s.userDao.RedisLockNickNameOP(ctx, md["Token-User"]); e != nil {
			log.Error(ctx, "[UpdateNickName] rate check failed", log.String("operator", md["Token-User"]), log.String(req.SrcTypeExtra, req.DynamicPassword), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		oauthid, e := util.OAuthVerifyCode(ctx, "UpdateNickName", req.SrcTypeExtra, req.DynamicPassword)
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		user, e := s.userDao.GetUserByOAuth(ctx, req.SrcTypeExtra, oauthid)
		if e != nil {
			log.Error(ctx, "[UpdateNickName] dao op failed", log.String("operator", md["Token-User"]), log.String(req.SrcTypeExtra, oauthid), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.UserID.Hex() != md["Token-User"] {
			log.Error(ctx, "[UpdateNickName] this is not the required oauth", log.String("operator", md["Token-User"]), log.String(req.SrcTypeExtra, oauthid))
			return nil, ecode.ErrOAuthWrong
		}
		//verify success
		if e := update(); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.UpdateNickNameResp{Step: "success"}, nil
	}

	if req.DynamicPassword != "" {
		//step2
		if e := s.userDao.RedisCheckCode(ctx, md["Token-User"], util.UpdateNickName, req.DynamicPassword, ""); e != nil {
			log.Error(ctx, "[UpdateNickName] redis op failed", log.String("operator", md["Token-User"]), log.String("code", req.DynamicPassword), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		//verify success
		if e := update(); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.UpdateNickNameResp{Step: "success"}, nil
	}
	//step1

	user, e := s.userDao.GetUser(ctx, operator)
	if e != nil {
		log.Error(ctx, "[UpdateNickName] dao op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if user.NickName == req.NewNickName {
		return &api.UpdateNickNameResp{Step: "success"}, nil
	}

	if req.SrcType == "tel" && user.Tel == "" {
		log.Error(ctx, "[UpdateNickName] missing tel,can't use tel to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.SrcType == "email" && user.Email == "" {
		log.Error(ctx, "[UpdateNickName] missing email,can't use email to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}

	if req.SrcType == "email" {
		e = s.sendcode(ctx, "UpdateNickName", req.SrcType, user.Email, md["Token-User"], util.UpdateNickName)
	} else {
		e = s.sendcode(ctx, "UpdateNickName", req.SrcType, user.Tel, md["Token-User"], util.UpdateNickName)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if req.SrcType == "email" {
		return &api.UpdateNickNameResp{Step: "oldverify", Receiver: util.MaskEmail(user.Email)}, nil
	}
	return &api.UpdateNickNameResp{Step: "oldverify", Receiver: util.MaskTel(user.Tel)}, nil
}

// DelNickName By OAuth
//
//	Step1:verify oauth belong to this account
//
// DelNickName By Dynamic Password
//
//	Step1:send dynamic password to email or tel
//	Step2:verify email's or tel's dynamic password
func (s *Service) DelNickName(ctx context.Context, req *api.DelNickNameReq) (*api.DelNickNameResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[DelNickName] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}
	update := func() (bool, error) {
		//update db and clean redis is async
		//the service's rolling update may happened between update db and clean redis
		//so we need to make this not happened
		if e := s.stop.Add(2); e != nil {
			if e == graceful.ErrClosing {
				return false, cerror.ErrServerClosing
			}
			return false, ecode.ErrBusy
		}
		var olduser *model.User
		if olduser, e = s.userDao.MongoUpdateUserNickName(ctx, operator, ""); e != nil {
			s.stop.DoneOne()
			s.stop.DoneOne()
			log.Error(ctx, "[DelNickName] db op failed", log.String("operator", md["Token-user"]), log.CError(e))
			return false, e
		}
		log.Info(ctx, "[DelNickName] success", log.String("operator", md["Token-User"]))
		if olduser.NickName != "" {
			go func() {
				if e := s.userDao.RedisDelUser(context.Background(), md["Token-User"]); e != nil {
					log.Error(ctx, "[DelNickName] clean redis failed", log.String("operator", md["Token-User"]), log.CError(e))
				}
				s.stop.DoneOne()
			}()
			go func() {
				if e := s.userDao.RedisDelUserNickNameIndex(context.Background(), olduser.NickName); e != nil {
					log.Error(ctx, "[DelNickName] clean redis failed", log.String("nick_name", olduser.NickName), log.CError(e))
				}
				s.stop.DoneOne()
			}()
		} else {
			s.stop.DoneOne()
			s.stop.DoneOne()
		}
		var final bool
		if olduser.Email == "" && olduser.IDCard == "" && olduser.Tel == "" && len(olduser.OAuths) == 0 {
			final = true
		}
		return final, nil
	}
	if req.SrcType == "oauth" {
		if req.SrcTypeExtra == "" || req.DynamicPassword == "" {
			return nil, ecode.ErrReq
		}
		if e := s.userDao.RedisLockNickNameOP(ctx, md["Token-User"]); e != nil {
			log.Error(ctx, "[DelNickName] rate check failed", log.String("operator", md["Token-User"]), log.String(req.SrcTypeExtra, req.DynamicPassword), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		oauthid, e := util.OAuthVerifyCode(ctx, "DelNickName", req.SrcTypeExtra, req.DynamicPassword)
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		user, e := s.userDao.GetUserByOAuth(ctx, req.SrcTypeExtra, oauthid)
		if e != nil {
			log.Error(ctx, "[DelNickName] dao op failed", log.String("operator", md["Token-User"]), log.String(req.SrcTypeExtra, oauthid), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.UserID.Hex() != md["Token-User"] {
			log.Error(ctx, "[DelNickName] this is not the required oauth", log.String("operator", md["Token-User"]), log.String(req.SrcTypeExtra, oauthid))
			return nil, ecode.ErrOAuthWrong
		}
		//verify success
		final, e := update()
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.DelNickNameResp{Step: "success", Final: final}, nil
	}
	if req.DynamicPassword != "" {
		//step2
		if e := s.userDao.RedisCheckCode(ctx, md["Token-User"], util.DelNickName, req.DynamicPassword, ""); e != nil {
			log.Error(ctx, "[DelNickName] redis op failed", log.String("operator", md["Token-User"]), log.String("code", req.DynamicPassword), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		//verify success
		final, e := update()
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.DelNickNameResp{Step: "success", Final: final}, nil
	}
	//step1
	user, e := s.userDao.GetUser(ctx, operator)
	if e != nil {
		log.Error(ctx, "[DelNickName] dao op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	var final bool
	if user.Email == "" && user.IDCard == "" && user.Tel == "" && len(user.OAuths) == 0 {
		final = true
	}
	if user.NickName == "" {
		return &api.DelNickNameResp{Step: "success", Final: final}, nil
	}

	if req.SrcType == "tel" && user.Tel == "" {
		log.Error(ctx, "[DelNickName] missing tel,can't use tel to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.SrcType == "email" && user.Email == "" {
		log.Error(ctx, "[DelNickName] missing email,can't use email to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.SrcType == "email" {
		e = s.sendcode(ctx, "DelNickName", req.SrcType, user.Email, md["Token-User"], util.DelNickName)
	} else {
		e = s.sendcode(ctx, "DelNickName", req.SrcType, user.Tel, md["Token-User"], util.DelNickName)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if req.SrcType == "email" {
		return &api.DelNickNameResp{Step: "oldverify", Final: final, Receiver: util.MaskEmail(user.Email)}, nil
	}
	return &api.DelNickNameResp{Step: "oldverify", Final: final, Receiver: util.MaskTel(user.Tel)}, nil
}

func (s *Service) IdcardDuplicateCheck(ctx context.Context, req *api.IdcardDuplicateCheckReq) (*api.IdcardDuplicateCheckResp, error) {
	md := metadata.GetMetadata(ctx)
	//redis lock
	if e := s.userDao.RedisLockDuplicateCheck(ctx, "idcard", md["Token-User"]); e != nil {
		log.Error(ctx, "[IdcardDuplicateCheck] redis op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	userid, e := s.userDao.GetUserIDCardIndex(ctx, req.Idcard)
	if e != nil && e != ecode.ErrUserNotExist {
		log.Error(ctx, "[IdcardDuplicateCheck] dao op failed", log.String("operator", md["Token-User"]), log.String("idcard", req.Idcard), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.IdcardDuplicateCheckResp{Duplicate: userid != ""}, nil
}

// UpdateIdcard By OAuth
//
//	Step1:verify oauth belong to this account
//
// UpdateIdCard By Dynamic Password
//
//	Step1:send dynamic password to email or tel
//	Step2:verify email's or tel's dynamic password
func (s *Service) UpdateIdcard(ctx context.Context, req *api.UpdateIdcardReq) (*api.UpdateIdcardResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[UpdateIdcard] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}
	update := func() error {
		//update db and clean redis is async
		//the service's rolling update may happened between update db and clean redis
		//so we need to make this not happened
		if e := s.stop.Add(3); e != nil {
			if e == graceful.ErrClosing {
				return cerror.ErrServerClosing
			}
			return ecode.ErrBusy
		}

		var olduser *model.User
		if olduser, e = s.userDao.MongoUpdateUserIDCard(ctx, operator, req.NewIdcard); e != nil {
			s.stop.DoneOne()
			s.stop.DoneOne()
			s.stop.DoneOne()
			log.Error(ctx, "[UpdateIdcard] db op failed", log.String("operator", md["Token-User"]), log.CError(e))
			return e
		}
		log.Info(ctx, "[UpdateIdcard] success", log.String("operator", md["Token-User"]), log.String("new_idcard", req.NewIdcard))
		if olduser.IDCard != req.NewIdcard {
			go func() {
				if e := s.userDao.RedisDelUser(context.Background(), md["Token-User"]); e != nil {
					log.Error(ctx, "[UpdateIdcard] clean redis failed", log.String("operator", md["Token-User"]), log.CError(e))
				}
				s.stop.DoneOne()
			}()
			go func() {
				if olduser.IDCard != "" {
					if e := s.userDao.RedisDelUserIDCardIndex(context.Background(), olduser.IDCard); e != nil {
						log.Error(ctx, "[UpdateIdcard] clean redis failed", log.String("idcard", olduser.IDCard), log.CError(e))
					}
				}
				s.stop.DoneOne()
			}()
			go func() {
				if req.NewIdcard != "" {
					if e := s.userDao.RedisDelUserIDCardIndex(context.Background(), req.NewIdcard); e != nil {
						log.Error(ctx, "[UpdateIdcard] clean redis failed", log.String("idcard", req.NewIdcard), log.CError(e))
					}
				}
				s.stop.DoneOne()
			}()
		} else {
			s.stop.DoneOne()
			s.stop.DoneOne()
			s.stop.DoneOne()
		}
		return nil
	}
	if req.SrcType == "oauth" {
		if req.SrcTypeExtra == "" || req.DynamicPassword == "" {
			return nil, ecode.ErrReq
		}
		if e := s.userDao.RedisLockIDCardOP(ctx, md["Token-User"]); e != nil {
			log.Error(ctx, "[UpdateIdcard] rate check failed", log.String("operator", md["Token-User"]), log.String(req.SrcTypeExtra, req.DynamicPassword), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		oauthid, e := util.OAuthVerifyCode(ctx, "UpdateIdcard", req.SrcTypeExtra, req.DynamicPassword)
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		user, e := s.userDao.GetUserByOAuth(ctx, req.SrcTypeExtra, oauthid)
		if e != nil {
			log.Error(ctx, "[UpdateIdcard] dao op failed", log.String("operator", md["Token-User"]), log.String(req.SrcTypeExtra, oauthid), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.UserID.Hex() != md["Token-User"] {
			log.Error(ctx, "[UpdateIdcard] this is not the required oauth", log.String("operator", md["Token-User"]), log.String(req.SrcTypeExtra, oauthid))
			return nil, ecode.ErrOAuthWrong
		}
		//verify success
		if e := update(); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.UpdateIdcardResp{Step: "success"}, nil
	}
	if req.DynamicPassword != "" {
		//step 2
		if e := s.userDao.RedisCheckCode(ctx, md["Token-User"], util.UpdateIDCard, req.DynamicPassword, ""); e != nil {
			log.Error(ctx, "[UpdateIdcard] redis op failed", log.String("operator", md["Token-User"]), log.String("code", req.DynamicPassword), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		//verify success
		if e := update(); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.UpdateIdcardResp{Step: "success"}, nil
	}
	//step 1
	user, e := s.userDao.GetUser(ctx, operator)
	if e != nil {
		log.Error(ctx, "[UpdateIdcard],dao op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if user.IDCard == req.NewIdcard {
		return &api.UpdateIdcardResp{Step: "success"}, nil
	}

	if req.SrcType == "tel" && user.Tel == "" {
		log.Error(ctx, "[UpdateIdcard] missing tel,can't use tel to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.SrcType == "email" && user.Email == "" {
		log.Error(ctx, "[UpdateIdcard] missing email,can't use email to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}

	if req.SrcType == "email" {
		e = s.sendcode(ctx, "UpdateIdcard", req.SrcType, user.Email, md["Token-User"], util.UpdateIDCard)
	} else {
		e = s.sendcode(ctx, "UpdateIdcard", req.SrcType, user.Tel, md["Token-User"], util.UpdateIDCard)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if req.SrcType == "email" {
		return &api.UpdateIdcardResp{Step: "oldverify", Receiver: util.MaskEmail(user.Email)}, nil
	}
	return &api.UpdateIdcardResp{Step: "oldverify", Receiver: util.MaskTel(user.Tel)}, nil
}

// DelIdcard By OAuth
//
//	Step1:verify oauth belong to this account
//
// DelIdcard By Dynamic Password
//
//	Step1:send dynamic password to email or tel
//	Step2:verify email's or tel's dynamic password
func (s *Service) DelIdcard(ctx context.Context, req *api.DelIdcardReq) (*api.DelIdcardResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[DelIdcard] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}
	update := func() (bool, error) {
		//update db and clean redis is async
		//the service's rolling update may happened between update db and clean redis
		//so we need to make this not happened
		if e := s.stop.Add(2); e != nil {
			if e == graceful.ErrClosing {
				return false, cerror.ErrServerClosing
			}
			return false, ecode.ErrBusy
		}
		var olduser *model.User
		if olduser, e = s.userDao.MongoUpdateUserIDCard(ctx, operator, ""); e != nil {
			s.stop.DoneOne()
			s.stop.DoneOne()
			log.Error(ctx, "[DelIdcard] db op failed", log.String("operator", md["Token-user"]), log.CError(e))
			return false, e
		}
		log.Info(ctx, "[DelIdcard] success", log.String("operator", md["Token-User"]))
		if olduser.IDCard != "" {
			go func() {
				if e := s.userDao.RedisDelUser(context.Background(), md["Token-User"]); e != nil {
					log.Error(ctx, "[DelIdcard] clean redis failed", log.String("operator", md["Token-User"]), log.CError(e))
				}
				s.stop.DoneOne()
			}()
			go func() {
				if e := s.userDao.RedisDelUserIDCardIndex(context.Background(), olduser.IDCard); e != nil {
					log.Error(ctx, "[DelIdcard] clean redis failed", log.String("idcard", olduser.IDCard), log.CError(e))
				}
				s.stop.DoneOne()
			}()
		} else {
			s.stop.DoneOne()
			s.stop.DoneOne()
		}
		var final bool
		if olduser.Email == "" && olduser.NickName == "" && olduser.Tel == "" && len(olduser.OAuths) == 0 {
			final = true
		}
		return final, nil
	}
	if req.SrcType == "oauth" {
		if req.SrcTypeExtra == "" || req.DynamicPassword == "" {
			return nil, ecode.ErrReq
		}
		if e := s.userDao.RedisLockIDCardOP(ctx, md["Token-User"]); e != nil {
			log.Error(ctx, "[DelIDCard] rate check failed", log.String("operator", md["Token-User"]), log.String(req.SrcTypeExtra, req.DynamicPassword), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		oauthid, e := util.OAuthVerifyCode(ctx, "DelIdcard", req.SrcTypeExtra, req.DynamicPassword)
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		user, e := s.userDao.GetUserByOAuth(ctx, req.SrcTypeExtra, oauthid)
		if e != nil {
			log.Error(ctx, "[DelIdcard] dao op failed", log.String("operator", md["Token-User"]), log.String(req.SrcTypeExtra, oauthid), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.UserID.Hex() != md["Token-User"] {
			log.Error(ctx, "[DelIdcard] this is not the required oauth", log.String("operator", md["Token-User"]), log.String(req.SrcTypeExtra, oauthid))
			return nil, ecode.ErrOAuthWrong
		}
		//verify success
		final, e := update()
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.DelIdcardResp{Step: "success", Final: final}, nil
	}
	if req.DynamicPassword != "" {
		//step2
		if e := s.userDao.RedisCheckCode(ctx, md["Token-User"], util.DelIDCard, req.DynamicPassword, ""); e != nil {
			log.Error(ctx, "[DelIdcard] redis op failed", log.String("operator", md["Token-User"]), log.String("code", req.DynamicPassword), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		//verify success
		final, e := update()
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.DelIdcardResp{Step: "success", Final: final}, nil
	}
	//step1
	user, e := s.userDao.GetUser(ctx, operator)
	if e != nil {
		log.Error(ctx, "[DelIdcard] dao op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	var final bool
	if user.Email == "" && user.NickName == "" && user.Tel == "" && len(user.OAuths) == 0 {
		final = true
	}
	if user.IDCard == "" {
		return &api.DelIdcardResp{Step: "success", Final: final}, nil
	}

	if req.SrcType == "tel" && user.Tel == "" {
		log.Error(ctx, "[DelIdcard] missing tel,can't use tel to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.SrcType == "email" && user.Email == "" {
		log.Error(ctx, "[DelIdcard] missing email,can't use email to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.SrcType == "email" {
		e = s.sendcode(ctx, "DelIdcard", req.SrcType, user.Email, md["Token-User"], util.DelIDCard)
	} else {
		e = s.sendcode(ctx, "DelIdcard", req.SrcType, user.Tel, md["Token-User"], util.DelIDCard)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if req.SrcType == "email" {
		return &api.DelIdcardResp{Step: "oldverify", Final: final, Receiver: util.MaskEmail(user.Email)}, nil
	}
	return &api.DelIdcardResp{Step: "oldverify", Final: final, Receiver: util.MaskTel(user.Tel)}, nil
}

func (s *Service) EmailDuplicateCheck(ctx context.Context, req *api.EmailDuplicateCheckReq) (*api.EmailDuplicateCheckResp, error) {
	md := metadata.GetMetadata(ctx)
	//redis lock
	if e := s.userDao.RedisLockDuplicateCheck(ctx, "email", md["Token-User"]); e != nil {
		log.Error(ctx, "[EmailDuplicateCheck] redis op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	userid, e := s.userDao.GetUserEmailIndex(ctx, req.Email)
	if e != nil && e != ecode.ErrUserNotExist {
		log.Error(ctx, "[EmailDuplicateCheck] dao op failed", log.String("operator", md["Token-User"]), log.String("email", req.Email), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.EmailDuplicateCheckResp{Duplicate: userid != ""}, nil
}

// TODO oauth
// UpdateTel Step 1:send dynamic password to old email or tel
// UpdateTel Step 2:verify old email's or tel's dynamic password and send dynamic password to new email
// UpdateTel Step 3:verify new email's dynamic and update
func (s *Service) UpdateEmail(ctx context.Context, req *api.UpdateEmailReq) (*api.UpdateEmailResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[UpdateEmail] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}
	if req.NewEmailDynamicPassword != "" {
		//step 3
		if e := s.userDao.RedisCheckCode(ctx, md["Token-User"], util.UpdateEmailStep2, req.NewEmailDynamicPassword, req.NewEmail); e != nil {
			log.Error(ctx, "[UpdateEmail] redis op failed", log.String("operator", md["Token-User"]), log.String("code", req.NewEmailDynamicPassword), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		//verify success

		//update db and clean redis is async
		//the service's rolling update may happened between update db and clean redis
		//so we need to make this not happened
		if e := s.stop.Add(3); e != nil {
			if e == graceful.ErrClosing {
				return nil, cerror.ErrServerClosing
			}
			return nil, ecode.ErrBusy
		}
		var olduser *model.User
		if olduser, e = s.userDao.MongoUpdateUserEmail(ctx, operator, req.NewEmail); e != nil {
			s.stop.DoneOne()
			s.stop.DoneOne()
			s.stop.DoneOne()
			log.Error(ctx, "[UpdateEmail] db op failed", log.String("operator", md["Token-User"]), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		log.Info(ctx, "[UpdateEmail] success", log.String("operator", md["Token-User"]), log.String("new_email", req.NewEmail))
		if olduser.Email != req.NewEmail {
			go func() {
				if e := s.userDao.RedisDelUser(context.Background(), md["Token-User"]); e != nil {
					log.Error(ctx, "[UpdateEmail] clean redis failed", log.String("operator", md["Token-User"]), log.CError(e))
				}
				s.stop.DoneOne()
			}()
			go func() {
				if olduser.Email != "" {
					if e := s.userDao.RedisDelUserEmailIndex(context.Background(), olduser.Email); e != nil {
						log.Error(ctx, "[UpdateEmail] clean redis failed", log.String("email", olduser.Email), log.CError(e))
					}
				}
				s.stop.DoneOne()
			}()
			go func() {
				if req.NewEmail != "" {
					if e := s.userDao.RedisDelUserEmailIndex(context.Background(), req.NewEmail); e != nil {
						log.Error(ctx, "[UpdateEmail] clean redis failed", log.String("email", req.NewEmail), log.CError(e))
					}
				}
				s.stop.DoneOne()
			}()
		} else {
			s.stop.DoneOne()
			s.stop.DoneOne()
			s.stop.DoneOne()
		}
		return &api.UpdateEmailResp{Step: "success"}, nil
	} else if req.DynamicPassword != "" {
		//step 2
		if e := s.userDao.RedisCodeCheckTimes(ctx, md["Token-User"], util.UpdateEmailStep2, req.NewEmail); e != nil && e != ecode.ErrCodeNotExist {
			log.Error(ctx, "[UpdateEmail] redis op failed", log.String("operator", md["Token-User"]), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		} else if e == nil {
			//if new email's code already send,we jump to step 3
			return &api.UpdateEmailResp{Step: "newverify", Receiver: util.MaskEmail(req.NewEmail)}, nil
		}

		if e := s.userDao.RedisCheckCode(ctx, md["Token-User"], util.UpdateEmailStep1, req.DynamicPassword, ""); e != nil {
			log.Error(ctx, "[UpdateEmail] redis op failed", log.String("operator", md["Token-User"]), log.String("code", req.DynamicPassword), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		//verify success
		if e := s.sendcode(ctx, "UpdateEmail", "email", req.NewEmail, md["Token-User"], util.UpdateEmailStep2); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.UpdateEmailResp{Step: "newverify", Receiver: util.MaskEmail(req.NewEmail)}, nil
	}
	//step 1
	if e := s.userDao.RedisCodeCheckTimes(ctx, md["Token-User"], util.UpdateEmailStep2, req.NewEmail); e != nil && e != ecode.ErrCodeNotExist {
		log.Error(ctx, "[UpdateEmail] redis op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	} else if e == nil {
		//if new email's code already send,we jump to step 3
		return &api.UpdateEmailResp{Step: "newverify", Receiver: util.MaskEmail(req.NewEmail)}, nil
	}

	user, e := s.userDao.GetUser(ctx, operator)
	if e != nil {
		log.Error(ctx, "[UpdateEmail] dao op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if user.Email == req.NewEmail {
		return &api.UpdateEmailResp{Step: "success"}, nil
	}
	if req.SrcType == "tel" && user.Tel == "" {
		log.Error(ctx, "[UpdateEmail] missing tel,can't use tel to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.SrcType == "email" && user.Email == "" {
		log.Error(ctx, "[UpdateEmail] missing email,can't use email to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}

	if req.SrcType == "email" {
		e = s.sendcode(ctx, "UpdateEmail", req.SrcType, user.Email, md["Token-User"], util.UpdateEmailStep1)
	} else {
		e = s.sendcode(ctx, "UpdateEmail", req.SrcType, user.Tel, md["Token-User"], util.UpdateEmailStep1)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if req.SrcType == "email" {
		return &api.UpdateEmailResp{Step: "oldverify", Receiver: util.MaskEmail(user.Email)}, nil
	}
	return &api.UpdateEmailResp{Step: "oldverify", Receiver: util.MaskTel(user.Tel)}, nil
}

// DelEmail By OAuth
//
//	Step1:verify oauth belong to this account
//
// DelEmail By Dynamic Password
//
//	Step1:send dynamic password to email or tel
//	Step2:verify email's or tel's dynamic password
func (s *Service) DelEmail(ctx context.Context, req *api.DelEmailReq) (*api.DelEmailResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[DelEmail] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}
	update := func() (bool, error) {
		//update db and clean redis is async
		//the service's rolling update may happened between update db and clean redis
		//so we need to make this not happened
		if e := s.stop.Add(2); e != nil {
			if e == graceful.ErrClosing {
				return false, cerror.ErrServerClosing
			}
			return false, ecode.ErrBusy
		}
		var olduser *model.User
		if olduser, e = s.userDao.MongoUpdateUserEmail(ctx, operator, ""); e != nil {
			s.stop.DoneOne()
			s.stop.DoneOne()
			log.Error(ctx, "[DelEmail] db op failed", log.String("operator", md["Token-user"]), log.CError(e))
			return false, e
		}
		log.Info(ctx, "[DelEmail] success", log.String("operator", md["Token-User"]))
		if olduser.Email != "" {
			go func() {
				if e := s.userDao.RedisDelUser(context.Background(), md["Token-User"]); e != nil {
					log.Error(ctx, "[DelEmail] clean redis failed", log.String("operator", md["Token-User"]), log.CError(e))
				}
				s.stop.DoneOne()
			}()
			go func() {
				if e := s.userDao.RedisDelUserEmailIndex(context.Background(), olduser.Email); e != nil {
					log.Error(ctx, "[DelEmail] clean redis failed", log.String("email", olduser.Email), log.CError(e))
				}
				s.stop.DoneOne()
			}()
		} else {
			s.stop.DoneOne()
			s.stop.DoneOne()
		}
		var final bool
		if olduser.IDCard == "" && olduser.NickName == "" && olduser.Tel == "" && len(olduser.OAuths) == 0 {
			final = true
		}
		return final, nil
	}
	if req.SrcType == "oauth" {
		if req.SrcTypeExtra == "" || req.DynamicPassword == "" {
			return nil, ecode.ErrReq
		}
		if e := s.userDao.RedisLockEmailOP(ctx, md["Token-User"]); e != nil {
			log.Error(ctx, "[DelEmail] rate check failed", log.String("operator", md["Token-User"]), log.String(req.SrcTypeExtra, req.DynamicPassword), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		oauthid, e := util.OAuthVerifyCode(ctx, "DelEmail", req.SrcTypeExtra, req.DynamicPassword)
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		user, e := s.userDao.GetUserByOAuth(ctx, req.SrcTypeExtra, oauthid)
		if e != nil {
			log.Error(ctx, "[DelEmail] dao op failed", log.String("operator", md["Token-User"]), log.String(req.SrcTypeExtra, oauthid), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.UserID.Hex() != md["Token-User"] {
			log.Error(ctx, "[DelEmail] this is not the required oauth", log.String("operator", md["Token-User"]), log.String(req.SrcTypeExtra, oauthid))
			return nil, ecode.ErrOAuthWrong
		}
		//verify success
		final, e := update()
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.DelEmailResp{Step: "success", Final: final}, nil
	}
	if req.DynamicPassword != "" {
		//step2
		if e := s.userDao.RedisCheckCode(ctx, md["Token-User"], util.DelEmail, req.DynamicPassword, ""); e != nil {
			log.Error(ctx, "[DelEmail] redis op failed", log.String("operator", md["Token-User"]), log.String("code", req.DynamicPassword), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		//verify success
		final, e := update()
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.DelEmailResp{Step: "success", Final: final}, nil
	}
	//step1
	user, e := s.userDao.GetUser(ctx, operator)
	if e != nil {
		log.Error(ctx, "[DelEmail] dao op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	var final bool
	if user.IDCard == "" && user.NickName == "" && user.Tel == "" && len(user.OAuths) == 0 {
		final = true
	}
	if user.Email == "" {
		return &api.DelEmailResp{Step: "success", Final: final}, nil
	}

	if req.SrcType == "tel" && user.Tel == "" {
		log.Error(ctx, "[DelEmail] missing tel,can't use tel to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.SrcType == "email" && user.Email == "" {
		log.Error(ctx, "[DelEmail] missing email,can't use email to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.SrcType == "email" {
		e = s.sendcode(ctx, "DelEmail", req.SrcType, user.Email, md["Token-User"], util.DelEmail)
	} else {
		e = s.sendcode(ctx, "DelEmail", req.SrcType, user.Tel, md["Token-User"], util.DelEmail)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if req.SrcType == "email" {
		return &api.DelEmailResp{Step: "oldverify", Final: final, Receiver: util.MaskEmail(user.Email)}, nil
	}
	return &api.DelEmailResp{Step: "oldverify", Final: final, Receiver: util.MaskTel(user.Tel)}, nil
}

func (s *Service) TelDuplicateCheck(ctx context.Context, req *api.TelDuplicateCheckReq) (*api.TelDuplicateCheckResp, error) {
	md := metadata.GetMetadata(ctx)
	//redis lock
	if e := s.userDao.RedisLockDuplicateCheck(ctx, "tel", md["Token-User"]); e != nil {
		log.Error(ctx, "[TelDuplicateCheck] redis op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	userid, e := s.userDao.GetUserTelIndex(ctx, req.Tel)
	if e != nil && e != ecode.ErrUserNotExist {
		log.Error(ctx, "[TelDuplicateCheck] dao op failed", log.String("operator", md["Token-User"]), log.String("tel", req.Tel), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.TelDuplicateCheckResp{Duplicate: userid != ""}, nil
}

// TODO oauth
// UpdateTel Step 1:send dynamic password to old email or tel
// UpdateTel Step 2:verify old email's or tel's dynamic password and send dynamic password to new tel
// UpdateTel Step 3:verify new tel's dynamic and update
func (s *Service) UpdateTel(ctx context.Context, req *api.UpdateTelReq) (*api.UpdateTelResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[UpdateTel] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}
	if req.NewTelDynamicPassword != "" {
		//step 3
		if e := s.userDao.RedisCheckCode(ctx, md["Token-User"], util.UpdateTelStep2, req.NewTelDynamicPassword, req.NewTel); e != nil {
			log.Error(ctx, "[UpdateTel] redis op failed", log.String("operator", md["Token-User"]), log.String("code", req.NewTelDynamicPassword), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		//verify success

		//update db and clean redis is async
		//the service's rolling update may happened between update db and clean redis
		//so we need to make this not happened
		if e := s.stop.Add(3); e != nil {
			if e == graceful.ErrClosing {
				return nil, cerror.ErrServerClosing
			}
			return nil, ecode.ErrBusy
		}
		var olduser *model.User
		if olduser, e = s.userDao.MongoUpdateUserTel(ctx, operator, req.NewTel); e != nil {
			s.stop.DoneOne()
			s.stop.DoneOne()
			s.stop.DoneOne()
			log.Error(ctx, "[UpdateTel] db op failed", log.String("operator", md["Token-User"]), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		log.Info(ctx, "[UpdateTel] success", log.String("operator", md["Token-User"]), log.String("new_tel", req.NewTel))
		if olduser.Tel != req.NewTel {
			go func() {
				if e := s.userDao.RedisDelUser(context.Background(), md["Token-User"]); e != nil {
					log.Error(ctx, "[UpdateTel] clean redis failed", log.String("operator", md["Token-User"]), log.CError(e))
				}
				s.stop.DoneOne()
			}()
			go func() {
				if olduser.Tel != "" {
					if e := s.userDao.RedisDelUserTelIndex(context.Background(), olduser.Tel); e != nil {
						log.Error(ctx, "[UpdateTel] clean redis failed", log.String("tel", olduser.Tel), log.CError(e))
					}
				}
				s.stop.DoneOne()
			}()
			go func() {
				if req.NewTel != "" {
					if e := s.userDao.RedisDelUserTelIndex(context.Background(), req.NewTel); e != nil {
						log.Error(ctx, "[UpdateTel] clean redis failed", log.String("tel", req.NewTel), log.CError(e))
					}
				}
				s.stop.DoneOne()
			}()
		} else {
			s.stop.DoneOne()
			s.stop.DoneOne()
			s.stop.DoneOne()
		}
		return &api.UpdateTelResp{Step: "success"}, nil
	} else if req.DynamicPassword != "" {
		//step 2
		if e := s.userDao.RedisCodeCheckTimes(ctx, md["Token-User"], util.UpdateTelStep2, req.NewTel); e != nil && e != ecode.ErrCodeNotExist {
			log.Error(ctx, "[UpdateTel] redis op failed", log.String("operator", md["Token-User"]), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		} else if e == nil {
			//if new tel's code already send,we jump to step 3
			return &api.UpdateTelResp{Step: "newverify", Receiver: util.MaskTel(req.NewTel)}, nil
		}

		if e := s.userDao.RedisCheckCode(ctx, md["Token-User"], util.UpdateTelStep1, req.DynamicPassword, ""); e != nil {
			log.Error(ctx, "[UpdateTel] redis op failed", log.String("operator", md["Token-User"]), log.String("code", req.DynamicPassword), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		//verify success
		if e := s.sendcode(ctx, "UpdateTel", "tel", req.NewTel, md["Token-User"], util.UpdateTelStep2); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.UpdateTelResp{Step: "newverify", Receiver: util.MaskTel(req.NewTel)}, nil
	}
	//step 1
	if e := s.userDao.RedisCodeCheckTimes(ctx, md["Token-User"], util.UpdateTelStep2, req.NewTel); e != nil && e != ecode.ErrCodeNotExist {
		log.Error(ctx, "[UpdateTel] redis op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	} else if e == nil {
		//if new tel's code already send,we jump to step 3
		return &api.UpdateTelResp{Step: "newverify", Receiver: util.MaskTel(req.NewTel)}, nil
	}

	user, e := s.userDao.GetUser(ctx, operator)
	if e != nil {
		log.Error(ctx, "[UpdateTel] dao op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if user.Tel == req.NewTel {
		return &api.UpdateTelResp{Step: "success"}, nil
	}
	if req.SrcType == "tel" && user.Tel == "" {
		log.Error(ctx, "[UpdateTel] missing tel,can't use tel to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.SrcType == "email" && user.Email == "" {
		log.Error(ctx, "[UpdateTel] missing email,can't use email to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}

	if req.SrcType == "email" {
		e = s.sendcode(ctx, "UpdateTel", req.SrcType, user.Email, md["Token-User"], util.UpdateTelStep1)
	} else {
		e = s.sendcode(ctx, "UpdateTel", req.SrcType, user.Tel, md["Token-User"], util.UpdateTelStep1)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if req.SrcType == "email" {
		return &api.UpdateTelResp{Step: "oldverify", Receiver: util.MaskEmail(user.Email)}, nil
	}
	return &api.UpdateTelResp{Step: "oldverify", Receiver: util.MaskTel(user.Tel)}, nil
}

// DelTel By OAuth
//
//	Step1:verify oauth belong to this account
//
// DelTel By Dynamic Password
//
//	Step1:send dynamic password to email or tel
//	Step2:verify email's or tel's dynamic password
func (s *Service) DelTel(ctx context.Context, req *api.DelTelReq) (*api.DelTelResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[DelTel] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}
	update := func() (bool, error) {
		//update db and clean redis is async
		//the service's rolling update may happened between update db and clean redis
		//so we need to make this not happened
		if e := s.stop.Add(2); e != nil {
			if e == graceful.ErrClosing {
				return false, cerror.ErrServerClosing
			}
			return false, ecode.ErrBusy
		}
		var olduser *model.User
		if olduser, e = s.userDao.MongoUpdateUserTel(ctx, operator, ""); e != nil {
			s.stop.DoneOne()
			s.stop.DoneOne()
			log.Error(ctx, "[DelTel] db op failed", log.String("operator", md["Token-user"]), log.CError(e))
			return false, e
		}
		log.Info(ctx, "[DelTel] success", log.String("operator", md["Token-User"]))
		if olduser.Tel != "" {
			go func() {
				if e := s.userDao.RedisDelUser(context.Background(), md["Token-User"]); e != nil {
					log.Error(ctx, "[DelTel] clean redis failed", log.String("operator", md["Token-User"]), log.CError(e))
				}
				s.stop.DoneOne()
			}()
			go func() {
				if e := s.userDao.RedisDelUserTelIndex(context.Background(), olduser.Tel); e != nil {
					log.Error(ctx, "[DelTel] clean redis failed", log.String("tel", olduser.Tel), log.CError(e))
				}
				s.stop.DoneOne()
			}()
		} else {
			s.stop.DoneOne()
			s.stop.DoneOne()
		}
		var final bool
		if olduser.IDCard == "" && olduser.NickName == "" && olduser.Email == "" && len(olduser.OAuths) == 0 {
			final = true
		}
		return final, nil
	}
	if req.SrcType == "oauth" {
		if req.SrcTypeExtra == "" || req.DynamicPassword == "" {
			return nil, ecode.ErrReq
		}
		if e := s.userDao.RedisLockTelOP(ctx, md["Token-User"]); e != nil {
			log.Error(ctx, "[DelTel] rate check failed", log.String("operator", md["Token-User"]), log.String(req.SrcTypeExtra, req.DynamicPassword), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		oauthid, e := util.OAuthVerifyCode(ctx, "DelTel", req.SrcTypeExtra, req.DynamicPassword)
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		user, e := s.userDao.GetUserByOAuth(ctx, req.SrcTypeExtra, oauthid)
		if e != nil {
			log.Error(ctx, "[DelTel] dao op failed", log.String("operator", md["Token-User"]), log.String(req.SrcTypeExtra, oauthid), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.UserID.Hex() != md["Token-User"] {
			log.Error(ctx, "[DelTel] this is not the required oauth", log.String("operator", md["Token-User"]), log.String(req.SrcTypeExtra, oauthid))
			return nil, ecode.ErrOAuthWrong
		}
		//verify success
		final, e := update()
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.DelTelResp{Step: "success", Final: final}, nil
	}
	if req.DynamicPassword != "" {
		//step2
		if e := s.userDao.RedisCheckCode(ctx, md["Token-User"], util.DelTel, req.DynamicPassword, ""); e != nil {
			log.Error(ctx, "[DelTel] redis op failed", log.String("operator", md["Token-User"]), log.String("code", req.DynamicPassword), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		//verify success
		final, e := update()
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.DelTelResp{Step: "success", Final: final}, nil
	}
	//step1
	user, e := s.userDao.GetUser(ctx, operator)
	if e != nil {
		log.Error(ctx, "[DelTel] dao op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	var final bool
	if user.IDCard == "" && user.NickName == "" && user.Email == "" && len(user.OAuths) == 0 {
		final = true
	}
	if user.Tel == "" {
		return &api.DelTelResp{Step: "success", Final: final}, nil
	}

	if req.SrcType == "tel" && user.Tel == "" {
		log.Error(ctx, "[DelTel] missing tel,can't use tel to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.SrcType == "email" && user.Email == "" {
		log.Error(ctx, "[DelTel] missing email,can't use email to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.SrcType == "email" {
		e = s.sendcode(ctx, "DelTel", req.SrcType, user.Email, md["Token-User"], util.DelTel)
	} else {
		e = s.sendcode(ctx, "DelTel", req.SrcType, user.Tel, md["Token-User"], util.DelTel)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if req.SrcType == "email" {
		return &api.DelTelResp{Step: "oldverify", Final: final, Receiver: util.MaskEmail(user.Email)}, nil
	}
	return &api.DelTelResp{Step: "oldverify", Final: final, Receiver: util.MaskTel(user.Tel)}, nil
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
