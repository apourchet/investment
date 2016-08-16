package test

import (
	"math"
	"testing"

	. "github.com/apourchet/investment"
)

const EPSILON = 0.000001

func assertDoubleEqual(t *testing.T, a, b float64) {
	if math.Abs(a-b) > EPSILON {
		t.Logf("%f != %f", a, b)
		t.Fail()
	}
}

func checkBalance(t *testing.T, a *Account, expectedBalance float64) {
	assertDoubleEqual(t, a.Balance, expectedBalance)
}

func checkRealizedPl(t *testing.T, a *Account, expected float64) {
	assertDoubleEqual(t, a.RealizedPl, expected)
}

// Simple test with result taken from online
// http://www.investopedia.com/articles/forex/12/calculating-profits-and-losses-of-forex-trades.asp
func TestLikeOnline1(t *testing.T) {
	a := NewAccount(0)
	Buy(a, "", 100000, 1.6240)
	Sell(a, "", 100000, 1.6255)
	checkBalance(t, a, 150.0)
	checkRealizedPl(t, a, 150.)
}

// Simple test with result taken from online
// http://www.investopedia.com/articles/forex/12/calculating-profits-and-losses-of-forex-trades.asp
func TestLikeOnline2(t *testing.T) {
	a := NewAccount(0)
	Buy(a, "", 100000, 1.6240)
	Sell(a, "", 100000, 1.6220)
	checkBalance(t, a, -200.0)
	checkRealizedPl(t, a, -200.)
}

func TestReduceSimpleLong(t *testing.T) {
	a := NewAccount(0)

	Buy(a, "", 10, 1)
	Sell(a, "", 5, 2)
	checkBalance(t, a, 0)
	checkRealizedPl(t, a, 5)
}

func TestReduceSimpleShort(t *testing.T) {
	a := NewAccount(0)

	Sell(a, "", 10, 2)
	Buy(a, "", 5, 1)
	checkBalance(t, a, -5)
	checkRealizedPl(t, a, 5)
}

func TestReduceTwiceLong(t *testing.T) {
	a := NewAccount(0)

	Buy(a, "", 10, 1)
	Sell(a, "", 2, 2)
	Sell(a, "", 3, 3)
	checkBalance(t, a, 3)
	checkRealizedPl(t, a, 8)
}

func TestReduceTwiceShort(t *testing.T) {
	a := NewAccount(0)

	Sell(a, "", 10, 5)
	Buy(a, "", 5, 2)
	Buy(a, "", 3, 1)
	checkBalance(t, a, 17)
	checkRealizedPl(t, a, 27)
}

func TestReduceTwiceClose(t *testing.T) {
	a := NewAccount(0)

	Sell(a, "", 10, 5)
	Buy(a, "", 5, 2)
	Buy(a, "", 3, 1)
	Buy(a, "", 2, 1)
	checkBalance(t, a, 35)

	Buy(a, "", 10, 5)
	Sell(a, "", 5, 2)
	Sell(a, "", 3, 1)
	Sell(a, "", 2, 1)

	checkBalance(t, a, 0)
	checkRealizedPl(t, a, 0)
}

func TestFlipLong(t *testing.T) {
	a := NewAccount(0)

	Buy(a, "", 20, 2)
	Sell(a, "", 40, 5)
	checkBalance(t, a, -40)
	checkRealizedPl(t, a, 60)
}

func TestFlipShort(t *testing.T) {
	a := NewAccount(0)

	Sell(a, "", 10, 5)
	Buy(a, "", 20, 2)
	checkBalance(t, a, 10)
	checkRealizedPl(t, a, 30)
}

func TestExtendClose(t *testing.T) {
	a := NewAccount(0)

	Buy(a, "", 10, 1)
	Buy(a, "", 10, 2) // 20 units @ 1.5
	Sell(a, "", 10, 1)
	Sell(a, "", 10, 2)
	checkBalance(t, a, 0)

	Sell(a, "", 10, 1)
	Sell(a, "", 10, 2)
	Buy(a, "", 10, 1)
	Buy(a, "", 10, 2) // 20 units @ 1.5
	checkBalance(t, a, 0)

	Buy(a, "", 10, 1)
	Buy(a, "", 20, 2) // 30 units @ 5/3
	assertDoubleEqual(t, a.OpenPositions[""].Value(), 50.)
	Sell(a, "", 10, 1)
	assertDoubleEqual(t, a.OpenPositions[""].Value(), 20.*5./3.)
	Sell(a, "", 20, 1)
	if a.OpenPositions[""] != nil {
		t.Fail()
	}
	checkBalance(t, a, -20.)
}
