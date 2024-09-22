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
	HandlerTimeout     map[string]map[string]ctime.Duration `json:"handler_timeout"`      //first key path,second key method(GET,POST,PUT,PATCH,DELETE,CRPC,GRPC),value timeout
	WebPathRewrite     map[string]map[string]string         `json:"web_path_rewrite"`     //first key method(GET,POST,PUT,PATCH,DELETE),second key origin url,value new url
	HandlerRate        publicmids.MultiPathRateConfigs      `json:"handler_rate"`         //key:path
	Accesses           publicmids.MultiPathAccessConfigs    `json:"accesses"`             //key:path
	TokenSecret        string                               `json:"token_secret"`         //if don't need token check,this can be ingored
	SessionTokenExpire ctime.Duration                       `json:"session_token_expire"` //if don't need session and token check,this can be ignored
	Service            *ServiceConfig                       `json:"service"`
}
type ServiceConfig struct {
	//add your config here
	SupportEmailService []string `json:"support_email_service"`

	//https://open.weixin.qq.com/connect/qrconnect?appid={APPID}&redirect_uri={REDIRECT_URI}&response_type=code&scope=snsapi_login&state={STATE}#wechat_redirect
	WeChatOauthUrl string `json:"wechat_oauth_url"`
	WeChatAppID    string `json:"wechat_appid"`
	WeChatSecret   string `json:"wechat_secret"`
}

// every time update AppConfig will call this function
func validateAppConfig(ac *AppConfig) {
}

// AC -
var AC *AppConfig

var watcher *fsnotify.Watcher

func initlocalapp(notice func(*AppConfig)) {
	data, e := os.ReadFile("./AppConfig.json")
	if e != nil {
		slog.ErrorContext(nil, "[config.local.app] read config file failed", slog.String("error", e.Error()))
		os.Exit(1)
	}
	AC = &AppConfig{}
	if e = json.Unmarshal(data, AC); e != nil {
		slog.ErrorContext(nil, "[config.local.app] config file format wrong", slog.String("error", e.Error()))
		os.Exit(1)
	}
	validateAppConfig(AC)
	slog.InfoContext(nil, "[config.local.app] update success", slog.Any("config", AC))
	if notice != nil {
		notice(AC)
	}
	watcher, e = fsnotify.NewWatcher()
	if e != nil {
		slog.ErrorContext(nil, "[config.local.app] create watcher for hot update failed", slog.String("error", e.Error()))
		os.Exit(1)
	}
	if e = watcher.Add("./"); e != nil {
		slog.ErrorContext(nil, "[config.local.app] create watcher for hot update failed", slog.String("error", e.Error()))
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
					slog.ErrorContext(nil, "[config.local.app] hot update read config file failed", slog.String("error", e.Error()))
					continue
				}
				c := &AppConfig{}
				if e = json.Unmarshal(data, c); e != nil {
					slog.ErrorContext(nil, "[config.local.app] hot update config file format wrong", slog.String("error", e.Error()))
					continue
				}
				validateAppConfig(c)
				AC = c
				slog.InfoContext(nil, "[config.local.app] update success", slog.Any("config", AC))
				if notice != nil {
					notice(AC)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				slog.ErrorContext(nil, "[config.local.app] hot update watcher failed", slog.String("error", err.Error()))
			}
		}
	}()
}
func initremoteapp(notice func(*AppConfig), wait chan *struct{}) (stopwatch func()) {
	return RemoteConfigSdk.Watch("AppConfig", func(key, keyvalue, keytype string) {
		//only support json
		if keytype != "json" {
			slog.ErrorContext(nil, "[config.remote.app] config data can only support json format")
			return
		}
		c := &AppConfig{}
		if e := json.Unmarshal(common.STB(keyvalue), c); e != nil {
			slog.ErrorContext(nil, "[config.remote.app] config data format wrong", slog.String("error", e.Error()))
			return
		}
		validateAppConfig(c)
		AC = c
		slog.InfoContext(nil, "[config.remote.app] update success", slog.Any("config", AC))
		if notice != nil {
			notice(AC)
		}
		select {
		case wait <- nil:
		default:
		}
	})
}
