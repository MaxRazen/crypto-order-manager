package order

import (
	"time"

	"cloud.google.com/go/datastore"
)

const (
	OrdersKind    = "orders"
	StatusNew     = "new"
	StatusPlaced  = "placed"
	StatusSkipped = "skipped"
	StatusFailed  = "failed"
)

type NewOrder struct {
	RuleId    string
	Pair      string
	Market    string
	Action    string
	Behavior  string
	Price     string
	Quantity  Quantity
	Deadlines []Deadline
}

type Order struct {
	DSKey       *datastore.Key `datastore:"__key__"`
	Status      string         `datastore:"status"`
	RuleId      string         `datastore:"rule_id"`
	Pair        string         `datastore:"pair"`
	Market      string         `datastore:"market"`
	Action      string         `datastore:"action,noindex"`
	Behavior    string         `datastore:"behavior,noindex"`
	Price       string         `datastore:"price,noindex"`
	Quantity    Quantity       `datastore:"quantity,noindex"`
	Deadlines   []Deadline     `datastore:"deadlines,noindex"`
	CreatedAt   time.Time      `datastore:"created_at,noindex"`
	CompletedAt time.Time      `datastore:"completed_at,noindex"`
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
