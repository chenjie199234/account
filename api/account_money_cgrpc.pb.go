// Code generated by protoc-gen-go-cgrpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-cgrpc v0.0.117<br />
// 	protoc              v5.27.2<br />
// source: api/account_money.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	cgrpc "github.com/chenjie199234/Corelib/cgrpc"
	log "github.com/chenjie199234/Corelib/log"
	grpc "google.golang.org/grpc"
)

var _CGrpcPathMoneyGetMoneyLogs = "/account.money/get_money_logs"

type MoneyCGrpcClient interface {
	GetMoneyLogs(context.Context, *GetMoneyLogsReq, ...grpc.CallOption) (*GetMoneyLogsResp, error)
}

type moneyCGrpcClient struct {
	cc grpc.ClientConnInterface
}

func NewMoneyCGrpcClient(cc grpc.ClientConnInterface) MoneyCGrpcClient {
	return &moneyCGrpcClient{cc: cc}
}

func (c *moneyCGrpcClient) GetMoneyLogs(ctx context.Context, req *GetMoneyLogsReq, opts ...grpc.CallOption) (*GetMoneyLogsResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(GetMoneyLogsResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathMoneyGetMoneyLogs, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}

type MoneyCGrpcServer interface {
	GetMoneyLogs(context.Context, *GetMoneyLogsReq) (*GetMoneyLogsResp, error)
}

func _Money_GetMoneyLogs_CGrpcHandler(handler func(context.Context, *GetMoneyLogsReq) (*GetMoneyLogsResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(GetMoneyLogsReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.money/get_money_logs] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.money/get_money_logs] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(GetMoneyLogsResp)
		}
		ctx.Write(resp)
	}
}
func RegisterMoneyCGrpcServer(engine *cgrpc.CGrpcServer, svc MoneyCGrpcServer, allmids map[string]cgrpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	engine.RegisterHandler("account.money", "get_money_logs", _Money_GetMoneyLogs_CGrpcHandler(svc.GetMoneyLogs))
}
