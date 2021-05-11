// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.8
// source: pkg/rpc/admin/service.proto

package admin

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	database "gosship/pkg/database"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetAuthChallenge struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetAuthChallenge) Reset() {
	*x = GetAuthChallenge{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_admin_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAuthChallenge) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAuthChallenge) ProtoMessage() {}

func (x *GetAuthChallenge) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_admin_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAuthChallenge.ProtoReflect.Descriptor instead.
func (*GetAuthChallenge) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_admin_service_proto_rawDescGZIP(), []int{0}
}

type Auth struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Auth) Reset() {
	*x = Auth{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_admin_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Auth) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Auth) ProtoMessage() {}

func (x *Auth) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_admin_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Auth.ProtoReflect.Descriptor instead.
func (*Auth) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_admin_service_proto_rawDescGZIP(), []int{1}
}

type UpdateUserFingerprint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UpdateUserFingerprint) Reset() {
	*x = UpdateUserFingerprint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_admin_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateUserFingerprint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUserFingerprint) ProtoMessage() {}

func (x *UpdateUserFingerprint) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_admin_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUserFingerprint.ProtoReflect.Descriptor instead.
func (*UpdateUserFingerprint) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_admin_service_proto_rawDescGZIP(), []int{2}
}

type SendMessageToRoom struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SendMessageToRoom) Reset() {
	*x = SendMessageToRoom{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_admin_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendMessageToRoom) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMessageToRoom) ProtoMessage() {}

func (x *SendMessageToRoom) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_admin_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMessageToRoom.ProtoReflect.Descriptor instead.
func (*SendMessageToRoom) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_admin_service_proto_rawDescGZIP(), []int{3}
}

type DumpMessages struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DumpMessages) Reset() {
	*x = DumpMessages{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_admin_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DumpMessages) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DumpMessages) ProtoMessage() {}

func (x *DumpMessages) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_admin_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DumpMessages.ProtoReflect.Descriptor instead.
func (*DumpMessages) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_admin_service_proto_rawDescGZIP(), []int{4}
}

type GetAuthChallenge_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetAuthChallenge_Request) Reset() {
	*x = GetAuthChallenge_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_admin_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAuthChallenge_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAuthChallenge_Request) ProtoMessage() {}

func (x *GetAuthChallenge_Request) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_admin_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAuthChallenge_Request.ProtoReflect.Descriptor instead.
func (*GetAuthChallenge_Request) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_admin_service_proto_rawDescGZIP(), []int{0, 0}
}

type GetAuthChallenge_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChallengeId string `protobuf:"bytes,1,opt,name=ChallengeId,proto3" json:"ChallengeId,omitempty"`
	Challenge   []byte `protobuf:"bytes,2,opt,name=Challenge,proto3" json:"Challenge,omitempty"`
}

func (x *GetAuthChallenge_Response) Reset() {
	*x = GetAuthChallenge_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_admin_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAuthChallenge_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAuthChallenge_Response) ProtoMessage() {}

func (x *GetAuthChallenge_Response) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_admin_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAuthChallenge_Response.ProtoReflect.Descriptor instead.
func (*GetAuthChallenge_Response) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_admin_service_proto_rawDescGZIP(), []int{0, 1}
}

func (x *GetAuthChallenge_Response) GetChallengeId() string {
	if x != nil {
		return x.ChallengeId
	}
	return ""
}

func (x *GetAuthChallenge_Response) GetChallenge() []byte {
	if x != nil {
		return x.Challenge
	}
	return nil
}

type Auth_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChallengeId string `protobuf:"bytes,1,opt,name=ChallengeId,proto3" json:"ChallengeId,omitempty"`
	Signature   []byte `protobuf:"bytes,2,opt,name=Signature,proto3" json:"Signature,omitempty"`
}

func (x *Auth_Request) Reset() {
	*x = Auth_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_admin_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Auth_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Auth_Request) ProtoMessage() {}

func (x *Auth_Request) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_admin_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Auth_Request.ProtoReflect.Descriptor instead.
func (*Auth_Request) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_admin_service_proto_rawDescGZIP(), []int{1, 0}
}

func (x *Auth_Request) GetChallengeId() string {
	if x != nil {
		return x.ChallengeId
	}
	return ""
}

func (x *Auth_Request) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type Auth_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionToken string `protobuf:"bytes,1,opt,name=SessionToken,proto3" json:"SessionToken,omitempty"`
}

