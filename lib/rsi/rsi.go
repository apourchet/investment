package rsi

import (
	"github.com/apourchet/investment/lib/sma"
)

type avg interface {
	Steps() int
	Compute() float64
	Step(float64) float64
}

type Rsi struct {
	PrevClose float64
	U         avg
	D         avg
}

// Implemented as described on wikipedia
// https://en.wikipedia.org/wiki/Relative_strength_index
func NewRsi(period int) *Rsi {
	rsi := &Rsi{}
	rsi.U = sma.NewSma(period)
	rsi.D = sma.NewSma(period)
	return rsi
}

func (rsi *Rsi) Step(val float64) float64 {
	if rsi.U.Steps() == 0 {
		rsi.PrevClose = val
	}
	if val > rsi.PrevClose {
		rsi.U.Step(val - rsi.PrevClose)
		rsi.D.Step(0)
	} else {
		rsi.U.Step(0)
		rsi.D.Step(rsi.PrevClose - val)
	}

	return 100. - 100./(1+rsi.rs())
}

func (rsi *Rsi) rs() float64 {
	return rsi.U.Compute() / rsi.D.Compute()
}
