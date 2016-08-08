package invt

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	pb "github.com/apourchet/investment/protos"
)

type Simulatable interface {
	pb.BrokerServer
	OnQuote(*pb.Quote)
	OnEnd()
}

type AccountSimulator struct {
	Account *Account
	orders  []*pb.Order
}

func NewAccountSimulator(balance float64) *AccountSimulator {
	return &AccountSimulator{CreateNewAccount(balance), make([]*pb.Order, 0)}
}

func (as *AccountSimulator) Buy(instrumentID pb.InstrumentID_ID, units int32, price float64) *AccountSimulator {
	o := &pb.Order{}
	o.Instrument = instrumentID
	o.Units = units
	o.Price = price
	o.Side = pb.OrderSide_BUY
	o.Type = pb.OrderType_MARKET
	as.orders = append(as.orders, o)
	return as
}

func (as *AccountSimulator) BuyNow(instrumentID pb.InstrumentID_ID, units int32, price float64) *AccountSimulator {
	o := &pb.Order{}
	o.Instrument = instrumentID
	o.Units = units
	o.Price = price
	o.Side = pb.OrderSide_BUY
	o.Type = pb.OrderType_MARKET
	as.orders = append(as.orders, o)
	as.Account.ProcessOrder(o)
	return as
}

func (as *AccountSimulator) Sell(instrumentID pb.InstrumentID_ID, units int32, price float64) *AccountSimulator {
	o := &pb.Order{}
	o.Instrument = instrumentID
	o.Units = units
	o.Price = price
	o.Side = pb.OrderSide_SELL
	o.Type = pb.OrderType_MARKET
	as.orders = append(as.orders, o)
	return as
}

func (as *AccountSimulator) SellNow(instrumentID pb.InstrumentID_ID, units int32, price float64) *AccountSimulator {
	o := &pb.Order{}
	o.Instrument = instrumentID
	o.Units = units
	o.Price = price
	o.Side = pb.OrderSide_SELL
	o.Type = pb.OrderType_MARKET
	as.Account.ProcessOrder(o)
	return as
}

func (as *AccountSimulator) Simulate() {
	for _, o := range as.orders {
		as.Account.ProcessOrder(o)
	}
}

func SimulateDataStream(b Simulatable, datafile string, milliStep int) {
	in, err := os.Open(datafile)
	if err != nil {
		fmt.Println("Could not open data file: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("Simulating: " + datafile)
	reader := csv.NewReader(in)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			fmt.Println("Simulation Ended")
			b.OnEnd()
			break
		}
		q := &pb.Quote{}
		q.Id = pb.InstrumentID_EURUSD
		q.Bid, err = strconv.ParseFloat(record[2], 64)
		q.Ask, err = strconv.ParseFloat(record[4], 64)
		// TODO
		// q.Time = date.ParseDate(record[0])

		b.OnQuote(q)
		time.Sleep(time.Millisecond * time.Duration(milliStep))
	}
}
