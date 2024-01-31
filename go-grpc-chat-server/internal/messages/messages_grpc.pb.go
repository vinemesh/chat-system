// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: messages.proto

package messages

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	MessageService_SendMessage_FullMethodName    = "/messages.MessageService/SendMessage"
	MessageService_StreamMessages_FullMethodName = "/messages.MessageService/StreamMessages"
)

// MessageServiceClient is the client API for MessageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessageServiceClient interface {
	SendMessage(ctx context.Context, in *PlayerMessage, opts ...grpc.CallOption) (*ResponseMessage, error)
	StreamMessages(ctx context.Context, opts ...grpc.CallOption) (MessageService_StreamMessagesClient, error)
}

type messageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageServiceClient(cc grpc.ClientConnInterface) MessageServiceClient {
	return &messageServiceClient{cc}
}

func (c *messageServiceClient) SendMessage(ctx context.Context, in *PlayerMessage, opts ...grpc.CallOption) (*ResponseMessage, error) {
	out := new(ResponseMessage)
	err := c.cc.Invoke(ctx, MessageService_SendMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageServiceClient) StreamMessages(ctx context.Context, opts ...grpc.CallOption) (MessageService_StreamMessagesClient, error) {
	stream, err := c.cc.NewStream(ctx, &MessageService_ServiceDesc.Streams[0], MessageService_StreamMessages_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &messageServiceStreamMessagesClient{stream}
	return x, nil
}

type MessageService_StreamMessagesClient interface {
	Send(*PlayerMessage) error
	Recv() (*ResponseMessage, error)
	grpc.ClientStream
}

type messageServiceStreamMessagesClient struct {
	grpc.ClientStream
}

func (x *messageServiceStreamMessagesClient) Send(m *PlayerMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *messageServiceStreamMessagesClient) Recv() (*ResponseMessage, error) {
	m := new(ResponseMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MessageServiceServer is the server API for MessageService service.
// All implementations must embed UnimplementedMessageServiceServer
// for forward compatibility
type MessageServiceServer interface {
	SendMessage(context.Context, *PlayerMessage) (*ResponseMessage, error)
	StreamMessages(MessageService_StreamMessagesServer) error
	mustEmbedUnimplementedMessageServiceServer()
}

// UnimplementedMessageServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMessageServiceServer struct {
}

func (UnimplementedMessageServiceServer) SendMessage(context.Context, *PlayerMessage) (*ResponseMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedMessageServiceServer) StreamMessages(MessageService_StreamMessagesServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamMessages not implemented")
}
func (UnimplementedMessageServiceServer) mustEmbedUnimplementedMessageServiceServer() {}

// UnsafeMessageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessageServiceServer will
// result in compilation errors.
type UnsafeMessageServiceServer interface {
	mustEmbedUnimplementedMessageServiceServer()
}

func RegisterMessageServiceServer(s grpc.ServiceRegistrar, srv MessageServiceServer) {
	s.RegisterService(&MessageService_ServiceDesc, srv)
}

func _MessageService_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlayerMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MessageService_SendMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).SendMessage(ctx, req.(*PlayerMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageService_StreamMessages_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MessageServiceServer).StreamMessages(&messageServiceStreamMessagesServer{stream})
}

type MessageService_StreamMessagesServer interface {
	Send(*ResponseMessage) error
	Recv() (*PlayerMessage, error)
	grpc.ServerStream
}

type messageServiceStreamMessagesServer struct {
	grpc.ServerStream
}

func (x *messageServiceStreamMessagesServer) Send(m *ResponseMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *messageServiceStreamMessagesServer) Recv() (*PlayerMessage, error) {
	m := new(PlayerMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MessageService_ServiceDesc is the grpc.ServiceDesc for MessageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MessageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "messages.MessageService",
	HandlerType: (*MessageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMessage",
			Handler:    _MessageService_SendMessage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamMessages",
			Handler:       _MessageService_StreamMessages_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "messages.proto",
}
