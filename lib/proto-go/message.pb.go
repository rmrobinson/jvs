// Code generated by protoc-gen-go.
// source: message.proto
// DO NOT EDIT!

package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Message struct {
	Id       string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	User     string `protobuf:"bytes,2,opt,name=user" json:"user,omitempty"`
	Contents string `protobuf:"bytes,3,opt,name=contents" json:"contents,omitempty"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto1.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

type SendMessageRequest struct {
	Message *Message `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
}

func (m *SendMessageRequest) Reset()                    { *m = SendMessageRequest{} }
func (m *SendMessageRequest) String() string            { return proto1.CompactTextString(m) }
func (*SendMessageRequest) ProtoMessage()               {}
func (*SendMessageRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func (m *SendMessageRequest) GetMessage() *Message {
	if m != nil {
		return m.Message
	}
	return nil
}

type SendMessageResponse struct {
	Error   string   `protobuf:"bytes,1,opt,name=error" json:"error,omitempty"`
	Message *Message `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
}

func (m *SendMessageResponse) Reset()                    { *m = SendMessageResponse{} }
func (m *SendMessageResponse) String() string            { return proto1.CompactTextString(m) }
func (*SendMessageResponse) ProtoMessage()               {}
func (*SendMessageResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{2} }

func (m *SendMessageResponse) GetMessage() *Message {
	if m != nil {
		return m.Message
	}
	return nil
}

type WatchMessagesRequest struct {
}

func (m *WatchMessagesRequest) Reset()                    { *m = WatchMessagesRequest{} }
func (m *WatchMessagesRequest) String() string            { return proto1.CompactTextString(m) }
func (*WatchMessagesRequest) ProtoMessage()               {}
func (*WatchMessagesRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{3} }

type WatchMessagesResponse struct {
	Message *Message `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
}

func (m *WatchMessagesResponse) Reset()                    { *m = WatchMessagesResponse{} }
func (m *WatchMessagesResponse) String() string            { return proto1.CompactTextString(m) }
func (*WatchMessagesResponse) ProtoMessage()               {}
func (*WatchMessagesResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{4} }

func (m *WatchMessagesResponse) GetMessage() *Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func init() {
	proto1.RegisterType((*Message)(nil), "proto.Message")
	proto1.RegisterType((*SendMessageRequest)(nil), "proto.SendMessageRequest")
	proto1.RegisterType((*SendMessageResponse)(nil), "proto.SendMessageResponse")
	proto1.RegisterType((*WatchMessagesRequest)(nil), "proto.WatchMessagesRequest")
	proto1.RegisterType((*WatchMessagesResponse)(nil), "proto.WatchMessagesResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for MessageManager service

type MessageManagerClient interface {
	SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error)
	WatchMessages(ctx context.Context, in *WatchMessagesRequest, opts ...grpc.CallOption) (MessageManager_WatchMessagesClient, error)
}

type messageManagerClient struct {
	cc *grpc.ClientConn
}

func NewMessageManagerClient(cc *grpc.ClientConn) MessageManagerClient {
	return &messageManagerClient{cc}
}

func (c *messageManagerClient) SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error) {
	out := new(SendMessageResponse)
	err := grpc.Invoke(ctx, "/proto.MessageManager/SendMessage", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageManagerClient) WatchMessages(ctx context.Context, in *WatchMessagesRequest, opts ...grpc.CallOption) (MessageManager_WatchMessagesClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_MessageManager_serviceDesc.Streams[0], c.cc, "/proto.MessageManager/WatchMessages", opts...)
	if err != nil {
		return nil, err
	}
	x := &messageManagerWatchMessagesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MessageManager_WatchMessagesClient interface {
	Recv() (*WatchMessagesResponse, error)
	grpc.ClientStream
}

type messageManagerWatchMessagesClient struct {
	grpc.ClientStream
}

func (x *messageManagerWatchMessagesClient) Recv() (*WatchMessagesResponse, error) {
	m := new(WatchMessagesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for MessageManager service

type MessageManagerServer interface {
	SendMessage(context.Context, *SendMessageRequest) (*SendMessageResponse, error)
	WatchMessages(*WatchMessagesRequest, MessageManager_WatchMessagesServer) error
}

func RegisterMessageManagerServer(s *grpc.Server, srv MessageManagerServer) {
	s.RegisterService(&_MessageManager_serviceDesc, srv)
}

func _MessageManager_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageManagerServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.MessageManager/SendMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageManagerServer).SendMessage(ctx, req.(*SendMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageManager_WatchMessages_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WatchMessagesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MessageManagerServer).WatchMessages(m, &messageManagerWatchMessagesServer{stream})
}

type MessageManager_WatchMessagesServer interface {
	Send(*WatchMessagesResponse) error
	grpc.ServerStream
}

type messageManagerWatchMessagesServer struct {
	grpc.ServerStream
}

func (x *messageManagerWatchMessagesServer) Send(m *WatchMessagesResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _MessageManager_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.MessageManager",
	HandlerType: (*MessageManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMessage",
			Handler:    _MessageManager_SendMessage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "WatchMessages",
			Handler:       _MessageManager_WatchMessages_Handler,
			ServerStreams: true,
		},
	},
	Metadata: fileDescriptor2,
}

func init() { proto1.RegisterFile("message.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 245 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x50, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0xc5, 0x81, 0x52, 0xb8, 0xaa, 0x19, 0x8e, 0x82, 0x82, 0x61, 0x40, 0x9e, 0x3a, 0x55, 0xa8,
	0xec, 0x48, 0x2c, 0x48, 0x0c, 0x65, 0x28, 0x42, 0xcc, 0xa1, 0x3d, 0x01, 0x03, 0x76, 0xf0, 0x39,
	0xff, 0xc4, 0x67, 0x12, 0xd9, 0x17, 0x94, 0x40, 0x84, 0x58, 0x92, 0xbb, 0xf7, 0xde, 0xbd, 0x7b,
	0x67, 0x98, 0xbe, 0x13, 0x73, 0xf9, 0x42, 0x8b, 0xca, 0xbb, 0xe0, 0x70, 0x14, 0x7f, 0xe6, 0x0e,
	0xc6, 0xab, 0x84, 0x63, 0x0e, 0xd9, 0xdb, 0xb6, 0x50, 0x17, 0x6a, 0x7e, 0xb8, 0x6e, 0x2a, 0x44,
	0xd8, 0xab, 0x99, 0x7c, 0x91, 0x45, 0x24, 0xd6, 0xa8, 0xe1, 0x60, 0xe3, 0x6c, 0x20, 0x1b, 0xb8,
	0xd8, 0x8d, 0xf8, 0x77, 0x6f, 0xae, 0x01, 0x1f, 0xc8, 0x6e, 0xc5, 0x6e, 0x4d, 0x1f, 0x35, 0x71,
	0xc0, 0x39, 0x8c, 0x65, 0x71, 0xb4, 0x9e, 0x2c, 0xf3, 0x14, 0x60, 0xd1, 0xea, 0x5a, 0xda, 0x3c,
	0xc2, 0x51, 0x6f, 0x9e, 0x2b, 0x67, 0x99, 0x70, 0x06, 0x23, 0xf2, 0xde, 0x79, 0x49, 0x96, 0x9a,
	0xae, 0x6d, 0xf6, 0xb7, 0xed, 0x09, 0xcc, 0x9e, 0xca, 0xb0, 0x79, 0x15, 0x82, 0x25, 0x98, 0xb9,
	0x81, 0xe3, 0x1f, 0xb8, 0x2c, 0xfc, 0x77, 0xe2, 0xe5, 0xa7, 0x82, 0x5c, 0xc0, 0x55, 0x69, 0x9b,
	0xaf, 0xc7, 0x5b, 0x98, 0x74, 0x8e, 0xc0, 0x53, 0x19, 0xfd, 0xfd, 0x30, 0x5a, 0x0f, 0x51, 0x29,
	0x82, 0xd9, 0xc1, 0x7b, 0x98, 0xf6, 0xd2, 0xe1, 0x99, 0xc8, 0x87, 0x6e, 0xd1, 0xe7, 0xc3, 0x64,
	0xeb, 0x76, 0xa9, 0x9e, 0xf7, 0xa3, 0xe0, 0xea, 0x2b, 0x00, 0x00, 0xff, 0xff, 0x02, 0xa0, 0x37,
	0x5a, 0x06, 0x02, 0x00, 0x00,
}