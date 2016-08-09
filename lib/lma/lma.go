package lma

type Lma struct {
	N      int
	Steps  int
	Values []float64
}

func NewLma(n int) *Lma {
	return &Lma{n, 0, make([]float64, n)}
}

func (l *Lma) Step(val float64) float64 {
	if l.Steps == 0 {
		for i := 0; i < l.N; i++ {
			l.Values[i] = val
		}
		return val
	}

	l.rotate()
	l.Values[0] = val
	return l.compute()
}

func (l *Lma) rotate() {
	for i := 1; i < l.N; i++ {
		l.Values[i] = l.Values[i-1]
	}
}

func (l *Lma) compute() float64 {
	s := 0.
	for i := 0; i < l.N; i++ {
		s += l.Values[i] * float64(l.N-i)
	}
	return s / float64(sum(l.N))
}

func sum(n int) int {
	s := 0
	for i := 1; i <= n; i++ {
		s += i
	}
	return s
}
