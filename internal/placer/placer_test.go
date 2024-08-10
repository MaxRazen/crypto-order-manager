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
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var contextArgument = mock.MatchedBy(func(_ context.Context) bool {
	return true
})

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
		RuleId:      "a002",
		Pair:        "SOL-USDT",
		Market:      "mock",
		Action:      market.ActionSell,
		Behavior:    "MARKET",
		Price:       "120.0",
		Quantity:    order.Quantity{Type: "FIXED", Value: "2"},
		Deadlines:   make([]order.Deadline, 0),
		CreatedAt:   time.Now(),
		CompletedAt: time.Now(),
	},
}

var placedOrders = []market.PlacedOrder{
	{
		DSKey:         storage.NewIDKey(market.PlacedOrderKind, 801),
		OrderId:       "20004",
		Market:        "mock",
		ClientOrderId: "101",
		Status:        "NEW",
		Deadlines:     make([]market.Deadline, 0),
	},
	{
		DSKey:         storage.NewIDKey(market.PlacedOrderKind, 802),
		OrderId:       "20005",
		Market:        "mock",
		ClientOrderId: "102",
		Status:        "NEW",
		Deadlines:     make([]market.Deadline, 0),
	},
}

var ordersCreationData = []market.OrderCreationData{
	{
		ClientOrderId: "701",
		Pair:          "BTC-USDT",
		Side:          market.ActionBuy,
		Type:          "LIMIT",
		Price:         20500.0,
		Quantity:      0.1,
	},
	{
		ClientOrderId: "702",
		Pair:          "SOL-USDT",
		Side:          market.ActionSell,
		Type:          "MARKET",
		Price:         120.0,
		Quantity:      2,
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
	orderRepo.On("FetchAllUnplaced", contextArgument).Once().Return(inputOrders, nil)

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
	orderRepo.On("FetchAllUnplaced", contextArgument).Once().Return([]order.Order{}, fmt.Errorf("authorization error"))

	placedOrderRepo := new(testpkg.PlacedOrderRepositoryMock)

	// -------------------------------------------------------------------------
	// Init markets collection

	mc := new(testpkg.MarketClientMock)
	markets := market.NewCollection()
	markets.Add(mc)

	// -------------------------------------------------------------------------
	// Placer setup & run

	orderPlacer := NewPlacementService(log, orderRepo, placedOrderRepo, markets)

	err := orderPlacer.Init(ctx)

	// -------------------------------------------------------------------------
	// Assert expectations

	assert.NotNil(t, err, "placer.Init returned an error")

	testpkg.VerboseOutput(t, logBuffer.String())
}

func TestRun(t *testing.T) {
	ctx := context.Background()

	// -------------------------------------------------------------------------
	// Init logger

	log, logBuffer := testpkg.NewLogger()

	// -------------------------------------------------------------------------
	// Init repository

	orderRepo := new(testpkg.OrderRepositoryMock)

	firstOrder := inputOrders[0]
	secondOrder := inputOrders[1]
	placedOrderFirst := placedOrders[0]
	placedOrderSecond := placedOrders[1]

	orderRepo.On("UpdateStatus", contextArgument, &firstOrder, order.StatusPlaced).Return(nil)
	orderRepo.On("UpdateStatus", contextArgument, &secondOrder, order.StatusPlaced).Return(nil)

	placedOrderRepo := new(testpkg.PlacedOrderRepositoryMock)
	placedOrderRepo.On("Create", contextArgument, &placedOrderFirst).Return(nil)
	placedOrderRepo.On("Create", contextArgument, &placedOrderSecond).Return(nil)

	// -------------------------------------------------------------------------
	// Init markets collection

	mc := new(testpkg.MarketClientMock)

	balances := []market.AssetBalance{
		{Asset: "USDT", Balance: "2050"},
		{Asset: "SOL", Balance: "500"},
	}

	mc.On("AvailableAssets", contextArgument, []string{"BTC", "USDT"}).Return(balances, nil)
	mc.On("AvailableAssets", contextArgument, []string{"SOL", "USDT"}).Return(balances, nil)

	symbolExchangeIno := []market.SymbolExchangeInfo{
		{Symbol: "SOL-USDT", MinQty: 0.001, MaxQty: 100000.0, StepSize: 0.001},
		{Symbol: "BTC-USDT", MinQty: 0.0001, MaxQty: 1000.0, StepSize: 0.0001},
	}
	exchangeInfo := market.ExchangeInfo{Market: "mock", Symbols: symbolExchangeIno}
	mc.On("GetExchangeInfo", contextArgument).Return(&exchangeInfo, nil).Twice()

	mc.On("PlaceOrder", contextArgument, mock.MatchedBy(func(o market.OrderCreationData) bool {
		return cmp.Equal(o, ordersCreationData[0])
	})).Return(&placedOrderFirst, nil).Once()

	mc.On("PlaceOrder", contextArgument, mock.MatchedBy(func(o market.OrderCreationData) bool {
		return cmp.Equal(o, ordersCreationData[1])
	})).Return(&placedOrderSecond, nil).Once()

	markets := market.NewCollection()
	markets.Add(mc)

	// -------------------------------------------------------------------------
	// Placer setup & run

	orderPlacer := NewPlacementService(log, orderRepo, placedOrderRepo, markets)

	trackerInputChan := make(chan market.PlacedOrder)

	go func() {
		orderPlacer.Run(ctx, trackerInputChan)
	}()

	placedOrders := make([]market.PlacedOrder, 0)
	go func() {
		for po := range trackerInputChan {
			placedOrders = append(placedOrders, po)
		}
	}()

	time.AfterFunc(10*time.Millisecond, func() {
		orderPlacer.Add(ctx, &firstOrder)
	})

	time.AfterFunc(15*time.Millisecond, func() {
		orderPlacer.Add(ctx, &secondOrder)
	})

	time.Sleep(50 * time.Millisecond)
	orderPlacer.Stop()

	thirdOrder := inputOrders[0]
	thirdOrder.DSKey.ID = 101
	thirdOrder.Pair = "GOLDUSDT"
	orderPlacer.Add(ctx, &thirdOrder)

	time.Sleep(10 * time.Millisecond)

	// -------------------------------------------------------------------------
	// Assert expectations

	assert.Equal(t, 2, len(placedOrders))
	assert.Equal(t, int32(0), orderPlacer.inProgress.Load())
	// TODO: check expectations by placed orders

	orderRepo.AssertExpectations(t)
	placedOrderRepo.AssertExpectations(t)
	mc.AssertExpectations(t)

	testpkg.VerboseOutput(t, logBuffer.String())
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
