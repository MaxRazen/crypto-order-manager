package app

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/MaxRazen/crypto-order-manager/internal/binance"
	"github.com/MaxRazen/crypto-order-manager/internal/config"
	"github.com/MaxRazen/crypto-order-manager/internal/market"
	"github.com/MaxRazen/crypto-order-manager/internal/storage"
)

type App struct {
	Storage *datastore.Client
	Markets map[string]market.MarketClient
}

func New(ctx context.Context, cfg *config.Config, envVars *config.EnvVariables) (*App, error) {
	// -------------------------------------------------------------------------
	// Init datastore by GCP

	ds, err := storage.New(ctx, storage.ClientOptions{
		ProjectId:      envVars.GCP_PROJECT_ID,
		ServiceKeyFile: envVars.GCP_SERVICE_KEY_FILE,
	})

	if err != nil {
		return nil, err
	}

	// -------------------------------------------------------------------------
	// Init market clients

	markets := make(map[string]market.MarketClient)
	markets["binance"] = binance.New(
		envVars.BINANCE_API_KEY,
		envVars.BINANCE_SECRET_KEY,
		envVars.BINANCE_BASE_URL,
		cfg.DryRun,
	)

	// -------------------------------------------------------------------------
	// Wrap everything into app

	app := App{
		Storage: ds,
		Markets: markets,
	}

	return &app, nil
}
