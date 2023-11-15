package util

import (
	"context"

	"github.com/chenjie199234/account/ecode"
)

// caller is the parent function name,this is used for the log
func OAuthVerifyCode(ctx context.Context, caller string, oauthservicename, code string) (oauthid string, e error) {
	switch oauthservicename {
	default:
		return "", ecode.ErrOAuthUnknown
	}
}
