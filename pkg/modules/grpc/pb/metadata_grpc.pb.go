// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.13.0
// source: github.com/dipdup-net/abi-indexer/pkg/modules/grpc/proto/metadata.proto

package pb

import (
	context "context"
	pb "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MetadataServiceClient is the client API for MetadataService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetadataServiceClient interface {
	SubscribeOnMetadata(ctx context.Context, in *pb.DefaultRequest, opts ...grpc.CallOption) (MetadataService_SubscribeOnMetadataClient, error)
	UnsubscribeFromMetadata(ctx context.Context, in *pb.UnsubscribeRequest, opts ...grpc.CallOption) (*pb.UnsubscribeResponse, error)
	GetMetadata(ctx context.Context, in *GetMetadataRequest, opts ...grpc.CallOption) (*Metadata, error)
	ListMetadata(ctx context.Context, in *ListMetadataRequest, opts ...grpc.CallOption) (*ListMetadataResponse, error)
	GetMetadataByMethodSinature(ctx context.Context, in *GetMetadataByMethodSinatureRequest, opts ...grpc.CallOption) (*ListMetadataResponse, error)
	GetMetadataByTopic(ctx context.Context, in *GetMetadataByTopicRequest, opts ...grpc.CallOption) (*ListMetadataResponse, error)
}

type metadataServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMetadataServiceClient(cc grpc.ClientConnInterface) MetadataServiceClient {
	return &metadataServiceClient{cc}
}

