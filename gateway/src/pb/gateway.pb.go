// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: gateway.proto

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

// 心跳, gateway 通知游戏服，每2秒发送
type PingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp int64 `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"` // 服务器时间戳
}

func (x *PingRequest) Reset() {
	*x = PingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gateway_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingRequest) ProtoMessage() {}

func (x *PingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gateway_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingRequest.ProtoReflect.Descriptor instead.
func (*PingRequest) Descriptor() ([]byte, []int) {
	return file_gateway_proto_rawDescGZIP(), []int{0}
}

func (x *PingRequest) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

type CommonHead struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GameId    string `protobuf:"bytes,1,opt,name=gameId,proto3" json:"gameId,omitempty"`
	RoomId    string `protobuf:"bytes,2,opt,name=roomId,proto3" json:"roomId,omitempty"`
	Uid       string `protobuf:"bytes,3,opt,name=uid,proto3" json:"uid,omitempty"`
	Pid       string `protobuf:"bytes,4,opt,name=pid,proto3" json:"pid,omitempty"`
	Sn        int64  `protobuf:"varint,5,opt,name=sn,proto3" json:"sn,omitempty"`
	Timestamp int64  `protobuf:"varint,6,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	ProtoName string `protobuf:"bytes,7,opt,name=protoName,proto3" json:"protoName,omitempty"`
	HostIp    string `protobuf:"bytes,8,opt,name=hostIp,proto3" json:"hostIp,omitempty"`
}

func (x *CommonHead) Reset() {
	*x = CommonHead{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gateway_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommonHead) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommonHead) ProtoMessage() {}

func (x *CommonHead) ProtoReflect() protoreflect.Message {
	mi := &file_gateway_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommonHead.ProtoReflect.Descriptor instead.
func (*CommonHead) Descriptor() ([]byte, []int) {
	return file_gateway_proto_rawDescGZIP(), []int{1}
}

func (x *CommonHead) GetGameId() string {
	if x != nil {
		return x.GameId
	}
	return ""
}

func (x *CommonHead) GetRoomId() string {
	if x != nil {
		return x.RoomId
	}
	return ""
}

func (x *CommonHead) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *CommonHead) GetPid() string {
	if x != nil {
		return x.Pid
	}
	return ""
}

func (x *CommonHead) GetSn() int64 {
	if x != nil {
		return x.Sn
	}
	return 0
}

func (x *CommonHead) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *CommonHead) GetProtoName() string {
	if x != nil {
		return x.ProtoName
	}
	return ""
}

func (x *CommonHead) GetHostIp() string {
	if x != nil {
		return x.HostIp
	}
	return ""
}

type GameCommonRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Head *CommonHead `protobuf:"bytes,1,opt,name=head,proto3" json:"head,omitempty"`
	Data []byte      `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *GameCommonRequest) Reset() {
	*x = GameCommonRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gateway_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GameCommonRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameCommonRequest) ProtoMessage() {}

func (x *GameCommonRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gateway_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GameCommonRequest.ProtoReflect.Descriptor instead.
func (*GameCommonRequest) Descriptor() ([]byte, []int) {
	return file_gateway_proto_rawDescGZIP(), []int{2}
}

func (x *GameCommonRequest) GetHead() *CommonHead {
	if x != nil {
		return x.Head
	}
	return nil
}

func (x *GameCommonRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type GameCommonResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code int32       `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg  string      `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Head *CommonHead `protobuf:"bytes,3,opt,name=head,proto3" json:"head,omitempty"`
	Data []byte      `protobuf:"bytes,6,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *GameCommonResponse) Reset() {
	*x = GameCommonResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gateway_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GameCommonResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameCommonResponse) ProtoMessage() {}

func (x *GameCommonResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gateway_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GameCommonResponse.ProtoReflect.Descriptor instead.
func (*GameCommonResponse) Descriptor() ([]byte, []int) {
	return file_gateway_proto_rawDescGZIP(), []int{3}
}

func (x *GameCommonResponse) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *GameCommonResponse) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *GameCommonResponse) GetHead() *CommonHead {
	if x != nil {
		return x.Head
	}
	return nil
}

func (x *GameCommonResponse) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type PushHead struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GameId    string `protobuf:"bytes,1,opt,name=gameId,proto3" json:"gameId,omitempty"`
	RoomId    string `protobuf:"bytes,2,opt,name=roomId,proto3" json:"roomId,omitempty"`
	Uid       string `protobuf:"bytes,3,opt,name=uid,proto3" json:"uid,omitempty"`
	Pid       string `protobuf:"bytes,4,opt,name=pid,proto3" json:"pid,omitempty"`
	Sn        int64  `protobuf:"varint,5,opt,name=sn,proto3" json:"sn,omitempty"`
	Timestamp int64  `protobuf:"varint,6,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	ProtoName string `protobuf:"bytes,7,opt,name=protoName,proto3" json:"protoName,omitempty"`
}

func (x *PushHead) Reset() {
	*x = PushHead{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gateway_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushHead) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushHead) ProtoMessage() {}

func (x *PushHead) ProtoReflect() protoreflect.Message {
	mi := &file_gateway_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushHead.ProtoReflect.Descriptor instead.
func (*PushHead) Descriptor() ([]byte, []int) {
	return file_gateway_proto_rawDescGZIP(), []int{4}
}

func (x *PushHead) GetGameId() string {
	if x != nil {
		return x.GameId
	}
	return ""
}

func (x *PushHead) GetRoomId() string {
	if x != nil {
		return x.RoomId
	}
	return ""
}

func (x *PushHead) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *PushHead) GetPid() string {
	if x != nil {
		return x.Pid
	}
	return ""
}

func (x *PushHead) GetSn() int64 {
	if x != nil {
		return x.Sn
	}
	return 0
}

func (x *PushHead) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *PushHead) GetProtoName() string {
	if x != nil {
		return x.ProtoName
	}
	return ""
}

type GamePushMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Head *PushHead `protobuf:"bytes,1,opt,name=head,proto3" json:"head,omitempty"`
	Data []byte    `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *GamePushMessage) Reset() {
	*x = GamePushMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gateway_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GamePushMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GamePushMessage) ProtoMessage() {}

func (x *GamePushMessage) ProtoReflect() protoreflect.Message {
	mi := &file_gateway_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GamePushMessage.ProtoReflect.Descriptor instead.
func (*GamePushMessage) Descriptor() ([]byte, []int) {
	return file_gateway_proto_rawDescGZIP(), []int{5}
}

func (x *GamePushMessage) GetHead() *PushHead {
	if x != nil {
		return x.Head
	}
	return nil
}

func (x *GamePushMessage) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

// 进入大厅
type LoginHallRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *LoginHallRequest) Reset() {
	*x = LoginHallRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gateway_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginHallRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginHallRequest) ProtoMessage() {}

func (x *LoginHallRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gateway_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginHallRequest.ProtoReflect.Descriptor instead.
func (*LoginHallRequest) Descriptor() ([]byte, []int) {
	return file_gateway_proto_rawDescGZIP(), []int{6}
}

// 进入大厅应答
type LoginHallResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code     int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg      string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	JumpGame bool   `protobuf:"varint,3,opt,name=jumpGame,proto3" json:"jumpGame,omitempty"`
	RoomId   string `protobuf:"bytes,4,opt,name=roomId,proto3" json:"roomId,omitempty"`
	Opt      string `protobuf:"bytes,5,opt,name=opt,proto3" json:"opt,omitempty"`
}

func (x *LoginHallResponse) Reset() {
	*x = LoginHallResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gateway_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginHallResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginHallResponse) ProtoMessage() {}

func (x *LoginHallResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gateway_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginHallResponse.ProtoReflect.Descriptor instead.
func (*LoginHallResponse) Descriptor() ([]byte, []int) {
	return file_gateway_proto_rawDescGZIP(), []int{7}
}

func (x *LoginHallResponse) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *LoginHallResponse) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *LoginHallResponse) GetJumpGame() bool {
	if x != nil {
		return x.JumpGame
	}
	return false
}

func (x *LoginHallResponse) GetRoomId() string {
	if x != nil {
		return x.RoomId
	}
	return ""
}

func (x *LoginHallResponse) GetOpt() string {
	if x != nil {
		return x.Opt
	}
	return ""
}

var File_gateway_proto protoreflect.FileDescriptor

var file_gateway_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x02, 0x70, 0x62, 0x22, 0x2b, 0x0a, 0x0b, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x22, 0xc4, 0x01, 0x0a, 0x0a, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x48, 0x65, 0x61, 0x64, 0x12,
	0x16, 0x0a, 0x06, 0x67, 0x61, 0x6d, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x67, 0x61, 0x6d, 0x65, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x6f, 0x6f, 0x6d, 0x49,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x6f, 0x6f, 0x6d, 0x49, 0x64, 0x12,
	0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x69,
	0x64, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x70, 0x69, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x73, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x02, 0x73, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x68, 0x6f, 0x73, 0x74, 0x49, 0x70, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x68, 0x6f, 0x73, 0x74, 0x49, 0x70, 0x22, 0x4b, 0x0a, 0x11, 0x47, 0x61, 0x6d, 0x65, 0x43,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x04,
	0x68, 0x65, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x62, 0x2e,
	0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x48, 0x65, 0x61, 0x64, 0x52, 0x04, 0x68, 0x65, 0x61, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x22, 0x72, 0x0a, 0x12, 0x47, 0x61, 0x6d, 0x65, 0x43, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67,
	0x12, 0x22, 0x0a, 0x04, 0x68, 0x65, 0x61, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e,
	0x2e, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x48, 0x65, 0x61, 0x64, 0x52, 0x04,
	0x68, 0x65, 0x61, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xaa, 0x01, 0x0a, 0x08, 0x50, 0x75, 0x73,
	0x68, 0x48, 0x65, 0x61, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x67, 0x61, 0x6d, 0x65, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x67, 0x61, 0x6d, 0x65, 0x49, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x72, 0x6f, 0x6f, 0x6d, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72,
	0x6f, 0x6f, 0x6d, 0x49, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x69, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x70, 0x69, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x73, 0x6e, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x73, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x47, 0x0a, 0x0f, 0x47, 0x61, 0x6d, 0x65, 0x50, 0x75, 0x73,
	0x68, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x20, 0x0a, 0x04, 0x68, 0x65, 0x61, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x62, 0x2e, 0x50, 0x75, 0x73, 0x68,
	0x48, 0x65, 0x61, 0x64, 0x52, 0x04, 0x68, 0x65, 0x61, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x12,
	0x0a, 0x10, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x48, 0x61, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x22, 0x7f, 0x0a, 0x11, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x48, 0x61, 0x6c, 0x6c, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d,
	0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x1a, 0x0a,
	0x08, 0x6a, 0x75, 0x6d, 0x70, 0x47, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x08, 0x6a, 0x75, 0x6d, 0x70, 0x47, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x6f, 0x6f,
	0x6d, 0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x6f, 0x6f, 0x6d, 0x49,
	0x64, 0x12, 0x10, 0x0a, 0x03, 0x6f, 0x70, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6f, 0x70, 0x74, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gateway_proto_rawDescOnce sync.Once
	file_gateway_proto_rawDescData = file_gateway_proto_rawDesc
)

func file_gateway_proto_rawDescGZIP() []byte {
	file_gateway_proto_rawDescOnce.Do(func() {
		file_gateway_proto_rawDescData = protoimpl.X.CompressGZIP(file_gateway_proto_rawDescData)
	})
	return file_gateway_proto_rawDescData
}

var file_gateway_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_gateway_proto_goTypes = []any{
	(*PingRequest)(nil),        // 0: pb.PingRequest
	(*CommonHead)(nil),         // 1: pb.CommonHead
	(*GameCommonRequest)(nil),  // 2: pb.GameCommonRequest
	(*GameCommonResponse)(nil), // 3: pb.GameCommonResponse
	(*PushHead)(nil),           // 4: pb.PushHead
	(*GamePushMessage)(nil),    // 5: pb.GamePushMessage
	(*LoginHallRequest)(nil),   // 6: pb.LoginHallRequest
	(*LoginHallResponse)(nil),  // 7: pb.LoginHallResponse
}
var file_gateway_proto_depIdxs = []int32{
	1, // 0: pb.GameCommonRequest.head:type_name -> pb.CommonHead
	1, // 1: pb.GameCommonResponse.head:type_name -> pb.CommonHead
	4, // 2: pb.GamePushMessage.head:type_name -> pb.PushHead
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_gateway_proto_init() }
func file_gateway_proto_init() {
	if File_gateway_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gateway_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*PingRequest); i {
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
		file_gateway_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*CommonHead); i {
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
		file_gateway_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*GameCommonRequest); i {
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
		file_gateway_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*GameCommonResponse); i {
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
		file_gateway_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*PushHead); i {
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
		file_gateway_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*GamePushMessage); i {
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
		file_gateway_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*LoginHallRequest); i {
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
		file_gateway_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*LoginHallResponse); i {
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
			RawDescriptor: file_gateway_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gateway_proto_goTypes,
		DependencyIndexes: file_gateway_proto_depIdxs,
		MessageInfos:      file_gateway_proto_msgTypes,
	}.Build()
	File_gateway_proto = out.File
	file_gateway_proto_rawDesc = nil
	file_gateway_proto_goTypes = nil
	file_gateway_proto_depIdxs = nil
}
