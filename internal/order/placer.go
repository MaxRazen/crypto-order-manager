package order

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
	"github.com/MaxRazen/crypto-order-manager/internal/ordergrpc"
)

const (
	placementTimeout = 5 * time.Second
)

var (
	errNoAvailableAssets = errors.New("no available assets")
)

type PlacementService struct {
	log        *logger.Logger
	placedRep  *market.Repository
	rep        *Repository
	markets    *market.Collection
	input      chan *Order
	shutdown   chan struct{}
	inProgress atomic.Int32
}

func NewPlacementService(log *logger.Logger, rep *Repository, placedRep *market.Repository, markets *market.Collection) *PlacementService {
	return &PlacementService{
		log:        log,
		rep:        rep,
		placedRep:  placedRep,
		markets:    markets,
		input:      make(chan *Order),
		shutdown:   make(chan struct{}),
		inProgress: atomic.Int32{},
	}
}

func (ps *PlacementService) Init(ctx context.Context) *PlacementService {
	orders, err := ps.rep.FetchAllUnplaced(ctx)
	if err != nil {
		ps.log.Error(ctx, "ps: initializing error", "error", err.Error())
		return ps
	}

	if len(orders) > 0 {
		ps.log.Info(ctx, fmt.Sprintf("ps: initializing... %d orders to be placed", len(orders)))
	}

	for _, o := range orders {
		ps.Add(ctx, &o)
	}

	return ps
}

func (ps *PlacementService) Run(ctx context.Context, trackerQueue chan<- market.PlacedOrder) {
	for {
		select {
		case order := <-ps.input:
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
	close(ps.input)
}

func (ps *PlacementService) Add(ctx context.Context, order *Order) {
	ps.log.Debug(ctx, "ps: adding a new order to the queue", "orderId", order.DSKey.ID)

	// check if the program is in shotdown process
	select {
	case _, ok := <-ps.shutdown:
		if ok {
			return
		}
	default:
	}

	ps.inProgress.Add(1)
	ps.input <- order
}

func (ps *PlacementService) placeOrderHandler(ctx context.Context, order *Order, trackerQueue chan<- market.PlacedOrder) {
	ctx, cancel := context.WithTimeout(ctx, placementTimeout)
	defer cancel()

	ps.log.Debug(ctx, "ps: placing the order", "orderId", order.DSKey.ID)
	defer ps.inProgress.Add(-1)

	placedOrder, err := ps.placeOrder(ctx, order)

	if err == nil {
		if err := ps.rep.UpdateStatus(ctx, order, StatusPlaced); err != nil {
			ps.log.Error(ctx, "ps: updating status error", "status", StatusPlaced, "orderId", order.DSKey.ID)
		}
		trackerQueue <- *placedOrder

		return
	}

	if err == errNoAvailableAssets {
		ps.log.Info(ctx, err.Error(), "order", order)

		if err := ps.rep.UpdateStatus(ctx, order, StatusSkipped); err != nil {
			ps.log.Error(ctx, "ps: updating status error", "status", StatusSkipped, "orderId", order.DSKey.ID)
		}

		return
	}

	if err := ps.rep.UpdateStatus(ctx, order, StatusFailed); err != nil {
		ps.log.Error(ctx, "ps: updating status error", "status", StatusFailed, "orderId", order.DSKey.ID)
	}
}

func (ps *PlacementService) placeOrder(ctx context.Context, order *Order) (*market.PlacedOrder, error) {
	// verify market client is available
	mc, ok := ps.markets.Get(order.Market)
	if !ok {
		return nil, fmt.Errorf("market is not supported '%s'", order.Market)
	}

	pair := mc.FormatPair(order.Pair)
	symbols := strings.Split(order.Pair, "-")
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
	if order.Action == market.ActionBuy {
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

	expectedPrice, _ := strconv.ParseFloat(order.Price, 64)

	// retrieve exchange info
	exchangeInfo, err := mc.GetExchangeInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("exchange info cannot be retrieved")
	}
	symbolExchangeInfo := exchangeInfo.Find(pair)
	if symbolExchangeInfo == nil {
		return nil, fmt.Errorf("symbol exchange info is not provided pair: %s", pair)
	}

	quantity := calculateDesiredQuantity(order.Quantity, balance, expectedPrice)
	quantity, err = calculateValidQuantity(quantity, symbolExchangeInfo)
	if err != nil {
		return nil, err
	}

	ps.log.Debug(ctx, "quantity is calculated", "qty", quantity)

	internalId := strconv.FormatInt(order.DSKey.ID, 10)

	ocd := market.OrderCreationData{
		ClientOrderId: internalId,
		Pair:          pair,
		Side:          order.Action,
		Type:          order.Behavior,
		Price:         expectedPrice,
		Quantity:      quantity,
	}

	ps.log.Debug(ctx, "order to be placed", "market.OrderCreationData", ocd)

	placedOrder, err := mc.PlaceOrder(ctx, ocd)
	if err != nil {
		return nil, fmt.Errorf("order cannot be placed")
	}

	err = ps.placedRep.Create(ctx, placedOrder)
	if err != nil {
		return nil, fmt.Errorf("placed order cannot be saved to storage")
	}

	ps.log.Info(ctx, "order is placed successfully", "order", placedOrder, "key", placedOrder.DSKey)

	return placedOrder, nil
}

func calculateDesiredQuantity(demandedQty Quantity, balance, price float64) float64 {
	val, _ := strconv.ParseFloat(demandedQty.Value, 64)

	// calculate quantity for fixed type. The value is taken without any calculation
	if demandedQty.Type == ordergrpc.QuantityType_FIXED.String() {
		return val
	}

	// calculate quantity for percent type
	return (balance / price) * val / 100
}

func calculateValidQuantity(desiredQty float64, exInfo *market.SymbolExchangeInfo) (float64, error) {
	if desiredQty < exInfo.MinQty {
		return 0, errNoEnoughAssets
	} else if desiredQty > exInfo.MaxQty {
		return 0, errQuantityLimit
	}

	q := desiredQty / exInfo.StepSize
	q = float64(int64(q)) * exInfo.StepSize

	return q, nil
}
