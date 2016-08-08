package invt

import "github.com/apourchet/investment/protos"

type Quote protos.Quote

type QuoteContext map[string]*Quote

func (qc *QuoteContext) Get(instrumentId string) *Quote {
	return (*qc)[instrumentId]
}
