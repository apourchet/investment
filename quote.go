package invt

import (
	"strconv"

	"github.com/apourchet/investment/lib/utils"
	"github.com/apourchet/investment/protos"
)

type Quote protos.Quote

type Candle protos.Candle

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

func (c *Candle) Proto() *protos.Candle {
	c1 := &protos.Candle{}
	c1.InstrumentId = c.InstrumentId
	c1.Close = c.Close
	c1.High = c.High
	c1.Low = c.Low
	c1.Open = c.Open
	c1.Status = c.Status
	c1.Time = c.Time
	return c1
}

func parseCandle(record []string) *Candle {
	c := &Candle{}
	c.InstrumentId = "EURUSD"
	v1, err1 := strconv.ParseFloat(record[2], 64)
	v2, err2 := strconv.ParseFloat(record[3], 64)
	v3, err3 := strconv.ParseFloat(record[4], 64)
	v4, err4 := strconv.ParseFloat(record[5], 64)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return nil
	}

	c.Time, _ = utils.ParseDateString(record)
	c.Open = v1
	c.High = v2
	c.Low = v3
	c.Close = v4
	return c
}

func parseQuote(record []string) *Quote {
	q := &Quote{}
	q.InstrumentId = "EURUSD"
	v, err := strconv.ParseFloat(record[2], 64)
	q.Bid = v
	if err != nil {
		return nil
	}

	q.Ask, err = strconv.ParseFloat(record[2], 64)
	q.Ask += 0.00025 // Adjust for arbitrary spread
	if err != nil {
		return nil
	}
	q.Time, _ = utils.ParseDateString(record)
	return q
}
