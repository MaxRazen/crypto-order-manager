package market

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/MaxRazen/crypto-order-manager/internal/storage"
)

type PlacedOrderDatastoreRepository struct {
	ds *datastore.Client
}

func NewRepository(ds *datastore.Client) *PlacedOrderDatastoreRepository {
	return &PlacedOrderDatastoreRepository{
		ds: ds,
	}
}

func (r *PlacedOrderDatastoreRepository) Create(ctx context.Context, o *PlacedOrder) error {
	dsk := storage.NewIDKey(PlacedOrderKind, storage.GenerateID())
	dsk, err := r.ds.Put(ctx, dsk, o)
	if err == nil {
		o.DSKey = dsk
	}
	return err
}

func (r *PlacedOrderDatastoreRepository) FetchAllUncompleted(ctx context.Context) ([]PlacedOrder, error) {
	var orders []PlacedOrder
	q := datastore.NewQuery(PlacedOrderKind).FilterField("status", "=", StatusNew)
	_, err := r.ds.GetAll(ctx, q, &orders)
	return orders, err
}

func (r *PlacedOrderDatastoreRepository) UpdateStatus(ctx context.Context, o *PlacedOrder, status string) error {
	_, err := r.ds.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		if err := tx.Get(o.DSKey, &o); err != nil {
			return err
		}

		o.Status = status

		if _, err := tx.Put(o.DSKey, o); err != nil {
			return err
		}
		return nil
	})

	return err
}
