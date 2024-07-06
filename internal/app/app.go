package app

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/MaxRazen/crypto-order-manager/internal/config"
	"github.com/MaxRazen/crypto-order-manager/internal/storage"
)

type App struct {
	Storage *datastore.Client
}

func New(ctx context.Context, cfg *config.Config, envVars *config.EnvVariables) (*App, error) {
	strg, err := storage.New(ctx, storage.ClientOptions{
		ProjectId:      envVars.GCP_PROJECT_ID,
		ServiceKeyFile: envVars.GCP_SERVICE_KEY_FILE,
	})

	if err != nil {
		return nil, err
	}

	app := App{
		Storage: strg,
	}

	return &app, nil
}
