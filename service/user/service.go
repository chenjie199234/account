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
	var user *model.User
	var e error
	switch req.SrcType {
	case "idcard":
		if req.PasswordType == "dynamic" {
			log.Error(ctx, "[Login] idcard can't use dynamic password", nil)
			return nil, ecode.ErrReq
		}
		if user, e = s.userDao.MongoGetUserByIDCard(ctx, req.Src); e != nil {
			log.Error(ctx, "[Login] db op failed", map[string]interface{}{"idcard": req.Src, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	case "nickname":
		if req.PasswordType == "dynamic" {
			log.Error(ctx, "[Login] nickname can't use dynamic password", nil)
			return nil, ecode.ErrReq
		}
		if user, e = s.userDao.MongoGetUserByNickName(ctx, req.Src); e != nil {
			log.Error(ctx, "[Login] db op failed", map[string]interface{}{"nickname": req.Src, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	case "tel":
		if req.PasswordType == "static" {
			if user, e = s.userDao.MongoGetUserByTel(ctx, req.Src); e != nil {
				log.Error(ctx, "[Login] db op failed", map[string]interface{}{"tel": req.Src, "error": e})
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
		} else if req.Password == "" {
			//TODO send code
		} else {
			//TODO verify code
		}
	case "email":
		if req.PasswordType == "static" {
			if user, e = s.userDao.MongoGetUserByEmail(ctx, req.Src); e != nil {
				log.Error(ctx, "[Login] db op failed", map[string]interface{}{"email": req.Src, "error": e})
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
		} else if req.Password == "" {
			//TODO send code
		} else {
			//TODO verify code
		}
	}
	if e := util.SignCheck(req.Password, user.Password); e != nil {
		if e == ecode.ErrSignCheckFailed {
			e = ecode.ErrPasswordWrong
		}
		log.Error(ctx, "[Login] sign check failed", map[string]interface{}{"src_type": req.SrcType, "src": req.Src, "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	//TODO set the puber
	token := publicmids.MakeToken(ctx, "", *config.EC.DeployEnv, *config.EC.RunEnv, user.UserID.Hex())
	return &api.LoginResp{
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
	}, nil
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
