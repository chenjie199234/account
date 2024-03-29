// Code generated by protoc-gen-go-cgrpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-cgrpc v0.0.110<br />
// 	protoc              v4.25.3<br />
// source: api/account_base.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	cgrpc "github.com/chenjie199234/Corelib/cgrpc"
	log "github.com/chenjie199234/Corelib/log"
	grpc "google.golang.org/grpc"
)

var _CGrpcPathBaseGetBaseInfo = "/account.base/get_base_info"

type BaseCGrpcClient interface {
	GetBaseInfo(context.Context, *GetBaseInfoReq, ...grpc.CallOption) (*GetBaseInfoResp, error)
}

type baseCGrpcClient struct {
	cc grpc.ClientConnInterface
}

func NewBaseCGrpcClient(cc grpc.ClientConnInterface) BaseCGrpcClient {
	return &baseCGrpcClient{cc: cc}
}

func (c *baseCGrpcClient) GetBaseInfo(ctx context.Context, req *GetBaseInfoReq, opts ...grpc.CallOption) (*GetBaseInfoResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(GetBaseInfoResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathBaseGetBaseInfo, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}

type BaseCGrpcServer interface {
	GetBaseInfo(context.Context, *GetBaseInfoReq) (*GetBaseInfoResp, error)
}

func _Base_GetBaseInfo_CGrpcHandler(handler func(context.Context, *GetBaseInfoReq) (*GetBaseInfoResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(GetBaseInfoReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.base/get_base_info] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.base/get_base_info] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(GetBaseInfoResp)
		}
		ctx.Write(resp)
	}
}
func RegisterBaseCGrpcServer(engine *cgrpc.CGrpcServer, svc BaseCGrpcServer, allmids map[string]cgrpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	engine.RegisterHandler("account.base", "get_base_info", _Base_GetBaseInfo_CGrpcHandler(svc.GetBaseInfo))
}
