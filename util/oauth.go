package util

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/chenjie199234/account/dao"

	"github.com/chenjie199234/Corelib/cerror"
)

type WeChatResponse struct {
	ErrCode int32  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	OpenID  string `json:"openid"`
}

func OAuthWeChatVerifyCode(ctx context.Context, appid, appsecret, code string) (oauthid string, e error) {
	query := "appid=" + appid + "&secret=" + appsecret + "&code=" + code + "&grant_type=authorization_code"
	r, err := dao.WeChatWebApi.Get(ctx, "/sns/oauth2/access_token", query, nil, nil)
	if err != nil {
		slog.ErrorContext(ctx, "[OAuthWeChatVerifyCode] send request failed", slog.String("code", code), slog.String("error", err.Error()))
		e = err
		return
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "[OAuthWeChatVerifyCode] read response failed", slog.String("code", code), slog.String("error", err.Error()))
		e = err
		return
	}
	resp := &WeChatResponse{}
	if err := json.Unmarshal(body, resp); e != nil {
		slog.ErrorContext(ctx, "[OAuthWeChatVerifyCode] decode response failed", slog.String("code", code), slog.String("error", err.Error()))
		e = err
		return
	}
	if resp.ErrCode != 0 {
		e = cerror.MakeCError(resp.ErrCode, 500, resp.ErrMsg)
		slog.ErrorContext(ctx, "[OAuthWeChatVerifyCode] oauth2 service provider return error",
			slog.String("code", code), slog.String("error", e.Error()))
		return
	}
	return resp.OpenID, nil
}
