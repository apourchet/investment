package test

import (
	"github.com/apourchet/investment"
	"github.com/apourchet/investment/protos"
	"testing"
	"log"
	"math"
)

const EPSILON = 0.000001

// Simple test with result taken from online
// http://www.investopedia.com/articles/forex/12/calculating-profits-and-losses-of-forex-trades.asp
func TestLikeOnline(t *testing.T) {
	sim := invt.NewAccountSimulator()
	sim.Account.Balance = 0
	sim.Buy(protos.InstrumentID_EURUSD, 100000, 1.6240)
	sim.Sell(protos.InstrumentID_EURUSD, 100000, 1.6255)
	sim.Simulate()
	if math.Abs(sim.Account.Balance - 150.0) > EPSILON {
		t.Fail()
	}
}
