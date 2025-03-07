package main

import (
	"context"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/chenjie199234/account/config"
	"github.com/chenjie199234/account/dao"
	_ "github.com/chenjie199234/account/model"
	"github.com/chenjie199234/account/server/xcrpc"
	"github.com/chenjie199234/account/server/xgrpc"
	"github.com/chenjie199234/account/server/xraw"
	"github.com/chenjie199234/account/server/xweb"
	"github.com/chenjie199234/account/service"

	"github.com/chenjie199234/Corelib/cotel"
	publicmids "github.com/chenjie199234/Corelib/mids"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/redis/go-redis/v9"
	_ "go.mongodb.org/mongo-driver/v2/mongo"
)

type LogHandler struct {
	slog.Handler
}

func (l *LogHandler) Handle(ctx context.Context, record slog.Record) error {
	if record.NumAttrs() > 0 {
		attrs := make([]slog.Attr, 0, record.NumAttrs())
		record.Attrs(func(a slog.Attr) bool {
			attrs = append(attrs, a)
			return true
		})
		if record.Message == "trace" {
			record.PC = 0
		}
		record = slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
		record.AddAttrs(slog.Attr{
			Key:   "msg_kvs",
			Value: slog.GroupValue(attrs...),
		})
	}
	if traceid := cotel.TraceIDFromContext(ctx); traceid != "" {
		record.AddAttrs(slog.String("traceid", traceid))
	}
	return l.Handler.Handle(ctx, record)
}

func main() {
	slog.SetDefault(slog.New(&LogHandler{
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
				if len(groups) == 0 && attr.Key == "function" {
					return slog.Attr{}
				}
				if len(groups) == 0 && attr.Key == slog.SourceKey {
					s := attr.Value.Any().(*slog.Source)
					if index := strings.Index(s.File, "corelib@v"); index != -1 {
						s.File = s.File[index:]
					} else if index = strings.Index(s.File, "Corelib@v"); index != -1 {
						s.File = s.File[index:]
					}
				}
				return attr
			},
		}),
	}))
	if e := cotel.Init(); e != nil {
		slog.Error("init cotel failed", slog.String("error", e.Error()))
		return
	}
	defer cotel.Stop()
	config.Init(func(ac *config.AppConfig) {
		//this is a notice callback every time appconfig changes
		//this function works in sync mode
		//don't write block logic inside this
		dao.UpdateAPI(ac)
		xcrpc.UpdateHandlerTimeout(ac.HandlerTimeout)
		xgrpc.UpdateHandlerTimeout(ac.HandlerTimeout)
		xweb.UpdateHandlerTimeout(ac.HandlerTimeout)
		xweb.UpdateWebPathRewrite(ac.WebPathRewrite)
		publicmids.UpdateRateConfig(ac.HandlerRate)
		publicmids.UpdateTokenConfig(ac.TokenSecret)
		publicmids.UpdateAccessConfig(ac.Accesses)
	})
	if rateredis := config.GetRedis("rate_redis"); rateredis != nil {
		publicmids.UpdateRateRedisInstance(rateredis)
	} else {
		slog.WarnContext(nil, "[main] rate redis missing,all rate check will be failed")
	}
	if sessionredis := config.GetRedis("session_redis"); sessionredis != nil {
		publicmids.UpdateSessionRedisInstance(sessionredis)
	} else {
		slog.WarnContext(nil, "[main] session redis missing,all session event will be failed")
	}
	//start the whole business service
	if e := service.StartService(); e != nil {
		slog.ErrorContext(nil, "[main] start service failed", slog.String("error", e.Error()))
		return
	}
	//start low level net service
	ch := make(chan os.Signal, 1)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		xcrpc.StartCrpcServer()
		select {
		case ch <- syscall.SIGTERM:
		default:
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		xweb.StartWebServer()
		select {
		case ch <- syscall.SIGTERM:
		default:
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		xgrpc.StartCGrpcServer()
		select {
		case ch <- syscall.SIGTERM:
		default:
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		xraw.StartRawServer()
		select {
		case ch <- syscall.SIGTERM:
		default:
		}
		wg.Done()
	}()
	//this is the server for pprof and prometheus(if METRIC env is prometheus)
	//this server should not be exposed to the public internet
	pserver := &http.Server{Addr: ":6060"}
	wg.Add(1)
	go func() {
		if h := cotel.GetPrometheusHandler(); h != nil {
			http.Handle("/metrics", h)
		}
		pserver.ListenAndServe()
		select {
		case ch <- syscall.SIGTERM:
		default:
		}
		wg.Done()
	}()
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-ch
	//stop the whole business service
	service.StopService()
	//stop low level net service
	wg.Add(1)
	go func() {
		xcrpc.StopCrpcServer(false)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		xweb.StopWebServer(false)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		xgrpc.StopCGrpcServer(false)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		xraw.StopRawServer()
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		pserver.Shutdown(context.Background())
		wg.Done()
	}()
	wg.Wait()
}
