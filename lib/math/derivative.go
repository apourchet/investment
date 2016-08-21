package invt_math

type FirstDerivative struct {
	Last  float64
	Value float64
	Steps int
}

type SecondDerivative struct {
	LastD   *FirstDerivative
	LastVal float64
	Value   float64
	Steps   int
}

func NewFirstDerivative() *FirstDerivative {
	return &FirstDerivative{0, 0, 0}
}

func NewSecondDerivative() *SecondDerivative {
	return &SecondDerivative{NewFirstDerivative(), 0, 0, 0}
}

func (f *FirstDerivative) Step(cur float64) float64 {
	if f.Steps == 0 {
		f.Value = cur
	} else {
		f.Value = cur - f.Last
	}
	f.Last = cur
	f.Steps += 1
	return f.Value
}

func (s *SecondDerivative) Step(cur float64) float64 {
	if s.Steps == 0 {
		s.Value = 0.0
		s.LastD.Step(cur)
		s.LastVal = cur
	} else {
		lastD := s.LastD.Value
		curD := s.LastD.Step(cur)
		s.Value = curD - lastD
	}
	s.LastVal = cur
	s.Steps += 1
	return s.Value
}
