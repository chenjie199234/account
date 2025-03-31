// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: api/account_status.proto

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

// req can be set with pbex extentions
type Pingreq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Timestamp     int64                  `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Pingreq) Reset() {
	*x = Pingreq{}
	mi := &file_api_account_status_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Pingreq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pingreq) ProtoMessage() {}

func (x *Pingreq) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_status_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pingreq.ProtoReflect.Descriptor instead.
func (*Pingreq) Descriptor() ([]byte, []int) {
	return file_api_account_status_proto_rawDescGZIP(), []int{0}
}

func (x *Pingreq) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

// resp's pbex extentions will be ignore
type Pingresp struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	ClientTimestamp int64                  `protobuf:"varint,1,opt,name=client_timestamp,json=clientTimestamp,proto3" json:"client_timestamp,omitempty"`
	ServerTimestamp int64                  `protobuf:"varint,2,opt,name=server_timestamp,json=serverTimestamp,proto3" json:"server_timestamp,omitempty"`
	TotalMem        uint64                 `protobuf:"varint,3,opt,name=total_mem,json=totalMem,proto3" json:"total_mem,omitempty"`
	CurMemUsage     uint64                 `protobuf:"varint,4,opt,name=cur_mem_usage,json=curMemUsage,proto3" json:"cur_mem_usage,omitempty"`
	MaxMemUsage     uint64                 `protobuf:"varint,5,opt,name=max_mem_usage,json=maxMemUsage,proto3" json:"max_mem_usage,omitempty"`
	CpuNum          float64                `protobuf:"fixed64,6,opt,name=cpu_num,json=cpuNum,proto3" json:"cpu_num,omitempty"`
	CurCpuUsage     float64                `protobuf:"fixed64,7,opt,name=cur_cpu_usage,json=curCpuUsage,proto3" json:"cur_cpu_usage,omitempty"`
	AvgCpuUsage     float64                `protobuf:"fixed64,8,opt,name=avg_cpu_usage,json=avgCpuUsage,proto3" json:"avg_cpu_usage,omitempty"`
	MaxCpuUsage     float64                `protobuf:"fixed64,9,opt,name=max_cpu_usage,json=maxCpuUsage,proto3" json:"max_cpu_usage,omitempty"`
	Host            string                 `protobuf:"bytes,10,opt,name=host,proto3" json:"host,omitempty"`
	Ip              string                 `protobuf:"bytes,11,opt,name=ip,proto3" json:"ip,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *Pingresp) Reset() {
	*x = Pingresp{}
	mi := &file_api_account_status_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Pingresp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pingresp) ProtoMessage() {}

func (x *Pingresp) ProtoReflect() protoreflect.Message {
	mi := &file_api_account_status_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pingresp.ProtoReflect.Descriptor instead.
func (*Pingresp) Descriptor() ([]byte, []int) {
	return file_api_account_status_proto_rawDescGZIP(), []int{1}
}

func (x *Pingresp) GetClientTimestamp() int64 {
	if x != nil {
		return x.ClientTimestamp
	}
	return 0
}

func (x *Pingresp) GetServerTimestamp() int64 {
	if x != nil {
		return x.ServerTimestamp
	}
	return 0
}

func (x *Pingresp) GetTotalMem() uint64 {
	if x != nil {
		return x.TotalMem
	}
	return 0
}

func (x *Pingresp) GetCurMemUsage() uint64 {
	if x != nil {
		return x.CurMemUsage
	}
	return 0
}

func (x *Pingresp) GetMaxMemUsage() uint64 {
	if x != nil {
		return x.MaxMemUsage
	}
	return 0
}

func (x *Pingresp) GetCpuNum() float64 {
	if x != nil {
		return x.CpuNum
	}
	return 0
}

func (x *Pingresp) GetCurCpuUsage() float64 {
	if x != nil {
		return x.CurCpuUsage
	}
	return 0
}

func (x *Pingresp) GetAvgCpuUsage() float64 {
	if x != nil {
		return x.AvgCpuUsage
	}
	return 0
}

func (x *Pingresp) GetMaxCpuUsage() float64 {
	if x != nil {
		return x.MaxCpuUsage
	}
	return 0
}

func (x *Pingresp) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *Pingresp) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

var File_api_account_status_proto protoreflect.FileDescriptor

const file_api_account_status_proto_rawDesc = "" +
	"\n" +
	"\x18api/account_status.proto\x12\aaccount\x1a\x0fpbex/pbex.proto\"-\n" +
	"\apingreq\x12\"\n" +
	"\ttimestamp\x18\x01 \x01(\x03B\x04\x80\x92N\x00R\ttimestamp\"\xee\x02\n" +
	"\bpingresp\x12)\n" +
	"\x10client_timestamp\x18\x01 \x01(\x03R\x0fclientTimestamp\x12)\n" +
	"\x10server_timestamp\x18\x02 \x01(\x03R\x0fserverTimestamp\x12\x1b\n" +
	"\ttotal_mem\x18\x03 \x01(\x04R\btotalMem\x12\"\n" +
	"\rcur_mem_usage\x18\x04 \x01(\x04R\vcurMemUsage\x12\"\n" +
	"\rmax_mem_usage\x18\x05 \x01(\x04R\vmaxMemUsage\x12\x17\n" +
	"\acpu_num\x18\x06 \x01(\x01R\x06cpuNum\x12\"\n" +
	"\rcur_cpu_usage\x18\a \x01(\x01R\vcurCpuUsage\x12\"\n" +
	"\ravg_cpu_usage\x18\b \x01(\x01R\vavgCpuUsage\x12\"\n" +
	"\rmax_cpu_usage\x18\t \x01(\x01R\vmaxCpuUsage\x12\x12\n" +
	"\x04host\x18\n" +
	" \x01(\tR\x04host\x12\x0e\n" +
	"\x02ip\x18\v \x01(\tR\x02ip2N\n" +
	"\x06status\x12D\n" +
	"\x04ping\x12\x10.account.pingreq\x1a\x11.account.pingresp\"\x17\x8a\x9fI\x03get\x8a\x9fI\x04crpc\x8a\x9fI\x04grpcB*Z(github.com/chenjie199234/account/api;apib\x06proto3"

var (
	file_api_account_status_proto_rawDescOnce sync.Once
	file_api_account_status_proto_rawDescData []byte
)

func file_api_account_status_proto_rawDescGZIP() []byte {
	file_api_account_status_proto_rawDescOnce.Do(func() {
		file_api_account_status_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_account_status_proto_rawDesc), len(file_api_account_status_proto_rawDesc)))
	})
	return file_api_account_status_proto_rawDescData
}

var file_api_account_status_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_account_status_proto_goTypes = []any{
	(*Pingreq)(nil),  // 0: account.pingreq
	(*Pingresp)(nil), // 1: account.pingresp
}
var file_api_account_status_proto_depIdxs = []int32{
	0, // 0: account.status.ping:input_type -> account.pingreq
	1, // 1: account.status.ping:output_type -> account.pingresp
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_account_status_proto_init() }
func file_api_account_status_proto_init() {
	if File_api_account_status_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_account_status_proto_rawDesc), len(file_api_account_status_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_account_status_proto_goTypes,
		DependencyIndexes: file_api_account_status_proto_depIdxs,
		MessageInfos:      file_api_account_status_proto_msgTypes,
	}.Build()
	File_api_account_status_proto = out.File
	file_api_account_status_proto_goTypes = nil
	file_api_account_status_proto_depIdxs = nil
}
