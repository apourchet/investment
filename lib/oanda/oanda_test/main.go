package main

import (
	"fmt"

	"github.com/apourchet/investment/lib/oanda"
)

func main() {
	broker := oanda.NewPracticeBroker()
	accounts, err := broker.GetAccounts(nil, nil, nil)
	fmt.Println(err)
	fmt.Println(accounts)
}
