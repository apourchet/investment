package main

import (
	"github.com/apourchet/investment/lib"
)

type DefaultBroker struct {
}

func (b DefaultBroker) GetQuote(pair string) invt.Data {
	return nil
}

func main() {
	broker := DefaultBroker{}
	invt.StartBroker(broker) // Blocking
}
