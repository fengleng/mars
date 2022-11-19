// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.17.3
// source: test.proto

package binding

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// The request message containing the user's name.
type HelloRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name         string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Sub          *Sub                   `protobuf:"bytes,2,opt,name=sub,proto3" json:"sub,omitempty"`
	UpdateMask   *fieldmaskpb.FieldMask `protobuf:"bytes,3,opt,name=update_mask,json=updateMask,proto3" json:"update_mask,omitempty"`
	OptInt32     *int32                 `protobuf:"varint,4,opt,name=opt_int32,json=optInt32,proto3,oneof" json:"opt_int32,omitempty"`
	OptInt64     *int64                 `protobuf:"varint,5,opt,name=opt_int64,json=optInt64,proto3,oneof" json:"opt_int64,omitempty"`
	OptString    *string                `protobuf:"bytes,6,opt,name=opt_string,json=optString,proto3,oneof" json:"opt_string,omitempty"`
	SubField     *Sub                   `protobuf:"bytes,7,opt,name=subField,proto3" json:"subField,omitempty"`
	TestRepeated []string               `protobuf:"bytes,8,rep,name=test_repeated,proto3" json:"test_repeated,omitempty"`
}

func (x *HelloRequest) Reset() {
	*x = HelloRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HelloRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloRequest) ProtoMessage() {}

func (x *HelloRequest) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloRequest.ProtoReflect.Descriptor instead.
func (*HelloRequest) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{0}
}

func (x *HelloRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *HelloRequest) GetSub() *Sub {
	if x != nil {
		return x.Sub
	}
	return nil
}

func (x *HelloRequest) GetUpdateMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.UpdateMask
	}
	return nil
}

func (x *HelloRequest) GetOptInt32() int32 {
	if x != nil && x.OptInt32 != nil {
		return *x.OptInt32
	}
	return 0
}

func (x *HelloRequest) GetOptInt64() int64 {
	if x != nil && x.OptInt64 != nil {
		return *x.OptInt64
	}
	return 0
}

func (x *HelloRequest) GetOptString() string {
	if x != nil && x.OptString != nil {
		return *x.OptString
	}
	return ""
}

func (x *HelloRequest) GetSubField() *Sub {
	if x != nil {
		return x.SubField
	}
	return nil
}

func (x *HelloRequest) GetTestRepeated() []string {
	if x != nil {
		return x.TestRepeated
	}
	return nil
}

type Sub struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,json=naming,proto3" json:"name,omitempty"`
}

func (x *Sub) Reset() {
	*x = Sub{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sub) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sub) ProtoMessage() {}

func (x *Sub) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sub.ProtoReflect.Descriptor instead.
func (*Sub) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{1}
}

func (x *Sub) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_test_proto protoreflect.FileDescriptor

var file_test_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x62, 0x69,
	0x6e, 0x64, 0x69, 0x6e, 0x67, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x6d, 0x61, 0x73,
	0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe2, 0x02, 0x0a, 0x0c, 0x48, 0x65, 0x6c, 0x6c,
	0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x03,
	0x73, 0x75, 0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x62, 0x69, 0x6e, 0x64,
	0x69, 0x6e, 0x67, 0x2e, 0x53, 0x75, 0x62, 0x52, 0x03, 0x73, 0x75, 0x62, 0x12, 0x3b, 0x0a, 0x0b,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4d, 0x61, 0x73, 0x6b, 0x52, 0x0a, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x73, 0x6b, 0x12, 0x20, 0x0a, 0x09, 0x6f, 0x70, 0x74,
	0x5f, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x48, 0x00, 0x52, 0x08,
	0x6f, 0x70, 0x74, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x88, 0x01, 0x01, 0x12, 0x20, 0x0a, 0x09, 0x6f,
	0x70, 0x74, 0x5f, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x48, 0x01,
	0x52, 0x08, 0x6f, 0x70, 0x74, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x88, 0x01, 0x01, 0x12, 0x22, 0x0a,
	0x0a, 0x6f, 0x70, 0x74, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x02, 0x52, 0x09, 0x6f, 0x70, 0x74, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x88, 0x01,
	0x01, 0x12, 0x28, 0x0a, 0x08, 0x73, 0x75, 0x62, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x62, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x53, 0x75,
	0x62, 0x52, 0x08, 0x73, 0x75, 0x62, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x24, 0x0a, 0x0d, 0x74,
	0x65, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x70, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x08, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x0d, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x70, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x6f, 0x70, 0x74, 0x5f, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x42,
	0x0c, 0x0a, 0x0a, 0x5f, 0x6f, 0x70, 0x74, 0x5f, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x42, 0x0d, 0x0a,
	0x0b, 0x5f, 0x6f, 0x70, 0x74, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x22, 0x1b, 0x0a, 0x03,
	0x53, 0x75, 0x62, 0x12, 0x14, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x6e, 0x61, 0x6d, 0x69, 0x6e, 0x67, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x2d, 0x6b, 0x72, 0x61, 0x74, 0x6f,
	0x73, 0x2f, 0x6b, 0x72, 0x61, 0x74, 0x6f, 0x73, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f,
	0x72, 0x74, 0x2f, 0x62, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_test_proto_rawDescOnce sync.Once
	file_test_proto_rawDescData = file_test_proto_rawDesc
)

func file_test_proto_rawDescGZIP() []byte {
	file_test_proto_rawDescOnce.Do(func() {
		file_test_proto_rawDescData = protoimpl.X.CompressGZIP(file_test_proto_rawDescData)
	})
	return file_test_proto_rawDescData
}

var file_test_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_test_proto_goTypes = []interface{}{
	(*HelloRequest)(nil),          // 0: binding.HelloRequest
	(*Sub)(nil),                   // 1: binding.Sub
	(*fieldmaskpb.FieldMask)(nil), // 2: google.protobuf.FieldMask
}
var file_test_proto_depIdxs = []int32{
	1, // 0: binding.HelloRequest.sub:type_name -> binding.Sub
	2, // 1: binding.HelloRequest.update_mask:type_name -> google.protobuf.FieldMask
	1, // 2: binding.HelloRequest.subField:type_name -> binding.Sub
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_test_proto_init() }
func file_test_proto_init() {
	if File_test_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_test_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HelloRequest); i {
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
		file_test_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Sub); i {
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
	file_test_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_test_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_test_proto_goTypes,
		DependencyIndexes: file_test_proto_depIdxs,
		MessageInfos:      file_test_proto_msgTypes,
	}.Build()
	File_test_proto = out.File
	file_test_proto_rawDesc = nil
	file_test_proto_goTypes = nil
	file_test_proto_depIdxs = nil
}