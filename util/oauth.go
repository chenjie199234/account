package util

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/chenjie199234/account/config"
	"github.com/chenjie199234/account/dao"
	"github.com/chenjie199234/account/ecode"

	"github.com/chenjie199234/Corelib/cerror"
)

type WeChatResponse struct {
	ErrCode int32  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	OpenID  string `json:"openid"`
}

// caller is the parent function name,this is used for the log
func OAuthVerifyCode(ctx context.Context, caller string, oauthservicename, code string) (oauthid string, e error) {
	c := config.AC.Service
	switch oauthservicename {
	case "wechat":
		if c.WeChatAppID == "" || c.WeChatSecret == "" {
			return "", ecode.ErrOAuthUnknown
		}
		query := "appid=" + c.WeChatAppID + "&secret=" + c.WeChatSecret + "&code=" + code + "&grant_type=authorization_code"
		r, err := dao.WeChatWebApi.Get(ctx, "/sns/oauth2/access_token", query, nil, nil)
		if err != nil {
			slog.ErrorContext(ctx, "[OAuthVerifyCode] call failed", slog.String("oauth_service", oauthservicename), slog.String("code", code), slog.String("error", err.Error()))
			e = err
			return
		}
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.ErrorContext(ctx, "[OAuthVerifyCode] read response body failed", slog.String("oauth_service", oauthservicename), slog.String("code", code), slog.String("error", err.Error()))
			e = err
			return
		}
		resp := &WeChatResponse{}
		if err := json.Unmarshal(body, resp); e != nil {
			slog.ErrorContext(ctx, "[OAuthVerifyCode] response body decode failed", slog.String("oauth_service", oauthservicename), slog.String("code", code), slog.String("error", err.Error()))
			e = err
			return
		}
		if resp.ErrCode != 0 {
			e = cerror.MakeCError(resp.ErrCode, 500, resp.ErrMsg)
			slog.ErrorContext(ctx, "[OAuthVerifyCode] oauth server return error", slog.String("oauth_service", oauthservicename), slog.String("code", code), slog.String("error", e.Error()))
			return
		}
		return resp.OpenID, nil
	default:
		return "", ecode.ErrOAuthUnknown
	}
}
