package atr

import (
	"github.com/apourchet/investment/lib/ema"
	"math"
)

type ATR struct {
	avg *ema.Ema
}

func NewAtr(n int) *ATR {
	return &ATR{ema.NewEma(ema.AlphaFromN(n))}
}

func (a *ATR) Step(high, low, prevClose float64) float64 {
	tr := trueRange(high, low, prevClose)
	return a.avg.Step(tr)
}

func (a *ATR) Compute() float64 {
	return a.avg.Compute()
}

// https://en.wikipedia.org/wiki/Average_true_range
func trueRange(high, low, prevClose float64) float64 {
	return math.Max(math.Max(high-low, high-prevClose), low-prevClose)
}
