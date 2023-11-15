package util

import (
	"context"

	"github.com/chenjie199234/account/ecode"
)

func OAuthVerifyCode(ctx context.Context, caller string, oauthservicename, code string) (oauthid string, e error) {
	switch oauthservicename {
	default:
		return "", ecode.ErrOAuthUnknown
	}
	return
}
