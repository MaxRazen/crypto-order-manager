package tracker

import (
	"context"
	"testing"
	"time"

	"github.com/MaxRazen/crypto-order-manager/internal/market"
	"github.com/MaxRazen/crypto-order-manager/internal/storage"
	"github.com/MaxRazen/crypto-order-manager/internal/testpkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var inputOrderFirst = market.PlacedOrder{
	DSKey:         storage.NewIDKey(market.PlacedOrderKind, storage.GenerateID()),
	OrderId:       "10001",
	Market:        "mock",
	ClientOrderId: "701",
	Status:        "NEW",
}

var inputOrderSecond = market.PlacedOrder{
	DSKey:         storage.NewIDKey(market.PlacedOrderKind, storage.GenerateID()),
	OrderId:       "10002",
	Market:        "mock",
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
	repo.On("FetchAllUncompleted", mock.Anything).Once().Return(dsResponse, nil)

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
	repo.On("UpdateStatus", mock.Anything, &inputOrderFirst, "FILLED").Return(nil)
	repo.On("UpdateStatus", mock.Anything, &inputOrderSecond, "CANCELED").Return(nil)

	// -------------------------------------------------------------------------
	// Init markets collection

	responseOrderFirst := inputOrderFirst
	responseOrderSecond := inputOrderSecond

	mc := new(testpkg.MarketClientMock)
	mc.On("GetOrder", mock.Anything, inputOrderFirst.OrderId).Return(&responseOrderFirst, nil)
	mc.On("GetOrder", mock.Anything, inputOrderSecond.OrderId).Return(&responseOrderSecond, nil)

	markets := market.NewCollection()
	markets.Add(mc)

	// -------------------------------------------------------------------------
	// Tracker setup

	tr := New(log, repo, markets, 20*time.Millisecond, 50*time.Millisecond)

	// -------------------------------------------------------------------------
	// Run tracker & do some work

	trackerErrors := make(chan error, 1)
	inputChan := tr.GetInputChan()
	go func() {
		trackerErrors <- tr.Run(ctx)
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

	err := <-trackerErrors
	if err != nil {
		t.Errorf("tracker finished with error: %v", err)
	}

	mc.AssertExpectations(t)
	repo.AssertExpectations(t)

	testpkg.VerboseOutput(t, logBuffer.String())
}
