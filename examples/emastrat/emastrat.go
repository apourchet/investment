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
	log.Println("Trader started")
	broker := def.GetClient()
	stream := getStream(broker)

	ema5 := ema.NewEma(ema.AlphaFromN(5))
	ema30 := ema.NewEma(ema.AlphaFromN(30))
	long := int32(0)
	longAt := 0.
	longAtTime := 0
	short := int32(0)
	shortAt := 0.
	shortAtTime := 0
	for {
		c1, err := stream.Recv()
		if err == io.EOF || c1 == nil {
			log.Println("Candle stream has ended.")
			return
		}
		c := invt.CandleFromProto(c1)
		steps := ema5.Steps
		ema5Diff := (ema5.ComputeNext(c.Close) - ema5.Compute())
		ema30Diff := (ema30.ComputeNext(c.Close) - ema30.Compute())

		if long > 0 && c.Close >= longAt+0.0020 {
			o := quickOrder(long, invt.SIDE_SELL_STR)
			resp, _ := broker.CreateOrder(context.Background(), o)
			log.Println("Closed at", resp.Price, (resp.Price - longAt), resp.Time)
			long = 0
			longAt = 0
			resp2, _ := broker.GetAccountInfo(context.Background(), &pb.AccountInfoReq{})
			log.Println("PL: ", resp2.Info)
		} else if long > 0 && c.Close >= longAt+0.0010 && steps-longAtTime >= 10 {
			o := quickOrder(long, invt.SIDE_SELL_STR)
			resp, _ := broker.CreateOrder(context.Background(), o)
			log.Println("Closed (-) at", resp.Price, (resp.Price - longAt), resp.Time)
			long = 0
			longAt = 0
			resp2, _ := broker.GetAccountInfo(context.Background(), &pb.AccountInfoReq{})
			log.Println("PL: ", resp2.Info)
		} else if long > 0 && (steps-longAtTime >= 60*5 || c.Close <= longAt) {
			o := quickOrder(long, invt.SIDE_SELL_STR)
			resp, _ := broker.CreateOrder(context.Background(), o)
			log.Println("Closed (--) at", resp.Price, (resp.Price - longAt), resp.Time)
			long = 0
			longAt = 0
			resp2, _ := broker.GetAccountInfo(context.Background(), &pb.AccountInfoReq{})
			log.Println("PL: ", resp2.Info)
		} else if short > 0 && c.Close <= shortAt-0.0020 {
			o := quickOrder(short, invt.SIDE_BUY_STR)
			resp, _ := broker.CreateOrder(context.Background(), o)
			log.Println("Closed at", resp.Price, (shortAt - resp.Price), resp.Time)
			short = 0
			shortAt = 0
			resp2, _ := broker.GetAccountInfo(context.Background(), &pb.AccountInfoReq{})
			log.Println("PL: ", resp2.Info)
		} else if short > 0 && c.Close <= shortAt-0.0010 && steps-shortAtTime >= 10 {
			o := quickOrder(short, invt.SIDE_BUY_STR)
			resp, _ := broker.CreateOrder(context.Background(), o)
			log.Println("Closed (-) at", resp.Price, -(resp.Price - shortAt), resp.Time)
			short = 0
			shortAt = 0
			resp2, _ := broker.GetAccountInfo(context.Background(), &pb.AccountInfoReq{})
			log.Println("PL: ", resp2.Info)
		} else if short > 0 && (steps-shortAtTime >= 60*5 || c.Close >= shortAt) {
			o := quickOrder(short, invt.SIDE_BUY_STR)
			resp, _ := broker.CreateOrder(context.Background(), o)
			log.Println("Closed (--) at", resp.Price, -(resp.Price - shortAt), resp.Time)
			short = 0
			shortAt = 0
			resp2, _ := broker.GetAccountInfo(context.Background(), &pb.AccountInfoReq{})
			log.Println("PL: ", resp2.Info)
		} else if ema5Diff > 0.0005 && ema30Diff < 0.0002 && short == 0 {
			o := quickOrder(1000, invt.SIDE_BUY_STR)
			resp, _ := broker.CreateOrder(context.Background(), o)
			log.Println("Went long at", resp.Price)
			if longAt == 0 {
				longAt = resp.Price
			}
			long += 1000
			longAtTime = steps
		} else if ema5Diff < -0.0005 && ema30Diff > -0.0002 && long == 0 {
			o := quickOrder(1000, invt.SIDE_SELL_STR)
			resp, _ := broker.CreateOrder(context.Background(), o)
			log.Println("Went short at", resp.Price)
			if shortAt == 0 {
				shortAt = resp.Price
			}
			short += 1000
			shortAtTime = steps
		}

		if ema5.Steps%500 == 0 {
			log.Println("Still going...", ema5.Steps)
		}

		ema5.Step(c.Close)
		ema30.Step(c.Close)
	}
}

func main() {
	log.SetOutput(os.Stdout)
	session = ix_session.NewSession(ix_session.DEFAULT_ADDRESS, "investment", "password", "testdb")
	datafile := "examples/data/largish.csv"
	if len(os.Args) >= 2 {
		datafile = os.Args[1]
	}

	simulator := invt.NewSimulator(invt.DATAFORMAT_CANDLE, datafile, 0)
	broker := invt.NewDefaultBroker(1027)
	go simulator.SimulateDataStream(broker)

	time.Sleep(time.Millisecond * 50)
	mine(broker)
}
