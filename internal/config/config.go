package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/ardanlabs/conf/v3"
	"github.com/joho/godotenv"
)

var ErrHelpWanted = conf.ErrHelpWanted

type Config struct {
	conf.Version
	Port   string `conf:"default:50051" desc:"GRPC Server port"`
	DryRun bool   `conf:"default:false" desc:"Run service in dry run mode. Received orders won't be places."`
}

type EnvVariables struct {
	GRPC_AUTHORIZATION_TOKEN string
	GCP_PROJECT_ID           string
	GCP_SERVICE_KEY_FILE     string
	BINANCE_BASE_URL         string
	BINANCE_API_KEY          string
	BINANCE_SECRET_KEY       string
}

const (
	GRPC_AUTHORIZATION_TOKEN = "GRPC_AUTHORIZATION_TOKEN"
	GCP_PROJECT_ID           = "GCP_PROJECT_ID"
	GCP_SERVICE_KEY_FILE     = "GCP_SERVICE_KEY_FILE"
	BINANCE_BASE_URL         = "BINANCE_BASE_URL"
	BINANCE_API_KEY          = "BINANCE_API_KEY"
	BINANCE_SECRET_KEY       = "BINANCE_SECRET_KEY"
)

func Parse(build string) (*Config, error) {
	cfg := &Config{
		Version: conf.Version{
			Build: build,
		},
	}
	help, err := conf.Parse("", cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil, ErrHelpWanted
		}
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}

func ParseEnvs(filename string) (*EnvVariables, error) {
	if filename != "" {
		err := godotenv.Load(filename)

		if err != nil {
			return nil, err
		}
	}

	envVars := EnvVariables{
		GRPC_AUTHORIZATION_TOKEN: os.Getenv(GRPC_AUTHORIZATION_TOKEN),
		GCP_PROJECT_ID:           os.Getenv(GCP_PROJECT_ID),
		GCP_SERVICE_KEY_FILE:     os.Getenv(GCP_SERVICE_KEY_FILE),
		BINANCE_BASE_URL:         os.Getenv(BINANCE_BASE_URL),
		BINANCE_API_KEY:          os.Getenv(BINANCE_API_KEY),
		BINANCE_SECRET_KEY:       os.Getenv(BINANCE_SECRET_KEY),
	}

	return &envVars, nil
}