func (x *Auth_Response) Reset() {
	*x = Auth_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_admin_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Auth_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Auth_Response) ProtoMessage() {}

func (x *Auth_Response) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_admin_service_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Auth_Response.ProtoReflect.Descriptor instead.
func (*Auth_Response) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_admin_service_proto_rawDescGZIP(), []int{1, 1}
}

func (x *Auth_Response) GetSessionToken() string {
	if x != nil {
		return x.SessionToken
	}
	return ""
}

type UpdateUserFingerprint_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionToken string `protobuf:"bytes,1,opt,name=SessionToken,proto3" json:"SessionToken,omitempty"`
	Username     string `protobuf:"bytes,2,opt,name=Username,proto3" json:"Username,omitempty"`
	Fingerprint  string `protobuf:"bytes,3,opt,name=Fingerprint,proto3" json:"Fingerprint,omitempty"`
}

func (x *UpdateUserFingerprint_Request) Reset() {
	*x = UpdateUserFingerprint_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_admin_service_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateUserFingerprint_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUserFingerprint_Request) ProtoMessage() {}

func (x *UpdateUserFingerprint_Request) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_admin_service_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUserFingerprint_Request.ProtoReflect.Descriptor instead.
func (*UpdateUserFingerprint_Request) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_admin_service_proto_rawDescGZIP(), []int{2, 0}
}

func (x *UpdateUserFingerprint_Request) GetSessionToken() string {
	if x != nil {
		return x.SessionToken
	}
	return ""
}

func (x *UpdateUserFingerprint_Request) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *UpdateUserFingerprint_Request) GetFingerprint() string {
	if x != nil {
		return x.Fingerprint
	}
	return ""
}

type UpdateUserFingerprint_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UpdateUserFingerprint_Response) Reset() {
	*x = UpdateUserFingerprint_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_admin_service_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateUserFingerprint_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUserFingerprint_Response) ProtoMessage() {}

func (x *UpdateUserFingerprint_Response) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_admin_service_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUserFingerprint_Response.ProtoReflect.Descriptor instead.
func (*UpdateUserFingerprint_Response) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_admin_service_proto_rawDescGZIP(), []int{2, 1}
}

type SendMessageToRoom_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionToken string `protobuf:"bytes,1,opt,name=SessionToken,proto3" json:"SessionToken,omitempty"`
	Room         string `protobuf:"bytes,2,opt,name=Room,proto3" json:"Room,omitempty"`
	Message      string `protobuf:"bytes,3,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (x *SendMessageToRoom_Request) Reset() {
	*x = SendMessageToRoom_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_admin_service_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendMessageToRoom_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMessageToRoom_Request) ProtoMessage() {}

func (x *SendMessageToRoom_Request) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_admin_service_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMessageToRoom_Request.ProtoReflect.Descriptor instead.
func (*SendMessageToRoom_Request) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_admin_service_proto_rawDescGZIP(), []int{3, 0}
}

func (x *SendMessageToRoom_Request) GetSessionToken() string {
	if x != nil {
		return x.SessionToken
	}
	return ""
}

func (x *SendMessageToRoom_Request) GetRoom() string {
	if x != nil {
		return x.Room
	}
	return ""
}

func (x *SendMessageToRoom_Request) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type SendMessageToRoom_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SendMessageToRoom_Response) Reset() {
	*x = SendMessageToRoom_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_admin_service_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendMessageToRoom_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMessageToRoom_Response) ProtoMessage() {}

func (x *SendMessageToRoom_Response) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_admin_service_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMessageToRoom_Response.ProtoReflect.Descriptor instead.
func (*SendMessageToRoom_Response) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_admin_service_proto_rawDescGZIP(), []int{3, 1}
}

type DumpMessages_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionToken string `protobuf:"bytes,1,opt,name=SessionToken,proto3" json:"SessionToken,omitempty"`
}

func (x *DumpMessages_Request) Reset() {
	*x = DumpMessages_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_admin_service_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DumpMessages_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DumpMessages_Request) ProtoMessage() {}

func (x *DumpMessages_Request) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_admin_service_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DumpMessages_Request.ProtoReflect.Descriptor instead.
func (*DumpMessages_Request) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_admin_service_proto_rawDescGZIP(), []int{4, 0}
}

func (x *DumpMessages_Request) GetSessionToken() string {
	if x != nil {
		return x.SessionToken
	}
	return ""
}

type DumpMessages_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message *database.MessageEntry `protobuf:"bytes,1,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (x *DumpMessages_Response) Reset() {
	*x = DumpMessages_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_admin_service_proto_msgTypes[14]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DumpMessages_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DumpMessages_Response) ProtoMessage() {}

