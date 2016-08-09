package lma

type Lma struct {
	N      int
	Steps  int
	Values []float64
}

func sum(a, b int) int {
	s := 0
	for i := a; i <= b; i++ {
		s += i
	}
	return s
}

func NewLma(n int) *Lma {
	return &Lma{n, 0, 0}
}

func (e *Lma) Step(val float64) float64 {
	// TODO
	return 0.
}
