package tracker

import (
	"context"
	"testing"
	"time"

	"github.com/MaxRazen/crypto-order-manager/internal/market"
	"github.com/MaxRazen/crypto-order-manager/internal/storage"
	"github.com/MaxRazen/crypto-order-manager/internal/testpkg"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	ctx := context.Background()

	// -------------------------------------------------------------------------
	// Init logger

	log, logBuffer := testpkg.NewLogger()

	// -------------------------------------------------------------------------
	// Init repository

	repo := new(testpkg.PlacedOrderRepositoryMock)

	// -------------------------------------------------------------------------
	// Prepare orders

	inputOrderFirst := market.PlacedOrder{
		DSKey:         storage.NewIDKey(market.PlacedOrderKind, storage.GenerateID()),
		OrderId:       "10001",
		Market:        "mock",
		ClientOrderId: "701",
		Status:        market.StatusNew,
	}

	inputOrderSecond := market.PlacedOrder{
		DSKey:         storage.NewIDKey(market.PlacedOrderKind, storage.GenerateID()),
		OrderId:       "10002",
		Market:        "mock",
		ClientOrderId: "702",
		Status:        market.StatusNew,
	}

	responseOrderFirst := inputOrderFirst
	responseOrderSecond := inputOrderSecond

	// -------------------------------------------------------------------------
	// Init markets collection

	mc := new(testpkg.MarketClientMock)
	mc.On("GetOrder", ctx, inputOrderFirst.OrderId).Return(&responseOrderFirst, nil)
	mc.On("GetOrder", ctx, inputOrderSecond.OrderId).Return(&responseOrderSecond, nil)

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

	time.Sleep(35 * time.Millisecond)

	go func() {
		inputChan <- inputOrderSecond
	}()

	time.Sleep(200 * time.Millisecond)
	tr.Stop(ctx)

	// -------------------------------------------------------------------------
	// Assert results

	err := <-trackerErrors
	assert.Nil(t, err, "tracker finished with error")

	t.Logf("\n----------- LOGS -----------\n%s----------- LOGS END -----------\n", logBuffer.String())
}
