package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/apourchet/investment"
	pb "github.com/apourchet/investment/protos"
)

type Strat func(broker pb.BrokerClient, stream pb.Broker_StreamQuotesClient)

func DownUpCSVStrategy(broker pb.BrokerClient, stream pb.Broker_StreamQuotesClient) {
	for {
		q, err := stream.Recv()
		if err == io.EOF || q == nil {
			return
		}
		if int(q.Bid) == 1 {
			// Buy a bunch
			o := &pb.OrderCreation{}
			o.Instrument = pb.InstrumentID_EURUSD
			o.Type = pb.OrderType_MARKET
			o.Side = pb.OrderSide_BUY
			o.Units = 100
			broker.CreateOrder(context.Background(), o)
		}
		if int(q.Bid) == 5 {
			// Sell a bunch
			o := &pb.OrderCreation{}
			o.Instrument = pb.InstrumentID_EURUSD
			o.Type = pb.OrderType_MARKET
			o.Side = pb.OrderSide_SELL
			o.Units = 100
			broker.CreateOrder(context.Background(), o)
		}
	}
}

func startTrader(strat Strat) {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	exitOnError(err)
	defer conn.Close()

	broker := pb.NewBrokerClient(conn)
	iid := &pb.InstrumentID{pb.InstrumentID_EURUSD}
	stream, err := broker.StreamQuotes(context.Background(), iid)
	exitOnError(err)

	strat(broker, stream)
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
	}
}

func main() {
	broker := invt.NewDefaultBroker()
	go broker.Start()
	time.Sleep(time.Millisecond * 50)

	milliStep := 200
	go startTrader(DownUpCSVStrategy)
	invt.SimulateDataStream(broker, "tests/testdata/updown.csv", milliStep)

	go startTrader(DownUpCSVStrategy)
	invt.SimulateDataStream(broker, "tests/testdata/downup.csv", milliStep)
}
