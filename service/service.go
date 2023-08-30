package service

import (
	"github.com/chenjie199234/account/dao"
	"github.com/chenjie199234/account/service/money"
	"github.com/chenjie199234/account/service/status"
	"github.com/chenjie199234/account/service/user"
)

// SvcStatus one specify sub service
var SvcStatus *status.Service
var SvcUser *user.Service
var SvcMoney *money.Service

// StartService start the whole service
func StartService() error {
	if e := dao.NewApi(); e != nil {
		return e
	}
	//start sub service
	SvcStatus = status.Start()
	SvcUser = user.Start()
	SvcMoney = money.Start()
	return nil
}

// StopService stop the whole service
func StopService() {
	//stop sub service
	SvcStatus.Stop()
	SvcUser.Stop()
	SvcMoney.Stop()
}
