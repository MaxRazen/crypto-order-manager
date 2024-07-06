package storage

import (
	"context"
	"crypto/sha256"
	"encoding/binary"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

const (
	Namespace = "crypto-tracker"
)

type ClientOptions struct {
	ProjectId      string
	ServiceKeyFile string
}

// Creates a new instance of datastore
func New(ctx context.Context, options ClientOptions) (*datastore.Client, error) {
	opts := option.WithCredentialsFile(options.ServiceKeyFile)

	return datastore.NewClient(ctx, options.ProjectId, opts)
}

// Creates a datastore key for a new record
func NewIDKey(kind string, id int64) *datastore.Key {
	return &datastore.Key{
		Kind:      kind,
		ID:        id,
		Namespace: Namespace,
	}
}

// Generates pseudo unique ID as int64 with min length 16 digits
func GenerateID() int64 {
	u := uuid.New()
	hash := sha256.Sum256(u[:])
	id := int64(binary.BigEndian.Uint64(hash[:8]))
	// Ensure the ID is greater than 0
	if id < 0 {
		id *= -1
	}
	// Ensure the ID has at least 16 digits
	if id < 10e15 {
		id += 10e15
	}
	return id
}
