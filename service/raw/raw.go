package raw

import (
	"context"
	"sync/atomic"
	// "log/slog"

	// "github.com/chenjie199234/account/config"
	// "github.com/chenjie199234/account/api"
	rawdao "github.com/chenjie199234/account/dao/raw"
	// "github.com/chenjie199234/account/ecode"

	"github.com/chenjie199234/Corelib/stream"
	"github.com/chenjie199234/Corelib/util/graceful"
)

// Service subservice for raw business
type Service struct {
	stop *graceful.Graceful

	rawDao   *rawdao.Dao
	instance atomic.Pointer[stream.Instance]
}

// Start -
func Start() (*Service, error) {
	return &Service{
		stop: graceful.New(),

		//rawDao: rawdao.NewDao(config.GetMysql("raw_mysql"), config.GetRedis("raw_redis"), config.GetMongo("raw_mongo")),
		rawDao: rawdao.NewDao(nil, nil, nil),
	}, nil
}

func (s *Service) SetStreamInstance(instance *stream.Instance) {
	s.instance.Store(instance)
}

func (s *Service) RawVerify(ctx context.Context, peerip string, peerVerifyData []byte) (response []byte, uniqueid string, success bool) {
	// instance := s.instance.Load()
	return nil, "", false
}

func (s *Service) RawOnline(ctx context.Context, p *stream.Peer) (success bool) {
	// instance := s.instance.Load()
	return false
}

func (s *Service) RawPingPong(p *stream.Peer) {
	// instance := s.instance.Load()
}

func (s *Service) RawUser(p *stream.Peer, userdata []byte) {
	// instance := s.instance.Load()
}

func (s *Service) RawOffline(p *stream.Peer) {
	// instance := s.instance.Load()
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
