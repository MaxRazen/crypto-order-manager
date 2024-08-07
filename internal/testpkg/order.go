package testpkg

import (
	"context"

	"github.com/MaxRazen/crypto-order-manager/internal/order"
	"github.com/stretchr/testify/mock"
)

type OrderRepositoryMock struct {
	mock.Mock
}

func (r *OrderRepositoryMock) FetchAllUnplaced(ctx context.Context) ([]order.Order, error) {
	args := r.Called(ctx)
	return args.Get(0).([]order.Order), args.Error(1)
}

func (r *OrderRepositoryMock) UpdateStatus(ctx context.Context, o *order.Order, status string) error {
	args := r.Called(ctx, o, status)
	return args.Error(0)
}
