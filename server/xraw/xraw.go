package xraw

import (
	"crypto/tls"
	"log/slog"
	"sync/atomic"

	"github.com/chenjie199234/account/config"
	"github.com/chenjie199234/account/service"

	"github.com/chenjie199234/Corelib/stream"
)

var s atomic.Pointer[stream.Instance]

func StartRawServer() {
	c := config.GetRawServerConfig()
	var tlsc *tls.Config
	if len(c.Certs) > 0 {
		certificates := make([]tls.Certificate, 0, len(c.Certs))
		for cert, key := range c.Certs {
			temp, e := tls.LoadX509KeyPair(cert, key)
			if e != nil {
				slog.Error("[xraw] load cert failed:", slog.String("cert", cert), slog.String("key", key), slog.String("error", e.Error()))
				return
			}
			certificates = append(certificates, temp)
		}
		tlsc = &tls.Config{Certificates: certificates}
	}
	server, _ := stream.NewInstance(&stream.InstanceConfig{
		HeartprobeInterval: c.HeartProbe.StdDuration(),
		ConnectTimeout:     c.ConnectTimeout.StdDuration(),
		ReadTimeout:        c.ReadTimeout.StdDuration(),
		WriteTimeout:       c.WriteTimeout.StdDuration(),
		IdleTimeout:        c.IdleTimeout.StdDuration(),
		GroupNum:           c.GroupNum,
		MaxMsgLen:          c.MaxMsgLen,
		VerifyFunc:         service.SvcRaw.RawVerify,
		OnlineFunc:         service.SvcRaw.RawOnline,
		PingPongFunc:       service.SvcRaw.RawPingPong,
		UserdataFunc:       service.SvcRaw.RawUser,
		OfflineFunc:        service.SvcRaw.RawOffline,
	})
	s.Store(server)

	service.SvcRaw.SetStreamInstance(server)

	if e := server.StartServer(":7000", tlsc); e != nil && e != stream.ErrServerClosed {
		slog.Error("[xraw] start server failed", slog.String("error", e.Error()))
		return
	}
	slog.Info("[xraw] server closed")
}

func StopRawServer() {
	tmps := s.Load()
	if tmps != nil {
		tmps.Stop()
	}
}
