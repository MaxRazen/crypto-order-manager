package placer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/MaxRazen/crypto-order-manager/internal/market"
	"github.com/MaxRazen/crypto-order-manager/internal/order"
	"github.com/MaxRazen/crypto-order-manager/internal/storage"
	"github.com/MaxRazen/crypto-order-manager/internal/testpkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var inputOrders = []order.Order{
	{
		DSKey:       storage.NewIDKey(order.OrdersKind, 701),
		Status:      order.StatusNew,
		RuleId:      "a001",
		Pair:        "BTC-USDT",
		Market:      "mock",
		Action:      market.ActionBuy,
		Behavior:    "LIMIT",
		Price:       "20500.0",
		Quantity:    order.Quantity{Type: "FIXED", Value: "0.1"},
		Deadlines:   make([]order.Deadline, 0),
		CreatedAt:   time.Now(),
		CompletedAt: time.Now(),
	},
	{
		DSKey:       storage.NewIDKey(order.OrdersKind, 702),
		Status:      order.StatusNew,
		RuleId:      "a001",
		Pair:        "SOL-USDT",
		Market:      "mock",
		Action:      market.ActionBuy,
		Behavior:    "MARKET",
		Price:       "120.0",
		Quantity:    order.Quantity{Type: "FIXED", Value: "2"},
		Deadlines:   make([]order.Deadline, 0),
		CreatedAt:   time.Now(),
		CompletedAt: time.Now(),
	},
}

func TestInit(t *testing.T) {
	ctx := context.Background()

	// -------------------------------------------------------------------------
	// Init logger

	log, logBuffer := testpkg.NewLogger()

	// -------------------------------------------------------------------------
	// Init repository

	orderRepo := new(testpkg.OrderRepositoryMock)
	orderRepo.On("FetchAllUnplaced", mock.Anything).Once().Return(inputOrders, nil)

	placedorderRepo := new(testpkg.PlacedOrderRepositoryMock)

	// -------------------------------------------------------------------------
	// Init markets collection

	mc := new(testpkg.MarketClientMock)
	markets := market.NewCollection()
	markets.Add(mc)

	// -------------------------------------------------------------------------
	// Placer setup & run

	orderPlacer := NewPlacementService(log, orderRepo, placedorderRepo, markets)

	var receivedOrders []*order.Order
	go func() {
		for o := range orderPlacer.input {
			receivedOrders = append(receivedOrders, o)
		}
	}()

	err := orderPlacer.Init(ctx)

	// -------------------------------------------------------------------------
	// Assert expectations

	if !assert.Nil(t, err, "placer.Init returned an error") {
		t.Logf("placer.Init error: %v", err)
	}
	assert.Equal(t, int32(0), orderPlacer.inProgress.Load())

	testpkg.VerboseOutput(t, logBuffer.String())
}

func TestInit_DatastoreError(t *testing.T) {
	ctx := context.Background()

	// -------------------------------------------------------------------------
	// Init logger

	log, logBuffer := testpkg.NewLogger()

	// -------------------------------------------------------------------------
	// Init repository

	orderRepo := new(testpkg.OrderRepositoryMock)
	orderRepo.On("FetchAllUnplaced", mock.Anything).Once().Return([]order.Order{}, fmt.Errorf("authorization error"))

	placedorderRepo := new(testpkg.PlacedOrderRepositoryMock)

	// -------------------------------------------------------------------------
	// Init markets collection

	mc := new(testpkg.MarketClientMock)
	markets := market.NewCollection()
	markets.Add(mc)

	// -------------------------------------------------------------------------
	// Placer setup & run

	orderPlacer := NewPlacementService(log, orderRepo, placedorderRepo, markets)

	err := orderPlacer.Init(ctx)

	// -------------------------------------------------------------------------
	// Assert expectations

	assert.NotNil(t, err, "placer.Init returned an error")

	testpkg.VerboseOutput(t, logBuffer.String())
}

func TestRun(t *testing.T) {

}

func TestCalculateDesiredQuantity(t *testing.T) {
	testCases := []struct {
		demandedQty order.Quantity
		balance     float64
		price       float64
		expeted     float64
	}{
		{
			demandedQty: order.Quantity{Type: "FIXED", Value: "5.5"},
			balance:     1000.0,
			price:       120.0,
			expeted:     5.5,
		},
		{
			demandedQty: order.Quantity{Type: "PERCENT", Value: "50"},
			balance:     768.0,
			price:       128.0,
			expeted:     3.0,
		},
	}

	for i, tc := range testCases {
		assert.Equal(
			t,
			tc.expeted,
			calculateDesiredQuantity(tc.demandedQty, tc.balance, tc.price),
			fmt.Sprintf("(%d) value assertion failed", i),
		)
	}
}

func TestCalculateValidQuantity(t *testing.T) {
	testCases := []struct {
		desiredQty    float64
		exchInfo      market.SymbolExchangeInfo
		expetedResult float64
		expectedError error
	}{
		{
			desiredQty: 0.001,
			exchInfo: market.SymbolExchangeInfo{
				Symbol:   "BTC",
				MinQty:   0.001,
				MaxQty:   100000.0,
				StepSize: 0.001,
			},
			expetedResult: 0.001,
			expectedError: nil,
		},
		{
			desiredQty: 0.0025,
			exchInfo: market.SymbolExchangeInfo{
				Symbol:   "BTC",
				MinQty:   0.001,
				MaxQty:   100000.0,
				StepSize: 0.001,
			},
			expetedResult: 0.002,
			expectedError: nil,
		},
		{
			desiredQty: 0.00001,
			exchInfo: market.SymbolExchangeInfo{
				Symbol:   "BTC",
				MinQty:   0.001,
				MaxQty:   100000.0,
				StepSize: 0.001,
			},
			expetedResult: 0.001,
			expectedError: errQuantityLimit,
		},
		{
			desiredQty: 100000.1,
			exchInfo: market.SymbolExchangeInfo{
				Symbol:   "BTC",
				MinQty:   0.001,
				MaxQty:   100000.0,
				StepSize: 0.001,
			},
			expetedResult: 100000.0,
			expectedError: errQuantityLimit,
		},
	}

	for i, tc := range testCases {
		res, err := calculateValidQuantity(tc.desiredQty, &tc.exchInfo)
		if err == nil {
			assert.Equal(t, tc.expetedResult, res, fmt.Sprintf("(%d) value assertion failed", i))
		} else {
			assert.Equal(t, tc.expectedError, err, fmt.Sprintf("(%d) error assertion failed", i))
		}
	}
}
