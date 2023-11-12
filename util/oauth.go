package util

import (
	"context"
)

func OAuthVerifyCode(ctx context.Context, caller string, oauthservicename, code string) (oauthid string, e error) {
	switch oauthservicename {
	case "wechat":
	}
	return
}
