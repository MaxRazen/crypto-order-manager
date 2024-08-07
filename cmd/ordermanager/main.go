package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/MaxRazen/crypto-order-manager/internal/app"
	"github.com/MaxRazen/crypto-order-manager/internal/config"
	"github.com/MaxRazen/crypto-order-manager/internal/grpcserver"
	"github.com/MaxRazen/crypto-order-manager/internal/logger"
)

var build string = "devonly"

func main() {

	// -------------------------------------------------------------------------
	// Init Logger

	logLevel := logger.LevelDebug
	if build != "devonly" {
		logLevel = logger.LevelInfo
	}

	log := logger.New(os.Stdout, logLevel)
	ctx := context.Background()

	// -------------------------------------------------------------------------
	// Run application

	if err := run(ctx, log); err != nil {
		println("exiting with error")
		log.Fatal(ctx, "startup", err)
	}

	os.Exit(0)
}

func run(ctx context.Context, log *logger.Logger) error {

	// -------------------------------------------------------------------------
	// Init Config

	cfg, err := config.Parse(build)
	if err == config.ErrHelpWanted {
		return nil
	} else if err != nil {
		return err
	}

	envVars, err := config.ParseEnvs(cfg.EnvPath)
	if err != nil {
		return err
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 1)
	trackerErrors := make(chan error, 1)

	// -------------------------------------------------------------------------
	// Init Application

	app, err := app.New(ctx, log, cfg, envVars)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	// -------------------------------------------------------------------------
	// Run Order placement service

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.OrderPlacer.Init(ctx); err != nil {
			log.Fatal(ctx, "ps: initializing error", "error", err.Error())
		}

		app.OrderPlacer.Run(ctx, app.OrderTracker.GetInputChan())
	}()

	// -------------------------------------------------------------------------
	// Run Order Tracker service

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.OrderTracker.Init(ctx); err != nil {
			log.Fatal(ctx, "tracker: initializing error", "error", err.Error())
		}

		trackerErrors <- app.OrderTracker.Run(ctx)
	}()

	// -------------------------------------------------------------------------
	// Run GRPC Server

	go func() {
		serverErrors <- grpcserver.Run(ctx, log, app, envVars.GRPC_AUTHORIZATION_TOKEN, cfg.Port)
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case err := <-trackerErrors:
		return fmt.Errorf("tracker error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		// TODO: shutdown grpc server and release resources
		app.OrderTracker.Stop(ctx)
		app.OrderPlacer.Stop()
		wg.Wait()

		// close storage
		if err := app.Storage.Close(); err != nil {
			return err
		}
	}

	return nil
}
