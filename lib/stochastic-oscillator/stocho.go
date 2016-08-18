package stocho

import "github.com/apourchet/investment/lib/sma"

type StochasticOscillator struct {
	values []float64 // Used for H5 and L5
	d      *sma.Sma
	dSlow  *sma.Sma
	Steps  int
}

func NewStochasticOscillator(n, d1, d2 int) *StochasticOscillator {
	return &StochasticOscillator{
		values: make([]float64, n),
		d:      sma.NewSma(d1),
		dSlow:  sma.NewSma(d2),
	}
}

func (so *StochasticOscillator) Step(val float64) float64 {
	if so.Steps == 0 {
		for i := range so.values {
			so.values[i] = val
		}
	}
	so.Steps += 1
	l := so.Low()
	h := so.High()
	k := (val - l) / (h - l)
	d := so.d.Step(k)
	return so.dSlow.Step(d)
}

func (so *StochasticOscillator) GetD() float64 {
	return so.d.Compute()
}

func (so *StochasticOscillator) GetDSlow() float64 {
	return so.dSlow.Compute()
}

// O(n) for getting minimum since we are using a ring
func (so *StochasticOscillator) Low() float64 {
	min := so.values[0]
	for _, f := range so.values {
		if min > f {
			min = f
		}
	}
	return min
}

// O(n) for getting maximum since we are using a ring
func (so *StochasticOscillator) High() float64 {
	max := so.values[0]
	for _, f := range so.values {
		if max < f {
			max = f
		}
	}
	return max
}

func (so *StochasticOscillator) rotate() {
	for i := len(so.values) - 1; i > 0; i-- {
		so.values[i] = so.values[i-1]
	}
}
