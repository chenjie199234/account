package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/chenjie199234/account/model"

	configsdk "github.com/chenjie199234/admin/sdk/config"
)

// EnvConfig can't hot update,all these data is from system env setting
// nil field means that system env not exist
type EnvConfig struct {
	ConfigType *int
	RunEnv     *string
	DeployEnv  *string
}

// EC -
var EC *EnvConfig

// RemoteConfigSdk -
var RemoteConfigSdk *configsdk.ConfigSdk

// notice is a sync function
// don't write block logic inside it
func Init(notice func(c *AppConfig)) {
	initenv()
	if EC.ConfigType != nil && *EC.ConfigType == 1 {
		tmer := time.NewTimer(time.Second * 2)
		waitapp := make(chan *struct{}, 1)
		waitsource := make(chan *struct{}, 1)
		initremoteapp(notice, waitapp)
		stopwatchsource := initremotesource(waitsource)
		appinit := false
		sourceinit := false
		for {
			select {
			case <-waitapp:
				appinit = true
			case <-waitsource:
				sourceinit = true
				stopwatchsource()
			case <-tmer.C:
				slog.ErrorContext(nil, "[config.Init] timeout")
				os.Exit(1)
			}
			if appinit && sourceinit {
				break
			}
		}
	} else {
		initlocalapp(notice)
		initlocalsource()
	}
}

func initenv() {
	EC = &EnvConfig{}
	if str, ok := os.LookupEnv("CONFIG_TYPE"); ok && str != "<CONFIG_TYPE>" && str != "" {
		configtype, e := strconv.Atoi(str)
		if e != nil || (configtype != 0 && configtype != 1 && configtype != 2) {
			slog.ErrorContext(nil, "[config.initenv] env CONFIG_TYPE must be number in [0,1,2]")
			os.Exit(1)
		}
		EC.ConfigType = &configtype
	} else {
		slog.WarnContext(nil, "[config.initenv] missing env CONFIG_TYPE")
	}
	if EC.ConfigType != nil && *EC.ConfigType == 1 {
		var e error
		if RemoteConfigSdk, e = configsdk.NewConfigSdk(model.Project, model.Group, model.Name, nil); e != nil {
			slog.ErrorContext(nil, "[config.initenv] new remote config sdk failed", slog.String("error", e.Error()))
			os.Exit(1)
		}
	}
	if str, ok := os.LookupEnv("RUN_ENV"); ok && str != "<RUN_ENV>" && str != "" {
		EC.RunEnv = &str
	} else {
		slog.WarnContext(nil, "[config.initenv] missing env RUN_ENV")
	}
	if str, ok := os.LookupEnv("DEPLOY_ENV"); ok && str != "<DEPLOY_ENV>" && str != "" {
		EC.DeployEnv = &str
	} else {
		slog.WarnContext(nil, "[config.initenv] missing env DEPLOY_ENV")
	}
}
