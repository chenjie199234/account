package util

import (
	"context"
)

func OauthVerifyCode(ctx context.Context, caller string, oauthservicename, code string) (oauthid string, e error) {
	switch oauthservicename {
	case "wechat":
	}
	return
}
