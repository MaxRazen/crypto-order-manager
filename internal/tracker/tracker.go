package tracker

import (
	"context"
	"sync"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/MaxRazen/crypto-order-manager/internal/logger"
	"github.com/MaxRazen/crypto-order-manager/internal/market"
)

type Tracker struct {
	mu     sync.Mutex
	ds     *datastore.Client
	log    *logger.Logger
	orders []market.PlacedOrder
	ticker *time.Ticker
}

func New(log *logger.Logger, ds *datastore.Client) *Tracker {
	return &Tracker{
		ds:     ds,
		log:    log,
		orders: make([]market.PlacedOrder, 0),
		ticker: time.NewTicker(time.Second * 30),
	}
}

func (t *Tracker) Put(o market.PlacedOrder) {
	t.mu.Lock()
	t.orders = append(t.orders, o)
	t.mu.Unlock()
}

func (t *Tracker) Start(ctx context.Context) error {
	for tcr := range t.ticker.C {
		t.log.Debug(ctx, "tracker :: tick :: "+tcr.Format(time.TimeOnly))
		t.tick()
	}

	return nil
}

func (t *Tracker) Stop(ctx context.Context) {
	t.log.Debug(ctx, "stoppping ticker")
	t.ticker.Stop()
}

func (t *Tracker) tick() {
	// TODO
}
