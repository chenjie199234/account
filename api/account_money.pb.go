// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: api/account_money.proto

//this is the proto package name,all proto in this project must use this name as the proto package name

package api

import (
	_ "github.com/chenjie199234/Corelib/pbex"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MoneyLog struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Action        string                 `protobuf:"bytes,2,opt,name=action,proto3" json:"action,omitempty"` //spend,recharge,refund
	UniqueId      string                 `protobuf:"bytes,3,opt,name=unique_id,json=uniqueId,proto3" json:"unique_id,omitempty"`
	SrcDst        string                 `protobuf:"bytes,4,opt,name=src_dst,json=srcDst,proto3" json:"src_dst,omitempty"`
	MoneyType     string                 `protobuf:"bytes,5,opt,name=money_type,json=moneyType,proto3" json:"money_type,omitempty"`
	MoneyAmount   uint32                 `protobuf:"varint,6,opt,name=money_amount,json=moneyAmount,proto3" json:"money_amount,omitempty"`
	Ctime         uint32                 `protobuf:"varint,7,opt,name=ctime,proto3" json:"ctime,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MoneyLog) Reset() {
	*x = MoneyLog{}
	mi := &file_api_account_money_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MoneyLog) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MoneyLog) ProtoMessage() {}

func (x *MoneyLog) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_money_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MoneyLog.ProtoReflect.Descriptor instead.
func (*MoneyLog) Descriptor() ([]byte, []int) {
	return file_api_account_money_proto_rawDescGZIP(), []int{0}
}

func (x *MoneyLog) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *MoneyLog) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

func (x *MoneyLog) GetUniqueId() string {
	if x != nil {
		return x.UniqueId
	}
	return ""
}

func (x *MoneyLog) GetSrcDst() string {
	if x != nil {
		return x.SrcDst
	}
	return ""
}

func (x *MoneyLog) GetMoneyType() string {
	if x != nil {
		return x.MoneyType
	}
	return ""
}

func (x *MoneyLog) GetMoneyAmount() uint32 {
	if x != nil {
		return x.MoneyAmount
	}
	return 0
}

func (x *MoneyLog) GetCtime() uint32 {
	if x != nil {
		return x.Ctime
	}
	return 0
}

type GetMoneyLogsReq struct {
	state   protoimpl.MessageState `protogen:"open.v1"`
	SrcType string                 `protobuf:"bytes,1,opt,name=src_type,json=srcType,proto3" json:"src_type,omitempty"`
	Src     string                 `protobuf:"bytes,2,opt,name=src,proto3" json:"src,omitempty"`
	// >0:return the required page's data
	//
	//0:return all logs
	StartTime     uint32 `protobuf:"varint,3,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"` //unit second
	EndTime       uint32 `protobuf:"varint,4,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`       //unit second
	Page          uint32 `protobuf:"varint,5,opt,name=page,proto3" json:"page,omitempty"`
	Action        string `protobuf:"bytes,6,opt,name=action,proto3" json:"action,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetMoneyLogsReq) Reset() {
	*x = GetMoneyLogsReq{}
	mi := &file_api_account_money_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMoneyLogsReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMoneyLogsReq) ProtoMessage() {}

func (x *GetMoneyLogsReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_money_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMoneyLogsReq.ProtoReflect.Descriptor instead.
func (*GetMoneyLogsReq) Descriptor() ([]byte, []int) {
	return file_api_account_money_proto_rawDescGZIP(), []int{1}
}

func (x *GetMoneyLogsReq) GetSrcType() string {
	if x != nil {
		return x.SrcType
	}
	return ""
}

func (x *GetMoneyLogsReq) GetSrc() string {
	if x != nil {
		return x.Src
	}
	return ""
}

func (x *GetMoneyLogsReq) GetStartTime() uint32 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *GetMoneyLogsReq) GetEndTime() uint32 {
	if x != nil {
		return x.EndTime
	}
	return 0
}

func (x *GetMoneyLogsReq) GetPage() uint32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *GetMoneyLogsReq) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

type GetMoneyLogsResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          uint32                 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Pagesize      uint32                 `protobuf:"varint,2,opt,name=pagesize,proto3" json:"pagesize,omitempty"`
	Totalsize     uint32                 `protobuf:"varint,3,opt,name=totalsize,proto3" json:"totalsize,omitempty"`
	Logs          []*MoneyLog            `protobuf:"bytes,4,rep,name=logs,proto3" json:"logs,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetMoneyLogsResp) Reset() {
	*x = GetMoneyLogsResp{}
	mi := &file_api_account_money_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMoneyLogsResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMoneyLogsResp) ProtoMessage() {}

