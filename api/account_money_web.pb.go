// Code generated by protoc-gen-go-web. DO NOT EDIT.
// version:
// 	protoc-gen-go-web v0.0.109<br />
// 	protoc            v4.25.3<br />
// source: api/account_money.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
	web "github.com/chenjie199234/Corelib/web"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	io "io"
	http "net/http"
	strconv "strconv"
	strings "strings"
)

var _WebPathMoneySelfMoneyLogs = "/account.money/self_money_logs"

type MoneyWebClient interface {
	SelfMoneyLogs(context.Context, *SelfMoneyLogsReq, http.Header) (*SelfMoneyLogsResp, error)
}

type moneyWebClient struct {
	cc *web.WebClient
}

func NewMoneyWebClient(c *web.WebClient) MoneyWebClient {
	return &moneyWebClient{cc: c}
}

func (c *moneyWebClient) SelfMoneyLogs(ctx context.Context, req *SelfMoneyLogsReq, header http.Header) (*SelfMoneyLogsResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathMoneySelfMoneyLogs, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(SelfMoneyLogsResp)
	if len(data) == 0 {
		return resp, nil
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-protobuf") {
		if e := proto.Unmarshal(data, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}

type MoneyWebServer interface {
	SelfMoneyLogs(context.Context, *SelfMoneyLogsReq) (*SelfMoneyLogsResp, error)
}

func _Money_SelfMoneyLogs_WebHandler(handler func(context.Context, *SelfMoneyLogsReq) (*SelfMoneyLogsResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(SelfMoneyLogsReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/account.money/self_money_logs] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/account.money/self_money_logs] unmarshal json body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/account.money/self_money_logs] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/account.money/self_money_logs] unmarshal proto body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			if e := ctx.ParseForm(); e != nil {
				log.Error(ctx, "[/account.money/self_money_logs] parse form failed", log.CError(e))
				ctx.Abort(cerror.ErrReq)
				return
			}
			// req.StartTime
			if form := ctx.GetForm("start_time"); len(form) != 0 {
				if num, e := strconv.ParseUint(form, 10, 32); e != nil {
					log.Error(ctx, "[/account.money/self_money_logs] data format wrong", log.String("field", "start_time"))
					ctx.Abort(cerror.ErrReq)
					return
				} else {
					req.StartTime = uint32(num)
				}
			}
			// req.EndTime
			if form := ctx.GetForm("end_time"); len(form) != 0 {
				if num, e := strconv.ParseUint(form, 10, 32); e != nil {
					log.Error(ctx, "[/account.money/self_money_logs] data format wrong", log.String("field", "end_time"))
					ctx.Abort(cerror.ErrReq)
					return
				} else {
					req.EndTime = uint32(num)
				}
			}
			// req.Page
			if form := ctx.GetForm("page"); len(form) != 0 {
				if num, e := strconv.ParseUint(form, 10, 32); e != nil {
					log.Error(ctx, "[/account.money/self_money_logs] data format wrong", log.String("field", "page"))
					ctx.Abort(cerror.ErrReq)
					return
				} else {
					req.Page = uint32(num)
				}
			}
			// req.Action
			if form := ctx.GetForm("action"); len(form) != 0 {
				req.Action = form
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/account.money/self_money_logs] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := cerror.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(SelfMoneyLogsResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func RegisterMoneyWebServer(router *web.Router, svc MoneyWebServer, allmids map[string]web.OutsideHandler) {
	// avoid lint
	_ = allmids
	{
		requiredMids := []string{"token"}
		mids := make([]web.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _Money_SelfMoneyLogs_WebHandler(svc.SelfMoneyLogs))
		router.Post(_WebPathMoneySelfMoneyLogs, mids...)
	}
}
