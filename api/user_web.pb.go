// Code generated by protoc-gen-go-web. DO NOT EDIT.
// version:
// 	protoc-gen-go-web v0.0.78<br />
// 	protoc            v4.24.1<br />
// source: api/user.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
	web "github.com/chenjie199234/Corelib/web"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	io "io"
	http "net/http"
	strings "strings"
)

var _WebPathUserLogin = "/account.user/login"
var _WebPathUserSelfUserInfo = "/account.user/self_user_info"
var _WebPathUserUpdateStaticPassword = "/account.user/update_static_password"
var _WebPathUserUpdateIdcard = "/account.user/update_idcard"
var _WebPathUserUpdateNickName = "/account.user/update_nick_name"
var _WebPathUserUpdateEmail = "/account.user/update_email"
var _WebPathUserUpdateTel = "/account.user/update_tel"

type UserWebClient interface {
	Login(context.Context, *LoginReq, http.Header) (*LoginResp, error)
	SelfUserInfo(context.Context, *SelfUserInfoReq, http.Header) (*SelfUserInfoResp, error)
	UpdateStaticPassword(context.Context, *UpdateStaticPasswordReq, http.Header) (*UpdateStaticPasswordResp, error)
	UpdateIdcard(context.Context, *UpdateIdcardReq, http.Header) (*UpdateIdcardResp, error)
	UpdateNickName(context.Context, *UpdateNickNameReq, http.Header) (*UpdateNickNameResp, error)
	UpdateEmail(context.Context, *UpdateEmailReq, http.Header) (*UpdateEmailResp, error)
	UpdateTel(context.Context, *UpdateTelReq, http.Header) (*UpdateTelResp, error)
}

type userWebClient struct {
	cc *web.WebClient
}

func NewUserWebClient(c *web.WebClient) UserWebClient {
	return &userWebClient{cc: c}
}

