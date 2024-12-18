// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: swap/stableswap/v1/query.proto

package stableswapv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Query_PositionsByProvider_FullMethodName          = "/swap.stableswap.v1.Query/PositionsByProvider"
	Query_BondedPositionsByProvider_FullMethodName    = "/swap.stableswap.v1.Query/BondedPositionsByProvider"
	Query_UnbondingPositionsByProvider_FullMethodName = "/swap.stableswap.v1.Query/UnbondingPositionsByProvider"
	Query_RewardsByProvider_FullMethodName            = "/swap.stableswap.v1.Query/RewardsByProvider"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	// Retrieves all the positions by a specific provider, including bonded/unbonded positions and rewards.
	PositionsByProvider(ctx context.Context, in *QueryPositionsByProvider, opts ...grpc.CallOption) (*QueryPositionsByProviderResponse, error)
	// Retrieves all the bonded positions by a specific provider.
	BondedPositionsByProvider(ctx context.Context, in *QueryBondedPositionsByProvider, opts ...grpc.CallOption) (*QueryBondedPositionsByProviderResponse, error)
	// Retrieves all the unbonding positions by a specific provider.
	UnbondingPositionsByProvider(ctx context.Context, in *QueryUnbondingPositionsByProvider, opts ...grpc.CallOption) (*QueryUnbondingPositionsByProviderResponse, error)
	// Retrieves all the rewards by a specific provider.
	RewardsByProvider(ctx context.Context, in *QueryRewardsByProvider, opts ...grpc.CallOption) (*QueryRewardsByProviderResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) PositionsByProvider(ctx context.Context, in *QueryPositionsByProvider, opts ...grpc.CallOption) (*QueryPositionsByProviderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryPositionsByProviderResponse)
	err := c.cc.Invoke(ctx, Query_PositionsByProvider_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) BondedPositionsByProvider(ctx context.Context, in *QueryBondedPositionsByProvider, opts ...grpc.CallOption) (*QueryBondedPositionsByProviderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryBondedPositionsByProviderResponse)
	err := c.cc.Invoke(ctx, Query_BondedPositionsByProvider_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) UnbondingPositionsByProvider(ctx context.Context, in *QueryUnbondingPositionsByProvider, opts ...grpc.CallOption) (*QueryUnbondingPositionsByProviderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryUnbondingPositionsByProviderResponse)
	err := c.cc.Invoke(ctx, Query_UnbondingPositionsByProvider_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) RewardsByProvider(ctx context.Context, in *QueryRewardsByProvider, opts ...grpc.CallOption) (*QueryRewardsByProviderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryRewardsByProviderResponse)
	err := c.cc.Invoke(ctx, Query_RewardsByProvider_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility.
type QueryServer interface {
	// Retrieves all the positions by a specific provider, including bonded/unbonded positions and rewards.
	PositionsByProvider(context.Context, *QueryPositionsByProvider) (*QueryPositionsByProviderResponse, error)
	// Retrieves all the bonded positions by a specific provider.
	BondedPositionsByProvider(context.Context, *QueryBondedPositionsByProvider) (*QueryBondedPositionsByProviderResponse, error)
	// Retrieves all the unbonding positions by a specific provider.
	UnbondingPositionsByProvider(context.Context, *QueryUnbondingPositionsByProvider) (*QueryUnbondingPositionsByProviderResponse, error)
	// Retrieves all the rewards by a specific provider.
	RewardsByProvider(context.Context, *QueryRewardsByProvider) (*QueryRewardsByProviderResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedQueryServer struct{}

func (UnimplementedQueryServer) PositionsByProvider(context.Context, *QueryPositionsByProvider) (*QueryPositionsByProviderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PositionsByProvider not implemented")
}
func (UnimplementedQueryServer) BondedPositionsByProvider(context.Context, *QueryBondedPositionsByProvider) (*QueryBondedPositionsByProviderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BondedPositionsByProvider not implemented")
}
func (UnimplementedQueryServer) UnbondingPositionsByProvider(context.Context, *QueryUnbondingPositionsByProvider) (*QueryUnbondingPositionsByProviderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnbondingPositionsByProvider not implemented")
}
func (UnimplementedQueryServer) RewardsByProvider(context.Context, *QueryRewardsByProvider) (*QueryRewardsByProviderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RewardsByProvider not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}
func (UnimplementedQueryServer) testEmbeddedByValue()               {}

// UnsafeQueryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueryServer will
// result in compilation errors.
type UnsafeQueryServer interface {
	mustEmbedUnimplementedQueryServer()
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	// If the following call pancis, it indicates UnimplementedQueryServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Query_ServiceDesc, srv)
}

func _Query_PositionsByProvider_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryPositionsByProvider)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).PositionsByProvider(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_PositionsByProvider_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).PositionsByProvider(ctx, req.(*QueryPositionsByProvider))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_BondedPositionsByProvider_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryBondedPositionsByProvider)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).BondedPositionsByProvider(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_BondedPositionsByProvider_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).BondedPositionsByProvider(ctx, req.(*QueryBondedPositionsByProvider))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_UnbondingPositionsByProvider_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryUnbondingPositionsByProvider)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).UnbondingPositionsByProvider(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_UnbondingPositionsByProvider_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).UnbondingPositionsByProvider(ctx, req.(*QueryUnbondingPositionsByProvider))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_RewardsByProvider_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRewardsByProvider)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).RewardsByProvider(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_RewardsByProvider_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).RewardsByProvider(ctx, req.(*QueryRewardsByProvider))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "swap.stableswap.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PositionsByProvider",
			Handler:    _Query_PositionsByProvider_Handler,
		},
		{
			MethodName: "BondedPositionsByProvider",
			Handler:    _Query_BondedPositionsByProvider_Handler,
		},
		{
			MethodName: "UnbondingPositionsByProvider",
			Handler:    _Query_UnbondingPositionsByProvider_Handler,
		},
		{
			MethodName: "RewardsByProvider",
			Handler:    _Query_RewardsByProvider_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "swap/stableswap/v1/query.proto",
}
