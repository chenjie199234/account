package xcrpc

import (
	"strings"
	"time"

	"github.com/chenjie199234/account/api"
	"github.com/chenjie199234/account/config"
	"github.com/chenjie199234/account/model"
	"github.com/chenjie199234/account/service"

	"github.com/chenjie199234/Corelib/crpc"
	"github.com/chenjie199234/Corelib/crpc/mids"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/util/ctime"
)

var s *crpc.CrpcServer

// StartCrpcServer -
func StartCrpcServer() {
	c := config.GetCrpcServerConfig()
	crpcc := &crpc.ServerConfig{
		ConnectTimeout: time.Duration(c.ConnectTimeout),
		GlobalTimeout:  time.Duration(c.GlobalTimeout),
		HeartPorbe:     time.Duration(c.HeartProbe),
		Certs:          c.Certs,
	}
	var e error
	if s, e = crpc.NewCrpcServer(crpcc, model.Project, model.Group, model.Name); e != nil {
		log.Error(nil, "[xcrpc] new server failed", map[string]interface{}{"error": e})
		return
	}
	UpdateHandlerTimeout(config.AC.HandlerTimeout)

	//this place can register global midwares
	//s.Use(globalmidwares)

	//you just need to register your service here
	api.RegisterStatusCrpcServer(s, service.SvcStatus, mids.AllMids())
	//example
	//api.RegisterExampleCrpcServer(s, service.SvcExample,mids.AllMids())
	api.RegisterUserCrpcServer(s, service.SvcUser, mids.AllMids())
	api.RegisterMoneyCrpcServer(s, service.SvcMoney, mids.AllMids())

	if e = s.StartCrpcServer(":9000"); e != nil && e != crpc.ErrServerClosed {
		log.Error(nil, "[xcrpc] start server failed", map[string]interface{}{"error": e})
		return
	}
	log.Info(nil, "[xcrpc] server closed", nil)
}

// UpdateHandlerTimeout -
// first key path,second key method,value timeout duration
func UpdateHandlerTimeout(hts map[string]map[string]ctime.Duration) {
	if s == nil {
		return
	}
	cc := make(map[string]time.Duration)
	for path, methods := range hts {
		for method, timeout := range methods {
			method = strings.ToUpper(method)
			if method == "CRPC" {
				cc[path] = timeout.StdDuration()
			}
		}
	}
	s.UpdateHandlerTimeout(cc)
}

// StopCrpcServer force - false(graceful),true(not graceful)
func StopCrpcServer(force bool) {
	if s != nil {
		s.StopCrpcServer(force)
	}
}
