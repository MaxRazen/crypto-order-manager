package testpkg

import (
	"bytes"
	"testing"

	"github.com/MaxRazen/crypto-order-manager/internal/logger"
)

func NewLogger() (*logger.Logger, *bytes.Buffer) {
	logBufer := new(bytes.Buffer)
	log := logger.New(logBufer, logger.LevelDebug)

	return log, logBufer
}

func VerboseOutput(t *testing.T, logStack string) {
	t.Logf("\n----------- LOGS -----------\n%s----------- LOGS END -----------\n", logStack)
}
