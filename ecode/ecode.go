package ecode

import (
	"net/http"

	"github.com/chenjie199234/Corelib/cerror"
)

var (
	ErrUnknown    = cerror.ErrUnknown    //10000 // http code 500
	ErrReq        = cerror.ErrReq        //10001 // http code 400
	ErrResp       = cerror.ErrResp       //10002 // http code 500
	ErrSystem     = cerror.ErrSystem     //10003 // http code 500
	ErrToken      = cerror.ErrToken      //10004 // http code 401
	ErrSession    = cerror.ErrSession    //10005 // http code 401
	ErrKey        = cerror.ErrKey        //10006 // http code 401
	ErrSign       = cerror.ErrSign       //10007 // http code 401
	ErrPermission = cerror.ErrPermission //10008 // http code 403
	ErrTooFast    = cerror.ErrTooFast    //10009 // http code 403
	ErrBan        = cerror.ErrBan        //10010 // http code 403
	ErrBusy       = cerror.ErrBusy       //10011 // http code 503
	ErrNotExist   = cerror.ErrNotExist   //10012 // http code 404

	ErrDBConflict      = cerror.MakeError(11001, http.StatusInternalServerError, "db data conflict")
	ErrRedisConflict   = cerror.MakeError(11002, http.StatusInternalServerError, "redis data conflict")
	ErrDBRedisConflict = cerror.MakeError(11003, http.StatusInternalServerError, "redis's data and db's data conflict")

	ErrCodeAlreadySend     = cerror.MakeError(20001, http.StatusBadRequest, "dynamic password already send,check your email or tel")
	ErrCodeNotExist        = cerror.MakeError(20002, http.StatusBadRequest, "dynamic password not exist,please get it again")
	ErrUserNotExist        = cerror.MakeError(20003, http.StatusBadRequest, "user not exist")
	ErrTelAlreadyUsed      = cerror.MakeError(20004, http.StatusBadRequest, "tel already used")
	ErrEmailAlreadyUsed    = cerror.MakeError(20005, http.StatusBadRequest, "email already used")
	ErrIDCardAlreadySetted = cerror.MakeError(20006, http.StatusBadRequest, "idcard already setted on this account")
	ErrIDCardAlreadyUsed   = cerror.MakeError(20007, http.StatusBadRequest, "idcard already used")
	ErrNickNameAlreadyUsed = cerror.MakeError(20008, http.StatusBadRequest, "nickname already used")
	ErrPasswordWrong       = cerror.MakeError(20009, http.StatusBadRequest, "password wrong")
	ErrDataBroken          = cerror.MakeError(20010, http.StatusBadRequest, "data broken")
	ErrSignCheckFailed     = cerror.MakeError(20011, http.StatusBadRequest, "sign check failed")
)

func ReturnEcode(originerror error, defaulterror *cerror.Error) error {
	if _, ok := originerror.(*cerror.Error); ok {
		return originerror
	}
	return defaulterror
}
