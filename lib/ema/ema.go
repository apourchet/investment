package ema

type Ema struct {
	Alpha float64
	Value float64
	Steps int
}

func NewEma(alpha float64) *Ema {
	return &Ema{alpha, 0, 0}
}

func AlphaFromN(N int) float64 {
	n := float64(N)
	return 2. / (n + 1.)
}

func (e *Ema) Step(val float64) float64 {
	if e.Steps == 0 {
		e.Value = val
	} else {
		e.Value = val*e.Alpha + e.Value*(1-e.Alpha)
	}
	e.Steps += 1
	return e.Value
}
