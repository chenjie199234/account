package base

import (
	"context"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/chenjie199234/account/api"
	"github.com/chenjie199234/account/config"
	basedao "github.com/chenjie199234/account/dao/base"
	"github.com/chenjie199234/account/ecode"
	"github.com/chenjie199234/account/model"
	"github.com/chenjie199234/account/util"

	// "github.com/chenjie199234/Corelib/web"
	// "github.com/chenjie199234/Corelib/cgrpc"
	// "github.com/chenjie199234/Corelib/crpc"
	"github.com/chenjie199234/Corelib/cerror"
	"github.com/chenjie199234/Corelib/cotel"
	"github.com/chenjie199234/Corelib/metadata"
	publicmids "github.com/chenjie199234/Corelib/mids"
	"github.com/chenjie199234/Corelib/util/graceful"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Service subservice for user business
type Service struct {
	stop *graceful.Graceful

	baseDao *basedao.Dao
}

// Start -
func Start() (*Service, error) {
	return &Service{
		stop: graceful.New(),

		baseDao: basedao.NewDao(nil, config.GetRedis("account_redis"), config.GetMongo("account_mongo")),
	}, nil
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
	code, dup, e := s.baseDao.RedisSetCode(ctx, operator, action, src)
	if e != nil {
		s.stop.DoneOne()
		slog.ErrorContext(ctx, "["+callerName+"] redis op failed", slog.String("operator", operator), slog.String(srctype, src), slog.String("error", e.Error()))
		return e
	}
	if dup {
		//if tel's or email's code already send,we jump to verify step
		s.stop.DoneOne()
		return nil
	}
	switch action {
	case util.Login:
		e = s.baseDao.RedisLockLoginDynamic(ctx, operator)
	case util.DelEmail:
		fallthrough
	case util.UpdateEmailStep1:
		e = s.baseDao.RedisLockEmailOP(ctx, operator)
	case util.UpdateEmailStep2:
		//this is controled by step1
	case util.DelTel:
		fallthrough
	case util.UpdateTelStep1:
		e = s.baseDao.RedisLockTelOP(ctx, operator)
	case util.UpdateTelStep2:
		//this is controled by step1
	case util.DelOAuth:
		fallthrough
	case util.UpdateOAuth:
		e = s.baseDao.RedisLockOAuthOP(ctx, operator)
	case util.ResetPassword:
		e = s.baseDao.RedisLockResetPassword(ctx, operator)
	default:
		s.stop.DoneOne()
		return ecode.ErrUnknownAction
	}
	if e != nil {
		slog.ErrorContext(ctx, "["+callerName+"] rate check failed", slog.String("operator", operator), slog.String(srctype, src), slog.String("error", e.Error()))
	} else if srctype == "email" {
		if e = util.SendEmailCode(ctx, src, code, action); e != nil {
			slog.ErrorContext(ctx, "["+callerName+"] send dynamic password failed", slog.String("operator", operator), slog.String(srctype, src), slog.String("error", e.Error()))
		}
	} else if e = util.SendTelCode(ctx, src, code, action); e != nil {
		slog.ErrorContext(ctx, "["+callerName+"] send dynamic password failed", slog.String("operator", operator), slog.String(srctype, src), slog.String("error", e.Error()))
	}
	if e == nil {
		slog.InfoContext(ctx, "["+callerName+"] send dynamic password success", slog.String("operator", operator), slog.String(srctype, src), slog.String("code", code))
		s.stop.DoneOne()
		return nil
	}
	//if rate check failed or send failed,clean redis code
	if ee := s.baseDao.RedisDelCode(ctx, operator, action); ee != nil {
		go func() {
			ctx := cotel.CloneTrace(ctx)
			if ee := s.baseDao.RedisDelCode(ctx, operator, action); ee != nil {
				slog.ErrorContext(ctx, "["+callerName+"] del redis code failed", slog.String("operator", operator), slog.String(srctype, src), slog.String("error", ee.Error()))
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
	switch req.GetSrcType() {
	case "user_id":
		if req.GetSrc() == "" {
			return nil, ecode.ErrReq
		}
		userid, e := bson.ObjectIDFromHex(req.GetSrc())
		if e != nil {
			slog.ErrorContext(ctx, "[BaseInfo] user_id format wrong",
				slog.String("user_id", req.GetSrc()), slog.String("error", e.Error()))
			return nil, ecode.ErrReq
		}
		if user, e = s.baseDao.GetUser(ctx, userid); e != nil {
			slog.ErrorContext(ctx, "[BaseInfo] dao op failed",
				slog.String("user_id", req.GetSrc()), slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	case "tel":
		if req.GetSrc() == "" {
			return nil, ecode.ErrReq
		}
		var e error
		if user, e = s.baseDao.GetUserByTel(ctx, req.GetSrc()); e != nil {
			slog.ErrorContext(ctx, "[BaseInfo] dao op failed",
				slog.String("tel", req.GetSrc()), slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	case "email":
		if req.GetSrc() == "" {
			return nil, ecode.ErrReq
		}
		var e error
		if user, e = s.baseDao.GetUserByEmail(ctx, req.GetSrc()); e != nil {
			slog.ErrorContext(ctx, "[BaseInfo] dao op failed",
				slog.String("email", req.GetSrc()), slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	case "idcard":
		if req.GetSrc() == "" {
			return nil, ecode.ErrReq
		}
		var e error
		if user, e = s.baseDao.GetUserByIDCard(ctx, req.GetSrc()); e != nil {
			slog.ErrorContext(ctx, "[BaseInfo] dao op failed",
				slog.String("idcard", req.GetSrc()), slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	}
	resp := &api.BaseInfoResp{}
	info := &api.BaseInfo{}
	info.SetUserId(user.UserID.Hex())
	info.SetIdcard(user.IDCard)
	info.SetTel(user.Tel)
	info.SetEmail(user.Email)
	info.SetMoney(user.Money)
	info.SetCtime(uint32(user.UserID.Timestamp().Unix()))
	info.SetBan(user.BReason)
	bind := make([]string, 0, len(user.OAuths))
	for oauth := range user.OAuths {
		bind = append(bind, oauth)
	}
	info.SetBindOauths(bind)
	resp.SetInfo(info)
	return resp, nil
}

func (s *Service) Ban(ctx context.Context, req *api.BanReq) (*api.BanResp, error) {
	var userid bson.ObjectID
	switch req.GetSrcType() {
	case "user_id":
		var e error
		userid, e = bson.ObjectIDFromHex(req.GetSrc())
		if e != nil {
			slog.ErrorContext(ctx, "[Ban] user_id format wrong",
				slog.String("user_id", req.GetSrc()), slog.String("error", e.Error()))
			return nil, ecode.ErrReq
		}
		if user, e := s.baseDao.GetUser(ctx, userid); e != nil {
			slog.ErrorContext(ctx, "[Ban] dao op failed",
				slog.String("user_id", req.GetSrc()), slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		} else if user.BReason == req.GetReason() {
			return &api.BanResp{}, nil
		}
	case "tel":
		if user, e := s.baseDao.GetUserByTel(ctx, req.GetSrc()); e != nil {
			slog.ErrorContext(ctx, "[Ban] dao op failed",
				slog.String("tel", req.GetSrc()), slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		} else if user.BReason == req.GetReason() {
			return &api.BanResp{}, nil
		} else {
			userid = user.UserID
		}
	case "email":
		if user, e := s.baseDao.GetUserByEmail(ctx, req.GetSrc()); e != nil {
			slog.ErrorContext(ctx, "[Ban] dao op failed",
				slog.String("email", req.GetSrc()), slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		} else if user.BReason == req.GetReason() {
			return &api.BanResp{}, nil
		} else {
			userid = user.UserID
		}
	case "idcard":
		if user, e := s.baseDao.GetUserByIDCard(ctx, req.GetSrc()); e != nil {
			slog.ErrorContext(ctx, "[Ban] dao op failed",
				slog.String("idcard", req.GetSrc()), slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		} else if user.BReason == req.GetReason() {
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
	if e := s.baseDao.MongoBanUser(ctx, userid, req.GetReason()); e != nil {
		s.stop.DoneOne()
		slog.ErrorContext(ctx, "[Ban] db op failed",
			slog.String(req.GetSrcType(), req.GetSrc()), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if req.GetSrcType() == "user_id" {
		slog.InfoContext(ctx, "[Ban] success", slog.String(req.GetSrcType(), req.GetSrc()))
	} else {
		slog.InfoContext(ctx, "[Ban] success", slog.String(req.GetSrcType(), req.GetSrc()), slog.String("user_id", userid.Hex()))
	}
	go func() {
		ctx := cotel.CloneTrace(ctx)
		if e := s.baseDao.RedisDelUser(ctx, userid.Hex()); e != nil {
			if req.GetSrcType() == "user_id" {
				slog.ErrorContext(ctx, "[Ban] clean redis failed",
					slog.String(req.GetSrcType(), req.GetSrc()), slog.String("error", e.Error()))
			} else {
				slog.ErrorContext(ctx, "[Ban] clean redis failed",
					slog.String(req.GetSrcType(), req.GetSrc()), slog.String("user_id", userid.Hex()), slog.String("error", e.Error()))
			}
		}
		s.stop.DoneOne()
	}()
	return &api.BanResp{}, nil
}

func (s *Service) Unban(ctx context.Context, req *api.UnbanReq) (*api.UnbanResp, error) {
	var userid bson.ObjectID
	switch req.GetSrcType() {
	case "user_id":
		var e error
		userid, e = bson.ObjectIDFromHex(req.GetSrc())
		if e != nil {
			slog.ErrorContext(ctx, "[Unban] user_id format wrong",
				slog.String("user_id", req.GetSrc()), slog.String("error", e.Error()))
			return nil, ecode.ErrReq
		}
		if user, e := s.baseDao.GetUser(ctx, userid); e != nil {
			slog.ErrorContext(ctx, "[Unban] dao op failed",
				slog.String("user_id", req.GetSrc()), slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		} else if user.BTime == 0 {
			return &api.UnbanResp{}, nil
		}
	case "tel":
		if user, e := s.baseDao.GetUserByTel(ctx, req.GetSrc()); e != nil {
			slog.ErrorContext(ctx, "[Unban] dao op failed",
				slog.String("tel", req.GetSrc()), slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		} else if user.BTime == 0 {
			return &api.UnbanResp{}, nil
		} else {
			userid = user.UserID
		}
	case "email":
		if user, e := s.baseDao.GetUserByEmail(ctx, req.GetSrc()); e != nil {
			slog.ErrorContext(ctx, "[Unban] dao op failed",
				slog.String("email", req.GetSrc()), slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		} else if user.BTime == 0 {
			return &api.UnbanResp{}, nil
		} else {
			userid = user.UserID
		}
	case "idcard":
		if user, e := s.baseDao.GetUserByIDCard(ctx, req.GetSrc()); e != nil {
			slog.ErrorContext(ctx, "[Unban] dao op failed",
				slog.String("idcard", req.GetSrc()), slog.String("error", e.Error()))
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
	if e := s.baseDao.MongoUnBanUser(ctx, userid); e != nil {
		s.stop.DoneOne()
		slog.ErrorContext(ctx, "[Unban] db op failed",
			slog.String(req.GetSrcType(), req.GetSrc()), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if req.GetSrcType() == "user_id" {
		slog.InfoContext(ctx, "[Unban] success", slog.String(req.GetSrcType(), req.GetSrc()))
	} else {
		slog.InfoContext(ctx, "[Unban] success", slog.String(req.GetSrcType(), req.GetSrc()), slog.String("user_id", userid.Hex()))
	}
	go func() {
		ctx := cotel.CloneTrace(ctx)
		if e := s.baseDao.RedisDelUser(ctx, userid.Hex()); e != nil {
			if req.GetSrcType() == "user_id" {
				slog.ErrorContext(ctx, "[Unban] clean redis failed",
					slog.String("user_id", userid.Hex()), slog.String("error", e.Error()))
			} else {
				slog.ErrorContext(ctx, "[Unban] clean redis failed",
					slog.String(req.GetSrcType(), req.GetSrc()), slog.String("user_id", userid.Hex()), slog.String("error", e.Error()))
			}
		}
		s.stop.DoneOne()
	}()
	return &api.UnbanResp{}, nil
}

func (s *Service) GetOauthUrl(ctx context.Context, req *api.GetOauthUrlReq) (*api.GetOauthUrlResp, error) {
	switch req.GetOauthServiceName() {
	case "wechat":
		if config.AC.Service.WeChatOauth2 == "" {
			return nil, ecode.ErrOAuthUnknown
		}
		resp := &api.GetOauthUrlResp{}
		resp.SetUrl(config.AC.Service.WeChatOauth2)
		return resp, nil
	default:
		return nil, ecode.ErrOAuthUnknown
	}
}
func (s *Service) Login(ctx context.Context, req *api.LoginReq) (*api.LoginResp, error) {
	if req.GetPasswordType() == "static" && req.GetPassword() == "" {
		slog.ErrorContext(ctx, "[Login] empty static password", slog.String(req.GetSrcType(), req.GetSrcTypeExtra()))
		return nil, ecode.ErrReq
	}
	var user *model.User
	var e error
	switch req.GetSrcType() {
	case "idcard":
		if req.GetPasswordType() == "dynamic" {
			slog.ErrorContext(ctx, "[Login] idcard can't use dynamic password")
			return nil, ecode.ErrReq
		}
		//static
		if user, e = s.baseDao.GetUserByIDCard(ctx, req.GetSrcTypeExtra()); e != nil {
			slog.ErrorContext(ctx, "[Login] dao op failed",
				slog.String("idcard", req.GetSrcTypeExtra()), slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	case "tel":
		if req.GetPasswordType() == "static" {
			if user, e = s.baseDao.GetUserByTel(ctx, req.GetSrcTypeExtra()); e != nil {
				slog.ErrorContext(ctx, "[Login] dao op failed", slog.String("tel", req.GetSrcTypeExtra()), slog.String("error", e.Error()))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
		} else if req.GetPassword() == "" {
			if e = s.sendcode(ctx, "Login", req.GetSrcType(), req.GetSrcTypeExtra(), req.GetSrcTypeExtra(), util.Login); e != nil {
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			resp := &api.LoginResp{}
			resp.SetStep("verify")
			return resp, nil
		} else {
			if e = s.baseDao.RedisCheckCode(ctx, req.GetSrcTypeExtra(), util.Login, req.GetPassword(), ""); e != nil {
				slog.ErrorContext(ctx, "[Login] redis op failed",
					slog.String("tel", req.GetSrcTypeExtra()), slog.String("code", req.GetPassword()), slog.String("error", e.Error()))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			if user, e = s.baseDao.GetOrCreateUserByTel(ctx, req.GetSrcTypeExtra()); e != nil {
				slog.ErrorContext(ctx, "[Login] dao op failed",
					slog.String("tel", req.GetSrcTypeExtra()), slog.String("error", e.Error()))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
		}
	case "email":
		emailServices := config.AC.Service.SupportEmailService
		support := false
		low := strings.ToLower(req.GetSrcTypeExtra())
		for _, v := range emailServices {
			if strings.HasSuffix(low, v) {
				support = true
				break
			}
		}
		if !support {
			return nil, ecode.ErrUnsupportedEmailService
		}
		if req.GetPasswordType() == "static" {
			if user, e = s.baseDao.GetUserByEmail(ctx, req.GetSrcTypeExtra()); e != nil {
				slog.ErrorContext(ctx, "[Login] dao op failed", slog.String("email", req.GetSrcTypeExtra()), slog.String("error", e.Error()))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
		} else if req.GetPassword() == "" {
			if e := s.sendcode(ctx, "Login", req.GetSrcType(), req.GetSrcTypeExtra(), req.GetSrcTypeExtra(), util.Login); e != nil {
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			resp := &api.LoginResp{}
			resp.SetStep("verify")
			return resp, nil
		} else {
			if e = s.baseDao.RedisCheckCode(ctx, req.GetSrcTypeExtra(), util.Login, req.GetPassword(), ""); e != nil {
				slog.ErrorContext(ctx, "[Login] redis op failed",
					slog.String("email", req.GetSrcTypeExtra()), slog.String("code", req.GetPassword()), slog.String("error", e.Error()))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			if user, e = s.baseDao.GetOrCreateUserByEmail(ctx, req.GetSrcTypeExtra()); e != nil {
				slog.ErrorContext(ctx, "[Login] dao op failed",
					slog.String("email", req.GetSrcTypeExtra()), slog.String("error", e.Error()))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
		}
	case "oauth":
		if req.GetPasswordType() == "static" {
			slog.ErrorContext(ctx, "[Login] oauth can't use static password")
			return nil, ecode.ErrReq
		}
		if req.GetPassword() == "" {
			return nil, ecode.ErrReq
		}
		var oauthid string
		switch req.GetSrcTypeExtra() {
		case "wechat":
			c := config.AC.Service
			oauthid, e = util.OAuthWeChatVerifyCode(ctx, c.WeChatAppID, c.WeChatAppSecret, req.GetPassword())
		default:
			return nil, ecode.ErrOAuthUnknown
		}
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user, e = s.baseDao.GetOrCreateUserByOAuth(ctx, req.GetSrcTypeExtra(), oauthid); e != nil {
			slog.ErrorContext(ctx, "[Login] dao op failed",
				slog.String(req.GetSrcTypeExtra(), oauthid), slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	}
	if user.BTime != 0 {
		return nil, ecode.ErrBan
	}
	if req.GetPasswordType() == "static" {
		if e := util.SignCheck(req.GetPassword(), user.Password); e != nil {
			if e == ecode.ErrSignCheckFailed {
				e = ecode.ErrPasswordWrong
			}
			slog.ErrorContext(ctx, "[Login] sign check failed",
				slog.String(req.GetSrcType(), req.GetSrcTypeExtra()), slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	}
	//TODO set the puber
	resp := &api.LoginResp{}
	info := &api.BaseInfo{}
	info.SetUserId(user.UserID.Hex())
	info.SetIdcard(util.MaskIDCard(user.IDCard))
	info.SetTel(util.MaskTel(user.Tel))
	info.SetEmail(util.MaskEmail(user.Email))
	info.SetCtime(uint32(user.UserID.Timestamp().Unix()))
	info.SetMoney(user.Money)
	bind := make([]string, 0, len(user.OAuths))
	for oauth := range user.OAuths {
		bind = append(bind, oauth)
	}
	info.SetBindOauths(bind)
	resp.SetInfo(info)
	if req.GetPasswordType() == "dynamic" &&
		(req.GetSrcType() == "email" || req.GetSrcType() == "tel") &&
		util.SignCheck("", user.Password) == nil {
		resp.SetStep("password")
	} else {
		resp.SetStep("success")
	}
	resp.SetToken(publicmids.MakeToken(ctx, "", *config.EC.DeployEnv, *config.EC.RunEnv, user.UserID.Hex(), "", config.AC.Service.TokenExpire.StdDuration()))
	resp.SetTokenexpire(uint64(time.Now().Add(config.AC.Service.TokenExpire.StdDuration() - time.Second).UnixNano()))
	slog.InfoContext(ctx, "[Login] success", slog.String("operator", user.UserID.Hex()))
	return resp, nil
}
func (s *Service) TemporaryToken(ctx context.Context, req *api.TemporaryTokenReq) (*api.TemporaryTokenResp, error) {
	md := metadata.GetMetadata(ctx)
	resp := &api.TemporaryTokenResp{}
	resp.SetToken(publicmids.MakeToken(ctx, md["Token-Puber"], md["Token-DeployEnv"], md["Token-RunEnv"], md["Token-User"], md["Token-Data"], time.Minute))
	resp.SetTokenexpire(uint64(time.Now().Add(time.Minute - time.Second).UnixNano()))
	return resp, nil
}
func (s *Service) SelfInfo(ctx context.Context, req *api.SelfInfoReq) (*api.SelfInfoResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[SelfInfo] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
	}
	var user *model.User
	if user, e = s.baseDao.GetUser(ctx, operator); e != nil {
		slog.ErrorContext(ctx, "[SelfInfo] dao op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if user.BTime != 0 {
		return nil, ecode.ErrBan
	}
	resp := &api.SelfInfoResp{}
	info := &api.BaseInfo{}
	info.SetUserId(user.UserID.Hex())
	info.SetIdcard(util.MaskIDCard(user.IDCard))
	info.SetTel(util.MaskTel(user.Tel))
	info.SetEmail(util.MaskEmail(user.Email))
	info.SetMoney(user.Money)
	info.SetCtime(uint32(user.UserID.Timestamp().Unix()))
	info.SetBan(user.BReason)
	bind := make([]string, 0, len(user.OAuths))
	for oauth := range user.OAuths {
		bind = append(bind, oauth)
	}
	info.SetBindOauths(bind)
	resp.SetInfo(info)
	return resp, nil
}
func (s *Service) UpdateStaticPassword(ctx context.Context, req *api.UpdateStaticPasswordReq) (*api.UpdateStaticPasswordResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[UpdateStaticPassword] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
	}
	if user, e := s.baseDao.GetUser(ctx, operator); e != nil {
		slog.ErrorContext(ctx, "[UpdateStaticPassword] dao op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
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
	if e := s.baseDao.RedisLockUpdatePassword(ctx, md["Token-User"]); e != nil {
		s.stop.DoneOne()
		slog.ErrorContext(ctx, "[UpdateStaticPassword] redis op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}

	if e := s.baseDao.MongoUpdateUserPassword(ctx, operator, req.GetOldStaticPassword(), req.GetNewStaticPassword()); e != nil {
		s.stop.DoneOne()
		slog.ErrorContext(ctx, "[UpdateStaticPassword] db op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[UpdateStaticPassword] success", slog.String("operator", md["Token-User"]))
	go func() {
		ctx := cotel.CloneTrace(ctx)
		if e := s.baseDao.RedisDelUser(ctx, md["Token-User"]); e != nil {
			slog.ErrorContext(ctx, "[UpdateStaticPassword] clean redis failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
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
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[ResetStaticPassword] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
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
		if e := s.baseDao.MongoResetUserPassword(ctx, operator); e != nil {
			slog.ErrorContext(ctx, "[ResetStaticPassword] db op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
			s.stop.DoneOne()
			return e
		}
		slog.InfoContext(ctx, "[ResetStaticPassword] success", slog.String("operator", md["Token-User"]))
		go func() {
			ctx := cotel.CloneTrace(ctx)
			if e := s.baseDao.RedisDelUser(ctx, md["Token-User"]); e != nil {
				slog.ErrorContext(ctx, "[ResetStaticPassword] clean redis failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
			}
			s.stop.DoneOne()
		}()
		return nil
	}
	if req.GetVerifySrcType() == "oauth" {
		if req.GetVerifySrcTypeExtra() == "" || req.GetVerifyDynamicPassword() == "" {
			return nil, ecode.ErrReq
		}
		if e := s.baseDao.RedisLockResetPassword(ctx, md["Token-User"]); e != nil {
			slog.ErrorContext(ctx, "[ResetStaticPassword] rate check failed",
				slog.String("operator", md["Token-User"]),
				slog.String(req.GetVerifySrcTypeExtra(), req.GetVerifyDynamicPassword()),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		var oauthid string
		switch req.GetVerifySrcTypeExtra() {
		case "wechat":
			c := config.AC.Service
			oauthid, e = util.OAuthWeChatVerifyCode(ctx, c.WeChatAppID, c.WeChatAppSecret, req.GetVerifyDynamicPassword())
		default:
			return nil, ecode.ErrOAuthUnknown
		}
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		user, e := s.baseDao.GetUserByOAuth(ctx, req.GetVerifySrcTypeExtra(), oauthid)
		if e != nil {
			slog.ErrorContext(ctx, "[ResetStaticPassword] dao op failed",
				slog.String("operator", md["Token-User"]),
				slog.String(req.GetVerifySrcTypeExtra(), oauthid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.UserID.Hex() != md["Token-User"] {
			slog.ErrorContext(ctx, "[ResetStaticPassword] this is not the required oauth",
				slog.String("operator", md["Token-User"]),
				slog.String(req.GetVerifySrcTypeExtra(), oauthid))
			return nil, ecode.ErrOAuthWrong
		}
		if util.SignCheck("", user.Password) == nil {
			resp := &api.ResetStaticPasswordResp{}
			resp.SetStep("success")
			return resp, nil
		}
		if user.BTime != 0 {
			return nil, ecode.ErrBan
		}
		//verify success
		if e := update(); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		resp := &api.ResetStaticPasswordResp{}
		resp.SetStep("success")
		return resp, nil
	}
	if req.GetVerifyDynamicPassword() != "" {
		//step2
		if e := s.baseDao.RedisCheckCode(ctx, md["Token-User"], util.ResetPassword, req.GetVerifyDynamicPassword(), ""); e != nil {
			slog.ErrorContext(ctx, "[ResetStaticPassword] redis op failed",
				slog.String("operator", md["Token-User"]),
				slog.String("code", req.GetVerifyDynamicPassword()),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		//verify success
		if e := update(); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		resp := &api.ResetStaticPasswordResp{}
		resp.SetStep("success")
		return resp, nil
	}
	//step1
	user, e := s.baseDao.GetUser(ctx, operator)
	if e != nil {
		slog.ErrorContext(ctx, "[ResetStaticPassword] dao op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if util.SignCheck("", user.Password) == nil {
		resp := &api.ResetStaticPasswordResp{}
		resp.SetStep("success")
		return resp, nil
	}
	if user.BTime != 0 {
		return nil, ecode.ErrBan
	}

	if req.GetVerifySrcType() == "tel" && user.Tel == "" {
		slog.ErrorContext(ctx, "[ResetStaticPassword] missing tel,can't use tel to receive dynamic password", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.GetVerifySrcType() == "email" && user.Email == "" {
		slog.ErrorContext(ctx, "[ResetStaticPassword] missing email,can't use email to receive dynamic password", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}

	if req.GetVerifySrcType() == "email" {
		e = s.sendcode(ctx, "ResetStaticPassword", "email", user.Email, md["Token-User"], util.ResetPassword)
	} else {
		e = s.sendcode(ctx, "ResetStaticPassword", "tel", user.Tel, md["Token-User"], util.ResetPassword)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.ResetStaticPasswordResp{}
	resp.SetStep("oldverify")
	if req.GetVerifySrcType() == "email" {
		resp.SetReceiver(util.MaskEmail(user.Email))
	} else {
		resp.SetReceiver(util.MaskTel(user.Tel))
	}
	return resp, nil
}

func (s *Service) IdcardDuplicateCheck(ctx context.Context, req *api.IdcardDuplicateCheckReq) (*api.IdcardDuplicateCheckResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[IdcardDuplicateCheck] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
	}
	if user, e := s.baseDao.GetUser(ctx, operator); e != nil {
		slog.ErrorContext(ctx, "[IdcardDuplicateCheck] dao op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	} else if user.BTime != 0 {
		return nil, ecode.ErrBan
	}
	//redis lock
	if e := s.baseDao.RedisLockDuplicateCheck(ctx, "idcard", md["Token-User"]); e != nil {
		slog.ErrorContext(ctx, "[IdcardDuplicateCheck] redis op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	userid, e := s.baseDao.GetUserIDCardIndex(ctx, req.GetIdcard())
	if e != nil && e != ecode.ErrUserNotExist {
		slog.ErrorContext(ctx, "[IdcardDuplicateCheck] dao op failed",
			slog.String("operator", md["Token-User"]), slog.String("idcard", req.GetIdcard()), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.IdcardDuplicateCheckResp{}
	resp.SetDuplicate(userid != "")
	return resp, nil
}

func (s *Service) SetIdcard(ctx context.Context, req *api.SetIdcardReq) (*api.SetIdcardResp, error) {
	match, _ := regexp.MatchString(`^[1-9]\d{5}(19|20)\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[Xx\d]$`, req.GetIdcard())
	if !match {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[SetIdcard] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
	}
	user, e := s.baseDao.GetUser(ctx, operator)
	if e != nil {
		slog.ErrorContext(ctx, "[SetIdcard] dao op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if user.IDCard == req.GetIdcard() {
		return &api.SetIdcardResp{}, nil
	}
	if user.IDCard != "" {
		return nil, ecode.ErrIDCardAlreadySetted
	}
	if user.BTime != 0 {
		return nil, ecode.ErrBan
	}
	//update db and clean redis is async
	//the service's rolling update may happened between update db and clean redis
	//so we need to make this not happened
	if e := s.stop.Add(2); e != nil {
		if e == graceful.ErrClosing {
			return nil, cerror.ErrServerClosing
		}
		return nil, ecode.ErrBusy
	}

	if _, e = s.baseDao.MongoUpdateUserIDCard(ctx, operator, req.GetIdcard()); e != nil {
		s.stop.DoneOne()
		s.stop.DoneOne()
		slog.ErrorContext(ctx, "[UpdateIdcard] db op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, e
	}
	slog.InfoContext(ctx, "[UpdateIdcard] success", slog.String("operator", md["Token-User"]), slog.String("new_idcard", req.GetIdcard()))
	go func() {
		ctx := cotel.CloneTrace(ctx)
		if e := s.baseDao.RedisDelUser(ctx, md["Token-User"]); e != nil {
			slog.ErrorContext(ctx, "[UpdateIdcard] clean redis failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		}
		s.stop.DoneOne()
	}()
	go func() {
		ctx := cotel.CloneTrace(ctx)
		if e := s.baseDao.RedisDelUserIDCardIndex(ctx, req.GetIdcard()); e != nil {
			slog.ErrorContext(ctx, "[UpdateIdcard] clean redis failed", slog.String("idcard", req.GetIdcard()), slog.String("error", e.Error()))
		}
		s.stop.DoneOne()
	}()
	return &api.SetIdcardResp{}, nil
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
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[UpdateOauth] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
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
		if olduser, e = s.baseDao.MongoUpdateUserOAuth(ctx, operator, req.GetNewOauthServiceName(), newoauthid); e != nil {
			s.stop.DoneOne()
			s.stop.DoneOne()
			s.stop.DoneOne()
			slog.ErrorContext(ctx, "[UpdateOauth] db op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
			return e
		}
		slog.InfoContext(ctx, "[UpdateOauth] success", slog.String("operator", md["Token-User"]), slog.String(req.GetNewOauthServiceName(), newoauthid))
		oldoauthid := olduser.OAuths[req.GetNewOauthServiceName()]
		if oldoauthid != newoauthid {
			go func() {
				ctx := cotel.CloneTrace(ctx)
				if e := s.baseDao.RedisDelUser(ctx, md["Token-User"]); e != nil {
					slog.ErrorContext(ctx, "[UpdateOauth] clean redis failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
				}
				s.stop.DoneOne()
			}()
			go func() {
				if oldoauthid != "" {
					ctx := cotel.CloneTrace(ctx)
					if e := s.baseDao.RedisDelUserOAuthIndex(ctx, req.GetNewOauthServiceName(), oldoauthid); e != nil {
						slog.ErrorContext(ctx, "[UpdateOauth] clean redis failed", slog.String(req.GetNewOauthServiceName(), oldoauthid), slog.String("error", e.Error()))
					}
				}
				s.stop.DoneOne()
			}()
			go func() {
				if newoauthid != "" {
					ctx := cotel.CloneTrace(ctx)
					if e := s.baseDao.RedisDelUserOAuthIndex(ctx, req.GetNewOauthServiceName(), newoauthid); e != nil {
						slog.ErrorContext(ctx, "[UpdateOauth] clean redis failed", slog.String(req.GetNewOauthServiceName(), newoauthid), slog.String("error", e.Error()))
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
	if req.GetVerifySrcType() == "oauth" {
		if req.GetVerifySrcTypeExtra() == "" ||
			req.GetVerifyDynamicPassword() == "" ||
			req.GetNewOauthServiceName() == "" ||
			req.GetNewOauthDynamicPassword() == "" {
			return nil, ecode.ErrReq
		}
		if e := s.baseDao.RedisLockOAuthOP(ctx, md["Token-User"]); e != nil {
			slog.ErrorContext(ctx, "[UpdateOauth] rate check failed",
				slog.String("operator", md["Token-User"]),
				slog.String(req.GetVerifySrcTypeExtra(), req.GetVerifyDynamicPassword()),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		var oauthid string
		switch req.GetVerifySrcTypeExtra() {
		case "wechat":
			c := config.AC.Service
			oauthid, e = util.OAuthWeChatVerifyCode(ctx, c.WeChatAppID, c.WeChatAppSecret, req.GetVerifyDynamicPassword())
		default:
			return nil, ecode.ErrOAuthUnknown
		}
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		user, e := s.baseDao.GetUserByOAuth(ctx, req.GetVerifySrcTypeExtra(), oauthid)
		if e != nil {
			slog.ErrorContext(ctx, "[UpdateOauth] dao op failed",
				slog.String("operator", md["Token-User"]),
				slog.String(req.GetVerifySrcTypeExtra(), oauthid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.UserID.Hex() != md["Token-User"] {
			slog.ErrorContext(ctx, "[UpdateOauth] this is not the required oauth",
				slog.String("operator", md["Token-User"]),
				slog.String(req.GetVerifySrcTypeExtra(), oauthid))
			return nil, ecode.ErrOAuthWrong
		}
		if user.BTime != 0 {
			return nil, ecode.ErrBan
		}
		//verify success
		//get the new oauth
		switch req.GetNewOauthServiceName() {
		case "wechat":
			c := config.AC.Service
			oauthid, e = util.OAuthWeChatVerifyCode(ctx, c.WeChatAppID, c.WeChatAppSecret, req.GetNewOauthDynamicPassword())
		default:
			return nil, ecode.ErrOAuthUnknown
		}
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.OAuths[req.GetNewOauthServiceName()] == oauthid {
			resp := &api.UpdateOauthResp{}
			resp.SetStep("success")
			return resp, nil
		}
		if e := update(oauthid); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		resp := &api.UpdateOauthResp{}
		resp.SetStep("success")
		return resp, nil
	}
	if req.GetVerifyDynamicPassword() != "" {
		if req.GetNewOauthServiceName() == "" || req.GetNewOauthDynamicPassword() == "" {
			return nil, ecode.ErrReq
		}
		if e := s.baseDao.RedisCheckCode(ctx, md["Token-User"], util.UpdateOAuth, req.GetVerifyDynamicPassword(), ""); e != nil {
			slog.ErrorContext(ctx, "[UpdateOauth] redis op failed",
				slog.String("operator", md["Token-User"]),
				slog.String("code", req.GetVerifyDynamicPassword()),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		//verify success
		//get the new oauth
		var oauthid string
		switch req.GetNewOauthServiceName() {
		case "wechat":
			c := config.AC.Service
			oauthid, e = util.OAuthWeChatVerifyCode(ctx, c.WeChatAppID, c.WeChatAppSecret, req.GetNewOauthDynamicPassword())
		default:
			return nil, ecode.ErrOAuthUnknown
		}
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if e := update(oauthid); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		resp := &api.UpdateOauthResp{}
		resp.SetStep("success")
		return resp, nil
	}
	//send dynamic password
	user, e := s.baseDao.GetUser(ctx, operator)
	if e != nil {
		slog.ErrorContext(ctx, "[UpdateOauth] dao op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if user.BTime != 0 {
		return nil, ecode.ErrBan
	}

	if req.GetVerifySrcType() == "email" && user.Email == "" {
		slog.ErrorContext(ctx, "[UpdateOauth] missing email,can't use email to receive dynamic password", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.GetVerifySrcType() == "tel" && user.Tel == "" {
		slog.ErrorContext(ctx, "[UpdateOauth] missing tel,can't use tel to receive dynamic password", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}

	if req.GetVerifySrcType() == "email" {
		e = s.sendcode(ctx, "UpdateOauth", "email", user.Email, md["Token-User"], util.UpdateOAuth)
	} else {
		e = s.sendcode(ctx, "UpdateOauth", "tel", user.Tel, md["Token-User"], util.UpdateOAuth)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.UpdateOauthResp{}
	resp.SetStep("oldverify")
	if req.GetVerifySrcType() == "email" {
		resp.SetReceiver(util.MaskEmail(user.Email))
	} else {
		resp.SetReceiver(util.MaskTel(user.Tel))
	}
	return resp, nil
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
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[DelOauth] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
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
		if olduser, e = s.baseDao.MongoUpdateUserOAuth(ctx, operator, req.GetDelOauthServiceName(), ""); e != nil {
			s.stop.DoneOne()
			s.stop.DoneOne()
			slog.ErrorContext(ctx, "[DelOauth] db op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
			return false, e
		}
		if oauthid := olduser.OAuths[req.GetDelOauthServiceName()]; oauthid != "" {
			go func() {
				ctx := cotel.CloneTrace(ctx)
				if e := s.baseDao.RedisDelUser(ctx, md["Token-User"]); e != nil {
					slog.ErrorContext(ctx, "[DelOauth] clean redis failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
				}
				s.stop.DoneOne()
			}()
			go func() {
				ctx := cotel.CloneTrace(ctx)
				if e := s.baseDao.RedisDelUserOAuthIndex(ctx, req.GetDelOauthServiceName(), oauthid); e != nil {
					slog.ErrorContext(ctx, "[DelOauth] clean redis failed",
						slog.String(req.GetDelOauthServiceName(), oauthid), slog.String("error", e.Error()))
				}
				s.stop.DoneOne()
			}()
		} else {
			s.stop.DoneOne()
			s.stop.DoneOne()
		}
		var final bool
		if olduser.Email == "" &&
			olduser.IDCard == "" &&
			olduser.Tel == "" &&
			len(olduser.OAuths) == 1 &&
			olduser.OAuths[req.GetDelOauthServiceName()] != "" {
			final = true
		}
		slog.InfoContext(ctx, "[DelOauth] success",
			slog.String("operator", md["Token-User"]), slog.String("oauth", req.GetDelOauthServiceName()), slog.Bool("final", final))
		return final, nil
	}
	if req.GetVerifySrcType() == "oauth" {
		if req.GetVerifySrcTypeExtra() == "" || req.GetVerifyDynamicPassword() == "" {
			return nil, ecode.ErrReq
		}
		if e := s.baseDao.RedisLockOAuthOP(ctx, md["Token-User"]); e != nil {
			slog.ErrorContext(ctx, "[DelOauth] rate check failed",
				slog.String("operator", md["Token-User"]),
				slog.String(req.GetVerifySrcTypeExtra(), req.GetVerifyDynamicPassword()),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		var oauthid string
		switch req.GetVerifySrcTypeExtra() {
		case "wechat":
			c := config.AC.Service
			oauthid, e = util.OAuthWeChatVerifyCode(ctx, c.WeChatAppID, c.WeChatAppSecret, req.GetVerifyDynamicPassword())
		default:
			return nil, ecode.ErrOAuthUnknown
		}
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		user, e := s.baseDao.GetUserByOAuth(ctx, req.GetVerifySrcTypeExtra(), oauthid)
		if e != nil {
			slog.ErrorContext(ctx, "[DelOauth] dao op failed",
				slog.String("operator", md["Token-User"]),
				slog.String(req.GetVerifySrcTypeExtra(), oauthid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.UserID.Hex() != md["Token-User"] {
			slog.ErrorContext(ctx, "[DelOauth] this is not the required oauth",
				slog.String("operator", md["Token-User"]),
				slog.String(req.GetVerifySrcTypeExtra(), oauthid))
			return nil, ecode.ErrOAuthWrong
		}
		if _, ok := user.OAuths[req.GetDelOauthServiceName()]; !ok {
			resp := &api.DelOauthResp{}
			resp.SetStep("success")
			resp.SetFinal(user.Email == "" && user.IDCard == "" && user.Tel == "" && len(user.OAuths) == 0)
			return resp, nil
		}
		if user.BTime != 0 {
			return nil, ecode.ErrBan
		}
		//verify success
		final, e := update()
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		resp := &api.DelOauthResp{}
		resp.SetStep("success")
		resp.SetFinal(final)
		return resp, nil
	}
	if req.GetVerifyDynamicPassword() != "" {
		//step2
		if e := s.baseDao.RedisCheckCode(ctx, md["Token-User"], util.DelOAuth, req.GetVerifyDynamicPassword(), ""); e != nil {
			slog.ErrorContext(ctx, "[DelOauth] redis op failed",
				slog.String("operator", md["Token-User"]),
				slog.String("code", req.GetVerifyDynamicPassword()),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		//verify success
		final, e := update()
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		resp := &api.DelOauthResp{}
		resp.SetStep("success")
		resp.SetFinal(final)
		return resp, nil
	}
	//step1
	user, e := s.baseDao.GetUser(ctx, operator)
	if e != nil {
		slog.ErrorContext(ctx, "[DelOauth] dao op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	var final bool
	if user.Email == "" && user.IDCard == "" && user.Tel == "" && len(user.OAuths) == 0 {
		final = true
	}
	if user.OAuths[req.GetDelOauthServiceName()] == "" {
		resp := &api.DelOauthResp{}
		resp.SetStep("success")
		resp.SetFinal(final)
		return resp, nil
	}
	if user.BTime != 0 {
		return nil, ecode.ErrBan
	}

	if req.GetVerifySrcType() == "tel" && user.Tel == "" {
		slog.ErrorContext(ctx, "[DelOauth] missing tel,can't use tel to receive dynamic password", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.GetVerifySrcType() == "email" && user.Email == "" {
		slog.ErrorContext(ctx, "[DelOauth] missing email,can't use email to receive dynamic password", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}

	if req.GetVerifySrcType() == "email" {
		e = s.sendcode(ctx, "DelOauth", "email", user.Email, md["Token-User"], util.DelOAuth)
	} else {
		e = s.sendcode(ctx, "DelOauth", "tel", user.Tel, md["Token-User"], util.DelOAuth)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.DelOauthResp{}
	resp.SetStep("oldverify")
	resp.SetFinal(final)
	if req.GetVerifySrcType() == "email" {
		resp.SetReceiver(util.MaskEmail(user.Email))
	} else {
		resp.SetReceiver(util.MaskTel(user.Tel))
	}
	return resp, nil
}

func (s *Service) EmailDuplicateCheck(ctx context.Context, req *api.EmailDuplicateCheckReq) (*api.EmailDuplicateCheckResp, error) {
	emailServices := config.AC.Service.SupportEmailService
	support := false
	low := strings.ToLower(req.GetEmail())
	for _, v := range emailServices {
		if strings.HasSuffix(low, v) {
			support = true
			break
		}
	}
	if !support {
		return nil, ecode.ErrUnsupportedEmailService
	}
	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[EmailDuplicateCheck] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
	}
	if user, e := s.baseDao.GetUser(ctx, operator); e != nil {
		slog.ErrorContext(ctx, "[EmailDuplicateCheck] dao op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	} else if user.BTime != 0 {
		return nil, ecode.ErrBan
	}
	//redis lock
	if e := s.baseDao.RedisLockDuplicateCheck(ctx, "email", md["Token-User"]); e != nil {
		slog.ErrorContext(ctx, "[EmailDuplicateCheck] redis op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	userid, e := s.baseDao.GetUserEmailIndex(ctx, req.GetEmail())
	if e != nil && e != ecode.ErrUserNotExist {
		slog.ErrorContext(ctx, "[EmailDuplicateCheck] dao op failed",
			slog.String("operator", md["Token-User"]), slog.String("email", req.GetEmail()), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.EmailDuplicateCheckResp{}
	resp.SetDuplicate(userid != "")
	return resp, nil
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
	emailServices := config.AC.Service.SupportEmailService
	support := false
	low := strings.ToLower(req.GetNewEmail())
	for _, v := range emailServices {
		if strings.HasSuffix(low, v) {
			support = true
			break
		}
	}
	if !support {
		return nil, ecode.ErrUnsupportedEmailService
	}
	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[UpdateEmail] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
	}
	if req.GetNewEmailDynamicPassword() != "" {
		//step final
		e = s.baseDao.RedisCheckCode(ctx, md["Token-User"], util.UpdateEmailStep2, req.GetNewEmailDynamicPassword(), req.GetNewEmail())
		if e != nil {
			slog.ErrorContext(ctx, "[UpdateEmail] redis op failed",
				slog.String("operator", md["Token-User"]),
				slog.String("code", req.GetNewEmailDynamicPassword()),
				slog.String("error", e.Error()))
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
		if olduser, e = s.baseDao.MongoUpdateUserEmail(ctx, operator, req.GetNewEmail()); e != nil {
			s.stop.DoneOne()
			s.stop.DoneOne()
			s.stop.DoneOne()
			slog.ErrorContext(ctx, "[UpdateEmail] db op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		slog.InfoContext(ctx, "[UpdateEmail] success",
			slog.String("operator", md["Token-User"]), slog.String("new_email", req.GetNewEmail()))
		if olduser.Email != req.GetNewEmail() {
			go func() {
				ctx := cotel.CloneTrace(ctx)
				if e := s.baseDao.RedisDelUser(ctx, md["Token-User"]); e != nil {
					slog.ErrorContext(ctx, "[UpdateEmail] clean redis failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
				}
				s.stop.DoneOne()
			}()
			go func() {
				if olduser.Email != "" {
					ctx := cotel.CloneTrace(ctx)
					if e := s.baseDao.RedisDelUserEmailIndex(ctx, olduser.Email); e != nil {
						slog.ErrorContext(ctx, "[UpdateEmail] clean redis failed", slog.String("email", olduser.Email), slog.String("error", e.Error()))
					}
				}
				s.stop.DoneOne()
			}()
			go func() {
				if req.GetNewEmail() != "" {
					ctx := cotel.CloneTrace(ctx)
					if e := s.baseDao.RedisDelUserEmailIndex(ctx, req.GetNewEmail()); e != nil {
						slog.ErrorContext(ctx, "[UpdateEmail] clean redis failed",
							slog.String("email", req.GetNewEmail()), slog.String("error", e.Error()))
					}
				}
				s.stop.DoneOne()
			}()
		} else {
			s.stop.DoneOne()
			s.stop.DoneOne()
			s.stop.DoneOne()
		}
		resp := &api.UpdateEmailResp{}
		resp.SetStep("success")
		return resp, nil
	}
	e = s.baseDao.RedisCodeCheckTimes(ctx, md["Token-User"], util.UpdateEmailStep2, req.GetNewEmail())
	if e != nil && e != ecode.ErrCodeNotExist {
		slog.ErrorContext(ctx, "[UpdateEmail] redis op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	} else if e == nil {
		//if new email's code already send,we jump to step final
		resp := &api.UpdateEmailResp{}
		resp.SetStep("newverify")
		resp.SetReceiver(util.MaskEmail(req.GetNewEmail()))
		return resp, nil
	}
	if req.GetVerifySrcType() == "oauth" {
		//step 1 when update by oauth
		if req.GetVerifySrcTypeExtra() == "" || req.GetVerifyDynamicPassword() == "" {
			return nil, ecode.ErrReq
		}
		if e := s.baseDao.RedisLockEmailOP(ctx, md["Token-User"]); e != nil {
			slog.ErrorContext(ctx, "[UpdateEmail] rate check failed",
				slog.String("operator", md["Token-User"]),
				slog.String(req.GetVerifySrcTypeExtra(), req.GetVerifyDynamicPassword()),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		var oauthid string
		switch req.GetVerifySrcTypeExtra() {
		case "wechat":
			c := config.AC.Service
			oauthid, e = util.OAuthWeChatVerifyCode(ctx, c.WeChatAppID, c.WeChatAppSecret, req.GetVerifyDynamicPassword())
		default:
			return nil, ecode.ErrOAuthUnknown
		}
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		user, e := s.baseDao.GetUserByOAuth(ctx, req.GetVerifySrcTypeExtra(), oauthid)
		if e != nil {
			slog.ErrorContext(ctx, "[UpdateEmail] dao op failed",
				slog.String("operator", md["Token-User"]),
				slog.String(req.GetVerifySrcTypeExtra(), oauthid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.UserID.Hex() != md["Token-User"] {
			slog.ErrorContext(ctx, "[UpdateEmail] this is not the required oauth",
				slog.String("operator", md["Token-User"]),
				slog.String(req.GetVerifySrcTypeExtra(), oauthid))
			return nil, ecode.ErrOAuthWrong
		}
		if user.Email == req.GetNewEmail() {
			resp := &api.UpdateEmailResp{}
			resp.SetStep("success")
			return resp, nil
		}
		if user.BTime != 0 {
			return nil, ecode.ErrBan
		}
		//verify success
		if e := s.sendcode(ctx, "UpdateEmail", "email", req.GetNewEmail(), md["Token-User"], util.UpdateEmailStep2); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		resp := &api.UpdateEmailResp{}
		resp.SetStep("newverify")
		resp.SetReceiver(util.MaskEmail(req.GetNewEmail()))
		return resp, nil
	}
	if req.GetVerifyDynamicPassword() != "" {
		//step 2 when update by dynamic password
		if e := s.baseDao.RedisCheckCode(ctx, md["Token-User"], util.UpdateEmailStep1, req.GetVerifyDynamicPassword(), ""); e != nil {
			slog.ErrorContext(ctx, "[UpdateEmail] redis op failed",
				slog.String("operator", md["Token-User"]),
				slog.String("code", req.GetVerifyDynamicPassword()),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		//verify success
		if e := s.sendcode(ctx, "UpdateEmail", "email", req.GetNewEmail(), md["Token-User"], util.UpdateEmailStep2); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		resp := &api.UpdateEmailResp{}
		resp.SetStep("newverify")
		resp.SetReceiver(util.MaskEmail(req.GetNewEmail()))
		return resp, nil
	}
	//step 1 when update by dynamic password
	user, e := s.baseDao.GetUser(ctx, operator)
	if e != nil {
		slog.ErrorContext(ctx, "[UpdateEmail] dao op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if user.Email == req.GetNewEmail() {
		resp := &api.UpdateEmailResp{}
		resp.SetStep("success")
		return resp, nil
	}
	if user.BTime != 0 {
		return nil, ecode.ErrBan
	}

	if req.GetVerifySrcType() == "tel" && user.Tel == "" {
		slog.ErrorContext(ctx, "[UpdateEmail] missing tel,can't use tel to receive dynamic password", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.GetVerifySrcType() == "email" && user.Email == "" {
		slog.ErrorContext(ctx, "[UpdateEmail] missing email,can't use email to receive dynamic password", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}

	if req.GetVerifySrcType() == "email" {
		e = s.sendcode(ctx, "UpdateEmail", "email", user.Email, md["Token-User"], util.UpdateEmailStep1)
	} else {
		e = s.sendcode(ctx, "UpdateEmail", "tel", user.Tel, md["Token-User"], util.UpdateEmailStep1)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.UpdateEmailResp{}
	resp.SetStep("oldverify")
	if req.GetVerifySrcType() == "email" {
		resp.SetReceiver(util.MaskEmail(user.Email))
	} else {
		resp.SetReceiver(util.MaskTel(user.Tel))
	}
	return resp, nil
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
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[DelEmail] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
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
		if olduser, e = s.baseDao.MongoUpdateUserEmail(ctx, operator, ""); e != nil {
			s.stop.DoneOne()
			s.stop.DoneOne()
			slog.ErrorContext(ctx, "[DelEmail] db op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
			return false, e
		}
		if olduser.Email != "" {
			go func() {
				ctx := cotel.CloneTrace(ctx)
				if e := s.baseDao.RedisDelUser(ctx, md["Token-User"]); e != nil {
					slog.ErrorContext(ctx, "[DelEmail] clean redis failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
				}
				s.stop.DoneOne()
			}()
			go func() {
				ctx := cotel.CloneTrace(ctx)
				if e := s.baseDao.RedisDelUserEmailIndex(ctx, olduser.Email); e != nil {
					slog.ErrorContext(ctx, "[DelEmail] clean redis failed", slog.String("email", olduser.Email), slog.String("error", e.Error()))
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
		slog.InfoContext(ctx, "[DelEmail] success", slog.String("operator", md["Token-User"]), slog.Bool("final", final))
		return final, nil
	}
	if req.GetVerifySrcType() == "oauth" {
		if req.GetVerifySrcTypeExtra() == "" || req.GetVerifyDynamicPassword() == "" {
			return nil, ecode.ErrReq
		}
		if e := s.baseDao.RedisLockEmailOP(ctx, md["Token-User"]); e != nil {
			slog.ErrorContext(ctx, "[DelEmail] rate check failed",
				slog.String("operator", md["Token-User"]),
				slog.String(req.GetVerifySrcTypeExtra(), req.GetVerifyDynamicPassword()),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		var oauthid string
		switch req.GetVerifySrcTypeExtra() {
		case "wechat":
			c := config.AC.Service
			oauthid, e = util.OAuthWeChatVerifyCode(ctx, c.WeChatAppID, c.WeChatAppSecret, req.GetVerifyDynamicPassword())
		default:
			return nil, ecode.ErrOAuthUnknown
		}
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		user, e := s.baseDao.GetUserByOAuth(ctx, req.GetVerifySrcTypeExtra(), oauthid)
		if e != nil {
			slog.ErrorContext(ctx, "[DelEmail] dao op failed",
				slog.String("operator", md["Token-User"]),
				slog.String(req.GetVerifySrcTypeExtra(), oauthid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.UserID.Hex() != md["Token-User"] {
			slog.ErrorContext(ctx, "[DelEmail] this is not the required oauth",
				slog.String("operator", md["Token-User"]),
				slog.String(req.GetVerifySrcTypeExtra(), oauthid))
			return nil, ecode.ErrOAuthWrong
		}
		if user.Email == "" {
			resp := &api.DelEmailResp{}
			resp.SetStep("success")
			resp.SetFinal(user.Email == "" && user.Tel == "" && user.IDCard == "" && len(user.OAuths) == 0)
			return resp, nil
		}
		if user.BTime != 0 {
			return nil, ecode.ErrBan
		}
		//verify success
		final, e := update()
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		resp := &api.DelEmailResp{}
		resp.SetStep("success")
		resp.SetFinal(final)
		return resp, nil
	}
	if req.GetVerifyDynamicPassword() != "" {
		//step2
		if e := s.baseDao.RedisCheckCode(ctx, md["Token-User"], util.DelEmail, req.GetVerifyDynamicPassword(), ""); e != nil {
			slog.ErrorContext(ctx, "[DelEmail] redis op failed",
				slog.String("operator", md["Token-User"]),
				slog.String("code", req.GetVerifyDynamicPassword()),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		//verify success
		final, e := update()
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		resp := &api.DelEmailResp{}
		resp.SetStep("success")
		resp.SetFinal(final)
		return resp, nil
	}
	//step1
	user, e := s.baseDao.GetUser(ctx, operator)
	if e != nil {
		slog.ErrorContext(ctx, "[DelEmail] dao op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	var final bool
	if user.IDCard == "" && user.Tel == "" && len(user.OAuths) == 0 {
		final = true
	}
	if user.Email == "" {
		resp := &api.DelEmailResp{}
		resp.SetStep("success")
		resp.SetFinal(final)
		return resp, nil
	}
	if user.BTime != 0 {
		return nil, ecode.ErrBan
	}

	if req.GetVerifySrcType() == "tel" && user.Tel == "" {
		slog.ErrorContext(ctx, "[DelEmail] missing tel,can't use tel to receive dynamic password", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.GetVerifySrcType() == "email" && user.Email == "" {
		slog.ErrorContext(ctx, "[DelEmail] missing email,can't use email to receive dynamic password", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}

	if req.GetVerifySrcType() == "email" {
		e = s.sendcode(ctx, "DelEmail", "email", user.Email, md["Token-User"], util.DelEmail)
	} else {
		e = s.sendcode(ctx, "DelEmail", "tel", user.Tel, md["Token-User"], util.DelEmail)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.DelEmailResp{}
	resp.SetStep("oldverify")
	resp.SetFinal(final)
	if req.GetVerifySrcType() == "email" {
		resp.SetReceiver(util.MaskEmail(user.Email))
	} else {
		resp.SetReceiver(util.MaskTel(user.Tel))
	}
	return resp, nil
}

func (s *Service) TelDuplicateCheck(ctx context.Context, req *api.TelDuplicateCheckReq) (*api.TelDuplicateCheckResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[TelDuplicateCheck] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
	}
	if user, e := s.baseDao.GetUser(ctx, operator); e != nil {
		slog.ErrorContext(ctx, "[TelDuplicateCheck] dao op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	} else if user.BTime != 0 {
		return nil, ecode.ErrBan
	}
	//redis lock
	if e := s.baseDao.RedisLockDuplicateCheck(ctx, "tel", md["Token-User"]); e != nil {
		slog.ErrorContext(ctx, "[TelDuplicateCheck] redis op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	userid, e := s.baseDao.GetUserTelIndex(ctx, req.GetTel())
	if e != nil && e != ecode.ErrUserNotExist {
		slog.ErrorContext(ctx, "[TelDuplicateCheck] dao op failed",
			slog.String("operator", md["Token-User"]), slog.String("tel", req.GetTel()), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.TelDuplicateCheckResp{}
	resp.SetDuplicate(userid != "")
	return resp, nil
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
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[UpdateTel] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
	}
	if req.GetNewTelDynamicPassword() != "" {
		//step final
		e = s.baseDao.RedisCheckCode(ctx, md["Token-User"], util.UpdateTelStep2, req.GetNewTelDynamicPassword(), req.GetNewTel())
		if e != nil {
			slog.ErrorContext(ctx, "[UpdateTel] redis op failed",
				slog.String("operator", md["Token-User"]),
				slog.String("code", req.GetNewTelDynamicPassword()),
				slog.String("error", e.Error()))
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
		if olduser, e = s.baseDao.MongoUpdateUserTel(ctx, operator, req.GetNewTel()); e != nil {
			s.stop.DoneOne()
			s.stop.DoneOne()
			s.stop.DoneOne()
			slog.ErrorContext(ctx, "[UpdateTel] db op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		slog.InfoContext(ctx, "[UpdateTel] success",
			slog.String("operator", md["Token-User"]), slog.String("new_tel", req.GetNewTel()))
		if olduser.Tel != req.GetNewTel() {
			go func() {
				ctx := cotel.CloneTrace(ctx)
				if e := s.baseDao.RedisDelUser(ctx, md["Token-User"]); e != nil {
					slog.ErrorContext(ctx, "[UpdateTel] clean redis failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
				}
				s.stop.DoneOne()
			}()
			go func() {
				if olduser.Tel != "" {
					ctx := cotel.CloneTrace(ctx)
					if e := s.baseDao.RedisDelUserTelIndex(ctx, olduser.Tel); e != nil {
						slog.ErrorContext(ctx, "[UpdateTel] clean redis failed", slog.String("tel", olduser.Tel), slog.String("error", e.Error()))
					}
				}
				s.stop.DoneOne()
			}()
			go func() {
				if req.GetNewTel() != "" {
					ctx := cotel.CloneTrace(ctx)
					if e := s.baseDao.RedisDelUserTelIndex(ctx, req.GetNewTel()); e != nil {
						slog.ErrorContext(ctx, "[UpdateTel] clean redis failed",
							slog.String("tel", req.GetNewTel()), slog.String("error", e.Error()))
					}
				}
				s.stop.DoneOne()
			}()
		} else {
			s.stop.DoneOne()
			s.stop.DoneOne()
			s.stop.DoneOne()
		}
		resp := &api.UpdateTelResp{}
		resp.SetStep("success")
		return resp, nil
	}
	e = s.baseDao.RedisCodeCheckTimes(ctx, md["Token-User"], util.UpdateTelStep2, req.GetNewTel())
	if e != nil && e != ecode.ErrCodeNotExist {
		slog.ErrorContext(ctx, "[UpdateTel] redis op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	} else if e == nil {
		//if new tel's code already send,we jump to step final
		resp := &api.UpdateTelResp{}
		resp.SetStep("newverify")
		resp.SetReceiver(util.MaskTel(req.GetNewTel()))
		return resp, nil
	}
	if req.GetVerifySrcType() == "oauth" {
		//step 1 when update by oauth
		if req.GetVerifySrcTypeExtra() == "" || req.GetVerifyDynamicPassword() == "" {
			return nil, ecode.ErrReq
		}
		if e := s.baseDao.RedisLockTelOP(ctx, md["Token-User"]); e != nil {
			slog.ErrorContext(ctx, "[UpdateTel] rate check failed",
				slog.String("operator", md["Token-User"]),
				slog.String(req.GetVerifySrcTypeExtra(), req.GetVerifyDynamicPassword()),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		var oauthid string
		switch req.GetVerifySrcTypeExtra() {
		case "wechat":
			c := config.AC.Service
			oauthid, e = util.OAuthWeChatVerifyCode(ctx, c.WeChatAppID, c.WeChatAppSecret, req.GetVerifyDynamicPassword())
		default:
			return nil, ecode.ErrOAuthUnknown
		}
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		user, e := s.baseDao.GetUserByOAuth(ctx, req.GetVerifySrcTypeExtra(), oauthid)
		if e != nil {
			slog.ErrorContext(ctx, "[UpdateTel] dao op failed",
				slog.String("operator", md["Token-User"]),
				slog.String(req.GetVerifySrcTypeExtra(), oauthid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.UserID.Hex() != md["Token-User"] {
			slog.ErrorContext(ctx, "[UpdateTel] this is not the required oauth",
				slog.String("operator", md["Token-User"]),
				slog.String(req.GetVerifySrcTypeExtra(), oauthid))
			return nil, ecode.ErrOAuthWrong
		}
		if user.Tel == req.GetNewTel() {
			resp := &api.UpdateTelResp{}
			resp.SetStep("success")
			return resp, nil
		}
		if user.BTime != 0 {
			return nil, ecode.ErrBan
		}
		//verify success
		if e := s.sendcode(ctx, "UpdateTel", "tel", req.GetNewTel(), md["Token-User"], util.UpdateTelStep2); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		resp := &api.UpdateTelResp{}
		resp.SetStep("newverify")
		resp.SetReceiver(util.MaskTel(req.GetNewTel()))
		return resp, nil
	}
	if req.GetVerifyDynamicPassword() != "" {
		//step 2 when update by dynamic password
		if e := s.baseDao.RedisCheckCode(ctx, md["Token-User"], util.UpdateTelStep1, req.GetVerifyDynamicPassword(), ""); e != nil {
			slog.ErrorContext(ctx, "[UpdateTel] redis op failed",
				slog.String("operator", md["Token-User"]),
				slog.String("code", req.GetVerifyDynamicPassword()),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		//verify success
		if e := s.sendcode(ctx, "UpdateTel", "tel", req.GetNewTel(), md["Token-User"], util.UpdateTelStep2); e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		resp := &api.UpdateTelResp{}
		resp.SetStep("newverify")
		resp.SetReceiver(util.MaskTel(req.GetNewTel()))
		return resp, nil
	}
	//step 1 when update by dynamic password
	user, e := s.baseDao.GetUser(ctx, operator)
	if e != nil {
		slog.ErrorContext(ctx, "[UpdateTel] dao op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if user.Tel == req.GetNewTel() {
		resp := &api.UpdateTelResp{}
		resp.SetStep("success")
		return resp, nil
	}
	if user.BTime != 0 {
		return nil, ecode.ErrBan
	}

	if req.GetVerifySrcType() == "tel" && user.Tel == "" {
		slog.ErrorContext(ctx, "[UpdateTel] missing tel,can't use tel to receive dynamic password", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.GetVerifySrcType() == "email" && user.Email == "" {
		slog.ErrorContext(ctx, "[UpdateTel] missing email,can't use email to receive dynamic password", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}

	if req.GetVerifySrcType() == "email" {
		e = s.sendcode(ctx, "UpdateTel", "email", user.Email, md["Token-User"], util.UpdateTelStep1)
	} else {
		e = s.sendcode(ctx, "UpdateTel", "tel", user.Tel, md["Token-User"], util.UpdateTelStep1)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.UpdateTelResp{}
	resp.SetStep("oldverify")
	if req.GetVerifySrcType() == "email" {
		resp.SetReceiver(util.MaskEmail(user.Email))
	} else {
		resp.SetReceiver(util.MaskTel(user.Tel))
	}
	return resp, nil
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
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[DelTel] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
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
		if olduser, e = s.baseDao.MongoUpdateUserTel(ctx, operator, ""); e != nil {
			s.stop.DoneOne()
			s.stop.DoneOne()
			slog.ErrorContext(ctx, "[DelTel] db op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
			return false, e
		}
		if olduser.Tel != "" {
			go func() {
				ctx := cotel.CloneTrace(ctx)
				if e := s.baseDao.RedisDelUser(ctx, md["Token-User"]); e != nil {
					slog.ErrorContext(ctx, "[DelTel] clean redis failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
				}
				s.stop.DoneOne()
			}()
			go func() {
				ctx := cotel.CloneTrace(ctx)
				if e := s.baseDao.RedisDelUserTelIndex(ctx, olduser.Tel); e != nil {
					slog.ErrorContext(ctx, "[DelTel] clean redis failed", slog.String("tel", olduser.Tel), slog.String("error", e.Error()))
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
		slog.InfoContext(ctx, "[DelTel] success", slog.String("operator", md["Token-User"]), slog.Bool("final", final))
		return final, nil
	}
	if req.GetVerifySrcType() == "oauth" {
		if req.GetVerifySrcTypeExtra() == "" || req.GetVerifyDynamicPassword() == "" {
			return nil, ecode.ErrReq
		}
		if e := s.baseDao.RedisLockTelOP(ctx, md["Token-User"]); e != nil {
			slog.ErrorContext(ctx, "[DelTel] rate check failed",
				slog.String("operator", md["Token-User"]),
				slog.String(req.GetVerifySrcTypeExtra(), req.GetVerifyDynamicPassword()),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		var oauthid string
		switch req.GetVerifySrcTypeExtra() {
		case "wechat":
			c := config.AC.Service
			oauthid, e = util.OAuthWeChatVerifyCode(ctx, c.WeChatAppID, c.WeChatAppSecret, req.GetVerifyDynamicPassword())
		default:
			return nil, ecode.ErrOAuthUnknown
		}
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		user, e := s.baseDao.GetUserByOAuth(ctx, req.GetVerifySrcTypeExtra(), oauthid)
		if e != nil {
			slog.ErrorContext(ctx, "[DelTel] dao op failed",
				slog.String("operator", md["Token-User"]),
				slog.String(req.GetVerifySrcTypeExtra(), oauthid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if user.UserID.Hex() != md["Token-User"] {
			slog.ErrorContext(ctx, "[DelTel] this is not the required oauth",
				slog.String("operator", md["Token-User"]),
				slog.String(req.GetVerifySrcTypeExtra(), oauthid))
			return nil, ecode.ErrOAuthWrong
		}
		if user.Tel == "" {
			resp := &api.DelTelResp{}
			resp.SetStep("success")
			resp.SetFinal(user.Email == "" && user.Tel == "" && user.IDCard == "" && len(user.OAuths) == 0)
			return resp, nil
		}
		if user.BTime != 0 {
			return nil, ecode.ErrBan
		}
		//verify success
		final, e := update()
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		resp := &api.DelTelResp{}
		resp.SetStep("success")
		resp.SetFinal(final)
		return resp, nil
	}
	if req.GetVerifyDynamicPassword() != "" {
		//step2
		if e := s.baseDao.RedisCheckCode(ctx, md["Token-User"], util.DelTel, req.GetVerifyDynamicPassword(), ""); e != nil {
			slog.ErrorContext(ctx, "[DelTel] redis op failed",
				slog.String("operator", md["Token-User"]),
				slog.String("code", req.GetVerifyDynamicPassword()),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		//verify success
		final, e := update()
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		resp := &api.DelTelResp{}
		resp.SetStep("success")
		resp.SetFinal(final)
		return resp, nil
	}
	//step1
	user, e := s.baseDao.GetUser(ctx, operator)
	if e != nil {
		slog.ErrorContext(ctx, "[DelTel] dao op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	var final bool
	if user.IDCard == "" && user.Email == "" && len(user.OAuths) == 0 {
		final = true
	}
	if user.Tel == "" {
		resp := &api.DelTelResp{}
		resp.SetStep("success")
		resp.SetFinal(final)
		return resp, nil
	}
	if user.BTime != 0 {
		return nil, ecode.ErrBan
	}

	if req.GetVerifySrcType() == "tel" && user.Tel == "" {
		slog.ErrorContext(ctx, "[DelTel] missing tel,can't use tel to receive dynamic password", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}
	if req.GetVerifySrcType() == "email" && user.Email == "" {
		slog.ErrorContext(ctx, "[DelTel] missing email,can't use email to receive dynamic password", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrReq
	}

	if req.GetVerifySrcType() == "email" {
		e = s.sendcode(ctx, "DelTel", "email", user.Email, md["Token-User"], util.DelTel)
	} else {
		e = s.sendcode(ctx, "DelTel", "tel", user.Tel, md["Token-User"], util.DelTel)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.DelTelResp{}
	resp.SetStep("oldverify")
	resp.SetFinal(final)
	if req.GetVerifySrcType() == "email" {
		resp.SetReceiver(util.MaskEmail(user.Email))
	} else {
		resp.SetReceiver(util.MaskTel(user.Tel))
	}
	return resp, nil
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
