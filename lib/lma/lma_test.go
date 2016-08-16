package lma_test

import (
	"testing"

	"math"

	. "github.com/apourchet/investment/lib/lma"
)

const EPSILON = 0.000001

func assertDoubleEqual(t *testing.T, a, b float64) {
	if math.Abs(a-b) > EPSILON {
		t.Logf("%f != %f", a, b)
		t.Fail()
	}
}

func TestBasic(t *testing.T) {
	lma := NewLma(4)
	lma.Step(1.)
	lma.Step(2.)
	lma.Step(3.)
	lma.Step(4.) // (1 + 4 + 9 + 16) = 30 / 10 = 3
	assertDoubleEqual(t, lma.Compute(), 3.)
}
