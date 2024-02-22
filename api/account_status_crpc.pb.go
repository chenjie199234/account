// Code generated by protoc-gen-go-crpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-crpc v0.0.98<br />
// 	protoc             v4.25.3<br />
// source: api/account_status.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	crpc "github.com/chenjie199234/Corelib/crpc"
	log "github.com/chenjie199234/Corelib/log"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
)

var _CrpcPathStatusPing = "/account.status/ping"

type StatusCrpcClient interface {
	// ping check server's health
	Ping(context.Context, *Pingreq) (*Pingresp, error)
}

type statusCrpcClient struct {
	cc *crpc.CrpcClient
}

func NewStatusCrpcClient(c *crpc.CrpcClient) StatusCrpcClient {
	return &statusCrpcClient{cc: c}
}

func (c *statusCrpcClient) Ping(ctx context.Context, req *Pingreq) (*Pingresp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathStatusPing, reqd)
	if e != nil {
		return nil, e
	}
	resp := new(Pingresp)
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

type StatusCrpcServer interface {
	// ping check server's health
	Ping(context.Context, *Pingreq) (*Pingresp, error)
}

func _Status_Ping_CrpcHandler(handler func(context.Context, *Pingreq) (*Pingresp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(Pingreq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/account.status/ping] json and proto format decode both failed")
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/account.status/ping] json and proto format decode both failed")
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
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
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func RegisterStatusCrpcServer(engine *crpc.CrpcServer, svc StatusCrpcServer, allmids map[string]crpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	engine.RegisterHandler("account.status", "ping", _Status_Ping_CrpcHandler(svc.Ping))
}
