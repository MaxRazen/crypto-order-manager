// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: protofiles/ordermanager.proto

package ordergrpc

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

// Enum for action type (sell or buy)
type ActionType int32

const (
	ActionType_UNKNOWN_ACTION ActionType = 0
	ActionType_SELL           ActionType = 1
	ActionType_BUY            ActionType = 2
)

// Enum value maps for ActionType.
var (
	ActionType_name = map[int32]string{
		0: "UNKNOWN_ACTION",
		1: "SELL",
		2: "BUY",
	}
	ActionType_value = map[string]int32{
		"UNKNOWN_ACTION": 0,
		"SELL":           1,
		"BUY":            2,
	}
)

func (x ActionType) Enum() *ActionType {
	p := new(ActionType)
	*p = x
	return p
}

func (x ActionType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ActionType) Descriptor() protoreflect.EnumDescriptor {
	return file_protofiles_ordermanager_proto_enumTypes[0].Descriptor()
}

func (ActionType) Type() protoreflect.EnumType {
	return &file_protofiles_ordermanager_proto_enumTypes[0]
}

func (x ActionType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ActionType.Descriptor instead.
func (ActionType) EnumDescriptor() ([]byte, []int) {
	return file_protofiles_ordermanager_proto_rawDescGZIP(), []int{0}
}

// Enum for engine behavior (market or limit)
type Behavior int32

const (
	Behavior_UNKNOWN_BEHAVIOR Behavior = 0
	Behavior_MARKET           Behavior = 1
	Behavior_LIMIT            Behavior = 2
)

// Enum value maps for Behavior.
var (
	Behavior_name = map[int32]string{
		0: "UNKNOWN_BEHAVIOR",
		1: "MARKET",
		2: "LIMIT",
	}
	Behavior_value = map[string]int32{
		"UNKNOWN_BEHAVIOR": 0,
		"MARKET":           1,
		"LIMIT":            2,
	}
)

func (x Behavior) Enum() *Behavior {
	p := new(Behavior)
	*p = x
	return p
}

func (x Behavior) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Behavior) Descriptor() protoreflect.EnumDescriptor {
	return file_protofiles_ordermanager_proto_enumTypes[1].Descriptor()
}

func (Behavior) Type() protoreflect.EnumType {
	return &file_protofiles_ordermanager_proto_enumTypes[1]
}

func (x Behavior) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Behavior.Descriptor instead.
func (Behavior) EnumDescriptor() ([]byte, []int) {
	return file_protofiles_ordermanager_proto_rawDescGZIP(), []int{1}
}

// Enum for amount type (scalar or percent)
type QuantityType int32

const (
	QuantityType_UNKNOWN_QUANTITY_TYPE QuantityType = 0
	QuantityType_FIXED                 QuantityType = 1
	QuantityType_PERCENT               QuantityType = 2
)

// Enum value maps for QuantityType.
var (
	QuantityType_name = map[int32]string{
		0: "UNKNOWN_QUANTITY_TYPE",
		1: "FIXED",
		2: "PERCENT",
	}
	QuantityType_value = map[string]int32{
		"UNKNOWN_QUANTITY_TYPE": 0,
		"FIXED":                 1,
		"PERCENT":               2,
	}
)

func (x QuantityType) Enum() *QuantityType {
	p := new(QuantityType)
	*p = x
	return p
}

