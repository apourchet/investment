package invt

import (
	"fmt"
)

type Strategy struct {
	f StratFunction
}

type StratFunction func(*Trader)

func NewStrategy(f StratFunction) *Strategy {
	return &Strategy{f}
}

func (s *Strategy) Start(tr *Trader) {
	fmt.Println("Strategy Start")
	go s.f(tr)
}
