package market

import (
	"context"

	"cloud.google.com/go/datastore"
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
	DSKey         *datastore.Key `datastore:"__key__"`
	OrderId       string         `datastore:"orderId"`
	Market        string         `datastore:"market"`
	ClientOrderId string         `datastore:"clientOrderId"`
	Status        string         `datastore:"status"`
	Deadlines     []Deadline     `datastore:"deadlines,noindex"`
}

type Deadline struct {
	Type   string `datastore:"type,noindex"`
	Value  string `datastore:"value,noindex"`
	Action string `datastore:"action,noindex"`
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
	GetOrder(ctx context.Context, orderId string) (*PlacedOrder, error)
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
