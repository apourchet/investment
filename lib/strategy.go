package invt

import (
	"fmt"
)

type StrategyInterface struct {
	QuoteRequest chan QuoteRequest
	BuyOrder     chan string
	SellOrder    chan string

	QuoteResponse chan Quote
}

type Strategy struct {
	f     StratFunction
	inter *StrategyInterface
}

type StratFunction func(*StrategyInterface)

func NewStrategy(f StratFunction) *Strategy {
	inter := StrategyInterface{}
	inter.QuoteRequest = make(chan QuoteRequest, 10)
	inter.BuyOrder = make(chan string, 10)
	inter.SellOrder = make(chan string, 10)
	inter.QuoteResponse = make(chan Quote, 10)
	return &Strategy{f, &inter}
}

func (s *Strategy) Start() {
	fmt.Println("Strategy Start")
	go s.f(s.inter)
}
