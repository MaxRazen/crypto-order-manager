package tracker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/MaxRazen/crypto-order-manager/internal/logger"
	"github.com/MaxRazen/crypto-order-manager/internal/market"
)

type Tracker struct {
	mu       sync.Mutex
	log      *logger.Logger
	rep      *market.Repository
	orders   []market.PlacedOrder
	ticker   *time.Ticker
	shutdown chan struct{}
	input    chan market.PlacedOrder
}

func New(log *logger.Logger, rep *market.Repository) *Tracker {
	return &Tracker{
		log:      log,
		rep:      rep,
		orders:   make([]market.PlacedOrder, 0),
		ticker:   time.NewTicker(30 * time.Second),
		shutdown: make(chan struct{}),
		input:    make(chan market.PlacedOrder),
	}
}

func (t *Tracker) Init(ctx context.Context) *Tracker {
	orders, err := t.rep.FetchAllUncompleted(ctx)
	if err != nil {
		t.log.Error(ctx, "tracker: initializing error", "error", err.Error())
		return t
	}

	if len(orders) > 0 {
		t.log.Info(ctx, fmt.Sprintf("ps: initializing... %d orders to be placed", len(orders)))
	}

	for _, o := range orders {
		t.put(o)
	}

	return t
}

func (t *Tracker) GetInputChan() chan<- market.PlacedOrder {
	return t.input
}

func (t *Tracker) Run(ctx context.Context) error {
	for {
		select {
		case order, ok := <-t.input:
			if ok {
				t.put(order)
			}
		case _, ok := <-t.shutdown:
			if ok {
				t.log.Debug(ctx, "ticker is stopped")
				return nil
			}
		case tcr := <-t.ticker.C:
			t.log.Debug(ctx, "tracker :: tick :: "+tcr.Format(time.TimeOnly))
			t.tick()
		}
	}
}

func (t *Tracker) Stop(ctx context.Context) {
	t.log.Debug(ctx, "stoppping ticker")
	t.shutdown <- struct{}{}
}

func (t *Tracker) tick() {
	// TODO
}

func (t *Tracker) put(o market.PlacedOrder) {
	t.mu.Lock()
	t.orders = append(t.orders, o)
	t.mu.Unlock()
}
