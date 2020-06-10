// Code generated by protoc-gen-go. DO NOT EDIT.
// source: messages.proto

/*
Package messages is a generated protocol buffer package.

It is generated from these files:
	messages.proto

It has these top-level messages:
	TaskDef
	Registration
	Ack
	Empty
	TaskStatusMsg
*/
package messages

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type TaskDef struct {
	RequestID  string   `protobuf:"bytes,1,opt,name=RequestID" json:"RequestID,omitempty"`
	Executable string   `protobuf:"bytes,2,opt,name=Executable" json:"Executable,omitempty"`
	Args       []string `protobuf:"bytes,3,rep,name=Args" json:"Args,omitempty"`
	Env        []string `protobuf:"bytes,4,rep,name=Env" json:"Env,omitempty"`
	Dir        string   `protobuf:"bytes,5,opt,name=Dir" json:"Dir,omitempty"`
}

func (m *TaskDef) Reset()                    { *m = TaskDef{} }
func (m *TaskDef) String() string            { return proto.CompactTextString(m) }
func (*TaskDef) ProtoMessage()               {}
func (*TaskDef) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *TaskDef) GetRequestID() string {
	if m != nil {
		return m.RequestID
	}
	return ""
}

func (m *TaskDef) GetExecutable() string {
	if m != nil {
		return m.Executable
	}
	return ""
}

func (m *TaskDef) GetArgs() []string {
	if m != nil {
		return m.Args
	}
	return nil
}

func (m *TaskDef) GetEnv() []string {
	if m != nil {
		return m.Env
	}
	return nil
}

func (m *TaskDef) GetDir() string {
	if m != nil {
		return m.Dir
	}
	return ""
}

type Registration struct {
	Ip   string `protobuf:"bytes,1,opt,name=ip" json:"ip,omitempty"`
	Port string `protobuf:"bytes,2,opt,name=port" json:"port,omitempty"`
}

func (m *Registration) Reset()                    { *m = Registration{} }
func (m *Registration) String() string            { return proto.CompactTextString(m) }
func (*Registration) ProtoMessage()               {}
func (*Registration) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Registration) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func (m *Registration) GetPort() string {
	if m != nil {
		return m.Port
	}
	return ""
}

type Ack struct {
	Ok bool `protobuf:"varint,1,opt,name=ok" json:"ok,omitempty"`
}

func (m *Ack) Reset()                    { *m = Ack{} }
func (m *Ack) String() string            { return proto.CompactTextString(m) }
func (*Ack) ProtoMessage()               {}
func (*Ack) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Ack) GetOk() bool {
	if m != nil {
		return m.Ok
	}
	return false
}

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type TaskStatusMsg struct {
	Status int32 `protobuf:"varint,1,opt,name=status" json:"status,omitempty"`
}

func (m *TaskStatusMsg) Reset()                    { *m = TaskStatusMsg{} }
func (m *TaskStatusMsg) String() string            { return proto.CompactTextString(m) }
func (*TaskStatusMsg) ProtoMessage()               {}
func (*TaskStatusMsg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *TaskStatusMsg) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func init() {
	proto.RegisterType((*TaskDef)(nil), "TaskDef")
	proto.RegisterType((*Registration)(nil), "Registration")
	proto.RegisterType((*Ack)(nil), "Ack")
	proto.RegisterType((*Empty)(nil), "Empty")
	proto.RegisterType((*TaskStatusMsg)(nil), "TaskStatusMsg")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for CronyTaskMaster service

type CronyTaskMasterClient interface {
	RunTask(ctx context.Context, in *TaskDef, opts ...grpc.CallOption) (*Ack, error)
	Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Ack, error)
}

type cronyTaskMasterClient struct {
	cc *grpc.ClientConn
}

func NewCronyTaskMasterClient(cc *grpc.ClientConn) CronyTaskMasterClient {
	return &cronyTaskMasterClient{cc}
}

