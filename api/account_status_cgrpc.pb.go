// Code generated by protoc-gen-go-cgrpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-cgrpc v0.0.81<br />
// 	protoc              v4.24.1<br />
// source: api/account_status.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	cgrpc "github.com/chenjie199234/Corelib/cgrpc"
	log "github.com/chenjie199234/Corelib/log"
	grpc "google.golang.org/grpc"
)

var _CGrpcPathStatusPing = "/account.status/ping"

type StatusCGrpcClient interface {
	// ping check server's health
	Ping(context.Context, *Pingreq, ...grpc.CallOption) (*Pingresp, error)
}

type statusCGrpcClient struct {
	cc grpc.ClientConnInterface
}

func NewStatusCGrpcClient(cc grpc.ClientConnInterface) StatusCGrpcClient {
	return &statusCGrpcClient{cc: cc}
}

func (c *statusCGrpcClient) Ping(ctx context.Context, req *Pingreq, opts ...grpc.CallOption) (*Pingresp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(Pingresp)
	if e := c.cc.Invoke(ctx, _CGrpcPathStatusPing, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}

type StatusCGrpcServer interface {
	// ping check server's health
	Ping(context.Context, *Pingreq) (*Pingresp, error)
}

func _Status_Ping_CGrpcHandler(handler func(context.Context, *Pingreq) (*Pingresp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(Pingreq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.status/ping] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.status/ping] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(Pingresp)
		}
		ctx.Write(resp)
	}
}
func RegisterStatusCGrpcServer(engine *cgrpc.CGrpcServer, svc StatusCGrpcServer, allmids map[string]cgrpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	engine.RegisterHandler("account.status", "ping", _Status_Ping_CGrpcHandler(svc.Ping))
}