func (c *userWebClient) Login(ctx context.Context, req *LoginReq, header http.Header) (*LoginResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathUserLogin, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(LoginResp)
	if len(data) == 0 {
		return resp, nil
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-protobuf") {
		if e := proto.Unmarshal(data, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *userWebClient) SelfUserInfo(ctx context.Context, req *SelfUserInfoReq, header http.Header) (*SelfUserInfoResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathUserSelfUserInfo, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(SelfUserInfoResp)
	if len(data) == 0 {
		return resp, nil
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-protobuf") {
		if e := proto.Unmarshal(data, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *userWebClient) UpdateStaticPassword(ctx context.Context, req *UpdateStaticPasswordReq, header http.Header) (*UpdateStaticPasswordResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathUserUpdateStaticPassword, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(UpdateStaticPasswordResp)
	if len(data) == 0 {
		return resp, nil
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-protobuf") {
		if e := proto.Unmarshal(data, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *userWebClient) UpdateIdcard(ctx context.Context, req *UpdateIdcardReq, header http.Header) (*UpdateIdcardResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathUserUpdateIdcard, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(UpdateIdcardResp)
	if len(data) == 0 {
		return resp, nil
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-protobuf") {
		if e := proto.Unmarshal(data, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *userWebClient) UpdateNickName(ctx context.Context, req *UpdateNickNameReq, header http.Header) (*UpdateNickNameResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathUserUpdateNickName, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(UpdateNickNameResp)
	if len(data) == 0 {
		return resp, nil
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-protobuf") {
		if e := proto.Unmarshal(data, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *userWebClient) UpdateEmail(ctx context.Context, req *UpdateEmailReq, header http.Header) (*UpdateEmailResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathUserUpdateEmail, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(UpdateEmailResp)
	if len(data) == 0 {
		return resp, nil
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-protobuf") {
		if e := proto.Unmarshal(data, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *userWebClient) UpdateTel(ctx context.Context, req *UpdateTelReq, header http.Header) (*UpdateTelResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathUserUpdateTel, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(UpdateTelResp)
	if len(data) == 0 {
		return resp, nil
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-protobuf") {
		if e := proto.Unmarshal(data, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}

type UserWebServer interface {
	Login(context.Context, *LoginReq) (*LoginResp, error)
	SelfUserInfo(context.Context, *SelfUserInfoReq) (*SelfUserInfoResp, error)
	UpdateStaticPassword(context.Context, *UpdateStaticPasswordReq) (*UpdateStaticPasswordResp, error)
	UpdateIdcard(context.Context, *UpdateIdcardReq) (*UpdateIdcardResp, error)
	UpdateNickName(context.Context, *UpdateNickNameReq) (*UpdateNickNameResp, error)
	UpdateEmail(context.Context, *UpdateEmailReq) (*UpdateEmailResp, error)
	UpdateTel(context.Context, *UpdateTelReq) (*UpdateTelResp, error)
}

func _User_Login_WebHandler(handler func(context.Context, *LoginReq) (*LoginResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(LoginReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/account.user/login]", map[string]interface{}{"error": e})
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/account.user/login]", map[string]interface{}{"error": e})
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/account.user/login]", map[string]interface{}{"error": e})
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/account.user/login]", map[string]interface{}{"error": e})
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			log.Error(ctx, "[/account.user/login]", map[string]interface{}{"error": "POST,PUT,PATCH only support application/json or application/x-protobuf"})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/login]", map[string]interface{}{"error": errstr})
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := cerror.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(LoginResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func _User_SelfUserInfo_WebHandler(handler func(context.Context, *SelfUserInfoReq) (*SelfUserInfoResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(SelfUserInfoReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/account.user/self_user_info]", map[string]interface{}{"error": e})
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/account.user/self_user_info]", map[string]interface{}{"error": e})
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/account.user/self_user_info]", map[string]interface{}{"error": e})
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/account.user/self_user_info]", map[string]interface{}{"error": e})
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			log.Error(ctx, "[/account.user/self_user_info]", map[string]interface{}{"error": "POST,PUT,PATCH only support application/json or application/x-protobuf"})
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := cerror.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(SelfUserInfoResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func _User_UpdateStaticPassword_WebHandler(handler func(context.Context, *UpdateStaticPasswordReq) (*UpdateStaticPasswordResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(UpdateStaticPasswordReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/account.user/update_static_password]", map[string]interface{}{"error": e})
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/account.user/update_static_password]", map[string]interface{}{"error": e})
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/account.user/update_static_password]", map[string]interface{}{"error": e})
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/account.user/update_static_password]", map[string]interface{}{"error": e})
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			log.Error(ctx, "[/account.user/update_static_password]", map[string]interface{}{"error": "POST,PUT,PATCH only support application/json or application/x-protobuf"})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/update_static_password]", map[string]interface{}{"error": errstr})
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := cerror.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(UpdateStaticPasswordResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func _User_UpdateIdcard_WebHandler(handler func(context.Context, *UpdateIdcardReq) (*UpdateIdcardResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(UpdateIdcardReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/account.user/update_idcard]", map[string]interface{}{"error": e})
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/account.user/update_idcard]", map[string]interface{}{"error": e})
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/account.user/update_idcard]", map[string]interface{}{"error": e})
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/account.user/update_idcard]", map[string]interface{}{"error": e})
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			log.Error(ctx, "[/account.user/update_idcard]", map[string]interface{}{"error": "POST,PUT,PATCH only support application/json or application/x-protobuf"})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/update_idcard]", map[string]interface{}{"error": errstr})
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := cerror.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(UpdateIdcardResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func _User_UpdateNickName_WebHandler(handler func(context.Context, *UpdateNickNameReq) (*UpdateNickNameResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(UpdateNickNameReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/account.user/update_nick_name]", map[string]interface{}{"error": e})
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/account.user/update_nick_name]", map[string]interface{}{"error": e})
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/account.user/update_nick_name]", map[string]interface{}{"error": e})
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/account.user/update_nick_name]", map[string]interface{}{"error": e})
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			log.Error(ctx, "[/account.user/update_nick_name]", map[string]interface{}{"error": "POST,PUT,PATCH only support application/json or application/x-protobuf"})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/update_nick_name]", map[string]interface{}{"error": errstr})
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := cerror.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(UpdateNickNameResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func _User_UpdateEmail_WebHandler(handler func(context.Context, *UpdateEmailReq) (*UpdateEmailResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(UpdateEmailReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/account.user/update_email]", map[string]interface{}{"error": e})
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/account.user/update_email]", map[string]interface{}{"error": e})
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/account.user/update_email]", map[string]interface{}{"error": e})
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/account.user/update_email]", map[string]interface{}{"error": e})
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			log.Error(ctx, "[/account.user/update_email]", map[string]interface{}{"error": "POST,PUT,PATCH only support application/json or application/x-protobuf"})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/update_email]", map[string]interface{}{"error": errstr})
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := cerror.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(UpdateEmailResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func _User_UpdateTel_WebHandler(handler func(context.Context, *UpdateTelReq) (*UpdateTelResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(UpdateTelReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/account.user/update_tel]", map[string]interface{}{"error": e})
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/account.user/update_tel]", map[string]interface{}{"error": e})
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/account.user/update_tel]", map[string]interface{}{"error": e})
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/account.user/update_tel]", map[string]interface{}{"error": e})
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			log.Error(ctx, "[/account.user/update_tel]", map[string]interface{}{"error": "POST,PUT,PATCH only support application/json or application/x-protobuf"})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/update_tel]", map[string]interface{}{"error": errstr})
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := cerror.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(UpdateTelResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func RegisterUserWebServer(engine *web.WebServer, svc UserWebServer, allmids map[string]web.OutsideHandler) {
	// avoid lint
	_ = allmids
	engine.Post(_WebPathUserLogin, _User_Login_WebHandler(svc.Login))
	{
		requiredMids := []string{"token"}
		mids := make([]web.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _User_SelfUserInfo_WebHandler(svc.SelfUserInfo))
		engine.Post(_WebPathUserSelfUserInfo, mids...)
	}
	{
		requiredMids := []string{"token"}
		mids := make([]web.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _User_UpdateStaticPassword_WebHandler(svc.UpdateStaticPassword))
		engine.Post(_WebPathUserUpdateStaticPassword, mids...)
	}
	{
		requiredMids := []string{"token"}
		mids := make([]web.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _User_UpdateIdcard_WebHandler(svc.UpdateIdcard))
		engine.Post(_WebPathUserUpdateIdcard, mids...)
	}
	{
		requiredMids := []string{"token"}
		mids := make([]web.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _User_UpdateNickName_WebHandler(svc.UpdateNickName))
		engine.Post(_WebPathUserUpdateNickName, mids...)
	}
	{
		requiredMids := []string{"token"}
		mids := make([]web.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _User_UpdateEmail_WebHandler(svc.UpdateEmail))
		engine.Post(_WebPathUserUpdateEmail, mids...)
	}
	{
		requiredMids := []string{"token"}
		mids := make([]web.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _User_UpdateTel_WebHandler(svc.UpdateTel))
		engine.Post(_WebPathUserUpdateTel, mids...)
	}
}
