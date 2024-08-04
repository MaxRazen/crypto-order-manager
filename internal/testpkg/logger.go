package testpkg

import (
	"bytes"

	"github.com/MaxRazen/crypto-order-manager/internal/logger"
)

func NewLogger() (*logger.Logger, *bytes.Buffer) {
	logBufer := new(bytes.Buffer)
	log := logger.New(logBufer, logger.LevelDebug)

	return log, logBufer
}
