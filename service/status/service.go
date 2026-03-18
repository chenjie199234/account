package status

import (
	"context"
	"time"

	// "github.com/chenjie199234/account/config"
	"github.com/chenjie199234/account/api"
	statusdao "github.com/chenjie199234/account/dao/status"
	"github.com/chenjie199234/account/model"
	// "github.com/chenjie199234/account/ecode"

	"github.com/chenjie199234/Corelib/cotel"
	"github.com/chenjie199234/Corelib/util/graceful"
	"github.com/chenjie199234/Corelib/util/host"
	// "github.com/chenjie199234/Corelib/cgrpc"
	// "github.com/chenjie199234/Corelib/crpc"
	// "github.com/chenjie199234/Corelib/log"
	// "github.com/chenjie199234/Corelib/web"
)

// Service subservice for status business
type Service struct {
	stop *graceful.Graceful

	statusDao *statusdao.Dao
}

// Start -
func Start() (*Service, error) {
	return &Service{
		stop: graceful.New(),

		//statusDao: statusdao.NewDao(config.GetSql("status_sql"), config.GetRedis("status_redis"), config.GetMongo("status_mongo")),
		statusDao: statusdao.NewDao(nil, nil, nil),
	}, nil
}

// Ping -
func (s *Service) Ping(ctx context.Context, in *api.Pingreq) (*api.Pingresp, error) {
	//if _, ok := ctx.(crpc.NoStreamServerContext); ok {
	//        log.Info("this is a crpc call")
	//}
	//if _, ok := ctx.(cgrpc.NoStreamServerContext); ok {
	//        log.Info("this is a cgrpc call")
	//}
	//if _, ok := ctx.(web.NoStreamServerContext); ok {
	//        log.Info("this is a web call")
	//}
	cpu, cpuu, cput, mem, memu, memt := cotel.GetCpuMemUsage()
	resp := &api.Pingresp{}
	resp.SetClientTimestamp(in.GetTimestamp())
	resp.SetServerTimestamp(time.Now().UnixNano())
	resp.SetHost(host.Hostname)
	resp.SetIp(host.Hostip)
	resp.SetCpuNum(cpu)
	resp.SetCpuUsage(cpuu)
	resp.SetCpuType(cput)
	resp.SetMemTotal(mem)
	resp.SetMemUsage(memu)
	resp.SetMemType(memt)
	resp.SetVersion(model.Version)
	return resp, nil
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