func (x QuantityType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (QuantityType) Descriptor() protoreflect.EnumDescriptor {
	return file_protofiles_ordermanager_proto_enumTypes[2].Descriptor()
}

func (QuantityType) Type() protoreflect.EnumType {
	return &file_protofiles_ordermanager_proto_enumTypes[2]
}

func (x QuantityType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use QuantityType.Descriptor instead.
func (QuantityType) EnumDescriptor() ([]byte, []int) {
	return file_protofiles_ordermanager_proto_rawDescGZIP(), []int{2}
}

// Enum for deadline type (time)
type DeadlineType int32

const (
	DeadlineType_UNKNOWN_DEADLINE_TYPE DeadlineType = 0
	DeadlineType_TIME                  DeadlineType = 1
)

// Enum value maps for DeadlineType.
var (
	DeadlineType_name = map[int32]string{
		0: "UNKNOWN_DEADLINE_TYPE",
		1: "TIME",
	}
	DeadlineType_value = map[string]int32{
		"UNKNOWN_DEADLINE_TYPE": 0,
		"TIME":                  1,
	}
)

func (x DeadlineType) Enum() *DeadlineType {
	p := new(DeadlineType)
	*p = x
	return p
}

func (x DeadlineType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DeadlineType) Descriptor() protoreflect.EnumDescriptor {
	return file_protofiles_ordermanager_proto_enumTypes[3].Descriptor()
}

func (DeadlineType) Type() protoreflect.EnumType {
	return &file_protofiles_ordermanager_proto_enumTypes[3]
}

func (x DeadlineType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DeadlineType.Descriptor instead.
func (DeadlineType) EnumDescriptor() ([]byte, []int) {
	return file_protofiles_ordermanager_proto_rawDescGZIP(), []int{3}
}

// Enum for deadline action (cancel or sellByMarket)
type DeadlineAction int32

const (
	DeadlineAction_UNKNOWN_DEADLINE_ACTION DeadlineAction = 0
	DeadlineAction_CANCEL                  DeadlineAction = 1
	DeadlineAction_SELL_BY_MARKET          DeadlineAction = 2
)

// Enum value maps for DeadlineAction.
var (
	DeadlineAction_name = map[int32]string{
		0: "UNKNOWN_DEADLINE_ACTION",
		1: "CANCEL",
		2: "SELL_BY_MARKET",
	}
	DeadlineAction_value = map[string]int32{
		"UNKNOWN_DEADLINE_ACTION": 0,
		"CANCEL":                  1,
		"SELL_BY_MARKET":          2,
	}
)

func (x DeadlineAction) Enum() *DeadlineAction {
	p := new(DeadlineAction)
	*p = x
	return p
}

func (x DeadlineAction) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DeadlineAction) Descriptor() protoreflect.EnumDescriptor {
	return file_protofiles_ordermanager_proto_enumTypes[4].Descriptor()
}

func (DeadlineAction) Type() protoreflect.EnumType {
	return &file_protofiles_ordermanager_proto_enumTypes[4]
}

func (x DeadlineAction) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DeadlineAction.Descriptor instead.
func (DeadlineAction) EnumDescriptor() ([]byte, []int) {
	return file_protofiles_ordermanager_proto_rawDescGZIP(), []int{4}
}

// Request message for creating an order
type CreateOrderRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pair      string      `protobuf:"bytes,1,opt,name=pair,proto3" json:"pair,omitempty"`
	Market    string      `protobuf:"bytes,2,opt,name=market,proto3" json:"market,omitempty"`
	Action    ActionType  `protobuf:"varint,3,opt,name=action,proto3,enum=ordermanager.ActionType" json:"action,omitempty"`
	Behavior  Behavior    `protobuf:"varint,4,opt,name=behavior,proto3,enum=ordermanager.Behavior" json:"behavior,omitempty"`
	Price     string      `protobuf:"bytes,5,opt,name=price,proto3" json:"price,omitempty"`
	Quantity  *Quantity   `protobuf:"bytes,6,opt,name=quantity,proto3" json:"quantity,omitempty"`
	Deadlines []*Deadline `protobuf:"bytes,7,rep,name=deadlines,proto3" json:"deadlines,omitempty"`
}

func (x *CreateOrderRequest) Reset() {
	*x = CreateOrderRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protofiles_ordermanager_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateOrderRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateOrderRequest) ProtoMessage() {}

func (x *CreateOrderRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protofiles_ordermanager_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateOrderRequest.ProtoReflect.Descriptor instead.
func (*CreateOrderRequest) Descriptor() ([]byte, []int) {
	return file_protofiles_ordermanager_proto_rawDescGZIP(), []int{0}
}

func (x *CreateOrderRequest) GetPair() string {
	if x != nil {
		return x.Pair
	}
	return ""
}

func (x *CreateOrderRequest) GetMarket() string {
	if x != nil {
		return x.Market
	}
	return ""
}

func (x *CreateOrderRequest) GetAction() ActionType {
	if x != nil {
		return x.Action
	}
	return ActionType_UNKNOWN_ACTION
}

func (x *CreateOrderRequest) GetBehavior() Behavior {
	if x != nil {
		return x.Behavior
	}
	return Behavior_UNKNOWN_BEHAVIOR
}

func (x *CreateOrderRequest) GetPrice() string {
	if x != nil {
		return x.Price
	}
	return ""
}

func (x *CreateOrderRequest) GetQuantity() *Quantity {
	if x != nil {
		return x.Quantity
	}
	return nil
}

func (x *CreateOrderRequest) GetDeadlines() []*Deadline {
	if x != nil {
		return x.Deadlines
	}
	return nil
}

// Response message for creating an order
type CreateOrderResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *CreateOrderResponse) Reset() {
	*x = CreateOrderResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protofiles_ordermanager_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateOrderResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateOrderResponse) ProtoMessage() {}

func (x *CreateOrderResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protofiles_ordermanager_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateOrderResponse.ProtoReflect.Descriptor instead.
func (*CreateOrderResponse) Descriptor() ([]byte, []int) {
	return file_protofiles_ordermanager_proto_rawDescGZIP(), []int{1}
}

func (x *CreateOrderResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *CreateOrderResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// Amount message to specify type and value
type Quantity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type  QuantityType `protobuf:"varint,1,opt,name=type,proto3,enum=ordermanager.QuantityType" json:"type,omitempty"`
	Value string       `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Quantity) Reset() {
	*x = Quantity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protofiles_ordermanager_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Quantity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Quantity) ProtoMessage() {}

func (x *Quantity) ProtoReflect() protoreflect.Message {
	mi := &file_protofiles_ordermanager_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Quantity.ProtoReflect.Descriptor instead.
func (*Quantity) Descriptor() ([]byte, []int) {
	return file_protofiles_ordermanager_proto_rawDescGZIP(), []int{2}
}

func (x *Quantity) GetType() QuantityType {
	if x != nil {
		return x.Type
	}
	return QuantityType_UNKNOWN_QUANTITY_TYPE
}

func (x *Quantity) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

// Deadline message to specify type, value, and action
type Deadline struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type DeadlineType `protobuf:"varint,1,opt,name=type,proto3,enum=ordermanager.DeadlineType" json:"type,omitempty"`
	// Type: TIME => the number of seconds
	Value  string         `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Action DeadlineAction `protobuf:"varint,3,opt,name=action,proto3,enum=ordermanager.DeadlineAction" json:"action,omitempty"`
}

func (x *Deadline) Reset() {
	*x = Deadline{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protofiles_ordermanager_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Deadline) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Deadline) ProtoMessage() {}

func (x *Deadline) ProtoReflect() protoreflect.Message {
	mi := &file_protofiles_ordermanager_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Deadline.ProtoReflect.Descriptor instead.
func (*Deadline) Descriptor() ([]byte, []int) {
	return file_protofiles_ordermanager_proto_rawDescGZIP(), []int{3}
}

func (x *Deadline) GetType() DeadlineType {
	if x != nil {
		return x.Type
	}
	return DeadlineType_UNKNOWN_DEADLINE_TYPE
}

func (x *Deadline) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *Deadline) GetAction() DeadlineAction {
	if x != nil {
		return x.Action
	}
	return DeadlineAction_UNKNOWN_DEADLINE_ACTION
}

var File_protofiles_ordermanager_proto protoreflect.FileDescriptor

var file_protofiles_ordermanager_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2f, 0x6f, 0x72, 0x64,
	0x65, 0x72, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0c, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x22, 0xa6, 0x02,
	0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x69, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x69, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x61, 0x72, 0x6b,
	0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74,
	0x12, 0x30, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x18, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e,
	0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x32, 0x0a, 0x08, 0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x6d, 0x61, 0x6e, 0x61,
	0x67, 0x65, 0x72, 0x2e, 0x42, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x52, 0x08, 0x62, 0x65,
	0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x32, 0x0a, 0x08,
	0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16,
	0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x51, 0x75,
	0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x12, 0x34, 0x0a, 0x09, 0x64, 0x65, 0x61, 0x64, 0x6c, 0x69, 0x6e, 0x65, 0x73, 0x18, 0x07, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x6d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x72, 0x2e, 0x44, 0x65, 0x61, 0x64, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x09, 0x64, 0x65, 0x61,
	0x64, 0x6c, 0x69, 0x6e, 0x65, 0x73, 0x22, 0x49, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07,
	0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x22, 0x50, 0x0a, 0x08, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x2e, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x51, 0x75, 0x61, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x22, 0x86, 0x01, 0x0a, 0x08, 0x44, 0x65, 0x61, 0x64, 0x6c, 0x69, 0x6e, 0x65,
	0x12, 0x2e, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a,
	0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x44, 0x65,
	0x61, 0x64, 0x6c, 0x69, 0x6e, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x34, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x6d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x44, 0x65, 0x61, 0x64, 0x6c, 0x69, 0x6e, 0x65, 0x41, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2a, 0x33, 0x0a, 0x0a,
	0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x0e, 0x55, 0x4e,
	0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x00, 0x12, 0x08,
	0x0a, 0x04, 0x53, 0x45, 0x4c, 0x4c, 0x10, 0x01, 0x12, 0x07, 0x0a, 0x03, 0x42, 0x55, 0x59, 0x10,
	0x02, 0x2a, 0x37, 0x0a, 0x08, 0x42, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x12, 0x14, 0x0a,
	0x10, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x42, 0x45, 0x48, 0x41, 0x56, 0x49, 0x4f,
	0x52, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x4d, 0x41, 0x52, 0x4b, 0x45, 0x54, 0x10, 0x01, 0x12,
	0x09, 0x0a, 0x05, 0x4c, 0x49, 0x4d, 0x49, 0x54, 0x10, 0x02, 0x2a, 0x41, 0x0a, 0x0c, 0x51, 0x75,
	0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x19, 0x0a, 0x15, 0x55, 0x4e,
	0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x51, 0x55, 0x41, 0x4e, 0x54, 0x49, 0x54, 0x59, 0x5f, 0x54,
	0x59, 0x50, 0x45, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x46, 0x49, 0x58, 0x45, 0x44, 0x10, 0x01,
	0x12, 0x0b, 0x0a, 0x07, 0x50, 0x45, 0x52, 0x43, 0x45, 0x4e, 0x54, 0x10, 0x02, 0x2a, 0x33, 0x0a,
	0x0c, 0x44, 0x65, 0x61, 0x64, 0x6c, 0x69, 0x6e, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x19, 0x0a,
	0x15, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x44, 0x45, 0x41, 0x44, 0x4c, 0x49, 0x4e,
	0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x54, 0x49, 0x4d, 0x45,
	0x10, 0x01, 0x2a, 0x4d, 0x0a, 0x0e, 0x44, 0x65, 0x61, 0x64, 0x6c, 0x69, 0x6e, 0x65, 0x41, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x17, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f,
	0x44, 0x45, 0x41, 0x44, 0x4c, 0x49, 0x4e, 0x45, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x10,
	0x00, 0x12, 0x0a, 0x0a, 0x06, 0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c, 0x10, 0x01, 0x12, 0x12, 0x0a,
	0x0e, 0x53, 0x45, 0x4c, 0x4c, 0x5f, 0x42, 0x59, 0x5f, 0x4d, 0x41, 0x52, 0x4b, 0x45, 0x54, 0x10,
	0x02, 0x32, 0x62, 0x0a, 0x0c, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x72, 0x12, 0x52, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72,
	0x12, 0x20, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x21, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x72, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0d, 0x5a, 0x0b, 0x2e, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72,
	0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protofiles_ordermanager_proto_rawDescOnce sync.Once
	file_protofiles_ordermanager_proto_rawDescData = file_protofiles_ordermanager_proto_rawDesc
)

func file_protofiles_ordermanager_proto_rawDescGZIP() []byte {
	file_protofiles_ordermanager_proto_rawDescOnce.Do(func() {
		file_protofiles_ordermanager_proto_rawDescData = protoimpl.X.CompressGZIP(file_protofiles_ordermanager_proto_rawDescData)
	})
	return file_protofiles_ordermanager_proto_rawDescData
}

var file_protofiles_ordermanager_proto_enumTypes = make([]protoimpl.EnumInfo, 5)
var file_protofiles_ordermanager_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_protofiles_ordermanager_proto_goTypes = []any{
	(ActionType)(0),             // 0: ordermanager.ActionType
	(Behavior)(0),               // 1: ordermanager.Behavior
	(QuantityType)(0),           // 2: ordermanager.QuantityType
	(DeadlineType)(0),           // 3: ordermanager.DeadlineType
	(DeadlineAction)(0),         // 4: ordermanager.DeadlineAction
	(*CreateOrderRequest)(nil),  // 5: ordermanager.CreateOrderRequest
	(*CreateOrderResponse)(nil), // 6: ordermanager.CreateOrderResponse
	(*Quantity)(nil),            // 7: ordermanager.Quantity
	(*Deadline)(nil),            // 8: ordermanager.Deadline
}
var file_protofiles_ordermanager_proto_depIdxs = []int32{
	0, // 0: ordermanager.CreateOrderRequest.action:type_name -> ordermanager.ActionType
	1, // 1: ordermanager.CreateOrderRequest.behavior:type_name -> ordermanager.Behavior
	7, // 2: ordermanager.CreateOrderRequest.quantity:type_name -> ordermanager.Quantity
	8, // 3: ordermanager.CreateOrderRequest.deadlines:type_name -> ordermanager.Deadline
	2, // 4: ordermanager.Quantity.type:type_name -> ordermanager.QuantityType
	3, // 5: ordermanager.Deadline.type:type_name -> ordermanager.DeadlineType
	4, // 6: ordermanager.Deadline.action:type_name -> ordermanager.DeadlineAction
	5, // 7: ordermanager.OrderManager.CreateOrder:input_type -> ordermanager.CreateOrderRequest
	6, // 8: ordermanager.OrderManager.CreateOrder:output_type -> ordermanager.CreateOrderResponse
	8, // [8:9] is the sub-list for method output_type
	7, // [7:8] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_protofiles_ordermanager_proto_init() }
func file_protofiles_ordermanager_proto_init() {
	if File_protofiles_ordermanager_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protofiles_ordermanager_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*CreateOrderRequest); i {
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
		file_protofiles_ordermanager_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*CreateOrderResponse); i {
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
		file_protofiles_ordermanager_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*Quantity); i {
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
		file_protofiles_ordermanager_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*Deadline); i {
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
			RawDescriptor: file_protofiles_ordermanager_proto_rawDesc,
			NumEnums:      5,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protofiles_ordermanager_proto_goTypes,
		DependencyIndexes: file_protofiles_ordermanager_proto_depIdxs,
		EnumInfos:         file_protofiles_ordermanager_proto_enumTypes,
		MessageInfos:      file_protofiles_ordermanager_proto_msgTypes,
	}.Build()
	File_protofiles_ordermanager_proto = out.File
	file_protofiles_ordermanager_proto_rawDesc = nil
	file_protofiles_ordermanager_proto_goTypes = nil
	file_protofiles_ordermanager_proto_depIdxs = nil
}
