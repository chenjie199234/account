package xweb

import (
	"crypto/tls"
	"log/slog"
	"sync/atomic"

	"github.com/chenjie199234/account/api"
	"github.com/chenjie199234/account/config"
	"github.com/chenjie199234/account/service"

	"github.com/chenjie199234/Corelib/util/ctime"
	"github.com/chenjie199234/Corelib/web"
	"github.com/chenjie199234/Corelib/web/mids"
)

var s atomic.Pointer[web.WebServer]
var r atomic.Pointer[web.Router]

func StartWebServer() {
	c := config.GetWebServerConfig()
	var tlsc *tls.Config
	if len(c.Certs) > 0 {
		certificates := make([]tls.Certificate, 0, len(c.Certs))
		for cert, key := range c.Certs {
			temp, e := tls.LoadX509KeyPair(cert, key)
			if e != nil {
				slog.Error("[xweb] load cert failed:", slog.String("cert", cert), slog.String("key", key), slog.String("error", e.Error()))
				return
			}
			certificates = append(certificates, temp)
		}
		tlsc = &tls.Config{Certificates: certificates}
	}
	server, e := web.NewWebServer(c.ServerConfig, tlsc)
	if e != nil {
		slog.Error("[xweb] new server failed", slog.String("error", e.Error()))
		return
	}
	s.Store(server)

	router, e := server.NewRouter()
	if e != nil {
		slog.Error("[xweb] new router failed", slog.String("error", e.Error()))
		return
	}
	r.Store(router)

	//this place can register global midwares
	//router.Use(globalmidwares)

	//example
	//api.RegisterExampleWebServer(router, service.SvcExample, mids.AllMids())
	//you need to register your service here
	api.RegisterStatusWebServer(router, service.SvcStatus, mids.AllMids())
	api.RegisterBaseWebServer(router, service.SvcBase, mids.AllMids())
	api.RegisterMoneyWebServer(router, service.SvcMoney, mids.AllMids())

	server.SetRouter(router)

	//path must be registered before this
	UpdateHandlerTimeout(config.AC.HandlerTimeout)
	UpdateWebPathRewrite(config.AC.WebPathRewrite)

	if e = server.StartWebServer(":8000"); e != nil && e != web.ErrServerClosed {
		slog.Error("[xweb] start server failed", slog.String("error", e.Error()))
		return
	}
	slog.Info("[xweb] server closed")
}

// first key:path,second key:method
func UpdateHandlerTimeout(timeout map[string]map[string]ctime.Duration) {
	tmpr := r.Load()
	if tmpr != nil {
		tmpr.UpdateHandlerTimeout(timeout)
	}
}

// first key:method,second key:origin url,value:new url
func UpdateWebPathRewrite(rewrite map[string]map[string]string) {
	tmpr := r.Load()
	if tmpr != nil {
		tmpr.UpdateHandlerRewrite(rewrite)
	}
}

func StopWebServer(force bool) {
	tmps := s.Load()
	if tmps != nil {
		tmps.StopWebServer(force)
	}
}
