package main

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/net/context"

	"time"

	"github.com/apourchet/investment"
	"github.com/apourchet/investment/lib/ema"
	tl "github.com/apourchet/investment/lib/tradelogger"
	pb "github.com/apourchet/investment/protos"
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
	ema5 := ema.NewEma(ema.AlphaFromN(30))
	lma := ema.NewEma(400)

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

		if ema5.Steps%2000 == 0 {
			req := &pb.AccountInfoReq{}
			resp, _ := broker.GetAccountInfo(context.Background(), req)
			fmt.Println(resp.Info.MarginAvail)
		}

		ema5.Step(q.Bid)
		lma.Step(q.Bid)

		if position == 0 && ema5.Value > lma.Value {
			if q.Ask < ema5.Value+0.0001 {
				o := quickOrder(3000, invt.StringOfSide(invt.SIDE_BUY))
				position = 1
				broker.CreateOrder(context.Background(), o)
			}
		} else if position == 1 && ema5.Value < lma.Value {
			if q.Bid > ema5.Value-0.0001 {
				o := quickOrder(3000, invt.StringOfSide(invt.SIDE_SELL))
				position = 0
				broker.CreateOrder(context.Background(), o)
			}
		}

	}
}

func main() {
	datafile := "examples/data/largish.csv"
	if len(os.Args) >= 2 {
		datafile = os.Args[1]
	}

	go tl.StartServer(1026, "logs/")
	time.Sleep(time.Millisecond * 50)

	logger := tl.NewLoggerClient("http://localhost:1026")
	invt.AddLogger(logger)

	broker := invt.NewDefaultBroker()
	invt.SimulateTradingScenario(broker, mine, datafile)
}
