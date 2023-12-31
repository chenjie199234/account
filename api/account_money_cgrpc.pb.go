// Code generated by protoc-gen-go-cgrpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-cgrpc v0.0.94<br />
// 	protoc              v4.25.1<br />
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
var _CGrpcPathMoneySelfMoneyLogs = "/account.money/self_money_logs"
var _CGrpcPathMoneyRechargeMoney = "/account.money/recharge_money"
var _CGrpcPathMoneySpendMoney = "/account.money/spend_money"
var _CGrpcPathMoneyRefundMoney = "/account.money/refund_money"

type MoneyCGrpcClient interface {
	GetUserMoneyLogs(context.Context, *GetUserMoneyLogsReq, ...grpc.CallOption) (*GetUserMoneyLogsResp, error)
	SelfMoneyLogs(context.Context, *SelfMoneyLogsReq, ...grpc.CallOption) (*SelfMoneyLogsResp, error)
	RechargeMoney(context.Context, *RechargeMoneyReq, ...grpc.CallOption) (*RechargeMoneyResp, error)
	SpendMoney(context.Context, *SpendMoneyReq, ...grpc.CallOption) (*SpendMoneyResp, error)
	RefundMoney(context.Context, *RefundMoneyReq, ...grpc.CallOption) (*RefundMoneyResp, error)
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
func (c *moneyCGrpcClient) SelfMoneyLogs(ctx context.Context, req *SelfMoneyLogsReq, opts ...grpc.CallOption) (*SelfMoneyLogsResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(SelfMoneyLogsResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathMoneySelfMoneyLogs, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *moneyCGrpcClient) RechargeMoney(ctx context.Context, req *RechargeMoneyReq, opts ...grpc.CallOption) (*RechargeMoneyResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(RechargeMoneyResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathMoneyRechargeMoney, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *moneyCGrpcClient) SpendMoney(ctx context.Context, req *SpendMoneyReq, opts ...grpc.CallOption) (*SpendMoneyResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(SpendMoneyResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathMoneySpendMoney, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *moneyCGrpcClient) RefundMoney(ctx context.Context, req *RefundMoneyReq, opts ...grpc.CallOption) (*RefundMoneyResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(RefundMoneyResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathMoneyRefundMoney, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}

type MoneyCGrpcServer interface {
	GetUserMoneyLogs(context.Context, *GetUserMoneyLogsReq) (*GetUserMoneyLogsResp, error)
	SelfMoneyLogs(context.Context, *SelfMoneyLogsReq) (*SelfMoneyLogsResp, error)
	RechargeMoney(context.Context, *RechargeMoneyReq) (*RechargeMoneyResp, error)
	SpendMoney(context.Context, *SpendMoneyReq) (*SpendMoneyResp, error)
	RefundMoney(context.Context, *RefundMoneyReq) (*RefundMoneyResp, error)
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
func _Money_SelfMoneyLogs_CGrpcHandler(handler func(context.Context, *SelfMoneyLogsReq) (*SelfMoneyLogsResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(SelfMoneyLogsReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.money/self_money_logs] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.money/self_money_logs] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(SelfMoneyLogsResp)
		}
		ctx.Write(resp)
	}
}
func _Money_RechargeMoney_CGrpcHandler(handler func(context.Context, *RechargeMoneyReq) (*RechargeMoneyResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(RechargeMoneyReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.money/recharge_money] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(RechargeMoneyResp)
		}
		ctx.Write(resp)
	}
}
func _Money_SpendMoney_CGrpcHandler(handler func(context.Context, *SpendMoneyReq) (*SpendMoneyResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(SpendMoneyReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.money/spend_money] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(SpendMoneyResp)
		}
		ctx.Write(resp)
	}
}
func _Money_RefundMoney_CGrpcHandler(handler func(context.Context, *RefundMoneyReq) (*RefundMoneyResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(RefundMoneyReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/account.money/refund_money] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(RefundMoneyResp)
		}
		ctx.Write(resp)
	}
}
func RegisterMoneyCGrpcServer(engine *cgrpc.CGrpcServer, svc MoneyCGrpcServer, allmids map[string]cgrpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	engine.RegisterHandler("account.money", "get_user_money_logs", _Money_GetUserMoneyLogs_CGrpcHandler(svc.GetUserMoneyLogs))
	engine.RegisterHandler("account.money", "self_money_logs", _Money_SelfMoneyLogs_CGrpcHandler(svc.SelfMoneyLogs))
	engine.RegisterHandler("account.money", "recharge_money", _Money_RechargeMoney_CGrpcHandler(svc.RechargeMoney))
	engine.RegisterHandler("account.money", "spend_money", _Money_SpendMoney_CGrpcHandler(svc.SpendMoney))
	engine.RegisterHandler("account.money", "refund_money", _Money_RefundMoney_CGrpcHandler(svc.RefundMoney))
}
