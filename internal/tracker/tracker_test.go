package tracker

import (
	"context"
	"testing"
	"time"

	"github.com/MaxRazen/crypto-order-manager/internal/deadline"
	"github.com/MaxRazen/crypto-order-manager/internal/market"
	"github.com/MaxRazen/crypto-order-manager/internal/storage"
	"github.com/MaxRazen/crypto-order-manager/internal/testpkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var contextArgument = mock.MatchedBy(func(_ context.Context) bool {
	return true
})

var inputOrderFirst = market.PlacedOrder{
	DSKey:         storage.NewIDKey(market.PlacedOrderKind, storage.GenerateID()),
	OrderId:       "10001",
	Market:        "mock",
	Symbol:        "BTCUSDT",
	ClientOrderId: "701",
	Status:        "NEW",
}

var inputOrderSecond = market.PlacedOrder{
	DSKey:         storage.NewIDKey(market.PlacedOrderKind, storage.GenerateID()),
	OrderId:       "10002",
	Market:        "mock",
	Symbol:        "BTCUSDT",
	ClientOrderId: "702",
	Status:        "NEW",
}

func TestInit(t *testing.T) {
	ctx := context.Background()

	// -------------------------------------------------------------------------
	// Init logger

	log, logBuffer := testpkg.NewLogger()

	// -------------------------------------------------------------------------
	// Init repository

	dsResponse := []market.PlacedOrder{inputOrderFirst, inputOrderSecond}

	repo := new(testpkg.PlacedOrderRepositoryMock)
	repo.On("FetchAllUncompleted", contextArgument).Once().Return(dsResponse, nil)

	// -------------------------------------------------------------------------
	// Init markets collection

	mc := new(testpkg.MarketClientMock)
	markets := market.NewCollection()
	markets.Add(mc)

	// -------------------------------------------------------------------------
	// Tracker setup

	tr := New(log, repo, markets, 20*time.Millisecond, 50*time.Millisecond)

	err := tr.Init(ctx)

	// -------------------------------------------------------------------------
	// Assert expectations

	if !assert.Nil(t, err, "tracker.Init returned an error") {
		t.Logf("tracker.Init error: %v", err)
	}

	assert.Equal(t, 2, len(tr.items))
	repo.AssertExpectations(t)
	mc.AssertExpectations(t)

	testpkg.VerboseOutput(t, logBuffer.String())
}

func TestRun(t *testing.T) {
	ctx := context.Background()

	// -------------------------------------------------------------------------
	// Init logger

	log, logBuffer := testpkg.NewLogger()

	// -------------------------------------------------------------------------
	// Init repository

	repo := new(testpkg.PlacedOrderRepositoryMock)
	repo.On("UpdateStatus", contextArgument, &inputOrderFirst, "FILLED").Return(nil)
	repo.On("UpdateStatus", contextArgument, &inputOrderSecond, "CANCELED").Return(nil)

	// -------------------------------------------------------------------------
	// Init markets collection

	responseOrderFirst := inputOrderFirst
	responseOrderSecond := inputOrderSecond

	mc := new(testpkg.MarketClientMock)
	mc.On("GetOrder", contextArgument, inputOrderFirst.Symbol, inputOrderFirst.OrderId).Return(&responseOrderFirst, nil)
	mc.On("GetOrder", contextArgument, inputOrderSecond.Symbol, inputOrderSecond.OrderId).Return(&responseOrderSecond, nil)

	markets := market.NewCollection()
	markets.Add(mc)

	// -------------------------------------------------------------------------
	// Tracker setup

	tr := New(log, repo, markets, 20*time.Millisecond, 50*time.Millisecond)

	// -------------------------------------------------------------------------
	// Run tracker & do some work

	inputChan := tr.GetInputChan()
	go func() {
		tr.Run(ctx)
	}()

	go func() {
		inputChan <- inputOrderFirst
	}()

	// -------------------------------------------------------------------------
	// Simulate asyncronius inputs/changes
	time.AfterFunc(35*time.Microsecond, func() {
		inputChan <- inputOrderSecond
	})
	time.AfterFunc(45*time.Millisecond, func() {
		responseOrderFirst.Status = "FILLED"
	})
	time.AfterFunc(160*time.Millisecond, func() {
		responseOrderSecond.Status = "CANCELED"
	})

	time.Sleep(200 * time.Millisecond)
	tr.Stop(ctx)

	// -------------------------------------------------------------------------
	// Assert expectations

	mc.AssertExpectations(t)
	repo.AssertExpectations(t)

	testpkg.VerboseOutput(t, logBuffer.String())
}

func TestRun_WithDeadline(t *testing.T) {
	ctx := context.Background()

	// -------------------------------------------------------------------------
	// Init logger

	log, logBuffer := testpkg.NewLogger()

	// -------------------------------------------------------------------------
	// Init repository

	firstOrder := inputOrderFirst
	firstOrder.PlacedAt = time.Now().Add(-40 * time.Millisecond)
	firstOrder.Deadlines = []deadline.Deadline{
		{
			Type:   deadline.TypeTime,
			Action: deadline.ActionCancel,
			Value:  "100",
		},
	}

	repo := new(testpkg.PlacedOrderRepositoryMock)
	repo.On("UpdateStatus", contextArgument, &firstOrder, "CANCELED").Return(nil).Once()

	// -------------------------------------------------------------------------
	// Init markets collection

	responseOrder := inputOrderFirst

	mc := new(testpkg.MarketClientMock)
	mc.On("GetOrder", contextArgument, firstOrder.Symbol, firstOrder.OrderId).Return(&responseOrder, nil)
	mc.On("CancelOrder", contextArgument, firstOrder.OrderId).Return(nil).Once().Run(func(args mock.Arguments) {
		responseOrder.Status = "CANCELED"
	})

	markets := market.NewCollection()
	markets.Add(mc)

	// -------------------------------------------------------------------------
	// Tracker setup

	tr := New(log, repo, markets, 20*time.Millisecond, 50*time.Millisecond)

	// -------------------------------------------------------------------------
	// Run tracker & do some work

	inputChan := tr.GetInputChan()
	go func() {
		tr.Run(ctx)
	}()

	// -------------------------------------------------------------------------
	// Simulate asyncronius inputs/changes
	time.AfterFunc(1*time.Microsecond, func() {
		inputChan <- firstOrder
	})

	time.Sleep(200 * time.Millisecond)
	tr.Stop(ctx)

	// -------------------------------------------------------------------------
	// Assert expectations

	mc.AssertExpectations(t)
	repo.AssertExpectations(t)

	testpkg.VerboseOutput(t, logBuffer.String())
}