func (x *DumpMessages_Response) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_admin_service_proto_msgTypes[14]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DumpMessages_Response.ProtoReflect.Descriptor instead.
func (*DumpMessages_Response) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_admin_service_proto_rawDescGZIP(), []int{4, 1}
}

func (x *DumpMessages_Response) GetMessage() *database.MessageEntry {
	if x != nil {
		return x.Message
	}
	return nil
}

var File_pkg_rpc_admin_service_proto protoreflect.FileDescriptor

var file_pkg_rpc_admin_service_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x6b, 0x67, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x23, 0x70,
	0x6b, 0x67, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x64, 0x61, 0x74, 0x61,
	0x62, 0x61, 0x73, 0x65, 0x5f, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x69, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x41, 0x75, 0x74, 0x68, 0x43, 0x68, 0x61,
	0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x1a, 0x09, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x4a, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x20, 0x0a,
	0x0b, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x49, 0x64, 0x12,
	0x1c, 0x0a, 0x09, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x09, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x22, 0x81, 0x01,
	0x0a, 0x04, 0x41, 0x75, 0x74, 0x68, 0x1a, 0x49, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x20, 0x0a, 0x0b, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x49, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67,
	0x65, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x1a, 0x2e, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a,
	0x0c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x22, 0x90, 0x01, 0x0a, 0x15, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72,
	0x46, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x1a, 0x6b, 0x0a, 0x07, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x55, 0x73,
	0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x55, 0x73,
	0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x46, 0x69, 0x6e, 0x67, 0x65, 0x72,
	0x70, 0x72, 0x69, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x46, 0x69, 0x6e,
	0x67, 0x65, 0x72, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x1a, 0x0a, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x7c, 0x0a, 0x11, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x54, 0x6f, 0x52, 0x6f, 0x6f, 0x6d, 0x1a, 0x5b, 0x0a, 0x07, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x53, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x52, 0x6f, 0x6f, 0x6d,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x52, 0x6f, 0x6f, 0x6d, 0x12, 0x18, 0x0a, 0x07,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x0a, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x72, 0x0a, 0x0c, 0x44, 0x75, 0x6d, 0x70, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x73, 0x1a, 0x2d, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a,
	0x0c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x1a, 0x33, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x27, 0x0a,
	0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d,
	0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0xe9, 0x02, 0x0a, 0x0c, 0x41, 0x64, 0x6d, 0x69, 0x6e,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x49, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x41, 0x75,
	0x74, 0x68, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x12, 0x19, 0x2e, 0x47, 0x65,
	0x74, 0x41, 0x75, 0x74, 0x68, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x2e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x75, 0x74, 0x68,
	0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x25, 0x0a, 0x04, 0x41, 0x75, 0x74, 0x68, 0x12, 0x0d, 0x2e, 0x41, 0x75, 0x74,
	0x68, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x41, 0x75, 0x74, 0x68,
	0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x58, 0x0a, 0x15, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x46, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x72, 0x69,
	0x6e, 0x74, 0x12, 0x1e, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x46,
	0x69, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x46,
	0x69, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x4c, 0x0a, 0x11, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x54, 0x6f, 0x52, 0x6f, 0x6f, 0x6d, 0x12, 0x1a, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x52, 0x6f, 0x6f, 0x6d, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x54, 0x6f, 0x52, 0x6f, 0x6f, 0x6d, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x3f, 0x0a, 0x0c, 0x44, 0x75, 0x6d, 0x70, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x73, 0x12, 0x15, 0x2e, 0x44, 0x75, 0x6d, 0x70, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x44, 0x75, 0x6d, 0x70, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x30, 0x01, 0x42, 0x0f, 0x5a, 0x0d, 0x70, 0x6b, 0x67, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_rpc_admin_service_proto_rawDescOnce sync.Once
	file_pkg_rpc_admin_service_proto_rawDescData = file_pkg_rpc_admin_service_proto_rawDesc
)

func file_pkg_rpc_admin_service_proto_rawDescGZIP() []byte {
	file_pkg_rpc_admin_service_proto_rawDescOnce.Do(func() {
		file_pkg_rpc_admin_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_rpc_admin_service_proto_rawDescData)
	})
	return file_pkg_rpc_admin_service_proto_rawDescData
}

