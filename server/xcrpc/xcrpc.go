package xcrpc

import (
	"crypto/tls"
	"log/slog"
	"sync/atomic"

	"github.com/chenjie199234/account/api"
	"github.com/chenjie199234/account/config"
	"github.com/chenjie199234/account/service"

	"github.com/chenjie199234/Corelib/crpc"
	"github.com/chenjie199234/Corelib/crpc/mids"
	"github.com/chenjie199234/Corelib/util/ctime"
)

var s atomic.Pointer[crpc.CrpcServer]

func StartCrpcServer() {
	c := config.GetCrpcServerConfig()
	var tlsc *tls.Config
	if len(c.Certs) > 0 {
		certificates := make([]tls.Certificate, 0, len(c.Certs))
		for cert, key := range c.Certs {
			temp, e := tls.LoadX509KeyPair(cert, key)
			if e != nil {
				slog.Error("[xcrpc] load cert failed:", slog.String("cert", cert), slog.String("key", key), slog.String("error", e.Error()))
				return
			}
			certificates = append(certificates, temp)
		}
		tlsc = &tls.Config{Certificates: certificates}
	}
	server, e := crpc.NewCrpcServer(c.ServerConfig, tlsc)
	if e != nil {
		slog.Error("[xcrpc] new server failed", slog.String("error", e.Error()))
		return
	}
	s.Store(server)

	//this place can register global midwares
	//server.Use(globalmidwares)

	//example
	//api.RegisterExampleCrpcServer(server, service.SvcExample,mids.AllMids())
	//you need to register your service here
	api.RegisterStatusCrpcServer(server, service.SvcStatus, mids.AllMids())
	api.RegisterBaseCrpcServer(server, service.SvcBase, mids.AllMids())
	api.RegisterMoneyCrpcServer(server, service.SvcMoney, mids.AllMids())

	//path must be registered before this
	UpdateHandlerTimeout(config.AC.HandlerTimeout)

	if e = server.StartCrpcServer(":9000"); e != nil && e != crpc.ErrServerClosed {
		slog.Error("[xcrpc] start server failed", slog.String("error", e.Error()))
		return
	}
	slog.Info("[xcrpc] server closed")
}

// first key:path,second key:method
func UpdateHandlerTimeout(timeout map[string]map[string]ctime.Duration) {
	tmps := s.Load()
	if tmps != nil {
		tmps.UpdateHandlerTimeout(timeout)
	}
}

func StopCrpcServer(force bool) {
	tmps := s.Load()
	if tmps != nil {
		tmps.StopCrpcServer(force)
	}
}
