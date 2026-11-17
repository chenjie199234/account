package xgrpc

import (
	"crypto/tls"
	"log/slog"
	"sync/atomic"

	"github.com/chenjie199234/account/api"
	"github.com/chenjie199234/account/config"
	"github.com/chenjie199234/account/service"

	"github.com/chenjie199234/Corelib/cgrpc"
	"github.com/chenjie199234/Corelib/cgrpc/mids"
	"github.com/chenjie199234/Corelib/util/ctime"
)

var s atomic.Pointer[cgrpc.CGrpcServer]

func StartCGrpcServer() {
	c := config.GetCGrpcServerConfig()
	var tlsc *tls.Config
	if len(c.Certs) > 0 {
		certificates := make([]tls.Certificate, 0, len(c.Certs))
		for cert, key := range c.Certs {
			temp, e := tls.LoadX509KeyPair(cert, key)
			if e != nil {
				slog.Error("[xgrpc] load cert failed:", slog.String("cert", cert), slog.String("key", key), slog.String("error", e.Error()))
				return
			}
			certificates = append(certificates, temp)
		}
		tlsc = &tls.Config{Certificates: certificates}
	}
	server, e := cgrpc.NewCGrpcServer(c.ServerConfig, tlsc)
	if e != nil {
		slog.Error("[xgrpc] new server failed", slog.String("error", e.Error()))
		return
	}
	s.Store(server)

	//this place can register global midwares
	//server.Use(globalmidwares)

	//example
	//api.RegisterExampleCGrpcServer(server, service.SvcExample, mids.AllMids())
	//you need to register your service here
	api.RegisterStatusCGrpcServer(server, service.SvcStatus, mids.AllMids())
	api.RegisterBaseCGrpcServer(server, service.SvcBase, mids.AllMids())
	api.RegisterMoneyCGrpcServer(server, service.SvcMoney, mids.AllMids())

	//path must be registered before this
	UpdateHandlerTimeout(config.AC.HandlerTimeout)

	if e = server.StartCGrpcServer(":10000"); e != nil && e != cgrpc.ErrServerClosed {
		slog.Error("[xgrpc] start server failed", slog.String("error", e.Error()))
		return
	}
	slog.Info("[xgrpc] server closed")
}

// first key:path,second key:method
func UpdateHandlerTimeout(timeout map[string]map[string]ctime.Duration) {
	tmps := s.Load()
	if tmps != nil {
		tmps.UpdateHandlerTimeout(timeout)
	}
}

func StopCGrpcServer(force bool) {
	tmps := s.Load()
	if tmps != nil {
		tmps.StopCGrpcServer(force)
	}
}
