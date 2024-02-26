// Code generated by protoc-gen-go-cgrpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-cgrpc v0.0.99<br />
// 	protoc              v4.25.3<br />
// source: api/account_money.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	cgrpc "github.com/chenjie199234/Corelib/cgrpc"
	log "github.com/chenjie199234/Corelib/log"
	grpc "google.golang.org/grpc"
)

var _CGrpcPathMoneyGetUserMoneyLogs = "/account.money/get_user_money_logs"

type MoneyCGrpcClient interface {
	GetUserMoneyLogs(context.Context, *GetUserMoneyLogsReq, ...grpc.CallOption) (*GetUserMoneyLogsResp, error)
}

type moneyCGrpcClient struct {
	cc grpc.ClientConnInterface
}

func NewMoneyCGrpcClient(cc grpc.ClientConnInterface) MoneyCGrpcClient {
	return &moneyCGrpcClient{cc: cc}
}

func (c *moneyCGrpcClient) GetUserMoneyLogs(ctx context.Context, req *GetUserMoneyLogsReq, opts ...grpc.CallOption) (*GetUserMoneyLogsResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(GetUserMoneyLogsResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathMoneyGetUserMoneyLogs, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}

type MoneyCGrpcServer interface {
	GetUserMoneyLogs(context.Context, *GetUserMoneyLogsReq) (*GetUserMoneyLogsResp, error)
}

func _Money_GetUserMoneyLogs_CGrpcHandler(handler func(context.Context, *GetUserMoneyLogsReq) (*GetUserMoneyLogsResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(GetUserMoneyLogsReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.money/get_user_money_logs] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.money/get_user_money_logs] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(GetUserMoneyLogsResp)
		}
		ctx.Write(resp)
	}
}
func RegisterMoneyCGrpcServer(engine *cgrpc.CGrpcServer, svc MoneyCGrpcServer, allmids map[string]cgrpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	engine.RegisterHandler("account.money", "get_user_money_logs", _Money_GetUserMoneyLogs_CGrpcHandler(svc.GetUserMoneyLogs))
}
