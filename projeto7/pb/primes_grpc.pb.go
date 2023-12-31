// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.1
// source: proto/primes.proto

package pb

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
	PrimeService_SeparatePrimeNumbers_FullMethodName = "/PrimeService/SeparatePrimeNumbers"
)

// PrimeServiceClient is the client API for PrimeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PrimeServiceClient interface {
	SeparatePrimeNumbers(ctx context.Context, in *NumbersRequest, opts ...grpc.CallOption) (*NumbersResponse, error)
}

type primeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPrimeServiceClient(cc grpc.ClientConnInterface) PrimeServiceClient {
	return &primeServiceClient{cc}
}

func (c *primeServiceClient) SeparatePrimeNumbers(ctx context.Context, in *NumbersRequest, opts ...grpc.CallOption) (*NumbersResponse, error) {
	out := new(NumbersResponse)
	err := c.cc.Invoke(ctx, PrimeService_SeparatePrimeNumbers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PrimeServiceServer is the server API for PrimeService service.
// All implementations must embed UnimplementedPrimeServiceServer
// for forward compatibility
type PrimeServiceServer interface {
	SeparatePrimeNumbers(context.Context, *NumbersRequest) (*NumbersResponse, error)
	mustEmbedUnimplementedPrimeServiceServer()
}

// UnimplementedPrimeServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPrimeServiceServer struct {
}

func (UnimplementedPrimeServiceServer) SeparatePrimeNumbers(context.Context, *NumbersRequest) (*NumbersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SeparatePrimeNumbers not implemented")
}
func (UnimplementedPrimeServiceServer) mustEmbedUnimplementedPrimeServiceServer() {}

// UnsafePrimeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PrimeServiceServer will
// result in compilation errors.
type UnsafePrimeServiceServer interface {
	mustEmbedUnimplementedPrimeServiceServer()
}

func RegisterPrimeServiceServer(s grpc.ServiceRegistrar, srv PrimeServiceServer) {
	s.RegisterService(&PrimeService_ServiceDesc, srv)
}

func _PrimeService_SeparatePrimeNumbers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NumbersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PrimeServiceServer).SeparatePrimeNumbers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PrimeService_SeparatePrimeNumbers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PrimeServiceServer).SeparatePrimeNumbers(ctx, req.(*NumbersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PrimeService_ServiceDesc is the grpc.ServiceDesc for PrimeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PrimeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "PrimeService",
	HandlerType: (*PrimeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SeparatePrimeNumbers",
			Handler:    _PrimeService_SeparatePrimeNumbers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/primes.proto",
}
