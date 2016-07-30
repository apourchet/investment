package main

import (
	"fmt"
	. "github.com/apourchet/investment/lib"
	"time"
)

func mystrat(s *StrategyInterface) {
	for {
		fmt.Println("Sending Quote Request")
		s.QuoteRequest <- QuoteRequest{"EURUSD", 0}
		p := <-s.QuoteResponse
		fmt.Println("QuoteResponse: " + p.String())
		time.Sleep(time.Second * 2)
	}
}

func main() {
	strat := NewStrategy(mystrat)

	trader := NewTrader(":1026", strat, 10000)
	trader.Start()
}
