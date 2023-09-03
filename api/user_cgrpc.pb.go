// Code generated by protoc-gen-go-cgrpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-cgrpc v0.0.78<br />
// 	protoc              v4.24.1<br />
// source: api/user.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	cgrpc "github.com/chenjie199234/Corelib/cgrpc"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
)

var _CGrpcPathUserGetUserInfo = "/account.user/get_user_info"
var _CGrpcPathUserLogin = "/account.user/login"
var _CGrpcPathUserSelfUserInfo = "/account.user/self_user_info"
var _CGrpcPathUserUpdateStaticPassword = "/account.user/update_static_password"
var _CGrpcPathUserUpdateIdcard = "/account.user/update_idcard"
var _CGrpcPathUserUpdateNickName = "/account.user/update_nick_name"
var _CGrpcPathUserUpdateEmail = "/account.user/update_email"
var _CGrpcPathUserUpdateTel = "/account.user/update_tel"

type UserCGrpcClient interface {
	GetUserInfo(context.Context, *GetUserInfoReq) (*GetUserInfoResp, error)
	Login(context.Context, *LoginReq) (*LoginResp, error)
	SelfUserInfo(context.Context, *SelfUserInfoReq) (*SelfUserInfoResp, error)
	UpdateStaticPassword(context.Context, *UpdateStaticPasswordReq) (*UpdateStaticPasswordResp, error)
	UpdateIdcard(context.Context, *UpdateIdcardReq) (*UpdateIdcardResp, error)
	UpdateNickName(context.Context, *UpdateNickNameReq) (*UpdateNickNameResp, error)
	UpdateEmail(context.Context, *UpdateEmailReq) (*UpdateEmailResp, error)
	UpdateTel(context.Context, *UpdateTelReq) (*UpdateTelResp, error)
}

type userCGrpcClient struct {
	cc *cgrpc.CGrpcClient
}

func NewUserCGrpcClient(c *cgrpc.CGrpcClient) UserCGrpcClient {
	return &userCGrpcClient{cc: c}
}

