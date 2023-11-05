// Code generated by protoc-gen-go-cgrpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-cgrpc v0.0.91<br />
// 	protoc              v4.24.4<br />
// source: api/account_user.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	cgrpc "github.com/chenjie199234/Corelib/cgrpc"
	log "github.com/chenjie199234/Corelib/log"
	grpc "google.golang.org/grpc"
)

var _CGrpcPathUserGetUserInfo = "/account.user/get_user_info"
var _CGrpcPathUserLogin = "/account.user/login"
var _CGrpcPathUserSelfUserInfo = "/account.user/self_user_info"
var _CGrpcPathUserUpdateStaticPassword = "/account.user/update_static_password"
var _CGrpcPathUserNickNameDuplicateCheck = "/account.user/nick_name_duplicate_check"
var _CGrpcPathUserUpdateNickName = "/account.user/update_nick_name"
var _CGrpcPathUserDelNickName = "/account.user/del_nick_name"
var _CGrpcPathUserIdcardDuplicateCheck = "/account.user/idcard_duplicate_check"
var _CGrpcPathUserUpdateIdcard = "/account.user/update_idcard"
var _CGrpcPathUserDelIdcard = "/account.user/del_idcard"
var _CGrpcPathUserEmailDuplicateCheck = "/account.user/email_duplicate_check"
var _CGrpcPathUserUpdateEmail = "/account.user/update_email"
var _CGrpcPathUserDelEmail = "/account.user/del_email"
var _CGrpcPathUserTelDuplicateCheck = "/account.user/tel_duplicate_check"
var _CGrpcPathUserUpdateTel = "/account.user/update_tel"
var _CGrpcPathUserDelTel = "/account.user/del_tel"

type UserCGrpcClient interface {
	GetUserInfo(context.Context, *GetUserInfoReq, ...grpc.CallOption) (*GetUserInfoResp, error)
	Login(context.Context, *LoginReq, ...grpc.CallOption) (*LoginResp, error)
	SelfUserInfo(context.Context, *SelfUserInfoReq, ...grpc.CallOption) (*SelfUserInfoResp, error)
	UpdateStaticPassword(context.Context, *UpdateStaticPasswordReq, ...grpc.CallOption) (*UpdateStaticPasswordResp, error)
	NickNameDuplicateCheck(context.Context, *NickNameDuplicateCheckReq, ...grpc.CallOption) (*NickNameDuplicateCheckResp, error)
	UpdateNickName(context.Context, *UpdateNickNameReq, ...grpc.CallOption) (*UpdateNickNameResp, error)
	DelNickName(context.Context, *DelNickNameReq, ...grpc.CallOption) (*DelNickNameResp, error)
	IdcardDuplicateCheck(context.Context, *IdcardDuplicateCheckReq, ...grpc.CallOption) (*IdcardDuplicateCheckResp, error)
	UpdateIdcard(context.Context, *UpdateIdcardReq, ...grpc.CallOption) (*UpdateIdcardResp, error)
	DelIdcard(context.Context, *DelIdcardReq, ...grpc.CallOption) (*DelIdcardResp, error)
	EmailDuplicateCheck(context.Context, *EmailDuplicateCheckReq, ...grpc.CallOption) (*EmailDuplicateCheckResp, error)
	UpdateEmail(context.Context, *UpdateEmailReq, ...grpc.CallOption) (*UpdateEmailResp, error)
	DelEmail(context.Context, *DelEmailReq, ...grpc.CallOption) (*DelEmailResp, error)
	TelDuplicateCheck(context.Context, *TelDuplicateCheckReq, ...grpc.CallOption) (*TelDuplicateCheckResp, error)
	UpdateTel(context.Context, *UpdateTelReq, ...grpc.CallOption) (*UpdateTelResp, error)
	DelTel(context.Context, *DelTelReq, ...grpc.CallOption) (*DelTelResp, error)
}

type userCGrpcClient struct {
	cc grpc.ClientConnInterface
}

func NewUserCGrpcClient(cc grpc.ClientConnInterface) UserCGrpcClient {
	return &userCGrpcClient{cc: cc}
}