func (c *cronyTaskMasterClient) RunTask(ctx context.Context, in *TaskDef, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := grpc.Invoke(ctx, "/CronyTaskMaster/RunTask", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cronyTaskMasterClient) Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := grpc.Invoke(ctx, "/CronyTaskMaster/Ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for CronyTaskMaster service

type CronyTaskMasterServer interface {
	RunTask(context.Context, *TaskDef) (*Ack, error)
	Ping(context.Context, *Empty) (*Ack, error)
}

func RegisterCronyTaskMasterServer(s *grpc.Server, srv CronyTaskMasterServer) {
	s.RegisterService(&_CronyTaskMaster_serviceDesc, srv)
}

func _CronyTaskMaster_RunTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskDef)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CronyTaskMasterServer).RunTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/CronyTaskMaster/RunTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CronyTaskMasterServer).RunTask(ctx, req.(*TaskDef))
	}
	return interceptor(ctx, in, info, handler)
}

func _CronyTaskMaster_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CronyTaskMasterServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/CronyTaskMaster/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CronyTaskMasterServer).Ping(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _CronyTaskMaster_serviceDesc = grpc.ServiceDesc{
	ServiceName: "CronyTaskMaster",
	HandlerType: (*CronyTaskMasterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RunTask",
			Handler:    _CronyTaskMaster_RunTask_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _CronyTaskMaster_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "messages.proto",
}

// Client API for CronyTaskSlave service

type CronyTaskSlaveClient interface {
	Register(ctx context.Context, in *Registration, opts ...grpc.CallOption) (*Empty, error)
	TaskStatus(ctx context.Context, in *TaskStatusMsg, opts ...grpc.CallOption) (*Ack, error)
}

type cronyTaskSlaveClient struct {
	cc *grpc.ClientConn
}

func NewCronyTaskSlaveClient(cc *grpc.ClientConn) CronyTaskSlaveClient {
	return &cronyTaskSlaveClient{cc}
}

func (c *cronyTaskSlaveClient) Register(ctx context.Context, in *Registration, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/CronyTaskSlave/Register", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cronyTaskSlaveClient) TaskStatus(ctx context.Context, in *TaskStatusMsg, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := grpc.Invoke(ctx, "/CronyTaskSlave/TaskStatus", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for CronyTaskSlave service

type CronyTaskSlaveServer interface {
	Register(context.Context, *Registration) (*Empty, error)
	TaskStatus(context.Context, *TaskStatusMsg) (*Ack, error)
}

func RegisterCronyTaskSlaveServer(s *grpc.Server, srv CronyTaskSlaveServer) {
	s.RegisterService(&_CronyTaskSlave_serviceDesc, srv)
}

func _CronyTaskSlave_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Registration)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CronyTaskSlaveServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/CronyTaskSlave/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CronyTaskSlaveServer).Register(ctx, req.(*Registration))
	}
	return interceptor(ctx, in, info, handler)
}

func _CronyTaskSlave_TaskStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskStatusMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CronyTaskSlaveServer).TaskStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/CronyTaskSlave/TaskStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CronyTaskSlaveServer).TaskStatus(ctx, req.(*TaskStatusMsg))
	}
	return interceptor(ctx, in, info, handler)
}

