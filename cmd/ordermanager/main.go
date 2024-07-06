package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/MaxRazen/crypto-order-manager/internal/config"
	"github.com/MaxRazen/crypto-order-manager/internal/grpcserver"
	"github.com/MaxRazen/crypto-order-manager/internal/logger"
)

var build string = "devonly"

func main() {

	// -------------------------------------------------------------------------
	// Init Logger

	logLevel := logger.LevelInfo
	if build != "devonly" {
		logLevel = logger.LevelDebug
	}

	log := logger.New(os.Stdout, logLevel)
	ctx := context.Background()

	// -------------------------------------------------------------------------
	// Run application

	if err := run(ctx, log); err != nil {
		log.Fatal(ctx, "startup", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {

	// -------------------------------------------------------------------------
	// Init Config

	cfg, err := config.Parse(build)
	if err == config.ErrHelpWanted {
		return nil
	} else if err != nil {
		log.Fatal(ctx, err.Error())
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 1)

	// -------------------------------------------------------------------------
	// Run GRPC Server

	go func() {
		serverErrors <- grpcserver.Run(ctx, log, cfg.Port)
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		// TODO: shutdown grpc server and release resources
	}

	return nil
}
