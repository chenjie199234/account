package money

import (
	"context"

	"github.com/chenjie199234/account/api"
	"github.com/chenjie199234/account/config"
	moneydao "github.com/chenjie199234/account/dao/money"
	"github.com/chenjie199234/account/ecode"

	// "github.com/chenjie199234/Corelib/cgrpc"
	// "github.com/chenjie199234/Corelib/crpc"
	// "github.com/chenjie199234/Corelib/log"
	// "github.com/chenjie199234/Corelib/web"
	"github.com/chenjie199234/Corelib/util/graceful"
)

// Service subservice for money business
type Service struct {
	stop *graceful.Graceful

	moneyDao *moneydao.Dao
}

// Start -
func Start() *Service {
	return &Service{
		stop: graceful.New(),

		moneyDao: moneydao.NewDao(config.GetSql("money_sql"), config.GetRedis("money_redis"), config.GetMongo("money_mongo")),
	}
}
func (s *Service) GetMoneyLogs(ctx context.Context, req *api.GetMoneyLogsReq) (*api.GetMoneyLogsResp, error) {

}
func (s *Service) RechargeMoney(ctx context.Context, req *api.RechargeMoneyReq) (*api.RechargeMoneyResp, error) {

}
func (s *Service) SpendMoney(ctx context.Context, req *api.SpendMoneyReq) (*api.SpendMoneyResp, error) {

}
func (s *Service) RefundMoney(ctx context.Context, req *api.RefundMoneyReq) (*api.RefundMoneyResp, error) {

}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
