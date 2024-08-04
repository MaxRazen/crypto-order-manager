package testpkg

import (
	"context"

	"github.com/MaxRazen/crypto-order-manager/internal/market"
	"github.com/stretchr/testify/mock"
)

// -------------------------------------------------------------------------
// Market Client mock

type MarketClientMock struct {
	mock.Mock
}

func (mc *MarketClientMock) Name() string {
	return "mock"
}

func (mc *MarketClientMock) FormatPair(pair string) string {
	return pair
}

func (mc *MarketClientMock) AvailableAssets(ctx context.Context, pairs []string) ([]market.AssetBalance, error) {
	args := mc.Called(ctx, pairs)
	return args.Get(0).([]market.AssetBalance), args.Error(1)
}

func (mc *MarketClientMock) GetPairPrice(ctx context.Context, pair string) (string, error) {
	args := mc.Called(ctx, pair)
	return args.String(0), args.Error(1)
}

func (mc *MarketClientMock) PlaceOrder(ctx context.Context, orderData market.OrderCreationData) (*market.PlacedOrder, error) {
	args := mc.Called(ctx, orderData)
	return args.Get(0).(*market.PlacedOrder), args.Error(1)
}

func (mc *MarketClientMock) CancelOrder(ctx context.Context, orderId string) error {
	args := mc.Called(ctx, orderId)
	return args.Error(0)
}

func (mc *MarketClientMock) GetOrder(ctx context.Context, orderId string) (*market.PlacedOrder, error) {
	args := mc.Called(ctx, orderId)
	return args.Get(0).(*market.PlacedOrder), args.Error(1)
}

func (mc *MarketClientMock) GetExchangeInfo(ctx context.Context) (*market.ExchangeInfo, error) {
	args := mc.Called(ctx)
	return args.Get(0).(*market.ExchangeInfo), args.Error(1)
}

func (mc *MarketClientMock) TranslatedStatus(status string) string {
	if status == "NEW" || status == "PARTIALLY_FILLED" {
		return market.StatusNew
	} else if status == "FILLED" {
		return market.StatusCompleted
	}
	return market.StatusCanceled
}

// -------------------------------------------------------------------------
// Market PlacedOrderRepository mock

type PlacedOrderRepositoryMock struct {
	mock.Mock
}

func (r *PlacedOrderRepositoryMock) Create(ctx context.Context, o *market.PlacedOrder) error {
	args := r.Called(ctx, o)
	return args.Error(0)
}

func (r *PlacedOrderRepositoryMock) FetchAllUncompleted(ctx context.Context) ([]market.PlacedOrder, error) {
	args := r.Called(ctx)
	return args.Get(0).([]market.PlacedOrder), args.Error(1)
}

func (r *PlacedOrderRepositoryMock) UpdateStatus(ctx context.Context, o *market.PlacedOrder, status string) error {
	args := r.Called(ctx, o, status)
	o.Status = status

	return args.Error(0)
}
