package placer

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/MaxRazen/crypto-order-manager/internal/logger"
	"github.com/MaxRazen/crypto-order-manager/internal/market"
	"github.com/MaxRazen/crypto-order-manager/internal/order"
	"github.com/MaxRazen/crypto-order-manager/internal/ordergrpc"
)

const (
	placementTimeout = 5 * time.Second
)

var (
	errNoAvailableAssets = errors.New("no available assets")
	errQuantityLimit     = errors.New("provided quantity value is too big or small")
)

type OrderRepository interface {
	FetchAllUnplaced(ctx context.Context) ([]order.Order, error)
	UpdateStatus(ctx context.Context, o *order.Order, status string) error
}

type PlacedOrderRepository interface {
	Create(ctx context.Context, po *market.PlacedOrder) error
}

type PlacementService struct {
	log         *logger.Logger
	orderRep    OrderRepository
	marketRep   PlacedOrderRepository
	markets     *market.Collection
	input       chan *order.Order
	shutdown    chan struct{}
	terminating atomic.Bool
	inProgress  atomic.Int32
}

func NewPlacementService(
	log *logger.Logger,
	orderRep OrderRepository,
	marketRep PlacedOrderRepository,
	markets *market.Collection,
) *PlacementService {
	return &PlacementService{
		log:         log,
		orderRep:    orderRep,
		marketRep:   marketRep,
		markets:     markets,
		input:       make(chan *order.Order),
		shutdown:    make(chan struct{}),
		terminating: atomic.Bool{},
		inProgress:  atomic.Int32{},
	}
}

func (ps *PlacementService) Init(ctx context.Context) error {
	orders, err := ps.orderRep.FetchAllUnplaced(ctx)
	if err != nil {
		return err
	}

	if len(orders) > 0 {
		ps.log.Info(ctx, fmt.Sprintf("ps: initializing... %d orders to be placed", len(orders)))
	}

	for _, o := range orders {
		ps.Add(ctx, &o)
	}

	return nil
}

func (ps *PlacementService) Run(ctx context.Context, trackerQueue chan<- market.PlacedOrder) {
	for {
		select {
		case order := <-ps.input:
			ps.inProgress.Add(1)
			go ps.placeOrderHandler(ctx, order, trackerQueue)
		case _, ok := <-ps.shutdown:
			if ok {
				ps.log.Debug(ctx, "ps: graceful shutdown signal received", "inProgress", ps.inProgress.Load())
				if ps.inProgress.Load() == 0 {
					close(trackerQueue)
				}
				return
			}
		}
	}
}

func (ps *PlacementService) Stop() {
	ps.shutdown <- struct{}{}
	ps.terminating.Store(true)
	close(ps.input)
}

func (ps *PlacementService) Add(ctx context.Context, order *order.Order) {
	ps.log.Debug(ctx, "ps: adding a new order to the queue", "orderId", order.DSKey.ID)

	if ps.terminating.Load() {
		return
	}

	ps.input <- order
}

func (ps *PlacementService) placeOrderHandler(ctx context.Context, o *order.Order, trackerQueue chan<- market.PlacedOrder) {
	ctx, cancel := context.WithTimeout(ctx, placementTimeout)
	defer cancel()

	ps.log.Debug(ctx, "ps: placing the order", "orderId", o.DSKey.ID)
	defer ps.inProgress.Add(-1)

	placedOrder, err := ps.placeOrder(ctx, o)

	if err == nil {
		if err := ps.orderRep.UpdateStatus(ctx, o, order.StatusPlaced); err != nil {
			ps.log.Error(ctx, "ps: updating status error", "status", order.StatusPlaced, "orderId", o.DSKey.ID)
		}
		trackerQueue <- *placedOrder

		return
	}

	if err == errNoAvailableAssets {
		ps.log.Info(ctx, err.Error(), "order", o)

		if err := ps.orderRep.UpdateStatus(ctx, o, order.StatusSkipped); err != nil {
			ps.log.Error(ctx, "ps: updating status error", "status", order.StatusSkipped, "orderId", o.DSKey.ID)
		}

		return
	} else {
		ps.log.Error(ctx, "ps: placing order error", "orderId", o.DSKey.ID, "error", err)
	}

	if err := ps.orderRep.UpdateStatus(ctx, o, order.StatusFailed); err != nil {
		ps.log.Error(ctx, "ps: updating status error", "status", order.StatusFailed, "orderId", o.DSKey.ID)
	}
}