func (x *GetMoneyLogsResp) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_money_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMoneyLogsResp.ProtoReflect.Descriptor instead.
func (*GetMoneyLogsResp) Descriptor() ([]byte, []int) {
	return file_api_account_money_proto_rawDescGZIP(), []int{2}
}

func (x *GetMoneyLogsResp) GetPage() uint32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *GetMoneyLogsResp) GetPagesize() uint32 {
	if x != nil {
		return x.Pagesize
	}
	return 0
}

func (x *GetMoneyLogsResp) GetTotalsize() uint32 {
	if x != nil {
		return x.Totalsize
	}
	return 0
}

func (x *GetMoneyLogsResp) GetLogs() []*MoneyLog {
	if x != nil {
		return x.Logs
	}
	return nil
}

type SelfMoneyLogsReq struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// >0:return the required page's data
	//
	//0:return all logs
	StartTime     uint32 `protobuf:"varint,1,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime       uint32 `protobuf:"varint,2,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	Page          uint32 `protobuf:"varint,3,opt,name=page,proto3" json:"page,omitempty"`
	Action        string `protobuf:"bytes,4,opt,name=action,proto3" json:"action,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SelfMoneyLogsReq) Reset() {
	*x = SelfMoneyLogsReq{}
	mi := &file_api_account_money_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SelfMoneyLogsReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SelfMoneyLogsReq) ProtoMessage() {}

func (x *SelfMoneyLogsReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_money_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SelfMoneyLogsReq.ProtoReflect.Descriptor instead.
func (*SelfMoneyLogsReq) Descriptor() ([]byte, []int) {
	return file_api_account_money_proto_rawDescGZIP(), []int{3}
}

func (x *SelfMoneyLogsReq) GetStartTime() uint32 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *SelfMoneyLogsReq) GetEndTime() uint32 {
	if x != nil {
		return x.EndTime
	}
	return 0
}

func (x *SelfMoneyLogsReq) GetPage() uint32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *SelfMoneyLogsReq) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

type SelfMoneyLogsResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          uint32                 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Pagesize      uint32                 `protobuf:"varint,2,opt,name=pagesize,proto3" json:"pagesize,omitempty"`
	Totalsize     uint32                 `protobuf:"varint,3,opt,name=totalsize,proto3" json:"totalsize,omitempty"`
	Logs          []*MoneyLog            `protobuf:"bytes,4,rep,name=logs,proto3" json:"logs,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SelfMoneyLogsResp) Reset() {
	*x = SelfMoneyLogsResp{}
	mi := &file_api_account_money_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SelfMoneyLogsResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SelfMoneyLogsResp) ProtoMessage() {}

func (x *SelfMoneyLogsResp) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_money_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SelfMoneyLogsResp.ProtoReflect.Descriptor instead.
func (*SelfMoneyLogsResp) Descriptor() ([]byte, []int) {
	return file_api_account_money_proto_rawDescGZIP(), []int{4}
}

func (x *SelfMoneyLogsResp) GetPage() uint32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *SelfMoneyLogsResp) GetPagesize() uint32 {
	if x != nil {
		return x.Pagesize
	}
	return 0
}

func (x *SelfMoneyLogsResp) GetTotalsize() uint32 {
	if x != nil {
		return x.Totalsize
	}
	return 0
}

func (x *SelfMoneyLogsResp) GetLogs() []*MoneyLog {
	if x != nil {
		return x.Logs
	}
	return nil
}

type RechargeMoneyReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RechargeMoneyReq) Reset() {
	*x = RechargeMoneyReq{}
	mi := &file_api_account_money_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RechargeMoneyReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RechargeMoneyReq) ProtoMessage() {}

func (x *RechargeMoneyReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_money_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RechargeMoneyReq.ProtoReflect.Descriptor instead.
func (*RechargeMoneyReq) Descriptor() ([]byte, []int) {
	return file_api_account_money_proto_rawDescGZIP(), []int{5}
}

type RechargeMoneyResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RechargeMoneyResp) Reset() {
	*x = RechargeMoneyResp{}
	mi := &file_api_account_money_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RechargeMoneyResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RechargeMoneyResp) ProtoMessage() {}

func (x *RechargeMoneyResp) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_money_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RechargeMoneyResp.ProtoReflect.Descriptor instead.
func (*RechargeMoneyResp) Descriptor() ([]byte, []int) {
	return file_api_account_money_proto_rawDescGZIP(), []int{6}
}

type SpendMoneyReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SpendMoneyReq) Reset() {
	*x = SpendMoneyReq{}
	mi := &file_api_account_money_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SpendMoneyReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SpendMoneyReq) ProtoMessage() {}

