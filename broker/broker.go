package main

import (
	. "github.com/apourchet/investment/lib"
	"time"
)

type DefaultBroker struct{}

func (b DefaultBroker) GetQuote(qr QuoteRequest) Quote {
	return Quote{42., 43., time.Now().UnixNano(), false, nil}
}

func main() {
	broker := DefaultBroker{}
	handler := BrokerHandler{broker}
	handler.Start()
}
