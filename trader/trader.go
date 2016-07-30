package trader

import (
	. "github.com/apourchet/investment/lib"
	"time"
)

func mystrat(ac ActionChannel) {
	for {
		ac <- QuoteAction()
		time.Sleep(time.Second * 2)
	}
}

func main() {
	strat := NewStrategy()

	trader := NewTrader(":1026", strat)
	trader.InitMargin(10000)
	trader.Start()
}
