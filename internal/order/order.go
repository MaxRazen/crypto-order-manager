package order

import "cloud.google.com/go/datastore"

const (
	IncomingOrdersKind = "incoming_orders"
	PlacedOrderKind    = "placed_orders"
)

type CreationData struct {
	DSKey     *datastore.Key `datastore:"__key__"`
	ID        int64          `datastore:"id"`
	Pair      string         `datastore:"pair"`
	Market    string         `datastore:"market"`
	Action    string         `datastore:"action,noindex"`
	Behavior  string         `datastore:"behavior,noindex"`
	Price     string         `datastore:"price,noindex"`
	Quantity  Quantity       `datastore:"quantity,noindex"`
	Deadlines []Deadline     `datastore:"deadlines,noindex"`
}

type Quantity struct {
	Type  string `datastore:"type,noindex"`
	Value string `datastore:"value,noindex"`
}

type Deadline struct {
	Type   string `datastore:"type,noindex"`
	Value  string `datastore:"value,noindex"`
	Action string `datastore:"action,noindex"`
}
