package sma

import (
	"github.com/apourchet/investment/lib/sliding-window"
)

type Sma struct {
	n      int
	values slidwin.SlidingWindow
	Steps  int
}

func NewSma(N int) *Sma {
	return &Sma{
		n:      N,
		values: slidwin.NewSlidingWindow(N),
		Steps:  0,
	}
}

func (s *Sma) Compute() float64 {
	return s.sum() / float64(s.n)
}

func (s *Sma) Step(val float64) float64 {
	if s.Steps == 0 {
		for i := range s.values {
			s.values[i] = val
		}
		s.Steps += 1
		return s.Compute()
	}
	s.values.Push(val)
	s.Steps += 1
	return s.Compute()
}

func (s *Sma) sum() float64 {
	sum := 0.
	for _, f := range s.values {
		sum += f.(float64)
	}
	return sum
}