func (c *metadataServiceClient) SubscribeOnMetadata(ctx context.Context, in *pb.DefaultRequest, opts ...grpc.CallOption) (MetadataService_SubscribeOnMetadataClient, error) {
	stream, err := c.cc.NewStream(ctx, &MetadataService_ServiceDesc.Streams[0], "/proto.MetadataService/SubscribeOnMetadata", opts...)
	if err != nil {
		return nil, err
	}
	x := &metadataServiceSubscribeOnMetadataClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MetadataService_SubscribeOnMetadataClient interface {
	Recv() (*SubscriptionMetadata, error)
	grpc.ClientStream
}

type metadataServiceSubscribeOnMetadataClient struct {
	grpc.ClientStream
}

func (x *metadataServiceSubscribeOnMetadataClient) Recv() (*SubscriptionMetadata, error) {
	m := new(SubscriptionMetadata)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *metadataServiceClient) UnsubscribeFromMetadata(ctx context.Context, in *pb.UnsubscribeRequest, opts ...grpc.CallOption) (*pb.UnsubscribeResponse, error) {
	out := new(pb.UnsubscribeResponse)
	err := c.cc.Invoke(ctx, "/proto.MetadataService/UnsubscribeFromMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metadataServiceClient) GetMetadata(ctx context.Context, in *GetMetadataRequest, opts ...grpc.CallOption) (*Metadata, error) {
	out := new(Metadata)
	err := c.cc.Invoke(ctx, "/proto.MetadataService/GetMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metadataServiceClient) ListMetadata(ctx context.Context, in *ListMetadataRequest, opts ...grpc.CallOption) (*ListMetadataResponse, error) {
	out := new(ListMetadataResponse)
	err := c.cc.Invoke(ctx, "/proto.MetadataService/ListMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metadataServiceClient) GetMetadataByMethodSinature(ctx context.Context, in *GetMetadataByMethodSinatureRequest, opts ...grpc.CallOption) (*ListMetadataResponse, error) {
	out := new(ListMetadataResponse)
	err := c.cc.Invoke(ctx, "/proto.MetadataService/GetMetadataByMethodSinature", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metadataServiceClient) GetMetadataByTopic(ctx context.Context, in *GetMetadataByTopicRequest, opts ...grpc.CallOption) (*ListMetadataResponse, error) {
	out := new(ListMetadataResponse)
	err := c.cc.Invoke(ctx, "/proto.MetadataService/GetMetadataByTopic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetadataServiceServer is the server API for MetadataService service.
// All implementations must embed UnimplementedMetadataServiceServer
// for forward compatibility
type MetadataServiceServer interface {
	SubscribeOnMetadata(*pb.DefaultRequest, MetadataService_SubscribeOnMetadataServer) error
	UnsubscribeFromMetadata(context.Context, *pb.UnsubscribeRequest) (*pb.UnsubscribeResponse, error)
	GetMetadata(context.Context, *GetMetadataRequest) (*Metadata, error)
	ListMetadata(context.Context, *ListMetadataRequest) (*ListMetadataResponse, error)
	GetMetadataByMethodSinature(context.Context, *GetMetadataByMethodSinatureRequest) (*ListMetadataResponse, error)
	GetMetadataByTopic(context.Context, *GetMetadataByTopicRequest) (*ListMetadataResponse, error)
	mustEmbedUnimplementedMetadataServiceServer()
}

// UnimplementedMetadataServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMetadataServiceServer struct {
}

func (UnimplementedMetadataServiceServer) SubscribeOnMetadata(*pb.DefaultRequest, MetadataService_SubscribeOnMetadataServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeOnMetadata not implemented")
}
func (UnimplementedMetadataServiceServer) UnsubscribeFromMetadata(context.Context, *pb.UnsubscribeRequest) (*pb.UnsubscribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnsubscribeFromMetadata not implemented")
}
func (UnimplementedMetadataServiceServer) GetMetadata(context.Context, *GetMetadataRequest) (*Metadata, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetadata not implemented")
}
func (UnimplementedMetadataServiceServer) ListMetadata(context.Context, *ListMetadataRequest) (*ListMetadataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMetadata not implemented")
}
func (UnimplementedMetadataServiceServer) GetMetadataByMethodSinature(context.Context, *GetMetadataByMethodSinatureRequest) (*ListMetadataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetadataByMethodSinature not implemented")
}
func (UnimplementedMetadataServiceServer) GetMetadataByTopic(context.Context, *GetMetadataByTopicRequest) (*ListMetadataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetadataByTopic not implemented")
}
func (UnimplementedMetadataServiceServer) mustEmbedUnimplementedMetadataServiceServer() {}

// UnsafeMetadataServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetadataServiceServer will
// result in compilation errors.
type UnsafeMetadataServiceServer interface {
	mustEmbedUnimplementedMetadataServiceServer()
}

func RegisterMetadataServiceServer(s grpc.ServiceRegistrar, srv MetadataServiceServer) {
	s.RegisterService(&MetadataService_ServiceDesc, srv)
}

func _MetadataService_SubscribeOnMetadata_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(pb.DefaultRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MetadataServiceServer).SubscribeOnMetadata(m, &metadataServiceSubscribeOnMetadataServer{stream})
}

type MetadataService_SubscribeOnMetadataServer interface {
	Send(*SubscriptionMetadata) error
	grpc.ServerStream
}

type metadataServiceSubscribeOnMetadataServer struct {
	grpc.ServerStream
}

func (x *metadataServiceSubscribeOnMetadataServer) Send(m *SubscriptionMetadata) error {
	return x.ServerStream.SendMsg(m)
}

func _MetadataService_UnsubscribeFromMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(pb.UnsubscribeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetadataServiceServer).UnsubscribeFromMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.MetadataService/UnsubscribeFromMetadata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetadataServiceServer).UnsubscribeFromMetadata(ctx, req.(*pb.UnsubscribeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetadataService_GetMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMetadataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetadataServiceServer).GetMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.MetadataService/GetMetadata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetadataServiceServer).GetMetadata(ctx, req.(*GetMetadataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetadataService_ListMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListMetadataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetadataServiceServer).ListMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.MetadataService/ListMetadata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetadataServiceServer).ListMetadata(ctx, req.(*ListMetadataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetadataService_GetMetadataByMethodSinature_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMetadataByMethodSinatureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetadataServiceServer).GetMetadataByMethodSinature(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.MetadataService/GetMetadataByMethodSinature",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetadataServiceServer).GetMetadataByMethodSinature(ctx, req.(*GetMetadataByMethodSinatureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetadataService_GetMetadataByTopic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMetadataByTopicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetadataServiceServer).GetMetadataByTopic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.MetadataService/GetMetadataByTopic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetadataServiceServer).GetMetadataByTopic(ctx, req.(*GetMetadataByTopicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MetadataService_ServiceDesc is the grpc.ServiceDesc for MetadataService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MetadataService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.MetadataService",
	HandlerType: (*MetadataServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UnsubscribeFromMetadata",
			Handler:    _MetadataService_UnsubscribeFromMetadata_Handler,
		},
		{
			MethodName: "GetMetadata",
			Handler:    _MetadataService_GetMetadata_Handler,
		},
		{
			MethodName: "ListMetadata",
			Handler:    _MetadataService_ListMetadata_Handler,
		},
		{
			MethodName: "GetMetadataByMethodSinature",
			Handler:    _MetadataService_GetMetadataByMethodSinature_Handler,
		},
		{
			MethodName: "GetMetadataByTopic",
			Handler:    _MetadataService_GetMetadataByTopic_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeOnMetadata",
			Handler:       _MetadataService_SubscribeOnMetadata_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "github.com/dipdup-net/abi-indexer/pkg/modules/grpc/proto/metadata.proto",
}
