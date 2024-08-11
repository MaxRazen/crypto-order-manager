package market

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/MaxRazen/crypto-order-manager/internal/deadline"
)

const (
	PlacedOrderKind = "placed_orders"
	ActionBuy       = "BUY"
	ActionSell      = "SELL"
	StatusNew       = "new"
	StatusCompleted = "completed"
	StatusCanceled  = "canceled"
)

type AssetBalance struct {
	Asset   string
	Balance string
}

type OrderCreationData struct {
	ClientOrderId string
	Pair          string
	Side          string
	Type          string
	Price         float64
	Quantity      float64
}

type PlacedOrder struct {
	DSKey         *datastore.Key      `datastore:"__key__"`
	Market        string              `datastore:"market"`
	Symbol        string              `datastore:"symbol"`
	OrderId       string              `datastore:"orderId"`
	ClientOrderId string              `datastore:"clientOrderId"`
	Status        string              `datastore:"status"`
	Deadlines     []deadline.Deadline `datastore:"deadlines,noindex"`
	PlacedAt      time.Time           `datastore:"created_at,noindex"`
}

type ExchangeInfo struct {
	Market  string
	Symbols []SymbolExchangeInfo
}

func (e *ExchangeInfo) Find(symbol string) *SymbolExchangeInfo {
	for _, s := range e.Symbols {
		if s.Symbol == symbol {
			return &s
		}
	}
	return nil
}

type SymbolExchangeInfo struct {
	Symbol   string
	MinQty   float64
	MaxQty   float64
	StepSize float64
}

type MarketClient interface {
	Name() string
	PriceFormatter
	BalanceFetcher
	PairPriceFetcher
	OrderPlacer
	OrderCanceler
	OrderFetcher
	ExchangeInfoFetcher
	StatusTranslator
}

type PriceFormatter interface {
	FormatPair(pair string) string
}

type BalanceFetcher interface {
	AvailableAssets(ctx context.Context, pairs []string) ([]AssetBalance, error)
}

type PairPriceFetcher interface {
	GetPairPrice(ctx context.Context, pair string) (string, error)
}

type OrderPlacer interface {
	PlaceOrder(ctx context.Context, orderData OrderCreationData) (*PlacedOrder, error)
}

type OrderCanceler interface {
	CancelOrder(ctx context.Context, orderId string) error
}

type OrderFetcher interface {
	GetOrder(ctx context.Context, pair, orderId string) (*PlacedOrder, error)
}

type ExchangeInfoFetcher interface {
	GetExchangeInfo(ctx context.Context) (*ExchangeInfo, error)
}

type StatusTranslator interface {
	TranslatedStatus(status string) string
}

type Collection struct {
	markets map[string]MarketClient
}

func NewCollection() *Collection {
	return &Collection{
		markets: make(map[string]MarketClient),
	}
}

func (c *Collection) Add(mc MarketClient) {
	c.markets[mc.Name()] = mc
}

func (c *Collection) Get(name string) (MarketClient, bool) {
	mc, ok := c.markets[name]
	return mc, ok
}
