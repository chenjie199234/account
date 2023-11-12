package ecode

import (
	"net/http"

	"github.com/chenjie199234/Corelib/cerror"
)

var (
	ErrServerClosing     = cerror.ErrServerClosing     //1000  // http code 449 Warning!! Client will retry on this error,be careful to use this error
	ErrDataConflict      = cerror.ErrDataConflict      //9001  // http code 500
	ErrDataBroken        = cerror.ErrDataBroken        //9002  // http code 500
	ErrDBDataConflict    = cerror.ErrDBDataConflict    //9101  // http code 500
	ErrDBDataBroken      = cerror.ErrDBDataBroken      //9102  // http code 500
	ErrCacheDataConflict = cerror.ErrCacheDataConflict //9201  // http code 500
	ErrCacheDataBroken   = cerror.ErrCacheDataBroken   //9202  // http code 500
	ErrMQDataBroken      = cerror.ErrMQDataBroken      //9301  // http code 500
	ErrUnknown           = cerror.ErrUnknown           //10000 // http code 500
	ErrReq               = cerror.ErrReq               //10001 // http code 400
	ErrResp              = cerror.ErrResp              //10002 // http code 500
	ErrSystem            = cerror.ErrSystem            //10003 // http code 500
	ErrToken             = cerror.ErrToken             //10004 // http code 401
	ErrSession           = cerror.ErrSession           //10005 // http code 401
	ErrAccessKey         = cerror.ErrAccessKey         //10006 // http code 401
	ErrAccessSign        = cerror.ErrAccessSign        //10007 // http code 401
	ErrPermission        = cerror.ErrPermission        //10008 // http code 403
	ErrTooFast           = cerror.ErrTooFast           //10009 // http code 403
	ErrBan               = cerror.ErrBan               //10010 // http code 403
	ErrBusy              = cerror.ErrBusy              //10011 // http code 503
	ErrNotExist          = cerror.ErrNotExist          //10012 // http code 404
	ErrPasswordWrong     = cerror.ErrPasswordWrong     //10013 // http code 400
	ErrPasswordLength    = cerror.ErrPasswordLength    //10014 // http code 400

	ErrUnknownAction       = cerror.MakeError(20000, http.StatusInternalServerError, "unknown action")
	ErrCodeNotExist        = cerror.MakeError(20001, http.StatusBadRequest, "dynamic password not exist,please get it again")
	ErrUserNotExist        = cerror.MakeError(20002, http.StatusBadRequest, "user not exist")
	ErrTelAlreadyUsed      = cerror.MakeError(20003, http.StatusBadRequest, "tel already used")
	ErrEmailAlreadyUsed    = cerror.MakeError(20004, http.StatusBadRequest, "email already used")
	ErrIDCardAlreadyUsed   = cerror.MakeError(20005, http.StatusBadRequest, "idcard already used")
	ErrNickNameAlreadyUsed = cerror.MakeError(20006, http.StatusBadRequest, "nickname already used")
	ErrSignCheckFailed     = cerror.MakeError(20007, http.StatusBadRequest, "sign check failed")
	ErrOAuthWrong          = cerror.MakeError(20008, http.StatusBadRequest, "oauth wrong")
)

func ReturnEcode(originerror error, defaulterror *cerror.Error) error {
	if _, ok := originerror.(*cerror.Error); ok {
		return originerror
	}
	return defaulterror
}
