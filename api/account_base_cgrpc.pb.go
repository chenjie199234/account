// Code generated by protoc-gen-go-cgrpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-cgrpc v0.0.124<br />
// 	protoc              v5.28.0<br />
// source: api/account_base.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	cgrpc "github.com/chenjie199234/Corelib/cgrpc"
	grpc "google.golang.org/grpc"
	slog "log/slog"
)

var _CGrpcPathBaseBaseInfo = "/account.base/base_info"
var _CGrpcPathBaseBan = "/account.base/ban"
var _CGrpcPathBaseUnban = "/account.base/unban"

type BaseCGrpcClient interface {
	// if the request if from web,only can get self's info,the src_type and src in request will be ignored,the user_id in token will be used
	BaseInfo(context.Context, *BaseInfoReq) (*BaseInfoResp, error)
	Ban(context.Context, *BanReq) (*BanResp, error)
	Unban(context.Context, *UnbanReq) (*UnbanResp, error)
}

type baseCGrpcClient struct {
	cc grpc.ClientConnInterface
}

func NewBaseCGrpcClient(cc grpc.ClientConnInterface) BaseCGrpcClient {
	return &baseCGrpcClient{cc: cc}
}

func (c *baseCGrpcClient) BaseInfo(ctx context.Context, req *BaseInfoReq) (*BaseInfoResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if errstr := req.Validate(); errstr != "" {
		slog.ErrorContext(ctx, "[/account.base/base_info] validate failed", slog.String("error", errstr))
		return nil, cerror.ErrReq
	}
	resp := new(BaseInfoResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathBaseBaseInfo, req, resp); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *baseCGrpcClient) Ban(ctx context.Context, req *BanReq) (*BanResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if errstr := req.Validate(); errstr != "" {
		slog.ErrorContext(ctx, "[/account.base/ban] validate failed", slog.String("error", errstr))
		return nil, cerror.ErrReq
	}
	resp := new(BanResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathBaseBan, req, resp); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *baseCGrpcClient) Unban(ctx context.Context, req *UnbanReq) (*UnbanResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if errstr := req.Validate(); errstr != "" {
		slog.ErrorContext(ctx, "[/account.base/unban] validate failed", slog.String("error", errstr))
		return nil, cerror.ErrReq
	}
	resp := new(UnbanResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathBaseUnban, req, resp); e != nil {
		return nil, e
	}
	return resp, nil
}

type BaseCGrpcServer interface {
	// if the request if from web,only can get self's info,the src_type and src in request will be ignored,the user_id in token will be used
	// Context is *cgrpc.NoStreamServerContext
	BaseInfo(context.Context, *BaseInfoReq) (*BaseInfoResp, error)
	// Context is *cgrpc.NoStreamServerContext
	Ban(context.Context, *BanReq) (*BanResp, error)
	// Context is *cgrpc.NoStreamServerContext
	Unban(context.Context, *UnbanReq) (*UnbanResp, error)
}

func _Base_BaseInfo_CGrpcHandler(handler func(context.Context, *BaseInfoReq) (*BaseInfoResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.ServerContext) {
		req := new(BaseInfoReq)
		if e := ctx.Read(req); e != nil {
			slog.ErrorContext(ctx, "[/account.base/base_info] decode failed", slog.String("error", e.Error()))
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			slog.ErrorContext(ctx, "[/account.base/base_info] validate failed", slog.String("error", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(cgrpc.NewNoStreamServerContext(ctx), req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(BaseInfoResp)
		}
		ctx.Write(resp)
	}
}
func _Base_Ban_CGrpcHandler(handler func(context.Context, *BanReq) (*BanResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.ServerContext) {
		req := new(BanReq)
		if e := ctx.Read(req); e != nil {
			slog.ErrorContext(ctx, "[/account.base/ban] decode failed", slog.String("error", e.Error()))
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			slog.ErrorContext(ctx, "[/account.base/ban] validate failed", slog.String("error", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(cgrpc.NewNoStreamServerContext(ctx), req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(BanResp)
		}
		ctx.Write(resp)
	}
}
func _Base_Unban_CGrpcHandler(handler func(context.Context, *UnbanReq) (*UnbanResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.ServerContext) {
		req := new(UnbanReq)
		if e := ctx.Read(req); e != nil {
			slog.ErrorContext(ctx, "[/account.base/unban] decode failed", slog.String("error", e.Error()))
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			slog.ErrorContext(ctx, "[/account.base/unban] validate failed", slog.String("error", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(cgrpc.NewNoStreamServerContext(ctx), req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(UnbanResp)
		}
		ctx.Write(resp)
	}
}
func RegisterBaseCGrpcServer(engine *cgrpc.CGrpcServer, svc BaseCGrpcServer, allmids map[string]cgrpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	{
		requiredMids := []string{"accesskey"}
		mids := make([]cgrpc.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _Base_BaseInfo_CGrpcHandler(svc.BaseInfo))
		engine.RegisterHandler("account.base", "base_info", false, false, mids...)
	}
	{
		requiredMids := []string{"accesskey"}
		mids := make([]cgrpc.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _Base_Ban_CGrpcHandler(svc.Ban))
		engine.RegisterHandler("account.base", "ban", false, false, mids...)
	}
	{
		requiredMids := []string{"accesskey"}
		mids := make([]cgrpc.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _Base_Unban_CGrpcHandler(svc.Unban))
		engine.RegisterHandler("account.base", "unban", false, false, mids...)
	}
}
