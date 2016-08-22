package slidwin

// This will allow us to easily move our training example windows
// This could also potentially be used for
// moving averages (eg: simple, linear)
// but also for a stochastic oscillator
type SlidingWindow []interface{}

func NewSlidingWindow(n int) SlidingWindow {
	return make([]interface{}, n)
}

// s is ordered from most recent to least recent
func (s SlidingWindow) Push(x interface{}) {
	s.rotate()
	s[0] = x
}

func (s SlidingWindow) rotate() {
	for i := len(s) - 1; i > 0; i-- {
		s[i] = s[i-1]
	}
}
