// Code generated by protoc-gen-go-crpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-crpc v0.0.101<br />
// 	protoc             v4.25.3<br />
// source: api/account_user.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	crpc "github.com/chenjie199234/Corelib/crpc"
	log "github.com/chenjie199234/Corelib/log"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
)

var _CrpcPathUserGetUserInfo = "/account.user/get_user_info"

type UserCrpcClient interface {
	GetUserInfo(context.Context, *GetUserInfoReq) (*GetUserInfoResp, error)
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
	respd, e := c.cc.Call(ctx, _CrpcPathUserGetUserInfo, reqd)
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

type UserCrpcServer interface {
	GetUserInfo(context.Context, *GetUserInfoReq) (*GetUserInfoResp, error)
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
					log.Error(ctx, "[/account.user/get_user_info] json and proto format decode both failed")
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/account.user/get_user_info] json and proto format decode both failed")
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
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
}
