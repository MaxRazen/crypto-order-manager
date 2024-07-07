package market

import "context"

const (
	ActionBuy  = "BUY"
	ActionSell = "SELL"
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
	OrderId       string `datastore:"orderId"`
	ClientOrderId string `datastore:"clientOrderId"`
	Status        string `datastore:"status"`
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
	GetOrder(ctx context.Context, orderId string) error
}

type ExchangeInfoFetcher interface {
	GetExchangeInfo(ctx context.Context) (*ExchangeInfo, error)
}
