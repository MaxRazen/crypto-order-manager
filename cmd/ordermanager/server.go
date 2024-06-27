package main

import (
	"context"
	"log"
	"net"
	"strings"

	"github.com/MaxRazen/crypto-order-manager/internal/order"
	"github.com/MaxRazen/crypto-order-manager/internal/ordergrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// server is used to implement ordermanager.OrderManagerServer
type server struct {
	ordergrpc.UnimplementedOrderManagerServer
}

// CreateOrder implements ordermanager.OrderManagerServer
func (s *server) CreateOrder(ctx context.Context, req *ordergrpc.CreateOrderRequest) (*ordergrpc.CreateOrderResponse, error) {
	// Log the received request
	log.Printf("Received CreateOrder request: %v", req)

	data, err := order.Validate(req)
	if err != nil {
		return &ordergrpc.CreateOrderResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	log.Printf("received request is validated: %v", data)

	// Here you would add your business logic to handle the order creation
	// For this example, we just return a success response
	return &ordergrpc.CreateOrderResponse{
		Success: true,
		Message: "Order created successfully",
	}, nil
}

func ensureValidToken(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	if !valid(md["authorization"]) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	return handler(ctx, req)
}

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	return token == "some-secret-token"
}

func runServer(port string) {
	lis, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(ensureValidToken),
	}

	s := grpc.NewServer(opts...)

	ordergrpc.RegisterOrderManagerServer(s, &server{})

	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
