package test

import (
	"math"
	"testing"

	"github.com/apourchet/investment"
	"github.com/apourchet/investment/protos"
)

const EPSILON = 0.000001

func checkBalance(t *testing.T, sim *invt.AccountSimulator, expectedBalance float64) {
	if math.Abs(sim.Account.Balance-expectedBalance) > EPSILON {
		t.Fail()
	}
}

// Simple test with result taken from online
// http://www.investopedia.com/articles/forex/12/calculating-profits-and-losses-of-forex-trades.asp
func TestLikeOnline1(t *testing.T) {
	sim := invt.NewAccountSimulator(0)
	sim.Buy(protos.InstrumentID_EURUSD, 100000, 1.6240)
	sim.Sell(protos.InstrumentID_EURUSD, 100000, 1.6255)
	sim.Simulate()
	checkBalance(t, sim, 150.0)
}

// Simple test with result taken from online
// http://www.investopedia.com/articles/forex/12/calculating-profits-and-losses-of-forex-trades.asp
func TestLikeOnline2(t *testing.T) {
	sim := invt.NewAccountSimulator(0)
	sim.Buy(protos.InstrumentID_EURUSD, 100000, 1.6240)
	sim.Sell(protos.InstrumentID_EURUSD, 100000, 1.6220)
	sim.Simulate()
	checkBalance(t, sim, -200.0)
}

func TestReduceSimpleLong(t *testing.T) {
	sim := invt.NewAccountSimulator(0)
	sim.Buy(protos.InstrumentID_EURUSD, 10, 1)
	sim.Sell(protos.InstrumentID_EURUSD, 5, 2)
	sim.Simulate()
	checkBalance(t, sim, 0)
}

func TestReduceSimpleShort(t *testing.T) {
	sim := invt.NewAccountSimulator(0)
	sim.SellNow(protos.InstrumentID_EURUSD, 10, 2)
	sim.BuyNow(protos.InstrumentID_EURUSD, 5, 1)
	checkBalance(t, sim, -5)
}

func TestReduceTwiceLong(t *testing.T) {
	sim := invt.NewAccountSimulator(0)
	sim.Buy(protos.InstrumentID_EURUSD, 10, 1)
	sim.Sell(protos.InstrumentID_EURUSD, 2, 2)
	sim.Sell(protos.InstrumentID_EURUSD, 3, 3)
	sim.Simulate()
	checkBalance(t, sim, 3)
}

func TestReduceTwiceShort(t *testing.T) {
	sim := invt.NewAccountSimulator(0)
	sim.SellNow(protos.InstrumentID_EURUSD, 10, 5) // -50
	sim.BuyNow(protos.InstrumentID_EURUSD, 5, 2)   // -10 (long 5 at 5 with +5*3)
	sim.BuyNow(protos.InstrumentID_EURUSD, 3, 1)   // (long 2 at 5 with +5*3+3*4)
	checkBalance(t, sim, 17)
}

func TestReduceTwiceClose(t *testing.T) {
	sim := invt.NewAccountSimulator(0)

	sim.SellNow(protos.InstrumentID_EURUSD, 10, 5) // -50
	sim.BuyNow(protos.InstrumentID_EURUSD, 5, 2)   // -10 (long 5 at 5 with +5*3)
	sim.BuyNow(protos.InstrumentID_EURUSD, 3, 1)   // (long 2 at 5 with +5*3+3*4)
	sim.BuyNow(protos.InstrumentID_EURUSD, 2, 1)   // (closed with +5*3+3*4+2*4)
	checkBalance(t, sim, 35)

	sim.BuyNow(protos.InstrumentID_EURUSD, 10, 5)
	sim.SellNow(protos.InstrumentID_EURUSD, 5, 2)
	sim.SellNow(protos.InstrumentID_EURUSD, 3, 1)
	sim.SellNow(protos.InstrumentID_EURUSD, 2, 1)
	checkBalance(t, sim, 0)
}

func TestFlipLong(t *testing.T) {
	sim := invt.NewAccountSimulator(0)
	sim.BuyNow(protos.InstrumentID_EURUSD, 20, 2)  // -40 (long 20@2)
	sim.SellNow(protos.InstrumentID_EURUSD, 40, 5) // -40 (short 20@5 with +20@3)
	checkBalance(t, sim, -40)
}

func TestFlipShort(t *testing.T) {
	sim := invt.NewAccountSimulator(0)
	sim.SellNow(protos.InstrumentID_EURUSD, 10, 5) // -50 (short 10@5)
	sim.BuyNow(protos.InstrumentID_EURUSD, 20, 2)  // +10 (long 10@2 with +10@3)
	checkBalance(t, sim, 10)
}
