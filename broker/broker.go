package main

import (
	. "github.com/apourchet/investment/lib"
	"strconv"
)

type DefaultQuote struct {
	i int
}

type DefaultBroker struct{}

func (q DefaultQuote) String() string {
	return strconv.Itoa(q.i)
}

func (b DefaultBroker) GetQuote(pair string, lookback int) Quote {
	return DefaultQuote{42}
}

func main() {
	broker := DefaultBroker{}
	handler := BrokerHandler{broker}
	handler.Start()
}
