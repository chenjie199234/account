// Code generated by protoc-gen-go-crpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-crpc v0.0.131<br />
// 	protoc             v5.28.2<br />
// source: api/account_money.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	crpc "github.com/chenjie199234/Corelib/crpc"
	proto "google.golang.org/protobuf/proto"
	slog "log/slog"
)

var _CrpcPathMoneyGetMoneyLogs = "/account.money/get_money_logs"

type MoneyCrpcClient interface {
	GetMoneyLogs(ctx context.Context, req *GetMoneyLogsReq) (resp *GetMoneyLogsResp, e error)
}

type moneyCrpcClient struct {
	cc *crpc.CrpcClient
}

func NewMoneyCrpcClient(c *crpc.CrpcClient) MoneyCrpcClient {
	return &moneyCrpcClient{cc: c}
}

func (c *moneyCrpcClient) GetMoneyLogs(ctx context.Context, req *GetMoneyLogsReq) (*GetMoneyLogsResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if errstr := req.Validate(); errstr != "" {
		slog.ErrorContext(ctx, "[/account.money/get_money_logs] request validate failed", slog.String("error", errstr))
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	var respbody []byte
	if e := c.cc.Call(ctx, _CrpcPathMoneyGetMoneyLogs, reqd, func(ctx *crpc.CallContext) error {
		var e error
		if respbody, e = ctx.Recv(); e != nil {
			slog.ErrorContext(ctx, "[/account.money/get_money_logs] read response failed", slog.String("error", e.Error()))
		}
		return e
	}); e != nil {
		return nil, e
	}
	resp := new(GetMoneyLogsResp)
	if len(respbody) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respbody, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}

type MoneyCrpcServer interface {
	// Context is *crpc.NoStreamServerContext
	GetMoneyLogs(context.Context, *GetMoneyLogsReq) (*GetMoneyLogsResp, error)
}

func _Money_GetMoneyLogs_CrpcHandler(handler func(context.Context, *GetMoneyLogsReq) (*GetMoneyLogsResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.ServerContext) {
		reqbody, e := ctx.Recv()
		if e != nil {
			slog.ErrorContext(ctx, "[/account.money/get_money_logs] read request failed", slog.String("error", e.Error()))
			ctx.Abort(e)
			return
		}
		req := new(GetMoneyLogsReq)
		if e := proto.Unmarshal(reqbody, req); e != nil {
			slog.ErrorContext(ctx, "[/account.money/get_money_logs] request decode failed", slog.String("error", e.Error()))
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			slog.ErrorContext(ctx, "[/account.money/get_money_logs] request validate failed", slog.String("error", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(crpc.NewNoStreamServerContext(ctx), req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(GetMoneyLogsResp)
		}
		respd, _ := proto.Marshal(resp)
		if e := ctx.Send(respd); e != nil {
			slog.ErrorContext(ctx, "[/account.money/get_money_logs] send response failed", slog.String("error", e.Error()))
		}
	}
}
func RegisterMoneyCrpcServer(engine *crpc.CrpcServer, svc MoneyCrpcServer, allmids map[string]crpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	engine.RegisterHandler("account.money", "get_money_logs", _Money_GetMoneyLogs_CrpcHandler(svc.GetMoneyLogs))
}
