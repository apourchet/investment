package main

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/net/context"

	"time"

	"log"

	"github.com/apourchet/investment"
	"github.com/apourchet/investment/lib/ema"
	"github.com/apourchet/investment/lib/influx-session"
	pb "github.com/apourchet/investment/protos"
)

var (
	session *ix_session.Session
)

func quickOrder(units int32, side string) *pb.OrderCreationReq {
	o := &pb.OrderCreationReq{}
	o.InstrumentId = "EURUSD"
	o.Type = invt.TYPE_MARKET
	o.Side = side
	o.Units = units
	return o
}

func getStream(broker pb.BrokerClient) pb.Broker_StreamCandleClient {
	req := &pb.StreamCandleReq{&pb.AuthToken{}, "EURUSD"}
	stream, err := broker.StreamCandle(context.Background(), req)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}
	return stream
}

func mine(def *invt.DefaultBroker) {
	fmt.Println("Trader started")
	broker := def.GetClient()
	stream := getStream(broker)

	ema5 := ema.NewEma(ema.AlphaFromN(5))
	ema30 := ema.NewEma(ema.AlphaFromN(30))
	for {
		c1, err := stream.Recv()
		if err == io.EOF || c1 == nil {
			fmt.Println("Candle stream has ended.")
			return
		}
		c := invt.CandleFromProto(c1)
		if ema5.Value < ema30.Value && ema5.ComputeNext(c.Close) > ema30.ComputeNext(c.Close) {
			o := quickOrder(10, invt.SIDE_BUY_STR)
			broker.CreateOrder(context.Background(), o)
			err = session.Write("order.buy", o, c.Timestamp)
		} else if ema5.Value > ema30.Value && ema5.ComputeNext(c.Close) < ema30.ComputeNext(c.Close) {
			o := quickOrder(10, invt.SIDE_SELL_STR)
			broker.CreateOrder(context.Background(), o)
			err = session.Write("order.sell", o, c.Timestamp)
		}
		if err != nil {
			log.Fatal("Error writing to influxdb:", err)
		}
		ema5.Step(c.Close)
		ema30.Step(c.Close)
		session.Write("candle", c, c.Timestamp)
		session.Write("ema.5", ema5, c.Timestamp)
		session.Write("ema.30", ema30, c.Timestamp)
	}
}

func main() {
	session = ix_session.NewSession(ix_session.DEFAULT_ADDRESS, "investment", "password", "testdb")
	datafile := "examples/data/largish.csv"
	if len(os.Args) >= 2 {
		datafile = os.Args[1]
	}

	simulator := invt.NewSimulator(invt.DATAFORMAT_CANDLE, datafile, 10)
	broker := invt.NewDefaultBroker(1027)
	go simulator.SimulateDataStream(broker)

	time.Sleep(time.Millisecond * 50)
	mine(broker)
}
