// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: gm.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GM_CODE int32

const (
	GM_CODE_GAME_NOTICE     GM_CODE = 0  //单个游戏公告
	GM_CODE_KICK_USER       GM_CODE = 1  // 踢用户
	GM_CODE_BLACK_LIST      GM_CODE = 2  // 加入黑名单
	GM_CODE_UNLOCK_BLACK    GM_CODE = 3  // 解除白名单
	GM_CODE_REWARD_ITEM     GM_CODE = 4  //奖励物品
	GM_CODE_ALL_GAME_NOTICE GM_CODE = 10 //所有游戏公告
)

// Enum value maps for GM_CODE.
var (
	GM_CODE_name = map[int32]string{
		0:  "GAME_NOTICE",
		1:  "KICK_USER",
		2:  "BLACK_LIST",
		3:  "UNLOCK_BLACK",
		4:  "REWARD_ITEM",
		10: "ALL_GAME_NOTICE",
	}
	GM_CODE_value = map[string]int32{
		"GAME_NOTICE":     0,
		"KICK_USER":       1,
		"BLACK_LIST":      2,
		"UNLOCK_BLACK":    3,
		"REWARD_ITEM":     4,
		"ALL_GAME_NOTICE": 10,
	}
)

func (x GM_CODE) Enum() *GM_CODE {
	p := new(GM_CODE)
	*p = x
	return p
}

func (x GM_CODE) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GM_CODE) Descriptor() protoreflect.EnumDescriptor {
	return file_gm_proto_enumTypes[0].Descriptor()
}

func (GM_CODE) Type() protoreflect.EnumType {
	return &file_gm_proto_enumTypes[0]
}

func (x GM_CODE) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GM_CODE.Descriptor instead.
func (GM_CODE) EnumDescriptor() ([]byte, []int) {
	return file_gm_proto_rawDescGZIP(), []int{0}
}

// 心跳,客户端无其他业务请求时，每2秒发送
type GmCodeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code     int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"` //
	GameId   string `protobuf:"bytes,2,opt,name=gameId,proto3" json:"gameId,omitempty"`
	Uid      string `protobuf:"bytes,3,opt,name=uid,proto3" json:"uid,omitempty"`
	Opt      string `protobuf:"bytes,4,opt,name=opt,proto3" json:"opt,omitempty"`
	Operator string `protobuf:"bytes,5,opt,name=operator,proto3" json:"operator,omitempty"`
}

func (x *GmCodeRequest) Reset() {
	*x = GmCodeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gm_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GmCodeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GmCodeRequest) ProtoMessage() {}

func (x *GmCodeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gm_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GmCodeRequest.ProtoReflect.Descriptor instead.
func (*GmCodeRequest) Descriptor() ([]byte, []int) {
	return file_gm_proto_rawDescGZIP(), []int{0}
}

func (x *GmCodeRequest) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *GmCodeRequest) GetGameId() string {
	if x != nil {
		return x.GameId
	}
	return ""
}

func (x *GmCodeRequest) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *GmCodeRequest) GetOpt() string {
	if x != nil {
		return x.Opt
	}
	return ""
}

func (x *GmCodeRequest) GetOperator() string {
	if x != nil {
		return x.Operator
	}
	return ""
}

type GmCodeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg  string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *GmCodeResponse) Reset() {
	*x = GmCodeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gm_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GmCodeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GmCodeResponse) ProtoMessage() {}

func (x *GmCodeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gm_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GmCodeResponse.ProtoReflect.Descriptor instead.
func (*GmCodeResponse) Descriptor() ([]byte, []int) {
	return file_gm_proto_rawDescGZIP(), []int{1}
}

func (x *GmCodeResponse) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *GmCodeResponse) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

var File_gm_proto protoreflect.FileDescriptor

var file_gm_proto_rawDesc = []byte{
	0x0a, 0x08, 0x67, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x22, 0x7b,
	0x0a, 0x0d, 0x47, 0x6d, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63,
	0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x67, 0x61, 0x6d, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x67, 0x61, 0x6d, 0x65, 0x49, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x75,
	0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x10, 0x0a,
	0x03, 0x6f, 0x70, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6f, 0x70, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x22, 0x36, 0x0a, 0x0e, 0x47,
	0x6d, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6d, 0x73, 0x67, 0x2a, 0x71, 0x0a, 0x07, 0x47, 0x4d, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x12, 0x0f,
	0x0a, 0x0b, 0x47, 0x41, 0x4d, 0x45, 0x5f, 0x4e, 0x4f, 0x54, 0x49, 0x43, 0x45, 0x10, 0x00, 0x12,
	0x0d, 0x0a, 0x09, 0x4b, 0x49, 0x43, 0x4b, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x10, 0x01, 0x12, 0x0e,
	0x0a, 0x0a, 0x42, 0x4c, 0x41, 0x43, 0x4b, 0x5f, 0x4c, 0x49, 0x53, 0x54, 0x10, 0x02, 0x12, 0x10,
	0x0a, 0x0c, 0x55, 0x4e, 0x4c, 0x4f, 0x43, 0x4b, 0x5f, 0x42, 0x4c, 0x41, 0x43, 0x4b, 0x10, 0x03,
	0x12, 0x0f, 0x0a, 0x0b, 0x52, 0x45, 0x57, 0x41, 0x52, 0x44, 0x5f, 0x49, 0x54, 0x45, 0x4d, 0x10,
	0x04, 0x12, 0x13, 0x0a, 0x0f, 0x41, 0x4c, 0x4c, 0x5f, 0x47, 0x41, 0x4d, 0x45, 0x5f, 0x4e, 0x4f,
	0x54, 0x49, 0x43, 0x45, 0x10, 0x0a, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2e, 0x2f, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gm_proto_rawDescOnce sync.Once
	file_gm_proto_rawDescData = file_gm_proto_rawDesc
)

func file_gm_proto_rawDescGZIP() []byte {
	file_gm_proto_rawDescOnce.Do(func() {
		file_gm_proto_rawDescData = protoimpl.X.CompressGZIP(file_gm_proto_rawDescData)
	})
	return file_gm_proto_rawDescData
}

var file_gm_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_gm_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_gm_proto_goTypes = []any{
	(GM_CODE)(0),           // 0: pb.GM_CODE
	(*GmCodeRequest)(nil),  // 1: pb.GmCodeRequest
	(*GmCodeResponse)(nil), // 2: pb.GmCodeResponse
}
var file_gm_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_gm_proto_init() }
func file_gm_proto_init() {
	if File_gm_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gm_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*GmCodeRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gm_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*GmCodeResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_gm_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gm_proto_goTypes,
		DependencyIndexes: file_gm_proto_depIdxs,
		EnumInfos:         file_gm_proto_enumTypes,
		MessageInfos:      file_gm_proto_msgTypes,
	}.Build()
	File_gm_proto = out.File
	file_gm_proto_rawDesc = nil
	file_gm_proto_goTypes = nil
	file_gm_proto_depIdxs = nil
}
