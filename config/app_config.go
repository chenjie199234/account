package config

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"

	publicmids "github.com/chenjie199234/Corelib/mids"
	"github.com/chenjie199234/Corelib/util/common"
	"github.com/chenjie199234/Corelib/util/ctime"
	"github.com/fsnotify/fsnotify"
)

// AppConfig can hot update
// this is the config used for this app
type AppConfig struct {
	HandlerTimeout map[string]map[string]ctime.Duration `json:"handler_timeout"`  //first key path,second key method(GET,POST,PUT,PATCH,DELETE,CRPC,GRPC),value timeout
	WebPathRewrite map[string]map[string]string         `json:"web_path_rewrite"` //first key method(GET,POST,PUT,PATCH,DELETE),second key origin url,value new url
	HandlerRate    publicmids.MultiPathRateConfigs      `json:"handler_rate"`     //key:path
	Accesses       publicmids.MultiPathAccessConfigs    `json:"accesses"`         //key:path
	TokenSecret    string                               `json:"token_secret"`
	Service        *ServiceConfig                       `json:"service"`
}
type ServiceConfig struct {
	//add your config here
	TokenExpire ctime.Duration `json:"token_expire"`

	SupportEmailService []string `json:"support_email_service"`

	//https://open.weixin.qq.com/connect/qrconnect?appid={APPID}&redirect_uri={REDIRECT_URI}&response_type=code&scope=snsapi_login&state={STATE}#wechat_redirect
	WeChatOauth2    string `json:"wechat_oauth2"`
	WeChatAppID     string `json:"wechat_app_id"`
	WeChatAppSecret string `json:"wechat_app_secret"`
}

// every time update AppConfig will call this function
func validateAppConfig(ac *AppConfig) {
	if ac.Service.WeChatOauth2 != "" && (ac.Service.WeChatAppID == "" || ac.Service.WeChatAppSecret == "") {
		slog.Error("[config.validateAppConfig] missing wechat_app_id or wechat_app_secret")
		os.Exit(1)
	}
}

// AC -
var AC *AppConfig

var watcher *fsnotify.Watcher

func initlocalapp(notice func(*AppConfig)) {
	data, e := os.ReadFile("./AppConfig.json")
	if e != nil {
		slog.Error("[config.local.app] read config file failed", slog.String("error", e.Error()))
		os.Exit(1)
	}
	AC = &AppConfig{}
	if e = json.Unmarshal(data, AC); e != nil {
		slog.Error("[config.local.app] config file format wrong", slog.String("error", e.Error()))
		os.Exit(1)
	}
	validateAppConfig(AC)
	slog.Info("[config.local.app] update success", slog.Any("config", AC))
	if notice != nil {
		notice(AC)
	}
	watcher, e = fsnotify.NewWatcher()
	if e != nil {
		slog.Error("[config.local.app] create watcher for hot update failed", slog.String("error", e.Error()))
		os.Exit(1)
	}
	if e = watcher.Add("./"); e != nil {
		slog.Error("[config.local.app] create watcher for hot update failed", slog.String("error", e.Error()))
		os.Exit(1)
	}
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if filepath.Base(event.Name) != "AppConfig.json" || (!event.Has(fsnotify.Create) && !event.Has(fsnotify.Write)) {
					continue
				}
				data, e := os.ReadFile("./AppConfig.json")
				if e != nil {
					slog.Error("[config.local.app] hot update read config file failed", slog.String("error", e.Error()))
					continue
				}
				c := &AppConfig{}
				if e = json.Unmarshal(data, c); e != nil {
					slog.Error("[config.local.app] hot update config file format wrong", slog.String("error", e.Error()))
					continue
				}
				validateAppConfig(c)
				AC = c
				slog.Info("[config.local.app] update success", slog.Any("config", AC))
				if notice != nil {
					notice(AC)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				slog.Error("[config.local.app] hot update watcher failed", slog.String("error", err.Error()))
			}
		}
	}()
}
func initremoteapp(notice func(*AppConfig), wait chan *struct{}) (stopwatch func()) {
	return RemoteConfigSdk.Watch("AppConfig", func(key, keyvalue, keytype string) {
		//only support json
		if keytype != "json" {
			slog.Error("[config.remote.app] config data can only support json format")
			return
		}
		c := &AppConfig{}
		if e := json.Unmarshal(common.STB(keyvalue), c); e != nil {
			slog.Error("[config.remote.app] config data format wrong", slog.String("error", e.Error()))
			return
		}
		validateAppConfig(c)
		slog.Info("[config.remote.app] update success", slog.Any("config", c))
		AC = c
		if notice != nil {
			notice(AC)
		}
		select {
		case wait <- nil:
		default:
		}
	})
}
