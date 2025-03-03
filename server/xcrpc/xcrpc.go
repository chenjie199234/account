package xcrpc

import (
	"crypto/tls"
	"log/slog"
	"sync/atomic"
	"unsafe"

	"github.com/chenjie199234/account/api"
	"github.com/chenjie199234/account/config"
	"github.com/chenjie199234/account/service"

	"github.com/chenjie199234/Corelib/crpc"
	"github.com/chenjie199234/Corelib/crpc/mids"
	"github.com/chenjie199234/Corelib/util/ctime"
)

var s *crpc.CrpcServer

// StartCrpcServer -
func StartCrpcServer() {
	c := config.GetCrpcServerConfig()
	var tlsc *tls.Config
	if len(c.Certs) > 0 {
		certificates := make([]tls.Certificate, 0, len(c.Certs))
		for cert, key := range c.Certs {
			temp, e := tls.LoadX509KeyPair(cert, key)
			if e != nil {
				slog.ErrorContext(nil, "[xcrpc] load cert failed:", slog.String("cert", cert), slog.String("key", key), slog.String("error", e.Error()))
				return
			}
			certificates = append(certificates, temp)
		}
		tlsc = &tls.Config{Certificates: certificates}
	}
	server, e := crpc.NewCrpcServer(c.ServerConfig, tlsc)
	if e != nil {
		slog.ErrorContext(nil, "[xcrpc] new server failed", slog.String("error", e.Error()))
		return
	}
	//avoid race when build/run in -race mode
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&s)), unsafe.Pointer(server))
	UpdateHandlerTimeout(config.AC.HandlerTimeout)

	//this place can register global midwares
	//server.Use(globalmidwares)

	//you just need to register your service here
	api.RegisterStatusCrpcServer(server, service.SvcStatus, mids.AllMids())
	api.RegisterBaseCrpcServer(server, service.SvcBase, mids.AllMids())
	api.RegisterMoneyCrpcServer(server, service.SvcMoney, mids.AllMids())
	//example
	//api.RegisterExampleCrpcServer(server, service.SvcExample,mids.AllMids())

	if e = server.StartCrpcServer(":9000"); e != nil && e != crpc.ErrServerClosed {
		slog.ErrorContext(nil, "[xcrpc] start server failed", slog.String("error", e.Error()))
		return
	}
	slog.InfoContext(nil, "[xcrpc] server closed")
}

// UpdateHandlerTimeout -
// first key path,second key method,value timeout duration
func UpdateHandlerTimeout(timeout map[string]map[string]ctime.Duration) {
	//avoid race when build/run in -race mode
	tmps := (*crpc.CrpcServer)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&s))))
	if tmps != nil {
		tmps.UpdateHandlerTimeout(timeout)
	}
}

// StopCrpcServer force - false(graceful),true(not graceful)
func StopCrpcServer(force bool) {
	//avoid race when build/run in -race mode
	tmps := (*crpc.CrpcServer)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&s))))
	if tmps != nil {
		tmps.StopCrpcServer(force)
	}
}
