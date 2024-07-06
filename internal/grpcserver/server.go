package grpcserver

import (
	"context"
	"net"
	"strings"

	"github.com/MaxRazen/crypto-order-manager/internal/logger"
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
	log *logger.Logger
	ctx context.Context
}

// CreateOrder implements ordermanager.OrderManagerServer
func (s *server) CreateOrder(ctx context.Context, req *ordergrpc.CreateOrderRequest) (*ordergrpc.CreateOrderResponse, error) {
	// Log the received request
	s.log.Debug(ctx, "received CreateOrder request", req)

	data, err := order.Validate(req)
	if err != nil {
		return &ordergrpc.CreateOrderResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	s.log.Debug(ctx, "received request is validated", data)

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

func Run(ctx context.Context, log *logger.Logger, port string) error {
	lis, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatal(ctx, "failed to listen", "error", err)
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(ensureValidToken),
	}

	s := grpc.NewServer(opts...)

	ordergrpc.RegisterOrderManagerServer(s, &server{
		log: log,
		ctx: ctx,
	})

	log.Info(ctx, "server listening at "+lis.Addr().String())

	return s.Serve(lis)
}
