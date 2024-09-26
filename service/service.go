package service

import (
	"github.com/chenjie199234/account/dao"
	"github.com/chenjie199234/account/service/base"
	"github.com/chenjie199234/account/service/money"
	"github.com/chenjie199234/account/service/raw"
	"github.com/chenjie199234/account/service/status"
)

// SvcStatus one specify sub service
var SvcStatus *status.Service

// SvcRaw one specify sub service
var SvcRaw *raw.Service

var SvcBase *base.Service
var SvcMoney *money.Service

// StartService start the whole service
func StartService() error {
	if e := dao.NewApi(); e != nil {
		return e
	}
	//start sub service
	SvcStatus = status.Start()
	SvcRaw = raw.Start()
	SvcBase = base.Start()
	SvcMoney = money.Start()
	return nil
}

// StopService stop the whole service
func StopService() {
	//stop sub service
	SvcStatus.Stop()
	SvcRaw.Stop()
	SvcBase.Stop()
	SvcMoney.Stop()
}
