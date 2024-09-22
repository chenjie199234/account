package xgrpc

import (
	"crypto/tls"
	"log/slog"

	"github.com/chenjie199234/account/api"
	"github.com/chenjie199234/account/config"
	"github.com/chenjie199234/account/model"
	"github.com/chenjie199234/account/service"

	"github.com/chenjie199234/Corelib/cgrpc"
	"github.com/chenjie199234/Corelib/cgrpc/mids"
	"github.com/chenjie199234/Corelib/util/ctime"
)

var s *cgrpc.CGrpcServer

// StartCGrpcServer -
func StartCGrpcServer() {
	c := config.GetCGrpcServerConfig()
	var tlsc *tls.Config
	if len(c.Certs) > 0 {
		certificates := make([]tls.Certificate, 0, len(c.Certs))
		for cert, key := range c.Certs {
			temp, e := tls.LoadX509KeyPair(cert, key)
			if e != nil {
				slog.ErrorContext(nil, "[xgrpc] load cert failed:", slog.String("cert", cert), slog.String("key", key), slog.String("error", e.Error()))
				return
			}
			certificates = append(certificates, temp)
		}
		tlsc = &tls.Config{Certificates: certificates}
	}
	var e error
	if s, e = cgrpc.NewCGrpcServer(c.ServerConfig, model.Project, model.Group, model.Name, tlsc); e != nil {
		slog.ErrorContext(nil, "[xgrpc] new server failed", slog.String("error", e.Error()))
		return
	}
	UpdateHandlerTimeout(config.AC.HandlerTimeout)

	//this place can register global midwares
	//s.Use(globalmidwares)

	//you just need to register your service here
	api.RegisterStatusCGrpcServer(s, service.SvcStatus, mids.AllMids())
	api.RegisterBaseCGrpcServer(s, service.SvcBase, mids.AllMids())
	api.RegisterMoneyCGrpcServer(s, service.SvcMoney, mids.AllMids())

	//example
	//api.RegisterExampleCGrpcServer(s, service.SvcExample, mids.AllMids())

	if e = s.StartCGrpcServer(":10000"); e != nil && e != cgrpc.ErrServerClosed {
		slog.ErrorContext(nil, "[xgrpc] start server failed", slog.String("error", e.Error()))
		return
	}
	slog.InfoContext(nil, "[xgrpc] server closed")
}

// UpdateHandlerTimeout -
// first key path,second key method,value timeout duration
func UpdateHandlerTimeout(timeout map[string]map[string]ctime.Duration) {
	if s != nil {
		s.UpdateHandlerTimeout(timeout)
	}
}

// StopCGrpcServer force - false(graceful),true(not graceful)
func StopCGrpcServer(force bool) {
	if s != nil {
		s.StopCGrpcServer(force)
	}
}
