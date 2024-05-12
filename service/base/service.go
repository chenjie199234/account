package base

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
	"github.com/chenjie199234/Corelib/cerror"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/log/trace"
	"github.com/chenjie199234/Corelib/metadata"
	publicmids "github.com/chenjie199234/Corelib/mids"
	"github.com/chenjie199234/Corelib/util/graceful"
	"github.com/chenjie199234/Corelib/web"
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
		s.stop.DoneOne()
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
		//this is controled by step1
	case util.DelOAuth:
		fallthrough
	case util.UpdateOAuth:
		e = s.userDao.RedisLockOAuthOP(ctx, operator)
	case util.DelIDCard:
		fallthrough
	case util.UpdateIDCard:
		e = s.userDao.RedisLockIDCardOP(ctx, operator)
	case util.ResetPassword:
		e = s.userDao.RedisLockResetPassword(ctx, operator)
	default:
		s.stop.DoneOne()
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
			ctx := trace.CloneSpan(ctx)
			if ee := s.userDao.RedisDelCode(ctx, operator, action); ee != nil {
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

func (s *Service) BaseInfo(ctx context.Context, req *api.BaseInfoReq) (*api.BaseInfoResp, error) {
	var user *model.User
	if _, ok := ctx.(*web.Context); ok {
		md := metadata.GetMetadata(ctx)
		operator, e := primitive.ObjectIDFromHex(md["Token-User"])
		if e != nil {
			log.Error(ctx, "[BaseInfo] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
			return nil, ecode.ErrToken
		}
		if user, e = s.userDao.GetUser(ctx, operator); e != nil {
			log.Error(ctx, "[BaseInfo] dao op failed", log.String("operator", md["Token-User"]), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	}
	if user == nil {
		switch req.SrcType {
		case "user_id":
			if req.Src == "" {
				return nil, ecode.ErrReq
			}
			userid, e := primitive.ObjectIDFromHex(req.Src)
			if e != nil {
				log.Error(ctx, "[BaseInfo] user_id format wrong", log.String("user_id", req.Src), log.CError(e))
				return nil, ecode.ErrReq
			}
			if user, e = s.userDao.GetUser(ctx, userid); e != nil {
				log.Error(ctx, "[BaseInfo] dao op failed", log.String("user_id", req.Src), log.CError(e))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
		case "tel":
			if req.Src == "" {
				return nil, ecode.ErrReq
			}
			var e error
			if user, e = s.userDao.GetUserByTel(ctx, req.Src); e != nil {
				log.Error(ctx, "[BaseInfo] dao op failed", log.String("tel", req.Src), log.CError(e))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
		case "email":
			if req.Src == "" {
				return nil, ecode.ErrReq
			}
			var e error
			if user, e = s.userDao.GetUserByEmail(ctx, req.Src); e != nil {
				log.Error(ctx, "[BaseInfo] dao op failed", log.String("email", req.Src), log.CError(e))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
		case "idcard":
			if req.Src == "" {
				return nil, ecode.ErrReq
			}
			var e error
			if user, e = s.userDao.GetUserByIDCard(ctx, req.Src); e != nil {
				log.Error(ctx, "[BaseInfo] dao op failed", log.String("idcard", req.Src), log.CError(e))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
		}
	}
	if user.BTime != 0 {
		return nil, ecode.ErrBan
	}
	resp := &api.BaseInfoResp{
		Info: &api.BaseInfo{
			UserId:     user.UserID.Hex(),
			Idcard:     user.IDCard,
			Tel:        user.Tel,
			Email:      user.Email,
			Money:      user.Money,
			BindOauths: make([]string, 0, len(user.OAuths)),
			Ctime:      uint32(user.UserID.Timestamp().Unix()),
			Ban:        user.BReason,
		},
	}
	if _, ok := ctx.(*web.Context); ok {
		resp.Info.Email = util.MaskEmail(resp.Info.Email)
		resp.Info.Tel = util.MaskTel(resp.Info.Tel)
	}
	for oauth := range user.OAuths {
		resp.Info.BindOauths = append(resp.Info.BindOauths, oauth)
	}
	return resp, nil
}

func (s *Service) Ban(ctx context.Context, req *api.BanReq) (*api.BanResp, error) {
	var userid primitive.ObjectID
	switch req.SrcType {
	case "user_id":
		var e error
		userid, e = primitive.ObjectIDFromHex(req.Src)
		if e != nil {
			log.Error(ctx, "[Ban] user_id format wrong", log.String("user_id", req.Src), log.CError(e))
			return nil, ecode.ErrReq
		}
		if user, e := s.userDao.GetUser(ctx, userid); e != nil {
			log.Error(ctx, "[Ban] dao op failed", log.String("user_id", req.Src), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		} else if user.BReason == req.Reason {
			return &api.BanResp{}, nil
		}
	case "tel":
		if user, e := s.userDao.GetUserByTel(ctx, req.Src); e != nil {
			log.Error(ctx, "[Ban] dao op failed", log.String("tel", req.Src), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		} else if user.BReason == req.Reason {
			return &api.BanResp{}, nil
		} else {
			userid = user.UserID
		}
	case "email":
		if user, e := s.userDao.GetUserByEmail(ctx, req.Src); e != nil {
			log.Error(ctx, "[Ban] dao op failed", log.String("email", req.Src), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		} else if user.BReason == req.Reason {
			return &api.BanResp{}, nil
		} else {
			userid = user.UserID
		}
	case "idcard":
		if user, e := s.userDao.GetUserByIDCard(ctx, req.Src); e != nil {
			log.Error(ctx, "[Ban] dao op failed", log.String("idcard", req.Src), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		} else if user.BReason == req.Reason {
			return &api.BanResp{}, nil
		} else {
			userid = user.UserID
		}
	}
	if e := s.stop.Add(1); e != nil {
		if e == graceful.ErrClosing {
			return nil, cerror.ErrServerClosing
		}
		return nil, ecode.ErrBusy
	}
	if e := s.userDao.MongoBanUser(ctx, userid, req.Reason); e != nil {
		s.stop.DoneOne()
		log.Error(ctx, "[Ban] db op failed", log.String(req.SrcType, req.Src), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[Ban] success", log.String(req.SrcType, req.Src))
	go func() {
		ctx := trace.CloneSpan(ctx)
		if e := s.userDao.RedisDelUser(ctx, userid.Hex()); e != nil {
			if req.SrcType == "user_id" {
				log.Error(ctx, "[Ban] clean redis failed", log.String(req.SrcType, req.Src), log.CError(e))
			} else {
				log.Error(ctx, "[Ban] clean redis failed", log.String(req.SrcType, req.Src), log.String("user_id", userid.Hex()), log.CError(e))
			}
		}
		s.stop.DoneOne()
	}()
	return &api.BanResp{}, nil
}

func (s *Service) Unban(ctx context.Context, req *api.UnbanReq) (*api.UnbanResp, error) {
	var userid primitive.ObjectID
	switch req.SrcType {
	case "user_id":
		var e error
		userid, e = primitive.ObjectIDFromHex(req.Src)
		if e != nil {
			log.Error(ctx, "[Unban] user_id format wrong", log.String("user_id", req.Src), log.CError(e))
			return nil, ecode.ErrReq
		}
		if user, e := s.userDao.GetUser(ctx, userid); e != nil {
			log.Error(ctx, "[Unban] dao op failed", log.String("user_id", req.Src), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		} else if user.BTime == 0 {
			return &api.UnbanResp{}, nil
		}
	case "tel":
		if user, e := s.userDao.GetUserByTel(ctx, req.Src); e != nil {
			log.Error(ctx, "[Unban] dao op failed", log.String("tel", req.Src), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		} else if user.BTime == 0 {
			return &api.UnbanResp{}, nil
		} else {
			userid = user.UserID
		}
	case "email":
		if user, e := s.userDao.GetUserByEmail(ctx, req.Src); e != nil {
			log.Error(ctx, "[Unban] dao op failed", log.String("email", req.Src), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		} else if user.BTime == 0 {
			return &api.UnbanResp{}, nil
		} else {
			userid = user.UserID
		}
	case "idcard":
		if user, e := s.userDao.GetUserByIDCard(ctx, req.Src); e != nil {
			log.Error(ctx, "[Unban] dao op failed", log.String("idcard", req.Src), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		} else if user.BTime == 0 {
			return &api.UnbanResp{}, nil
		} else {
			userid = user.UserID
		}
	}
	if e := s.stop.Add(1); e != nil {
		if e == graceful.ErrClosing {
			return nil, cerror.ErrServerClosing
		}
		return nil, ecode.ErrBusy
	}
	if e := s.userDao.MongoUnBanUser(ctx, userid); e != nil {
		s.stop.DoneOne()
		log.Error(ctx, "[Unban] db op failed", log.String(req.SrcType, req.Src), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[Unban] success", log.String(req.SrcType, req.Src))
	go func() {
		ctx := trace.CloneSpan(ctx)
		if e := s.userDao.RedisDelUser(ctx, userid.Hex()); e != nil {
			if req.SrcType == "user_id" {
				log.Error(ctx, "[Unban] clean redis failed", log.String("user_id", userid.Hex()), log.CError(e))
			} else {
				log.Error(ctx, "[Unban] clean redis failed", log.String(req.SrcType, req.Src), log.String("user_id", userid.Hex()), log.CError(e))
			}
		}
		s.stop.DoneOne()
	}()
	return &api.UnbanResp{}, nil
}

func (s *Service) GetOauthUrl(ctx context.Context, req *api.GetOauthUrlReq) (*api.GetOauthUrlResp, error) {
	switch req.OauthServiceName {
	case "wechat":
		if config.AC.Service.WeChatOauthUrl == "" {
			return nil, ecode.ErrOAuthUnknown
		}
		return &api.GetOauthUrlResp{Url: config.AC.Service.WeChatOauthUrl}, nil
	default:
		return nil, ecode.ErrOAuthUnknown
	}
}
func (s *Service) Login(ctx context.Context, req *api.LoginReq) (*api.LoginResp, error) {
	if req.PasswordType == "static" && req.Password == "" {
		log.Error(ctx, "[Login] empty static password", log.String(req.SrcType, req.SrcTypeExtra))
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
		if user, e = s.userDao.GetUserByIDCard(ctx, req.SrcTypeExtra); e != nil {
			log.Error(ctx, "[Login] dao op failed", log.String("idcard", req.SrcTypeExtra), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	case "tel":
		if req.PasswordType == "static" {
			if user, e = s.userDao.GetUserByTel(ctx, req.SrcTypeExtra); e != nil {
				log.Error(ctx, "[Login] dao op failed", log.String("tel", req.SrcTypeExtra), log.CError(e))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
		} else if req.Password == "" {
			if e = s.sendcode(ctx, "Login", req.SrcType, req.SrcTypeExtra, req.SrcTypeExtra, util.Login); e != nil {
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			return &api.LoginResp{Step: "verify"}, nil
		} else {
			if e = s.userDao.RedisCheckCode(ctx, req.SrcTypeExtra, util.Login, req.Password, ""); e != nil {
				log.Error(ctx, "[Login] redis op failed", log.String("tel", req.SrcTypeExtra), log.String("code", req.Password), log.CError(e))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			if user, e = s.userDao.GetOrCreateUserByTel(ctx, req.SrcTypeExtra); e != nil {
				log.Error(ctx, "[Login] dao op failed", log.String("tel", req.SrcTypeExtra), log.CError(e))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
		}
	case "email":
		if req.PasswordType == "static" {
			if user, e = s.userDao.GetUserByEmail(ctx, req.SrcTypeExtra); e != nil {
				log.Error(ctx, "[Login] dao op failed", log.String("email", req.SrcTypeExtra), log.CError(e))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
		} else if req.Password == "" {
			if e := s.sendcode(ctx, "Login", req.SrcType, req.SrcTypeExtra, req.SrcTypeExtra, util.Login); e != nil {
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			return &api.LoginResp{Step: "verify"}, nil
		} else {
			if e = s.userDao.RedisCheckCode(ctx, req.SrcTypeExtra, util.Login, req.Password, ""); e != nil {
				log.Error(ctx, "[Login] redis op failed", log.String("email", req.SrcTypeExtra), log.String("code", req.Password), log.CError(e))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			if user, e = s.userDao.GetOrCreateUserByEmail(ctx, req.SrcTypeExtra); e != nil {
				log.Error(ctx, "[Login] dao op failed", log.String("email", req.SrcTypeExtra), log.CError(e))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
		}
	case "oauth":
		if req.PasswordType == "static" {
			log.Error(ctx, "[Login] oauth can't use static password")
			return nil, ecode.ErrReq
		}
		if req.Password == "" {
			return nil, ecode.ErrReq
		}
		oauthid, e := util.OAuthVerifyCode(ctx, "Login", req.SrcTypeExtra, req.Password)
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user, e = s.userDao.GetOrCreateUserByOAuth(ctx, req.SrcTypeExtra, oauthid); e != nil {
			log.Error(ctx, "[Login] dao op failed", log.String(req.SrcTypeExtra, oauthid), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	}
	if user.BTime != 0 {
		return nil, ecode.ErrBan
	}
	if req.PasswordType == "static" {
		if e := util.SignCheck(req.Password, user.Password); e != nil {
			if e == ecode.ErrSignCheckFailed {
				e = ecode.ErrPasswordWrong
			}
			log.Error(ctx, "[Login] sign check failed", log.String(req.SrcType, req.SrcTypeExtra), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	}
	//TODO set the puber
	token := publicmids.MakeToken(ctx, "", *config.EC.DeployEnv, *config.EC.RunEnv, user.UserID.Hex(), "")
	resp := &api.LoginResp{
		Token: token,
		Info: &api.BaseInfo{
			UserId:     user.UserID.Hex(),
			Idcard:     util.MaskIDCard(user.IDCard),
			Tel:        util.MaskTel(user.Tel),
			Email:      util.MaskEmail(user.Email),
			Ctime:      uint32(user.UserID.Timestamp().Unix()),
			BindOauths: make([]string, 0, len(user.OAuths)),
			Money:      user.Money,
		},
		Step: "success",
	}
	for oauth := range user.OAuths {
		resp.Info.BindOauths = append(resp.Info.BindOauths, oauth)
	}
	if req.PasswordType == "dynamic" && util.SignCheck("", user.Password) == nil {
		//this is a new account
		resp.Step = "password"
	}
	log.Info(ctx, "[Login] success", log.String("operator", user.UserID.Hex()))
	return resp, nil
}

func (s *Service) UpdateStaticPassword(ctx context.Context, req *api.UpdateStaticPasswordReq) (*api.UpdateStaticPasswordResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[UpdateStaticPassword] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}
	if user, e := s.userDao.GetUser(ctx, operator); e != nil {
		log.Error(ctx, "[UpdateStaticPassword] dao op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	} else if user.BTime != 0 {
		return nil, ecode.ErrBan
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
		ctx := trace.CloneSpan(ctx)
		if e := s.userDao.RedisDelUser(ctx, md["Token-User"]); e != nil {
			log.Error(ctx, "[UpdateStaticPassword] clean redis failed", log.String("operator", md["Token-User"]), log.CError(e))
		}
		s.stop.DoneOne()
	}()
	return &api.UpdateStaticPasswordResp{}, nil
}

// ResetStaticPassword By OAuth
//
//	Step 1:verify oauth belong to this account
//
// ResetStaticPassword By Dynamic Password
//
//	Step 1:send dynamic password to email to tel
//	Step 2:verify email's or tel's dynamic password
func (s *Service) ResetStaticPassword(ctx context.Context, req *api.ResetStaticPasswordReq) (*api.ResetStaticPasswordResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-user"])
	if e != nil {
		log.Error(ctx, "[ResetStaticPassword] operator's token format wrong", log.String("operator", md["Token-user"]), log.CError(e))
		return nil, ecode.ErrToken
	}
	update := func() error {
		//update db and clean redis is async
		//the service's rolling update may happened between update db and clean redis
		//so we need to make this not happened
		if e := s.stop.Add(1); e != nil {
			if e == graceful.ErrClosing {
				return ecode.ErrServerClosing
			}
			return ecode.ErrBusy
		}
		if e := s.userDao.MongoResetUserPassword(ctx, operator); e != nil {
			log.Error(ctx, "[ResetStaticPassword] db op failed", log.String("operator", md["Token-user"]), log.CError(e))
			s.stop.DoneOne()
			return e
		}
		log.Info(ctx, "[ResetStaticPassword] success", log.String("operator", md["Token-User"]))
		go func() {
			ctx := trace.CloneSpan(ctx)
			if e := s.userDao.RedisDelUser(ctx, md["Token-User"]); e != nil {
				log.Error(ctx, "[ResetStaticPassword] clean redis failed", log.String("operator", md["Token-User"]), log.CError(e))
			}
			s.stop.DoneOne()
		}()
		return nil
	}
	if req.VerifySrcType == "oauth" {
		if req.VerifySrcTypeExtra == "" || req.VerifyDynamicPassword == "" {
			return nil, ecode.ErrReq
		}
		if e := s.userDao.RedisLockResetPassword(ctx, md["Token-User"]); e != nil {
			log.Error(ctx, "[ResetStaticPassword] rate check failed",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, req.VerifyDynamicPassword),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		oauthid, e := util.OAuthVerifyCode(ctx, "ResetStaticPassword", req.VerifySrcTypeExtra, req.VerifyDynamicPassword)
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		user, e := s.userDao.GetUserByOAuth(ctx, req.VerifySrcTypeExtra, oauthid)
		if e != nil {
			log.Error(ctx, "[ResetStaticPassword] dao op failed",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, oauthid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.UserID.Hex() != md["Token-User"] {
			log.Error(ctx, "[ResetStaticPassword] this is not the required oauth",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, oauthid))
			return nil, ecode.ErrOAuthWrong
		}
		if user.BTime != 0 {
			return nil, ecode.ErrBan
		}
		//verify success
		if e := update(); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.ResetStaticPasswordResp{Step: "success"}, nil
	}
	if req.VerifyDynamicPassword != "" {
		//step2
		if e := s.userDao.RedisCheckCode(ctx, md["Token-User"], util.ResetPassword, req.VerifyDynamicPassword, ""); e != nil {
			log.Error(ctx, "[ResetStaticPassword] redis op failed",
				log.String("operator", md["Token-User"]),
				log.String("code", req.VerifyDynamicPassword),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		//verify success
		if e := update(); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.ResetStaticPasswordResp{Step: "success"}, nil
	}
	//step1
	user, e := s.userDao.GetUser(ctx, operator)
	if e != nil {
		log.Error(ctx, "[ResetStaticPassword] dao op failed", log.String("operator", md["Token-user"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if user.BTime != 0 {
		return nil, ecode.ErrBan
	}

	if req.VerifySrcType == "tel" && user.Tel == "" {
		log.Error(ctx, "[ResetStaticPassword] missing tel,can't use tel to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.VerifySrcType == "email" && user.Email == "" {
		log.Error(ctx, "[ResetStaticPassword] missing email,can't use email to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}

	if req.VerifySrcType == "email" {
		e = s.sendcode(ctx, "ResetStaticPassword", "email", user.Email, md["Token-User"], util.ResetPassword)
	} else {
		e = s.sendcode(ctx, "ResetStaticPassword", "tel", user.Tel, md["Token-User"], util.ResetPassword)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if req.VerifySrcType == "email" {
		return &api.ResetStaticPasswordResp{Step: "oldverify", Receiver: util.MaskEmail(user.Email)}, nil
	}
	return &api.ResetStaticPasswordResp{Step: "oldverify", Receiver: util.MaskTel(user.Tel)}, nil
}

// UpdateOAuth By OAuth
//
//	Step 1:verify oauth belong to this account and verify the new oauth
//
// UpdateOAuth By Dynamic Password
//
//	Step 1:send dynamic password to old email or tel
//	Step 2:verify old email's or tel's dynamic password and verify the new oauth
func (s *Service) UpdateOauth(ctx context.Context, req *api.UpdateOauthReq) (*api.UpdateOauthResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[UpdateOauth] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}
	update := func(newoauthid string) error {
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
		if olduser, e = s.userDao.MongoUpdateUserOAuth(ctx, operator, req.NewOauthServiceName, newoauthid); e != nil {
			s.stop.DoneOne()
			s.stop.DoneOne()
			s.stop.DoneOne()
			log.Error(ctx, "[UpdateOauth] db op failed", log.String("operator", md["Token-User"]), log.CError(e))
			return e
		}
		log.Info(ctx, "[UpdateOauth] success", log.String("operator", md["Token-User"]), log.String(req.NewOauthServiceName, newoauthid))
		oldoauthid := olduser.OAuths[req.NewOauthServiceName]
		if oldoauthid != newoauthid {
			go func() {
				ctx := trace.CloneSpan(ctx)
				if e := s.userDao.RedisDelUser(ctx, md["Token-User"]); e != nil {
					log.Error(ctx, "[UpdateOauth] clean redis failed", log.String("operator", md["Token-User"]), log.CError(e))
				}
				s.stop.DoneOne()
			}()
			go func() {
				if oldoauthid != "" {
					ctx := trace.CloneSpan(ctx)
					if e := s.userDao.RedisDelUserOAuthIndex(ctx, req.NewOauthServiceName, oldoauthid); e != nil {
						log.Error(ctx, "[UpdateOauth] clean redis failed", log.String(req.NewOauthServiceName, oldoauthid), log.CError(e))
					}
				}
				s.stop.DoneOne()
			}()
			go func() {
				if newoauthid != "" {
					ctx := trace.CloneSpan(ctx)
					if e := s.userDao.RedisDelUserOAuthIndex(ctx, req.NewOauthServiceName, newoauthid); e != nil {
						log.Error(ctx, "[UpdateOauth] clean redis failed", log.String(req.NewOauthServiceName, newoauthid), log.CError(e))
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
	if req.VerifySrcType == "oauth" {
		if req.VerifySrcTypeExtra == "" || req.VerifyDynamicPassword == "" || req.NewOauthServiceName == "" || req.NewOauthDynamicPassword == "" {
			return nil, ecode.ErrReq
		}
		if e := s.userDao.RedisLockOAuthOP(ctx, md["Token-User"]); e != nil {
			log.Error(ctx, "[UpdateOauth] rate check failed",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, req.VerifyDynamicPassword),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		oauthid, e := util.OAuthVerifyCode(ctx, "UpdateOauth", req.VerifySrcTypeExtra, req.VerifyDynamicPassword)
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		user, e := s.userDao.GetUserByOAuth(ctx, req.VerifySrcTypeExtra, oauthid)
		if e != nil {
			log.Error(ctx, "[UpdateOauth] dao op failed",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, oauthid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.UserID.Hex() != md["Token-User"] {
			log.Error(ctx, "[UpdateOauth] this is not the required oauth",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, oauthid))
			return nil, ecode.ErrOAuthWrong
		}
		if user.BTime != 0 {
			return nil, ecode.ErrBan
		}
		//verify success
		//get the new oauth
		oauthid, e = util.OAuthVerifyCode(ctx, "UpdateOauth", req.NewOauthServiceName, req.NewOauthDynamicPassword)
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.OAuths[req.NewOauthServiceName] == oauthid {
			return &api.UpdateOauthResp{Step: "success"}, nil
		}
		if e := update(oauthid); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.UpdateOauthResp{Step: "success"}, nil
	}
	if req.VerifyDynamicPassword != "" {
		if req.NewOauthServiceName == "" || req.NewOauthDynamicPassword == "" {
			return nil, ecode.ErrReq
		}
		if e := s.userDao.RedisCheckCode(ctx, md["Token-User"], util.UpdateOAuth, req.VerifyDynamicPassword, ""); e != nil {
			log.Error(ctx, "[UpdateOauth] redis op failed",
				log.String("operator", md["Token-User"]),
				log.String("code", req.VerifyDynamicPassword),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		//verify success
		//get the new oauth
		oauthid, e := util.OAuthVerifyCode(ctx, "UpdateOauth", req.NewOauthServiceName, req.NewOauthDynamicPassword)
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if e := update(oauthid); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.UpdateOauthResp{Step: "success"}, nil
	}
	//send dynamic password
	user, e := s.userDao.GetUser(ctx, operator)
	if e != nil {
		log.Error(ctx, "[UpdateOauth] dao op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if user.BTime != 0 {
		return nil, ecode.ErrBan
	}

	if req.VerifySrcType == "email" && user.Email == "" {
		log.Error(ctx, "[UpdateOauth] missing email,can't use email to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.VerifySrcType == "tel" && user.Tel == "" {
		log.Error(ctx, "[UpdateOauth] missing tel,can't use tel to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}

	if req.VerifySrcType == "email" {
		e = s.sendcode(ctx, "UpdateOauth", "email", user.Email, md["Token-User"], util.UpdateOAuth)
	} else {
		e = s.sendcode(ctx, "UpdateOauth", "tel", user.Tel, md["Token-User"], util.UpdateOAuth)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if req.VerifySrcType == "email" {
		return &api.UpdateOauthResp{Step: "oldverify", Receiver: util.MaskEmail(user.Email)}, nil
	}
	return &api.UpdateOauthResp{Step: "oldverify", Receiver: util.MaskTel(user.Tel)}, nil
}

// DelOauth By OAuth
//
//	Step 1:verify oauth belong to this account
//
// DelOauth By Dynamic Password
//
//	Step 1:send dynamic password to email or tel
//	Step 2:verify email's or tel's dynamic password
func (s *Service) DelOauth(ctx context.Context, req *api.DelOauthReq) (*api.DelOauthResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[DelOauth] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
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
		if olduser, e = s.userDao.MongoUpdateUserOAuth(ctx, operator, req.DelOauthServiceName, ""); e != nil {
			s.stop.DoneOne()
			s.stop.DoneOne()
			log.Error(ctx, "[DelOauth] db op failed", log.String("operator", md["Token-user"]), log.CError(e))
			return false, e
		}
		log.Info(ctx, "[DelOauth] success", log.String("operator", md["Token-User"]), log.String("oauth", req.DelOauthServiceName))
		if oauthid := olduser.OAuths[req.DelOauthServiceName]; oauthid != "" {
			go func() {
				ctx := trace.CloneSpan(ctx)
				if e := s.userDao.RedisDelUser(ctx, md["Token-User"]); e != nil {
					log.Error(ctx, "[DelOauth] clean redis failed", log.String("operator", md["Token-User"]), log.CError(e))
				}
				s.stop.DoneOne()
			}()
			go func() {
				ctx := trace.CloneSpan(ctx)
				if e := s.userDao.RedisDelUserOAuthIndex(ctx, req.DelOauthServiceName, oauthid); e != nil {
					log.Error(ctx, "[DelOauth] clean redis failed", log.String(req.DelOauthServiceName, oauthid), log.CError(e))
				}
				s.stop.DoneOne()
			}()
		} else {
			s.stop.DoneOne()
			s.stop.DoneOne()
		}
		var final bool
		if olduser.Email == "" && olduser.IDCard == "" && olduser.Tel == "" && len(olduser.OAuths) == 1 && olduser.OAuths[req.DelOauthServiceName] != "" {
			final = true
		}
		return final, nil
	}
	if req.VerifySrcType == "oauth" {
		if req.VerifySrcTypeExtra == "" || req.VerifyDynamicPassword == "" {
			return nil, ecode.ErrReq
		}
		if e := s.userDao.RedisLockOAuthOP(ctx, md["Token-User"]); e != nil {
			log.Error(ctx, "[DelOauth] rate check failed",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, req.VerifyDynamicPassword),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		oauthid, e := util.OAuthVerifyCode(ctx, "DelOauth", req.VerifySrcTypeExtra, req.VerifyDynamicPassword)
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		user, e := s.userDao.GetUserByOAuth(ctx, req.VerifySrcTypeExtra, oauthid)
		if e != nil {
			log.Error(ctx, "[DelOauth] dao op failed",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, oauthid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.UserID.Hex() != md["Token-User"] {
			log.Error(ctx, "[DelOauth] this is not the required oauth",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, oauthid))
			return nil, ecode.ErrOAuthWrong
		}
		if _, ok := user.OAuths[req.DelOauthServiceName]; !ok {
			return &api.DelOauthResp{Step: "success", Final: user.Email == "" && user.IDCard == "" && user.Tel == "" && len(user.OAuths) == 0}, nil
		}
		if user.BTime != 0 {
			return nil, ecode.ErrBan
		}
		//verify success
		final, e := update()
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.DelOauthResp{Step: "success", Final: final}, nil
	}
	if req.VerifyDynamicPassword != "" {
		//step2
		if e := s.userDao.RedisCheckCode(ctx, md["Token-User"], util.DelOAuth, req.VerifyDynamicPassword, ""); e != nil {
			log.Error(ctx, "[DelOauth] redis op failed",
				log.String("operator", md["Token-User"]),
				log.String("code", req.VerifyDynamicPassword),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		//verify success
		final, e := update()
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.DelOauthResp{Step: "success", Final: final}, nil
	}
	//step1
	user, e := s.userDao.GetUser(ctx, operator)
	if e != nil {
		log.Error(ctx, "[DelOauth] dao op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	var final bool
	if user.Email == "" && user.IDCard == "" && user.Tel == "" && len(user.OAuths) == 0 {
		final = true
	}
	if user.OAuths[req.DelOauthServiceName] == "" {
		return &api.DelOauthResp{Step: "success", Final: final}, nil
	}
	if user.BTime != 0 {
		return nil, ecode.ErrBan
	}

	if req.VerifySrcType == "tel" && user.Tel == "" {
		log.Error(ctx, "[DelOauth] missing tel,can't use tel to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.VerifySrcType == "email" && user.Email == "" {
		log.Error(ctx, "[DelOauth] missing email,can't use email to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}

	if req.VerifySrcType == "email" {
		e = s.sendcode(ctx, "DelOauth", "email", user.Email, md["Token-User"], util.DelOAuth)
	} else {
		e = s.sendcode(ctx, "DelOauth", "tel", user.Tel, md["Token-User"], util.DelOAuth)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if req.VerifySrcType == "email" {
		return &api.DelOauthResp{Step: "oldverify", Final: final, Receiver: util.MaskEmail(user.Email)}, nil
	}
	return &api.DelOauthResp{Step: "oldverify", Final: final, Receiver: util.MaskTel(user.Tel)}, nil
}

func (s *Service) IdcardDuplicateCheck(ctx context.Context, req *api.IdcardDuplicateCheckReq) (*api.IdcardDuplicateCheckResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[IdcardDuplicateCheck] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}
	if user, e := s.userDao.GetUser(ctx, operator); e != nil {
		log.Error(ctx, "[IdcardDuplicateCheck] dao op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	} else if user.BTime != 0 {
		return nil, ecode.ErrBan
	}
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
//	Step 1:verify oauth belong to this account
//
// UpdateIdCard By Dynamic Password
//
//	Step 1:send dynamic password to email or tel
//	Step 2:verify email's or tel's dynamic password
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
				ctx := trace.CloneSpan(ctx)
				if e := s.userDao.RedisDelUser(ctx, md["Token-User"]); e != nil {
					log.Error(ctx, "[UpdateIdcard] clean redis failed", log.String("operator", md["Token-User"]), log.CError(e))
				}
				s.stop.DoneOne()
			}()
			go func() {
				if olduser.IDCard != "" {
					ctx := trace.CloneSpan(ctx)
					if e := s.userDao.RedisDelUserIDCardIndex(ctx, olduser.IDCard); e != nil {
						log.Error(ctx, "[UpdateIdcard] clean redis failed", log.String("idcard", olduser.IDCard), log.CError(e))
					}
				}
				s.stop.DoneOne()
			}()
			go func() {
				if req.NewIdcard != "" {
					ctx := trace.CloneSpan(ctx)
					if e := s.userDao.RedisDelUserIDCardIndex(ctx, req.NewIdcard); e != nil {
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
	if req.VerifySrcType == "oauth" {
		if req.VerifySrcTypeExtra == "" || req.VerifyDynamicPassword == "" {
			return nil, ecode.ErrReq
		}
		if e := s.userDao.RedisLockIDCardOP(ctx, md["Token-User"]); e != nil {
			log.Error(ctx, "[UpdateIdcard] rate check failed",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, req.VerifyDynamicPassword),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		oauthid, e := util.OAuthVerifyCode(ctx, "UpdateIdcard", req.VerifySrcTypeExtra, req.VerifyDynamicPassword)
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		user, e := s.userDao.GetUserByOAuth(ctx, req.VerifySrcTypeExtra, oauthid)
		if e != nil {
			log.Error(ctx, "[UpdateIdcard] dao op failed",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, oauthid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.UserID.Hex() != md["Token-User"] {
			log.Error(ctx, "[UpdateIdcard] this is not the required oauth",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, oauthid))
			return nil, ecode.ErrOAuthWrong
		}
		if user.IDCard == req.NewIdcard {
			return &api.UpdateIdcardResp{Step: "success"}, nil
		}
		if user.BTime != 0 {
			return nil, ecode.ErrBan
		}
		//verify success
		if e := update(); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.UpdateIdcardResp{Step: "success"}, nil
	}
	if req.VerifyDynamicPassword != "" {
		//step 2
		if e := s.userDao.RedisCheckCode(ctx, md["Token-User"], util.UpdateIDCard, req.VerifyDynamicPassword, ""); e != nil {
			log.Error(ctx, "[UpdateIdcard] redis op failed",
				log.String("operator", md["Token-User"]),
				log.String("code", req.VerifyDynamicPassword),
				log.CError(e))
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
	if user.BTime != 0 {
		return nil, ecode.ErrBan
	}

	if req.VerifySrcType == "tel" && user.Tel == "" {
		log.Error(ctx, "[UpdateIdcard] missing tel,can't use tel to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.VerifySrcType == "email" && user.Email == "" {
		log.Error(ctx, "[UpdateIdcard] missing email,can't use email to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}

	if req.VerifySrcType == "email" {
		e = s.sendcode(ctx, "UpdateIdcard", "email", user.Email, md["Token-User"], util.UpdateIDCard)
	} else {
		e = s.sendcode(ctx, "UpdateIdcard", "tel", user.Tel, md["Token-User"], util.UpdateIDCard)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if req.VerifySrcType == "email" {
		return &api.UpdateIdcardResp{Step: "oldverify", Receiver: util.MaskEmail(user.Email)}, nil
	}
	return &api.UpdateIdcardResp{Step: "oldverify", Receiver: util.MaskTel(user.Tel)}, nil
}

// DelIdcard By OAuth
//
//	Step 1:verify oauth belong to this account
//
// DelIdcard By Dynamic Password
//
//	Step 1:send dynamic password to email or tel
//	Step 2:verify email's or tel's dynamic password
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
				ctx := trace.CloneSpan(ctx)
				if e := s.userDao.RedisDelUser(ctx, md["Token-User"]); e != nil {
					log.Error(ctx, "[DelIdcard] clean redis failed", log.String("operator", md["Token-User"]), log.CError(e))
				}
				s.stop.DoneOne()
			}()
			go func() {
				ctx := trace.CloneSpan(ctx)
				if e := s.userDao.RedisDelUserIDCardIndex(ctx, olduser.IDCard); e != nil {
					log.Error(ctx, "[DelIdcard] clean redis failed", log.String("idcard", olduser.IDCard), log.CError(e))
				}
				s.stop.DoneOne()
			}()
		} else {
			s.stop.DoneOne()
			s.stop.DoneOne()
		}
		var final bool
		if olduser.Email == "" && olduser.Tel == "" && len(olduser.OAuths) == 0 {
			final = true
		}
		return final, nil
	}
	if req.VerifySrcType == "oauth" {
		if req.VerifySrcTypeExtra == "" || req.VerifyDynamicPassword == "" {
			return nil, ecode.ErrReq
		}
		if e := s.userDao.RedisLockIDCardOP(ctx, md["Token-User"]); e != nil {
			log.Error(ctx, "[DelIDCard] rate check failed",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, req.VerifyDynamicPassword),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		oauthid, e := util.OAuthVerifyCode(ctx, "DelIdcard", req.VerifySrcTypeExtra, req.VerifyDynamicPassword)
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		user, e := s.userDao.GetUserByOAuth(ctx, req.VerifySrcTypeExtra, oauthid)
		if e != nil {
			log.Error(ctx, "[DelIdcard] dao op failed",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, oauthid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.UserID.Hex() != md["Token-User"] {
			log.Error(ctx, "[DelIdcard] this is not the required oauth",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, oauthid))
			return nil, ecode.ErrOAuthWrong
		}
		if user.IDCard == "" {
			return &api.DelIdcardResp{Step: "success", Final: user.Email == "" && user.Tel == "" && user.IDCard == "" && len(user.OAuths) == 0}, nil
		}
		if user.BTime != 0 {
			return nil, ecode.ErrBan
		}
		//verify success
		final, e := update()
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.DelIdcardResp{Step: "success", Final: final}, nil
	}
	if req.VerifyDynamicPassword != "" {
		//step2
		if e := s.userDao.RedisCheckCode(ctx, md["Token-User"], util.DelIDCard, req.VerifyDynamicPassword, ""); e != nil {
			log.Error(ctx, "[DelIdcard] redis op failed",
				log.String("operator", md["Token-User"]),
				log.String("code", req.VerifyDynamicPassword),
				log.CError(e))
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
	if user.Email == "" && user.Tel == "" && len(user.OAuths) == 0 {
		final = true
	}
	if user.IDCard == "" {
		return &api.DelIdcardResp{Step: "success", Final: final}, nil
	}
	if user.BTime != 0 {
		return nil, ecode.ErrBan
	}

	if req.VerifySrcType == "tel" && user.Tel == "" {
		log.Error(ctx, "[DelIdcard] missing tel,can't use tel to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.VerifySrcType == "email" && user.Email == "" {
		log.Error(ctx, "[DelIdcard] missing email,can't use email to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}

	if req.VerifySrcType == "email" {
		e = s.sendcode(ctx, "DelIdcard", "email", user.Email, md["Token-User"], util.DelIDCard)
	} else {
		e = s.sendcode(ctx, "DelIdcard", "tel", user.Tel, md["Token-User"], util.DelIDCard)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if req.VerifySrcType == "email" {
		return &api.DelIdcardResp{Step: "oldverify", Final: final, Receiver: util.MaskEmail(user.Email)}, nil
	}
	return &api.DelIdcardResp{Step: "oldverify", Final: final, Receiver: util.MaskTel(user.Tel)}, nil
}

func (s *Service) EmailDuplicateCheck(ctx context.Context, req *api.EmailDuplicateCheckReq) (*api.EmailDuplicateCheckResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[EmailDuplicateCheck] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}
	if user, e := s.userDao.GetUser(ctx, operator); e != nil {
		log.Error(ctx, "[EmailDuplicateCheck] dao op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	} else if user.BTime != 0 {
		return nil, ecode.ErrBan
	}
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

// UpdateEmail By OAuth
//
//	Step 1:verify oauth belong to this account and send dynamic password to new email
//	Step final:verify new email's dynamic password and update
//
// UpdateEmail By Dynamic Password
//
//	Step 1:send dynamic password to old email or tel
//	Step 2:verify old email's or tel's dynamic password and send dynamic password to new email
//	Step final:verify new email's dynamic password and update
func (s *Service) UpdateEmail(ctx context.Context, req *api.UpdateEmailReq) (*api.UpdateEmailResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[UpdateEmail] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}
	if req.NewEmailDynamicPassword != "" {
		//step final
		if e := s.userDao.RedisCheckCode(ctx, md["Token-User"], util.UpdateEmailStep2, req.NewEmailDynamicPassword, req.NewEmail); e != nil {
			log.Error(ctx, "[UpdateEmail] redis op failed",
				log.String("operator", md["Token-User"]),
				log.String("code", req.NewEmailDynamicPassword),
				log.CError(e))
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
				ctx := trace.CloneSpan(ctx)
				if e := s.userDao.RedisDelUser(ctx, md["Token-User"]); e != nil {
					log.Error(ctx, "[UpdateEmail] clean redis failed", log.String("operator", md["Token-User"]), log.CError(e))
				}
				s.stop.DoneOne()
			}()
			go func() {
				if olduser.Email != "" {
					ctx := trace.CloneSpan(ctx)
					if e := s.userDao.RedisDelUserEmailIndex(ctx, olduser.Email); e != nil {
						log.Error(ctx, "[UpdateEmail] clean redis failed", log.String("email", olduser.Email), log.CError(e))
					}
				}
				s.stop.DoneOne()
			}()
			go func() {
				if req.NewEmail != "" {
					ctx := trace.CloneSpan(ctx)
					if e := s.userDao.RedisDelUserEmailIndex(ctx, req.NewEmail); e != nil {
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
	}
	if e := s.userDao.RedisCodeCheckTimes(ctx, md["Token-User"], util.UpdateEmailStep2, req.NewEmail); e != nil && e != ecode.ErrCodeNotExist {
		log.Error(ctx, "[UpdateEmail] redis op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	} else if e == nil {
		//if new email's code already send,we jump to step final
		return &api.UpdateEmailResp{Step: "newverify", Receiver: util.MaskEmail(req.NewEmail)}, nil
	}
	if req.VerifySrcType == "oauth" {
		//step 1 when update by oauth
		if req.VerifySrcTypeExtra == "" || req.VerifyDynamicPassword == "" {
			return nil, ecode.ErrReq
		}
		if e := s.userDao.RedisLockEmailOP(ctx, md["Token-User"]); e != nil {
			log.Error(ctx, "[UpdateEmail] rate check failed",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, req.VerifyDynamicPassword),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		oauthid, e := util.OAuthVerifyCode(ctx, "UpdateEmail", req.VerifySrcTypeExtra, req.VerifyDynamicPassword)
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		user, e := s.userDao.GetUserByOAuth(ctx, req.VerifySrcTypeExtra, oauthid)
		if e != nil {
			log.Error(ctx, "[UpdateEmail] dao op failed",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, oauthid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.UserID.Hex() != md["Token-User"] {
			log.Error(ctx, "[UpdateEmail] this is not the required oauth",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, oauthid))
			return nil, ecode.ErrOAuthWrong
		}
		if user.Email == req.NewEmail {
			return &api.UpdateEmailResp{Step: "success"}, nil
		}
		if user.BTime != 0 {
			return nil, ecode.ErrBan
		}
		//verify success
		if e := s.sendcode(ctx, "UpdateEmail", "email", req.NewEmail, md["Token-User"], util.UpdateEmailStep2); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.UpdateEmailResp{Step: "newverify", Receiver: util.MaskEmail(req.NewEmail)}, nil
	}
	if req.VerifyDynamicPassword != "" {
		//step 2 when update by dynamic password
		if e := s.userDao.RedisCheckCode(ctx, md["Token-User"], util.UpdateEmailStep1, req.VerifyDynamicPassword, ""); e != nil {
			log.Error(ctx, "[UpdateEmail] redis op failed",
				log.String("operator", md["Token-User"]),
				log.String("code", req.VerifyDynamicPassword),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		//verify success
		if e := s.sendcode(ctx, "UpdateEmail", "email", req.NewEmail, md["Token-User"], util.UpdateEmailStep2); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.UpdateEmailResp{Step: "newverify", Receiver: util.MaskEmail(req.NewEmail)}, nil
	}
	//step 1 when update by dynamic password
	user, e := s.userDao.GetUser(ctx, operator)
	if e != nil {
		log.Error(ctx, "[UpdateEmail] dao op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if user.Email == req.NewEmail {
		return &api.UpdateEmailResp{Step: "success"}, nil
	}
	if user.BTime != 0 {
		return nil, ecode.ErrBan
	}

	if req.VerifySrcType == "tel" && user.Tel == "" {
		log.Error(ctx, "[UpdateEmail] missing tel,can't use tel to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.VerifySrcType == "email" && user.Email == "" {
		log.Error(ctx, "[UpdateEmail] missing email,can't use email to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}

	if req.VerifySrcType == "email" {
		e = s.sendcode(ctx, "UpdateEmail", "email", user.Email, md["Token-User"], util.UpdateEmailStep1)
	} else {
		e = s.sendcode(ctx, "UpdateEmail", "tel", user.Tel, md["Token-User"], util.UpdateEmailStep1)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if req.VerifySrcType == "email" {
		return &api.UpdateEmailResp{Step: "oldverify", Receiver: util.MaskEmail(user.Email)}, nil
	}
	return &api.UpdateEmailResp{Step: "oldverify", Receiver: util.MaskTel(user.Tel)}, nil
}

// DelEmail By OAuth
//
//	Step 1:verify oauth belong to this account
//
// DelEmail By Dynamic Password
//
//	Step 1:send dynamic password to email or tel
//	Step 2:verify email's or tel's dynamic password
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
				ctx := trace.CloneSpan(ctx)
				if e := s.userDao.RedisDelUser(ctx, md["Token-User"]); e != nil {
					log.Error(ctx, "[DelEmail] clean redis failed", log.String("operator", md["Token-User"]), log.CError(e))
				}
				s.stop.DoneOne()
			}()
			go func() {
				ctx := trace.CloneSpan(ctx)
				if e := s.userDao.RedisDelUserEmailIndex(ctx, olduser.Email); e != nil {
					log.Error(ctx, "[DelEmail] clean redis failed", log.String("email", olduser.Email), log.CError(e))
				}
				s.stop.DoneOne()
			}()
		} else {
			s.stop.DoneOne()
			s.stop.DoneOne()
		}
		var final bool
		if olduser.IDCard == "" && olduser.Tel == "" && len(olduser.OAuths) == 0 {
			final = true
		}
		return final, nil
	}
	if req.VerifySrcType == "oauth" {
		if req.VerifySrcTypeExtra == "" || req.VerifyDynamicPassword == "" {
			return nil, ecode.ErrReq
		}
		if e := s.userDao.RedisLockEmailOP(ctx, md["Token-User"]); e != nil {
			log.Error(ctx, "[DelEmail] rate check failed",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, req.VerifyDynamicPassword),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		oauthid, e := util.OAuthVerifyCode(ctx, "DelEmail", req.VerifySrcTypeExtra, req.VerifyDynamicPassword)
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		user, e := s.userDao.GetUserByOAuth(ctx, req.VerifySrcTypeExtra, oauthid)
		if e != nil {
			log.Error(ctx, "[DelEmail] dao op failed",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, oauthid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.UserID.Hex() != md["Token-User"] {
			log.Error(ctx, "[DelEmail] this is not the required oauth",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, oauthid))
			return nil, ecode.ErrOAuthWrong
		}
		if user.Email == "" {
			return &api.DelEmailResp{Step: "success", Final: user.Email == "" && user.Tel == "" && user.IDCard == "" && len(user.OAuths) == 0}, nil
		}
		if user.BTime != 0 {
			return nil, ecode.ErrBan
		}
		//verify success
		final, e := update()
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.DelEmailResp{Step: "success", Final: final}, nil
	}
	if req.VerifyDynamicPassword != "" {
		//step2
		if e := s.userDao.RedisCheckCode(ctx, md["Token-User"], util.DelEmail, req.VerifyDynamicPassword, ""); e != nil {
			log.Error(ctx, "[DelEmail] redis op failed",
				log.String("operator", md["Token-User"]),
				log.String("code", req.VerifyDynamicPassword),
				log.CError(e))
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
	if user.IDCard == "" && user.Tel == "" && len(user.OAuths) == 0 {
		final = true
	}
	if user.Email == "" {
		return &api.DelEmailResp{Step: "success", Final: final}, nil
	}
	if user.BTime != 0 {
		return nil, ecode.ErrBan
	}

	if req.VerifySrcType == "tel" && user.Tel == "" {
		log.Error(ctx, "[DelEmail] missing tel,can't use tel to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.VerifySrcType == "email" && user.Email == "" {
		log.Error(ctx, "[DelEmail] missing email,can't use email to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}

	if req.VerifySrcType == "email" {
		e = s.sendcode(ctx, "DelEmail", "email", user.Email, md["Token-User"], util.DelEmail)
	} else {
		e = s.sendcode(ctx, "DelEmail", "tel", user.Tel, md["Token-User"], util.DelEmail)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if req.VerifySrcType == "email" {
		return &api.DelEmailResp{Step: "oldverify", Final: final, Receiver: util.MaskEmail(user.Email)}, nil
	}
	return &api.DelEmailResp{Step: "oldverify", Final: final, Receiver: util.MaskTel(user.Tel)}, nil
}

func (s *Service) TelDuplicateCheck(ctx context.Context, req *api.TelDuplicateCheckReq) (*api.TelDuplicateCheckResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[TelDuplicateCheck] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}
	if user, e := s.userDao.GetUser(ctx, operator); e != nil {
		log.Error(ctx, "[TelDuplicateCheck] dao op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	} else if user.BTime != 0 {
		return nil, ecode.ErrBan
	}
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

// UpdateTel By OAuth
//
//	Step 1:verify oauth belong to this account and send dynamic password to new tel
//	Step final:verify new tel's dynamic password and update
//
// UpdateTel By Dynamic Password
//
//	Step 1:send dynamic password to old email or tel
//	Step 2:verify old email's or tel's dynamic password and send dynamic password to new tel
//	Step final:verify new tel's dynamic password and update
func (s *Service) UpdateTel(ctx context.Context, req *api.UpdateTelReq) (*api.UpdateTelResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[UpdateTel] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}
	if req.NewTelDynamicPassword != "" {
		//step final
		if e := s.userDao.RedisCheckCode(ctx, md["Token-User"], util.UpdateTelStep2, req.NewTelDynamicPassword, req.NewTel); e != nil {
			log.Error(ctx, "[UpdateTel] redis op failed",
				log.String("operator", md["Token-User"]),
				log.String("code", req.NewTelDynamicPassword),
				log.CError(e))
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
				ctx := trace.CloneSpan(ctx)
				if e := s.userDao.RedisDelUser(ctx, md["Token-User"]); e != nil {
					log.Error(ctx, "[UpdateTel] clean redis failed", log.String("operator", md["Token-User"]), log.CError(e))
				}
				s.stop.DoneOne()
			}()
			go func() {
				if olduser.Tel != "" {
					ctx := trace.CloneSpan(ctx)
					if e := s.userDao.RedisDelUserTelIndex(ctx, olduser.Tel); e != nil {
						log.Error(ctx, "[UpdateTel] clean redis failed", log.String("tel", olduser.Tel), log.CError(e))
					}
				}
				s.stop.DoneOne()
			}()
			go func() {
				if req.NewTel != "" {
					ctx := trace.CloneSpan(ctx)
					if e := s.userDao.RedisDelUserTelIndex(ctx, req.NewTel); e != nil {
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
	}
	if e := s.userDao.RedisCodeCheckTimes(ctx, md["Token-User"], util.UpdateTelStep2, req.NewTel); e != nil && e != ecode.ErrCodeNotExist {
		log.Error(ctx, "[UpdateTel] redis op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	} else if e == nil {
		//if new tel's code already send,we jump to step final
		return &api.UpdateTelResp{Step: "newverify", Receiver: util.MaskTel(req.NewTel)}, nil
	}
	if req.VerifySrcType == "oauth" {
		//step 1 when update by oauth
		if req.VerifySrcTypeExtra == "" || req.VerifyDynamicPassword == "" {
			return nil, ecode.ErrReq
		}
		if e := s.userDao.RedisLockTelOP(ctx, md["Token-User"]); e != nil {
			log.Error(ctx, "[UpdateTel] rate check failed",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, req.VerifyDynamicPassword),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		oauthid, e := util.OAuthVerifyCode(ctx, "UpdateTel", req.VerifySrcTypeExtra, req.VerifyDynamicPassword)
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		user, e := s.userDao.GetUserByOAuth(ctx, req.VerifySrcTypeExtra, oauthid)
		if e != nil {
			log.Error(ctx, "[UpdateTel] dao op failed",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, oauthid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.UserID.Hex() != md["Token-User"] {
			log.Error(ctx, "[UpdateTel] this is not the required oauth",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, oauthid))
			return nil, ecode.ErrOAuthWrong
		}
		if user.Tel == req.NewTel {
			return &api.UpdateTelResp{Step: "success"}, nil
		}
		if user.BTime != 0 {
			return nil, ecode.ErrBan
		}
		//verify success
		if e := s.sendcode(ctx, "UpdateTel", "tel", req.NewTel, md["Token-User"], util.UpdateTelStep2); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.UpdateTelResp{Step: "newverify", Receiver: util.MaskTel(req.NewTel)}, nil
	}
	if req.VerifyDynamicPassword != "" {
		//step 2 when update by dynamic password
		if e := s.userDao.RedisCheckCode(ctx, md["Token-User"], util.UpdateTelStep1, req.VerifyDynamicPassword, ""); e != nil {
			log.Error(ctx, "[UpdateTel] redis op failed",
				log.String("operator", md["Token-User"]),
				log.String("code", req.VerifyDynamicPassword),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		//verify success
		if e := s.sendcode(ctx, "UpdateTel", "tel", req.NewTel, md["Token-User"], util.UpdateTelStep2); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.UpdateTelResp{Step: "newverify", Receiver: util.MaskTel(req.NewTel)}, nil
	}
	//step 1 when update by dynamic password
	user, e := s.userDao.GetUser(ctx, operator)
	if e != nil {
		log.Error(ctx, "[UpdateTel] dao op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if user.Tel == req.NewTel {
		return &api.UpdateTelResp{Step: "success"}, nil
	}
	if user.BTime != 0 {
		return nil, ecode.ErrBan
	}

	if req.VerifySrcType == "tel" && user.Tel == "" {
		log.Error(ctx, "[UpdateTel] missing tel,can't use tel to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.VerifySrcType == "email" && user.Email == "" {
		log.Error(ctx, "[UpdateTel] missing email,can't use email to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}

	if req.VerifySrcType == "email" {
		e = s.sendcode(ctx, "UpdateTel", "email", user.Email, md["Token-User"], util.UpdateTelStep1)
	} else {
		e = s.sendcode(ctx, "UpdateTel", "tel", user.Tel, md["Token-User"], util.UpdateTelStep1)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if req.VerifySrcType == "email" {
		return &api.UpdateTelResp{Step: "oldverify", Receiver: util.MaskEmail(user.Email)}, nil
	}
	return &api.UpdateTelResp{Step: "oldverify", Receiver: util.MaskTel(user.Tel)}, nil
}

// DelTel By OAuth
//
//	Step 1:verify oauth belong to this account
//
// DelTel By Dynamic Password
//
//	Step 1:send dynamic password to email or tel
//	Step 2:verify email's or tel's dynamic password
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
				ctx := trace.CloneSpan(ctx)
				if e := s.userDao.RedisDelUser(ctx, md["Token-User"]); e != nil {
					log.Error(ctx, "[DelTel] clean redis failed", log.String("operator", md["Token-User"]), log.CError(e))
				}
				s.stop.DoneOne()
			}()
			go func() {
				ctx := trace.CloneSpan(ctx)
				if e := s.userDao.RedisDelUserTelIndex(ctx, olduser.Tel); e != nil {
					log.Error(ctx, "[DelTel] clean redis failed", log.String("tel", olduser.Tel), log.CError(e))
				}
				s.stop.DoneOne()
			}()
		} else {
			s.stop.DoneOne()
			s.stop.DoneOne()
		}
		var final bool
		if olduser.IDCard == "" && olduser.Email == "" && len(olduser.OAuths) == 0 {
			final = true
		}
		return final, nil
	}
	if req.VerifySrcType == "oauth" {
		if req.VerifySrcTypeExtra == "" || req.VerifyDynamicPassword == "" {
			return nil, ecode.ErrReq
		}
		if e := s.userDao.RedisLockTelOP(ctx, md["Token-User"]); e != nil {
			log.Error(ctx, "[DelTel] rate check failed",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, req.VerifyDynamicPassword),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		oauthid, e := util.OAuthVerifyCode(ctx, "DelTel", req.VerifySrcTypeExtra, req.VerifyDynamicPassword)
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		user, e := s.userDao.GetUserByOAuth(ctx, req.VerifySrcTypeExtra, oauthid)
		if e != nil {
			log.Error(ctx, "[DelTel] dao op failed",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, oauthid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.UserID.Hex() != md["Token-User"] {
			log.Error(ctx, "[DelTel] this is not the required oauth",
				log.String("operator", md["Token-User"]),
				log.String(req.VerifySrcTypeExtra, oauthid))
			return nil, ecode.ErrOAuthWrong
		}
		if user.Tel == "" {
			return &api.DelTelResp{Step: "success", Final: user.Email == "" && user.Tel == "" && user.IDCard == "" && len(user.OAuths) == 0}, nil
		}
		if user.BTime != 0 {
			return nil, ecode.ErrBan
		}
		//verify success
		final, e := update()
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		return &api.DelTelResp{Step: "success", Final: final}, nil
	}
	if req.VerifyDynamicPassword != "" {
		//step2
		if e := s.userDao.RedisCheckCode(ctx, md["Token-User"], util.DelTel, req.VerifyDynamicPassword, ""); e != nil {
			log.Error(ctx, "[DelTel] redis op failed",
				log.String("operator", md["Token-User"]),
				log.String("code", req.VerifyDynamicPassword),
				log.CError(e))
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
	if user.IDCard == "" && user.Email == "" && len(user.OAuths) == 0 {
		final = true
	}
	if user.Tel == "" {
		return &api.DelTelResp{Step: "success", Final: final}, nil
	}
	if user.BTime != 0 {
		return nil, ecode.ErrBan
	}

	if req.VerifySrcType == "tel" && user.Tel == "" {
		log.Error(ctx, "[DelTel] missing tel,can't use tel to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.VerifySrcType == "email" && user.Email == "" {
		log.Error(ctx, "[DelTel] missing email,can't use email to receive dynamic password", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}

	if req.VerifySrcType == "email" {
		e = s.sendcode(ctx, "DelTel", "email", user.Email, md["Token-User"], util.DelTel)
	} else {
		e = s.sendcode(ctx, "DelTel", "tel", user.Tel, md["Token-User"], util.DelTel)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if req.VerifySrcType == "email" {
		return &api.DelTelResp{Step: "oldverify", Final: final, Receiver: util.MaskEmail(user.Email)}, nil
	}
	return &api.DelTelResp{Step: "oldverify", Final: final, Receiver: util.MaskTel(user.Tel)}, nil
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
