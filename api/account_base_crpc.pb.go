// Code generated by protoc-gen-go-crpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-crpc v0.0.135<br />
// 	protoc             v6.30.2<br />
// source: api/account_base.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	crpc "github.com/chenjie199234/Corelib/crpc"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	slog "log/slog"
)

var _CrpcPathBaseBaseInfo = "/account.base/base_info"
var _CrpcPathBaseBan = "/account.base/ban"
var _CrpcPathBaseUnban = "/account.base/unban"

type BaseCrpcClient interface {
	BaseInfo(ctx context.Context, req *BaseInfoReq) (resp *BaseInfoResp, e error)
	Ban(ctx context.Context, req *BanReq) (resp *BanResp, e error)
	Unban(ctx context.Context, req *UnbanReq) (resp *UnbanResp, e error)
}

type baseCrpcClient struct {
	cc *crpc.CrpcClient
}

func NewBaseCrpcClient(c *crpc.CrpcClient) BaseCrpcClient {
	return &baseCrpcClient{cc: c}
}

func (c *baseCrpcClient) BaseInfo(ctx context.Context, req *BaseInfoReq) (*BaseInfoResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if errstr := req.Validate(); errstr != "" {
		slog.ErrorContext(ctx, "[/account.base/base_info] request validate failed", slog.String("error", errstr))
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	var respbody []byte
	var encoder crpc.Encoder
	if e := c.cc.Call(ctx, _CrpcPathBaseBaseInfo, reqd, crpc.Encoder_Protobuf, func(ctx *crpc.CallContext) error {
		var e error
		if respbody, encoder, e = ctx.Recv(); e != nil {
			slog.ErrorContext(ctx, "[/account.base/base_info] read response failed", slog.String("error", e.Error()))
		}
		return e
	}); e != nil {
		return nil, e
	}
	resp := new(BaseInfoResp)
	if len(respbody) == 0 {
		return resp, nil
	}
	switch encoder {
	case crpc.Encoder_Protobuf:
		if e := proto.Unmarshal(respbody, resp); e != nil {
			return nil, cerror.ErrResp
		}
	case crpc.Encoder_Json:
		if e := protojson.Unmarshal(respbody, resp); e != nil {
			return nil, cerror.ErrResp
		}
	default:
		slog.ErrorContext(ctx, "[/account.base/base_info] unknown response encoder")
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *baseCrpcClient) Ban(ctx context.Context, req *BanReq) (*BanResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if errstr := req.Validate(); errstr != "" {
		slog.ErrorContext(ctx, "[/account.base/ban] request validate failed", slog.String("error", errstr))
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	var respbody []byte
	var encoder crpc.Encoder
	if e := c.cc.Call(ctx, _CrpcPathBaseBan, reqd, crpc.Encoder_Protobuf, func(ctx *crpc.CallContext) error {
		var e error
		if respbody, encoder, e = ctx.Recv(); e != nil {
			slog.ErrorContext(ctx, "[/account.base/ban] read response failed", slog.String("error", e.Error()))
		}
		return e
	}); e != nil {
		return nil, e
	}
	resp := new(BanResp)
	if len(respbody) == 0 {
		return resp, nil
	}
	switch encoder {
	case crpc.Encoder_Protobuf:
		if e := proto.Unmarshal(respbody, resp); e != nil {
			return nil, cerror.ErrResp
		}
	case crpc.Encoder_Json:
		if e := protojson.Unmarshal(respbody, resp); e != nil {
			return nil, cerror.ErrResp
		}
	default:
		slog.ErrorContext(ctx, "[/account.base/ban] unknown response encoder")
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *baseCrpcClient) Unban(ctx context.Context, req *UnbanReq) (*UnbanResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if errstr := req.Validate(); errstr != "" {
		slog.ErrorContext(ctx, "[/account.base/unban] request validate failed", slog.String("error", errstr))
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	var respbody []byte
	var encoder crpc.Encoder
	if e := c.cc.Call(ctx, _CrpcPathBaseUnban, reqd, crpc.Encoder_Protobuf, func(ctx *crpc.CallContext) error {
		var e error
		if respbody, encoder, e = ctx.Recv(); e != nil {
			slog.ErrorContext(ctx, "[/account.base/unban] read response failed", slog.String("error", e.Error()))
		}
		return e
	}); e != nil {
		return nil, e
	}
	resp := new(UnbanResp)
	if len(respbody) == 0 {
		return resp, nil
	}
	switch encoder {
	case crpc.Encoder_Protobuf:
		if e := proto.Unmarshal(respbody, resp); e != nil {
			return nil, cerror.ErrResp
		}
	case crpc.Encoder_Json:
		if e := protojson.Unmarshal(respbody, resp); e != nil {
			return nil, cerror.ErrResp
		}
	default:
		slog.ErrorContext(ctx, "[/account.base/unban] unknown response encoder")
		return nil, cerror.ErrResp
	}
	return resp, nil
}

type BaseCrpcServer interface {
	// Context is *crpc.NoStreamServerContext
	BaseInfo(context.Context, *BaseInfoReq) (*BaseInfoResp, error)
	// Context is *crpc.NoStreamServerContext
	Ban(context.Context, *BanReq) (*BanResp, error)
	// Context is *crpc.NoStreamServerContext
	Unban(context.Context, *UnbanReq) (*UnbanResp, error)
}

func _Base_BaseInfo_CrpcHandler(handler func(context.Context, *BaseInfoReq) (*BaseInfoResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.ServerContext) {
		reqbody, encoder, e := ctx.Recv()
		if e != nil {
			slog.ErrorContext(ctx, "[/account.base/base_info] read request failed", slog.String("error", e.Error()))
			ctx.Abort(e)
			return
		}
		req := new(BaseInfoReq)
		switch encoder {
		case crpc.Encoder_Protobuf:
			if e := proto.Unmarshal(reqbody, req); e != nil {
				slog.ErrorContext(ctx, "[/account.base/base_info] request decode failed", slog.String("error", e.Error()))
				ctx.Abort(cerror.ErrReq)
				return
			}
		case crpc.Encoder_Json:
			if e := protojson.Unmarshal(reqbody, req); e != nil {
				slog.ErrorContext(ctx, "[/account.base/base_info] request decode failed", slog.String("error", e.Error()))
				ctx.Abort(cerror.ErrReq)
				return
			}
		default:
			slog.ErrorContext(ctx, "[/account.base/base_info] request encoder unknown")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			slog.ErrorContext(ctx, "[/account.base/base_info] request validate failed", slog.String("error", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(crpc.NewNoStreamServerContext(ctx), req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(BaseInfoResp)
		}
		var respd []byte
		switch encoder {
		case crpc.Encoder_Protobuf:
			respd, _ = proto.Marshal(resp)
		case crpc.Encoder_Json:
			respd, _ = protojson.Marshal(resp)
		}
		if e := ctx.Send(respd, encoder); e != nil {
			slog.ErrorContext(ctx, "[/account.base/base_info] send response failed", slog.String("error", e.Error()))
		}
	}
}
func _Base_Ban_CrpcHandler(handler func(context.Context, *BanReq) (*BanResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.ServerContext) {
		reqbody, encoder, e := ctx.Recv()
		if e != nil {
			slog.ErrorContext(ctx, "[/account.base/ban] read request failed", slog.String("error", e.Error()))
			ctx.Abort(e)
			return
		}
		req := new(BanReq)
		switch encoder {
		case crpc.Encoder_Protobuf:
			if e := proto.Unmarshal(reqbody, req); e != nil {
				slog.ErrorContext(ctx, "[/account.base/ban] request decode failed", slog.String("error", e.Error()))
				ctx.Abort(cerror.ErrReq)
				return
			}
		case crpc.Encoder_Json:
			if e := protojson.Unmarshal(reqbody, req); e != nil {
				slog.ErrorContext(ctx, "[/account.base/ban] request decode failed", slog.String("error", e.Error()))
				ctx.Abort(cerror.ErrReq)
				return
			}
		default:
			slog.ErrorContext(ctx, "[/account.base/ban] request encoder unknown")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			slog.ErrorContext(ctx, "[/account.base/ban] request validate failed", slog.String("error", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(crpc.NewNoStreamServerContext(ctx), req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(BanResp)
		}
		var respd []byte
		switch encoder {
		case crpc.Encoder_Protobuf:
			respd, _ = proto.Marshal(resp)
		case crpc.Encoder_Json:
			respd, _ = protojson.Marshal(resp)
		}
		if e := ctx.Send(respd, encoder); e != nil {
			slog.ErrorContext(ctx, "[/account.base/ban] send response failed", slog.String("error", e.Error()))
		}
	}
}
func _Base_Unban_CrpcHandler(handler func(context.Context, *UnbanReq) (*UnbanResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.ServerContext) {
		reqbody, encoder, e := ctx.Recv()
		if e != nil {
			slog.ErrorContext(ctx, "[/account.base/unban] read request failed", slog.String("error", e.Error()))
			ctx.Abort(e)
			return
		}
		req := new(UnbanReq)
		switch encoder {
		case crpc.Encoder_Protobuf:
			if e := proto.Unmarshal(reqbody, req); e != nil {
				slog.ErrorContext(ctx, "[/account.base/unban] request decode failed", slog.String("error", e.Error()))
				ctx.Abort(cerror.ErrReq)
				return
			}
		case crpc.Encoder_Json:
			if e := protojson.Unmarshal(reqbody, req); e != nil {
				slog.ErrorContext(ctx, "[/account.base/unban] request decode failed", slog.String("error", e.Error()))
				ctx.Abort(cerror.ErrReq)
				return
			}
		default:
			slog.ErrorContext(ctx, "[/account.base/unban] request encoder unknown")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			slog.ErrorContext(ctx, "[/account.base/unban] request validate failed", slog.String("error", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(crpc.NewNoStreamServerContext(ctx), req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(UnbanResp)
		}
		var respd []byte
		switch encoder {
		case crpc.Encoder_Protobuf:
			respd, _ = proto.Marshal(resp)
		case crpc.Encoder_Json:
			respd, _ = protojson.Marshal(resp)
		}
		if e := ctx.Send(respd, encoder); e != nil {
			slog.ErrorContext(ctx, "[/account.base/unban] send response failed", slog.String("error", e.Error()))
		}
	}
}
func RegisterBaseCrpcServer(engine *crpc.CrpcServer, svc BaseCrpcServer, allmids map[string]crpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	{
		requiredMids := []string{"accesskey"}
		mids := make([]crpc.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _Base_BaseInfo_CrpcHandler(svc.BaseInfo))
		engine.RegisterHandler("account.base", "base_info", mids...)
	}
	{
		requiredMids := []string{"accesskey"}
		mids := make([]crpc.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _Base_Ban_CrpcHandler(svc.Ban))
		engine.RegisterHandler("account.base", "ban", mids...)
	}
	{
		requiredMids := []string{"accesskey"}
		mids := make([]crpc.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _Base_Unban_CrpcHandler(svc.Unban))
		engine.RegisterHandler("account.base", "unban", mids...)
	}
}
