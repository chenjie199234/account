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
func (s *Service) GetMoneyLogs(ctx context.Context, req *api.GetMoneyLogsReq) (*api.GetMoneyLogsResp, error) {
	if req.EndTime < req.StartTime {
		return nil, ecode.ErrReq
	}
	var userid primitive.ObjectID
	switch req.SrcType {
	case "user_id":
		var e error
		if userid, e = primitive.ObjectIDFromHex(req.Src); e != nil {
			log.Error(ctx, "[GetMoneyLogs] userid format wrong", log.String("user_id", req.Src), log.CError(e))
			return nil, ecode.ErrReq
		}
	case "tel":
		useridstr, e := s.userDao.GetUserTelIndex(ctx, req.Src)
		if e != nil {
			log.Error(ctx, "[GetMoneyLogs] dao op failed", log.String("tel", req.Src), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if userid, e = primitive.ObjectIDFromHex(useridstr); e != nil {
			log.Error(ctx, "[GetMoneyLogs] userid format wrong", log.String("tel", req.Src), log.CError(e))
			return nil, ecode.ErrSystem
		}
	case "email":
		useridstr, e := s.userDao.GetUserEmailIndex(ctx, req.Src)
		if e != nil {
			log.Error(ctx, "[GetMoneyLogs] dao op failed", log.String("email", req.Src), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if userid, e = primitive.ObjectIDFromHex(useridstr); e != nil {
			log.Error(ctx, "[GetMoneyLogs] userid format wrong", log.String("email", req.Src), log.CError(e))
			return nil, ecode.ErrSystem
		}
	case "idcard":
		useridstr, e := s.userDao.GetUserIDCardIndex(ctx, req.Src)
		if e != nil {
			log.Error(ctx, "[GetMoneyLogs] dao op failed", log.String("idcard", req.Src), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if userid, e = primitive.ObjectIDFromHex(useridstr); e != nil {
			log.Error(ctx, "[GetMoneyLogs] userid format wrong", log.String("idcard", req.Src), log.CError(e))
			return nil, ecode.ErrSystem
		}
	}
	logs, totalsize, page, e := s.moneyDao.GetMoneyLogs(ctx, userid, req.Action, req.StartTime, req.EndTime, moneydao.DefaultMoneyLogsPageSize, req.Page)
	if e != nil {
		log.Error(ctx, "[GetMoneyLogs] dao op failed", log.String(req.SrcType, req.Src), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.GetMoneyLogsResp{
		Page:      page,
		Pagesize:  moneydao.DefaultMoneyLogsPageSize,
		Totalsize: totalsize,
		Logs:      make([]*api.MoneyLog, 0, len(logs)),
	}
	if resp.Page == 0 {
		resp.Pagesize = resp.Totalsize
	}
	for _, v := range logs {
		resp.Logs = append(resp.Logs, &api.MoneyLog{
			UserId:      v.UserID.Hex(),
			Action:      v.Action,
			UniqueId:    v.UniqueID,
			SrcDst:      v.SrcDst,
			MoneyType:   v.MoneyType,
			MoneyAmount: v.MoneyAmount,
			Ctime:       uint32(v.LogID.Timestamp().Unix()),
		})
	}
	return resp, nil
}
func (s *Service) SelfMoneyLogs(ctx context.Context, req *api.SelfMoneyLogsReq) (*api.SelfMoneyLogsResp, error) {
	if req.EndTime < req.StartTime {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[SelfMoneyLogs] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}
	logs, totalsize, page, e := s.moneyDao.GetMoneyLogs(ctx, operator, req.Action, req.StartTime, req.EndTime, moneydao.DefaultMoneyLogsPageSize, req.Page)
	if e != nil {
		log.Error(ctx, "[SelfMoneyLogs] dao op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.SelfMoneyLogsResp{
		Page:      page,
		Pagesize:  moneydao.DefaultMoneyLogsPageSize,
		Totalsize: totalsize,
		Logs:      make([]*api.MoneyLog, 0, len(logs)),
	}
	if resp.Page == 0 {
		resp.Pagesize = resp.Totalsize
	}
	for _, v := range logs {
		resp.Logs = append(resp.Logs, &api.MoneyLog{
			UserId:      v.UserID.Hex(),
			Action:      v.Action,
			UniqueId:    v.UniqueID,
			SrcDst:      v.SrcDst,
			MoneyType:   v.MoneyType,
			MoneyAmount: v.MoneyAmount,
			Ctime:       uint32(v.LogID.Timestamp().Unix()),
		})
	}
	return resp, nil
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