func (x *SpendMoneyReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_money_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SpendMoneyReq.ProtoReflect.Descriptor instead.
func (*SpendMoneyReq) Descriptor() ([]byte, []int) {
	return file_api_account_money_proto_rawDescGZIP(), []int{7}
}

type SpendMoneyResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SpendMoneyResp) Reset() {
	*x = SpendMoneyResp{}
	mi := &file_api_account_money_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SpendMoneyResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SpendMoneyResp) ProtoMessage() {}

func (x *SpendMoneyResp) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_money_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SpendMoneyResp.ProtoReflect.Descriptor instead.
func (*SpendMoneyResp) Descriptor() ([]byte, []int) {
	return file_api_account_money_proto_rawDescGZIP(), []int{8}
}

type RefundMoneyReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RefundMoneyReq) Reset() {
	*x = RefundMoneyReq{}
	mi := &file_api_account_money_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RefundMoneyReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefundMoneyReq) ProtoMessage() {}

func (x *RefundMoneyReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_money_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefundMoneyReq.ProtoReflect.Descriptor instead.
func (*RefundMoneyReq) Descriptor() ([]byte, []int) {
	return file_api_account_money_proto_rawDescGZIP(), []int{9}
}

type RefundMoneyResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RefundMoneyResp) Reset() {
	*x = RefundMoneyResp{}
	mi := &file_api_account_money_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RefundMoneyResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefundMoneyResp) ProtoMessage() {}

func (x *RefundMoneyResp) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_money_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefundMoneyResp.ProtoReflect.Descriptor instead.
func (*RefundMoneyResp) Descriptor() ([]byte, []int) {
	return file_api_account_money_proto_rawDescGZIP(), []int{10}
}

var File_api_account_money_proto protoreflect.FileDescriptor

const file_api_account_money_proto_rawDesc = "" +
	"\n" +
	"\x17api/account_money.proto\x12\aaccount\x1a\x0fpbex/pbex.proto\"\xca\x01\n" +
	"\tmoney_log\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\x16\n" +
	"\x06action\x18\x02 \x01(\tR\x06action\x12\x1b\n" +
	"\tunique_id\x18\x03 \x01(\tR\buniqueId\x12\x17\n" +
	"\asrc_dst\x18\x04 \x01(\tR\x06srcDst\x12\x1d\n" +
	"\n" +
	"money_type\x18\x05 \x01(\tR\tmoneyType\x12!\n" +
	"\fmoney_amount\x18\x06 \x01(\rR\vmoneyAmount\x12\x14\n" +
	"\x05ctime\x18\a \x01(\rR\x05ctime\"\x88\x02\n" +
	"\x12get_money_logs_req\x12@\n" +
	"\bsrc_type\x18\x01 \x01(\tB%ҐN\auser_idҐN\x03telҐN\x05emailҐN\x06idcardR\asrcType\x12\x16\n" +
	"\x03src\x18\x02 \x01(\tB\x04\xf0\x90N\x00R\x03src\x12#\n" +
	"\n" +
	"start_time\x18\x03 \x01(\rB\x04\xa0\x91N\x00R\tstartTime\x12\x1f\n" +
	"\bend_time\x18\x04 \x01(\rB\x04\xa0\x91N\x00R\aendTime\x12\x12\n" +
	"\x04page\x18\x05 \x01(\rR\x04page\x12>\n" +
	"\x06action\x18\x06 \x01(\tB&ҐN\x05spendҐN\brechargeҐN\x06refundҐN\x03allR\x06action\"\x8b\x01\n" +
	"\x13get_money_logs_resp\x12\x12\n" +
	"\x04page\x18\x01 \x01(\rR\x04page\x12\x1a\n" +
	"\bpagesize\x18\x02 \x01(\rR\bpagesize\x12\x1c\n" +
	"\ttotalsize\x18\x03 \x01(\rR\ttotalsize\x12&\n" +
	"\x04logs\x18\x04 \x03(\v2\x12.account.money_logR\x04logs\"\xaf\x01\n" +
	"\x13self_money_logs_req\x12#\n" +
	"\n" +
	"start_time\x18\x01 \x01(\rB\x04\xa0\x91N\x00R\tstartTime\x12\x1f\n" +
	"\bend_time\x18\x02 \x01(\rB\x04\xa0\x91N\x00R\aendTime\x12\x12\n" +
	"\x04page\x18\x03 \x01(\rR\x04page\x12>\n" +
	"\x06action\x18\x04 \x01(\tB&ҐN\x05spendҐN\brechargeҐN\x06refundҐN\x03allR\x06action\"\x8c\x01\n" +
	"\x14self_money_logs_resp\x12\x12\n" +
	"\x04page\x18\x01 \x01(\rR\x04page\x12\x1a\n" +
	"\bpagesize\x18\x02 \x01(\rR\bpagesize\x12\x1c\n" +
	"\ttotalsize\x18\x03 \x01(\rR\ttotalsize\x12&\n" +
	"\x04logs\x18\x04 \x03(\v2\x12.account.money_logR\x04logs\"\x14\n" +
	"\x12recharge_money_req\"\x15\n" +
	"\x13recharge_money_resp\"\x11\n" +
	"\x0fspend_money_req\"\x12\n" +
	"\x10spend_money_resp\"\x12\n" +
	"\x10refund_money_req\"\x13\n" +
	"\x11refund_money_resp2\xa1\x03\n" +
	"\x05money\x12]\n" +
	"\x0eget_money_logs\x12\x1b.account.get_money_logs_req\x1a\x1c.account.get_money_logs_resp\"\x10\x8a\x9fI\x04crpc\x8a\x9fI\x04grpc\x12a\n" +
	"\x0fself_money_logs\x12\x1c.account.self_money_logs_req\x1a\x1d.account.self_money_logs_resp\"\x11\x8a\x9fI\x04post\x92\x9fI\x05token\x12K\n" +
	"\x0erecharge_money\x12\x1b.account.recharge_money_req\x1a\x1c.account.recharge_money_resp\x12B\n" +
	"\vspend_money\x12\x18.account.spend_money_req\x1a\x19.account.spend_money_resp\x12E\n" +
	"\frefund_money\x12\x19.account.refund_money_req\x1a\x1a.account.refund_money_respB*Z(github.com/chenjie199234/account/api;apib\x06proto3"

