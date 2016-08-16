package ema

type Ema struct {
	Alpha float64
	Steps int
	value float64
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
		e.value = val
	} else {
		e.value = val*e.Alpha + e.value*(1-e.Alpha)
	}
	e.Steps += 1
	return e.value
}

func (e *Ema) Compute() float64 {
	return e.value
}

func (e *Ema) ComputeNext(val float64) float64 {
	if e.Steps == 0 {
		return val
	}
	return val*e.Alpha + e.value*(1-e.Alpha)
}
