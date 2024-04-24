package util

import (
	"context"
	"encoding/json"
	"io"

	"github.com/chenjie199234/account/config"
	"github.com/chenjie199234/account/dao"
	"github.com/chenjie199234/account/ecode"

	"github.com/chenjie199234/Corelib/cerror"
	"github.com/chenjie199234/Corelib/log"
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
		query := "appid=" + c.WeChatAppID + "&secret=" + c.WeChatSecret + "&code=" + code + "&grant_type=authorization_code"
		r, err := dao.WeChatWebApi.Get(ctx, "/sns/oauth2/access_token", query, nil, nil)
		if err != nil {
			log.Error(ctx, "[OAuthVerifyCode] call failed", log.String("oauth_service", oauthservicename), log.String("code", code), log.CError(err))
			e = err
			return
		}
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error(ctx, "[OAuthVerifyCode] read response body failed", log.String("oauth_service", oauthservicename), log.String("code", code), log.CError(err))
			e = err
			return
		}
		resp := &WeChatResponse{}
		if err := json.Unmarshal(body, resp); e != nil {
			log.Error(ctx, "[OAuthVerifyCode] response body decode failed", log.String("oauth_service", oauthservicename), log.String("code", code), log.CError(err))
			e = err
			return
		}
		if resp.ErrCode != 0 {
			e = cerror.MakeError(resp.ErrCode, 500, resp.ErrMsg)
			log.Error(ctx, "[OAuthVerifyCode] oauth server return error", log.String("oauth_service", oauthservicename), log.String("code", code), log.CError(e))
			return
		}
		return resp.OpenID, nil
	default:
		return "", ecode.ErrOAuthUnknown
	}
}
