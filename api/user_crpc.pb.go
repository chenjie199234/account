// Code generated by protoc-gen-go-crpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-crpc v0.0.78<br />
// 	protoc             v4.24.1<br />
// source: api/user.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	crpc "github.com/chenjie199234/Corelib/crpc"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
)

var _CrpcPathUserGetUserInfo = "/account.user/get_user_info"
var _CrpcPathUserLogin = "/account.user/login"
var _CrpcPathUserSelfUserInfo = "/account.user/self_user_info"
var _CrpcPathUserUpdateStaticPassword = "/account.user/update_static_password"
var _CrpcPathUserUpdateIdcard = "/account.user/update_idcard"
var _CrpcPathUserUpdateNickName = "/account.user/update_nick_name"
var _CrpcPathUserUpdateEmail = "/account.user/update_email"
var _CrpcPathUserUpdateTel = "/account.user/update_tel"

type UserCrpcClient interface {
	GetUserInfo(context.Context, *GetUserInfoReq) (*GetUserInfoResp, error)
	Login(context.Context, *LoginReq) (*LoginResp, error)
	SelfUserInfo(context.Context, *SelfUserInfoReq) (*SelfUserInfoResp, error)
	UpdateStaticPassword(context.Context, *UpdateStaticPasswordReq) (*UpdateStaticPasswordResp, error)
	UpdateIdcard(context.Context, *UpdateIdcardReq) (*UpdateIdcardResp, error)
	UpdateNickName(context.Context, *UpdateNickNameReq) (*UpdateNickNameResp, error)
	UpdateEmail(context.Context, *UpdateEmailReq) (*UpdateEmailResp, error)
	UpdateTel(context.Context, *UpdateTelReq) (*UpdateTelResp, error)
}

type userCrpcClient struct {
	cc *crpc.CrpcClient
}

func NewUserCrpcClient(c *crpc.CrpcClient) UserCrpcClient {
	return &userCrpcClient{cc: c}
}

