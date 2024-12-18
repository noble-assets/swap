// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: swap/v1/query.proto

package swapv1

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
	Query_Paused_FullMethodName       = "/swap.v1.Query/Paused"
	Query_Pools_FullMethodName        = "/swap.v1.Query/Pools"
	Query_Pool_FullMethodName         = "/swap.v1.Query/Pool"
	Query_SimulateSwap_FullMethodName = "/swap.v1.Query/SimulateSwap"
	Query_Rates_FullMethodName        = "/swap.v1.Query/Rates"
	Query_Rate_FullMethodName         = "/swap.v1.Query/Rate"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	// Retrieves a list of the currently paused Pools.
	Paused(ctx context.Context, in *QueryPaused, opts ...grpc.CallOption) (*QueryPausedResponse, error)
	// Retrieves the details of all Pools.
	Pools(ctx context.Context, in *QueryPools, opts ...grpc.CallOption) (*QueryPoolsResponse, error)
	// Retrieves details of a specific Pool.
	Pool(ctx context.Context, in *QueryPool, opts ...grpc.CallOption) (*QueryPoolResponse, error)
	// Simulates a token swap simulation.
	SimulateSwap(ctx context.Context, in *QuerySimulateSwap, opts ...grpc.CallOption) (*MsgSwapResponse, error)
	// Retrieves exchange rates for all tokens, with the optionality of filtering by algorithm.
	Rates(ctx context.Context, in *QueryRates, opts ...grpc.CallOption) (*QueryRatesResponse, error)
	// Retrieves exchange rates for a specific token, with the optionality of filtering by algorithm.
	Rate(ctx context.Context, in *QueryRate, opts ...grpc.CallOption) (*QueryRateResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Paused(ctx context.Context, in *QueryPaused, opts ...grpc.CallOption) (*QueryPausedResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryPausedResponse)
	err := c.cc.Invoke(ctx, Query_Paused_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Pools(ctx context.Context, in *QueryPools, opts ...grpc.CallOption) (*QueryPoolsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryPoolsResponse)
	err := c.cc.Invoke(ctx, Query_Pools_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Pool(ctx context.Context, in *QueryPool, opts ...grpc.CallOption) (*QueryPoolResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryPoolResponse)
	err := c.cc.Invoke(ctx, Query_Pool_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) SimulateSwap(ctx context.Context, in *QuerySimulateSwap, opts ...grpc.CallOption) (*MsgSwapResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MsgSwapResponse)
	err := c.cc.Invoke(ctx, Query_SimulateSwap_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Rates(ctx context.Context, in *QueryRates, opts ...grpc.CallOption) (*QueryRatesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryRatesResponse)
	err := c.cc.Invoke(ctx, Query_Rates_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Rate(ctx context.Context, in *QueryRate, opts ...grpc.CallOption) (*QueryRateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryRateResponse)
	err := c.cc.Invoke(ctx, Query_Rate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility.
type QueryServer interface {
	// Retrieves a list of the currently paused Pools.
	Paused(context.Context, *QueryPaused) (*QueryPausedResponse, error)
	// Retrieves the details of all Pools.
	Pools(context.Context, *QueryPools) (*QueryPoolsResponse, error)
	// Retrieves details of a specific Pool.
	Pool(context.Context, *QueryPool) (*QueryPoolResponse, error)
	// Simulates a token swap simulation.
	SimulateSwap(context.Context, *QuerySimulateSwap) (*MsgSwapResponse, error)
	// Retrieves exchange rates for all tokens, with the optionality of filtering by algorithm.
	Rates(context.Context, *QueryRates) (*QueryRatesResponse, error)
	// Retrieves exchange rates for a specific token, with the optionality of filtering by algorithm.
	Rate(context.Context, *QueryRate) (*QueryRateResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedQueryServer struct{}

func (UnimplementedQueryServer) Paused(context.Context, *QueryPaused) (*QueryPausedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Paused not implemented")
}
func (UnimplementedQueryServer) Pools(context.Context, *QueryPools) (*QueryPoolsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Pools not implemented")
}
func (UnimplementedQueryServer) Pool(context.Context, *QueryPool) (*QueryPoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Pool not implemented")
}
func (UnimplementedQueryServer) SimulateSwap(context.Context, *QuerySimulateSwap) (*MsgSwapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SimulateSwap not implemented")
}
func (UnimplementedQueryServer) Rates(context.Context, *QueryRates) (*QueryRatesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Rates not implemented")
}
func (UnimplementedQueryServer) Rate(context.Context, *QueryRate) (*QueryRateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Rate not implemented")
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

func _Query_Paused_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryPaused)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Paused(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Paused_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Paused(ctx, req.(*QueryPaused))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Pools_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryPools)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Pools(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Pools_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Pools(ctx, req.(*QueryPools))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Pool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryPool)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Pool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Pool_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Pool(ctx, req.(*QueryPool))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_SimulateSwap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySimulateSwap)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).SimulateSwap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_SimulateSwap_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).SimulateSwap(ctx, req.(*QuerySimulateSwap))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Rates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRates)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Rates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Rates_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Rates(ctx, req.(*QueryRates))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Rate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Rate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Rate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Rate(ctx, req.(*QueryRate))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "swap.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Paused",
			Handler:    _Query_Paused_Handler,
		},
		{
			MethodName: "Pools",
			Handler:    _Query_Pools_Handler,
		},
		{
			MethodName: "Pool",
			Handler:    _Query_Pool_Handler,
		},
		{
			MethodName: "SimulateSwap",
			Handler:    _Query_SimulateSwap_Handler,
		},
		{
			MethodName: "Rates",
			Handler:    _Query_Rates_Handler,
		},
		{
			MethodName: "Rate",
			Handler:    _Query_Rate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "swap/v1/query.proto",
}
