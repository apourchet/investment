package lma

type Lma struct {
	N      int
	Steps  int
	values []float64
}

func NewLma(n int) *Lma {
	return &Lma{n, 0, make([]float64, n)}
}

func (l *Lma) Step(val float64) float64 {
	if l.Steps == 0 {
		for i := 0; i < l.N; i++ {
			l.values[i] = val
		}
	}

	l.rotate()
	l.values[0] = val
	l.Steps += 1
	return l.Compute()
}

func (l *Lma) Compute() float64 {
	s := 0.
	for i, f := range l.values {
		s += f * float64(l.N-i)
	}
	return s / float64(sum(l.N))
}

func (l *Lma) rotate() {
	for i := l.N - 1; i > 0; i-- {
		l.values[i] = l.values[i-1]
	}
}

func sum(n int) int {
	s := 0
	for i := 1; i <= n; i++ {
		s += i
	}
	return s
}
