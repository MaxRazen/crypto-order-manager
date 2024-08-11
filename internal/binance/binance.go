package binance

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"

	"github.com/MaxRazen/crypto-order-manager/internal/market"
	"github.com/MaxRazen/crypto-order-manager/internal/utils"

	connector "github.com/binance/binance-connector-go"
)

func New(apiKey, secretKey, baseUrl string, dryRun bool) *BinanceClient {
	bc := BinanceClient{
		client: connector.NewClient(apiKey, secretKey, baseUrl),
		dryRun: dryRun,
	}
	return &bc
}

type BinanceClient struct {
	client *connector.Client
	dryRun bool
}

func (bc *BinanceClient) Name() string {
	return "binance"
}

func (bc *BinanceClient) FormatPair(pair string) string {
	return strings.ToUpper(strings.ReplaceAll(pair, "-", ""))
}

func (bc *BinanceClient) AvailableAssets(ctx context.Context, pairs []string) ([]market.AssetBalance, error) {
	resp, err := bc.client.NewGetAccountService().Do(ctx)
	if err != nil {
		return nil, err
	}

	assets := make([]market.AssetBalance, 0)
	for _, v := range resp.Balances {
		if !utils.InSlice(pairs, v.Asset) {
			continue
		}
		assets = append(assets, market.AssetBalance{
			Asset:   v.Asset,
			Balance: v.Free,
		})
	}
	return assets, nil
}

func (bc *BinanceClient) GetPairPrice(ctx context.Context, pair string) (string, error) {
	resp, err := bc.client.NewTickerPriceService().Symbol(pair).Do(ctx)
	if err != nil {
		return "", err
	}

	return resp.Price, nil
}

func (bc *BinanceClient) PlaceOrder(ctx context.Context, orderData market.OrderCreationData) (*market.PlacedOrder, error) {
	// whenever dryRun=true a test order will be placed
	if bc.dryRun {
		_, err := bc.client.NewTestNewOrder().
			Symbol(orderData.Pair).
			Side(orderData.Side).
			OrderType(orderData.Type).
			Price(orderData.Price).
			Quantity(orderData.Quantity).
			NewClientOrderId(orderData.ClientOrderId).
			TimeInForce("GTC").
			Do(ctx)
		if err != nil {
			return nil, err
		}
		placedOrderData := market.PlacedOrder{
			Market:        bc.Name(),
			Symbol:        orderData.Pair,
			OrderId:       fmt.Sprintf("100%d%d", rand.IntN(1000), time.Now().UnixMilli()),
			ClientOrderId: orderData.ClientOrderId,
			Status:        "ACTIVE",
			PlacedAt:      time.Now(),
		}
		return &placedOrderData, nil
	}

	if true { // TODO: temp prevention mechanism
		return nil, errors.New("woops! the real order is going to be placed. terminating")
	}

	// place a real order
	respRaw, err := bc.client.NewCreateOrderService().
		Symbol(orderData.Pair).
		Side(orderData.Side).
		Type(orderData.Type).
		Price(orderData.Price).
		Quantity(orderData.Quantity).
		NewClientOrderId(orderData.ClientOrderId).
		TimeInForce("GTC").
		Do(ctx)
	if err != nil {
		return nil, err
	}
	resp, ok := respRaw.(connector.CreateOrderResponseFULL)
	if !ok {
		return nil, errors.New("response type cannot be converted to specified type")
	}

	placedOrderData := market.PlacedOrder{
		Market:        bc.Name(),
		Symbol:        resp.Symbol,
		OrderId:       strconv.FormatInt(resp.OrderId, 10),
		ClientOrderId: resp.ClientOrderId,
		Status:        resp.Status,
		PlacedAt:      time.Now(),
	}
	return &placedOrderData, nil
}

func (bc *BinanceClient) CancelOrder(ctx context.Context, orderId string) error {
	ordId, err := strconv.ParseInt(orderId, 10, 64)
	if err != nil {
		return err
	}
	_, err = bc.client.NewCancelOrderService().OrderId(ordId).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (bc *BinanceClient) GetOrder(ctx context.Context, pair, orderId string) (*market.PlacedOrder, error) {
	ordId, err := strconv.ParseInt(orderId, 10, 64)
	if err != nil {
		return nil, err
	}
	resp, err := bc.client.NewGetOrderService().
		Symbol("").
		OrderId(ordId).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	po := market.PlacedOrder{
		Market:        bc.Name(),
		Symbol:        resp.Symbol,
		OrderId:       strconv.FormatInt(resp.OrderId, 10),
		ClientOrderId: resp.ClientOrderId,
		Status:        resp.Status,
		PlacedAt:      time.UnixMilli(int64(resp.Time)),
	}
	return &po, nil
}

func (bc *BinanceClient) GetExchangeInfo(ctx context.Context) (*market.ExchangeInfo, error) {
	info, err := bc.client.NewExchangeInfoService().Do(ctx)
	if err != nil {
		return nil, err
	}

	symbolsInfo := make([]market.SymbolExchangeInfo, 0, len(info.Symbols))
	for _, s := range info.Symbols {
		sInfo := market.SymbolExchangeInfo{Symbol: s.Symbol}

		for _, f := range s.Filters {
			if f.FilterType == "LOT_SIZE" {
				sInfo.MaxQty = parseFloat(f.MaxQty)
				sInfo.MinQty = parseFloat(f.MinQty)
				sInfo.StepSize = parseFloat(f.StepSize)
			}
		}

		symbolsInfo = append(symbolsInfo, sInfo)
	}

	exchangeInfo := market.ExchangeInfo{
		Market:  bc.Name(),
		Symbols: symbolsInfo,
	}

	return &exchangeInfo, nil
}

func (bc *BinanceClient) TranslatedStatus(status string) string {
	// NEW
	// PARTIALLY_FILLED
	// FILLED
	// CANCELED
	// PENDING_CANCEL
	// REJECTED
	// EXPIRED
	// EXPIRED_IN_MATCH
	if status == "NEW" || status == "PARTIALLY_FILLED" {
		return market.StatusNew
	} else if status == "FILLED" {
		return market.StatusCompleted
	}
	return market.StatusCanceled
}

func parseFloat(val string) float64 {
	v, _ := strconv.ParseFloat(val, 64)
	return v
}