var _CronyTaskSlave_serviceDesc = grpc.ServiceDesc{
	ServiceName: "CronyTaskSlave",
	HandlerType: (*CronyTaskSlaveServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _CronyTaskSlave_Register_Handler,
		},
		{
			MethodName: "TaskStatus",
			Handler:    _CronyTaskSlave_TaskStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "messages.proto",
}

func init() { proto.RegisterFile("messages.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 307 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x91, 0x5d, 0x4b, 0xfb, 0x30,
	0x14, 0xc6, 0xb7, 0xb6, 0x7b, 0x3b, 0xfc, 0xd7, 0xbf, 0x1c, 0x70, 0x94, 0x29, 0x32, 0xa2, 0xe0,
	0xae, 0x72, 0x31, 0x3f, 0xc1, 0xb0, 0x15, 0xbc, 0x18, 0x48, 0xe6, 0x9d, 0x57, 0x59, 0x89, 0x25,
	0x74, 0x6b, 0x6a, 0x92, 0x0e, 0x77, 0xe5, 0x57, 0x97, 0xc4, 0xce, 0xcd, 0xbb, 0xe7, 0xfc, 0xfa,
	0x9c, 0x97, 0xa7, 0x81, 0x78, 0x27, 0x8c, 0xe1, 0x85, 0x30, 0xb4, 0xd6, 0xca, 0x2a, 0xf2, 0x05,
	0x83, 0x57, 0x6e, 0xca, 0x54, 0xbc, 0xe3, 0x35, 0x8c, 0x98, 0xf8, 0x68, 0x84, 0xb1, 0xcf, 0x69,
	0xd2, 0x9d, 0x75, 0xe7, 0x23, 0x76, 0x02, 0x78, 0x03, 0x90, 0x7d, 0x8a, 0xbc, 0xb1, 0x7c, 0xb3,
	0x15, 0x49, 0xe0, 0x3f, 0x9f, 0x11, 0x44, 0x88, 0x96, 0xba, 0x30, 0x49, 0x38, 0x0b, 0xe7, 0x23,
	0xe6, 0x35, 0x5e, 0x40, 0x98, 0x55, 0xfb, 0x24, 0xf2, 0xc8, 0x49, 0x47, 0x52, 0xa9, 0x93, 0x9e,
	0x6f, 0x77, 0x92, 0x2c, 0xe0, 0x1f, 0x13, 0x85, 0x34, 0x56, 0x73, 0x2b, 0x55, 0x85, 0x31, 0x04,
	0xb2, 0x6e, 0xd7, 0x07, 0xb2, 0x76, 0x73, 0x6b, 0xa5, 0x6d, 0xbb, 0xd1, 0x6b, 0x72, 0x09, 0xe1,
	0x32, 0x2f, 0x9d, 0x55, 0x95, 0xde, 0x3a, 0x64, 0x81, 0x2a, 0xc9, 0x00, 0x7a, 0xd9, 0xae, 0xb6,
	0x07, 0x72, 0x0f, 0x63, 0x17, 0x6a, 0x6d, 0xb9, 0x6d, 0xcc, 0xca, 0x14, 0x38, 0x81, 0xbe, 0xf1,
	0x85, 0x77, 0xf7, 0x58, 0x5b, 0x2d, 0x9e, 0xe0, 0xff, 0xa3, 0x56, 0xd5, 0xc1, 0xb9, 0x57, 0xdc,
	0x58, 0xa1, 0xf1, 0x0a, 0x06, 0xac, 0xa9, 0x1c, 0xc0, 0x21, 0x6d, 0x7f, 0xcd, 0x34, 0xa2, 0xcb,
	0xbc, 0x24, 0x1d, 0x9c, 0x40, 0xf4, 0x22, 0xab, 0x02, 0xfb, 0xd4, 0x2f, 0x3a, 0xf2, 0xc5, 0x1b,
	0xc4, 0xbf, 0x73, 0xd6, 0x5b, 0xbe, 0x17, 0x78, 0x0b, 0xc3, 0x9f, 0x58, 0x42, 0xe3, 0x98, 0x9e,
	0x27, 0x9c, 0xb6, 0xcd, 0xa4, 0x83, 0x77, 0x00, 0xa7, 0x3b, 0x31, 0xa6, 0x7f, 0x8e, 0x3e, 0x0e,
	0xdf, 0xf4, 0xfd, 0x4b, 0x3d, 0x7c, 0x07, 0x00, 0x00, 0xff, 0xff, 0x73, 0x3a, 0x1a, 0xe1, 0xbb,
	0x01, 0x00, 0x00,
}
