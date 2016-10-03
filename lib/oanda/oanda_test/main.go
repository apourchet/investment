package main

import (
	"fmt"

	"github.com/apourchet/investment/lib/oanda"
	"github.com/apourchet/investment/protos"
)

func main() {
	broker := oanda.NewPracticeBroker()
	accounts, err := broker.GetAccounts(nil, nil, nil)
	fmt.Println(err)
	fmt.Println(accounts)
	prices, err := broker.GetPrices(nil, &protos.PriceListReq{[]string{"EUR_USD"}}, nil)
	fmt.Println(err)
	fmt.Println(prices)
}
