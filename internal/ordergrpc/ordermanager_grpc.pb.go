// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.27.1
// source: protofiles/ordermanager.proto

package ordergrpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	OrderManager_CreateOrder_FullMethodName = "/ordermanager.OrderManager/CreateOrder"
)

// OrderManagerClient is the client API for OrderManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrderManagerClient interface {
	CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error)
}

type orderManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderManagerClient(cc grpc.ClientConnInterface) OrderManagerClient {
	return &orderManagerClient{cc}
}

func (c *orderManagerClient) CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateOrderResponse)
	err := c.cc.Invoke(ctx, OrderManager_CreateOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderManagerServer is the server API for OrderManager service.
// All implementations must embed UnimplementedOrderManagerServer
// for forward compatibility
type OrderManagerServer interface {
	CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error)
	mustEmbedUnimplementedOrderManagerServer()
}

// UnimplementedOrderManagerServer must be embedded to have forward compatible implementations.
type UnimplementedOrderManagerServer struct {
}

func (UnimplementedOrderManagerServer) CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrder not implemented")
}
func (UnimplementedOrderManagerServer) mustEmbedUnimplementedOrderManagerServer() {}

// UnsafeOrderManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrderManagerServer will
// result in compilation errors.
type UnsafeOrderManagerServer interface {
	mustEmbedUnimplementedOrderManagerServer()
}

func RegisterOrderManagerServer(s grpc.ServiceRegistrar, srv OrderManagerServer) {
	s.RegisterService(&OrderManager_ServiceDesc, srv)
}

func _OrderManager_CreateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderManagerServer).CreateOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderManager_CreateOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderManagerServer).CreateOrder(ctx, req.(*CreateOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OrderManager_ServiceDesc is the grpc.ServiceDesc for OrderManager service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OrderManager_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ordermanager.OrderManager",
	HandlerType: (*OrderManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateOrder",
			Handler:    _OrderManager_CreateOrder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protofiles/ordermanager.proto",
}
