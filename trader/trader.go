package main

import (
	"fmt"
	. "github.com/apourchet/investment/lib"
	"time"
)

func mystrat(tr *Trader) {
	for {
		q := tr.Broker.GetQuote(QuoteRequest{"EURUSD", 0})
		fmt.Println("QuoteResponse: " + q.String())
		time.Sleep(time.Second * 5)
	}
}

func main() {
	strat := NewStrategy(mystrat)

	trader := NewTrader(":1026", strat, 10000)
	trader.Start()
	for trader.Margin != 0 {
		fmt.Printf("Margin: %d\n", trader.Margin)
		time.Sleep(time.Minute)
	}
}
