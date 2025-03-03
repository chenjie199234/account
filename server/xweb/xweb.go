package xweb

import (
	"crypto/tls"
	"log/slog"
	"sync/atomic"
	"unsafe"

	"github.com/chenjie199234/account/api"
	"github.com/chenjie199234/account/config"
	"github.com/chenjie199234/account/service"

	"github.com/chenjie199234/Corelib/util/ctime"
	"github.com/chenjie199234/Corelib/web"
	"github.com/chenjie199234/Corelib/web/mids"
)

var s *web.WebServer

// StartWebServer -
func StartWebServer() {
	c := config.GetWebServerConfig()
	var tlsc *tls.Config
	if len(c.Certs) > 0 {
		certificates := make([]tls.Certificate, 0, len(c.Certs))
		for cert, key := range c.Certs {
			temp, e := tls.LoadX509KeyPair(cert, key)
			if e != nil {
				slog.ErrorContext(nil, "[xweb] load cert failed:", slog.String("cert", cert), slog.String("key", key), slog.String("error", e.Error()))
				return
			}
			certificates = append(certificates, temp)
		}
		tlsc = &tls.Config{Certificates: certificates}
	}
	server, e := web.NewWebServer(c.ServerConfig, tlsc)
	if e != nil {
		slog.ErrorContext(nil, "[xweb] new server failed", slog.String("error", e.Error()))
		return
	}
	//avoid race when build/run in -race mode
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&s)), unsafe.Pointer(server))
	UpdateHandlerTimeout(config.AC.HandlerTimeout)
	UpdateWebPathRewrite(config.AC.WebPathRewrite)

	r := server.NewRouter()

	//this place can register global midwares
	//r.Use(globalmidwares)

	//you just need to register your service here
	api.RegisterStatusWebServer(r, service.SvcStatus, mids.AllMids())
	api.RegisterBaseWebServer(r, service.SvcBase, mids.AllMids())
	api.RegisterMoneyWebServer(r, service.SvcMoney, mids.AllMids())
	//example
	//api.RegisterExampleWebServer(r, service.SvcExample, mids.AllMids())

	server.SetRouter(r)
	if e = server.StartWebServer(":8000"); e != nil && e != web.ErrServerClosed {
		slog.ErrorContext(nil, "[xweb] start server failed", slog.String("error", e.Error()))
		return
	}
	slog.InfoContext(nil, "[xweb] server closed")
}

// UpdateHandlerTimeout -
// first key path,second key method,value timeout duration
func UpdateHandlerTimeout(timeout map[string]map[string]ctime.Duration) {
	//avoid race when build/run in -race mode
	tmps := (*web.WebServer)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&s))))
	if tmps != nil {
		tmps.UpdateHandlerTimeout(timeout)
	}
}

// UpdateWebPathRewrite -
// first key method,second key origin url,value rewrite url
func UpdateWebPathRewrite(rewrite map[string]map[string]string) {
	//avoid race when build/run in -race mode
	tmps := (*web.WebServer)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&s))))
	if tmps != nil {
		tmps.UpdateHandlerRewrite(rewrite)
	}
}

// StopWebServer force - false(graceful),true(not graceful)
func StopWebServer(force bool) {
	//avoid race when build/run in -race mode
	tmps := (*web.WebServer)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&s))))
	if tmps != nil {
		tmps.StopWebServer(force)
	}
}
