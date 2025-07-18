// Code generated by protoc-gen-go-cgrpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-cgrpc v0.0.139<br />
// 	protoc              v6.31.0<br />
// source: api/account_money.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	cgrpc "github.com/chenjie199234/Corelib/cgrpc"
	grpc "google.golang.org/grpc"
	slog "log/slog"
)

var _CGrpcPathMoneyGetMoneyLogs = "/account.money/get_money_logs"

type MoneyCGrpcClient interface {
	GetMoneyLogs(context.Context, *GetMoneyLogsReq) (*GetMoneyLogsResp, error)
}

type moneyCGrpcClient struct {
	cc grpc.ClientConnInterface
}

func NewMoneyCGrpcClient(cc grpc.ClientConnInterface) MoneyCGrpcClient {
	return &moneyCGrpcClient{cc: cc}
}

func (c *moneyCGrpcClient) GetMoneyLogs(ctx context.Context, req *GetMoneyLogsReq) (*GetMoneyLogsResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if errstr := req.Validate(); errstr != "" {
		slog.ErrorContext(ctx, "[/account.money/get_money_logs] validate failed", slog.String("error", errstr))
		return nil, cerror.ErrReq
	}
	resp := new(GetMoneyLogsResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathMoneyGetMoneyLogs, req, resp); e != nil {
		return nil, e
	}
	return resp, nil
}

type MoneyCGrpcServer interface {
	// Context is *cgrpc.NoStreamServerContext
	GetMoneyLogs(context.Context, *GetMoneyLogsReq) (*GetMoneyLogsResp, error)
}

func _Money_GetMoneyLogs_CGrpcHandler(handler func(context.Context, *GetMoneyLogsReq) (*GetMoneyLogsResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.ServerContext) {
		req := new(GetMoneyLogsReq)
		if e := ctx.Read(req); e != nil {
			slog.ErrorContext(ctx, "[/account.money/get_money_logs] decode failed", slog.String("error", e.Error()))
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			slog.ErrorContext(ctx, "[/account.money/get_money_logs] validate failed", slog.String("error", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(cgrpc.NewNoStreamServerContext(ctx), req)
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
	engine.RegisterHandler("account.money", "get_money_logs", false, false, _Money_GetMoneyLogs_CGrpcHandler(svc.GetMoneyLogs))
}
