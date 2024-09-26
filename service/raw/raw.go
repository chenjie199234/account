package raw

import (
	"context"
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
	instance *stream.Instance
}

// Start -
func Start() *Service {
	return &Service{
		stop: graceful.New(),

		//rawDao: rawdao.NewDao(config.GetMysql("raw_mysql"), config.GetRedis("raw_redis"), config.GetMongo("raw_mongo")),
		rawDao: rawdao.NewDao(nil, nil, nil),
	}
}

func (s *Service) SetStreamInstance(instance *stream.Instance) {
	s.instance = instance
}

func (s *Service) RawVerify(ctx context.Context, peerVerifyData []byte) (response []byte, uniqueid string, success bool) {
	return nil, "", false
}

func (s *Service) RawOnline(ctx context.Context, p *stream.Peer) (success bool) {
	return false
}

func (s *Service) RawPingPong(p *stream.Peer) {
}

func (s *Service) RawUser(p *stream.Peer, userdata []byte) {
}

func (s *Service) RawOffline(p *stream.Peer) {
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
