package main

import (
	"fmt"
	"time"

	. "github.com/apourchet/investment/lib"
)

func mystrat(tr *Trader) {
	broker := tr.Broker
	for {
		q := broker.GetQuote(QuoteRequest{"EURUSD", 0, ""})
		fmt.Println("QuoteResponse: " + q.String())
		time.Sleep(time.Second * 5)
	}
}

func main() {
	trader := NewTrader(":1026", mystrat, 10000)
	trader.Start()
	for trader.Margin != 0 {
		fmt.Printf("Margin: %d\n", trader.Margin)
		time.Sleep(time.Minute)
	}
}
