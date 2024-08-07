package tracker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/MaxRazen/crypto-order-manager/internal/logger"
	"github.com/MaxRazen/crypto-order-manager/internal/market"
)

type PlacedOrderRepository interface {
	FetchAllUncompleted(ctx context.Context) ([]market.PlacedOrder, error)
	UpdateStatus(ctx context.Context, o *market.PlacedOrder, status string) error
}

type trackingItem struct {
	order     market.PlacedOrder
	checkedAt time.Time
	done      bool
}

type Tracker struct {
	mu       sync.RWMutex
	log      *logger.Logger
	rep      PlacedOrderRepository
	markets  *market.Collection
	items    []*trackingItem
	ticker   *time.Ticker
	intval   time.Duration
	shutdown chan struct{}
	input    chan market.PlacedOrder
}

func New(log *logger.Logger, rep PlacedOrderRepository, markets *market.Collection, tickInterval, checkInterval time.Duration) *Tracker {
	return &Tracker{
		log:      log,
		rep:      rep,
		markets:  markets,
		items:    make([]*trackingItem, 0),
		ticker:   time.NewTicker(tickInterval),
		intval:   checkInterval,
		shutdown: make(chan struct{}),
		input:    make(chan market.PlacedOrder),
	}
}

func (t *Tracker) Init(ctx context.Context) error {
	orders, err := t.rep.FetchAllUncompleted(ctx)
	if err != nil {
		return err
	}

	if len(orders) > 0 {
		t.log.Info(ctx, fmt.Sprintf("tracker: initializing... %d orders to be placed", len(orders)))
	}

	for _, o := range orders {
		t.put(o)
	}

	return nil
}

func (t *Tracker) GetInputChan() chan<- market.PlacedOrder {
	return t.input
}

func (t *Tracker) Run(ctx context.Context) error {
	for {
		select {
		case order, ok := <-t.input:
			if ok {
				t.log.Debug(ctx, "adding a new order to the queue")
				t.put(order)
			}
		case _, ok := <-t.shutdown:
			if ok {
				t.log.Debug(ctx, "ticker is stopped")
				return nil
			}
		case tcr := <-t.ticker.C:
			t.log.Debug(ctx, "tracker :: tick :: "+tcr.Format(time.TimeOnly))
			t.tick(ctx)
		}
	}
}

func (t *Tracker) Stop(ctx context.Context) {
	t.log.Debug(ctx, "stoppping ticker")
	t.shutdown <- struct{}{}
}

func (t *Tracker) tick(ctx context.Context) {
	for idx, itm := range t.items {
		if itm.done {
			if _, err := t.popByIndex(idx); err != nil {
				t.log.Error(ctx, "tracker: tick cleanup", "error", err)
			}
		}
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	now := time.Now()
	for _, itm := range t.items {
		if itm.checkedAt.Add(t.intval).Before(now) {
			continue
		}

		go t.check(ctx, itm)
	}
}

func (t *Tracker) check(ctx context.Context, ti *trackingItem) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	t.log.Debug(ctx, "tracker: order to be checked: #"+ti.order.ClientOrderId)

	mc, ok := t.markets.Get(ti.order.Market)
	if !ok {
		t.log.Error(ctx, fmt.Sprintf("tracker: market is not supported '%s'", ti.order.Market))
		ti.done = true
		return
	}

	po, err := mc.GetOrder(ctx, ti.order.OrderId)
	if err != nil {
		t.log.Error(ctx, "tracker: order obtain error", "error", err)
		ti.done = true
		return
	}

	if ti.order.Status != po.Status {
		t.rep.UpdateStatus(ctx, &ti.order, po.Status)
		t.log.Debug(ctx, "tracker: order status has been changed", "old", ti.order.Status, "new", po.Status)
	}

	status := mc.TranslatedStatus(po.Status)
	if status == market.StatusCompleted || status == market.StatusCanceled {
		t.log.Debug(ctx, "tracker: order has been canceled or completed", "orderId", ti.order.ClientOrderId)
		ti.done = true
	}

	ti.checkedAt = time.Now()
}

func (t *Tracker) put(o market.PlacedOrder) {
	item := trackingItem{
		order:     o,
		checkedAt: time.Now(),
	}
	t.mu.Lock()
	t.items = append(t.items, &item)
	t.mu.Unlock()
}

func (t *Tracker) popByIndex(idx int) (*trackingItem, error) {
	if idx < 0 || idx >= len(t.items) {
		return nil, fmt.Errorf("index out of range")
	}

	item := t.items[idx]

	t.mu.Lock()
	t.items = append(t.items[:idx], t.items[idx+1:]...)
	t.mu.Unlock()

	return item, nil
}
