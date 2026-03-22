package service

//Warning!!Don't add comments in this file
//this file will be updated automaticly when create sub service
//however,golang's ast package can't handle ast tree with comments well(checked 1.26.0)

import (
	"github.com/chenjie199234/account/dao"
	"github.com/chenjie199234/account/service/base"
	"github.com/chenjie199234/account/service/money"
	"github.com/chenjie199234/account/service/raw"
	"github.com/chenjie199234/account/service/status"
)

var SvcMoney *money.Service

var SvcBase *base.Service

var SvcStatus *status.Service

var SvcRaw *raw.Service

func StartService() error {
	var e error
	if e = dao.NewApi(); e != nil {
		return e
	}
	if SvcStatus, e = status.Start(); e != nil {
		return e
	}
	if SvcRaw, e = raw.Start(); e != nil {
		return e
	}
	if SvcBase, e = base.Start(); e != nil {
		return e
	}
	if SvcMoney, e = money.Start(); e != nil {
		return e
	}
	return nil
}

func StopService() {
	SvcStatus.Stop()
	SvcRaw.Stop()
	SvcBase.Stop()
	SvcMoney.Stop()
}
