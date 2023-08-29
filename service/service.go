package service

import (
	"github.com/chenjie199234/account/dao"
	"github.com/chenjie199234/account/service/status"
)

// SvcStatus one specify sub service
var SvcStatus *status.Service

// StartService start the whole service
func StartService() error {
	if e := dao.NewApi(); e != nil {
		return e
	}
	//start sub service
	SvcStatus = status.Start()
	return nil
}

// StopService stop the whole service
func StopService() {
	//stop sub service
	SvcStatus.Stop()
}