package invt

import (
	"fmt"
)

type Trader struct {
	Broker   Broker
	Strategy *Strategy
	Margin   int
}

func NewTrader(brokerUrl string, strat *Strategy, margin int) *Trader {
	tr := Trader{}
	tr.Broker = NewBrokerClient(brokerUrl)
	tr.Strategy = strat
	tr.Margin = margin
	return &tr
}

func (tr *Trader) Start() {
	fmt.Println("Trader Start")
	tr.Strategy.Start(tr)
}
