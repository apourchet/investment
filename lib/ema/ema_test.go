package ema_test

import (
	"testing"

	. "github.com/apourchet/investment/lib/ema"
)

func TestBasic(t *testing.T) {
	ema := NewEma(AlphaFromN(10))
	ema.Step(1.)
	ema.Step(1.)
	ema.Step(1.)
	ema.Step(1.)
	ema.Step(1.)
	if ema.Compute() != 1. {
		t.Fail()
	}
}
