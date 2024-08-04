package app

import (
	"context"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/MaxRazen/crypto-order-manager/internal/binance"
	"github.com/MaxRazen/crypto-order-manager/internal/config"
	"github.com/MaxRazen/crypto-order-manager/internal/logger"
	"github.com/MaxRazen/crypto-order-manager/internal/market"
	"github.com/MaxRazen/crypto-order-manager/internal/order"
	"github.com/MaxRazen/crypto-order-manager/internal/storage"
	"github.com/MaxRazen/crypto-order-manager/internal/tracker"
)

type App struct {
	Storage      *datastore.Client
	Markets      *market.Collection
	OrderPlacer  *order.PlacementService
	OrderTracker *tracker.Tracker
}

func New(ctx context.Context, log *logger.Logger, cfg *config.Config, envVars *config.EnvVariables) (*App, error) {
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
	// Init repositories

	orderRepo := order.NewRepository(ds)
	placedOrderRepo := market.NewRepository(ds)

	// -------------------------------------------------------------------------
	// Init market clients

	markets := market.NewCollection()
	markets.Add(binance.New(
		envVars.BINANCE_API_KEY,
		envVars.BINANCE_SECRET_KEY,
		envVars.BINANCE_BASE_URL,
		cfg.DryRun,
	))

	// -------------------------------------------------------------------------
	// Init Order placement service

	ordPlacer := order.NewPlacementService(log, orderRepo, placedOrderRepo, markets)

	// -------------------------------------------------------------------------
	// Init Order tracker service

	tickInterval, _ := strconv.Atoi(envVars.TICKER_TICK_INTERVAL)
	checkInterval, _ := strconv.Atoi(envVars.TICKER_CHECK_INTERVAL)
	ordTracker := tracker.New(log, placedOrderRepo, markets, time.Duration(tickInterval)*time.Second, time.Duration(checkInterval)*time.Second)

	// -------------------------------------------------------------------------
	// Wrap everything into app

	app := App{
		Storage:      ds,
		Markets:      markets,
		OrderPlacer:  ordPlacer,
		OrderTracker: ordTracker,
	}

	return &app, nil
}