func (c *userCrpcClient) GetUserInfo(ctx context.Context, req *GetUserInfoReq) (*GetUserInfoResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathUserGetUserInfo, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(GetUserInfoResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if len(respd) >= 2 && respd[0] == '{' && respd[len(respd)-1] == '}' {
		if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(respd, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *userCrpcClient) Login(ctx context.Context, req *LoginReq) (*LoginResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathUserLogin, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(LoginResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if len(respd) >= 2 && respd[0] == '{' && respd[len(respd)-1] == '}' {
		if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(respd, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *userCrpcClient) SelfUserInfo(ctx context.Context, req *SelfUserInfoReq) (*SelfUserInfoResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathUserSelfUserInfo, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(SelfUserInfoResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if len(respd) >= 2 && respd[0] == '{' && respd[len(respd)-1] == '}' {
		if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(respd, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *userCrpcClient) UpdateStaticPassword(ctx context.Context, req *UpdateStaticPasswordReq) (*UpdateStaticPasswordResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathUserUpdateStaticPassword, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(UpdateStaticPasswordResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if len(respd) >= 2 && respd[0] == '{' && respd[len(respd)-1] == '}' {
		if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(respd, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *userCrpcClient) UpdateIdcard(ctx context.Context, req *UpdateIdcardReq) (*UpdateIdcardResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathUserUpdateIdcard, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(UpdateIdcardResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if len(respd) >= 2 && respd[0] == '{' && respd[len(respd)-1] == '}' {
		if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(respd, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *userCrpcClient) UpdateNickName(ctx context.Context, req *UpdateNickNameReq) (*UpdateNickNameResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathUserUpdateNickName, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(UpdateNickNameResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if len(respd) >= 2 && respd[0] == '{' && respd[len(respd)-1] == '}' {
		if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(respd, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *userCrpcClient) UpdateEmail(ctx context.Context, req *UpdateEmailReq) (*UpdateEmailResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathUserUpdateEmail, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(UpdateEmailResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if len(respd) >= 2 && respd[0] == '{' && respd[len(respd)-1] == '}' {
		if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(respd, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *userCrpcClient) UpdateTel(ctx context.Context, req *UpdateTelReq) (*UpdateTelResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathUserUpdateTel, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(UpdateTelResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if len(respd) >= 2 && respd[0] == '{' && respd[len(respd)-1] == '}' {
		if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(respd, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}

type UserCrpcServer interface {
	GetUserInfo(context.Context, *GetUserInfoReq) (*GetUserInfoResp, error)
	Login(context.Context, *LoginReq) (*LoginResp, error)
	SelfUserInfo(context.Context, *SelfUserInfoReq) (*SelfUserInfoResp, error)
	UpdateStaticPassword(context.Context, *UpdateStaticPasswordReq) (*UpdateStaticPasswordResp, error)
	UpdateIdcard(context.Context, *UpdateIdcardReq) (*UpdateIdcardResp, error)
	UpdateNickName(context.Context, *UpdateNickNameReq) (*UpdateNickNameResp, error)
	UpdateEmail(context.Context, *UpdateEmailReq) (*UpdateEmailResp, error)
	UpdateTel(context.Context, *UpdateTelReq) (*UpdateTelResp, error)
}

func _User_GetUserInfo_CrpcHandler(handler func(context.Context, *GetUserInfoReq) (*GetUserInfoResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(GetUserInfoReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/account.user/get_user_info] json and proto format decode both failed", nil)
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/account.user/get_user_info] json and proto format decode both failed", nil)
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
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
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func _User_Login_CrpcHandler(handler func(context.Context, *LoginReq) (*LoginResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(LoginReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/account.user/login] json and proto format decode both failed", nil)
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/account.user/login] json and proto format decode both failed", nil)
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
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
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func _User_SelfUserInfo_CrpcHandler(handler func(context.Context, *SelfUserInfoReq) (*SelfUserInfoResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(SelfUserInfoReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/account.user/self_user_info] json and proto format decode both failed", nil)
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/account.user/self_user_info] json and proto format decode both failed", nil)
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(SelfUserInfoResp)
		}
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func _User_UpdateStaticPassword_CrpcHandler(handler func(context.Context, *UpdateStaticPasswordReq) (*UpdateStaticPasswordResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(UpdateStaticPasswordReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/account.user/update_static_password] json and proto format decode both failed", nil)
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/account.user/update_static_password] json and proto format decode both failed", nil)
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
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
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func _User_UpdateIdcard_CrpcHandler(handler func(context.Context, *UpdateIdcardReq) (*UpdateIdcardResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(UpdateIdcardReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/account.user/update_idcard] json and proto format decode both failed", nil)
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/account.user/update_idcard] json and proto format decode both failed", nil)
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
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
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func _User_UpdateNickName_CrpcHandler(handler func(context.Context, *UpdateNickNameReq) (*UpdateNickNameResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(UpdateNickNameReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/account.user/update_nick_name] json and proto format decode both failed", nil)
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/account.user/update_nick_name] json and proto format decode both failed", nil)
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
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
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func _User_UpdateEmail_CrpcHandler(handler func(context.Context, *UpdateEmailReq) (*UpdateEmailResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(UpdateEmailReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/account.user/update_email] json and proto format decode both failed", nil)
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/account.user/update_email] json and proto format decode both failed", nil)
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
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
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func _User_UpdateTel_CrpcHandler(handler func(context.Context, *UpdateTelReq) (*UpdateTelResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(UpdateTelReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/account.user/update_tel] json and proto format decode both failed", nil)
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/account.user/update_tel] json and proto format decode both failed", nil)
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
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
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func RegisterUserCrpcServer(engine *crpc.CrpcServer, svc UserCrpcServer, allmids map[string]crpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	engine.RegisterHandler("account.user", "get_user_info", _User_GetUserInfo_CrpcHandler(svc.GetUserInfo))
	engine.RegisterHandler("account.user", "login", _User_Login_CrpcHandler(svc.Login))
	engine.RegisterHandler("account.user", "self_user_info", _User_SelfUserInfo_CrpcHandler(svc.SelfUserInfo))
	engine.RegisterHandler("account.user", "update_static_password", _User_UpdateStaticPassword_CrpcHandler(svc.UpdateStaticPassword))
	engine.RegisterHandler("account.user", "update_idcard", _User_UpdateIdcard_CrpcHandler(svc.UpdateIdcard))
	engine.RegisterHandler("account.user", "update_nick_name", _User_UpdateNickName_CrpcHandler(svc.UpdateNickName))
	engine.RegisterHandler("account.user", "update_email", _User_UpdateEmail_CrpcHandler(svc.UpdateEmail))
	engine.RegisterHandler("account.user", "update_tel", _User_UpdateTel_CrpcHandler(svc.UpdateTel))
}
