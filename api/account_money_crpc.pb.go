// Code generated by protoc-gen-go-crpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-crpc v0.0.92<br />
// 	protoc             v4.24.4<br />
// source: api/account_money.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	crpc "github.com/chenjie199234/Corelib/crpc"
	log "github.com/chenjie199234/Corelib/log"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
)

var _CrpcPathMoneyGetUserMoneyLogs = "/account.money/get_user_money_logs"
var _CrpcPathMoneySelfMoneyLogs = "/account.money/self_money_logs"
var _CrpcPathMoneyRechargeMoney = "/account.money/recharge_money"
var _CrpcPathMoneySpendMoney = "/account.money/spend_money"
var _CrpcPathMoneyRefundMoney = "/account.money/refund_money"

type MoneyCrpcClient interface {
	GetUserMoneyLogs(context.Context, *GetUserMoneyLogsReq) (*GetUserMoneyLogsResp, error)
	SelfMoneyLogs(context.Context, *SelfMoneyLogsReq) (*SelfMoneyLogsResp, error)
	RechargeMoney(context.Context, *RechargeMoneyReq) (*RechargeMoneyResp, error)
	SpendMoney(context.Context, *SpendMoneyReq) (*SpendMoneyResp, error)
	RefundMoney(context.Context, *RefundMoneyReq) (*RefundMoneyResp, error)
}

type moneyCrpcClient struct {
	cc *crpc.CrpcClient
}

func NewMoneyCrpcClient(c *crpc.CrpcClient) MoneyCrpcClient {
	return &moneyCrpcClient{cc: c}
}

func (c *moneyCrpcClient) GetUserMoneyLogs(ctx context.Context, req *GetUserMoneyLogsReq) (*GetUserMoneyLogsResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathMoneyGetUserMoneyLogs, reqd)
	if e != nil {
		return nil, e
	}
	resp := new(GetUserMoneyLogsResp)
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
func (c *moneyCrpcClient) SelfMoneyLogs(ctx context.Context, req *SelfMoneyLogsReq) (*SelfMoneyLogsResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathMoneySelfMoneyLogs, reqd)
	if e != nil {
		return nil, e
	}
	resp := new(SelfMoneyLogsResp)
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
func (c *moneyCrpcClient) RechargeMoney(ctx context.Context, req *RechargeMoneyReq) (*RechargeMoneyResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathMoneyRechargeMoney, reqd)
	if e != nil {
		return nil, e
	}
	resp := new(RechargeMoneyResp)
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
func (c *moneyCrpcClient) SpendMoney(ctx context.Context, req *SpendMoneyReq) (*SpendMoneyResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathMoneySpendMoney, reqd)
	if e != nil {
		return nil, e
	}
	resp := new(SpendMoneyResp)
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
func (c *moneyCrpcClient) RefundMoney(ctx context.Context, req *RefundMoneyReq) (*RefundMoneyResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathMoneyRefundMoney, reqd)
	if e != nil {
		return nil, e
	}
	resp := new(RefundMoneyResp)
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

type MoneyCrpcServer interface {
	GetUserMoneyLogs(context.Context, *GetUserMoneyLogsReq) (*GetUserMoneyLogsResp, error)
	SelfMoneyLogs(context.Context, *SelfMoneyLogsReq) (*SelfMoneyLogsResp, error)
	RechargeMoney(context.Context, *RechargeMoneyReq) (*RechargeMoneyResp, error)
	SpendMoney(context.Context, *SpendMoneyReq) (*SpendMoneyResp, error)
	RefundMoney(context.Context, *RefundMoneyReq) (*RefundMoneyResp, error)
}

func _Money_GetUserMoneyLogs_CrpcHandler(handler func(context.Context, *GetUserMoneyLogsReq) (*GetUserMoneyLogsResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(GetUserMoneyLogsReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/account.money/get_user_money_logs] json and proto format decode both failed")
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/account.money/get_user_money_logs] json and proto format decode both failed")
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
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
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func _Money_SelfMoneyLogs_CrpcHandler(handler func(context.Context, *SelfMoneyLogsReq) (*SelfMoneyLogsResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(SelfMoneyLogsReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/account.money/self_money_logs] json and proto format decode both failed")
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/account.money/self_money_logs] json and proto format decode both failed")
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
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
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func _Money_RechargeMoney_CrpcHandler(handler func(context.Context, *RechargeMoneyReq) (*RechargeMoneyResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(RechargeMoneyReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/account.money/recharge_money] json and proto format decode both failed")
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/account.money/recharge_money] json and proto format decode both failed")
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(RechargeMoneyResp)
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
func _Money_SpendMoney_CrpcHandler(handler func(context.Context, *SpendMoneyReq) (*SpendMoneyResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(SpendMoneyReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/account.money/spend_money] json and proto format decode both failed")
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/account.money/spend_money] json and proto format decode both failed")
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(SpendMoneyResp)
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
func _Money_RefundMoney_CrpcHandler(handler func(context.Context, *RefundMoneyReq) (*RefundMoneyResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(RefundMoneyReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/account.money/refund_money] json and proto format decode both failed")
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/account.money/refund_money] json and proto format decode both failed")
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(RefundMoneyResp)
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
func RegisterMoneyCrpcServer(engine *crpc.CrpcServer, svc MoneyCrpcServer, allmids map[string]crpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	engine.RegisterHandler("account.money", "get_user_money_logs", _Money_GetUserMoneyLogs_CrpcHandler(svc.GetUserMoneyLogs))
	engine.RegisterHandler("account.money", "self_money_logs", _Money_SelfMoneyLogs_CrpcHandler(svc.SelfMoneyLogs))
	engine.RegisterHandler("account.money", "recharge_money", _Money_RechargeMoney_CrpcHandler(svc.RechargeMoney))
	engine.RegisterHandler("account.money", "spend_money", _Money_SpendMoney_CrpcHandler(svc.SpendMoney))
	engine.RegisterHandler("account.money", "refund_money", _Money_RefundMoney_CrpcHandler(svc.RefundMoney))
}
