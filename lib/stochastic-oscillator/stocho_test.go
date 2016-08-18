package stocho_test

import (
	"testing"

	. "github.com/apourchet/investment/lib/stochastic-oscillator"
	"github.com/stretchr/testify/assert"
)

func TestNewStochasticOscillator(t *testing.T) {
	so := NewStochasticOscillator(5, 3, 3)
	so.Step(1.)
	so.Step(1.)
	so.Step(1.)
	so.Step(1.)
	so.Step(1.)
	assert.InEpsilon(t, 1., so.GetD(), 0.0001)
}
