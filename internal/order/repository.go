package order

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
)

type Repository struct {
	ds *datastore.Client
}

func NewRepository(ds *datastore.Client) *Repository {
	return &Repository{
		ds: ds,
	}
}
func (r *Repository) FetchAllUnplaced(ctx context.Context) ([]Order, error) {
	var orders []Order
	q := datastore.NewQuery(OrdersKind).FilterField("status", "=", StatusNew)
	_, err := r.ds.GetAll(ctx, q, &orders)
	return orders, err
}

func (r *Repository) UpdateStatus(ctx context.Context, o *Order, status string) error {
	_, err := r.ds.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		if err := tx.Get(o.DSKey, &o); err != nil {
			return err
		}

		o.Status = status
		if status == StatusPlaced {
			o.CompletedAt = time.Now()
		}

		if _, err := tx.Put(o.DSKey, o); err != nil {
			return err
		}
		return nil
	})

	return err
}
