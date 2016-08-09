package lma

type Lma struct {
	N      int
	Steps  int
	values []float64
	Value  float64
}

func NewLma(n int) *Lma {
	return &Lma{n, 0, make([]float64, n), 0}
}

func (l *Lma) Step(val float64) {
	if l.Steps == 0 {
		for i := 0; i < l.N; i++ {
			l.values[i] = val
		}
		l.Value = val
	}

	l.rotate()
	l.values[0] = val
	l.Value = l.compute()
}

func (l *Lma) rotate() {
	for i := 1; i < l.N; i++ {
		l.values[i] = l.values[i-1]
	}
}

func (l *Lma) compute() float64 {
	s := 0.
	for i := 0; i < l.N; i++ {
		s += l.values[i] * float64(l.N-i)
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
