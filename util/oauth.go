package util

import (
	"context"
)

func OauthVerifyCode(ctx context.Context, oauthservicename, code string) (oauthid string) {
	switch oauthservicename {
	case "wechat":
	}
	return
}
