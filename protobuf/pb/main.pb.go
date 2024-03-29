// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v3.11.2
// source: main.proto

package pb

import (
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	any1 "github.com/golang/protobuf/ptypes/any"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type Gender int32

const (
	Gender_GENDER_DEFAULT Gender = 0 // 枚举值中有相同值时，后面的都是别名，需要声明allow_alias，否则报错
	Gender_GENDER_UNKNOWN Gender = 0 // 建议：枚举变量名都要带上枚举类名，因为枚举变量名是全局唯一的
	Gender_GENDER_FEMALE  Gender = 1
	Gender_GENDER_MALE    Gender = 2
)

// Enum value maps for Gender.
var (
	Gender_name = map[int32]string{
		0: "GENDER_DEFAULT",
		// Duplicate value: 0: "GENDER_UNKNOWN",
		1: "GENDER_FEMALE",
		2: "GENDER_MALE",
	}
	Gender_value = map[string]int32{
		"GENDER_DEFAULT": 0,
		"GENDER_UNKNOWN": 0,
		"GENDER_FEMALE":  1,
		"GENDER_MALE":    2,
	}
)

func (x Gender) Enum() *Gender {
	p := new(Gender)
	*p = x
	return p
}

func (x Gender) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Gender) Descriptor() protoreflect.EnumDescriptor {
	return file_main_proto_enumTypes[0].Descriptor()
}

func (Gender) Type() protoreflect.EnumType {
	return &file_main_proto_enumTypes[0]
}

func (x Gender) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Gender.Descriptor instead.
func (Gender) EnumDescriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{0}
}

type Subject int32

const (
	Subject_SUBJECT_UNKNOWN Subject = 0 // 注意：protobufV3为了减少无效数据的传输会放弃传输零值数据，比如枚举的0，bool的false，字符串的""。这是设计意图。如果造成困扰可以改用protobufV2并使用required限定。要用v3就避免在枚举中有效值使用零值。
	Subject_SUBJECT_MATH    Subject = 1
	Subject_SUBJECT_ENGLISH Subject = 2
)

// Enum value maps for Subject.
var (
	Subject_name = map[int32]string{
		0: "SUBJECT_UNKNOWN",
		1: "SUBJECT_MATH",
		2: "SUBJECT_ENGLISH",
	}
	Subject_value = map[string]int32{
		"SUBJECT_UNKNOWN": 0,
		"SUBJECT_MATH":    1,
		"SUBJECT_ENGLISH": 2,
	}
)

func (x Subject) Enum() *Subject {
	p := new(Subject)
	*p = x
	return p
}

func (x Subject) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Subject) Descriptor() protoreflect.EnumDescriptor {
	return file_main_proto_enumTypes[1].Descriptor()
}

func (Subject) Type() protoreflect.EnumType {
	return &file_main_proto_enumTypes[1]
}

func (x Subject) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Subject.Descriptor instead.
func (Subject) EnumDescriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{1}
}

type DefaultRemark int32

const (
	DefaultRemark_REMARK_EMPTY     DefaultRemark = 0
	DefaultRemark_REMARK_EXCELLENT DefaultRemark = 1
	DefaultRemark_REMARK_GOOD      DefaultRemark = 2
	DefaultRemark_REMARK_BAD       DefaultRemark = 3
)

// Enum value maps for DefaultRemark.
var (
	DefaultRemark_name = map[int32]string{
		0: "REMARK_EMPTY",
		1: "REMARK_EXCELLENT",
		2: "REMARK_GOOD",
		3: "REMARK_BAD",
	}
	DefaultRemark_value = map[string]int32{
		"REMARK_EMPTY":     0,
		"REMARK_EXCELLENT": 1,
		"REMARK_GOOD":      2,
		"REMARK_BAD":       3,
	}
)

func (x DefaultRemark) Enum() *DefaultRemark {
	p := new(DefaultRemark)
	*p = x
	return p
}

func (x DefaultRemark) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DefaultRemark) Descriptor() protoreflect.EnumDescriptor {
	return file_main_proto_enumTypes[2].Descriptor()
}

func (DefaultRemark) Type() protoreflect.EnumType {
	return &file_main_proto_enumTypes[2]
}