var (
	file_api_account_money_proto_rawDescOnce sync.Once
	file_api_account_money_proto_rawDescData []byte
)

func file_api_account_money_proto_rawDescGZIP() []byte {
	file_api_account_money_proto_rawDescOnce.Do(func() {
		file_api_account_money_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_account_money_proto_rawDesc), len(file_api_account_money_proto_rawDesc)))
	})
	return file_api_account_money_proto_rawDescData
}

var file_api_account_money_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_api_account_money_proto_goTypes = []any{
	(*MoneyLog)(nil),          // 0: account.money_log
	(*GetMoneyLogsReq)(nil),   // 1: account.get_money_logs_req
	(*GetMoneyLogsResp)(nil),  // 2: account.get_money_logs_resp
	(*SelfMoneyLogsReq)(nil),  // 3: account.self_money_logs_req
	(*SelfMoneyLogsResp)(nil), // 4: account.self_money_logs_resp
	(*RechargeMoneyReq)(nil),  // 5: account.recharge_money_req
	(*RechargeMoneyResp)(nil), // 6: account.recharge_money_resp
	(*SpendMoneyReq)(nil),     // 7: account.spend_money_req
	(*SpendMoneyResp)(nil),    // 8: account.spend_money_resp
	(*RefundMoneyReq)(nil),    // 9: account.refund_money_req
	(*RefundMoneyResp)(nil),   // 10: account.refund_money_resp
}
var file_api_account_money_proto_depIdxs = []int32{
	0,  // 0: account.get_money_logs_resp.logs:type_name -> account.money_log
	0,  // 1: account.self_money_logs_resp.logs:type_name -> account.money_log
	1,  // 2: account.money.get_money_logs:input_type -> account.get_money_logs_req
	3,  // 3: account.money.self_money_logs:input_type -> account.self_money_logs_req
	5,  // 4: account.money.recharge_money:input_type -> account.recharge_money_req
	7,  // 5: account.money.spend_money:input_type -> account.spend_money_req
	9,  // 6: account.money.refund_money:input_type -> account.refund_money_req
	2,  // 7: account.money.get_money_logs:output_type -> account.get_money_logs_resp
	4,  // 8: account.money.self_money_logs:output_type -> account.self_money_logs_resp
	6,  // 9: account.money.recharge_money:output_type -> account.recharge_money_resp
	8,  // 10: account.money.spend_money:output_type -> account.spend_money_resp
	10, // 11: account.money.refund_money:output_type -> account.refund_money_resp
	7,  // [7:12] is the sub-list for method output_type
	2,  // [2:7] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_api_account_money_proto_init() }
func file_api_account_money_proto_init() {
	if File_api_account_money_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_account_money_proto_rawDesc), len(file_api_account_money_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_account_money_proto_goTypes,
		DependencyIndexes: file_api_account_money_proto_depIdxs,
		MessageInfos:      file_api_account_money_proto_msgTypes,
	}.Build()
	File_api_account_money_proto = out.File
	file_api_account_money_proto_goTypes = nil
	file_api_account_money_proto_depIdxs = nil
}
