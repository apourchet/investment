package test

import (
	"math"
	"testing"

	"github.com/apourchet/investment"
	"github.com/apourchet/investment/protos"
)

const EPSILON = 0.000001

// Simple test with result taken from online
// http://www.investopedia.com/articles/forex/12/calculating-profits-and-losses-of-forex-trades.asp
func TestLikeOnline1(t *testing.T) {
	sim := invt.NewAccountSimulator()
	sim.Account.Balance = 0
	sim.Buy(protos.InstrumentID_EURUSD, 100000, 1.6240)
	sim.Sell(protos.InstrumentID_EURUSD, 100000, 1.6255)
	sim.Simulate()
	if math.Abs(sim.Account.Balance-150.0) > EPSILON {
		t.Fail()
	}
}

// Simple test with result taken from online
// http://www.investopedia.com/articles/forex/12/calculating-profits-and-losses-of-forex-trades.asp
func TestLikeOnline2(t *testing.T) {
	sim := invt.NewAccountSimulator()
	sim.Account.Balance = 0
	sim.Buy(protos.InstrumentID_EURUSD, 100000, 1.6240)
	sim.Sell(protos.InstrumentID_EURUSD, 100000, 1.6220)
	sim.Simulate()
	if math.Abs(sim.Account.Balance-(-200.0)) > EPSILON {
		t.Fail()
	}
}

func TestReduceSimpleLong(t *testing.T) {
	sim := invt.NewAccountSimulator()
	sim.Account.Balance = 0
	sim.Buy(protos.InstrumentID_EURUSD, 10, 1)
	sim.Sell(protos.InstrumentID_EURUSD, 5, 2)
	sim.Simulate()
	if math.Abs(sim.Account.Balance-0) > EPSILON {
		t.Fail()
	}
}

func TestReduceSimpleShort(t *testing.T) {
	sim := invt.NewAccountSimulator()
	sim.Account.Balance = 0
	sim.SellNow(protos.InstrumentID_EURUSD, 10, 2)
	sim.BuyNow(protos.InstrumentID_EURUSD, 5, 1)
	if math.Abs(sim.Account.Balance+5) > EPSILON {
		t.Fail()
	}
}

func TestReduceTwiceLong(t *testing.T) {
	sim := invt.NewAccountSimulator()
	sim.Account.Balance = 0
	sim.Buy(protos.InstrumentID_EURUSD, 10, 1)
	sim.Sell(protos.InstrumentID_EURUSD, 2, 2)
	sim.Sell(protos.InstrumentID_EURUSD, 3, 3)
	sim.Simulate()
	if math.Abs(sim.Account.Balance-3) > EPSILON {
		t.Fail()
	}
}

func TestReduceTwiceShort(t *testing.T) {
	sim := invt.NewAccountSimulator()
	sim.Account.Balance = 0
	sim.SellNow(protos.InstrumentID_EURUSD, 10, 5) // -50
	sim.BuyNow(protos.InstrumentID_EURUSD, 5, 2)   // -10 (long 5 at 5)
	sim.BuyNow(protos.InstrumentID_EURUSD, 3, 2)   // (long 2 at 5)
	if math.Abs(sim.Account.Balance-14) > EPSILON {
		t.Fail()
	}
}