func (x DefaultRemark) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DefaultRemark.Descriptor instead.
func (DefaultRemark) EnumDescriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{2}
}

// 用枚举实现字典(dict)功能
// 在不同语言的实现中获取此选项的方式不同。
// 比如Java中是`UnitType.KM_PER_HOUR.getValueDescriptor().getOptions().getExtension(MyOuterClass.name);`
// Go中，t := pb.PerformanceType_PERFORMANCE_TYPE_ATTENDANCE
// tt, err := proto.GetExtension(proto.MessageV1(t.Descriptor().Values().ByNumber(t.Number()).Options()), pb.E_CustomName)
type PerformanceType int32

const (
	PerformanceType_PERFORMANCE_TYPE_ATTENDANCE PerformanceType = 0 // 在枚举中想要得到string类型的值，可以指定自定义选项`custom_name`的值。
	PerformanceType_PERFORMANCE_TYPE_MIDTERM    PerformanceType = 1
)

// Enum value maps for PerformanceType.
var (
	PerformanceType_name = map[int32]string{
		0: "PERFORMANCE_TYPE_ATTENDANCE",
		1: "PERFORMANCE_TYPE_MIDTERM",
	}
	PerformanceType_value = map[string]int32{
		"PERFORMANCE_TYPE_ATTENDANCE": 0,
		"PERFORMANCE_TYPE_MIDTERM":    1,
	}
)

func (x PerformanceType) Enum() *PerformanceType {
	p := new(PerformanceType)
	*p = x
	return p
}

func (x PerformanceType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PerformanceType) Descriptor() protoreflect.EnumDescriptor {
	return file_main_proto_enumTypes[3].Descriptor()
}

func (PerformanceType) Type() protoreflect.EnumType {
	return &file_main_proto_enumTypes[3]
}

func (x PerformanceType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PerformanceType.Descriptor instead.
func (PerformanceType) EnumDescriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{3}
}

type GetStudentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetStudentRequest) Reset() {
	*x = GetStudentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_main_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetStudentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStudentRequest) ProtoMessage() {}

func (x *GetStudentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStudentRequest.ProtoReflect.Descriptor instead.
func (*GetStudentRequest) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{0}
}

func (x *GetStudentRequest) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type GetStudentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Student *Student `protobuf:"bytes,1,opt,name=student,proto3" json:"student,omitempty"`
}

func (x *GetStudentResponse) Reset() {
	*x = GetStudentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_main_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetStudentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStudentResponse) ProtoMessage() {}

func (x *GetStudentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStudentResponse.ProtoReflect.Descriptor instead.
func (*GetStudentResponse) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{1}
}

func (x *GetStudentResponse) GetStudent() *Student {
	if x != nil {
		return x.Student
	}
	return nil
}

type Student struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         uint32      `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name       string      `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`                                // 内置类型
	Gender     Gender      `protobuf:"varint,3,opt,name=gender,proto3,enum=proto.Gender" json:"gender,omitempty"`         // 自定义枚举类型
	Scores     []*Score    `protobuf:"bytes,4,rep,name=scores,proto3" json:"scores,omitempty"`                            // 自定义类型Score的多重结构
	PassStatus bool        `protobuf:"varint,5,opt,name=pass_status,json=passStatus,proto3" json:"pass_status,omitempty"` // 是否通过考试
	Homeworks  []*Homework `protobuf:"bytes,6,rep,name=homeworks,proto3" json:"homeworks,omitempty"`                      // 个人作业信息
	Address    *Address    `protobuf:"bytes,7,opt,name=address,proto3" json:"address,omitempty"`                          // 导入其他包中的地址信息。需要用`package_name.Type`作为类型，而不是`Address`
	// Types that are assignable to Remark:
	//
	//	*Student_DefaultRemark
	//	*Student_OtherRemark
	Remark            isStudent_Remark     `protobuf_oneof:"remark"`
	Special           *any1.Any            `protobuf:"bytes,10,opt,name=special,proto3" json:"special,omitempty"`                                                                                                                                        // 特长：pb扩展类型any，允许任何类型的值。类型同样使用`package_name.Type`而不是`Any`
	UsualPerformances map[string]float64   `protobuf:"bytes,11,rep,name=usual_performances,json=usualPerformances,proto3" json:"usual_performances,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"fixed64,2,opt,name=value,proto3"` // 平时表现：map类型：自定义类型->内置类型。另外：Key in map fields cannot be enum types.
	LastUpdated       *timestamp.Timestamp `protobuf:"bytes,12,opt,name=last_updated,json=lastUpdated,proto3" json:"last_updated,omitempty"`
}