func (c *userCGrpcClient) GetUserInfo(ctx context.Context, req *GetUserInfoReq, opts ...grpc.CallOption) (*GetUserInfoResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(GetUserInfoResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathUserGetUserInfo, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) Login(ctx context.Context, req *LoginReq, opts ...grpc.CallOption) (*LoginResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(LoginResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathUserLogin, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) SelfUserInfo(ctx context.Context, req *SelfUserInfoReq, opts ...grpc.CallOption) (*SelfUserInfoResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(SelfUserInfoResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathUserSelfUserInfo, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) UpdateStaticPassword(ctx context.Context, req *UpdateStaticPasswordReq, opts ...grpc.CallOption) (*UpdateStaticPasswordResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(UpdateStaticPasswordResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathUserUpdateStaticPassword, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) NickNameDuplicateCheck(ctx context.Context, req *NickNameDuplicateCheckReq, opts ...grpc.CallOption) (*NickNameDuplicateCheckResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(NickNameDuplicateCheckResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathUserNickNameDuplicateCheck, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) UpdateNickName(ctx context.Context, req *UpdateNickNameReq, opts ...grpc.CallOption) (*UpdateNickNameResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(UpdateNickNameResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathUserUpdateNickName, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) DelNickName(ctx context.Context, req *DelNickNameReq, opts ...grpc.CallOption) (*DelNickNameResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(DelNickNameResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathUserDelNickName, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) IdcardDuplicateCheck(ctx context.Context, req *IdcardDuplicateCheckReq, opts ...grpc.CallOption) (*IdcardDuplicateCheckResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(IdcardDuplicateCheckResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathUserIdcardDuplicateCheck, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) UpdateIdcard(ctx context.Context, req *UpdateIdcardReq, opts ...grpc.CallOption) (*UpdateIdcardResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(UpdateIdcardResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathUserUpdateIdcard, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) DelIdcard(ctx context.Context, req *DelIdcardReq, opts ...grpc.CallOption) (*DelIdcardResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(DelIdcardResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathUserDelIdcard, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) EmailDuplicateCheck(ctx context.Context, req *EmailDuplicateCheckReq, opts ...grpc.CallOption) (*EmailDuplicateCheckResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(EmailDuplicateCheckResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathUserEmailDuplicateCheck, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) UpdateEmail(ctx context.Context, req *UpdateEmailReq, opts ...grpc.CallOption) (*UpdateEmailResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(UpdateEmailResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathUserUpdateEmail, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) DelEmail(ctx context.Context, req *DelEmailReq, opts ...grpc.CallOption) (*DelEmailResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(DelEmailResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathUserDelEmail, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) TelDuplicateCheck(ctx context.Context, req *TelDuplicateCheckReq, opts ...grpc.CallOption) (*TelDuplicateCheckResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(TelDuplicateCheckResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathUserTelDuplicateCheck, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) UpdateTel(ctx context.Context, req *UpdateTelReq, opts ...grpc.CallOption) (*UpdateTelResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(UpdateTelResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathUserUpdateTel, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) DelTel(ctx context.Context, req *DelTelReq, opts ...grpc.CallOption) (*DelTelResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(DelTelResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathUserDelTel, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}

type UserCGrpcServer interface {
	GetUserInfo(context.Context, *GetUserInfoReq) (*GetUserInfoResp, error)
	Login(context.Context, *LoginReq) (*LoginResp, error)
	SelfUserInfo(context.Context, *SelfUserInfoReq) (*SelfUserInfoResp, error)
	UpdateStaticPassword(context.Context, *UpdateStaticPasswordReq) (*UpdateStaticPasswordResp, error)
	NickNameDuplicateCheck(context.Context, *NickNameDuplicateCheckReq) (*NickNameDuplicateCheckResp, error)
	UpdateNickName(context.Context, *UpdateNickNameReq) (*UpdateNickNameResp, error)
	DelNickName(context.Context, *DelNickNameReq) (*DelNickNameResp, error)
	IdcardDuplicateCheck(context.Context, *IdcardDuplicateCheckReq) (*IdcardDuplicateCheckResp, error)
	UpdateIdcard(context.Context, *UpdateIdcardReq) (*UpdateIdcardResp, error)
	DelIdcard(context.Context, *DelIdcardReq) (*DelIdcardResp, error)
	EmailDuplicateCheck(context.Context, *EmailDuplicateCheckReq) (*EmailDuplicateCheckResp, error)
	UpdateEmail(context.Context, *UpdateEmailReq) (*UpdateEmailResp, error)
	DelEmail(context.Context, *DelEmailReq) (*DelEmailResp, error)
	TelDuplicateCheck(context.Context, *TelDuplicateCheckReq) (*TelDuplicateCheckResp, error)
	UpdateTel(context.Context, *UpdateTelReq) (*UpdateTelResp, error)
	DelTel(context.Context, *DelTelReq) (*DelTelResp, error)
}

func _User_GetUserInfo_CGrpcHandler(handler func(context.Context, *GetUserInfoReq) (*GetUserInfoResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(GetUserInfoReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.user/get_user_info] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/get_user_info] validate failed", log.String("validate", errstr))
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
			log.Error(ctx, "[/account.user/login] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/login] validate failed", log.String("validate", errstr))
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
			log.Error(ctx, "[/account.user/self_user_info] decode failed")
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
			log.Error(ctx, "[/account.user/update_static_password] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/update_static_password] validate failed", log.String("validate", errstr))
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
func _User_NickNameDuplicateCheck_CGrpcHandler(handler func(context.Context, *NickNameDuplicateCheckReq) (*NickNameDuplicateCheckResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(NickNameDuplicateCheckReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.user/nick_name_duplicate_check] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/nick_name_duplicate_check] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(NickNameDuplicateCheckResp)
		}
		ctx.Write(resp)
	}
}
func _User_UpdateNickName_CGrpcHandler(handler func(context.Context, *UpdateNickNameReq) (*UpdateNickNameResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(UpdateNickNameReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.user/update_nick_name] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/update_nick_name] validate failed", log.String("validate", errstr))
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
func _User_DelNickName_CGrpcHandler(handler func(context.Context, *DelNickNameReq) (*DelNickNameResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(DelNickNameReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.user/del_nick_name] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/del_nick_name] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(DelNickNameResp)
		}
		ctx.Write(resp)
	}
}
func _User_IdcardDuplicateCheck_CGrpcHandler(handler func(context.Context, *IdcardDuplicateCheckReq) (*IdcardDuplicateCheckResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(IdcardDuplicateCheckReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.user/idcard_duplicate_check] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/idcard_duplicate_check] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(IdcardDuplicateCheckResp)
		}
		ctx.Write(resp)
	}
}
func _User_UpdateIdcard_CGrpcHandler(handler func(context.Context, *UpdateIdcardReq) (*UpdateIdcardResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(UpdateIdcardReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.user/update_idcard] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/update_idcard] validate failed", log.String("validate", errstr))
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
func _User_DelIdcard_CGrpcHandler(handler func(context.Context, *DelIdcardReq) (*DelIdcardResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(DelIdcardReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.user/del_idcard] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/del_idcard] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(DelIdcardResp)
		}
		ctx.Write(resp)
	}
}
func _User_EmailDuplicateCheck_CGrpcHandler(handler func(context.Context, *EmailDuplicateCheckReq) (*EmailDuplicateCheckResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(EmailDuplicateCheckReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.user/email_duplicate_check] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/email_duplicate_check] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(EmailDuplicateCheckResp)
		}
		ctx.Write(resp)
	}
}
func _User_UpdateEmail_CGrpcHandler(handler func(context.Context, *UpdateEmailReq) (*UpdateEmailResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(UpdateEmailReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.user/update_email] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/update_email] validate failed", log.String("validate", errstr))
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
func _User_DelEmail_CGrpcHandler(handler func(context.Context, *DelEmailReq) (*DelEmailResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(DelEmailReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.user/del_email] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/del_email] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(DelEmailResp)
		}
		ctx.Write(resp)
	}
}
func _User_TelDuplicateCheck_CGrpcHandler(handler func(context.Context, *TelDuplicateCheckReq) (*TelDuplicateCheckResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(TelDuplicateCheckReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.user/tel_duplicate_check] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/tel_duplicate_check] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(TelDuplicateCheckResp)
		}
		ctx.Write(resp)
	}
}
func _User_UpdateTel_CGrpcHandler(handler func(context.Context, *UpdateTelReq) (*UpdateTelResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(UpdateTelReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.user/update_tel] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/update_tel] validate failed", log.String("validate", errstr))
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
func _User_DelTel_CGrpcHandler(handler func(context.Context, *DelTelReq) (*DelTelResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(DelTelReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.user/del_tel] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.user/del_tel] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(DelTelResp)
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
	engine.RegisterHandler("account.user", "nick_name_duplicate_check", _User_NickNameDuplicateCheck_CGrpcHandler(svc.NickNameDuplicateCheck))
	engine.RegisterHandler("account.user", "update_nick_name", _User_UpdateNickName_CGrpcHandler(svc.UpdateNickName))
	engine.RegisterHandler("account.user", "del_nick_name", _User_DelNickName_CGrpcHandler(svc.DelNickName))
	engine.RegisterHandler("account.user", "idcard_duplicate_check", _User_IdcardDuplicateCheck_CGrpcHandler(svc.IdcardDuplicateCheck))
	engine.RegisterHandler("account.user", "update_idcard", _User_UpdateIdcard_CGrpcHandler(svc.UpdateIdcard))
	engine.RegisterHandler("account.user", "del_idcard", _User_DelIdcard_CGrpcHandler(svc.DelIdcard))
	engine.RegisterHandler("account.user", "email_duplicate_check", _User_EmailDuplicateCheck_CGrpcHandler(svc.EmailDuplicateCheck))
	engine.RegisterHandler("account.user", "update_email", _User_UpdateEmail_CGrpcHandler(svc.UpdateEmail))
	engine.RegisterHandler("account.user", "del_email", _User_DelEmail_CGrpcHandler(svc.DelEmail))
	engine.RegisterHandler("account.user", "tel_duplicate_check", _User_TelDuplicateCheck_CGrpcHandler(svc.TelDuplicateCheck))
	engine.RegisterHandler("account.user", "update_tel", _User_UpdateTel_CGrpcHandler(svc.UpdateTel))
	engine.RegisterHandler("account.user", "del_tel", _User_DelTel_CGrpcHandler(svc.DelTel))
}
