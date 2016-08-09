package test

import "time"
import (
	"testing"

	. "github.com/apourchet/investment/lib/tradelogger"
)

func TestMain(t *testing.M) {
	logger := NewLogger("test", "/tmp/testlog")
	logger.Log(&Item{time.Now(), "TestTag", "TestMessage"})
}