func (x *Student) Reset() {
	*x = Student{}
	if protoimpl.UnsafeEnabled {
		mi := &file_main_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Student) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Student) ProtoMessage() {}

func (x *Student) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Student.ProtoReflect.Descriptor instead.
func (*Student) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{2}
}

func (x *Student) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Student) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Student) GetGender() Gender {
	if x != nil {
		return x.Gender
	}
	return Gender_GENDER_DEFAULT
}

func (x *Student) GetScores() []*Score {
	if x != nil {
		return x.Scores
	}
	return nil
}

func (x *Student) GetPassStatus() bool {
	if x != nil {
		return x.PassStatus
	}
	return false
}

func (x *Student) GetHomeworks() []*Homework {
	if x != nil {
		return x.Homeworks
	}
	return nil
}

func (x *Student) GetAddress() *Address {
	if x != nil {
		return x.Address
	}
	return nil
}

func (m *Student) GetRemark() isStudent_Remark {
	if m != nil {
		return m.Remark
	}
	return nil
}

func (x *Student) GetDefaultRemark() DefaultRemark {
	if x, ok := x.GetRemark().(*Student_DefaultRemark); ok {
		return x.DefaultRemark
	}
	return DefaultRemark_REMARK_EMPTY
}

func (x *Student) GetOtherRemark() string {
	if x, ok := x.GetRemark().(*Student_OtherRemark); ok {
		return x.OtherRemark
	}
	return ""
}

func (x *Student) GetSpecial() *any1.Any {
	if x != nil {
		return x.Special
	}
	return nil
}

func (x *Student) GetUsualPerformances() map[string]float64 {
	if x != nil {
		return x.UsualPerformances
	}
	return nil
}

func (x *Student) GetLastUpdated() *timestamp.Timestamp {
	if x != nil {
		return x.LastUpdated
	}
	return nil
}

type isStudent_Remark interface {
	isStudent_Remark()
}

type Student_DefaultRemark struct {
	DefaultRemark DefaultRemark `protobuf:"varint,8,opt,name=default_remark,json=defaultRemark,proto3,enum=proto.DefaultRemark,oneof"` // one_of中的变量名不能是出现过的`remark`，会报错`already defined`
}

type Student_OtherRemark struct {
	OtherRemark string `protobuf:"bytes,9,opt,name=other_remark,json=otherRemark,proto3,oneof"`
}

func (*Student_DefaultRemark) isStudent_Remark() {}

func (*Student_OtherRemark) isStudent_Remark() {}

type Score struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Subject Subject `protobuf:"varint,1,opt,name=subject,proto3,enum=proto.Subject" json:"subject,omitempty"`
	Score   float64 `protobuf:"fixed64,2,opt,name=score,proto3" json:"score,omitempty"`
}

func (x *Score) Reset() {
	*x = Score{}
	if protoimpl.UnsafeEnabled {
		mi := &file_main_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Score) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Score) ProtoMessage() {}

func (x *Score) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Score.ProtoReflect.Descriptor instead.
func (*Score) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{3}
}

func (x *Score) GetSubject() Subject {
	if x != nil {
		return x.Subject
	}
	return Subject_SUBJECT_UNKNOWN
}

func (x *Score) GetScore() float64 {
	if x != nil {
		return x.Score
	}
	return 0
}

type Homework struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id  uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Url string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *Homework) Reset() {
	*x = Homework{}
	if protoimpl.UnsafeEnabled {
		mi := &file_main_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Homework) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Homework) ProtoMessage() {}

func (x *Homework) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Homework.ProtoReflect.Descriptor instead.
func (*Homework) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{4}
}

func (x *Homework) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Homework) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

