package candelizer

import "github.com/apourchet/investment/lib/sliding-window"

type CandleInterface interface {
	Open() float64
	Close() float64
	High() float64
	Low() float64
}

type Candelizer struct {
	values slidwin.SlidingWindow
	Steps  int
}

func NewCandelizer(n int) *Candelizer {
	return &Candelizer{slidwin.NewSlidingWindow(n), 0}
}

func (c *Candelizer) Step(val float64) CandleInterface {
	if c.Steps == 0 {
		for i := range c.values {
			c.values[i] = val
		}
	}
	c.Steps += 1
	c.values.Push(val)
	return c
}

func (c *Candelizer) Open() float64 {
	return c.values[len(c.values)-1].(float64)
}

func (c *Candelizer) Close() float64 {
	return c.values[0].(float64)
}

// O(n) for getting minimum since we are using a ring
func (c *Candelizer) Low() float64 {
	min := c.values[0].(float64)
	for _, x := range c.values {
		f := x.(float64)
		if min > f {
			min = f
		}
	}
	return min
}

// O(n) for getting maximum since we are using a ring
func (c *Candelizer) High() float64 {
	max := c.values[0].(float64)
	for _, x := range c.values {
		f := x.(float64)
		if max < f {
			max = f
		}
	}
	return max
}

func (c *Candelizer) rotate() {
	for i := len(c.values) - 1; i > 0; i-- {
		c.values[i] = c.values[i-1]
	}
}
