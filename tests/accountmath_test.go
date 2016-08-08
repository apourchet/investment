package test

import (
	"math"
	"testing"

	. "github.com/apourchet/investment"
)

const EPSILON = 0.000001

func checkBalance(t *testing.T, a *Account, expectedBalance float64) {
	if math.Abs(a.Balance-expectedBalance) > EPSILON {
		t.Fail()
	}
}

// Simple test with result taken from online
// http://www.investopedia.com/articles/forex/12/calculating-profits-and-losses-of-forex-trades.asp
func TestLikeOnline1(t *testing.T) {
	a := NewAccount(0)
	Buy(a, "", 100000, 1.6240)
	Sell(a, "", 100000, 1.6255)
	checkBalance(t, a, 150.0)
}

// Simple test with result taken from online
// http://www.investopedia.com/articles/forex/12/calculating-profits-and-losses-of-forex-trades.asp
func TestLikeOnline2(t *testing.T) {
	a := NewAccount(0)
	Buy(a, "", 100000, 1.6240)
	Sell(a, "", 100000, 1.6220)
	checkBalance(t, a, -200.0)
}

func TestReduceSimpleLong(t *testing.T) {
	a := NewAccount(0)

	Buy(a, "", 10, 1)
	Sell(a, "", 5, 2)
	checkBalance(t, a, 0)
}

func TestReduceSimpleShort(t *testing.T) {
	a := NewAccount(0)

	Sell(a, "", 10, 2)
	Buy(a, "", 5, 1)
	checkBalance(t, a, -5)
}

func TestReduceTwiceLong(t *testing.T) {
	a := NewAccount(0)

	Buy(a, "", 10, 1)
	Sell(a, "", 2, 2)
	Sell(a, "", 3, 3)
	checkBalance(t, a, 3)
}

func TestReduceTwiceShort(t *testing.T) {
	a := NewAccount(0)

	Sell(a, "", 10, 5)
	Buy(a, "", 5, 2)
	Buy(a, "", 3, 1)
	checkBalance(t, a, 17)
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
}

func TestFlipLong(t *testing.T) {
	a := NewAccount(0)

	Buy(a, "", 20, 2)
	Sell(a, "", 40, 5)
	checkBalance(t, a, -40)
}

func TestFlipShort(t *testing.T) {
	a := NewAccount(0)

	Sell(a, "", 10, 5)
	Buy(a, "", 20, 2)
	checkBalance(t, a, 10)
}
