// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.25.3
// source: defaultResponse.proto

package miniresolverproto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ResultStatus int32

const (
	ResultStatus_Error    ResultStatus = 0
	ResultStatus_OK       ResultStatus = 1
	ResultStatus_Warning  ResultStatus = 2
	ResultStatus_NotFound ResultStatus = 3
)

// Enum value maps for ResultStatus.
var (
	ResultStatus_name = map[int32]string{
		0: "Error",
		1: "OK",
		2: "Warning",
		3: "NotFound",
	}
	ResultStatus_value = map[string]int32{
		"Error":    0,
		"OK":       1,
		"Warning":  2,
		"NotFound": 3,
	}
)

func (x ResultStatus) Enum() *ResultStatus {
	p := new(ResultStatus)
	*p = x
	return p
}

func (x ResultStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ResultStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_defaultResponse_proto_enumTypes[0].Descriptor()
}

func (ResultStatus) Type() protoreflect.EnumType {
	return &file_defaultResponse_proto_enumTypes[0]
}

func (x ResultStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ResultStatus.Descriptor instead.
func (ResultStatus) EnumDescriptor() ([]byte, []int) {
	return file_defaultResponse_proto_rawDescGZIP(), []int{0}
}

type DefaultResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status  ResultStatus `protobuf:"varint,1,opt,name=status,proto3,enum=miniresolverproto.ResultStatus" json:"status,omitempty"`
	Message string       `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data    *anypb.Any   `protobuf:"bytes,3,opt,name=data,proto3,oneof" json:"data,omitempty"`
}

func (x *DefaultResponse) Reset() {
	*x = DefaultResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_defaultResponse_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DefaultResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DefaultResponse) ProtoMessage() {}

func (x *DefaultResponse) ProtoReflect() protoreflect.Message {
	mi := &file_defaultResponse_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DefaultResponse.ProtoReflect.Descriptor instead.
func (*DefaultResponse) Descriptor() ([]byte, []int) {
	return file_defaultResponse_proto_rawDescGZIP(), []int{0}
}

func (x *DefaultResponse) GetStatus() ResultStatus {
	if x != nil {
		return x.Status
	}
	return ResultStatus_Error
}

func (x *DefaultResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *DefaultResponse) GetData() *anypb.Any {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_defaultResponse_proto protoreflect.FileDescriptor

var file_defaultResponse_proto_rawDesc = []byte{
	0x0a, 0x15, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x6d, 0x69, 0x6e, 0x69, 0x72, 0x65, 0x73,
	0x6f, 0x6c, 0x76, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9c, 0x01, 0x0a, 0x0f, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x6d, 0x69, 0x6e, 0x69,
	0x72, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2d, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79,
	0x48, 0x00, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f,
	0x64, 0x61, 0x74, 0x61, 0x2a, 0x3c, 0x0a, 0x0c, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x00, 0x12,
	0x06, 0x0a, 0x02, 0x4f, 0x4b, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x57, 0x61, 0x72, 0x6e, 0x69,
	0x6e, 0x67, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x4e, 0x6f, 0x74, 0x46, 0x6f, 0x75, 0x6e, 0x64,
	0x10, 0x03, 0x42, 0x85, 0x01, 0x0a, 0x19, 0x63, 0x68, 0x2e, 0x75, 0x6e, 0x69, 0x62, 0x61, 0x73,
	0x2e, 0x75, 0x62, 0x2e, 0x6d, 0x69, 0x6e, 0x69, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x72,
	0x42, 0x11, 0x4d, 0x69, 0x6e, 0x69, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x72, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x6a, 0x65, 0x34, 0x2f, 0x6d, 0x69, 0x6e, 0x69, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x76,
	0x65, 0x72, 0x2f, 0x76, 0x32, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x69, 0x6e, 0x69, 0x72, 0x65,
	0x73, 0x6f, 0x6c, 0x76, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0xa2, 0x02, 0x03, 0x55, 0x42,
	0x42, 0xaa, 0x02, 0x16, 0x55, 0x6e, 0x69, 0x62, 0x61, 0x73, 0x2e, 0x55, 0x42, 0x2e, 0x4d, 0x69,
	0x6e, 0x69, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_defaultResponse_proto_rawDescOnce sync.Once
	file_defaultResponse_proto_rawDescData = file_defaultResponse_proto_rawDesc
)

func file_defaultResponse_proto_rawDescGZIP() []byte {
	file_defaultResponse_proto_rawDescOnce.Do(func() {
		file_defaultResponse_proto_rawDescData = protoimpl.X.CompressGZIP(file_defaultResponse_proto_rawDescData)
	})
	return file_defaultResponse_proto_rawDescData
}

var file_defaultResponse_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_defaultResponse_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_defaultResponse_proto_goTypes = []interface{}{
	(ResultStatus)(0),       // 0: miniresolverproto.ResultStatus
	(*DefaultResponse)(nil), // 1: miniresolverproto.DefaultResponse
	(*anypb.Any)(nil),       // 2: google.protobuf.Any
}
var file_defaultResponse_proto_depIdxs = []int32{
	0, // 0: miniresolverproto.DefaultResponse.status:type_name -> miniresolverproto.ResultStatus
	2, // 1: miniresolverproto.DefaultResponse.data:type_name -> google.protobuf.Any
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_defaultResponse_proto_init() }
func file_defaultResponse_proto_init() {
	if File_defaultResponse_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_defaultResponse_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DefaultResponse); i {
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
	file_defaultResponse_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_defaultResponse_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_defaultResponse_proto_goTypes,
		DependencyIndexes: file_defaultResponse_proto_depIdxs,
		EnumInfos:         file_defaultResponse_proto_enumTypes,
		MessageInfos:      file_defaultResponse_proto_msgTypes,
	}.Build()
	File_defaultResponse_proto = out.File
	file_defaultResponse_proto_rawDesc = nil
	file_defaultResponse_proto_goTypes = nil
	file_defaultResponse_proto_depIdxs = nil
}
