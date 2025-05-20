package raw

import (
	"context"
	"sync/atomic"
	"unsafe"
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
func Start() (*Service, error) {
	return &Service{
		stop: graceful.New(),

		//rawDao: rawdao.NewDao(config.GetMysql("raw_mysql"), config.GetRedis("raw_redis"), config.GetMongo("raw_mongo")),
		rawDao: rawdao.NewDao(nil, nil, nil),
	}, nil
}

func (s *Service) SetStreamInstance(instance *stream.Instance) {
	//avoid race when build/run in -race mode
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&s.instance)), unsafe.Pointer(instance))
}

func (s *Service) RawVerify(ctx context.Context, peerVerifyData []byte) (response []byte, uniqueid string, success bool) {
	//avoid race when build/run in -race mode
	// instance := (*stream.Instance)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&s.instance))))
	return nil, "", false
}

func (s *Service) RawOnline(ctx context.Context, p *stream.Peer) (success bool) {
	//avoid race when build/run in -race mode
	// instance := (*stream.Instance)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&s.instance))))
	return false
}

func (s *Service) RawPingPong(p *stream.Peer) {
	//avoid race when build/run in -race mode
	// instance := (*stream.Instance)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&s.instance))))
}

func (s *Service) RawUser(p *stream.Peer, userdata []byte) {
	//avoid race when build/run in -race mode
	// instance := (*stream.Instance)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&s.instance))))
}

func (s *Service) RawOffline(p *stream.Peer) {
	//avoid race when build/run in -race mode
	// instance := (*stream.Instance)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&s.instance))))
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
