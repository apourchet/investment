package main

import (
	"log"

	"os"

	"github.com/apourchet/investment"
	"github.com/apourchet/investment/lib/ema"
	"github.com/apourchet/investment/lib/influx-session"
)

type Trader struct {
	account *invt.Account
	in      chan *invt.Quote
}

type MyPoint struct {
	Bid   float64
	Ema5  float64
	Ema30 float64
}

var (
	db *ix_session.Session
)

func NewTrader() *Trader {
	return &Trader{invt.NewAccount(10000), make(chan *invt.Quote)}
}

func (t *Trader) Start() error {
	log.Println("Trader starting")
	ema5 := ema.NewEma(ema.AlphaFromN(5))
	ema30 := ema.NewEma(ema.AlphaFromN(30))

	longAtTime := 0
	shortAtTime := 0
	longAt := 0.
	shortAt := 0.
	short := int32(0)
	long := int32(0)
	for q := range t.in {
		if q == nil {
			break
		}
		pt := MyPoint{}
		pt.Bid = q.Bid

		steps := ema5.Steps
		ema5Diff := (ema5.ComputeNext(q.Bid) - ema5.Value)
		ema30Diff := (ema30.ComputeNext(q.Bid) - ema30.Value)
		sold := false
		bought := false
		if long > 0 && q.Bid >= longAt+0.0020 {
			invt.TradeQuote(t.account, q, long, invt.SIDE_SELL)
			sold = true
			long = 0
			longAt = 0
		} else if long > 0 && q.Bid >= longAt+0.0010 && steps-longAtTime >= 10 {
			invt.TradeQuote(t.account, q, long, invt.SIDE_SELL)
			sold = true
			long = 0
			longAt = 0
		} else if long > 0 && (steps-longAtTime >= 60*5 || q.Bid <= longAt) {
			invt.TradeQuote(t.account, q, long, invt.SIDE_SELL)
			sold = true
			long = 0
			longAt = 0
		} else if short > 0 && q.Bid <= shortAt-0.0020 {
			invt.TradeQuote(t.account, q, short, invt.SIDE_BUY)
			bought = true
			short = 0
			shortAt = 0
		} else if short > 0 && q.Bid <= shortAt-0.0010 && steps-shortAtTime >= 10 {
			invt.TradeQuote(t.account, q, short, invt.SIDE_BUY)
			bought = true
			short = 0
			shortAt = 0
		} else if short > 0 && (steps-shortAtTime >= 60*5 || q.Bid >= shortAt) {
			invt.TradeQuote(t.account, q, short, invt.SIDE_BUY)
			bought = true
			short = 0
			shortAt = 0
		} else if ema5Diff > 0.0005 && ema30Diff < 0.0002 && short == 0 {
			invt.TradeQuote(t.account, q, 1000, invt.SIDE_BUY)
			bought = true
			if longAt == 0 {
				longAt = q.Price(invt.SIDE_BUY_STR)
			}
			long += 1000
			longAtTime = steps
		} else if ema5Diff < -0.0005 && ema30Diff > -0.0002 && long == 0 {
			invt.TradeQuote(t.account, q, 1000, invt.SIDE_SELL)
			sold = true
			if shortAt == 0 {
				shortAt = q.Price(invt.SIDE_SELL_STR)
			}
			short += 1000
			shortAtTime = steps
		}
		if sold {
			db.Write("moment", struct{ OrderSell float64 }{q.Bid}, q.Timestamp)
		} else if bought {
			db.Write("moment", struct{ OrderBuy float64 }{q.Ask}, q.Timestamp)
		}

		ema5.Step(q.Bid)
		ema30.Step(q.Bid)
		pt.Ema5 = ema5.Value
		pt.Ema30 = ema30.Value
		db.Write("moment", pt, q.Timestamp)
	}
	return nil
}

func (t *Trader) OnEnd() {
	db.Flush()
	log.Printf("%+v\n", t.account.Stats)
	log.Printf("%+v\n", t.account)
	t.in <- nil
}

func (t *Trader) OnData(record []string, format invt.DataFormat) {
	if format == invt.DATAFORMAT_QUOTE {
		q := invt.ParseQuoteFromRecord("EURUSD", record)
		t.in <- q
	} else if format == invt.DATAFORMAT_CANDLE {
		c := invt.ParseCandleFromRecord("EURUSD", record)
		q := &invt.Quote{}
		q.Bid = c.Close
		q.Ask = c.Close + 0.00025
		q.InstrumentId = c.InstrumentId
		q.Timestamp = c.Timestamp
		t.in <- q
	}
}

func main() {
	db = ix_session.NewSession(ix_session.DEFAULT_ADDRESS, "investment", "password", "testdb")
	datafile := "examples/data/largish.csv"
	if len(os.Args) >= 2 {
		datafile = os.Args[1]
	}
	trader := NewTrader()
	simulator := invt.NewSimulator(invt.DATAFORMAT_CANDLE, datafile, 0)
	simulator.SimulateDataStream(trader)
}
