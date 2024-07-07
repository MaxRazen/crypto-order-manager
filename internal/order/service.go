package order

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/MaxRazen/crypto-order-manager/internal/app"
	"github.com/MaxRazen/crypto-order-manager/internal/logger"
	"github.com/MaxRazen/crypto-order-manager/internal/market"
	"github.com/MaxRazen/crypto-order-manager/internal/ordergrpc"
	"github.com/MaxRazen/crypto-order-manager/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errQuantityLimit  = errors.New("provided quantity value is too big")
	errNoEnoughAssets = errors.New("provided quantity value is too small")
)

func Validate(req *ordergrpc.CreateOrderRequest) (*CreationData, error) {
	// validate Pair property
	pair := req.GetPair()
	if match, err := regexp.MatchString("^([A-Z]+)-([A-Z]+)$", pair); err != nil || !match {
		return nil, status.Error(codes.InvalidArgument, "pair does not match the expected pattern")
	}

	// validate Market property
	market := req.GetMarket()
	if market == "" {
		return nil, status.Error(codes.InvalidArgument, "market must be specified")
	}

	// validate Action propery
	action := req.GetAction().String()
	if req.GetAction() == ordergrpc.ActionType_UNKNOWN_ACTION {
		return nil, status.Error(codes.InvalidArgument, "action must be specified")
	}

	// validate Behavior propery
	behavior := req.GetBehavior().String()
	if req.GetBehavior() == ordergrpc.Behavior_UNKNOWN_BEHAVIOR {
		return nil, status.Error(codes.InvalidArgument, "behavior must be specified")
	}

	// validate Price property
	price := req.GetPrice()
	if match, err := regexp.MatchString("^(\\d+)(\\.\\d+)?$", price); err != nil || !match {
		return nil, status.Error(codes.InvalidArgument, "price does not match the expected pattern")
	}

	// validate Quantity properties
	quantity := req.GetQuantity()
	quantityType := quantity.GetType().String()
	quantityValue := quantity.GetValue()
	if req.GetQuantity().GetType() == ordergrpc.QuantityType_UNKNOWN_QUANTITY_TYPE {
		return nil, status.Error(codes.InvalidArgument, "quantity.type must be specified")
	}
	if match, err := regexp.MatchString("^(\\d+)(\\.\\d+)?$", price); err != nil || !match {
		return nil, status.Error(codes.InvalidArgument, "quantity.value does not match the expected pattern")
	}

	// validate Deadline properties
	reqDeadlines := req.GetDeadlines()
	deadlines := make([]Deadline, 0, len(reqDeadlines))
	for idx, d := range reqDeadlines {
		// validate Deadline.type property
		dType := d.Type.String()
		if d.Type == ordergrpc.DeadlineType_UNKNOWN_DEADLINE_TYPE {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("deadlines[%v].type must be specified", idx))
		}

		// validate Deadline.value property
		dValue := d.Value
		if dValue == "" {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("deadlines[%v].value must be specified", idx))
		}

		// validate Deadline.action property
		dAction := d.Action.String()
		if d.Action == ordergrpc.DeadlineAction_UNKNOWN_DEADLINE_ACTION {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("deadlines[%v].action must be specified", idx))
		}

		deadlines = append(deadlines, Deadline{
			Type:   dType,
			Value:  dValue,
			Action: dAction,
		})
	}

	data := CreationData{
		Pair:     pair,
		Market:   market,
		Action:   action,
		Behavior: behavior,
		Price:    price,
		Quantity: Quantity{
			Type:  quantityType,
			Value: quantityValue,
		},
		Deadlines: deadlines,
	}

	return &data, nil
}

func SaveOrder(ctx context.Context, log *logger.Logger, app *app.App, data *CreationData) error {
	dsk := storage.NewIDKey(IncomingOrdersKind, storage.GenerateID())

	dsk, err := app.Storage.Put(ctx, dsk, data)
	if err != nil {
		return err
	}
	data.DSKey = dsk

	return nil
}

func PlaceOrder(ctx context.Context, log *logger.Logger, app *app.App, data *CreationData) {
	// verify market cliend is available
	mc, ok := app.Markets[data.Market]
	if !ok {
		log.Error(ctx, "market is not supported", "market", data.Market)
		return
	}

	pair := mc.FormatPair(data.Pair)
	symbols := strings.Split(data.Pair, "-")
	assetsBySymbol := make(map[string]float64)

	log.Debug(ctx, "assets to be fetched", "symbols", symbols)

	assets, err := mc.AvailableAssets(ctx, symbols)
	if err != nil {
		log.Error(ctx, "assets cannot be fetched", "error", err.Error())
		return
	}
	for _, asset := range assets {
		b, _ := strconv.ParseFloat(asset.Balance, 64)
		assetsBySymbol[asset.Asset] = b
	}

	operationSymbol := symbols[0]
	if data.Action == market.ActionBuy {
		operationSymbol = symbols[1]
	}

	balance, ok := assetsBySymbol[operationSymbol]
	if !ok {
		log.Error(ctx, "stake symbol is not received from account")
		return
	}

	if balance == 0.0 {
		log.Info(ctx, "no available assets")
		return
	}

	log.Debug(ctx, "account assets are fetched and valid")

	expectedPrice, _ := strconv.ParseFloat(data.Price, 64)

	// retrieve exchange info
	exchangeInfo, err := mc.GetExchangeInfo(ctx)
	if err != nil {
		log.Error(ctx, "exchange info cannot be retrieved", "error", err)
		return
	}
	symbolExchangeInfo := exchangeInfo.Find(pair)
	if symbolExchangeInfo == nil {
		log.Error(ctx, "symbol exchange info is not provided", "pair", pair)
		return
	}

	quantity := calculateDesiredQuantity(data.Quantity, balance, expectedPrice)
	quantity, err = calculateValidQuantity(quantity, symbolExchangeInfo)
	if err != nil {
		log.Error(ctx, err.Error())
		return
	}

	log.Debug(ctx, "quantity is calculated", "qty", quantity)

	internalId := strconv.FormatInt(data.DSKey.ID, 10)

	ocd := market.OrderCreationData{
		ClientOrderId: internalId,
		Pair:          pair,
		Side:          data.Action,
		Type:          data.Behavior,
		Price:         expectedPrice,
		Quantity:      quantity,
	}

	log.Debug(ctx, "order to be placed", "market.OrderCreationData", ocd)

	placedOrder, err := mc.PlaceOrder(ctx, ocd)
	if err != nil {
		log.Error(ctx, "order cannot be placed", "error", err)
		return
	}

	sKey := storage.NewIDKey(PlacedOrderKind, storage.GenerateID())
	_, err = app.Storage.Put(ctx, sKey, placedOrder)
	if err != nil {
		log.Error(ctx, "placed order cannot be saved to storage", "order", placedOrder)
		return
	}

	log.Info(ctx, "order is placed successfully", "order", placedOrder, "key", sKey)
}

func calculateDesiredQuantity(demandedQty Quantity, balance, price float64) float64 {
	val, _ := strconv.ParseFloat(demandedQty.Value, 64)

	// calculate quantity for fixed type. The value is taken without any calculation
	if demandedQty.Type == ordergrpc.QuantityType_FIXED.String() {
		return val
	}

	// calculate quantity for percent type
	// quantity := (balance / price) * val / 100 * 10e5
	// return float64(int64(quantity)) / 10e5
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