func (ps *PlacementService) placeOrder(ctx context.Context, o *order.Order) (*market.PlacedOrder, error) {
	// verify market client is available
	mc, ok := ps.markets.Get(o.Market)
	if !ok {
		return nil, fmt.Errorf("market is not supported '%s'", o.Market)
	}

	pair := mc.FormatPair(o.Pair)
	symbols := strings.Split(o.Pair, "-")
	assetsBySymbol := make(map[string]float64)

	ps.log.Debug(ctx, "assets to be fetched", "symbols", symbols)

	assets, err := mc.AvailableAssets(ctx, symbols)
	if err != nil {
		return nil, fmt.Errorf("assets cannot be fetched: %s", err.Error())
	}
	for _, asset := range assets {
		b, _ := strconv.ParseFloat(asset.Balance, 64)
		assetsBySymbol[asset.Asset] = b
	}

	operationSymbol := symbols[0]
	if o.Action == market.ActionBuy {
		operationSymbol = symbols[1]
	}

	balance, ok := assetsBySymbol[operationSymbol]
	if !ok {
		return nil, fmt.Errorf("stake symbol is not received from account")
	}

	if balance == 0.0 {
		return nil, errNoAvailableAssets
	}

	ps.log.Debug(ctx, "account assets are fetched and valid")

	expectedPrice, _ := strconv.ParseFloat(o.Price, 64)

	// retrieve exchange info
	exchangeInfo, err := mc.GetExchangeInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("exchange info cannot be retrieved")
	}
	symbolExchangeInfo := exchangeInfo.Find(pair)
	if symbolExchangeInfo == nil {
		return nil, fmt.Errorf("symbol exchange info is not provided pair: %s", pair)
	}

	quantity := calculateDesiredQuantity(o.Quantity, balance, expectedPrice)
	quantity, err = calculateValidQuantity(quantity, symbolExchangeInfo)
	if err != nil {
		return nil, err
	}

	ps.log.Debug(ctx, "quantity is calculated", "qty", quantity)

	internalId := strconv.FormatInt(o.DSKey.ID, 10)

	ocd := market.OrderCreationData{
		ClientOrderId: internalId,
		Pair:          pair,
		Side:          o.Action,
		Type:          o.Behavior,
		Price:         expectedPrice,
		Quantity:      quantity,
	}

	ps.log.Debug(ctx, "order to be placed", "market.OrderCreationData", ocd)

	placedOrder, err := mc.PlaceOrder(ctx, ocd)
	if err != nil {
		return nil, fmt.Errorf("order cannot be placed")
	}

	ps.log.Info(ctx, "order is placed", "order", placedOrder)

	err = ps.marketRep.Create(ctx, placedOrder)
	if err != nil {
		return nil, fmt.Errorf("placed order cannot be saved to storage")
	}

	// -------------------------------------------------------------------------
	// Copy deadlines

	deadlines := make([]market.Deadline, 0, len(o.Deadlines))
	for _, dl := range o.Deadlines {
		deadlines = append(deadlines, market.Deadline{
			Type:   dl.Type,
			Value:  dl.Value,
			Action: dl.Action,
		})
	}
	placedOrder.Deadlines = deadlines

	ps.log.Info(ctx, "order is placed successfully", "orderId", placedOrder.ClientOrderId, "key", placedOrder.DSKey.ID)

	return placedOrder, nil
}

func calculateDesiredQuantity(demandedQty order.Quantity, balance, price float64) float64 {
	val, _ := strconv.ParseFloat(demandedQty.Value, 64)

	// calculate quantity for fixed type. The value is taken without any calculation
	if demandedQty.Type == ordergrpc.QuantityType_FIXED.String() {
		return val
	}

	// calculate quantity for percent type
	return (balance / price) * val / 100
}

func calculateValidQuantity(desiredQty float64, exInfo *market.SymbolExchangeInfo) (float64, error) {
	if desiredQty < exInfo.MinQty || desiredQty > exInfo.MaxQty {
		return 0, errQuantityLimit
	}

	q := desiredQty / exInfo.StepSize
	q = float64(int64(q)) * exInfo.StepSize

	return q, nil
}
