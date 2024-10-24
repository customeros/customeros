// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: event_completion.proto

package event_completion_grpc_service

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// EventCompletionGrpcServiceClient is the client API for EventCompletionGrpcService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventCompletionGrpcServiceClient interface {
	NotifyEventProcessed(ctx context.Context, in *NotifyEventProcessedRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type eventCompletionGrpcServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEventCompletionGrpcServiceClient(cc grpc.ClientConnInterface) EventCompletionGrpcServiceClient {
	return &eventCompletionGrpcServiceClient{cc}
}

func (c *eventCompletionGrpcServiceClient) NotifyEventProcessed(ctx context.Context, in *NotifyEventProcessedRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/EventCompletionGrpcService/NotifyEventProcessed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventCompletionGrpcServiceServer is the server API for EventCompletionGrpcService service.
// All implementations should embed UnimplementedEventCompletionGrpcServiceServer
// for forward compatibility
type EventCompletionGrpcServiceServer interface {
	NotifyEventProcessed(context.Context, *NotifyEventProcessedRequest) (*emptypb.Empty, error)
}

// UnimplementedEventCompletionGrpcServiceServer should be embedded to have forward compatible implementations.
type UnimplementedEventCompletionGrpcServiceServer struct {
}

func (UnimplementedEventCompletionGrpcServiceServer) NotifyEventProcessed(context.Context, *NotifyEventProcessedRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyEventProcessed not implemented")
}

// UnsafeEventCompletionGrpcServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventCompletionGrpcServiceServer will
// result in compilation errors.
type UnsafeEventCompletionGrpcServiceServer interface {
	mustEmbedUnimplementedEventCompletionGrpcServiceServer()
}

func RegisterEventCompletionGrpcServiceServer(s grpc.ServiceRegistrar, srv EventCompletionGrpcServiceServer) {
	s.RegisterService(&EventCompletionGrpcService_ServiceDesc, srv)
}

func _EventCompletionGrpcService_NotifyEventProcessed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotifyEventProcessedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventCompletionGrpcServiceServer).NotifyEventProcessed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/EventCompletionGrpcService/NotifyEventProcessed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventCompletionGrpcServiceServer).NotifyEventProcessed(ctx, req.(*NotifyEventProcessedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EventCompletionGrpcService_ServiceDesc is the grpc.ServiceDesc for EventCompletionGrpcService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EventCompletionGrpcService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "EventCompletionGrpcService",
	HandlerType: (*EventCompletionGrpcServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NotifyEventProcessed",
			Handler:    _EventCompletionGrpcService_NotifyEventProcessed_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "event_completion.proto",
}