func (c *userCGrpcClient) GetUserInfo(ctx context.Context, req *GetUserInfoReq) (*GetUserInfoResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(GetUserInfoResp)
	if e := c.cc.Call(ctx, _CGrpcPathUserGetUserInfo, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) Login(ctx context.Context, req *LoginReq) (*LoginResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(LoginResp)
	if e := c.cc.Call(ctx, _CGrpcPathUserLogin, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) SelfUserInfo(ctx context.Context, req *SelfUserInfoReq) (*SelfUserInfoResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(SelfUserInfoResp)
	if e := c.cc.Call(ctx, _CGrpcPathUserSelfUserInfo, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) UpdateStaticPassword(ctx context.Context, req *UpdateStaticPasswordReq) (*UpdateStaticPasswordResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(UpdateStaticPasswordResp)
	if e := c.cc.Call(ctx, _CGrpcPathUserUpdateStaticPassword, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) UpdateIdcard(ctx context.Context, req *UpdateIdcardReq) (*UpdateIdcardResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(UpdateIdcardResp)
	if e := c.cc.Call(ctx, _CGrpcPathUserUpdateIdcard, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) UpdateNickName(ctx context.Context, req *UpdateNickNameReq) (*UpdateNickNameResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(UpdateNickNameResp)
	if e := c.cc.Call(ctx, _CGrpcPathUserUpdateNickName, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) UpdateEmail(ctx context.Context, req *UpdateEmailReq) (*UpdateEmailResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(UpdateEmailResp)
	if e := c.cc.Call(ctx, _CGrpcPathUserUpdateEmail, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) UpdateTel(ctx context.Context, req *UpdateTelReq) (*UpdateTelResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(UpdateTelResp)
	if e := c.cc.Call(ctx, _CGrpcPathUserUpdateTel, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}

type UserCGrpcServer interface {
	GetUserInfo(context.Context, *GetUserInfoReq) (*GetUserInfoResp, error)
	Login(context.Context, *LoginReq) (*LoginResp, error)
	SelfUserInfo(context.Context, *SelfUserInfoReq) (*SelfUserInfoResp, error)
	UpdateStaticPassword(context.Context, *UpdateStaticPasswordReq) (*UpdateStaticPasswordResp, error)
	UpdateIdcard(context.Context, *UpdateIdcardReq) (*UpdateIdcardResp, error)
	UpdateNickName(context.Context, *UpdateNickNameReq) (*UpdateNickNameResp, error)
	UpdateEmail(context.Context, *UpdateEmailReq) (*UpdateEmailResp, error)
	UpdateTel(context.Context, *UpdateTelReq) (*UpdateTelResp, error)
}

func _User_GetUserInfo_CGrpcHandler(handler func(context.Context, *GetUserInfoReq) (*GetUserInfoResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(GetUserInfoReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.user/get_user_info]", map[string]interface{}{"error": e})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/get_user_info]", map[string]interface{}{"error": errstr})
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(GetUserInfoResp)
		}
		ctx.Write(resp)
	}
}
func _User_Login_CGrpcHandler(handler func(context.Context, *LoginReq) (*LoginResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(LoginReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.user/login]", map[string]interface{}{"error": e})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/login]", map[string]interface{}{"error": errstr})
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(LoginResp)
		}
		ctx.Write(resp)
	}
}
func _User_SelfUserInfo_CGrpcHandler(handler func(context.Context, *SelfUserInfoReq) (*SelfUserInfoResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(SelfUserInfoReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.user/self_user_info]", map[string]interface{}{"error": e})
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(SelfUserInfoResp)
		}
		ctx.Write(resp)
	}
}
func _User_UpdateStaticPassword_CGrpcHandler(handler func(context.Context, *UpdateStaticPasswordReq) (*UpdateStaticPasswordResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(UpdateStaticPasswordReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.user/update_static_password]", map[string]interface{}{"error": e})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/update_static_password]", map[string]interface{}{"error": errstr})
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(UpdateStaticPasswordResp)
		}
		ctx.Write(resp)
	}
}
func _User_UpdateIdcard_CGrpcHandler(handler func(context.Context, *UpdateIdcardReq) (*UpdateIdcardResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(UpdateIdcardReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.user/update_idcard]", map[string]interface{}{"error": e})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/update_idcard]", map[string]interface{}{"error": errstr})
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(UpdateIdcardResp)
		}
		ctx.Write(resp)
	}
}
func _User_UpdateNickName_CGrpcHandler(handler func(context.Context, *UpdateNickNameReq) (*UpdateNickNameResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(UpdateNickNameReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.user/update_nick_name]", map[string]interface{}{"error": e})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/update_nick_name]", map[string]interface{}{"error": errstr})
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(UpdateNickNameResp)
		}
		ctx.Write(resp)
	}
}
func _User_UpdateEmail_CGrpcHandler(handler func(context.Context, *UpdateEmailReq) (*UpdateEmailResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(UpdateEmailReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.user/update_email]", map[string]interface{}{"error": e})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/update_email]", map[string]interface{}{"error": errstr})
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(UpdateEmailResp)
		}
		ctx.Write(resp)
	}
}
func _User_UpdateTel_CGrpcHandler(handler func(context.Context, *UpdateTelReq) (*UpdateTelResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(UpdateTelReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.user/update_tel]", map[string]interface{}{"error": e})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/update_tel]", map[string]interface{}{"error": errstr})
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(UpdateTelResp)
		}
		ctx.Write(resp)
	}
}
func RegisterUserCGrpcServer(engine *cgrpc.CGrpcServer, svc UserCGrpcServer, allmids map[string]cgrpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	engine.RegisterHandler("account.user", "get_user_info", _User_GetUserInfo_CGrpcHandler(svc.GetUserInfo))
	engine.RegisterHandler("account.user", "login", _User_Login_CGrpcHandler(svc.Login))
	engine.RegisterHandler("account.user", "self_user_info", _User_SelfUserInfo_CGrpcHandler(svc.SelfUserInfo))
	engine.RegisterHandler("account.user", "update_static_password", _User_UpdateStaticPassword_CGrpcHandler(svc.UpdateStaticPassword))
	engine.RegisterHandler("account.user", "update_idcard", _User_UpdateIdcard_CGrpcHandler(svc.UpdateIdcard))
	engine.RegisterHandler("account.user", "update_nick_name", _User_UpdateNickName_CGrpcHandler(svc.UpdateNickName))
	engine.RegisterHandler("account.user", "update_email", _User_UpdateEmail_CGrpcHandler(svc.UpdateEmail))
	engine.RegisterHandler("account.user", "update_tel", _User_UpdateTel_CGrpcHandler(svc.UpdateTel))
}
