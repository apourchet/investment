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

var (
	totalTrades = 0
)

func quickOrder(units int32, side string) *pb.OrderCreationReq {
	o := &pb.OrderCreationReq{}
	o.InstrumentId = "EURUSD"
	o.Type = invt.TYPE_MARKET
	o.Side = side
	o.Units = units
	return o
}

func mine(broker pb.BrokerClient, stream pb.Broker_StreamPricesClient) {
	ema5 := ema.NewEma(ema.AlphaFromN(8))
	ema30 := ema.NewEma(ema.AlphaFromN(50))
	position := 0 // ema5 < ema30
	for {
		q, err := stream.Recv()
		if err == io.EOF || q == nil {
			if position == 1 {
				// Close position
				o := quickOrder(3000, invt.StringOfSide(invt.SIDE_SELL))
				broker.CreateOrder(context.Background(), o)
			}
			return
		}

		if ema5.Steps%10000 == 0 {
			req := &pb.AccountInfoReq{}
			resp, _ := broker.GetAccountInfo(context.Background(), req)
			fmt.Println(resp.Info.Balance, totalTrades)
		}

		ema5.Step(q.Bid)
		ema30.Step(q.Bid)

		if position == 0 && ema5.Value > ema30.Value {
			o := quickOrder(3000, invt.StringOfSide(invt.SIDE_BUY))
			position = 1
			broker.CreateOrder(context.Background(), o)
			totalTrades++
		} else if position == 1 && ema5.Value < ema30.Value {
			o := quickOrder(3000, invt.StringOfSide(invt.SIDE_SELL))
			position = 0
			broker.CreateOrder(context.Background(), o)
			totalTrades++
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

	time.Sleep(time.Millisecond * 500) // Let last changes kick in
	req := &pb.AccountInfoReq{}
	resp, _ := broker.GetAccountInfo(context.Background(), req)
	fmt.Println(resp.Info.Balance)
}
