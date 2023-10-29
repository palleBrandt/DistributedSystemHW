// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.3
// source: proto/template.proto

package proto

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

// ChittyChatClient is the client API for ChittyChat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChittyChatClient interface {
	Chat(ctx context.Context, opts ...grpc.CallOption) (ChittyChat_ChatClient, error)
	Join(ctx context.Context, in *Client, opts ...grpc.CallOption) (ChittyChat_JoinClient, error)
}

type chittyChatClient struct {
	cc grpc.ClientConnInterface
}

func NewChittyChatClient(cc grpc.ClientConnInterface) ChittyChatClient {
	return &chittyChatClient{cc}
}

func (c *chittyChatClient) Chat(ctx context.Context, opts ...grpc.CallOption) (ChittyChat_ChatClient, error) {
	stream, err := c.cc.NewStream(ctx, &ChittyChat_ServiceDesc.Streams[0], "/proto.ChittyChat/Chat", opts...)
	if err != nil {
		return nil, err
	}
	x := &chittyChatChatClient{stream}
	return x, nil
}

type ChittyChat_ChatClient interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ClientStream
}

type chittyChatChatClient struct {
	grpc.ClientStream
}

func (x *chittyChatChatClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *chittyChatChatClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *chittyChatClient) Join(ctx context.Context, in *Client, opts ...grpc.CallOption) (ChittyChat_JoinClient, error) {
	stream, err := c.cc.NewStream(ctx, &ChittyChat_ServiceDesc.Streams[1], "/proto.ChittyChat/Join", opts...)
	if err != nil {
		return nil, err
	}
	x := &chittyChatJoinClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ChittyChat_JoinClient interface {
	Recv() (*Message, error)
	grpc.ClientStream
}

type chittyChatJoinClient struct {
	grpc.ClientStream
}

func (x *chittyChatJoinClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ChittyChatServer is the server API for ChittyChat service.
// All implementations must embed UnimplementedChittyChatServer
// for forward compatibility
type ChittyChatServer interface {
	Chat(ChittyChat_ChatServer) error
	Join(*Client, ChittyChat_JoinServer) error
	mustEmbedUnimplementedChittyChatServer()
}

// UnimplementedChittyChatServer must be embedded to have forward compatible implementations.
type UnimplementedChittyChatServer struct {
}

func (UnimplementedChittyChatServer) Chat(ChittyChat_ChatServer) error {
	return status.Errorf(codes.Unimplemented, "method Chat not implemented")
}
func (UnimplementedChittyChatServer) Join(*Client, ChittyChat_JoinServer) error {
	return status.Errorf(codes.Unimplemented, "method Join not implemented")
}
func (UnimplementedChittyChatServer) mustEmbedUnimplementedChittyChatServer() {}

// UnsafeChittyChatServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChittyChatServer will
// result in compilation errors.
type UnsafeChittyChatServer interface {
	mustEmbedUnimplementedChittyChatServer()
}

func RegisterChittyChatServer(s grpc.ServiceRegistrar, srv ChittyChatServer) {
	s.RegisterService(&ChittyChat_ServiceDesc, srv)
}

func _ChittyChat_Chat_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChittyChatServer).Chat(&chittyChatChatServer{stream})
}

type ChittyChat_ChatServer interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type chittyChatChatServer struct {
	grpc.ServerStream
}

func (x *chittyChatChatServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *chittyChatChatServer) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ChittyChat_Join_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Client)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChittyChatServer).Join(m, &chittyChatJoinServer{stream})
}

type ChittyChat_JoinServer interface {
	Send(*Message) error
	grpc.ServerStream
}

type chittyChatJoinServer struct {
	grpc.ServerStream
}

func (x *chittyChatJoinServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

// ChittyChat_ServiceDesc is the grpc.ServiceDesc for ChittyChat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChittyChat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ChittyChat",
	HandlerType: (*ChittyChatServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Chat",
			Handler:       _ChittyChat_Chat_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Join",
			Handler:       _ChittyChat_Join_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/template.proto",
}