var file_main_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptor.EnumValueOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         123456789,
		Name:          "proto.custom_name",
		Tag:           "bytes,123456789,opt,name=custom_name",
		Filename:      "main.proto",
	},
}

// Extension fields to descriptor.EnumValueOptions.
var (
	// optional string custom_name = 123456789;
	E_CustomName = &file_main_proto_extTypes[0] // 为枚举增加一个自定义选项，顺序为超级大
)

var File_main_proto protoreflect.FileDescriptor

var file_main_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x0d, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x23, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x02, 0x69, 0x64, 0x22, 0x3e, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x53, 0x74, 0x75, 0x64,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x07, 0x73,
	0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x07, 0x73, 0x74,
	0x75, 0x64, 0x65, 0x6e, 0x74, 0x22, 0xef, 0x04, 0x0a, 0x07, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x06, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x52, 0x06, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x24, 0x0a, 0x06,
	0x73, 0x63, 0x6f, 0x72, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x52, 0x06, 0x73, 0x63, 0x6f, 0x72,
	0x65, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x61, 0x73, 0x73, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x70, 0x61, 0x73, 0x73, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x2d, 0x0a, 0x09, 0x68, 0x6f, 0x6d, 0x65, 0x77, 0x6f, 0x72, 0x6b, 0x73,
	0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x48,
	0x6f, 0x6d, 0x65, 0x77, 0x6f, 0x72, 0x6b, 0x52, 0x09, 0x68, 0x6f, 0x6d, 0x65, 0x77, 0x6f, 0x72,
	0x6b, 0x73, 0x12, 0x2a, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x2e, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x3d,
	0x0a, 0x0e, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x72, 0x65, 0x6d, 0x61, 0x72, 0x6b,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44,
	0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x52, 0x65, 0x6d, 0x61, 0x72, 0x6b, 0x48, 0x00, 0x52, 0x0d,
	0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x52, 0x65, 0x6d, 0x61, 0x72, 0x6b, 0x12, 0x23, 0x0a,
	0x0c, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x5f, 0x72, 0x65, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0b, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x52, 0x65, 0x6d, 0x61,
	0x72, 0x6b, 0x12, 0x2e, 0x0a, 0x07, 0x73, 0x70, 0x65, 0x63, 0x69, 0x61, 0x6c, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x07, 0x73, 0x70, 0x65, 0x63, 0x69,
	0x61, 0x6c, 0x12, 0x54, 0x0a, 0x12, 0x75, 0x73, 0x75, 0x61, 0x6c, 0x5f, 0x70, 0x65, 0x72, 0x66,
	0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x2e, 0x55,
	0x73, 0x75, 0x61, 0x6c, 0x50, 0x65, 0x72, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x63, 0x65, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x11, 0x75, 0x73, 0x75, 0x61, 0x6c, 0x50, 0x65, 0x72, 0x66,
	0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x12, 0x3d, 0x0a, 0x0c, 0x6c, 0x61, 0x73, 0x74,
	0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0b, 0x6c, 0x61, 0x73, 0x74,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x1a, 0x44, 0x0a, 0x16, 0x55, 0x73, 0x75, 0x61, 0x6c,
	0x50, 0x65, 0x72, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x08, 0x0a,
	0x06, 0x72, 0x65, 0x6d, 0x61, 0x72, 0x6b, 0x22, 0x47, 0x0a, 0x05, 0x53, 0x63, 0x6f, 0x72, 0x65,
	0x12, 0x28, 0x0a, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x75, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x52, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63,
	0x6f, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65,
	0x22, 0x2c, 0x0a, 0x08, 0x48, 0x6f, 0x6d, 0x65, 0x77, 0x6f, 0x72, 0x6b, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x2a, 0x58,
	0x0a, 0x06, 0x47, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x0e, 0x47, 0x45, 0x4e, 0x44,
	0x45, 0x52, 0x5f, 0x44, 0x45, 0x46, 0x41, 0x55, 0x4c, 0x54, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e,
	0x47, 0x45, 0x4e, 0x44, 0x45, 0x52, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00,
	0x12, 0x11, 0x0a, 0x0d, 0x47, 0x45, 0x4e, 0x44, 0x45, 0x52, 0x5f, 0x46, 0x45, 0x4d, 0x41, 0x4c,
	0x45, 0x10, 0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x47, 0x45, 0x4e, 0x44, 0x45, 0x52, 0x5f, 0x4d, 0x41,
	0x4c, 0x45, 0x10, 0x02, 0x1a, 0x02, 0x10, 0x01, 0x2a, 0x60, 0x0a, 0x07, 0x53, 0x75, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x12, 0x13, 0x0a, 0x0f, 0x53, 0x55, 0x42, 0x4a, 0x45, 0x43, 0x54, 0x5f, 0x55,
	0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x1c, 0x0a, 0x0c, 0x53, 0x55, 0x42, 0x4a,
	0x45, 0x43, 0x54, 0x5f, 0x4d, 0x41, 0x54, 0x48, 0x10, 0x01, 0x1a, 0x0a, 0xaa, 0xd1, 0xf9, 0xd6,
	0x03, 0x04, 0x4d, 0x41, 0x54, 0x48, 0x12, 0x22, 0x0a, 0x0f, 0x53, 0x55, 0x42, 0x4a, 0x45, 0x43,
	0x54, 0x5f, 0x45, 0x4e, 0x47, 0x4c, 0x49, 0x53, 0x48, 0x10, 0x02, 0x1a, 0x0d, 0xaa, 0xd1, 0xf9,
	0xd6, 0x03, 0x07, 0x45, 0x4e, 0x47, 0x4c, 0x49, 0x53, 0x48, 0x2a, 0x58, 0x0a, 0x0d, 0x44, 0x65,
	0x66, 0x61, 0x75, 0x6c, 0x74, 0x52, 0x65, 0x6d, 0x61, 0x72, 0x6b, 0x12, 0x10, 0x0a, 0x0c, 0x52,
	0x45, 0x4d, 0x41, 0x52, 0x4b, 0x5f, 0x45, 0x4d, 0x50, 0x54, 0x59, 0x10, 0x00, 0x12, 0x14, 0x0a,
	0x10, 0x52, 0x45, 0x4d, 0x41, 0x52, 0x4b, 0x5f, 0x45, 0x58, 0x43, 0x45, 0x4c, 0x4c, 0x45, 0x4e,
	0x54, 0x10, 0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x52, 0x45, 0x4d, 0x41, 0x52, 0x4b, 0x5f, 0x47, 0x4f,
	0x4f, 0x44, 0x10, 0x02, 0x12, 0x0e, 0x0a, 0x0a, 0x52, 0x45, 0x4d, 0x41, 0x52, 0x4b, 0x5f, 0x42,
	0x41, 0x44, 0x10, 0x03, 0x2a, 0x71, 0x0a, 0x0f, 0x50, 0x65, 0x72, 0x66, 0x6f, 0x72, 0x6d, 0x61,
	0x6e, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x31, 0x0a, 0x1b, 0x50, 0x45, 0x52, 0x46, 0x4f,
	0x52, 0x4d, 0x41, 0x4e, 0x43, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x41, 0x54, 0x54, 0x45,
	0x4e, 0x44, 0x41, 0x4e, 0x43, 0x45, 0x10, 0x00, 0x1a, 0x10, 0xaa, 0xd1, 0xf9, 0xd6, 0x03, 0x0a,
	0x41, 0x54, 0x54, 0x45, 0x4e, 0x44, 0x41, 0x4e, 0x43, 0x45, 0x12, 0x2b, 0x0a, 0x18, 0x50, 0x45,
	0x52, 0x46, 0x4f, 0x52, 0x4d, 0x41, 0x4e, 0x43, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4d,
	0x49, 0x44, 0x54, 0x45, 0x52, 0x4d, 0x10, 0x01, 0x1a, 0x0d, 0xaa, 0xd1, 0xf9, 0xd6, 0x03, 0x07,
	0x4d, 0x49, 0x44, 0x54, 0x45, 0x52, 0x4d, 0x32, 0x51, 0x0a, 0x0a, 0x53, 0x74, 0x75, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x43, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x53, 0x74, 0x75, 0x64,
	0x65, 0x6e, 0x74, 0x12, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x53,
	0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x3a, 0x45, 0x0a, 0x0b, 0x63, 0x75,
	0x73, 0x74, 0x6f, 0x6d, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6e, 0x75, 0x6d,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x95, 0x9a, 0xef,
	0x3a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x4e, 0x61, 0x6d,
	0x65, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2e, 0x2f, 0x70, 0x62, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_main_proto_rawDescOnce sync.Once
	file_main_proto_rawDescData = file_main_proto_rawDesc
)

