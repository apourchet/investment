package sma_test

import (
	"testing"

	"math"

	. "github.com/apourchet/investment/lib/sma"
)

const EPSILON = 0.000001

func assertDoubleEqual(t *testing.T, a, b float64) {
	if math.Abs(a-b) > EPSILON {
		t.Logf("%f != %f", a, b)
		t.Fail()
	}
}

func TestBasic(t *testing.T) {
	sma := NewSma(10)
	sma.Step(1.)
	sma.Step(4.)
	sma.Step(7.)
	sma.Step(10.)
	sma.Step(9.)
	sma.Step(3.)
	sma.Step(8.)
	sma.Step(5.)
	sma.Step(6.)
	sma.Step(2.)
	assertDoubleEqual(t, sma.Compute(), 5.5)
}
