package invt

import (
	"fmt"
)

type Strategy func(*Trader)

func (s Strategy) Start(tr *Trader) {
	fmt.Println("Strategy Start")
	go s(tr)
}