func file_main_proto_rawDescGZIP() []byte {
	file_main_proto_rawDescOnce.Do(func() {
		file_main_proto_rawDescData = protoimpl.X.CompressGZIP(file_main_proto_rawDescData)
	})
	return file_main_proto_rawDescData
}

var file_main_proto_enumTypes = make([]protoimpl.EnumInfo, 4)
var file_main_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_main_proto_goTypes = []interface{}{
	(Gender)(0),                         // 0: proto.Gender
	(Subject)(0),                        // 1: proto.Subject
	(DefaultRemark)(0),                  // 2: proto.DefaultRemark
	(PerformanceType)(0),                // 3: proto.PerformanceType
	(*GetStudentRequest)(nil),           // 4: proto.GetStudentRequest
	(*GetStudentResponse)(nil),          // 5: proto.GetStudentResponse
	(*Student)(nil),                     // 6: proto.Student
	(*Score)(nil),                       // 7: proto.Score
	(*Homework)(nil),                    // 8: proto.Homework
	nil,                                 // 9: proto.Student.UsualPerformancesEntry
	(*Address)(nil),                     // 10: address.Address
	(*any1.Any)(nil),                    // 11: google.protobuf.Any
	(*timestamp.Timestamp)(nil),         // 12: google.protobuf.Timestamp
	(*descriptor.EnumValueOptions)(nil), // 13: google.protobuf.EnumValueOptions
}
var file_main_proto_depIdxs = []int32{
	6,  // 0: proto.GetStudentResponse.student:type_name -> proto.Student
	0,  // 1: proto.Student.gender:type_name -> proto.Gender
	7,  // 2: proto.Student.scores:type_name -> proto.Score
	8,  // 3: proto.Student.homeworks:type_name -> proto.Homework
	10, // 4: proto.Student.address:type_name -> address.Address
	2,  // 5: proto.Student.default_remark:type_name -> proto.DefaultRemark
	11, // 6: proto.Student.special:type_name -> google.protobuf.Any
	9,  // 7: proto.Student.usual_performances:type_name -> proto.Student.UsualPerformancesEntry
	12, // 8: proto.Student.last_updated:type_name -> google.protobuf.Timestamp
	1,  // 9: proto.Score.subject:type_name -> proto.Subject
	13, // 10: proto.custom_name:extendee -> google.protobuf.EnumValueOptions
	4,  // 11: proto.StuService.GetStudent:input_type -> proto.GetStudentRequest
	5,  // 12: proto.StuService.GetStudent:output_type -> proto.GetStudentResponse
	12, // [12:13] is the sub-list for method output_type
	11, // [11:12] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	10, // [10:11] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_main_proto_init() }
func file_main_proto_init() {
	if File_main_proto != nil {
		return
	}
	file_address_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_main_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetStudentRequest); i {
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
		file_main_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetStudentResponse); i {
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
		file_main_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Student); i {
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
		file_main_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Score); i {
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
		file_main_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Homework); i {
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
	file_main_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*Student_DefaultRemark)(nil),
		(*Student_OtherRemark)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_main_proto_rawDesc,
			NumEnums:      4,
			NumMessages:   6,
			NumExtensions: 1,
			NumServices:   1,
		},
		GoTypes:           file_main_proto_goTypes,
		DependencyIndexes: file_main_proto_depIdxs,
		EnumInfos:         file_main_proto_enumTypes,
		MessageInfos:      file_main_proto_msgTypes,
		ExtensionInfos:    file_main_proto_extTypes,
	}.Build()
	File_main_proto = out.File
	file_main_proto_rawDesc = nil
	file_main_proto_goTypes = nil
	file_main_proto_depIdxs = nil
}
