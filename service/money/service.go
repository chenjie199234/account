package money

import (
	"context"
	"log/slog"

	"github.com/chenjie199234/account/api"
	"github.com/chenjie199234/account/config"
	basedao "github.com/chenjie199234/account/dao/base"
	moneydao "github.com/chenjie199234/account/dao/money"
	"github.com/chenjie199234/account/ecode"

	// "github.com/chenjie199234/Corelib/cgrpc"
	// "github.com/chenjie199234/Corelib/crpc"
	// "github.com/chenjie199234/Corelib/web"
	"github.com/chenjie199234/Corelib/metadata"
	"github.com/chenjie199234/Corelib/util/graceful"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Service subservice for money business
type Service struct {
	stop *graceful.Graceful

	baseDao  *basedao.Dao
	moneyDao *moneydao.Dao
}

// Start -
func Start() (*Service, error) {
	return &Service{
		stop: graceful.New(),

		baseDao:  basedao.NewDao(nil, config.GetRedis("account_redis"), config.GetMongo("account_mongo")),
		moneyDao: moneydao.NewDao(nil, config.GetRedis("account_redis"), config.GetMongo("account_mongo")),
	}, nil
}
func (s *Service) GetMoneyLogs(ctx context.Context, req *api.GetMoneyLogsReq) (*api.GetMoneyLogsResp, error) {
	if req.EndTime < req.StartTime {
		return nil, ecode.ErrReq
	}
	var userid bson.ObjectID
	switch req.SrcType {
	case "user_id":
		var e error
		if userid, e = bson.ObjectIDFromHex(req.Src); e != nil {
			slog.ErrorContext(ctx, "[GetMoneyLogs] userid format wrong", slog.String("user_id", req.Src), slog.String("error", e.Error()))
			return nil, ecode.ErrReq
		}
	case "tel":
		useridstr, e := s.baseDao.GetUserTelIndex(ctx, req.Src)
		if e != nil {
			slog.ErrorContext(ctx, "[GetMoneyLogs] dao op failed", slog.String("tel", req.Src), slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if userid, e = bson.ObjectIDFromHex(useridstr); e != nil {
			slog.ErrorContext(ctx, "[GetMoneyLogs] userid format wrong", slog.String("tel", req.Src), slog.String("error", e.Error()))
			return nil, ecode.ErrSystem
		}
	case "email":
		useridstr, e := s.baseDao.GetUserEmailIndex(ctx, req.Src)
		if e != nil {
			slog.ErrorContext(ctx, "[GetMoneyLogs] dao op failed", slog.String("email", req.Src), slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if userid, e = bson.ObjectIDFromHex(useridstr); e != nil {
			slog.ErrorContext(ctx, "[GetMoneyLogs] userid format wrong", slog.String("email", req.Src), slog.String("error", e.Error()))
			return nil, ecode.ErrSystem
		}
	case "idcard":
		useridstr, e := s.baseDao.GetUserIDCardIndex(ctx, req.Src)
		if e != nil {
			slog.ErrorContext(ctx, "[GetMoneyLogs] dao op failed", slog.String("idcard", req.Src), slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if userid, e = bson.ObjectIDFromHex(useridstr); e != nil {
			slog.ErrorContext(ctx, "[GetMoneyLogs] userid format wrong", slog.String("idcard", req.Src), slog.String("error", e.Error()))
			return nil, ecode.ErrSystem
		}
	}
	logs, totalsize, page, e := s.moneyDao.GetMoneyLogs(ctx, userid, req.Action, req.StartTime, req.EndTime, moneydao.DefaultMoneyLogsPageSize, req.Page)
	if e != nil {
		slog.ErrorContext(ctx, "[GetMoneyLogs] dao op failed", slog.String(req.SrcType, req.Src), slog.String("error", e.Error()))
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
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[SelfMoneyLogs] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
	}
	logs, totalsize, page, e := s.moneyDao.GetMoneyLogs(ctx, operator, req.Action, req.StartTime, req.EndTime, moneydao.DefaultMoneyLogsPageSize, req.Page)
	if e != nil {
		slog.ErrorContext(ctx, "[SelfMoneyLogs] dao op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
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
