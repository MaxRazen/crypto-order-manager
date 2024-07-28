package grpcserver

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/MaxRazen/crypto-order-manager/internal/app"
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
	app *app.App
	ctx context.Context // global context to use in delayed gorutines
}

// CreateOrder implements ordermanager.OrderManagerServer
func (s *server) CreateOrder(ctx context.Context, req *ordergrpc.CreateOrderRequest) (*ordergrpc.CreateOrderResponse, error) {
	s.log.Debug(ctx, "received CreateOrder request", "request", req)

	// extracting & validating data from request
	data, err := order.Validate(req)
	if err != nil {
		s.log.Debug(ctx, "validation failed", "error", err)
		return &ordergrpc.CreateOrderResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	s.log.Debug(ctx, "received request is validated", "data", data)

	// creating an order
	ord, err := order.Create(ctx, s.log, s.app.Storage, data)
	if err != nil {
		s.log.Error(ctx, "order cannot be saved", "error", err)
		return &ordergrpc.CreateOrderResponse{
			Success: false,
			Message: "Order cannot be processed",
		}, nil
	}

	// putting order to the queue for placement in background
	s.app.OrderPlacer.Add(ctx, ord)

	return &ordergrpc.CreateOrderResponse{
		Success: true,
		Message: "Order creation request is accepted successfully",
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
	return token == authorizationToken
}

var authorizationToken string

func Run(ctx context.Context, log *logger.Logger, app *app.App, authToken, port string) error {
	authorizationToken = authToken
	lis, err := net.Listen("tcp", ":"+port)

	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(ensureValidToken),
	}

	s := grpc.NewServer(opts...)

	ordergrpc.RegisterOrderManagerServer(s, &server{
		log: log,
		app: app,
		ctx: ctx,
	})

	log.Info(ctx, "server listening at "+lis.Addr().String())

	return s.Serve(lis)
}
