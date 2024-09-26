package xraw

import (
	"crypto/tls"
	"log/slog"

	"github.com/chenjie199234/account/config"
	"github.com/chenjie199234/account/service"

	"github.com/chenjie199234/Corelib/stream"
)

var s *stream.Instance

func StartRawServer() {
	c := config.GetRawServerConfig()
	var tlsc *tls.Config
	if len(c.Certs) > 0 {
		certificates := make([]tls.Certificate, 0, len(c.Certs))
		for cert, key := range c.Certs {
			temp, e := tls.LoadX509KeyPair(cert, key)
			if e != nil {
				slog.ErrorContext(nil, "[xraw] load cert failed:", slog.String("cert", cert), slog.String("key", key), slog.String("error", e.Error()))
				return
			}
			certificates = append(certificates, temp)
		}
		tlsc = &tls.Config{Certificates: certificates}
	}
	s, _ = stream.NewInstance(&stream.InstanceConfig{
		TcpC:               &stream.TcpConfig{ConnectTimeout: c.ConnectTimeout.StdDuration(), MaxMsgLen: c.MaxMsgLen},
		HeartprobeInterval: c.HeartProbe.StdDuration(),
		GroupNum:           c.GroupNum,
		VerifyFunc:         service.SvcRaw.RawVerify,
		OnlineFunc:         service.SvcRaw.RawOnline,
		PingPongFunc:       service.SvcRaw.RawPingPong,
		UserdataFunc:       service.SvcRaw.RawUser,
		OfflineFunc:        service.SvcRaw.RawOffline,
	})

	service.SvcRaw.SetStreamInstance(s)

	if e := s.StartServer(":7000", tlsc); e != nil && e != stream.ErrServerClosed {
		slog.ErrorContext(nil, "[xraw] start server failed", slog.String("error", e.Error()))
		return
	}
	slog.InfoContext(nil, "[xraw] server closed")
}

func StopRawServer() {
	if s != nil {
		s.Stop()
	}
}
