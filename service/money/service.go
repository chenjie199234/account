package money

import (
	"context"

	"github.com/chenjie199234/account/api"
	"github.com/chenjie199234/account/config"
	moneydao "github.com/chenjie199234/account/dao/money"
	userdao "github.com/chenjie199234/account/dao/user"
	"github.com/chenjie199234/account/ecode"

	// "github.com/chenjie199234/Corelib/cgrpc"
	// "github.com/chenjie199234/Corelib/crpc"
	// "github.com/chenjie199234/Corelib/web"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/metadata"
	"github.com/chenjie199234/Corelib/util/graceful"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service subservice for money business
type Service struct {
	stop *graceful.Graceful

	userDao  *userdao.Dao
	moneyDao *moneydao.Dao
}

// Start -
func Start() *Service {
	return &Service{
		stop: graceful.New(),

		userDao:  userdao.NewDao(nil, config.GetRedis("account_redis"), config.GetMongo("account_mongo")),
		moneyDao: moneydao.NewDao(nil, config.GetRedis("account_redis"), config.GetMongo("account_mongo")),
	}
}
func (s *Service) GetUserMoneyLogs(ctx context.Context, req *api.GetUserMoneyLogsReq) (*api.GetUserMoneyLogsResp, error) {
	var userid primitive.ObjectID
	switch req.SrcType {
	case "user_id":
		var e error
		if userid, e = primitive.ObjectIDFromHex(req.Src); e != nil {
			log.Error(ctx, "[GetUserMoneyLogs] userid format wrong", map[string]interface{}{"user_id": req.Src, "error": e})
			return nil, ecode.ErrReq
		}
	case "tel":
		user, e := s.userDao.MongoGetUserByTel(ctx, req.Src)
		if e != nil {
			log.Error(ctx, "[GetUserMoneyLogs] db op failed", map[string]interface{}{"tel": req.Src, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		userid = user.UserID
	case "email":
		user, e := s.userDao.MongoGetUserByEmail(ctx, req.Src)
		if e != nil {
			log.Error(ctx, "[GetUserMoneyLogs] db op failed", map[string]interface{}{"email": req.Src, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		userid = user.UserID
	case "idcard":
		user, e := s.userDao.MongoGetUserByIDCard(ctx, req.Src)
		if e != nil {
			log.Error(ctx, "[GetUserMoneyLogs] db op failed", map[string]interface{}{"idcard": req.Src, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		userid = user.UserID
	case "nickname":
		user, e := s.userDao.MongoGetUserByNickName(ctx, req.Src)
		if e != nil {
			log.Error(ctx, "[GetUserMoneyLogs] db op failed", map[string]interface{}{"nick_name": req.Src, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		userid = user.UserID
	}
	logs, page, totalsize, e := s.moneyDao.MongoGetMoneyLogs(ctx, userid, req.Action, 20, req.Page)
	if e != nil {
		log.Error(ctx, "[GetUserMoneyLogs] db op failed", map[string]interface{}{req.SrcType: req.Src, "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	apilogs := make([]*api.MoneyLog, 0, len(logs))
	for _, v := range logs {
		apilogs = append(apilogs, &api.MoneyLog{
			UserId:      v.UserID.Hex(),
			Action:      v.Action,
			UniqueId:    v.UniqueID,
			SrcDst:      v.SrcDst,
			MoneyType:   v.MoneyType,
			MoneyAmount: v.MoneyAmount,
			Ctime:       uint32(v.LogID.Timestamp().Unix()),
		})
	}
	return &api.GetUserMoneyLogsResp{
		Page:      page,
		Pagesize:  20,
		Totalsize: totalsize,
		Logs:      apilogs,
	}, nil
}
func (s *Service) SelfMoneyLogs(ctx context.Context, req *api.SelfMoneyLogsReq) (*api.SelfMoneyLogsResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[SelfMoneyLogs] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}
	logs, page, totalsize, e := s.moneyDao.MongoGetMoneyLogs(ctx, operator, req.Action, 20, req.Page)
	if e != nil {
		log.Error(ctx, "[SelfMoneyLogs] db op failed", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	apilogs := make([]*api.MoneyLog, 0, len(logs))
	for _, v := range logs {
		apilogs = append(apilogs, &api.MoneyLog{
			UserId:      v.UserID.Hex(),
			Action:      v.Action,
			UniqueId:    v.UniqueID,
			SrcDst:      v.SrcDst,
			MoneyType:   v.MoneyType,
			MoneyAmount: v.MoneyAmount,
			Ctime:       uint32(v.LogID.Timestamp().Unix()),
		})
	}
	return &api.SelfMoneyLogsResp{
		Page:      page,
		Pagesize:  20,
		Totalsize: totalsize,
		Logs:      apilogs,
	}, nil
}
func (s *Service) RechargeMoney(ctx context.Context, req *api.RechargeMoneyReq) (*api.RechargeMoneyResp, error) {
	//TODO
	return &api.RechargeMoneyResp{}, nil
}
func (s *Service) SpendMoney(ctx context.Context, req *api.SpendMoneyReq) (*api.SpendMoneyResp, error) {
	//TODO
	return &api.SpendMoneyResp{}, nil
}
func (s *Service) RefundMoney(ctx context.Context, req *api.RefundMoneyReq) (*api.RefundMoneyResp, error) {
	//TODO
	return &api.RefundMoneyResp{}, nil
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
