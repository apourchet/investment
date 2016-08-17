package sma

type Sma struct {
	n      int
	values []float64
	Steps  int
}

func NewSma(N int) *Sma {
	return &Sma{n: N, values: make([]float64, N), Steps: 0}
}

func (s *Sma) Compute() float64 {
	return s.sum() / float64(s.n)
}

func (s *Sma) Step(val float64) float64 {
	if s.Steps == 0 {
		for i := range s.values {
			s.values[i] = val
		}
	}
	s.rotate()
	s.values[0] = val
	s.Steps += 1
	return s.Compute()
}

func (s *Sma) sum() float64 {
	sum := 0.
	for _, f := range s.values {
		sum += f
	}
	return sum
}

func (s *Sma) rotate() {
	for i := len(s.values) - 1; i > 0; i-- {
		s.values[i] = s.values[i-1]
	}
}
