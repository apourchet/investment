package invt

import "github.com/apourchet/investment/protos"

type Quote protos.Quote

type QuoteContext map[string]*Quote

func (qc *QuoteContext) Get(instrumentId string) *Quote {
	return (*qc)[instrumentId]
}

func (q *Quote) Proto() *protos.Quote {
	q1 := &protos.Quote{}
	q1.InstrumentId = q.InstrumentId
	q1.Ask = q.Ask
	q1.Bid = q.Bid
	q1.Status = q.Status
	q1.Time = q.Time
	return q1
}