var file_pkg_rpc_admin_service_proto_msgTypes = make([]protoimpl.MessageInfo, 15)
var file_pkg_rpc_admin_service_proto_goTypes = []interface{}{
	(*GetAuthChallenge)(nil),               // 0: GetAuthChallenge
	(*Auth)(nil),                           // 1: Auth
	(*UpdateUserFingerprint)(nil),          // 2: UpdateUserFingerprint
	(*SendMessageToRoom)(nil),              // 3: SendMessageToRoom
	(*DumpMessages)(nil),                   // 4: DumpMessages
	(*GetAuthChallenge_Request)(nil),       // 5: GetAuthChallenge.Request
	(*GetAuthChallenge_Response)(nil),      // 6: GetAuthChallenge.Response
	(*Auth_Request)(nil),                   // 7: Auth.Request
	(*Auth_Response)(nil),                  // 8: Auth.Response
	(*UpdateUserFingerprint_Request)(nil),  // 9: UpdateUserFingerprint.Request
	(*UpdateUserFingerprint_Response)(nil), // 10: UpdateUserFingerprint.Response
	(*SendMessageToRoom_Request)(nil),      // 11: SendMessageToRoom.Request
	(*SendMessageToRoom_Response)(nil),     // 12: SendMessageToRoom.Response
	(*DumpMessages_Request)(nil),           // 13: DumpMessages.Request
	(*DumpMessages_Response)(nil),          // 14: DumpMessages.Response
	(*database.MessageEntry)(nil),          // 15: MessageEntry
}
var file_pkg_rpc_admin_service_proto_depIdxs = []int32{
	15, // 0: DumpMessages.Response.Message:type_name -> MessageEntry
	5,  // 1: AdminService.GetAuthChallenge:input_type -> GetAuthChallenge.Request
	7,  // 2: AdminService.Auth:input_type -> Auth.Request
	9,  // 3: AdminService.UpdateUserFingerprint:input_type -> UpdateUserFingerprint.Request
	11, // 4: AdminService.SendMessageToRoom:input_type -> SendMessageToRoom.Request
	13, // 5: AdminService.DumpMessages:input_type -> DumpMessages.Request
	6,  // 6: AdminService.GetAuthChallenge:output_type -> GetAuthChallenge.Response
	8,  // 7: AdminService.Auth:output_type -> Auth.Response
	10, // 8: AdminService.UpdateUserFingerprint:output_type -> UpdateUserFingerprint.Response
	12, // 9: AdminService.SendMessageToRoom:output_type -> SendMessageToRoom.Response
	14, // 10: AdminService.DumpMessages:output_type -> DumpMessages.Response
	6,  // [6:11] is the sub-list for method output_type
	1,  // [1:6] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_pkg_rpc_admin_service_proto_init() }
func file_pkg_rpc_admin_service_proto_init() {
	if File_pkg_rpc_admin_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_rpc_admin_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAuthChallenge); i {
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
		file_pkg_rpc_admin_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Auth); i {
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
		file_pkg_rpc_admin_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateUserFingerprint); i {
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
		file_pkg_rpc_admin_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendMessageToRoom); i {
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
		file_pkg_rpc_admin_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DumpMessages); i {
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
		file_pkg_rpc_admin_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAuthChallenge_Request); i {
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
		file_pkg_rpc_admin_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAuthChallenge_Response); i {
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
		file_pkg_rpc_admin_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Auth_Request); i {
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
		file_pkg_rpc_admin_service_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Auth_Response); i {
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
		file_pkg_rpc_admin_service_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateUserFingerprint_Request); i {
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
		file_pkg_rpc_admin_service_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateUserFingerprint_Response); i {
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
		file_pkg_rpc_admin_service_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendMessageToRoom_Request); i {
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
		file_pkg_rpc_admin_service_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendMessageToRoom_Response); i {
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
		file_pkg_rpc_admin_service_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DumpMessages_Request); i {
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
		file_pkg_rpc_admin_service_proto_msgTypes[14].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DumpMessages_Response); i {
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
			RawDescriptor: file_pkg_rpc_admin_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   15,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_rpc_admin_service_proto_goTypes,
		DependencyIndexes: file_pkg_rpc_admin_service_proto_depIdxs,
		MessageInfos:      file_pkg_rpc_admin_service_proto_msgTypes,
	}.Build()
	File_pkg_rpc_admin_service_proto = out.File
	file_pkg_rpc_admin_service_proto_rawDesc = nil
	file_pkg_rpc_admin_service_proto_goTypes = nil
	file_pkg_rpc_admin_service_proto_depIdxs = nil
}
