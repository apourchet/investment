package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/apourchet/investment"
	"github.com/apourchet/investment/lib/ema"
	pb "github.com/apourchet/investment/protos"
)

type Strat func(broker pb.BrokerClient, stream pb.Broker_StreamPricesClient)

func mine(broker pb.BrokerClient, stream pb.Broker_StreamPricesClient) {
	ema5 := ema.NewEma(ema.AlphaFromN(8))
	ema30 := ema.NewEma(ema.AlphaFromN(34))
	position := 0 // ema5 < ema30
	for {
		q, err := stream.Recv()
		if err == io.EOF || q == nil {
			if position == 1 {
				// Close position
				o := &pb.OrderCreationReq{}
				o.InstrumentId = "EURUSD"
				o.Type = invt.TYPE_MARKET
				o.Side = invt.StringOfSide(invt.SIDE_SELL)
				o.Units = 3000
				position = 1
				broker.CreateOrder(context.Background(), o)
			}
			return
		}
		ema5.Step(q.Bid)
		ema30.Step(q.Bid)
		if position == 0 && ema5.Value > ema30.Value {
			o := &pb.OrderCreationReq{}
			o.InstrumentId = "EURUSD"
			o.Type = invt.TYPE_MARKET
			o.Side = invt.StringOfSide(invt.SIDE_BUY)
			o.Units = 3000
			position = 1
			broker.CreateOrder(context.Background(), o)
		} else if position == 1 && ema5.Value < ema30.Value {
			o := &pb.OrderCreationReq{}
			o.InstrumentId = "EURUSD"
			o.Type = invt.TYPE_MARKET
			o.Side = invt.StringOfSide(invt.SIDE_SELL)
			o.Units = 3000
			position = 0
			broker.CreateOrder(context.Background(), o)
		}
	}

}

func startTrader(strat Strat) {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	exitOnError(err)
	defer conn.Close()

	broker := pb.NewBrokerClient(conn)
	req := &pb.StreamPricesReq{}
	req.InstrumentId = "EURUSD"
	stream, err := broker.StreamPrices(context.Background(), req)
	exitOnError(err)

	strat(broker, stream)
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
	}
}

func main() {
	datafile := "examples/data/largish.csv"
	if len(os.Args) >= 2 {
		datafile = os.Args[1]
	}

	broker := invt.NewDefaultBroker()
	go broker.Start()
	time.Sleep(time.Millisecond * 50)

	milliStep := 1
	go startTrader(mine)
	invt.SimulateDataStream(broker, datafile, milliStep)

	time.Sleep(time.Millisecond) // Let last changes kick in
	req := &pb.AccountInfoReq{}
	resp, _ := broker.GetAccountInfo(context.Background(), req)
	fmt.Println(resp.Info.Balance)
}
