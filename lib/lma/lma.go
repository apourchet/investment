package lma

import "github.com/apourchet/investment/lib/sliding-window"

type Lma struct {
	N      int
	Steps  int
	values slidwin.SlidingWindow
}

func NewLma(n int) *Lma {
	return &Lma{
		N:      n,
		Steps:  0,
		values: slidwin.NewSlidingWindow(n),
	}
}

func (l *Lma) Step(val float64) float64 {
	if l.Steps == 0 {
		for i := 0; i < l.N; i++ {
			l.values[i] = val
		}
		l.Steps += 1
		return l.Compute()
	}
	l.values.Push(val)
	l.Steps += 1
	return l.Compute()
}

func (l *Lma) Compute() float64 {
	s := 0.
	for i, f := range l.values {
		s += f.(float64) * float64(l.N-i)
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
