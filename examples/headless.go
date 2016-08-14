package main

import (
	"log"

	"os"

	"github.com/apourchet/investment"
	"github.com/apourchet/investment/lib/ema"
)

type Trader struct {
	account *invt.Account
	in      chan *invt.Quote
}

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
		steps := ema5.Steps
		ema5Diff := (ema5.ComputeNext(q.Bid) - ema5.Value)
		ema30Diff := (ema30.ComputeNext(q.Bid) - ema30.Value)
		if long > 0 && q.Bid >= longAt+0.0020 {
			invt.TradeQuote(t.account, q, long, invt.SIDE_SELL)
			long = 0
			longAt = 0
		} else if long > 0 && q.Bid >= longAt+0.0010 && steps-longAtTime >= 10 {
			invt.TradeQuote(t.account, q, long, invt.SIDE_SELL)
			long = 0
			longAt = 0
		} else if long > 0 && (steps-longAtTime >= 60*5 || q.Bid <= longAt) {
			invt.TradeQuote(t.account, q, long, invt.SIDE_SELL)
			long = 0
			longAt = 0
		} else if short > 0 && q.Bid <= shortAt-0.0020 {
			invt.TradeQuote(t.account, q, short, invt.SIDE_BUY)
			short = 0
			shortAt = 0
		} else if short > 0 && q.Bid <= shortAt-0.0010 && steps-shortAtTime >= 10 {
			invt.TradeQuote(t.account, q, short, invt.SIDE_BUY)
			short = 0
			shortAt = 0
		} else if short > 0 && (steps-shortAtTime >= 60*5 || q.Bid >= shortAt) {
			invt.TradeQuote(t.account, q, short, invt.SIDE_BUY)
			short = 0
			shortAt = 0
		} else if ema5Diff > 0.0005 && ema30Diff < 0.0002 && short == 0 {
			invt.TradeQuote(t.account, q, 1000, invt.SIDE_BUY)
			if longAt == 0 {
				longAt = q.Price(invt.SIDE_BUY_STR)
			}
			long += 1000
			longAtTime = steps
		} else if ema5Diff < -0.0005 && ema30Diff > -0.0002 && long == 0 {
			invt.TradeQuote(t.account, q, 1000, invt.SIDE_SELL)
			if shortAt == 0 {
				shortAt = q.Price(invt.SIDE_SELL_STR)
			}
			short += 1000
			shortAtTime = steps
		}
		ema5.Step(q.Bid)
		ema30.Step(q.Bid)
	}
	return nil
}

func (t *Trader) OnEnd() {
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
		t.in <- q
	}
}

func main() {
	datafile := "examples/data/largish.csv"
	if len(os.Args) >= 2 {
		datafile = os.Args[1]
	}
	trader := NewTrader()
	simulator := invt.NewSimulator(invt.DATAFORMAT_CANDLE, datafile, 0)
	simulator.SimulateDataStream(trader)
}
